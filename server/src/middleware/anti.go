// middleware/anti.go
package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/bits-and-blooms/bloom/v3"
	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

const (
	BloomCapacity = 1000000 // 布隆过滤器容量
	BloomFPRate   = 0.01    // 误判率 1%
	RateLimit     = 10      // 每秒 1 次
	Burst         = 20      // 突发 5 次
	DailyLimit    = 1000    // 每个 IP 24 小时内最多 1000 次

	RequestTokenSecret = "secret-key-rao-bang-lin!" // ← 必须修改！
	RequestTokenWindow = 60                         // 秒，前后 60 秒有效（共 2 分钟窗口）
)

// IPLimiter 每个 IP 的限流状态
type IPLimiter struct {
	limiter     *rate.Limiter
	dailyCount  int
	windowStart time.Time // 当前 24 小时窗口起始时间
	lastAccess  time.Time // 最后访问时间（用于清理）
}

var (
	limiterMap = make(map[string]*IPLimiter)
	limiterMu  sync.RWMutex

	bloomFilter *bloom.BloomFilter
	bloomMu     sync.RWMutex

	once sync.Once
)

// 获取或创建限流器（线程安全）
func getOrCreateLimiter(ip string) *IPLimiter {
	limiterMu.RLock()
	il, exists := limiterMap[ip]
	limiterMu.RUnlock()

	if !exists {
		limiterMu.Lock()
		defer limiterMu.Unlock()
		// 双重检查
		if il, exists = limiterMap[ip]; !exists {
			now := time.Now()
			il = &IPLimiter{
				limiter:     rate.NewLimiter(rate.Limit(RateLimit), Burst),
				dailyCount:  0,
				windowStart: now,
				lastAccess:  now,
			}
			limiterMap[ip] = il
		}
	}
	return il
}

// 检查 24 小时滑动窗口限制
func dailyLimit(ip string) bool {
	il := getOrCreateLimiter(ip)

	limiterMu.Lock()
	defer limiterMu.Unlock()

	now := time.Now()
	il.lastAccess = now

	// 如果超过 24 小时，重置窗口
	if now.Sub(il.windowStart) >= 24*time.Hour {
		il.dailyCount = 0
		il.windowStart = now
	}

	if il.dailyCount >= DailyLimit {
		return false
	}

	il.dailyCount++
	return true
}

// 初始化布隆过滤器
func initBloom() {
	bloomFilter = bloom.NewWithEstimates(BloomCapacity, BloomFPRate)
	fmt.Println("✅ Bloom filter initialized (in-memory)")
}

// 后台定期清理不活跃 IP（每小时一次）
func cleanup() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for range ticker.C {
			limiterMu.Lock()
			now := time.Now()
			for ip, il := range limiterMap {
				if now.Sub(il.lastAccess) > time.Hour {
					delete(limiterMap, ip)
				}
			}
			limiterMu.Unlock()
		}
	}()
}

// 将 IP 加入布隆过滤器（封禁）
func addBlockIP(ip string) {
	bloomMu.Lock()
	defer bloomMu.Unlock()
	if bloomFilter != nil {
		bloomFilter.AddString(ip)
	}
}

// 判断 IP 是否被封禁
func isBlockIP(ip string) bool {
	bloomMu.RLock()
	defer bloomMu.RUnlock()
	if bloomFilter == nil {
		return false
	}
	return bloomFilter.TestString(ip)
}

// verifyToken 仅校验时间窗口 + 签名，无状态、无缓存
func verifyToken(tokenStr string) bool {
	if tokenStr == "" {
		return false
	}

	// Base64URL 解码（兼容无 padding）
	data, err := base64.RawURLEncoding.DecodeString(tokenStr)
	if err != nil || len(data) < 40 { // 8 + 32 = 40 字节
		return false
	}

	tsBytes := data[:8]
	hmacSig := data[8:]

	// 解析时间戳（大端序）
	var timestamp int64
	for i := 0; i < 8; i++ {
		timestamp = (timestamp << 8) | int64(tsBytes[i])
	}

	// 获取当前时间（秒）
	now := time.Now().Unix()

	// 检查是否在 [now - 60, now + 60] 范围内
	if timestamp < now-int64(RequestTokenWindow) || timestamp > now+int64(RequestTokenWindow) {
		return false
	}

	// 验证 HMAC 签名
	mac := hmac.New(sha256.New, []byte(RequestTokenSecret))
	mac.Write(tsBytes)
	expectedSig := mac.Sum(nil)

	return hmac.Equal(hmacSig, expectedSig)
}

// AntiMiddleware Gin 中间件
func AntiMiddleware() gin.HandlerFunc {
	once.Do(func() {
		initBloom()
		cleanup()
	})

	return func(c *gin.Context) {
		// === 新增：第一步校验 X-Request-Token ===
		token := c.GetHeader("X-Request-Token")
		if !verifyToken(token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Access denied",
				"data":    nil,
			})
			return
		}

		ip := c.ClientIP()

		// 1. 24 小时请求限制（滑动窗口）
		if !dailyLimit(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Service unavailable",
				"data":    nil,
			})
			return
		}

		// 2. 布隆过滤器拦截已知恶意 IP
		if isBlockIP(ip) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    http.StatusForbidden,
				"message": "Access denied",
				"data":    nil,
			})
			return
		}

		// 3. User-Agent 检查（封禁常见爬虫关键词）
		ua := strings.ToLower(c.GetHeader("User-Agent"))
		badKeywords := []string{"python", "scrapy", "curl", "httpclient", "java", "bot", "spider", "crawl"}
		for _, kw := range badKeywords {
			if strings.Contains(ua, kw) {
				addBlockIP(ip) // 自动封禁该 IP
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": "Access denied",
					"data":    nil,
				})
				return
			}
		}

		// 4. Referer 白名单校验（允许空 Referer）
		referer := c.GetHeader("Referer")
		if referer != "" {
			valid := strings.Contains(referer, "localhost") ||
				strings.Contains(referer, "ma3ok.com") ||
				strings.Contains(referer, "mageg.cn")
			if !valid {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"code":    http.StatusForbidden,
					"message": "Access denied",
					"data":    nil,
				})
				return
			}
		}

		// 5. 秒级速率限制
		il := getOrCreateLimiter(ip)
		if !il.limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":    http.StatusTooManyRequests,
				"message": "Service unavailable",
				"data":    nil,
			})
			return
		}

		c.Next()
	}
}
