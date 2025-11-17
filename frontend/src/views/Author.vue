<template>
    <div class="flex min-h-screen bg-gray-50">
        <!-- 侧边导航栏 -->
        <div class="w-64 bg-gray-800 text-white shadow-lg fixed inset-y-0 left-0 transform transition-transform duration-300 ease-in-out z-50"
            :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'">
            <div class="h-16 flex items-center justify-center bg-gray-900">
                <i class="fas fa-book text-primary text-2xl mr-2"></i>
                <span class="text-xl font-bold">作者后台</span>
            </div>
            <nav class="mt-4">
                <div v-for="item in menuItems" :key="item.page"
                    class="flex items-center px-6 py-3 cursor-pointer transition-colors duration-200"
                    :class="activePage === item.page ? 'bg-primary text-white' : 'text-gray-300 hover:bg-gray-700'"
                    @click="activePage = item.page; isSidebarOpen = false">
                    <i :class="['fas', item.icon, 'mr-3']"></i>
                    <span>{{ item.text }}</span>
                </div>
            </nav>
        </div>

        <!-- 遮罩层 -->
        <div v-if="isSidebarOpen" class="fixed inset-0 bg-black bg-opacity-50 z-40 md:hidden"
            @click="isSidebarOpen = false"></div>

        <!-- 主内容区 -->
        <div class="flex-1 flex flex-col overflow-hidden" :class="isSidebarOpen ? 'ml-0 md:ml-64' : 'ml-0 md:ml-64'">
            <!-- 顶部导航 -->
            <header
                class="h-16 bg-white border-b border-gray-200 flex items-center justify-between px-4 md:px-6 shadow-sm">
                <!-- 移动端菜单按钮 -->
                <button class="md:hidden text-gray-600 mr-3" @click="isSidebarOpen = true">
                    <i class="fas fa-bars text-xl"></i>
                </button>
                <div class="text-gray-600 text-sm md:text-base">作者后台 / {{ getPageTitle() }}</div>
                <div class="flex items-center">
                    <img v-if="user.avatar" :src="user.avatar" :alt="user.nickname"
                        class="w-8 h-8 md:w-10 md:h-10 rounded-full mr-2 md:mr-3 object-cover">
                    <div v-else
                        class="w-8 h-8 md:w-10 md:h-10 rounded-full bg-primary flex items-center justify-center text-white mr-2 md:mr-3">
                        <i class="fas fa-user text-sm md:text-base"></i>
                    </div>
                    <span class="text-gray-700 text-sm md:text-base">{{ user.nickname || user.username }}</span>
                </div>
            </header>

            <!-- 内容区域 -->
            <main class="flex-1 p-6 overflow-y-auto">
                <!-- 小说列表页面 -->
                <div v-if="activePage === 'book-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-xl font-semibold text-gray-800">小说列表</h2>
                        <button class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors"
                            @click="activePage = 'book-add'">
                            <i class="fas fa-plus mr-2"></i>发布小说
                        </button>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="min-w-full">
                            <thead>
                                <tr class="bg-gray-50">
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm">
                                        书名</th>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm hidden md:table-cell">
                                        分类</th>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm hidden sm:table-cell">
                                        点击量</th>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm hidden md:table-cell">
                                        昨日订阅</th>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm hidden lg:table-cell">
                                        更新时间</th>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm hidden sm:table-cell">
                                        总字数</th>
                                    <th
                                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm">
                                        操作</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200">
                                <template v-if="books.length > 0">
                                    <tr v-for="book in paginatedBooks" :key="book.id"
                                        class="hover:bg-gray-50 transition-colors">
                                        <td class="px-4 py-3 whitespace-nowrap">
                                            <div class="flex items-center">
                                                <img :src="book.cover" alt="封面"
                                                    class="w-10 h-14 md:w-12 md:h-16 object-cover rounded mr-2 md:mr-3">
                                                <div>
                                                    <div class="font-medium text-gray-900 text-sm md:text-base">{{
                                                        book.title }}</div>
                                                    <div class="text-xs text-gray-500">{{ book.author }}</div>
                                                </div>
                                            </div>
                                        </td>
                                        <td
                                            class="px-4 py-3 whitespace-nowrap text-sm text-gray-500 hidden md:table-cell">
                                            {{ book.category }}
                                        </td>
                                        <td
                                            class="px-4 py-3 whitespace-nowrap text-sm text-gray-500 hidden sm:table-cell">
                                            {{ book.views }}</td>
                                        <td
                                            class="px-4 py-3 whitespace-nowrap text-sm text-gray-500 hidden md:table-cell">
                                            {{ book.subscriptions
                                            }}</td>
                                        <td
                                            class="px-4 py-3 whitespace-nowrap text-sm text-gray-500 hidden lg:table-cell">
                                            {{ book.updateTime }}
                                        </td>
                                        <td
                                            class="px-4 py-3 whitespace-nowrap text-sm text-gray-500 hidden sm:table-cell">
                                            {{ book.wordCount }}
                                        </td>
                                        <td class="px-4 py-3 whitespace-nowrap text-sm">
                                            <div class="flex flex-wrap gap-2">
                                                <button
                                                    class="text-green-600 hover:text-green-800 text-xs md:text-sm">章节管理</button>
                                                <button
                                                    class="text-red-600 hover:text-red-800 text-xs md:text-sm">删除</button>
                                            </div>
                                        </td>
                                    </tr>
                                </template>
                                <tr v-else>
                                    <td colspan="7" class="px-4 py-8 text-center text-gray-500">
                                        <i class="fas fa-book-open text-2xl mb-2"></i>
                                        <p>暂无小说数据</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 分页导航 -->
                    <div class="mt-4 flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
                        v-if="books.length > 0">
                        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p class="text-sm text-gray-700">
                                    显示 <span class="font-medium">{{ (currentBookPage - 1) * bookPageSize + 1 }}</span> 到
                                    <span class="font-medium">{{ Math.min(currentBookPage * bookPageSize, books.length)
                                        }}</span> 条，共
                                    <span class="font-medium">{{ books.length }}</span> 条结果
                                </p>
                            </div>
                            <div>
                                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                                    aria-label="Pagination">
                                    <button @click="goToPrevBookPage" :disabled="currentBookPage === 1"
                                        class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">上一页</span>
                                        <i class="fas fa-chevron-left"></i>
                                    </button>
                                    <!-- 动态生成页码 -->
                                    <button v-for="page in totalBookPages" :key="page" @click="goToBookPage(page)"
                                        :class="currentBookPage === page
                                            ? 'relative inline-flex items-center z-10 bg-primary border-primary px-4 py-2 text-sm font-medium text-white focus:outline-offset-0'
                                            : 'relative inline-flex items-center border-gray-300 px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'">
                                        {{ page }}
                                    </button>
                                    <button @click="goToNextBookPage" :disabled="currentBookPage === totalBookPages"
                                        class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">下一页</span>
                                        <i class="fas fa-chevron-right"></i>
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 发布小说页面 -->
                <div v-if="activePage === 'book-add'" class="bg-white rounded-lg shadow-sm p-4 md:p-6">
                    <div class="flex justify-between items-center mb-4 md:mb-6 pb-3 md:pb-4 border-b border-gray-200">
                        <h2 class="text-lg md:text-xl font-semibold text-gray-800">发布小说</h2>
                        <button
                            class="px-3 py-1.5 md:px-4 md:py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors text-sm"
                            @click="activePage = 'book-list'">
                            返回列表
                        </button>
                    </div>

                    <form class="space-y-4 md:space-y-6">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">作品方向</label>
                                <select
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                                    <option value="0">男频</option>
                                    <option value="1">女频</option>
                                </select>
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">分类</label>
                                <select
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                                    <option value="1">玄幻奇幻</option>
                                    <option value="2">武侠仙侠</option>
                                    <option value="3">都市言情</option>
                                    <option value="4">历史军事</option>
                                    <option value="5">科幻灵异</option>
                                    <option value="6">游戏竞技</option>
                                </select>
                            </div>

                            <div class="md:col-span-2">
                                <label class="block text-sm font-medium text-gray-700 mb-1">小说名</label>
                                <input type="text" placeholder="请输入小说名"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                            </div>

                            <div>
                                <label class="block text-sm font-medium text-gray-700 mb-1">小说封面</label>
                                <div
                                    class="border-2 border-dashed border-gray-300 rounded-md p-3 md:p-4 text-center cursor-pointer hover:border-primary transition-colors">
                                    <i class="fas fa-upload text-3xl md:text-4xl text-gray-400 mb-1 md:mb-2"></i>
                                    <p class="text-xs md:text-sm text-gray-500">点击上传封面</p>
                                </div>
                            </div>

                            <div class="md:col-span-2">
                                <label class="block text-sm font-medium text-gray-700 mb-1">小说介绍</label>
                                <textarea placeholder="请输入小说介绍" rows="3 md:rows=5"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 resize-vertical text-sm"></textarea>
                            </div>
                        </div>

                        <div class="flex flex-wrap gap-3 md:space-x-4">
                            <button type="button"
                                class="flex-1 min-w-[100px] px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm">提交</button>
                            <button type="button"
                                class="flex-1 min-w-[100px] px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors text-sm">取消</button>
                        </div>
                    </form>
                </div>

                <!-- 章节列表页面 -->
                <div v-if="activePage === 'chapter-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <div class="flex items-center space-x-3">
                            <h2 class="text-xl font-semibold text-gray-800">章节列表</h2>
                            <select v-model="selectedBook"
                                class="px-3 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-primary/50"
                                @change="handleBookChange">
                                <option v-for="book in books" :key="book.id" :value="book">{{ book.title }}</option>
                            </select>
                        </div>
                        <button class="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors"
                            @click="activePage = 'chapter-add'" :disabled="!selectedBook">
                            <i class="fas fa-plus mr-2"></i>新建章节
                        </button>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="min-w-full">
                            <thead>
                                <tr class="bg-gray-50">
                                    <th
                                        class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:px-4 md:py-3 md:text-sm">
                                        章节名</th>
                                    <th
                                        class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:px-4 md:py-3 md:text-sm hidden sm:table-cell">
                                        更新时间</th>
                                    <th
                                        class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:px-4 md:py-3 md:text-sm">
                                        是否收费</th>
                                    <th
                                        class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:px-4 md:py-3 md:text-sm">
                                        操作</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200">
                                <template v-if="chapters.length > 0">
                                    <tr v-for="chapter in paginatedChapters" :key="chapter.chapterId"
                                        class="hover:bg-gray-50 transition-colors">
                                        <td
                                            class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900 md:px-4 md:py-3">
                                            {{
                                                chapter.title }}</td>
                                        <td
                                            class="px-3 py-2 whitespace-nowrap text-sm text-gray-500 md:px-4 md:py-3 hidden sm:table-cell">
                                            {{ chapter.updateTime
                                            }}</td>
                                        <td class="px-3 py-2 whitespace-nowrap md:px-4 md:py-3">
                                            <span :class="[
                                                'px-2 py-1 text-xs rounded-full',
                                                chapter.isVip ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'
                                            ]">
                                                {{ chapter.isVip ? '收费' : '免费' }}
                                            </span>
                                        </td>
                                        <td class="px-3 py-2 whitespace-nowrap text-sm md:px-4 md:py-3">
                                            <div class="flex flex-wrap gap-2">
                                                <button class="text-blue-600 hover:text-blue-800">编辑</button>
                                                <button class="text-red-600 hover:text-red-800">删除</button>
                                            </div>
                                        </td>
                                    </tr>
                                </template>
                                <tr v-else>
                                    <td colspan="4" class="px-4 py-8 text-center text-gray-500">
                                        <i class="fas fa-file-alt text-2xl mb-2"></i>
                                        <p>暂无章节数据</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 分页导航 -->
                    <div class="mt-4 flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
                        v-if="chapters.length > 0">
                        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p class="text-sm text-gray-700">
                                    显示 <span class="font-medium">{{ (currentChapterPage - 1) * chapterPageSize + 1
                                        }}</span> 到 <span class="font-medium">{{ Math.min(currentChapterPage *
                                        chapterPageSize, chapters.length) }}</span>
                                    条，共
                                    <span class="font-medium">{{ chapters.length }}</span> 条结果
                                </p>
                            </div>
                            <div>
                                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                                    aria-label="Pagination">
                                    <button @click="goToPrevChapterPage" :disabled="currentChapterPage === 1"
                                        class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">上一页</span>
                                        <i class="fas fa-chevron-left"></i>
                                    </button>
                                    <!-- 动态生成页码 -->
                                    <button v-for="page in totalChapterPages" :key="page" @click="goToChapterPage(page)"
                                        :class="currentChapterPage === page
                                            ? 'relative inline-flex items-center z-10 bg-primary border-primary px-4 py-2 text-sm font-medium text-white focus:outline-offset-0'
                                            : 'relative inline-flex items-center border-gray-300 px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'">
                                        {{ page }}
                                    </button>
                                    <button @click="goToNextChapterPage"
                                        :disabled="currentChapterPage === totalChapterPages"
                                        class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">下一页</span>
                                        <i class="fas fa-chevron-right"></i>
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 发布章节页面 -->
                <div v-if="activePage === 'chapter-add'" class="bg-white rounded-lg shadow-sm p-4 md:p-6">
                    <div class="flex justify-between items-center mb-4 md:mb-6 pb-3 md:pb-4 border-b border-gray-200">
                        <div class="flex items-center space-x-3">
                            <h2 class="text-lg md:text-xl font-semibold text-gray-800">发布章节</h2>
                            <select v-model="selectedBook"
                                class="px-3 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-primary/50">
                                <option v-for="book in books" :key="book.id" :value="book">{{ book.title }}</option>
                            </select>
                        </div>
                        <button
                            class="px-3 py-1.5 md:px-4 md:py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors text-sm"
                            @click="activePage = 'chapter-list'">
                            返回列表
                        </button>
                    </div>

                    <form class="space-y-4 md:space-y-6">
                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-1">章节名</label>
                            <input type="text" placeholder="请输入章节名"
                                class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2 md:mb-3">章节内容</label>
                            <div class="flex flex-wrap gap-2 mb-3 md:gap-3 md:mb-4">
                                <button type="button"
                                    class="flex-1 min-w-[80px] px-3 py-1.5 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm">
                                    <i class="fas fa-magic mr-1"></i>AI扩写
                                </button>
                                <button type="button"
                                    class="flex-1 min-w-[80px] px-3 py-1.5 bg-green-500 text-white rounded-md hover:bg-green-600 transition-colors text-sm">
                                    <i class="fas fa-compress-alt mr-1"></i>AI缩写
                                </button>
                                <button type="button"
                                    class="flex-1 min-w-[80px] px-3 py-1.5 bg-yellow-500 text-white rounded-md hover:bg-yellow-600 transition-colors text-sm">
                                    <i class="fas fa-pen-alt mr-1"></i>AI续写
                                </button>
                                <button type="button"
                                    class="flex-1 min-w-[80px] px-3 py-1.5 bg-red-500 text-white rounded-md hover:bg-red-600 transition-colors text-sm">
                                    <i class="fas fa-paint-brush mr-1"></i>AI润色
                                </button>
                            </div>
                            <textarea placeholder="请输入章节内容" rows="6 md:rows=10"
                                class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 resize-vertical text-sm"></textarea>
                        </div>

                        <div>
                            <label class="block text-sm font-medium text-gray-700 mb-2">是否收费</label>
                            <div class="flex space-x-4 md:space-x-6">
                                <label class="flex items-center">
                                    <input type="radio" name="isVip" value="0" checked class="text-primary">
                                    <span class="ml-2 text-gray-700 text-sm">免费</span>
                                </label>
                                <label class="flex items-center">
                                    <input type="radio" name="isVip" value="1" class="text-primary">
                                    <span class="ml-2 text-gray-700 text-sm">收费</span>
                                </label>
                            </div>
                        </div>

                        <div class="flex flex-wrap gap-3 md:space-x-4">
                            <button type="button"
                                class="flex-1 min-w-[100px] px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm">提交</button>
                            <button type="button"
                                class="flex-1 min-w-[100px] px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors text-sm">取消</button>
                        </div>
                    </form>
                </div>

                <!-- 数据统计页面 -->
                <div v-if="activePage === 'data-analysis'" class="bg-white rounded-lg shadow-sm p-4 md:p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-lg md:text-xl font-semibold text-gray-800">数据统计</h2>
                        <div class="flex items-center space-x-2">
                            <select
                                class="px-3 py-1 border border-gray-300 rounded text-sm focus:outline-none focus:ring-2 focus:ring-primary/50">
                                <option>近7天</option>
                                <option>近30天</option>
                                <option>近90天</option>
                                <option>全部</option>
                            </select>
                        </div>
                    </div>

                    <!-- 统计概览卡片 -->
                    <div class="grid grid-cols-2 md:grid-cols-2 lg:grid-cols-4 gap-3 md:gap-4 mb-6">
                        <div class="bg-blue-50 p-3 md:p-4 rounded-lg border border-blue-100">
                            <div class="flex items-center justify-between">
                                <div>
                                    <p class="text-xs md:text-sm text-blue-600 mb-1">总点击量</p>
                                    <p class="text-lg md:text-2xl font-bold text-blue-800">14,729,890</p>
                                    <p class="text-xs text-blue-500 mt-1"><i class="fas fa-arrow-up mr-1"></i>12.5% 较上月
                                    </p>
                                </div>
                                <div
                                    class="w-10 h-10 md:w-12 md:h-12 bg-blue-100 rounded-full flex items-center justify-center text-blue-600">
                                    <i class="fas fa-eye text-base md:text-xl"></i>
                                </div>
                            </div>
                        </div>

                        <div class="bg-green-50 p-3 md:p-4 rounded-lg border border-green-100">
                            <div class="flex items-center justify-between">
                                <div>
                                    <p class="text-xs md:text-sm text-green-600 mb-1">总订阅量</p>
                                    <p class="text-lg md:text-2xl font-bold text-green-800">34,832</p>
                                    <p class="text-xs text-green-500 mt-1"><i class="fas fa-arrow-up mr-1"></i>8.3% 较上月
                                    </p>
                                </div>
                                <div
                                    class="w-10 h-10 md:w-12 md:h-12 bg-green-100 rounded-full flex items-center justify-center text-green-600">
                                    <i class="fas fa-bookmark text-base md:text-xl"></i>
                                </div>
                            </div>
                        </div>

                        <div class="bg-purple-50 p-3 md:p-4 rounded-lg border border-purple-100">
                            <div class="flex items-center justify-between">
                                <div>
                                    <p class="text-xs md:text-sm text-purple-600 mb-1">总字数</p>
                                    <p class="text-lg md:text-2xl font-bold text-purple-800">1,400万字</p>
                                    <p class="text-xs text-purple-500 mt-1"><i class="fas fa-arrow-up mr-1"></i>5.2% 较上月
                                    </p>
                                </div>
                                <div
                                    class="w-10 h-10 md:w-12 md:h-12 bg-purple-100 rounded-full flex items-center justify-center text-purple-600">
                                    <i class="fas fa-font text-base md:text-xl"></i>
                                </div>
                            </div>
                        </div>

                        <div class="bg-orange-50 p-3 md:p-4 rounded-lg border border-orange-100">
                            <div class="flex items-center justify-between">
                                <div>
                                    <p class="text-xs md:text-sm text-orange-600 mb-1">总章节数</p>
                                    <p class="text-lg md:text-2xl font-bold text-orange-800">1,245</p>
                                    <p class="text-xs text-orange-500 mt-1"><i class="fas fa-arrow-up mr-1"></i>15.8%
                                        较上月</p>
                                </div>
                                <div
                                    class="w-10 h-10 md:w-12 md:h-12 bg-orange-100 rounded-full flex items-center justify-center text-orange-600">
                                    <i class="fas fa-list-alt text-base md:text-xl"></i>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 近期点击趋势 -->
                    <div class="mb-6">
                        <h3 class="text-base md:text-lg font-semibold text-gray-800 mb-4">近期点击趋势</h3>
                        <div class="bg-gray-50 p-3 md:p-4 rounded-lg">
                            <div
                                class="h-48 md:h-64 flex items-end justify-center space-x-2 md:space-x-4 overflow-x-auto pb-2">
                                <div v-for="item in clickTrend" :key="item.date"
                                    class="flex flex-col items-center flex-shrink-0">
                                    <div class="w-6 md:w-10 bg-primary rounded-t transition-all duration-300 hover:bg-primary/80"
                                        :style="{ height: normalizeClickValue(item.value) + 'px' }">
                                        <div
                                            class="absolute -top-8 left-1/2 transform -translate-x-1/2 bg-gray-800 text-white text-xs px-2 py-1 rounded opacity-0 hover:opacity-100 transition-opacity whitespace-nowrap">
                                            {{ item.value.toLocaleString() }}
                                        </div>
                                    </div>
                                    <span class="text-xs text-gray-500 mt-2 truncate w-16">{{ item.date }}</span>
                                </div>
                            </div>
                        </div>
                    </div>

                    <!-- 作品数据分布 -->
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
                        <!-- 作品点击量分布 -->
                        <div>
                            <h3 class="text-base md:text-lg font-semibold text-gray-800 mb-4">作品点击量分布</h3>
                            <div class="bg-gray-50 p-3 md:p-4 rounded-lg">
                                <div class="space-y-3">
                                    <div v-for="book in books" :key="book.id" class="flex items-center">
                                        <div class="w-24 truncate text-xs md:text-sm text-gray-700">{{ book.title }}
                                        </div>
                                        <div class="flex-1 mx-3 md:mx-4">
                                            <div class="h-1.5 md:h-2 bg-gray-200 rounded-full overflow-hidden">
                                                <div class="h-full bg-primary rounded-full transition-all duration-300"
                                                    :style="{ width: (book.clickCount / maxViews * 100) + '%' }"></div>
                                            </div>
                                        </div>
                                        <div class="w-16 md:w-20 text-right text-xs md:text-sm text-gray-600">{{
                                            (book.clickCount /
                                            10000).toFixed(1) }}万</div>
                                    </div>
                                </div>
                            </div>
                        </div>

                        <!-- 作品订阅量分布 -->
                        <div>
                            <h3 class="text-base md:text-lg font-semibold text-gray-800 mb-4">作品订阅量分布</h3>
                            <div class="bg-gray-50 p-3 md:p-4 rounded-lg">
                                <div class="space-y-3">
                                    <div v-for="book in books" :key="book.id" class="flex items-center">
                                        <div class="w-24 truncate text-xs md:text-sm text-gray-700">{{ book.title }}
                                        </div>
                                        <div class="flex-1 mx-3 md:mx-4">
                                            <div class="h-1.5 md:h-2 bg-gray-200 rounded-full overflow-hidden">
                                                <div class="h-full bg-green-500 rounded-full transition-all duration-300"
                                                    :style="{ width: (book.subscriptions / maxSubscriptions * 100) + '%' }">
                                                </div>
                                            </div>
                                        </div>
                                        <div class="w-16 md:w-20 text-right text-xs md:text-sm text-gray-600">{{
                                            book.subscriptions }}
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 账户设置页面 -->
                <div v-if="activePage === 'account-settings'" class="bg-white rounded-lg shadow-sm p-4 md:p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-lg md:text-xl font-semibold text-gray-800">账户设置</h2>
                    </div>

                    <form class="space-y-4 md:space-y-6">
                        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 md:gap-6">
                            <div>
                                <label class="block text-xs md:text-sm font-medium text-gray-700 mb-1">用户名</label>
                                <input type="text" v-model="user.username"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                            </div>

                            <div>
                                <label class="block text-xs md:text-sm font-medium text-gray-700 mb-1">邮箱</label>
                                <input type="email" v-model="user.email"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                            </div>

                            <div>
                                <label class="block text-xs md:text-sm font-medium text-gray-700 mb-1">手机号</label>
                                <input type="tel"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                            </div>

                            <div>
                                <label class="block text-xs md:text-sm font-medium text-gray-700 mb-1">性别</label>
                                <select v-model="user.sex"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm">
                                    <option value="0">男</option>
                                    <option value="1">女</option>
                                    <option value="2">保密</option>
                                </select>
                            </div>

                            <div class="md:col-span-2">
                                <label class="block text-xs md:text-sm font-medium text-gray-700 mb-1">个人简介</label>
                                <textarea v-model="user.desc" placeholder="请输入个人简介" rows="3 md:rows-4"
                                    class="w-full px-3 py-1.5 md:px-4 md:py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 resize-vertical text-sm"></textarea>
                            </div>
                        </div>

                        <div class="flex flex-wrap gap-3 md:space-x-4">
                            <button type="button"
                                class="flex-1 min-w-[100px] px-4 py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm">保存设置</button>
                            <button type="button"
                                class="flex-1 min-w-[100px] px-4 py-2 bg-gray-200 text-gray-700 rounded-md hover:bg-gray-300 transition-colors text-sm">取消</button>
                        </div>
                    </form>
                </div>
            </main>
        </div>
    </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useUserStore } from '@/stores/user.js'

