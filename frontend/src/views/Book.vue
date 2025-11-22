<template>
    <div class="hidden md:block">
        <Header />
    </div>

    <!-- 小屏幕 -->
    <div
        :class="['block md:hidden min-h-screen w-full overflow-x-hidden', currentTheme === 'night' ? 'night' : '', currentTheme === 'eye-protect' ? 'eye-protect' : '']">
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
        <div v-if="showControls" class="flex items-center justify-between px-4 pt-2 theme-border theme-control-bg">
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
                <button @click="setFontSize('large')"
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="[fontSize === 'large' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                        currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">大</button>
                <button @click="setFontSize('medium')"
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="[fontSize === 'medium' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                        currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">中</button>
                <button @click="setFontSize('small')"
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="[fontSize === 'small' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                        currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">小</button>
            </div>
            <div class="flex items-center gap-2">
                <button @click="fontType = fontType === 'handwriting' ? 'default' : 'handwriting'"
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="[fontType === 'handwriting' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                        currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">
                    手写
                </button>
                <button @click="showPinyin = !showPinyin"
                    class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                    :class="[showPinyin ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                            'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                        currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">
                    注音
                </button>
            </div>
            <button class="w-10 h-10 flex items-center justify-center rounded-full shadow-sm transition-all"
                :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white' :
                    currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-400 to-green-500 text-white' :
                        'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500'">
                <i class="fas fa-volume-up text-lg"></i>
            </button>
        </div>

        <!-- 阅读区域 -->
        <main class="p-4 overflow-x-hidden"
            :class="showControls ? 'min-h-[calc(100vh-224px)]' : 'min-h-[calc(100vh-124px)]'">
            <div class="theme-text leading-relaxed break-words"
                :style="{ fontSize: fontSize === 'large' ? '1.5rem' : fontSize === 'small' ? '0.875rem' : '1rem', lineHeight: fontSize === 'large' ? '2rem' : '1.8rem', fontFamily: fontType === 'handwriting' ? 'Yozai, Microsoft YaHei, sans-serif' : 'Microsoft YaHei, sans-serif' }">
                <p v-for="(paragraph, index) in chapterContent" :key="index" class="mb-4" v-html="addPinyin(paragraph)">
                </p>
            </div>
        </main>

        <!-- 章节导航 -->
        <div class="mt-4 mb-8">
            <div class="flex justify-around items-center px-4">
                <button @click="previousChapter" :disabled="!hasPreviousChapter"
                    class="text-sm theme-text disabled:text-gray-500">
                    上一章
                </button>
                <button @click="goToc" class="text-sm theme-text">
                    目录
                </button>
                <button @click="nextChapter" :disabled="!hasNextChapter"
                    class="text-sm theme-text disabled:text-gray-500">
                    下一章
                </button>
            </div>
        </div>

        <!-- 底部导航栏 -->
        <!-- <ToolBar :showControls="showControls" /> -->
    </div>

    <!-- 大屏 -->
    <div
        :class="['hidden md:block relative min-h-screen w-5xl my-4 mx-auto', currentTheme === 'night' ? 'night' : currentTheme === 'eye-protect' ? 'eye-protect' : 'bg-[#ebe5d8]']">
        <!-- 面包屑导航 -->
        <div class="max-w-7xl mx-auto px-4 py-3 text-sm theme-text">
            <a href="/" class="hover:text-[#469b75]">首页</a>
            <span class="mx-2">></span>
            <a :href="`/class?category=${bookData.category}`" class="hover:text-[#469b75]">{{ bookData.category }}</a>
            <span class="mx-2">></span>
            <a :href="`/book/${bookId}`" class="hover:text-[#469b75]">{{ bookData.title }}</a>
        </div>

        <!-- 阅读控制栏 -->
        <div v-if="showControls" class=" w-fit mx-auto px-4 pt-2 pb-4 theme-border theme-control-bg">
            <div class="flex items-center justify-between gap-4">
                <div class="flex items-center gap-4">
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
                <span class="flex items-center gap-4 theme-text">字体：</span>
                <div class="flex items-center gap-4">
                    <button @click="setFontSize('large')"
                        class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                        :class="[fontSize === 'large' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                                'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                            currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                                currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                    'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">大</button>
                    <button @click="setFontSize('medium')"
                        class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                        :class="[fontSize === 'medium' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                                'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                            currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                                currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                    'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">中</button>
                    <button @click="setFontSize('small')"
                        class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                        :class="[fontSize === 'small' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                                'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                            currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                                currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                    'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">小</button>
                </div>
                <div class="flex items-center gap-4">
                    <button @click="fontType = fontType === 'handwriting' ? 'default' : 'handwriting'"
                        class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                        :class="[fontType === 'handwriting' ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                                'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                            currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                                currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                    'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">
                        手写
                    </button>
                    <button @click="showPinyin = !showPinyin"
                        class="px-2 py-1 rounded-full text-sm font-medium shadow-sm transition-all"
                        :class="[showPinyin ? (currentTheme === 'night' ? 'bg-gradient-to-r from-purple-600 to-purple-700 text-white border border-purple-500' :
                            currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-500 to-green-600 text-white border border-green-400' :
                                'bg-gradient-to-r from-[#469b75] to-[#3d8766] text-white border border-[#3d8766]') :
                            currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white border border-gray-600' :
                                currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-100 to-green-200 text-green-800 border border-green-300' :
                                    'bg-gradient-to-r from-white to-gray-100 text-gray-700 border border-gray-200 hover:from-gray-100 hover:to-gray-200']">
                        注音
                    </button>
                </div>
                <button class="w-10 h-10 flex items-center justify-center rounded-full shadow-sm transition-all"
                    :class="currentTheme === 'night' ? 'bg-gradient-to-r from-gray-700 to-gray-800 text-white' :
                        currentTheme === 'eye-protect' ? 'bg-gradient-to-r from-green-400 to-green-500 text-white' :
                            'bg-gradient-to-r from-gray-300 to-gray-400 text-gray-700 hover:from-gray-400 hover:to-gray-500'">
                    <i class="fas fa-volume-up text-lg"></i>
                </button>
            </div>
        </div>

        <!-- 阅读区域 -->
        <main class="w-full pt-6 ">
            <div class="mx-auto w-fit  ">
                <!-- 章节标题 -->
                <h2 class="text-xl font-bold theme-text text-center mb-6">{{ chapterTitle }}</h2>

                <!-- 章节内容 -->
                <div class="theme-text leading-relaxed break-words px-16"
                    :style="{ fontSize: fontSize === 'large' ? '1.5rem' : fontSize === 'small' ? '0.875rem' : '1rem', lineHeight: fontSize === 'large' ? '2rem' : '1.8rem', fontFamily: fontType === 'handwriting' ? 'Yozai, Microsoft YaHei, sans-serif' : 'Microsoft YaHei, sans-serif' }">
                    <p v-for="(paragraph, index) in chapterContent" :key="index" class="mb-4"
                        v-html="addPinyin(paragraph)">
                    </p>
                </div>


            </div>

            <!-- 章节导航 -->
            <div class=" absolute w-full bottom-2">
                <div class=" w-fit mx-auto flex justify-between mt-10 pt-6 gap-8 ">
                    <button @click="previousChapter" :disabled="!hasPreviousChapter"
                        :class="['px-4 py-2 rounded-md transition-colors disabled:cursor-not-allowed',
                            currentTheme === 'night' ? 'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-700 disabled:text-gray-500' :
                                currentTheme === 'eye-protect' ? 'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-300 disabled:text-gray-500' :
                                    'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-300 disabled:text-gray-500']">
                        上一章
                    </button>

                    <button @click="goToc"
                        :class="['px-4 py-2 rounded-md transition-colors disabled:cursor-not-allowed',
                            currentTheme === 'night' ? 'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-700 disabled:text-gray-500' :
                                currentTheme === 'eye-protect' ? 'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-300 disabled:text-gray-500' :
                                    'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-300 disabled:text-gray-500']">
                        目录
                    </button>

                    <button @click="nextChapter" :disabled="!hasNextChapter"
                        :class="['px-4 py-2 rounded-md transition-colors disabled:cursor-not-allowed',
                            currentTheme === 'night' ? 'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-700 disabled:text-gray-500' :
                                currentTheme === 'eye-protect' ? 'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-300 disabled:text-gray-500' :
                                    'bg-[#469b75] hover:bg-[#3d8766] text-white disabled:bg-gray-300 disabled:text-gray-500']">
                        下一章
                    </button>
                </div>
            </div>

        </main>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import ToolBar from '@/components/ToolBar.vue'
