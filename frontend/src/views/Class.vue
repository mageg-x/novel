<template>
    <!-- 小屏幕 -->
    <div class="block md:hidden">
        <!-- 顶部导航 -->
        <div class="sticky top-0 z-50 bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-lg">
            <div class="flex items-center justify-between h-14 px-4">
                <button @click="goBack"
                    class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <span class="text-lg font-semibold">全部作品</span>
                <button @click="goHome"
                    class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
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
                    @click="showFilters = !showFilters">
                    <span class="flex items-center">
                        <i class="fas fa-filter mr-2 text-emerald-500"></i>
                        筛选条件
                    </span>
                    <i
                        :class="['fas', showFilters ? 'fa-chevron-up' : 'fa-chevron-down', 'text-gray-400 transition-transform duration-200']"></i>
                </button>

                <!-- 已选筛选条件 -->
                <div v-if="!showFilters && getSelectedFilters().length > 0" class="px-4 pb-3 overflow-x-auto">
                    <div class="flex flex-wrap gap-2">
                        <div v-for="(filter, idx) in getSelectedFilters()" :key="idx"
                            class="flex items-center bg-emerald-50 text-emerald-700 text-xs px-3 py-1 rounded-full shadow-sm">
                            <span class="mr-1">{{ filter.text }}</span>
                            <i class="fas fa-times text-xs cursor-pointer hover:text-emerald-900 transition-colors"
                                @click="clearFilter(filter)"></i>
                        </div>
                    </div>
                </div>

                <!-- 筛选面板内容 -->
                <div v-if="showFilters" class="p-4 space-y-4 max-h-80 overflow-y-auto">
                    <!-- 作品分类筛选 -->
                    <div>
                        <h3 class="text-sm font-semibold text-gray-700 mb-2">{{ filters[1].label }}</h3>
                        <div class="flex overflow-x-auto pb-2 -mx-2 px-2 space-x-2">
                            <button v-for="(option, optIdx) in filters[1].options" :key="optIdx" :class="[
                                'flex-shrink-0 px-4 py-2 rounded-full text-sm transition-all duration-200',
                                option.active
                                    ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md'
                                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                            ]" @click="selectFilter(filters[1], option)">
                                {{ option.text }}
                            </button>
                        </div>
                    </div>

                    <!-- 其他筛选条件 -->
                    <div class="space-y-3">
                        <div v-for="(filter, idx) in filters.slice(0, 1).concat(filters.slice(2))" :key="idx">
                            <h3 class="text-sm font-semibold text-gray-700 mb-2">{{ filter.label }}</h3>
                            <div class="grid grid-cols-2 gap-2">
                                <button v-for="(option, optIdx) in filter.options" :key="optIdx" :class="[
                                    'text-left px-3 py-2 rounded-md text-sm transition-all duration-200',
                                    option.active
                                        ? 'bg-emerald-100 text-emerald-700 font-medium'
                                        : 'bg-gray-100 text-gray-600 hover:bg-gray-200'
                                ]" @click="selectFilter(filter, option)">
                                    {{ option.text }}
                                </button>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 书籍列表 -->
            <div class="container mx-auto px-4">
                <div v-for="(book, index) in filteredBooks" :key="book.id"
                    class="bg-white rounded-2xl shadow-md p-4 mb-3 transition-all duration-300 hover:shadow-lg hover:-translate-y-1 cursor-pointer"
                    @click="navigateToBook(book.id)">
                    <!-- 排名与书籍信息 -->
                    <div class="flex items-center">
                        <!-- 排名 -->
                        <div class="mr-3">
                            <span :class="{
                                'bg-red-500': index === 0,
                                'bg-orange-500': index === 1,
                                'bg-yellow-500': index === 2,
                                'bg-blue-500': index > 2 && index < 10,
                                'bg-gray-300': index >= 10
                            }"
                                class="flex items-center justify-center w-9 h-9 text-white rounded-full text-base font-bold shadow-sm">
                                {{ index + 1 }}
                            </span>
                        </div>

                        <!-- 书籍信息 -->
                        <div class="flex-1">
                            <div class="flex items-center justify-between mb-1">
                                <div class="flex items-center">
                                    <span
                                        class="text-xs px-2 py-0.5 bg-emerald-100 text-emerald-700 rounded-full font-medium mr-2">
                                        [{{ book.category }}]
                                    </span>
                                    <span class="text-gray-500 text-xs font-medium">{{ book.author }}</span>
                                </div>
                                <span class="text-sm text-gray-600 font-medium ml-2 whitespace-nowrap">
                                    {{ (book.wordCount / 10000).toFixed(2) }}万
                                </span>
                            </div>

                            <h3
                                class="text-base font-semibold text-gray-800 mb-1 line-clamp-1 hover:text-emerald-600 transition-colors">
                                {{ book.title }}
                            </h3>

                            <p class="text-xs text-gray-500 line-clamp-1">{{ book.latestChapterTitle }}</p>
                        </div>
                    </div>
                </div>
            </div>

            <!-- 分页 -->
            <div class="flex justify-center items-center space-x-1 py-8 px-4">
                <button @click="goToPrevPage"
                    class="px-3 py-2 text-gray-500 hover:text-emerald-500 rounded-md hover:bg-gray-100 transition-colors"
                    :disabled="currentPage === 1">
                    <i class="fas fa-chevron-left"></i>
                </button>

                <button v-for="page in Math.min(totalPages, 5)" :key="page" :class="[
                    'px-3 py-2 rounded-md text-sm transition-all duration-200',
                    page === currentPage
                        ? 'bg-gradient-to-r from-emerald-500 to-teal-500 text-white shadow-md'
                        : 'text-gray-600 hover:bg-gray-100'
                ]" @click="goToPage(page)">
                    {{ page }}
                </button>

                <button @click="goToNextPage"
                    class="px-3 py-2 text-gray-500 hover:text-emerald-500 rounded-md hover:bg-gray-100 transition-colors"
                    :disabled="currentPage === totalPages">
                    <i class="fas fa-chevron-right"></i>
                </button>
            </div>
        </div>

        <!-- 底部导航 -->
        <!-- <ToolBar :activeTab="activeTab" @tab-change="handleTabChange" /> -->
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
                    <div class="leading-[100%] w-20 text-gray-500 pt-1">{{ filter.label }}：</div>
                    <div class="flex-1 flex flex-wrap gap-3">
                        <a v-for="(option, optIdx) in filter.options" :key="optIdx" href="#"
                            @click.prevent="selectFilter(filter, option)" :class="[
                                option.active
                                    ? 'bg-green-100 text-green-700 hover:bg-green-200 px-3 py-1 rounded-full'
                                    : 'text-gray-600 hover:text-green-500 hover:bg-gray-100 px-3 py-1 rounded-full'
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
                        <tr v-for="(book, index) in filteredBooks" :key="book.id" class="hover:bg-gray-50 border-0">
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
                            <td class="px-4 py-3 text-gray-600">{{ book.category }}</td>
                            <td class="px-4 py-3">
                                <router-link :to="'/book/' + book.id"
                                    class="text-gray-800 hover:text-green-500 cursor-pointer">{{ book.title
                                    }}</router-link>
                            </td>
                            <td class="px-4 py-3 text-gray-600">{{ book.latestChapterTitle }}</td>
                            <td class="px-4 py-3 text-gray-600">{{ book.author }}</td>
                            <td class="px-4 py-3 text-gray-600">{{ (book.wordCount / 10000).toFixed(2) }}万</td>
                        </tr>
                    </tbody>
                </table>
            </div>

            <!-- Pagination -->
            <div class="flex justify-center items-center space-x-2 py-8 text-sm">
                <button @click="goToPrevPage" class="px-3 py-1 text-gray-500 hover:text-green-500"
                    :disabled="currentPage === 1">上一页</button>

                <button v-for="page in Math.min(totalPages, 5)" :key="page" @click="goToPage(page)" :class="[
                    page === currentPage
                        ? 'px-3 py-1 bg-green-500 text-white rounded'
                        : 'px-3 py-1 text-gray-600 hover:bg-gray-100 rounded'
                ]">
                    {{ page }}
                </button>

                <button @click="goToNextPage" class="px-3 py-1 text-gray-500 hover:text-green-500"
                    :disabled="currentPage === totalPages">下一页</button>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import Header from '@/components/Header.vue';
