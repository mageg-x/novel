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
                <div v-for="(item, index) in displayedBooks" :key="index" class="flex items-center bg-white rounded-2xl shadow-md mb-4 overflow-hidden border border-gray-100 transition-all duration-300 hover:shadow-xl hover:-translate-y-1">
                    <!-- 封面缩略图 -->
                    <div class="w-20 h-28 flex-shrink-0 overflow-hidden">
                        <img :src="item.book.cover" :alt="item.book.title" class="w-full h-full object-cover transition-transform duration-500 hover:scale-110" />
                    </div>
                    
                    <!-- 书籍信息 -->
                    <div class="flex-1 p-4">
                        <div class="text-base font-semibold text-gray-900 mb-2 line-clamp-2 leading-tight">{{ item.book.title }}</div>
                        <div class="text-sm text-emerald-600 mb-2 line-clamp-1 font-medium">{{ item.book.latestChapterTitle }}</div>
                        <div class="text-xs text-gray-500 mb-3">{{ activeTab === 'shelf' ? item.addTime : item.updateTime }}</div>
                        <div class="flex items-center gap-4 pt-2">
                            <button @click="navigateToBook(item.book.id)" class="text-sm text-emerald-600 hover:text-emerald-800 font-medium underline cursor-pointer transition-colors duration-200">
                                继续阅读
                            </button>
                            <button v-if="activeTab === 'shelf'" @click="removeFromShelf(item.book.id, $event)" class="text-sm text-red-500 hover:text-red-700 font-medium underline cursor-pointer transition-colors duration-200">
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
                    <tr v-for="(item, index) in displayedBooks" :key="index" class="hover:bg-gray-50 transition-all duration-150 hover:shadow-sm">
                        <td class="px-6 py-4 whitespace-nowrap">
                            <img :src="item.book.cover" :alt="item.book.title" class="w-14 h-18 object-cover rounded-md shadow-sm" />
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 font-medium">{{ item.book.category }}</td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900 hover:text-emerald-700 transition-colors cursor-pointer" @click="navigateToBook(item.book.id)">
                            {{ item.book.title }}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-emerald-600 hover:text-emerald-800 transition-colors">
                            {{ item.book.latestChapterTitle }}
                        </td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ activeTab === 'shelf' ? item.addTime : item.updateTime }}</td>
                        <td class="px-6 py-4 whitespace-nowrap text-sm">
                            <button @click="navigateToBook(item.book.id)" class="text-emerald-600 hover:text-emerald-800 font-medium underline cursor-pointer mr-3 transition-colors duration-200">
                                继续阅读
                            </button>
                            <button v-if="activeTab === 'shelf'" @click="removeFromShelf(item.book.id, $event)" class="text-red-500 hover:text-red-700 font-medium underline cursor-pointer transition-colors duration-200">
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
import { ref, computed, watch, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';
import { shelfAPI, historyAPI } from '@/api/services';

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

// ———————— API数据：最近阅读 ————————
const recentBooks = ref([]);

// ———————— API数据：我的书架 ————————
const shelfBooks = ref([]);

// 计算当前显示的书籍
const displayedBooks = computed(() => {
    return activeTab.value === 'recent' ? recentBooks.value : shelfBooks.value
});

// 导航到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`);
};

// 移除书架
const removeFromShelf = async (bookId, event) => {
    event.stopPropagation();
    try {
        // 临时使用默认userId=1
        await shelfAPI.removeFromShelf(1, bookId);
        shelfBooks.value = shelfBooks.value.filter(item => item.bookId !== bookId);
    } catch (error) {
        console.error('移除书架失败:', error);
    }
};

// 获取数据
onMounted(async () => {
    try {
        // 临时使用默认userId=1
        // 获取最近阅读数据
        const recentResponse = await historyAPI.getHistory(1);
        recentBooks.value = recentResponse.data;
        
        // 获取书架数据
        const shelfResponse = await shelfAPI.getShelf(1);
        shelfBooks.value = shelfResponse.data;
    } catch (error) {
        console.error('获取数据失败:', error);
    }
});
</script>