<template>
    <!-- 小屏幕 -->
    <div class="block md:hidden">
        <!-- 顶部导航 -->
        <div class="sticky top-0 z-50 bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-lg">
            <div class="flex items-center justify-between h-14 px-4">
                <button @click="goBack" class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <span class="text-lg font-semibold">排行榜列表</span>
                <button @click="goHome" class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fa-solid fa-house"></i>
                </button>
            </div>
        </div>
        
        <!-- 内容区域 -->
        <div class="bg-gradient-to-br from-gray-50 to-gray-100 min-h-screen pt-4 pb-20">
            <!-- 排行榜类型切换 -->
            <div class="px-4 mb-4">
                <div class="flex space-x-2 bg-white rounded-xl shadow-sm p-1">
                    <button 
                        v-for="type in ['click', 'new', 'update', 'comment']" 
                        :key="type"
                        @click="currentRankType = type"
                        class="flex-1 py-2 px-3 rounded-lg transition-all duration-300 text-sm font-medium"
                        :class="currentRankType === type 
                            ? 'bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-md transform scale-105' 
                            : 'text-gray-600 hover:bg-gray-100'"
                    >
                        {{ type === 'click' ? '点击榜' : 
                           type === 'new' ? '新书榜' : 
                           type === 'update' ? '更新榜' : '评论榜' }}
                    </button>
                </div>
            </div>

            <!-- 排行榜列表 -->
            <div class="px-4">
                <div 
                    v-for="(book, index) in books" 
                    :key="index"
                    class="bg-white rounded-2xl shadow-md p-4 mb-3 transition-all duration-300 hover:shadow-lg hover:-translate-y-1 cursor-pointer"
                    @click="navigateToBook(book.bookId)"
                >
                    <!-- 排名与书籍信息 -->
                    <div class="flex items-center">
                        <!-- 排名 -->
                        <div class="mr-3">
                            <span 
                                :class="{
                                    'bg-red-500': index === 0,
                                    'bg-orange-500': index === 1,
                                    'bg-yellow-500': index === 2,
                                    'bg-blue-500': index > 2 && index < 10,
                                    'bg-gray-300': index >= 10
                                }" 
                                class="flex items-center justify-center w-9 h-9 text-white rounded-full text-base font-bold shadow-sm"
                            >
                                {{ index + 1 }}
                            </span>
                        </div>

                        <!-- 封面图片 -->
                        <div class="mr-3">
                            <img 
                                :src="book.cover" 
                                alt="{{ book.title }}" 
                                class="w-16 h-22 rounded-lg object-cover shadow-md transition-transform duration-300 hover:scale-105"
                            />
                        </div>

                        <!-- 书籍信息 -->
                        <div class="flex-1">
                            <div class="flex items-center justify-between mb-1">
                                <div class="flex items-center">
                                    <span class="text-xs px-2 py-0.5 bg-emerald-100 text-emerald-700 rounded-full font-medium mr-2">
                                        [{{ book.category }}]
                                    </span>
                                    <span class="text-gray-500 text-xs font-medium">{{ book.author }}</span>
                                </div>
                                <span class="text-sm text-gray-600 font-medium ml-2 whitespace-nowrap">
                                    {{ book.wordCount }}万
                                </span>
                            </div>
                            
                            <h3 class="text-base font-semibold text-gray-800 mb-1 line-clamp-1 hover:text-emerald-600 transition-colors">
                                {{ book.title }}
                            </h3>
                            
                            <p class="text-xs text-gray-500 line-clamp-1">{{ book.latestChapter }}</p>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 底部导航 -->
        <ToolBar :activeTab="activeTab" @tab-change="handleTabChange" />
    </div>
    <!-- 大屏 -->
    <div class="hidden md:block container md:w-5xl  mx-auto ">
        <div class="hidden md:block">
            <Header />
        </div>   
        <main class="max-w-7xl mx-auto py-8">
            <div class="flex gap-6">
                <!-- 左侧排行榜内容 -->
                <div class="flex-1 bg-white rounded-lg shadow">
                    <h1 class="text-2xl font-bold text-gray-800 px-6 py-4 ">{{ currentRankTitle }}</h1>

                    <!-- 表头 -->
                    <div class="grid grid-cols-13 gap-4 px-6 py-3 bg-gray-100 text-gray-600 text-sm">
                        <div class="col-span-1">排名</div>
                        <div class="col-span-2">类别</div>
                        <div class="col-span-3">书名</div>
                        <div class="col-span-4">最新章节</div>
                        <div class="col-span-2">作者</div>
                        <div class="col-span-1 text-right">字数</div>
                    </div>

                    <!-- 排行榜列表 -->
                    <div class="divide-y divide-gray-100">
                        <div v-for="(book, index) in books" :key="index"
                            class="grid grid-cols-13 gap-4 px-6 py-3 hover:bg-gray-50 items-center transition-colors duration-150">
                            <div class="col-span-1 overflow-hidden">
                                <span :class="{ 'bg-red-500': index === 0, 'bg-orange-500': index === 1, 'bg-yellow-500': index === 2, 'bg-blue-500': index > 2 && index < 10, 'bg-gray-300': index >= 10 }"
                                    class="inline-flex items-center justify-center w-7 h-7 text-white rounded-full text-sm font-bold">{{ index + 1 }}</span>
                            </div>
                            <div class="col-span-2 overflow-hidden text-gray-500 text-sm whitespace-nowrap pr-2">{{ `[${book.category}]` }}</div>
                            <div class="col-span-3 whitespace-nowrap overflow-hidden text-ellipsis"><router-link :to="'/book/' + book.bookId" class="text-gray-600 text-sm hover:text-blue-800 hover:underline transition-colors cursor-pointer">{{ book.title}}</router-link></div>
                            <div class="col-span-4 whitespace-nowrap overflow-hidden text-ellipsis text-gray-600 text-sm">{{ book.latestChapter }}</div>
                            <div class="col-span-2 whitespace-nowrap overflow-hidden text-ellipsis text-gray-600 text-sm">{{ book.author }}</div>
                            <div class="col-span-1 overflow-hidden text-gray-500 text-sm text-right whitespace-nowrap">{{ book.wordCount }}万</div>
                        </div>
                    </div>
                </div>

                <!-- 右侧标签栏 -->
                <div class="w-52">
                    <div class="bg-white rounded-xl shadow-md hover:shadow-lg transition-shadow duration-200 overflow-hidden">
                        <h2 class="text-base font-semibold text-gray-800 px-4 py-3 border-b border-gray-100 bg-gray-50">排行榜</h2>
                        <ul class="py-1">
                            <li><a href="#" @click.prevent="currentRankType = 'click'"
                                    class="block px-4 py-2.5 transition-all duration-150" :class="{'text-emerald-600 bg-emerald-50 border-l-4 border-emerald-600': currentRankType === 'click', 'text-gray-700 hover:bg-gray-100 hover:text-emerald-600': currentRankType !== 'click'}">点击榜</a>
                            </li>
                            <li><a href="#" @click.prevent="currentRankType = 'new'" class="block px-4 py-2.5 transition-all duration-150" :class="{'text-emerald-600 bg-emerald-50 border-l-4 border-emerald-600': currentRankType === 'new', 'text-gray-700 hover:bg-gray-100 hover:text-emerald-600': currentRankType !== 'new'}">新书榜</a></li>
                            <li><a href="#" @click.prevent="currentRankType = 'update'" class="block px-4 py-2.5 transition-all duration-150" :class="{'text-emerald-600 bg-emerald-50 border-l-4 border-emerald-600': currentRankType === 'update', 'text-gray-700 hover:bg-gray-100 hover:text-emerald-600': currentRankType !== 'update'}">更新榜</a></li>
                            <li><a href="#" @click.prevent="currentRankType = 'comment'" class="block px-4 py-2.5 transition-all duration-150" :class="{'text-emerald-600 bg-emerald-50 border-l-4 border-emerald-600': currentRankType === 'comment', 'text-gray-700 hover:bg-gray-100 hover:text-emerald-600': currentRankType !== 'comment'}">评论榜</a></li>
                        </ul>
                    </div>
                </div>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';

