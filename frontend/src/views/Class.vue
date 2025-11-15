<template>
        <!-- 小屏幕 -->
    <div class="block md:hidden">
        <!-- 顶部导航 -->
        <div class="sticky top-0 z-50 bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-lg">
            <div class="flex items-center justify-between h-14 px-4">
                <button @click="goBack" class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <span class="text-lg font-semibold">全部作品</span>
                <button @click="goHome" class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fa-solid fa-house"></i>
                </button>
            </div>
        </div>
        
        <!-- 内容区域 -->
        <div class="bg-gradient-to-br from-gray-50 to-gray-100 min-h-screen pb-20">
            <!-- 筛选条件（可折叠整体面板） -->
            <div class="sticky top-14 z-40 bg-white shadow-md mb-4">
                <!-- 筛选面板顶部 -->
                <button 
                    class="w-full flex items-center justify-between p-4 text-sm font-medium text-gray-800 hover:bg-gray-50 transition-colors"
                    @click="showFilters = !showFilters"
                >
                    <span class="flex items-center">
                        <i class="fas fa-filter mr-2 text-emerald-500"></i>
                        筛选条件
                    </span>
                    <i :class="['fas', showFilters ? 'fa-chevron-up' : 'fa-chevron-down', 'text-gray-400 transition-transform duration-200']"></i>
                </button>
                
                <!-- 已选筛选条件 -->
                <div v-if="!showFilters && getSelectedFilters().length > 0" class="px-4 pb-3 overflow-x-auto">
                    <div class="flex flex-wrap gap-2">
                        <div 
                            v-for="(filter, idx) in getSelectedFilters()" 
                            :key="idx"
                            class="flex items-center bg-emerald-50 text-emerald-700 text-xs px-3 py-1 rounded-full shadow-sm"
                        >
                            <span class="mr-1">{{ filter.text }}</span>
                            <i class="fas fa-times text-xs cursor-pointer hover:text-emerald-900 transition-colors" @click="clearFilter(filter)"></i>
                        </div>
                    </div>
                </div>
                
                <!-- 筛选面板内容 -->
                <div v-if="showFilters" class="p-4 space-y-4 max-h-80 overflow-y-auto">
                    <!-- 作品分类筛选 -->
                    <div>
                        <h3 class="text-sm font-semibold text-gray-700 mb-2">{{ filters[1].label }}</h3>
                        <div class="flex overflow-x-auto pb-2 -mx-2 px-2 space-x-2">
                            <button 
                                v-for="(option, optIdx) in filters[1].options" 
                                :key="optIdx"
                                :class="[
                                    'flex-shrink-0 px-4 py-2 rounded-full text-sm transition-all duration-200',
                                    option.active 
                                        ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md'
                                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                                ]"
                                @click="selectFilter(filters[1], option)"
                            >
                                {{ option.text }}
                            </button>
                        </div>
                    </div>
                    
                    <!-- 其他筛选条件 -->
                    <div class="space-y-3">
                        <div v-for="(filter, idx) in filters.slice(0, 1).concat(filters.slice(2))" :key="idx">
                            <h3 class="text-sm font-semibold text-gray-700 mb-2">{{ filter.label }}</h3>
                            <div class="grid grid-cols-2 gap-2">
                                <button 
                                    v-for="(option, optIdx) in filter.options" 
                                    :key="optIdx"
                                    :class="[
                                        'text-left px-3 py-2 rounded-md text-sm transition-all duration-200',
                                        option.active 
                                            ? 'bg-emerald-100 text-emerald-700 font-medium'
                                            : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                    ]"
                                    @click="selectFilter(filter, option)"
                                >
                                    {{ option.text }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- 书籍列表 -->
            <div class="container mx-auto px-4">
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
                                    {{ book.wordCount }}
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
            
            <!-- 分页 -->
            <div class="flex justify-center items-center space-x-1 py-8 px-4">
                <button class="px-3 py-2 text-gray-500 hover:text-emerald-500 rounded-md hover:bg-gray-100 transition-colors">
                    <i class="fas fa-chevron-left"></i>
                </button>
                
                <button 
                    v-for="page in [1, 2, 3, 4, 5]" 
                    :key="page"
                    :class="[
                        'px-3 py-2 rounded-md text-sm transition-all duration-200',
                        page === 1 
                            ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md'
                            : 'text-gray-600 hover:bg-gray-100'
                    ]"
                >
                    {{ page }}
                </button>
                
                <span class="px-3 py-2 text-gray-400">...</span>
                
                <button class="px-3 py-2 text-gray-600 hover:bg-gray-100 rounded-md transition-colors">41</button>
                
                <button class="px-3 py-2 text-gray-500 hover:text-emerald-500 rounded-md hover:bg-gray-100 transition-colors">
                    <i class="fas fa-chevron-right"></i>
                </button>
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

        <main class="max-w-7xl mx-auto  py-6">
            <!-- Filters Section -->
            <div class="bg-white rounded-lg shadow-sm p-6 mb-6">
                <div v-for="(filter, idx) in filters" :key="idx" class="flex items-start mb-4 text-sm">
                    <div class="w-20 text-gray-500 pt-1">{{ filter.label }}：</div>
                    <div class="flex-1 flex flex-wrap gap-3">
                        <a v-for="(option, optIdx) in filter.options" :key="optIdx" href="#" :class="[
                            option.active
                                ? 'text-green-500 hover:text-green-600'
                                : 'text-gray-600 hover:text-green-500'
                        ]">
                            {{ option.text }}
                        </a>
                    </div>
                </div>
            </div>

            <!-- Books Table -->
            <div class="bg-white rounded-lg shadow-sm overflow-hidden">
                <table class="w-full text-sm">
                    <thead class="bg-gray-50 ">
                        <tr class="text-gray-500">
                            <th class="px-4 py-3 text-left max-w-16">序号</th>
                            <th class="px-4 py-3 text-left max-w-32">类别</th>
                            <th class="px-4 py-3 text-left">书名</th>
                            <th class="px-4 py-3 text-left">最新章节</th>
                            <th class="px-4 py-3 text-left max-w-32">作者</th>
                            <th class="px-4 py-3 text-left max-w-24">字数</th>
                        </tr>
                    </thead>
                    <tbody class="divide-y">
                        <tr v-for="(book, index) in books" :key="index" class="hover:bg-gray-50 border-0">
                            <td class="px-4 py-3">
                                <span :class="[
                                    'inline-block w-6 h-6 text-white text-center rounded',
                                    index === 0 ? 'bg-red-500' : '',
                                    index === 1 ? 'bg-orange-500' : '',
                                    index === 2 ? 'bg-yellow-500' : '',
                                    index >= 3 && index <= 9 ? 'bg-green-500' : '',
                                    index >= 10 ? 'bg-gray-400' : ''
                                ]">
                                    {{ index + 1 }}
                                </span>
                            </td>
                            <td class="px-4 py-3 text-gray-600">[{{ book.category }}]</td>
                            <td class="px-4 py-3">
                                <router-link :to="'/book/' + book.bookId" class="text-gray-800 hover:text-green-500 cursor-pointer">{{ book.title }}</router-link>
                            </td>
                            <td class="px-4 py-3 text-gray-600">{{ book.latestChapter }}</td>
                            <td class="px-4 py-3 text-gray-600">{{ book.author }}</td>
                            <td class="px-4 py-3 text-gray-600">{{ book.wordCount }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Pagination -->
            <div class="flex justify-center items-center space-x-2 py-8 text-sm">
                <a href="#" class="px-3 py-1 text-gray-500 hover:text-green-500">上一页</a>

                <a v-for="page in [1, 2, 3, 4, 5]" :key="page" href="#" :class="[
                    page === 1
                        ? 'px-3 py-1 bg-green-500 text-white rounded'
                        : 'px-3 py-1 text-gray-600 hover:bg-gray-100 rounded'
                ]">
                    {{ page }}
                </a>

                <span class="px-3 py-1 text-gray-400">...</span>

                <a href="#" class="px-3 py-1 text-gray-600 hover:bg-gray-100 rounded">41</a>

                <a href="#" class="px-3 py-1 text-gray-500 hover:text-green-500">下一页</a>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted  } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';

