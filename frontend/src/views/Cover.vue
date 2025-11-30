<template>
    <div class="hidden md:block">
        <Header />
    </div>

    <!-- 小屏幕 -->
    <div class="block md:hidden">
        <!-- 顶部导航 -->
        <div class="sticky top-0 z-50 bg-[#469b75] text-white shadow-md">
            <div class="flex items-center justify-between h-12 px-4">
                <button @click="goBack" class="text-white text-xl cursor-pointer">
                    <i class="fas fa-arrow-left"></i>
                </button>
                <span class="font-medium">{{ book.title }}</span>
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
                    <img :src="book.cover" :alt="book.title" class="w-28 h-40 object-cover rounded-md shadow" />
                </div>

                <!-- 右侧书籍详情列表 -->
                <div class="flex-1">
                    <div class="text-sm text-gray-600 mb-2">作者：{{ book.author }}</div>
                    <div class="text-sm text-gray-600 mb-2">类别：{{ book.category }}</div>
                    <div class="text-sm text-gray-600 mb-2">状态：{{ book.status === 'serializing' ? '连载中' : book.status
                        === 'completed' ? '完结' : book.status }}</div>
                    <div class="text-sm text-gray-600 mb-2">更新：25-11-13</div>
                    <div class="text-sm text-gray-600 mb-2">评分：6.5分</div>
                    <div class="text-sm text-gray-600">点击：{{ book.clickCount }}</div>
                </div>
            </div>

            <!-- 书籍简介 -->
            <div class="text-sm text-gray-700 leading-relaxed mb-4">
                <div class="text-gray-800 font-medium mb-2">简介：</div>
                <div v-for="(paragraph, index) in book.description" :key="index" class="mb-2">
                    {{ paragraph }}
                </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex gap-3 mb-4 justify-evenly">
                <button class="  bg-[#469b75] hover:bg-[#3d8766] text-white py-2 px-4 rounded-4xl font-medium"
                    @click="handleRead">
                    开始阅读
                </button>
                <button
                    class=" border-2 border-[#469b75] text-[#469b75] hover:bg-[#469b75] hover:text-white py-2 px-4  rounded-4xl font-medium">
                    加入书架
                </button>
            </div>
        </div>

        <!-- 最新章节 -->
        <div class="bg-white mt-3 px-4 py-3">
            <div class="flex items-center justify-between mb-3 pb-2 border-b">
                <h3 class="font-bold text-gray-800">最新章节</h3>
                <div class="text-xs text-gray-500">更新：{{ book.updateTime ? new
                    Date(book.updateTime).toLocaleDateString('zh-CN').replace(/\//g, '-') : '未知' }}</div>
            </div>

            <div class="space-y-3">
                <div class="border-l-2 border-[#469b75] pl-3 py-1" v-for="(chapter, index) in chapters.slice(-3).reverse()"
                    :key="index">
                    <router-link :to="`/book/${book.id}/${chapter.chapterNo}`"
                        class="text-sm text-gray-800 mb-1 font-medium cursor-pointer hover:text-[#469b75] transition-colors">第{{ chapter.chapterNo }}章 {{
                            chapter.title }}</router-link>
                </div>
            </div>

            <div class="text-center mt-4 mb-16">
                <a @click="goToToc" class="text-[#469b75] text-sm cursor-pointer">查看完整目录 ></a>
            </div>
        </div>

        <!-- 底部导航 -->
        <!-- <ToolBar :showControls="true" /> -->

    </div>
    <!-- 大屏 -->
    <div class="hidden md:block container md:w-5xl  mx-auto px-4 py-8">
        <!-- Breadcrumb -->
        <div class="max-w-7xl mx-auto px-4 py-4 text-sm text-gray-600">
            <a href="#" class="hover:text-[#469b75]">{{ breadcrumb.home }}</a>
            <span class="mx-2">></span>
            <a href="#" class="hover:text-[#469b75]">{{ breadcrumb.category }}</a>
            <span class="mx-2">></span>
            <span class="text-gray-800">{{ book.title }}</span>
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
                                <img :src="book.cover" :alt="book.title"
                                    class="w-40 h-56 rounded-lg shadow-lg object-cover" />
                            </div>

                            <!-- Book Details -->
                            <div class="flex-1">
                                <div class="mb-4">
                                    <h1 class="text-2xl font-bold text-gray-800 mb-2">
                                        {{ book.title }}
                                    </h1>
                                    <p class="text-gray-600">{{ book.author }} 著</p>
                                </div>

                                <div class="flex flex-wrap gap-4 text-sm text-gray-600 mb-4">
                                    <span>类别：<span class="text-gray-800">{{ book.category }}</span></span>
                                    <span>状态：<span class="text-gray-800">{{ book.status === 'serializing' ? '连载中' :
                                        book.status === 'completed' ? '完结' : book.status }}</span></span>
                                    <span>点击量：<span class="text-gray-800">{{ book.clickCount }}</span></span>
                                    <span>总字数：<span class="text-gray-800">{{ book.wordCount }}</span></span>
                                </div>

                                <div class="text-sm text-gray-700 leading-relaxed mb-4 max-h-32 overflow-hidden">
                                    <p v-for="(paragraph, index) in book.description" :key="index">
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
                                <span class="text-gray-500 text-sm">({{ chapters.length }}章)</span>
                            </div>
                            <a @click="goToToc" class="text-[#469b75] hover:underline text-sm cursor-pointer">目录</a>
                        </div>

                        <div class="space-y-3" v-for="(chapter, index) in chapters.slice(-3).reverse()" :key="index">
                            <div class="flex justify-between items-start">
                                <router-link :to="`/book/${book.id}/${chapter.chapterNo}`"
                                    class=" max-w-[25rem] text-gray-800 hover:text-[#469b75] flex-1 my-2 whitespace-nowrap overflow-hidden text-ellipsis">
                                    第{{chapter.chapterNo}}章 {{ chapter.title }}
                                </router-link>
                                <span class="text-gray-500 text-sm ml-4 whitespace-nowrap">更新时间：{{ chapter.updateTime ?
                                    new Date(chapter.updateTime*1000).toLocaleDateString('zh-CN') : '未知' }}</span>
                            </div>
                        </div>
                    </div>

                    <!-- Comments Section -->
                    <div class="bg-white rounded-lg shadow p-6">
                        <div class="flex items-center justify-between mb-4 pb-4 border-b">
                            <div>
                                <h3 class="text-xl font-bold text-gray-800 inline-block mr-2">作品评论区</h3>
                                <span class="text-gray-500 text-sm">({{ comments.length }}条)</span>
                            </div>
                            <a href="#" class="text-[#469b75] hover:underline text-sm">发表评论</a>
                        </div>

                        <!-- Comment Items -->
                        <div class="flex gap-4 mb-6"
                            v-for="(comment, index) in showAllComments ? comments : comments.slice(0, 5)" :key="index">
                            <div class="flex-shrink-0">
                                <img :src="comment.user.avatar" :alt="comment.user.nickname"
                                    class="w-12 h-12 rounded-full object-cover" />
                            </div>
                            <div class="flex-1">
                                <div class="mb-2">
                                    <span class="font-medium text-gray-800">{{ comment.user.nickname }}</span>
                                    <span class="text-gray-500 text-sm ml-2">{{ comment.user.location }}</span>
                                </div>
                                <p class="text-gray-700 mb-2">{{ comment.content }}</p>
                                <div class="flex items-center gap-4 text-sm text-gray-500">
                                    <span>{{ comment.createTime ? new Date(comment.createTime).toLocaleString('zh-CN') :
                                        '未知' }}</span>
                                    <button class="hover:text-[#469b75]">↩️回复 ({{ comment.replyCount }})</button>
                                    <button class="hover:text-[#469b75]">👍赞 ({{ comment.likeCount }})</button>
                                    <button class="hover:text-[#469b75]">👎踩 ({{ comment.dislikeCount }})</button>
                                </div>
                            </div>
                        </div>

                        <div class="text-center mb-4" v-if="comments.length > 5">
                            <a href="#" class="text-[#469b75] hover:underline text-sm"
                                @click.prevent="showAllComments = !showAllComments">
                                {{ showAllComments ? '收起评论 ⏫' : '查看全部评论 ⏬' }}
                            </a>
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
                            <img :src="author.avatar" :alt="author.nickname"
                                class="w-16 h-16 rounded-full object-cover flex-shrink-0" />
                            <div>
                                <span class="inline-block bg-[#469b75] text-white text-xs px-2 py-1 rounded mb-1">{{
                                    author.userType || '作者' }}</span>
                                <h4 class="font-bold text-gray-800">{{ author.nickname }}</h4>
                            </div>
                        </div>
                        <div>
                            <h5 class="font-medium text-gray-800 mb-2">作者有话说</h5>
                            <p class="text-sm text-gray-600">
                                {{ author.desc || '暂无简介' }}
                            </p>
                        </div>
                    </div>

                    <!-- Related Books -->
                    <div class="bg-white rounded-lg shadow p-6">
                        <h3 class="text-xl font-bold text-gray-800 mb-4 pb-4 border-b">同类推荐</h3>

                        <div class="space-y-4">
                            <div v-for="(book, index) in relatedBooks" :key="index"
                                class="flex gap-3 cursor-pointer hover:bg-gray-50 p-1 rounded transition-colors"
                                @click="navigateToBook(book.id)">
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

    <!-- 通知组件 -->
    <Notice :visible="notice.visible" :type="notice.type" :title="notice.title" :message="notice.message"
        @close="closeNotice" />
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import Header from '@/components/Header.vue'
import ToolBar from '@/components/ToolBar.vue'
import Notice from '@/components/Notice.vue'
import { bookAPI, userAPI } from '@/api/services'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

// 回退和首页跳转方法
const goBack = () => {
    router.back()
}

const goHome = () => {
    router.push('/')
}

const breadcrumb = ref({
    home: '阁林小说',
    category: ''
})
// Reactive data
const commentContent = ref('')
const book = ref({})
const relatedBooks = ref([])
const author = ref({})
const comments = ref([])
const chapters = ref([])
// 添加一个响应式变量来控制是否显示全部评论
const showAllComments = ref(false)

// Notice 状态
const notice = ref({
    visible: false,
    type: 'info',
    title: '提示',
    message: ''
})

// 显示通知
const showNotice = (type, title, message) => {
    notice.value = {
        visible: true,
        type,
        title,
        message
    }
}

// 关闭通知
const closeNotice = () => {
    notice.value.visible = false
}

// 获取书籍详情
const fetchBookDetail = async () => {
    const bookId = route.params.id
    try {
        // 获取书籍详情
        const bookResponse = await bookAPI.getById(bookId)
        const bookData = bookResponse.data

        // 将description字符串转换为数组格式以适配模板
        if (typeof bookData.description === 'string') {
            bookData.description = bookData.description.split('\n').filter(line => line.trim())
        }

        // 更新独立的breadcrumb变量
        breadcrumb.value.category = bookData.category || ''

        // 设置书籍信息，不再包含breadcrumb
        book.value = bookData

        // 获取相关书籍
        const relatedResponse = await bookAPI.getRelatedBooks(bookId)
        relatedBooks.value = relatedResponse.data

        // 获取章节列表
        const chaptersResponse = await bookAPI.getChapters(bookId)
        chapters.value = chaptersResponse?.data?.chapters || []

        // 获取评论
        const commentsResponse = await bookAPI.getComments(bookId)
        comments.value = commentsResponse.data

        // 获取作者信息
        const authorResponse = await userAPI.getByName(bookData.author)
        author.value = authorResponse.data
        console.log('Author:', author.value)

    } catch (error) {
        console.error('Failed to fetch book details:', error)
    }
}

onMounted(() => {
    fetchBookDetail()
})

// 监听路由参数变化，重新加载数据
watch(
    () => route.params.id,
    (newId, oldId) => {
        if (newId !== oldId) {
            fetchBookDetail()
        }
    }
)

// Methods
const handleRead = () => {
    console.log('Read book')
    // 构建路由路径: /book/bookid/chapterno
    if (chapters.value.length > 0) {
        const bookId = book.value.id
        const chapterNo = chapters.value[0].chapterNo
        router.push(`/book/${bookId}/${chapterNo}`)
    }
}

const handleAddToShelf = () => {
    if (!userStore.isLoggedIn.value) {
        showNotice('info', '提示', '请先登录')

        // 延迟跳转，让用户看到提示
        // setTimeout(() => {
        //     router.push('/login')
        // }, 1000)
        return
    }
    console.log('Add book to shelf')
    // 实际的加入书架逻辑
    showNotice('success', '成功', '已加入书架')
}

const submitComment = () => {
    if (!userStore.isLoggedIn.value) {
        showNotice('info', '提示', '请先登录')

        // 延迟跳转，让用户看到提示
        // setTimeout(() => {
        //     router.push('/login')
        // }, 1000)
        return
    }
    if (!commentContent.value.trim()) return
    console.log('Submit comment:', commentContent.value)
    // 实际的提交评论逻辑
    commentContent.value = ''
    showNotice('success', '成功', '评论已发表')
}

const navigateToBook = (bookId) => {
    console.log('Navigate to book:', bookId)
    router.push(`/book/${bookId}`)
}

const goToToc = () => {
    const bookId = book.value.id
    router.push(`/book/${bookId}/toc`)
}
</script>

<style scoped>
/* Optional: add custom styles if needed */
</style>