// 侧边栏状态管理
const isSidebarOpen = ref(false)

// 侧边栏菜单项
const menuItems = [
    { page: 'book-list', text: '小说列表', icon: 'fa-list' },
    { page: 'book-add', text: '发布小说', icon: 'fa-plus' },
    { page: 'chapter-list', text: '章节管理', icon: 'fa-file-alt' },
    { page: 'chapter-add', text: '发布章节', icon: 'fa-edit' },
    { page: 'data-analysis', text: '数据统计', icon: 'fa-chart-bar' },
    { page: 'account-settings', text: '账户设置', icon: 'fa-user-cog' }
]

// 当前活跃页面
const activePage = ref('book-list')

// 导入Vue Router
import { useRoute } from 'vue-router';

// 导入API服务
import { bookAPI, authorAPI } from '@/api/services';

// 获取用户信息
const userStore = useUserStore();

// 使用当前登录用户的ID作为作者ID
const authorId = computed(() => userStore.userId.value);

// 模拟书籍数据
const books = ref([]);
// 小说列表分页相关
const currentBookPage = ref(1);
const bookPageSize = ref(10);
const totalBookPages = computed(() => {
    return Math.ceil(books.value.length / bookPageSize.value);
});
// 分页后的小说列表
const paginatedBooks = computed(() => {
    const start = (currentBookPage.value - 1) * bookPageSize.value;
    const end = start + bookPageSize.value;
    return books.value.slice(start, end);
});

