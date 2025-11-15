<template>

    <!-- 小屏幕 -->
    <div class="block md:hidden">
        <!-- 顶部导航 -->
        <div class="sticky top-0 z-50 bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-lg">
            <div class="flex items-center justify-between h-14 px-4">
                <button @click="goBack" class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <span class="text-lg font-semibold">{{ activeTab === 'shelf' ? '我的书架' : '我的阅读' }}</span>
                <button @click="goHome" class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fa-solid fa-house"></i>
                </button>
            </div>
        </div>
        
        <!-- 内容区域 -->
        <div class="pb-16 bg-gradient-to-b from-gray-50 to-white">
            <div class="px-3 py-4">
                <div v-for="(book, index) in displayedBooks" :key="index" class="flex items-center bg-white rounded-2xl shadow-md mb-4 overflow-hidden border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
                    <!-- 封面缩略图 -->
                    <div class="w-20 h-28 flex-shrink-0 overflow-hidden">
                        <img :src="book.cover" :alt="book.title" class="w-full h-full object-cover transition-transform duration-500 hover:scale-110" />
                    </div>
                    
                    <!-- 书籍信息 -->
                    <div class="flex-1 p-4">
                        <div class="text-base font-semibold text-gray-900 mb-2 line-clamp-2 leading-tight">{{ book.title }}</div>
                        <div class="text-sm text-emerald-600 mb-2 line-clamp-1 font-medium">{{ book.latestChapter }}</div>
                        <div class="text-xs text-gray-500 mb-3">{{ book.updateTime }}</div>
                        <div class="flex items-center gap-4 pt-2">
                            <button class="text-sm text-emerald-600 hover:text-emerald-800 font-medium underline cursor-pointer transition-colors duration-200">
                                继续阅读
                            </button>
                            <button v-if="activeTab === 'shelf'" class="text-sm text-red-500 hover:text-red-700 font-medium underline cursor-pointer transition-colors duration-200">
                                移除书架
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        </div>

        <!-- 底部导航 -->
        <ToolBar :activeTab="activeTab" @tab-change="handleTabChange" />
    </div>

    <!-- 大屏 -->
    <div class="hidden md:block container md:w-5xl  mx-auto">
        <div class="hidden md:block">
            <Header />
        </div>
        <!-- 标签页 -->
        <div class="flex items-center gap-2 mt-6 mb-4 p-1 bg-gray-50 rounded-lg shadow-sm">
            <button 
                :class="{
                    'bg-white text-emerald-600 font-semibold': activeTab === 'shelf',
                    'text-gray-600 hover:text-gray-800': activeTab !== 'shelf'
                }" 
                @click="activeTab = 'shelf'"
                class="px-5 py-2 rounded-md transition-all duration-200 hover:shadow-sm"
            >
                <i class="fas fa-book mr-2"></i> 我的书架
            </button>
            <button 
                :class="{
                    'bg-white text-emerald-600 font-semibold': activeTab === 'recent',
                    'text-gray-600 hover:text-gray-800': activeTab !== 'recent'
                }" 
                @click="activeTab = 'recent'"
                class="px-5 py-2 rounded-md transition-all duration-200 hover:shadow-sm"
            >
                <i class="fas fa-clock mr-2"></i> 最近阅读
            </button>
        </div>

        <!-- 书籍列表 -->
        <div class="bg-white rounded-xl shadow-md overflow-hidden border border-gray-100">
            <table class="min-w-full divide-y divide-gray-100">
                <thead class="bg-gradient-to-r from-emerald-50 to-teal-50">
                    <tr>
                        <th class="px-6 py-4 text-center text-xs font-semibold text-gray-700 uppercase tracking-wider">
                            封面
                        </th>
                        <th class="px-6 py-4 text-center text-xs font-semibold text-gray-700 uppercase tracking-wider">
                            类别
                        </th>
                        <th class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                            书名
                        </th>
                        <th class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                            最新章节
                        </th>
                        <th class="px-6 py-4 text-center text-xs font-semibold text-gray-700 uppercase tracking-wider">
                            更新时间
                        </th>
                        <th class="px-6 py-4 text-center text-xs font-semibold text-gray-700 uppercase tracking-wider">
                            操作
                        </th>
                    </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-50">
                    <tr v-for="(book, index) in displayedBooks" :key="index" class="hover:bg-gray-50 transition-all duration-150 hover:shadow-sm">
                        <td class="px-6 py-4 whitespace-nowrap">
                            <img :src="book.cover" :alt="book.title" class="w-14 h-18 object-cover rounded-md shadow-sm" />
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-medium">{{ book.category }}</td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 hover:text-emerald-700 transition-colors">
                            {{ book.title }}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-emerald-600 hover:text-emerald-800 transition-colors">
                            {{ book.latestChapter }}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ book.updateTime }}</td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm">
                            <button class="text-emerald-600 hover:text-emerald-800 font-medium underline cursor-pointer mr-3 transition-colors duration-200">
                                继续阅读
                            </button>
                            <button v-if="activeTab === 'shelf'" class="text-red-500 hover:text-red-700 font-medium underline cursor-pointer transition-colors duration-200">
                                移除书架
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

        <!-- 分页（简化版） -->
        <div class="flex justify-center items-center mt-6 space-x-1 text-sm">
            <button class="text-gray-500 hover:text-gray-700 px-3 py-1 rounded">上一页</button>
            <button class="bg-emerald-500 text-white px-3 py-1 rounded">1</button>
            <button class="text-gray-500 hover:text-gray-700 px-3 py-1 rounded">2</button>
            <button class="text-gray-500 hover:text-gray-700 px-3 py-1 rounded">下一页</button>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';

