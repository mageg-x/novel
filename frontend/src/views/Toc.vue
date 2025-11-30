<template>
    <!-- 小屏幕 -->
    <div
        :class="['block md:hidden min-h-screen', currentTheme === 'night' ? 'night' : '', currentTheme === 'eye-protect' ? 'eye-protect' : '']">
        <!-- 顶部导航栏 -->
        <header class="sticky top-0 z-50 bg-[#469b75] text-white shadow-md">
            <div class="flex items-center justify-between px-4 py-3">
                <button @click="goBack" class="text-white text-xl cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <h1 class="text-base font-bold truncate mx-2">{{ chapterTitle }}</h1>
                <button @click="goHome" class="text-white text-xl cursor-pointer">
                    <i class="fas fa-home"></i>
                </button>
            </div>
        </header>

        <!-- 阅读控制栏 -->
        <div :class="[
            'sticky top-[56px] z-40 flex items-center justify-between px-4 py-2 theme-border transition-transform duration-300',
            showControls ? 'translate-y-0' : '-translate-y-full',
            'theme-control-bg'
        ]">
            <div class="flex items-center gap-2">
                <button @click="toggleNightMode"
                    :class="['px-3 py-1 rounded-full text-sm font-medium shadow-sm transition-all',
                        currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white' : 'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500']">
                    关灯
                </button>
                <button @click="toggleEyeProtect"
                    :class="['px-3 py-1 rounded-full text-sm font-medium shadow-sm transition-all',
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white' : 'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500']">
                    护眼
                </button>
            </div>
            <div class="flex items-center gap-2">
                <span class="text-sm theme-text">字体：</span>
                <button class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                            'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200'">大</button>
                <button class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]'">中</button>
                <button class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                            'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200'">小</button>
            </div>
            <button class="w-10 h-10 flex items-center justify-center rounded-full shadow-sm transition-all"
                :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white' :
                    currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-400 to-green-500 text-white' :
                        'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500'">
                <i class="fas fa-volume-up text-lg"></i>
            </button>
        </div>

        <!-- 直达底部按钮 -->
        <div class="fixed right-4 bottom-20 z-30">
            <button @click="scrollToBottom"
                class="w-12 h-12 bg-[#469b75] text-white rounded-full shadow-lg flex items-center justify-center hover:bg-[#3d8766] transition-colors"
                aria-label="直达底部">
                <i class="fas fa-arrow-down text-lg"></i>
            </button>
        </div>

        <!-- 章节列表 -->
        <div class="p-4 ">
            <div :class="['rounded-lg shadow-md p-4 mb-4',
                currentTheme === 'night' ? 'bg-[#2d2d2d] theme-border' :
                    currentTheme === 'eye-protect' ? 'bg-[#b3e5cc] theme-border' : 'bg-white']">
                <h2 :class="['text-lg font-semibold mb-4',
                    currentTheme === 'night' ? 'text-white' :currentTheme === 'eye-protect' ? 'text-[#333333]' : 'text-gray-500']">
                    目录：
                </h2>
                <div class="space-y-2">
                    <a v-for="(chapter, pageIndex) in currentPageChapters" :key="pageIndex" @click="goToChapter((currentPage - 1) * pageSize + pageIndex)" :class="['flex items-center justify-between cursor-pointer px-3 py-3 rounded transition-colors border',
                        currentTheme === 'night' ?
                            'text-white hover:text-emerald-400 hover:bg-[#3d3d3d] border-[#404040]' :
                            currentTheme === 'eye-protect' ?
                                'text-[#333333] hover:text-[#3d8766] hover:bg-[#a5d6a7] border-[#81c784]' :
                                'text-gray-700 hover:text-emerald-600 hover:bg-gray-50 border-gray-100']">
                        <span class="flex-grow truncate pr-2">第{{ numToChinese(chapter.chapterNo) }}章 {{ chapter.title }}</span>
                        <span
                            :class="['text-sm min-w-12 text-right', chapter.isVip ? 'text-red-500' : 'text-emerald-600']">
                            {{ chapter.isVip ? '[VIP]' : '[免费]' }}
                        </span>
                    </a>
                </div>
                
                <!-- 小屏幕分页控件 -->
                <div v-if="totalPages > 1" class="mt-6 flex justify-center items-center space-x-2">
                    <button @click="prevPage" :disabled="currentPage === 1" 
                        :class="['px-4 py-2 rounded-full text-sm', 
                            currentTheme === 'night' ? 
                                'bg-gray-800 text-white hover:bg-gray-700 disabled:opacity-50' :
                                currentTheme === 'eye-protect' ?
                                    'bg-green-100 text-green-800 hover:bg-green-200 disabled:opacity-50' :
                                    'bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50']">
                        上一页
                    </button>
                    <span :class="['px-4 py-2 rounded-full text-sm font-medium', 
                        currentTheme === 'night' ? 'bg-[#469b75] text-white' :
                        currentTheme === 'eye-protect' ? 'bg-green-500 text-white' :
                        'bg-[#469b75] text-white']">
                        {{ currentPage }}/{{ totalPages }}
                    </span>
                    <button @click="nextPage" :disabled="currentPage === totalPages" 
                        :class="['px-4 py-2 rounded-full text-sm', 
                            currentTheme === 'night' ? 
                                'bg-gray-800 text-white hover:bg-gray-700 disabled:opacity-50' :
                                currentTheme === 'eye-protect' ?
                                    'bg-green-100 text-green-800 hover:bg-green-200 disabled:opacity-50' :
                                    'bg-gray-100 text-gray-700 hover:bg-gray-200 disabled:opacity-50']">
                        下一页
                    </button>
                </div>
            </div>
        </div>

        <!-- 底部导航栏 -->
        <!-- <ToolBar :showControls="showControls" /> -->
    </div>

    <!-- 大屏 -->
    <div class="hidden md:block container md:w-5xl  mx-auto ">
        <div class="hidden md:block">
            <Header />
        </div>
        <!-- Breadcrumb -->
        <div class="bg-gray-100 ">
            <div class="max-w-7xl mx-auto px-4 py-3 text-sm text-gray-600">
                <a href="#" class="hover:text-emerald-600">阁林小说</a>
                <span class="mx-2">></span>
                <a href="#" class="hover:text-emerald-600">{{ bookData.category }}</a>
                <span class="mx-2">></span>
                <a href="#" class="hover:text-emerald-600">{{ bookData.title }}</a>
                <span class="mx-2">></span>
                <span>作品目录</span>
            </div>
        </div>

        <!-- Main Content -->
        <main class="max-w-7xl mx-auto px-4 py-8 bg-white ">
            <h1 class="text-3xl font-bold text-center text-gray-800 mb-6">
                {{ bookData.title }}
            </h1>

            <div class="flex justify-center gap-8 text-sm text-gray-600 mb-8">
                <span>作者：{{ bookData.author }}</span>
                <span>分类：{{ bookData.category }}</span>
                <span>状态：{{ bookData.status }}</span>
                <span>总点击：{{ bookData.clickCount }}</span>
                <span>总字数：{{ bookData.wordCount }}</span>
            </div>

            <h2 class="text-xl font-semibold text-gray-500 bg-gray-100 mb-6">目录：</h2>

            <div class="grid grid-cols-1 md:grid-cols-2  gap-x-6 gap-y-3">
                <a v-for="(chapter, pageIndex) in currentPageChapters" :key="pageIndex" @click="goToChapter((currentPage - 1) * pageSize + pageIndex)" :class="['flex items-center justify-between cursor-pointer pl-0 pr-3 py-2 rounded transition-colors',
                    currentTheme === 'night' ?
                        'text-white hover:text-emerald-400 hover:bg-gray-800' :
                        currentTheme === 'eye-protect' ?
                            'text-[#333333] hover:text-[#3d8766] hover:bg-[#a5d6a7]' :
                            'text-gray-700 hover:text-emerald-600 hover:bg-gray-50']">
                    <span class="flex-grow truncate  pr-2">第{{ numToChinese(chapter.chapterNo) }}章 {{ chapter.title }}</span>
                    <span :class="['text-sm min-w-16 text-right', chapter.isVip ? 'text-red-500' : 'text-emerald-600']">
                        {{ chapter.isVip ? '[VIP]' : '[免费]' }}
                    </span>
                </a>
            </div>
            
            <!-- 大屏幕分页控件 -->
            <div v-if="totalPages > 1" class="mt-8 flex justify-center items-center space-x-4">
                <button @click="prevPage" :disabled="currentPage === 1" 
                    class="px-6 py-2 rounded-md shadow-sm text-sm font-medium transition-colors bg-[#469b75] text-white hover:bg-[#3d8766] disabled:opacity-50">
                    上一页
                </button>
                
                <!-- 页码按钮，只显示当前页前后各3页 -->
                <template v-for="page in visiblePages" :key="page">
                    <button v-if="page === '...'" disabled 
                        class="px-4 py-2 rounded-md text-sm text-gray-500">
                        {{ page }}
                    </button>
                    <button v-else @click="goToPage(page)" 
                        :class="['px-4 py-2 rounded-md shadow-sm text-sm font-medium transition-colors', 
                            page === currentPage ? 
                                'bg-[#469b75] text-white' : 
                                'bg-white text-gray-700 hover:bg-gray-50']">
                        {{ page }}
                    </button>
                </template>
                
                <button @click="nextPage" :disabled="currentPage === totalPages" 
                    class="px-6 py-2 rounded-md shadow-sm text-sm font-medium transition-colors bg-[#469b75] text-white hover:bg-[#3d8766] disabled:opacity-50">
                    下一页
                </button>
                
                <div class="ml-4 text-sm text-gray-600">
                    共 {{ totalPages }} 页，当前第 {{ currentPage }} 页
                </div>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';
