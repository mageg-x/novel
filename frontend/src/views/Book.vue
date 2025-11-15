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
        <ToolBar :showControls="showControls" />
    </div>

    <!-- 大屏 -->
    <div
        :class="['hidden md:block relative min-h-screen w-5xl my-4 mx-auto', currentTheme === 'night' ? 'night' : currentTheme === 'eye-protect' ? 'eye-protect' : 'bg-[#ebe5d8]']">
        <!-- 面包屑导航 -->
        <div class="max-w-7xl mx-auto px-4 py-3 text-sm theme-text">
            <a href="#" class="hover:text-[#469b75]">首页</a>
            <span class="mx-2">></span>
            <a href="#" class="hover:text-[#469b75]">玄幻网络</a>
            <span class="mx-2">></span>
            <a href="#" class="hover:text-[#469b75]">叛出家族后，转身投靠魔族女帝</a>
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

const route = useRoute()
const router = useRouter()

// 主题状态
const currentTheme = ref('') // '' 表示默认主题, 'night' 表示深色主题, 'eye-protect' 表示护眼主题

// 滚动显示控制状态
const showControls = ref(false)
const lastScrollTop = ref(0)
const scrollThreshold = 50 // 增加阈值，减少轻微滚动的影响

// 字体大小控制
const fontSize = ref('medium') // 'small', 'medium', 'large'
// 字体类型控制
const fontType = ref('default') // 'default', 'handwriting'
// 注音控制
const showPinyin = ref(false)

// 切换字体大小
const setFontSize = (size) => {
    fontSize.value = size;
};

// 注音功能
import pinyin from 'pinyin'

// 为文字添加注音
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



// 切换深色主题
const toggleNightMode = () => {
    if (currentTheme.value === 'night') {
        currentTheme.value = '' // 恢复默认主题
    } else {
        currentTheme.value = 'night'
        // 切换深色主题时，取消护眼主题
        document.body.classList.remove('eye-protect')
    }
}

// 切换护眼主题
const toggleEyeProtect = () => {
    if (currentTheme.value === 'eye-protect') {
        currentTheme.value = '' // 恢复默认主题
    } else {
        currentTheme.value = 'eye-protect'
        // 切换护眼主题时，取消深色主题
        document.body.classList.remove('night')
    }
}

// 获取路由参数
const bookId = route.params.bookId
const chapterId = route.params.chapterId

// 模拟书籍数据
const bookData = {
    title: '叛出家族后，转身投靠魔族女帝',
    chapters: [
        {
            id: '1334332598936240128',
            title: '第一卷 第1082章 忍疼割爱',
            content: [
                '千幻长老看好龙千跃，认为他有潜力成为下一代的宗门领袖。',
                '在宗门大会上，千幻长老宣布了一个重要决定，让所有弟子都感到震惊。',
                '龙千跃听到这个消息后，心中五味杂陈，不知道该如何面对这个突如其来的机遇。',
                '他知道，这不仅是对他实力的认可，更是一份沉甸甸的责任。',
                '在接下来的日子里，龙千跃开始了更加刻苦的修炼，为了不辜负千幻长老的期望。'
            ]
        },
        {
            id: '1334332598936240129',
            title: '第一卷 第1081章 宗门大会',
            content: [
                '宗门大会即将召开，各弟子都在紧张准备着自己的参赛项目。',
                '龙千跃也不例外，他每天都在修炼场中刻苦训练，希望能在大会上取得好成绩。',
                '随着大会日期的临近，宗门内的气氛变得越来越紧张。',
                '终于，大会的日子到了，所有弟子都聚集在宗门广场上，等待着大会的开始。',
                '宗主宣布大会开始后，各弟子依次上台展示自己的实力。'
            ]
        },
        {
            id: '1334332598936240130',
            title: '第一卷 第1080章 神秘强者',
            content: [
                '叶楚尘在修炼中遇到了一位神秘的强者，对方似乎对他的血脉很感兴趣。',
                '这位神秘强者提出要收叶楚尘为徒，传授他更高级的修炼功法。',
                '叶楚尘陷入了两难的境地，不知道是否应该接受这位神秘强者的邀请。',
                '经过深思熟虑，叶楚尘最终决定接受邀请，跟随这位神秘强者学习。',
                '在神秘强者的指导下，叶楚尘的修为突飞猛进，很快就超越了其他弟子。'
            ]
        }
    ]
}

