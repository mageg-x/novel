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
                    @click="navigateToBook(book.book.id)"
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
                                :src="book.book.cover" 
                                alt="{{ book.book.title }}" 
                                class="w-16 h-22 rounded-lg object-cover shadow-md transition-transform duration-300 hover:scale-105"
                            />
                        </div>

                        <!-- 书籍信息 -->
                        <div class="flex-1">
                            <div class="flex items-center justify-between mb-1">
                                <div class="flex items-center">
                                    <span class="text-xs px-2 py-0.5 bg-emerald-100 text-emerald-700 rounded-full font-medium mr-2">
                                        [{{ book.book.category }}]
                                    </span>
                                    <span class="text-gray-500 text-xs font-medium">{{ book.book.author }}</span>
                                </div>
                                <span class="text-sm text-gray-600 font-medium ml-2 whitespace-nowrap">
                                    {{ (book.book.wordCount / 10000).toFixed(1) }}万
                                </span>
                            </div>
                            
                            <h3 class="text-base font-semibold text-gray-800 mb-1 line-clamp-1 hover:text-emerald-600 transition-colors">
                                {{ book.book.title }}
                            </h3>
                            
                            <p class="text-xs text-gray-500 line-clamp-1">{{ book.book.latestChapterTitle }}</p>
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
                        <div class="col-span-3">最新章节</div>
                        <div class="col-span-2">作者</div>
                        <div class="col-span-2 text-right">字数</div>
                    </div>

                    <!-- 排行榜列表 -->
                    <div class="divide-y divide-gray-100">
                        <div v-for="(book, index) in books" :key="index"
                        class="grid grid-cols-13 gap-4 px-6 py-3 hover:bg-gray-50 items-center transition-colors duration-150">
                        <div class="col-span-1 overflow-hidden">
                            <span :class="{ 'bg-red-500': index === 0, 'bg-orange-500': index === 1, 'bg-yellow-500': index === 2, 'bg-blue-500': index > 2 && index < 10, 'bg-gray-300': index >= 10 }"
                                class="inline-flex items-center justify-center w-7 h-7 text-white rounded-full text-sm font-bold">{{ index + 1 }}</span>
                        </div>
                        <div class="col-span-2 overflow-hidden text-gray-500 text-sm whitespace-nowrap pr-2">{{ `[${book.book.category}]` }}</div>
                        <div class="col-span-3 whitespace-nowrap overflow-hidden text-ellipsis"><router-link :to="'/book/' + book.book.id" class="text-gray-600 text-sm hover:text-blue-800 hover:underline transition-colors cursor-pointer">{{ book.book.title}}</router-link></div>
                        <div class="col-span-3 whitespace-nowrap overflow-hidden text-ellipsis text-gray-600 text-sm">{{ book.book.latestChapterTitle }}</div>
                        <div class="col-span-2 whitespace-nowrap overflow-hidden text-ellipsis text-gray-600 text-sm">{{ book.book.author }}</div>
                        <div class="col-span-2 overflow-hidden text-gray-500 text-sm text-right whitespace-nowrap">{{ (book.book.wordCount / 10000).toFixed(1) }}万</div>
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

// 排行榜数据 - 从API获取
const clickRankData = ref([]);
const newBookRankData = ref([]);
const updateRankData = ref([]);
const commentRankData = ref([]);

// 辅助函数：生成随机封面图片（API返回数据可能需要）
const getRandomCover = (id) => {
    // 使用picsum.photos生成不同的封面图片
    return `https://picsum.photos/id/${(id % 100) + 10}/120/180`
}

// 导入API服务
import { rankAPI } from '@/api/services';

// 根据当前选中的排行榜类型计算要显示的书籍列表
const books = computed(() => {
    switch (currentRankType.value) {
        case 'click':
            return clickRankData.value;
        case 'new':
            return newBookRankData.value;
        case 'update':
            return updateRankData.value;
        case 'comment':
            return commentRankData.value;
        default:
            return [];
    }
})

// 根据当前选中的排行榜类型显示对应的标题
const currentRankTitle = computed(() => {
    switch (currentRankType.value) {
        case 'click':
            return '点击榜';
        case 'new':
            return '新书榜';
        case 'update':
            return '更新榜';
        case 'comment':
            return '评论榜';
        default:
            return '';
    }
})

// 切换排行榜类型
const changeRankType = (type) => {
    currentRankType.value = type
}

// 跳转到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`)
}



// 获取所有排行榜数据
onMounted(async () => {
    try {
        // 获取点击榜数据
        const clickRankResponse = await rankAPI.getClickRank();
        clickRankData.value = clickRankResponse.data;
        
        // 获取新书榜数据
        const newBookRankResponse = await rankAPI.getNewBookRank();
        newBookRankData.value = newBookRankResponse.data;
        
        // 获取更新榜数据
        const updateRankResponse = await rankAPI.getUpdateRank();
        updateRankData.value = updateRankResponse.data;
        
        // 获取评论榜数据
        const commentRankResponse = await rankAPI.getCommentRank();
        commentRankData.value = commentRankResponse.data;
    } catch (error) {
        console.error('获取排行榜数据失败:', error);
    }
});
</script>

<style scoped>
/* 如需局部样式可在此添加 */
</style>