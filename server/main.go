package main

import (
	"github.com/sirupsen/logrus"

	"github.com/mageg-x/novel/src/log"
	"github.com/mageg-x/novel/src/router"
	"github.com/mageg-x/novel/src/service"
)

var (
	logger  = log.GetLogger("novel")
	verbose = 4
)

func init() {
	// 初始化日志系统
	log.Init(&log.Config{
		LogDir:     "./logs",
		MaxSize:    10,
		MaxBackups: 7,
		MaxAge:     30,
		Compress:   true,
	})

	setLogLevel()
	logger.Info("日志系统初始化完成")
}

// setLogLevel 根据 verbose 级别设置日志级别
func setLogLevel() {
	var level logrus.Level

	switch verbose {
	case 0:
		level = logrus.ErrorLevel // 默认错误级别
	case 1:
		level = logrus.WarnLevel // -v: 警告级别
	case 2:
		level = logrus.InfoLevel // -vv: 信息级别
	case 3:
		level = logrus.DebugLevel // -vvv: 调试级别
	default:
		level = logrus.TraceLevel // -vvvv 或更多: 跟踪级别
	}
	// fmt.Println("Log level:", level.String())
	logger.SetLevel(level)

	logger.Debugf("verbose level: %d, log level: %s", verbose, level.String())
}

func main() {
	logger.Info("开始启动小说阅读系统...")

	// 初始化数据库
	if err := service.InitDB(); err != nil {
		logger.Fatalf("数据库初始化失败: %v", err)
		return
	}
	logger.Info("数据库初始化完成")

	// 设置路由
	r := router.SetupRouter()
	logger.Info("路由配置完成")

	// 启动服务器
	serverAddr := ":8000"
	logger.Infof("服务器启动成功，监听地址: http://localhost%s", serverAddr)

	if err := r.Run(serverAddr); err != nil {
		logger.Fatalf("服务器启动失败: %v", err)
	}
}