import ToolBar from '@/components/ToolBar.vue';
import { bookAPI } from '@/api/services';

// 路由和导航
const router = useRouter()

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

// 状态管理
const showFilters = ref(false);
const currentPage = ref(1);
const pageSize = ref(20);
const books = ref([]);

// 分页相关
const totalPages = computed(() => {
    return Math.ceil(filteredBooksWithoutPagination.value.length / pageSize.value);
});

// 分页控制函数
const goToPage = (page) => {
    if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page;
    }
};

const goToPrevPage = () => {
    if (currentPage.value > 1) {
        currentPage.value--;
    }
};

const goToNextPage = () => {
    if (currentPage.value < totalPages.value) {
        currentPage.value++;
    }
};

// 筛选条件
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

// 筛选条件操作
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

const clearFilter = (filterToClear) => {
    filters.value.forEach(filter => {
        const optionIndex = filter.options.findIndex(option => option === filterToClear);
        if (optionIndex !== -1) {
            filter.options[optionIndex].active = false;
            const unlimitedOption = filter.options.find(option => option.text === '不限');
            if (unlimitedOption) {
                unlimitedOption.active = true;
            }
        }
    });
};

const selectFilter = (filter, option) => {
    filter.options.forEach(opt => {
        opt.active = false;
    });
    option.active = true;
    currentPage.value = 1;
};