// 获取路由信息
const route = useRoute()
const router = useRouter()

// 导航方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}

// 当前选中的排行榜类型
const currentRankType = ref('click')

// 组件挂载时处理URL锚点
onMounted(() => {
    // 滚动到页面顶部
    window.scrollTo(0, 0)
    
    // 检查URL锚点
    const hash = route.hash.replace('#', '')
    if (hash) {
        // 根据锚点设置对应的排行榜类型
        switch (hash) {
            case 'newest':
                currentRankType.value = 'new'
                break
            case 'updatest':
                currentRankType.value = 'update'
                break
            case 'clickest':
                currentRankType.value = 'click'
                break
            default:
                currentRankType.value = 'click'
        }
    }
})

// 辅助函数：生成随机封面图片
const getRandomCover = (id) => {
    // 使用picsum.photos生成不同的封面图片
    return `https://picsum.photos/id/${(id % 100) + 10}/120/180`
}

// 点击榜数据
const clickRankData = ref([
    {
        bookId: '2001',
        category: '网游竞技',
        title: '不是主角，身怀修罗场有...',
        latestChapter: '第780章 那我只要证据了',
        author: '减乃',
        wordCount: '263.08',
        cover: getRandomCover(1)
    },
    {
        bookId: '2002',
        category: '科幻灵异',
        title: '星河之上',
        latestChapter: '第六百五十章，我的命好苦啊！',
        author: '柳下挥',
        wordCount: '222.22',
        cover: getRandomCover(2)
    },
    {
        bookId: '2003',
        category: '武侠仙侠',
        title: '我在修仙界万古长青',
        latestChapter: '第510章 中期妖体，天狐之境...',
        author: '铁紫店',
        wordCount: '247.96',
        cover: getRandomCover(3)
    },
    {
        bookId: '2004',
        category: '历史军事',
        title: '驾崩',
        latestChapter: '1257之后，有些东西想说...',
        author: '傻爸的青蛙',
        wordCount: '632.55',
        cover: getRandomCover(4)
    },
    {
        bookId: '2005',
        category: '女生频道',
        title: '六年后，我携四个幼崽回...',
        latestChapter: '第2748章 真表现能走了',
        author: '相思一瞬',
        wordCount: '584.92',
        cover: getRandomCover(5)
    },
    {
        bookId: '2006',
        category: '武侠仙侠',
        title: '我成了女魔头的心魔',
        latestChapter: '第423章 娘娘恩宠！送我离去...',
        author: '金枪雨来',
        wordCount: '159.88',
        cover: getRandomCover(6)
    },
    {
        bookId: '2007',
        category: '都市言情',
        title: '龙王苏醒，美女总裁求婚证',
        latestChapter: '第1206章 了结（大结局）',
        author: '九品医物',
        wordCount: '297.49',
        cover: getRandomCover(7)
    },
    {
        bookId: '2008',
        category: '武侠仙侠',
        title: '独步成仙',
        latestChapter: '5925章 镇魂甲俯地！',
        author: '鲨个橘子',
        wordCount: '1335.54',
        cover: getRandomCover(8)
    },
    {
        bookId: '2009',
        category: '武侠仙侠',
        title: '仙子别怕，我是跨子',
        latestChapter: '完结感言！',
        author: '爱公子',
        wordCount: '232.35',
        cover: getRandomCover(9)
    },
    {
        bookId: '2010',
        category: '历史军事',
        title: '大满之日，我将未将要业...',
        latestChapter: '第1098章 大军到处，震惊分明',
        author: '砑武准',
        wordCount: '284.10',
        cover: getRandomCover(10)
    },
    {
        bookId: '2011',
        category: '女生频道',
        title: '七零：靠第基锅绑财大佬...',
        latestChapter: '第一章 第100章 大结局',
        author: '我宝螺解了',
        wordCount: '14.69',
        cover: getRandomCover(11)
    },
    {
        bookId: '2012',
        category: '历史军事',
        title: '年工科技',
        latestChapter: '四十四百零三二章 去莫选的他穴...',
        author: '止天龙',
        wordCount: '916.11',
        cover: getRandomCover(12)
    },
    {
        bookId: '2013',
        category: '武侠仙侠',
        title: '肠肚不得，我干乱世证长生',
        latestChapter: '第1152章 娘娘散资子',
        author: '剑圣',
        wordCount: '272.70',
        cover: getRandomCover(13)
    },
    {
        bookId: '2014',
        category: '武侠仙侠',
        title: '我在凡人科学修仙',
        latestChapter: '新有《龙珠炎务：我等你这绳...',
        author: '冶育于',
        wordCount: '768.36',
        cover: getRandomCover(14)
    },
    {
        bookId: '2015',
        category: '都市言情',
        title: '特战之王',
        latestChapter: '第十四章：切入点',
        author: '小算',
        wordCount: '1098.35',
        cover: getRandomCover(15)
    },
    {
        bookId: '2016',
        category: '玄幻奇幻',
        title: '淘藏葫芦快',
        latestChapter: '第2709章 十贷大会',
        author: '鱼划负',
        wordCount: '829.48',
        cover: getRandomCover(16)
    },
    {
        bookId: '2017',
        category: '都市言情',
        title: '宫闹扶摇',
        latestChapter: '第2960章 我纳急求果卷台你我...',
        author: '雪管神花',
        wordCount: '713.72',
        cover: getRandomCover(17)
    },
    {
        bookId: '2018',
        category: '历史军事',
        title: '乱战三国之春宙纪险',
        latestChapter: '第1880章三城试新怀结，去...',
        author: '九菁筝窝',
        wordCount: '562.77',
        cover: getRandomCover(18)
    },
    {
        bookId: '2019',
        category: '历史军事',
        title: '隔壁的机遇',
        latestChapter: '第八百六十一章 元庆未来【大...',
        author: '高月',
        wordCount: '197.83',
        cover: getRandomCover(19)
    },
    {
        bookId: '2020',
        category: '科幻灵异',
        title: '柯南里的老爷人',
        latestChapter: '3561【害怕的妹妹】决开我呆...',
        author: '仙丹',
        wordCount: '738.77',
        cover: getRandomCover(20)
    }
])

