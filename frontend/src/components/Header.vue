<template>
    <searchBar.define>
        <!-- 搜索框 -->
        <div class="relative w-full">
            <input v-model="searchQuery" type="text" placeholder="书名、作者、关键字"
                class="w-full py-2 px-4 pr-10 rounded-full border border-gray-300 focus:outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary transition-all"
                @keyup.enter="handleSearch" />
            <button class="absolute right-3 top-1/2 transform -translate-y-1/2 text-primary"
                @click="handleSearch">
                <i class="fa-solid fa-magnifying-glass"></i>
            </button>
        </div>
    </searchBar.define>

    <!-- 顶部导航栏 -->
    <header class="bg-white shadow-sm sticky top-0 z-50 md:w-5xl mx-auto">
        <!-- 小屏幕第一行：Logo、书架图标、用户图标 -->
        <div class="container mx-auto px-4 py-4 flex items-center justify-between">
            <router-link to="/" class="flex items-center space-x-2 cursor-pointer">
                <i class="fa-solid fa-book text-primary text-2xl"></i>
                <span class="text-xl font-bold text-primary">阁林小说</span>
            </router-link>

            <div class="min-w-80 hidden md:block">
                <div class="flex items-center flex-1 max-w-md mx-auto">
                    <!-- 搜索框 -->
                    <searchBar.reuse />
                </div>
            </div>

            <div class="flex items-center space-x-4">
                <!-- 小屏幕显示书架图标 -->
                <router-link to="/history" class="md:hidden text-gray-600 hover:text-primary transition-colors">
                    <i class="fa-solid fa-swatchbook text-primary"></i>
                </router-link>
                <!-- 小屏幕显示人像图标 -->
                <a href="#" class="md:hidden text-gray-600 hover:text-primary transition-colors">
                    <i class="fa-solid fa-user"></i>
                </a>
                <!-- 大屏幕显示完整用户操作区 -->
                <router-link to="/history"
                    class="hidden md:inline-flex items-center text-gray-600 hover:text-primary transition-colors">
                    <i class="fa-solid fa-swatchbook text-primary mr-1"></i> 我的书架
                </router-link>
                <span class="hidden md:inline-block text-gray-300 mr-4">|</span>
                <a href="#" class="hidden md:inline-flex text-gray-600 hover:text-primary transition-colors">登录</a>
                <a href="#" class="hidden md:inline-flex text-gray-600 hover:text-dark transition-colors">注册</a>
            </div>
        </div>

        <!-- 小屏幕第二行：搜索框，独占一行 -->
        <div class="container mx-auto px-4 pb-4 md:hidden">
            <searchBar.reuse />
        </div>

        <!-- 主导航 -->
        <nav class="bg-primary text-white hidden md:block">
            <div class="container mx-auto  w-5xl">
                <ul class="flex space-x-1 md:space-x-6 py-2 overflow-x-auto">
                    <li>
                        <router-link to="/" class="px-9 py-1 rounded hover:bg-secondary transition-colors">首页</router-link>
                    </li>
                    <li>
                        <router-link to="/class" class="px-9 py-1 rounded hover:bg-secondary transition-colors">全部作品</router-link>
                    </li>
                    <li>
                        <router-link to="/rank" class="px-9 py-1 rounded hover:bg-secondary transition-colors">排行榜</router-link>
                    </li>
                    <li>
                        <router-link to="/admin" class="px-9 py-1 rounded hover:bg-secondary transition-colors">作者专区</router-link>
                    </li>
                </ul>
            </div>
        </nav>
    </header>
</template>

<script setup>
import { ref } from 'vue'
import { createReusableTemplate } from '@vueuse/core'

const searchBar = createReusableTemplate()
// 响应式搜索输入
const searchQuery = ref('')

// 搜索处理函数（可后续对接路由或 API）
const handleSearch = () => {
    if (searchQuery.value.trim()) {
        console.log('搜索关键词:', searchQuery.value)
        // 例如：router.push(`/search?q=${encodeURIComponent(searchQuery.value)}`)
    }
}
</script>