const router = useRouter()

// 回退和首页跳转方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}

// 底部导航栏相关
const activeTab = ref('class')

const handleTabChange = (tab) => {
    activeTab.value = tab
    router.push(`/${tab}`)
}

// 导航到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`)
}

// ===== 模拟筛选条件数据（完整还原）=====
// 控制筛选面板的显示与隐藏
const showFilters = ref(false);

const filters = ref([
    {
        label: '作品频道',
        options: [
            { text: '男频', active: true },
            { text: '女频', active: false }
        ]
    },
    {
        label: '作品分类',
        options: [
            { text: '不限', active: true },
            { text: '玄幻奇幻', active: false },
            { text: '武侠仙侠', active: false },
            { text: '都市言情', active: false },
            { text: '历史军事', active: false },
            { text: '科幻灵异', active: false },
            { text: '网游竞技', active: false }
        ]
    },
    {
        label: '是否完结',
        options: [
            { text: '不限', active: true },
            { text: '连载中', active: false },
            { text: '已完结', active: false }
        ]
    },
    {
        label: '作品字数',
        options: [
            { text: '不限', active: true },
            { text: '30万字以下', active: false },
            { text: '30-50万字', active: false },
            { text: '50-100万字', active: false },
            { text: '100万字以上', active: false }
        ]
    },
    {
        label: '更新时间',
        options: [
            { text: '不限', active: true },
            { text: '三日内', active: false },
            { text: '七日内', active: false },
            { text: '半月内', active: false },
            { text: '一月内', active: false }
        ]
    },
    {
        label: '排序方式',
        options: [
            { text: '不限', active: true },
            { text: '更新时间', active: false },
            { text: '总字数', active: false },
            { text: '点击量', active: false }
        ]
    }
])