const route = useRoute()
const router = useRouter()

// 回退和首页跳转方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}
// 从URL参数获取tab值，默认使用'shelf'
const activeTab = ref(route.query.tab === 'recent' || route.query.tab === 'shelf' ? route.query.tab : 'shelf')

// 监听路由参数变化，更新activeTab
watch(() => route.query.tab, (newTab) => {
    if (newTab === 'recent' || newTab === 'shelf') {
        activeTab.value = newTab
    }
})

// 监听activeTab变化，更新URL参数
watch(activeTab, (newTab) => {
    router.replace({ query: { tab: newTab } })
})

// 处理底部导航切换
const handleTabChange = (tab) => {
    if (tab === 'recent' || tab === 'shelf') {
        activeTab.value = tab
    }
}

// ———————— 模拟数据：最近阅读（30条） ————————
const recentBooks = ref([
    { category: '[网游竞技]', title: '燃烬之龙', latestChapter: '第四百零七章 海洛薇茨', updateTime: '11/08 01:57', cover: 'https://picsum.photos/id/31/120/160' },
    { category: '[网游竞技]', title: '退队，然后捡到问题美少女', latestChapter: '710.圣都怪事专版', updateTime: '11/14 00:09', cover: 'https://picsum.photos/id/32/120/160' },
    { category: '[女生频道]', title: '影视诸天从知否开始', latestChapter: '第四十章 走着瞧', updateTime: '07/23 02:59', cover: 'https://picsum.photos/id/33/120/160' },
    { category: '[女生频道]', title: '穿成农家长姐，开局养全家', latestChapter: '第243章 抢救', updateTime: '08/03 03:22', cover: 'https://picsum.photos/id/34/120/160' },
    { category: '[都市言情]', title: '龙王苏醒，美女总裁求婚证', latestChapter: '第1206章 了结（大结局）', updateTime: '11/13 22:15', cover: 'https://picsum.photos/id/35/120/160' },
    { category: '[武侠仙侠]', title: '我在修仙界万古长青', latestChapter: '第510章 中期妖体，天狐之境...', updateTime: '11/14 08:30', cover: 'https://picsum.photos/id/36/120/160' },
    { category: '[科幻灵异]', title: '星河之上', latestChapter: '第六百五十章，我的命好苦啊！', updateTime: '11/12 19:45', cover: 'https://picsum.photos/id/37/120/160' },
    { category: '[历史军事]', title: '驾崩', latestChapter: '1257之后，有些东西想说...', updateTime: '11/10 14:20', cover: 'https://picsum.photos/id/38/120/160' },
    { category: '[女生频道]', title: '六年后，我携四个幼崽回...', latestChapter: '第2748章 真表现能走了', updateTime: '11/14 07:12', cover: 'https://picsum.photos/id/39/120/160' },
    { category: '[武侠仙侠]', title: '独步成仙', latestChapter: '5925章 镇魂甲俯地！', updateTime: '11/14 05:00', cover: 'https://picsum.photos/id/40/120/160' },
    { category: '[玄幻奇幻]', title: '淘藏葫芦快', latestChapter: '第2709章 十贷大会', updateTime: '11/13 23:10', cover: 'https://picsum.photos/id/41/120/160' },
    { category: '[都市言情]', title: '特战之王', latestChapter: '第十四章：切入点', updateTime: '11/11 16:40', cover: 'https://picsum.photos/id/42/120/160' },
    { category: '[历史军事]', title: '我娘子天下第一', latestChapter: '第七百零二章', updateTime: '11/09 12:05', cover: 'https://picsum.photos/id/43/120/160' },
    { category: '[科幻灵异]', title: '柯南里的老爷人', latestChapter: '3561【害怕的妹妹】决开我呆...', updateTime: '11/14 01:30', cover: 'https://picsum.photos/id/44/120/160' },
    { category: '[网游竞技]', title: '精灵：这个训练家太野了', latestChapter: '新书#2开，感冬侧鬼无获封不下', updateTime: '11/12 10:25', cover: 'https://picsum.photos/id/45/120/160' },
    { category: '[女生频道]', title: '快穿年代女配', latestChapter: '第1448章 穿越上千老去角清...', updateTime: '11/13 18:55', cover: 'https://picsum.photos/id/46/120/160' },
    { category: '[武侠仙侠]', title: '我成了女魔头的心魔', latestChapter: '第423章 娘娘恩宠！送我离去...', updateTime: '11/11 09:40', cover: 'https://picsum.photos/id/47/120/160' },
    { category: '[历史军事]', title: '乱战三国之春宙纪险', latestChapter: '第1880章三城试新怀结，去...', updateTime: '11/10 20:15', cover: 'https://picsum.photos/id/48/120/160' },
    { category: '[女生频道]', title: '系统选王角的我加入猎天群', latestChapter: '第二十章五十三章 本末倒屋', updateTime: '11/12 06:50', cover: 'https://picsum.photos/id/49/120/160' },
    { category: '[玄幻奇幻]', title: '苟在初圣魔门当入材', latestChapter: '第一十一百三十四章 手发选螺...', updateTime: '11/13 14:30', cover: 'https://picsum.photos/id/50/120/160' },
    { category: '[都市言情]', title: '宫闹扶摇', latestChapter: '第2960章 我纳急求果卷台你我...', updateTime: '11/14 02:20', cover: 'https://picsum.photos/id/51/120/160' },
    { category: '[历史军事]', title: '隔壁的机遇', latestChapter: '第八百六十一章 元庆未来【大...', updateTime: '11/09 22:40', cover: 'https://picsum.photos/id/52/120/160' },
    { category: '[网游竞技]', title: '斗破面具人', latestChapter: '第176章 线章', updateTime: '11/11 11:05', cover: 'https://picsum.photos/id/53/120/160' },
    { category: '[女生频道]', title: '七零：靠第基锅绑财大佬...', latestChapter: '第一章 第100章 大结局', updateTime: '11/08 17:30', cover: 'https://picsum.photos/id/54/120/160' },
    { category: '[科幻灵异]', title: '重生来世：开局中奖300...', latestChapter: '第2141章：减，还是要继封的！', updateTime: '11/13 08:15', cover: 'https://picsum.photos/id/55/120/160' },
    { category: '[武侠仙侠]', title: '肠肚不得，我干乱世证长生', latestChapter: '第1152章 娘娘散资子', updateTime: '11/12 13:20', cover: 'https://picsum.photos/id/56/120/160' },
    { category: '[历史军事]', title: '年工科技', latestChapter: '四十四百零三二章 去莫选的他穴...', updateTime: '11/10 09:10', cover: 'https://picsum.photos/id/57/120/160' },
    { category: '[武侠仙侠]', title: '仙子别怕，我是跨子', latestChapter: '完结感言！', updateTime: '11/07 21:00', cover: 'https://picsum.photos/id/58/120/160' },
    { category: '[都市言情]', title: '最强赘婿', latestChapter: '第3201章 隐藏身份曝光', updateTime: '11/14 03:45', cover: 'https://picsum.photos/id/59/120/160' },
    { category: '[玄幻奇幻]', title: '万古神帝', latestChapter: '第4120章 天庭崩塌', updateTime: '11/13 19:20', cover: 'https://picsum.photos/id/60/120/160' }
])

