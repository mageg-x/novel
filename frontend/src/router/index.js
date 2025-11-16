import { createRouter, createWebHistory } from 'vue-router'
import Home from '@/views/Home.vue'
import Cover from '@/views/Cover.vue'
import Book from '@/views/Book.vue'
import Toc from '@/views/Toc.vue'
import Class from '@/views/Class.vue'
import Rank from '@/views/Rank.vue'
import History from '@/views/History.vue'
import Author from '@/views/Author.vue'

const routes = [
    {
        path: '/',
        name: 'home',
        component: Home
    },
    {
        path: '/class',
        name: 'class',
        component: Class
    },
    {
        path: '/rank',
        name: 'rank',
        component: Rank
    },
    {
        path: '/history',
        name: 'history',
        component: History
    },
    {
        path: '/book/:id',
        name: 'cover',
        component: Cover
    },
    {
        path: '/book/:bookId/:chapterId',
        name: 'book',
        component: Book
    },
    {
        path: '/book/:bookId/toc',
        name: 'toc',
        component: Toc
    },
    {
        path: '/author/:authorId',
        name: 'author',
        component: Author
    }
]

const router = createRouter({
    history: createWebHistory(),
    routes
})

export default router