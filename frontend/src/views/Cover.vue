<template>
    <div class="hidden md:block">
        <Header  />
    </div>
    
    <!-- 小屏幕 -->
    <div class="block md:hidden">
        <!-- 顶部导航 -->
        <div class="sticky top-0 z-50 bg-[#469b75] text-white shadow-md">
            <div class="flex items-center justify-between h-12 px-4">
                <button @click="goBack" class="text-white text-xl cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <span class="font-medium">{{ bookData.title }}</span>
                <button @click="goHome" class="text-white text-xl cursor-pointer">
                    <i class="fa-solid fa-house"></i>
                </button>
            </div>
        </div>

        <!-- 书籍信息区域 -->
        <div class="px-4 py-3 bg-white">          
            <div class="flex gap-8 mb-4">
                <!-- 左侧封面图片 -->
                <div class="flex-shrink-0">
                    <img :src="bookData.cover" :alt="bookData.title" class="w-28 h-40 object-cover rounded-md shadow" />
                </div>
                
                <!-- 右侧书籍详情列表 -->
                <div class="flex-1">
                    <div class="text-sm text-gray-600 mb-2">作者：{{ bookData.author }}</div>
                    <div class="text-sm text-gray-600 mb-2">类别：{{ bookData.category }}</div>
                    <div class="text-sm text-gray-600 mb-2">状态：{{ bookData.status }}</div>
                    <div class="text-sm text-gray-600 mb-2">更新：25-11-13</div>
                    <div class="text-sm text-gray-600 mb-2">评分：6.5分</div>
                    <div class="text-sm text-gray-600">点击：424</div>
                </div>
            </div>

            <!-- 书籍简介 -->
            <div class="text-sm text-gray-700 leading-relaxed mb-4">
                <div class="text-gray-800 font-medium mb-2">简介：</div>
                <div v-for="(paragraph, index) in bookData.description" :key="index" class="mb-2">
                    {{ paragraph }}
                </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex gap-3 mb-4 justify-evenly">
                <button class="  bg-[#469b75] hover:bg-[#3d8766] text-white py-2 px-4 rounded-4xl font-medium" @click="handleRead">
                    开始阅读
                </button>
                <button class=" border-2 border-[#469b75] text-[#469b75] hover:bg-[#469b75] hover:text-white py-2 px-4  rounded-4xl font-medium">
                    加入书架
                </button>
            </div>
        </div>

        <!-- 最新章节 -->
        <div class="bg-white mt-3 px-4 py-3">
            <div class="flex items-center justify-between mb-3 pb-2 border-b">
                <h3 class="font-bold text-gray-800">最新章节</h3>
                <div class="text-xs text-gray-500">更新：25-11-13 01:23:49</div>
            </div>

            <div class="space-y-3">
                <div class="border-l-2 border-[#469b75] pl-3 py-1" v-for="(chapter, index) in bookData.latestChapters" :key="index">
                    <router-link :to="`/book/${bookData.bookId}/${chapter.chapterId}`" class="text-sm text-gray-800 mb-1 font-medium cursor-pointer hover:text-[#469b75] transition-colors">{{ chapter.title }}</router-link>
                    <div class="text-xs text-gray-600">{{ chapter.content }}</div>
                </div>
            </div>

            <div class="text-center mt-4 mb-16">
                <a @click="goToToc" class="text-[#469b75] text-sm cursor-pointer">查看完整目录 ></a>
            </div>
        </div>

        <!-- 底部导航 -->
        <ToolBar :showControls="true" />
        
    </div>
    <!-- 大屏 -->
    <div class="hidden md:block container md:w-5xl  mx-auto px-4 py-8">
        <!-- Breadcrumb -->
        <div class="max-w-7xl mx-auto px-4 py-4 text-sm text-gray-600">
            <a href="#" class="hover:text-[#469b75]">{{ bookData.breadcrumb.home }}</a>
            <span class="mx-2">></span>
            <a href="#" class="hover:text-[#469b75]">{{ bookData.breadcrumb.category }}</a>
            <span class="mx-2">></span>
            <span class="text-gray-800">{{ bookData.title }}</span>
        </div>

        <!-- Main Content -->
        <main class="max-w-7xl mx-auto px-4 pb-12">
            <div class="flex gap-6">
                <!-- Left Column -->
                <div class="flex-1">
                    <!-- Book Info Card -->
                    <div class="bg-white rounded-lg shadow p-6 mb-6">
                        <div class="flex gap-6">
                            <!-- Book Cover -->
                            <div class="flex-shrink-0">
                                <img :src="bookData.cover" :alt="bookData.title"
                                    class="w-40 h-56 rounded-lg shadow-lg object-cover" />
                            </div>

                            <!-- Book Details -->
                            <div class="flex-1">
                                <div class="mb-4">
                                    <h1 class="text-2xl font-bold text-gray-800 mb-2">
                                        {{ bookData.title }}
                                    </h1>
                                    <p class="text-gray-600">{{ bookData.author }} 著</p>
                                </div>

                                <div class="flex flex-wrap gap-4 text-sm text-gray-600 mb-4">
                                    <span>类别：<span class="text-gray-800">{{ bookData.category }}</span></span>
                                    <span>状态：<span class="text-gray-800">{{ bookData.status }}</span></span>
                                    <span>总点击：<span class="text-gray-800">{{ bookData.totalClicks }}</span></span>
                                    <span>总字数：<span class="text-gray-800">{{ bookData.totalWords }}</span></span>
                                </div>

                                <div class="text-sm text-gray-700 leading-relaxed mb-4 max-h-32 overflow-hidden">
                                    <p v-for="(paragraph, index) in bookData.description" :key="index">
                                        {{ paragraph }}
                                    </p>
                                </div>

                                <div class="flex gap-4">
                                    <button
                                        class="bg-[#469b75] hover:bg-[#3d8766] text-white px-8 py-2 rounded-md font-medium"
                                        @click="handleRead">
                                        点击阅读
                                    </button>
                                    <button
                                        class="border-2 border-[#469b75] text-[#469b75] hover:bg-[#469b75] hover:text-white px-8 py-2 rounded-md font-medium"
                                        @click="handleAddToShelf">
                                        加入书架
                                    </button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- Latest Chapters -->
                    <div class="bg-white rounded-lg shadow p-6 mb-6">
                        <div class="flex items-center justify-between mb-4 pb-4 border-b">
                            <div>
                                <h3 class="text-xl font-bold text-gray-800 inline-block mr-2">最新章节</h3>
                                <span class="text-gray-500 text-sm">({{ bookData.totalChapters }}章)</span>
                            </div>
                            <a @click="goToToc" class="text-[#469b75] hover:underline text-sm cursor-pointer">全部目录</a>
                        </div>

                        <div class="space-y-3" v-for="(chapter, index) in bookData.latestChapters" :key="index">
                            <div class="flex justify-between items-start">
                                <a href="#" class="text-gray-800 hover:text-[#469b75] flex-1">
                                    {{ chapter.title }}
                                </a>
                                <span class="text-gray-500 text-sm ml-4 whitespace-nowrap">更新时间：{{ chapter.updateTime
                                    }}</span>
                            </div>
                            <div class="text-gray-600 text-sm pl-4">
                                {{ chapter.content }}
                            </div>
                        </div>
                    </div>

                    <!-- Comments Section -->
                    <div class="bg-white rounded-lg shadow p-6">
                        <div class="flex items-center justify-between mb-4 pb-4 border-b">
                            <div>
                                <h3 class="text-xl font-bold text-gray-800 inline-block mr-2">作品评论区</h3>
                                <span class="text-gray-500 text-sm">({{ bookData.comments.length }}条)</span>
                            </div>
                            <a href="#" class="text-[#469b75] hover:underline text-sm">发表评论</a>
                        </div>

                        <!-- Comment Items -->
                        <div class="flex gap-4 mb-6" v-for="(comment, index) in bookData.comments" :key="index">
                            <div class="flex-shrink-0">
                                <img :src="comment.avatar" :alt="comment.userName"
                                    class="w-12 h-12 rounded-full object-cover" />
                            </div>
                            <div class="flex-1">
                                <div class="mb-2">
                                    <span class="font-medium text-gray-800">{{ comment.userName }}</span>
                                    <span class="text-gray-500 text-sm ml-2">{{ comment.userLocation }}</span>
                                </div>
                                <p class="text-gray-700 mb-2">{{ comment.content }}</p>
                                <div class="flex items-center gap-4 text-sm text-gray-500">
                                    <span>{{ comment.createTime }}</span>
                                    <button class="hover:text-[#469b75]">回复 ({{ comment.replyCount }})</button>
                                    <button class="hover:text-[#469b75]">赞 ({{ comment.likeCount }})</button>
                                    <button class="hover:text-[#469b75]">踩 ({{ comment.dislikeCount }})</button>
                                </div>
                            </div>
                        </div>

                        <div class="text-center mb-4">
                            <a href="#" class="text-[#469b75] hover:underline text-sm">查看全部评论 ></a>
                        </div>

                        <!-- Comment Form -->
                        <div class="border-t pt-4">
                            <h4 class="text-base font-medium text-gray-800 mb-3">发表评论</h4>
                            <textarea v-model="commentContent" placeholder="我来说两句..."
                                class="w-full border border-gray-300 rounded-md p-3 focus:outline-none focus:border-[#469b75] resize-none"
                                rows="4" maxlength="1000"></textarea>
                            <div class="flex justify-between items-center mt-2">
                                <span class="text-sm text-gray-500">{{ commentContent.length }}/1000 字</span>
                                <button class="bg-[#469b75] hover:bg-[#3d8766] text-white px-6 py-2 rounded-md"
                                    :disabled="!commentContent.trim()" @click="submitComment">
                                    发表
                                </button>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Right Sidebar -->
                <div class="w-80 space-y-6">
                    <!-- Author Info -->
                    <div class="bg-white rounded-lg shadow p-6">
                        <div class="flex gap-4 mb-4 pb-4 border-b">
                            <img :src="bookData.authorInfo.avatar" :alt="bookData.authorInfo.name"
                                class="w-16 h-16 rounded-full object-cover flex-shrink-0" />
                            <div>
                                <span class="inline-block bg-[#469b75] text-white text-xs px-2 py-1 rounded mb-1">{{
                                    bookData.authorInfo.type }}</span>
                                <h4 class="font-bold text-gray-800">{{ bookData.authorInfo.name }}</h4>
                            </div>
                        </div>
                        <div>
                            <h5 class="font-medium text-gray-800 mb-2">作者有话说</h5>
                            <p class="text-sm text-gray-600">
                                {{ bookData.authorInfo.message }}
                            </p>
                        </div>
                    </div>

                    <!-- Related Books -->
                    <div class="bg-white rounded-lg shadow p-6">
                        <h3 class="text-xl font-bold text-gray-800 mb-4 pb-4 border-b">同类推荐</h3>

                        <div class="space-y-4">
                            <div v-for="(book, index) in relatedBooks" :key="index"
                                class="flex gap-3 cursor-pointer hover:bg-gray-50 p-1 rounded transition-colors"
                                @click="navigateToBook(index + 1)">
                                <img :src="book.cover" :alt="book.title"
                                    class="w-16 h-20 object-cover rounded flex-shrink-0" />
                                <div class="flex-1 min-w-0">
                                    <h5 class="font-medium text-gray-800 text-sm mb-1 truncate">{{ book.title }}</h5>
                                    <p class="text-xs text-gray-600 line-clamp-3">{{ book.description }}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </main>
    </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import ToolBar from '@/components/ToolBar.vue'