// 新书榜数据
const newBookRankData = ref([
    {
        bookId: '3001',
        category: '都市言情',
        title: '都市医神：我的贴身校花',
        latestChapter: '第100章 新的开始',
        author: '都市小能手',
        wordCount: '25.36',
        cover: getRandomCover(21)
    },
    {
        bookId: '3002',
        category: '玄幻奇幻',
        title: '九转成神：我有混沌系统',
        latestChapter: '第50章 突破元婴',
        author: '玄幻大师',
        wordCount: '12.89',
        cover: getRandomCover(22)
    },
    {
        bookId: '3003',
        category: '武侠仙侠',
        title: '剑起云涌：一剑破万法',
        latestChapter: '第30章 剑派大比',
        author: '剑仙传人',
        wordCount: '8.54',
        cover: getRandomCover(23)
    },
    {
        bookId: '3004',
        category: '历史军事',
        title: '大明：开局成为锦衣卫',
        latestChapter: '第45章 侦破大案',
        author: '历史爱好者',
        wordCount: '18.22',
        cover: getRandomCover(24)
    },
    {
        bookId: '3005',
        category: '女生频道',
        title: '总裁的贴身小秘书',
        latestChapter: '第60章 公司危机',
        author: '言情女王',
        wordCount: '15.77',
        cover: getRandomCover(25)
    },
    {
        bookId: '3006',
        category: '科幻灵异',
        title: '星际穿越：我的太空帝国',
        latestChapter: '第35章 外星文明',
        author: '科幻迷',
        wordCount: '10.33',
        cover: getRandomCover(26)
    },
    {
        bookId: '3007',
        category: '网游竞技',
        title: '电竞王者：从零开始',
        latestChapter: '第55章 战队组建',
        author: '游戏大神',
        wordCount: '22.45',
        cover: getRandomCover(27)
    },
    {
        bookId: '3008',
        category: '武侠仙侠',
        title: '仙尊重生：废柴逆袭',
        latestChapter: '第40章 打脸时刻',
        author: '修仙专家',
        wordCount: '16.88',
        cover: getRandomCover(28)
    },
    {
        bookId: '3009',
        category: '都市言情',
        title: '重生之商业帝国',
        latestChapter: '第65章 上市计划',
        author: '商业奇才',
        wordCount: '28.91',
        cover: getRandomCover(29)
    },
    {
        bookId: '3010',
        category: '历史军事',
        title: '三国：我是曹操之子',
        latestChapter: '第70章 官渡之战',
        author: '三国迷',
        wordCount: '33.56',
        cover: getRandomCover(30)
    }
])

