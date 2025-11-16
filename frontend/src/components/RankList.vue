<template>
    <div class="w-full">
        <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">
            {{ title }}
        </h3>
        <div class="bg-white rounded-xl shadow-sm py-4 h-full">
            <ul class="space-y-3">
                <li v-for="(book, index) in books" :key="index"
                    class="flex items-center gap-3 p-2 hover:bg-gray-50 rounded card-hover cursor-pointer"
                    @click="navigateToBook(book.id)">
                    <span :class="[
                        'w-6 h-6  rounded-sm flex items-center justify-center text-sm',
                        index < 3
                            ? 'bg-primary text-white'
                            : 'bg-gray-600 text-white'
                    ]">
                        {{ index + 1 }}
                    </span>

                    <img v-if="book.cover" :src="book.cover" :alt="book.title" class="w-10 h-14 object-cover rounded" />

                    <div>
                        <h4 class="text-sm font-medium">{{ book.title }}</h4>
                        <p v-if="book.description" class="text-xs text-gray-500">
                            {{ book.description }}
                        </p>
                    </div>
                </li>
            </ul>
            <router-link v-if="showMore" :to="moreLink"
                class="text-center block text-primary text-sm mt-4 hover:underline">查看更多</router-link>
        </div>
    </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const props = defineProps({
    // 榜单标题，例如：'新书榜单'、'月票排行'、'热门推荐'
    title: {
        type: String,
        default: '排行榜'
    },
    // 榜单书籍列表
    books: {
        type: Array,
        required: true,
        default: () => [],
        validator(value) {
            return value.every(book =>
                typeof book.title === 'string' &&
                (typeof book.description === 'string' || book.description === undefined) &&
                (typeof book.cover === 'string' || book.cover === undefined)
            )
        }
    },
    // 是否显示"查看更多"按钮
    showMore: {
        type: Boolean,
        default: false
    },
    // "查看更多"按钮的跳转链接
    moreLink: {
        type: String,
        default: ''
    }
})

// 导航到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`)
}
</script>

<style scoped>
.card-hover:hover {
    transform: translateY(-2px);
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.08);
    transition: all 0.2s ease;
    cursor: pointer;
}
</style>