import { bookAPI } from '@/api/services';
import { numToChinese } from '../utils/tiny';

const route = useRoute();
const router = useRouter();

// 小屏幕相关状态
const currentTheme = ref('default');
const chapterTitle = ref('');
const showControls = ref(true);
const lastScrollTop = ref(0);
const scrollThreshold = 10;

const breadcrumb = ref({
    home: '阁林小说',
    category: ''
});
// 书籍数据
const bookData = ref({});
const chapters = ref([]);
// 分页相关状态
const currentPage = ref(1);
const pageSize = ref(60);

// 计算当前页显示的章节
const currentPageChapters = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    return chapters.value.slice(start, end);
});

// 计算总页数
const totalPages = computed(() => {
    const total = Math.ceil(chapters.value.length / pageSize.value);
    console.log('总页数:', total);
    return total;
});

// 计算可见页码，显示当前页前后各3页
const visiblePages = computed(() => {
    const pages = [];
    const total = totalPages.value;
    const current = currentPage.value;
    
    if (total <= 7) {
        // 总页数小于等于7，显示所有页码
        for (let i = 1; i <= total; i++) {
            pages.push(i);
        }
    } else {
        // 总页数大于7，显示当前页前后各3页，其余用省略号
        if (current <= 4) {
            // 当前页在1-4之间，显示1-7页
            for (let i = 1; i <= 5; i++) {
                pages.push(i);
            }
            pages.push('...');
            pages.push(total);
        } else if (current >= total - 3) {
            // 当前页在最后4页，显示最后7页
            pages.push(1);
            pages.push('...');
            for (let i = total - 4; i <= total; i++) {
                pages.push(i);
            }
        } else {
            // 当前页在中间，显示前后各3页
            pages.push(1);
            pages.push('...');
            for (let i = current - 2; i <= current + 2; i++) {
                pages.push(i);
            }
            pages.push('...');
            pages.push(total);
        }
    }
    
    return pages;
});