// 当前选中的书籍
const selectedBook = ref(null);

// 模拟章节数据
const chapters = ref([]);
// 章节列表分页相关
const currentChapterPage = ref(1);
const chapterPageSize = ref(10);
const totalChapterPages = computed(() => {
    return Math.ceil(chapters.value.length / chapterPageSize.value);
});
// 分页后的章节列表
const paginatedChapters = computed(() => {
    const start = (currentChapterPage.value - 1) * chapterPageSize.value;
    const end = start + chapterPageSize.value;
    return chapters.value.slice(start, end);
});

// 点击趋势数据
const clickTrend = ref([]);

// 计算最大订阅量用于图表显示
const maxSubscriptions = computed(() => {
    return books.value.length > 0 ? Math.max(...books.value.map(book => book.subscriptions)) : 0
})

// 计算最大点击趋势值用于归一化柱状图高度
const maxClickTrend = computed(() => {
    return clickTrend.value.length > 0 ? Math.max(...clickTrend.value.map(item => item.value)) : 0
})

// 归一化点击趋势值到指定高度范围
const normalizeClickValue = (value) => {
    const maxHeight = 180 // 最大高度为180px
    return maxClickTrend.value > 0 ? (value / maxClickTrend.value) * maxHeight : 0
}

// 获取页面标题
const getPageTitle = () => {
    const item = menuItems.find(item => item.page === activePage.value)
    return item ? item.text : ''
}