const router = useRouter()

// 回退和首页跳转方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}

// Reactive data
const commentContent = ref('')

// Mock book data
const bookData = {
    bookId: 1,
    title: '叛出家族后，转身投靠魔族女帝',
    author: '锦言公举',
    cover: 'https://ext.same-assets.com/4081703551/1339035738.jpeg',
    category: '玄幻奇幻',
    status: '连载中',
    totalClicks: '418,567',
    totalWords: '2,705,761',
    totalChapters: '1082',
    description: [
        '【玄幻+无敌+反穿越+杀伐果断+绝不回头+剑道+吞噬流】为了让家族晋升不朽世家，叶楚尘不惜与魔族女帝林书望大战，拼得筋脉尽毁，神体破碎，修为尽失。',
        '叶家成功晋升为不朽世家，掌管九十九条灵脉，而他楚尘被废后却受尽父亲和七个姐姐嘲讽，被剥夺圣子之位...',
        '就在他万念俱灰之际，偶然获得了魔族女帝林书望的传承，从此踏上了一条逆袭之路。'
    ],
    breadcrumb: {
        home: '小说精品屋',
        category: '玄幻奇幻'
    },
    latestChapters: [
        {
            chapterId: '1334332598936240128',
            title: '第一卷 第1082章 忍疼割爱',
            updateTime: '25/11/13 01:23:49',
            content: '千幻长老看好龙千跃，认为他有潜力成为下一代的宗门领袖。'
        },
        {
            chapterId: '1334332598936240129',
            title: '第一卷 第1081章 宗门大会',
            updateTime: '25/11/12 23:15:32',
            content: '宗门大会即将召开，各弟子都在紧张准备着自己的参赛项目。'
        },
        {
            chapterId: '1334332598936240130',
            title: '第一卷 第1080章 神秘强者',
            updateTime: '25/11/12 18:45:17',
            content: '叶楚尘在修炼中遇到了一位神秘的强者，对方似乎对他的血脉很感兴趣。'
        }
    ],
    comments: [
        {
            avatar: 'https://api.dicebear.com/9.x/avataaars/svg?seed=50',
            userName: '1356****324',
            userLocation: '湖北读者',
            content: '可以试试看，这本书的设定还挺新颖的。',
            createTime: '4个月前',
            replyCount: 0,
            likeCount: 12,
            dislikeCount: 1
        },
        {
            avatar: 'https://api.dicebear.com/9.x/avataaars/svg?seed=12',
            userName: '2234****5678',
            userLocation: '广东读者',
            content: '作者更新很稳定，剧情也越来越精彩了！',
            createTime: '3个月前',
            replyCount: 2,
            likeCount: 25,
            dislikeCount: 0
        },
        {
            avatar: 'https://api.dicebear.com/9.x/avataaars/svg?seed=43',
            userName: '6789****1234',
            userLocation: '北京读者',
            content: '主角的性格我很喜欢，杀伐果断，不拖泥带水。',
            createTime: '2个月前',
            replyCount: 1,
            likeCount: 18,
            dislikeCount: 0
        },
        {
            avatar: 'https://api.dicebear.com/9.x/avataaars/svg?seed=44',
            userName: '9876****5432',
            userLocation: '上海读者',
            content: '期待接下来的剧情发展，希望作者能保持这个节奏。',
            createTime: '1个月前',
            replyCount: 0,
            likeCount: 9,
            dislikeCount: 0
        }
    ],
    authorInfo: {
        name: '锦言公举',
        avatar: 'https://api.dicebear.com/9.x/avataaars/svg?seed=45',
        type: '签约作家',
        message: '亲亲们，你们的支持是我最大的动力！求点击、求推荐、求书评哦！'
    }
}