// 切换到指定页
const goToPage = (page) => {
    if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page;
        // 滚动到顶部
        window.scrollTo({ top: 0, behavior: 'smooth' });
    }
};

// 上一页
const prevPage = () => {
    goToPage(currentPage.value - 1);
};

// 下一页
const nextPage = () => {
    goToPage(currentPage.value + 1);
};

// 小屏幕相关方法
const goBack = () => {
    router.back();
};

const goHome = () => {
    router.push('/');
};

const toggleNightMode = () => {
    currentTheme.value = currentTheme.value === 'night' ? 'default' : 'night';
};

const toggleEyeProtect = () => {
    currentTheme.value = currentTheme.value === 'eye-protect' ? 'default' : 'eye-protect';
};

const scrollToBottom = () => {
    window.scrollTo({
        top: document.body.scrollHeight,
        behavior: 'smooth'
    });
};

// 跳转到章节阅读页面
const goToChapter = (index) => {
    // 使用章节的实际ID，不再构建模拟的chapterId
    const chapter = chapters.value[index];
    router.push(`/book/${bookData.value.id}/${chapter.chapterNo}`);
};

// 滚动监听处理函数
const handleScroll = () => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;

    // 只有当滚动超过阈值时才触发显示/隐藏逻辑
    if (Math.abs(scrollTop - lastScrollTop.value) > scrollThreshold) {
        // 向上滚动时显示控制栏，向下滚动时隐藏
        showControls.value = scrollTop < lastScrollTop.value;
        lastScrollTop.value = scrollTop;
    }
};

// 获取书籍数据
onMounted(async () => {
    // 添加滚动监听
    window.addEventListener('scroll', handleScroll);
    // 获取路由参数
    const bookId = route.params.bookId;

    try {
        // 从API获取书籍详情
        const bookResponse = await bookAPI.getById(bookId);
        bookData.value = bookResponse.data;

        // 设置小屏幕导航栏标题
        chapterTitle.value = bookData.value.title;

        // 从API获取章节列表
        const chaptersResponse = await bookAPI.getChapters(bookId);
        chapters.value = chaptersResponse?.data?.chapters || [];
    } catch (error) {
        console.error('获取书籍或章节数据失败:', error);
    }
});

// 监听路由参数变化，重新加载数据
watch(
    () => route.params.bookId,
    (newId, oldId) => {
        if (newId !== oldId) {
            fetchBookData();
        }
    }
);

// 封装数据获取函数，便于复用
const fetchBookData = async () => {
    const bookId = route.params.bookId;
    try {
        // 先重置状态
        chapters.value = [];
        currentPage.value = 1;
        
        // 从API获取书籍详情
        const bookResponse = await bookAPI.getById(bookId);
        bookData.value = bookResponse.data;

        // 设置小屏幕导航栏标题
        chapterTitle.value = bookData.value.title;

        // 从API获取章节列表
        const chaptersResponse = await bookAPI.getChapters(bookId);
        // 使用响应式的方式更新chapters数组
        chapters.value = chaptersResponse?.data?.chapters || [];
        
        console.log('章节数据已加载，共', chapters.value.length, '章');
        console.log('计算得到的总页数:', totalPages.value);
    } catch (error) {
        console.error('获取书籍或章节数据失败:', error);
    }
};

// 组件卸载时移除滚动监听
onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
});
</script>

<style scoped>
/* 可选：添加局部样式 */
</style>