// 小说列表分页控制方法
const goToBookPage = (page) => {
    if (page >= 1 && page <= totalBookPages.value) {
        currentBookPage.value = page;
    }
};

const goToPrevBookPage = () => {
    if (currentBookPage.value > 1) {
        currentBookPage.value--;
    }
};

const goToNextBookPage = () => {
    if (currentBookPage.value < totalBookPages.value) {
        currentBookPage.value++;
    }
};

// 章节列表分页控制方法
const goToChapterPage = (page) => {
    if (page >= 1 && page <= totalChapterPages.value) {
        currentChapterPage.value = page;
    }
};

const goToPrevChapterPage = () => {
    if (currentChapterPage.value > 1) {
        currentChapterPage.value--;
    }
};

const goToNextChapterPage = () => {
    if (currentChapterPage.value < totalChapterPages.value) {
        currentChapterPage.value++;
    }
};

// 处理书籍选择变化
const handleBookChange = async () => {
    // 重置章节列表页码
    currentChapterPage.value = 1;
    try {
        // 根据选中的书籍ID获取章节列表
        const chaptersResponse = await bookAPI.getChapters(selectedBook.value.id);
        chapters.value = chaptersResponse.data;
    } catch (error) {
        console.error('获取章节列表失败:', error);
    }
};