import { bookAPI } from '@/api/services'

// 路由和导航
const route = useRoute()
const router = useRouter()
const bookId = route.params.bookId
const chapterId = route.params.chapterId

// 状态管理
const currentTheme = ref('') // '' 表示默认主题, 'night' 表示深色主题, 'eye-protect' 表示护眼主题
const showControls = ref(false)
const lastScrollTop = ref(0)
const scrollThreshold = 50
const fontSize = ref('medium') // 'small', 'medium', 'large'
const fontType = ref('default') // 'default', 'handwriting'
const showPinyin = ref(false)

// 数据
const bookData = ref({
    title: '',
    chapters: []
})

const currentChapter = ref({
    id: '',
    title: '',
    content: []
})

const loading = ref(true)
const error = ref(null)

// 计算属性
const bookTitle = computed(() => bookData.value.title)
const chapterTitle = computed(() => currentChapter.value.title)
const chapterContent = computed(() => {
    const content = currentChapter.value.content
    if (!content || typeof content !== 'string') return []
    return content.split(/\n\n+/).filter(paragraph => paragraph.trim())
})
const hasPreviousChapter = computed(() => {
    const currentIndex = bookData.value.chapters.findIndex(chapter => chapter.chapterId === chapterId)
    return currentIndex > 0
})
const hasNextChapter = computed(() => {
    const currentIndex = bookData.value.chapters.findIndex(chapter => chapter.chapterId === chapterId)
    return currentIndex < bookData.value.chapters.length - 1
})