// 数据获取
const fetchBooks = async () => {
    try {
        const response = await bookAPI.getAll();
        books.value = response.data.books;
        currentPage.value = 1;
    } catch (error) {
        console.error('获取书籍列表失败:', error);
    }
};

// 计算属性
const filteredBooksWithoutPagination = computed(() => {
    let result = [...books.value];

    // 获取所有选中的筛选条件
    const selectedFilters = {};
    filters.value.forEach(filter => {
        selectedFilters[filter.label] = filter.options.find(option => option.active && option.text !== '不限');
    });

    // 作品分类筛选
    if (selectedFilters['作品分类']) {
        result = result.filter(book => book.category === selectedFilters['作品分类'].text);
    }

    // 是否完结筛选
    if (selectedFilters['是否完结']) {
        const statusText = selectedFilters['是否完结'].text;
        result = result.filter(book => {
            if (statusText === '连载中') {
                return book.status === 'serializing';
            } else if (statusText === '已完结') {
                return book.status === 'completed';
            }
            return true;
        });
    }

    // 作品字数筛选
    if (selectedFilters['作品字数']) {
        const wordCountOption = selectedFilters['作品字数'].text;
        result = result.filter(book => {
            const wordCount = book.wordCount;
            if (wordCountOption === '30万字以下') {
                return wordCount < 300000;
            } else if (wordCountOption === '30-50万字') {
                return wordCount >= 300000 && wordCount < 500000;
            } else if (wordCountOption === '50-100万字') {
                return wordCount >= 500000 && wordCount < 1000000;
            } else if (wordCountOption === '100万字以上') {
                return wordCount >= 1000000;
            }
            return true;
        });
    }

    // 排序方式筛选
    if (selectedFilters['排序方式']) {
        const sortOption = selectedFilters['排序方式'].text;
        switch (sortOption) {
            case '更新时间':
                result.sort((a, b) => new Date(b.updatedAt) - new Date(a.updatedAt));
                break;
            case '总字数':
                result.sort((a, b) => b.wordCount - a.wordCount);
                break;
            case '点击量':
                result.sort((a, b) => b.clickCount - a.clickCount);
                break;
        }
    }

    return result;
});

const filteredBooks = computed(() => {
    const start = (currentPage.value - 1) * pageSize.value;
    const end = start + pageSize.value;
    return filteredBooksWithoutPagination.value.slice(start, end);
});

// 生命周期
onMounted(() => {
    fetchBooks();
});

onUnmounted(() => {
    // 清理工作
});
</script>

<style scoped>
/* 如需局部样式可在此添加 */
</style>