// 更新榜数据
const updateRankData = ref([
    {
        bookId: '4001',
        category: '武侠仙侠',
        title: '独步成仙',
        latestChapter: '5926章 新的征程',
        author: '鲨个橘子',
        wordCount: '1336.21',
        cover: getRandomCover(31)
    },
    {
        bookId: '4002',
        category: '都市言情',
        title: '特战之王',
        latestChapter: '第十五章：行动开始',
        author: '小算',
        wordCount: '1098.89',
        cover: getRandomCover(32)
    },
    {
        bookId: '4003',
        category: '玄幻奇幻',
        title: '淘藏葫芦快',
        latestChapter: '第2710章 神秘宝藏',
        author: '鱼划负',
        wordCount: '829.95',
        cover: getRandomCover(33)
    },
    {
        bookId: '4004',
        category: '历史军事',
        title: '年工科技',
        latestChapter: '四十四百零三四章 新的发现',
        author: '止天龙',
        wordCount: '916.58',
        cover: getRandomCover(34)
    },
    {
        bookId: '4005',
        category: '科幻灵异',
        title: '柯南里的老爷人',
        latestChapter: '3562【真相大白】',
        author: '仙丹',
        wordCount: '739.22',
        cover: getRandomCover(35)
    },
    {
        bookId: '4006',
        category: '武侠仙侠',
        title: '我在凡人科学修仙',
        latestChapter: '新章节：突破金丹',
        author: '冶育于',
        wordCount: '768.89',
        cover: getRandomCover(36)
    },
    {
        bookId: '4007',
        category: '都市言情',
        title: '宫闹扶摇',
        latestChapter: '第2961章 新的挑战',
        author: '雪管神花',
        wordCount: '714.25',
        cover: getRandomCover(37)
    },
    {
        bookId: '4008',
        category: '历史军事',
        title: '乱战三国之春宙纪险',
        latestChapter: '第1881章 新的策略',
        author: '九菁筝窝',
        wordCount: '563.24',
        cover: getRandomCover(38)
    },
    {
        bookId: '4009',
        category: '女生频道',
        title: '六年后，我携四个幼崽回...',
        latestChapter: '第2749章 新的计划',
        author: '相思一瞬',
        wordCount: '585.37',
        cover: getRandomCover(39)
    },
    {
        bookId: '4010',
        category: '历史军事',
        title: '驾崩',
        latestChapter: '1258章 新的局势',
        author: '傻爸的青蛙',
        wordCount: '633.12',
        cover: getRandomCover(40)
    }
])