// 当前章节索引
const currentChapterIndex = ref(bookData.chapters.findIndex(chapter => chapter.id === chapterId))

// 计算属性
const bookTitle = computed(() => bookData.title)
const chapterTitle = computed(() => bookData.chapters[currentChapterIndex.value].title)
const chapterContent = computed(() => bookData.chapters[currentChapterIndex.value].content)
const hasPreviousChapter = computed(() => currentChapterIndex.value > 0)
const hasNextChapter = computed(() => currentChapterIndex.value < bookData.chapters.length - 1)

// 方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}

const goToc = (chapterId) => {
    router.push(`/book/${bookId}/toc`)
}

const previousChapter = () => {
    if (hasPreviousChapter.value) {
        currentChapterIndex.value--
        const previousChapter = bookData.chapters[currentChapterIndex.value]
        router.push(`/book/${bookId}/${previousChapter.id}`)
    }
}

const nextChapter = () => {
    if (hasNextChapter.value) {
        currentChapterIndex.value++
        const nextChapter = bookData.chapters[currentChapterIndex.value]
        router.push(`/book/${bookId}/${nextChapter.id}`)
    }
}

// 节流函数 - 修正实现，确保每隔指定时间执行一次
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

// 滚动监听处理函数（带节流）
const handleScroll = throttle(() => {
    const scrollTop = window.pageYOffset || document.documentElement.scrollTop;
    // 只有当滚动超过阈值时才触发显示/隐藏逻辑
    if (Math.abs(scrollTop - lastScrollTop.value) > scrollThreshold) {
        // 计算新的状态
        const newShowControls = scrollTop < lastScrollTop.value;
        // 只有当状态实际改变时才更新，避免频繁切换导致的闪烁
        if (newShowControls !== showControls.value) {
            showControls.value = newShowControls;
        }
        lastScrollTop.value = scrollTop;
    }
}, 100);

// 组件挂载时添加滚动监听
onMounted(() => {
    window.addEventListener('scroll', handleScroll);
});

// 组件卸载时移除滚动监听
onUnmounted(() => {
    window.removeEventListener('scroll', handleScroll);
});
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
    /* 允许文本正常换行 */
    white-space: normal;
    /* 确保ruby元素作为整体参与排版 */
    line-height: normal;
    /* 增加ruby元素间的间隙 */
    margin: 0 1px;
}

/* 使用更具体的选择器覆盖浏览器默认样式 */
.theme-text ruby > rt {
    font-size: 0.8em !important;
    line-height: 1;
    color: inherit;
    opacity: 0.8;
    /* 优化对齐和防止叠加 */
    text-align: center;
    /* 增加拼音之间的水平间隙 */
    padding: 0 2px;
    /* 增加拼音与汉字之间的垂直间隙 */
    margin-bottom: 4px;
    /* 行间隙 */
    margin-top: 8px;
    /* 使用浏览器原生的ruby-text定位 */
    display: ruby-text;
}

/* 优化带注音文本的整体段落排版 */
.theme-text p {
    /* 增加行高以容纳注音并保持行间距 */
    line-height: 2.2;
    /* 确保文本对齐整齐 */
    text-align: justify;
    /* 控制段落间距 */
    margin-bottom: 1em;
    /* 增加段落内边距，提高整体阅读体验 */
    padding: 0.2em 0;
}

.theme-text rp {
    display: none;
}
</style>