// ———————— 模拟数据：我的书架（20条） ————————
const shelfBooks = ref([
    { category: '[都市言情]', title: '最强医圣', latestChapter: '第2890章 神秘组织', updateTime: '11/14 06:10', cover: 'https://ext.same-assets.com/4081703551/959896919.jpeg' },
    { category: '[武侠仙侠]', title: '凡人修仙传', latestChapter: '第2105章 化神之战', updateTime: '11/13 20:30', cover: 'https://ext.same-assets.com/4081703551/1339035738.jpeg' },
    { category: '[女生频道]', title: '嫡女重生：王爷请接招', latestChapter: '第1892章 大婚风波', updateTime: '11/12 15:40', cover: 'https://ext.same-assets.com/4081703551/594530825.png' },
    { category: '[历史军事]', title: '回到明朝当王爷', latestChapter: '第1560章 皇权更迭', updateTime: '11/11 18:20', cover: 'https://ext.same-assets.com/4081703551/3622865948.jpeg' },
    { category: '[科幻灵异]', title: '末日蟑螂', latestChapter: '第3420章 新人类崛起', updateTime: '11/10 23:50', cover: 'https://picsum.photos/id/11/120/160' },
    { category: '[网游竞技]', title: '全职高手', latestChapter: '第1780章 荣耀全明星', updateTime: '11/09 16:15', cover: 'https://picsum.photos/id/12/120/160' },
    { category: '[玄幻奇幻]', title: '斗破苍穹', latestChapter: '第1650章 萧炎归来', updateTime: '11/08 12:40', cover: 'https://picsum.photos/id/13/120/160' },
    { category: '[都市言情]', title: '校花的贴身高手', latestChapter: '第9870章 隐藏BOSS', updateTime: '11/14 04:20', cover: 'https://picsum.photos/id/14/120/160' },
    { category: '[武侠仙侠]', title: '遮天', latestChapter: '第1800章 成仙路开启', updateTime: '11/13 10:05', cover: 'https://picsum.photos/id/15/120/160' },
    { category: '[女生频道]', title: '重生之将门毒后', latestChapter: '第1200章 凤凰涅槃', updateTime: '11/12 09:30', cover: 'https://picsum.photos/id/16/120/160' },
    { category: '[历史军事]', title: '大秦帝国', latestChapter: '第890章 统一六国', updateTime: '11/11 14:10', cover: 'https://picsum.photos/id/17/120/160' },
    { category: '[科幻灵异]', title: '吞噬星空', latestChapter: '第1450章 宇宙级强者', updateTime: '11/10 17:25', cover: 'https://picsum.photos/id/18/120/160' },
    { category: '[网游竞技]', title: '英雄联盟：我的时代', latestChapter: '第620章 世界赛冠军', updateTime: '11/09 08:55', cover: 'https://picsum.photos/id/19/120/160' },
    { category: '[玄幻奇幻]', title: '完美世界', latestChapter: '第2000章 石昊成帝', updateTime: '11/08 21:30', cover: 'https://picsum.photos/id/20/120/160' },
    { category: '[都市言情]', title: '神医凰后', latestChapter: '第1500章 凤凰觉醒', updateTime: '11/14 00:45', cover: 'https://picsum.photos/id/21/120/160' },
    { category: '[武侠仙侠]', title: '一念永恒', latestChapter: '第1300章 白小纯飞升', updateTime: '11/13 05:20', cover: 'https://picsum.photos/id/22/120/160' },
    { category: '[女生频道]', title: '锦绣未央', latestChapter: '第980章 权谋巅峰', updateTime: '11/12 22:10', cover: 'https://picsum.photos/id/23/120/160' },
    { category: '[历史军事]', title: '三国之超级召唤系统', latestChapter: '第1100章 诸葛亮出山', updateTime: '11/11 07:40', cover: 'https://picsum.photos/id/24/120/160' },
    { category: '[科幻灵异]', title: '全球进化', latestChapter: '第780章 人类反击', updateTime: '11/10 11:15', cover: 'https://picsum.photos/id/25/120/160' },
    { category: '[网游竞技]', title: '超神机械师', latestChapter: '第1200章 星海终局', updateTime: '11/09 13:00', cover: 'https://picsum.photos/id/26/120/160' }
])

// 计算当前显示的书籍
const displayedBooks = computed(() => {
    return activeTab.value === 'recent' ? recentBooks.value : shelfBooks.value
})
</script>