// 用户信息
const user = ref({});
// 统计数据
const maxViews = ref(0);
const maxWords = ref(0);
const maxSubs = ref(0);
const maxIncs = ref(0);

// 获取作者数据
const fetchAuthorData = async () => {
    try {
        // 获取当前作者的书籍列表数据
        const booksResponse = await authorAPI.getBooks(authorId.value);
        books.value = booksResponse.data;

        // 如果有书籍，设置默认选中的书籍并获取章节列表
        if (booksResponse.data.length > 0) {
            selectedBook.value = booksResponse.data[0];
            const chaptersResponse = await bookAPI.getChapters(booksResponse.data[0].id);
            chapters.value = chaptersResponse.data;
        }

        // 获取作者统计数据
        const statsResponse = await authorAPI.getStats(authorId.value);
        // 使用从API获取的统计数据
        maxViews.value = statsResponse.data.totalClicks;
        maxWords.value = statsResponse.data.totalWords;
        maxSubs.value = statsResponse.data.bookCount; // 临时用书籍数量代替订阅量
        maxIncs.value = statsResponse.data.bookCount; // 临时用书籍数量代替收藏量

        // 模拟点击趋势数据（实际项目中应该有对应的API）
        clickTrend.value = [
            { date: '1日', value: 12500 },
            { date: '2日', value: 14200 },
            { date: '3日', value: 13800 },
            { date: '4日', value: 16500 },
            { date: '5日', value: 18200 },
            { date: '6日', value: 15800 },
            { date: '7日', value: 21300 },
        ];

        // 获取作者用户信息
        const userResponse = await authorAPI.getInfo(authorId.value);
        user.value = userResponse.data;
    } catch (error) {
        console.error('获取数据失败:', error);
    }
};

// 获取数据
onMounted(async () => {
    // 直接获取作者数据（authorId已经从userStore中获取）
    await fetchAuthorData();
});
</script> 