// 评论榜数据
const commentRankData = ref([
    {
        bookId: '5001',
        category: '武侠仙侠',
        title: '我成了女魔头的心魔',
        latestChapter: '第423章 娘娘恩宠！送我离去...',
        author: '金枪雨来',
        wordCount: '159.88',
        cover: getRandomCover(41)
    },
    {
        bookId: '5002',
        category: '网游竞技',
        title: '不是主角，身怀修罗场有...',
        latestChapter: '第780章 那我只要证据了',
        author: '减乃',
        wordCount: '263.08',
        cover: getRandomCover(42)
    },
    {
        bookId: '5003',
        category: '科幻灵异',
        title: '星河之上',
        latestChapter: '第六百五十章，我的命好苦啊！',
        author: '柳下挥',
        wordCount: '222.22',
        cover: getRandomCover(43)
    },
    {
        bookId: '5004',
        category: '武侠仙侠',
        title: '仙子别怕，我是跨子',
        latestChapter: '完结感言！',
        author: '爱公子',
        wordCount: '232.35',
        cover: getRandomCover(44)
    },
    {
        bookId: '5005',
        category: '女生频道',
        title: '六年后，我携四个幼崽回...',
        latestChapter: '第2748章 真表现能走了',
        author: '相思一瞬',
        wordCount: '584.92',
        cover: getRandomCover(45)
    },
    {
        bookId: '5006',
        category: '都市言情',
        title: '龙王苏醒，美女总裁求婚证',
        latestChapter: '第1206章 了结（大结局）',
        author: '九品医物',
        wordCount: '297.49',
        cover: getRandomCover(46)
    },
    {
        bookId: '5007',
        category: '武侠仙侠',
        title: '我在修仙界万古长青',
        latestChapter: '第510章 中期妖体，天狐之境...',
        author: '铁紫店',
        wordCount: '247.96',
        cover: getRandomCover(47)
    },
    {
        bookId: '5008',
        category: '历史军事',
        title: '大满之日，我将未将要业...',
        latestChapter: '第1098章 大军到处，震惊分明',
        author: '砑武准',
        wordCount: '284.10',
        cover: getRandomCover(48)
    },
    {
        bookId: '5009',
        category: '武侠仙侠',
        title: '肠肚不得，我干乱世证长生',
        latestChapter: '第1152章 娘娘散资子',
        author: '剑圣',
        wordCount: '272.70',
        cover: getRandomCover(49)
    },
    {
        bookId: '5010',
        category: '历史军事',
        title: '隔壁的机遇',
        latestChapter: '第八百六十一章 元庆未来【大...',
        author: '高月',
        wordCount: '197.83',
        cover: getRandomCover(50)
    }
])

// 根据当前选中类型返回对应的数据
const books = computed(() => {
    switch(currentRankType.value) {
        case 'new':
            return newBookRankData.value
        case 'update':
            return updateRankData.value
        case 'comment':
            return commentRankData.value
        case 'click':
        default:
            return clickRankData.value
    }
})

// 获取当前排行榜标题
const currentRankTitle = computed(() => {
    switch(currentRankType.value) {
        case 'new':
            return '新书榜'
        case 'update':
            return '更新榜'
        case 'comment':
            return '评论榜'
        case 'click':
        default:
            return '点击榜'
    }
})

// 底部导航栏相关
const activeTab = ref('')
const handleTabChange = (tab) => {
    activeTab.value = tab
}

// 导航到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`)
}
</script>

<style scoped>
/* 如需局部样式可在此添加 */
</style>