// 字体和主题控制
const setFontSize = (size) => {
    fontSize.value = size;
};

const toggleNightMode = () => {
    if (currentTheme.value === 'night') {
        currentTheme.value = ''
    } else {
        currentTheme.value = 'night'
        document.body.classList.remove('eye-protect')
    }
}

const toggleEyeProtect = () => {
    if (currentTheme.value === 'eye-protect') {
        currentTheme.value = ''
    } else {
        currentTheme.value = 'eye-protect'
        document.body.classList.remove('night')
    }
}

// 注音功能
import pinyin from 'pinyin'

const addPinyin = (text) => {
    if (!showPinyin.value) {
        return text
    }

    return text.replace(/[\u4e00-\u9fa5]/g, (char) => {
        const pinyinArray = pinyin(char, { style: pinyin.STYLE_TONE })
        const pinyinText = pinyinArray[0][0]
        return `<ruby>${char}<rt>${pinyinText}</rt></ruby>`
    })
}

// 导航方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}

const goToc = () => {
    router.push(`/book/${bookId}/toc`)
}

const previousChapter = () => {
    if (hasPreviousChapter.value) {
        const currentIndex = bookData.value.chapters.findIndex(chapter => chapter.chapterId === chapterId)
        const previousChapter = bookData.value.chapters[currentIndex - 1]
        router.push(`/book/${bookId}/${previousChapter.chapterId}`)
    }
}

const nextChapter = () => {
    if (hasNextChapter.value) {
        const currentIndex = bookData.value.chapters.findIndex(chapter => chapter.chapterId === chapterId)
        const nextChapter = bookData.value.chapters[currentIndex + 1]
        router.push(`/book/${bookId}/${nextChapter.chapterId}`)
    }
}

// 数据获取
const fetchChapterContent = async () => {
    try {
        loading.value = true
        error.value = null

        const [bookInfoResponse, chaptersResponse] = await Promise.all([
            bookAPI.getById(bookId),
            bookAPI.getChapters(bookId)
        ])

        bookData.value = bookInfoResponse.data
        bookData.value.chapters = chaptersResponse.data

        const chapterContentResponse = await bookAPI.getChapter(bookId, chapterId)
        currentChapter.value = chapterContentResponse.data
    } catch (err) {
        console.error('Failed to fetch chapter content:', err)
        error.value = '获取章节内容失败，请稍后重试'
    } finally {
        loading.value = false
    }
}

// 滚动处理
const throttle = (func, delay) => {
    let lastCall = 0;
    return function (...args) {
        const now = Date.now();
        if (now - lastCall >= delay) {
            func.apply(this, args);
            lastCall = now;
        }
    };
};

const handleScroll = throttle(() => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    if (Math.abs(scrollTop - lastScrollTop.value) > scrollThreshold) {
        const newShowControls = scrollTop < lastScrollTop.value;
        if (newShowControls !== showControls.value) {
            showControls.value = newShowControls;
        }
        lastScrollTop.value = scrollTop;
    }
}, 100);

// 生命周期
onMounted(() => {
    fetchChapterContent()
    window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll)
})
</script>

<style>
/* 自定义阅读样式 */
@media (min-width: 768px) {
    main {
        padding-bottom: 100px;
    }
}

/* 注音样式 */
.theme-text ruby {
    display: ruby;
    text-align: center;
    white-space: normal;
    line-height: normal;
    margin: 0 1px;
}

.theme-text ruby>rt {
    font-size: 0.8em !important;
    line-height: 1;
    color: inherit;
    opacity: 0.8;
    text-align: center;
    padding: 0 2px;
    margin-bottom: 4px;
    margin-top: 8px;
    display: ruby-text;
}

.theme-text p {
    line-height: 2.2;
    text-align: justify;
    margin-bottom: 1em;
    padding: 0.2em 0;
}

.theme-text rp {
    display: none;
}
</style>