// 获取所有选中的筛选条件
const getSelectedFilters = () => {
    const selected = [];
    filters.value.forEach(filter => {
        filter.options.forEach(option => {
            if (option.active && option.text !== '不限') {
                selected.push(option);
            }
        });
    });
    return selected;
};

// 清除指定的筛选条件
const clearFilter = (filterToClear) => {
    filters.value.forEach(filter => {
        // 找到包含该选项的筛选条件组
        const optionIndex = filter.options.findIndex(option => option === filterToClear);
        if (optionIndex !== -1) {
            // 将当前选项设置为非激活
            filter.options[optionIndex].active = false;
            // 找到"不限"选项并设置为激活
            const unlimitedOption = filter.options.find(option => option.text === '不限');
            if (unlimitedOption) {
                unlimitedOption.active = true;
            }
        }
    });
};

// 选择筛选条件
const selectFilter = (filter, option) => {
    // 将该筛选条件组下的所有选项设置为非激活状态
    filter.options.forEach(opt => {
        opt.active = false;
    });
    // 将选中的选项设置为激活状态
    option.active = true;
};

// ===== 模拟书籍数据（20条完整还原）=====
const books = ref([
    {
        bookId: '1001',
        category: '都市言情',
        title: '故梦难填',
        latestChapter: '第258章 未来',
        author: '山与沙',
        wordCount: '30.17万'
    },
    {
        bookId: '1002',
        category: '都市言情',
        title: '直播捡漏：我眼中却是是...',
        latestChapter: '第 1626章 宁见神山水双叠画',
        author: '命取每日万更',
        wordCount: '451.16万'
    },
    {
        bookId: '1003',
        category: '科幻灵异',
        title: '假期兼职被抓，同事治国...',
        latestChapter: '第483章 将军之问，枢密朝阁',
        author: '荣朗拿的主师',
        wordCount: '126.34万'
    },
    {
        bookId: '1004',
        category: '玄幻奇幻',
        title: '荒古武神',
        latestChapter: '第3545章 开天',
        author: '化十',
        wordCount: '846.43万'
    },
    {
        bookId: '1005',
        category: '武侠仙侠',
        title: '正经道爷，女居士请自重',
        latestChapter: '第671章 冰融城，天地造宝',
        author: '黑山书生',
        wordCount: '140.56万'
    },
    {
        bookId: '1006',
        category: '网游竞技',
        title: '足球：国足开始，来刚足...',
        latestChapter: '第3038章 土为国已者落！连击...',
        author: '不正芬',
        wordCount: '834.31万'
    },
    {
        bookId: '1007',
        category: '历史军事',
        title: '家父慷慨祖',
        latestChapter: '第二六六章 引狼入室',
        author: '今宵朝吹月',
        wordCount: '63.44万'
    },
    {
        bookId: '1008',
        category: '历史军事',
        title: '我，大楚最狂太子',
        latestChapter: '第一卷 第161章 大场局',
        author: '验牛王子',
        wordCount: '25.87万'
    },
    {
        bookId: '1009',
        category: '科幻灵异',
        title: '强制匹配：炸实力带高危...',
        latestChapter: '第三十一章双双杀',
        author: '浅海光影',
        wordCount: '6.41万'
    },
    {
        bookId: '1010',
        category: '科幻灵异',
        title: '末日：我真不是土世',
        latestChapter: '章外2 午前的生活',
        author: '草梦莲缘',
        wordCount: '259.14万'
    },
    {
        bookId: '1011',
        category: '历史军事',
        title: '大清要完',
        latestChapter: '第1037章 后记——天因之梦之一...',
        author: '大罗罗',
        wordCount: '273.25万'
    },
    {
        bookId: '1012',
        category: '武侠仙侠',
        title: '修仙：从怀孕那郎开始',
        latestChapter: '第1005章 毒场',
        author: '一林雕时',
        wordCount: '205.83万'
    },
    {
        bookId: '1013',
        category: '历史军事',
        title: '制胜记',
        latestChapter: '第一千三章 校石万千攻',
        author: '无开',
        wordCount: '338.50万'
    },
    {
        bookId: '1014',
        category: '武侠仙侠',
        title: '高武来世：从不修炼开始...',
        latestChapter: '第658章 一块据盾席',
        author: '太古肥龙',
        wordCount: '264.14万'
    },
    {
        bookId: '1015',
        category: '科幻灵异',
        title: '灭世洪水，我在浮桥之上...',
        latestChapter: '第539章 绝处长米，突破百万1',
        author: '西宁编',
        wordCount: '86.81万'
    },
    {
        bookId: '1016',
        category: '都市言情',
        title: '我的系统能是废物',
        latestChapter: '第 344 章 水团的妹妹真地',
        author: '指明宁击',
        wordCount: '96.57万'
    },
    {
        bookId: '1017',
        category: '科幻灵异',
        title: '从机械师开始无敌转世',
        latestChapter: '第七章+八章 续续（大结局）',
        author: '宇由无羊欧十',
        wordCount: '278.74万'
    },
    {
        bookId: '1018',
        category: '都市言情',
        title: '别问我捡鬼',
        latestChapter: '第2478章 调音变化，为为素组',
        author: '陈青锋',
        wordCount: '580.00万'
    },
    {
        bookId: '1019',
        category: '武侠仙侠',
        title: '九兽剑师',
        latestChapter: '第473章 样火煨质',
        author: '钉言',
        wordCount: '74.06万'
    },
    {
        bookId: '1020',
        category: '历史军事',
        title: '提前是准游戏世界，开局...',
        latestChapter: '第978章 魔灭的来宗！',
        author: '谨筵花花兔',
        wordCount: '221.80万'
    }
])
</script>

<style scoped>
/* 如需局部样式可在此添加 */
</style>