// Mock related books
const relatedBooks = [
    {
        title: '穿成极品老太，整...',
        description:
            '原名：极品老太整顿家门：开局喜提七孝子白贞下界后傻了——原主上一世为了这七个白眼狼...',
        cover: 'https://ext.same-assets.com/4081703551/959896919.jpeg'
    },
    {
        title: '玄幻：从猎户到万...',
        description:
            '妖魔乱世，凡人如蝼蚁！王煊穿越异界，意外激活【深红面板】...',
        cover: 'https://ext.same-assets.com/4081703551/1339035738.jpeg'
    },
    {
        title: '我想成为影之实力者',
        description:
            '「吾名暗影。乃潜伏于暗影之中，狩猎暗影之人……」少年席德憧憬着...',
        cover: 'https://ext.same-assets.com/4081703551/594530825.png'
    },
    {
        title: '天剑神狱',
        description:
            '太古时期，诸圣争霸，乱天动地，世界湮灭！无上剑主、至高妖皇...',
        cover: 'https://ext.same-assets.com/4081703551/3622865948.jpeg'
    }
]

// Methods
const handleRead = () => {
    // 构建路由路径: /book/bookid/chapterid
    const bookId = bookData.bookId
    const chapterId = bookData.latestChapters[0].chapterId
    router.push(`/book/${bookId}/${chapterId}`)
}

const handleAddToShelf = () => {
    console.log('Add book to shelf')
}

const submitComment = () => {
    if (!commentContent.value.trim()) return
    console.log('Submit comment:', commentContent.value)
    commentContent.value = ''
}

const navigateToBook = (bookId) => {
    console.log('Navigate to book:', bookId)
    router.push(`/book/${bookId}`)
}

const goToToc = () => {
    const bookId = bookData.bookId
    router.push(`/book/${bookId}/toc`)
}
</script>

<style scoped>
/* Optional: add custom styles if needed */
</style>