<template>
    <div class="flex min-h-screen bg-gray-50">
        <!-- 侧边导航栏 -->
        <div class="w-64 bg-gray-800 text-white shadow-lg fixed inset-y-0 left-0 transform transition-transform duration-300 ease-in-out z-50"
            :class="isSidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'">
            <div class="h-16 flex items-center justify-center bg-gray-900">
                <i class="fas fa-shield-alt text-primary text-2xl mr-2"></i>
                <span class="text-xl font-bold">管理后台</span>
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
        <div v-if="isSidebarOpen" class="fixed inset-0 bg-black opacity-30 z-40 md:hidden"
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
                <div class="text-gray-600 text-sm md:text-base">管理后台 / {{ getPageTitle() }}</div>
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
                <!-- 排行榜管理页面 -->
                <div v-if="activePage === 'rank-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-xl font-semibold text-gray-800">排行榜管理</h2>
                        <div class="flex gap-2">
                            <button @click="updateRank()"
                                class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-dark transition-colors">
                                <i class="fas fa-save mr-2"></i>保存
                            </button>
                            <button @click="showAddRankModal = true"
                                class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-dark transition-colors">
                                <i class="fas fa-plus mr-2"></i>添加
                            </button>
                        </div>
                    </div>
                    
                    <!-- 排行榜类型过滤 -->
                    <div class="flex flex-wrap gap-2 mb-4">
                        <button
                            v-for="type in rankTypes"
                            :key="type.value"
                            @click="currentRankType = type.value"
                            :class="currentRankType === type.value ? 'bg-primary text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'"
                            class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
                        >
                            {{ type.label }}
                        </button>
                    </div>

                    <!-- 添加排行榜模态框 -->
                    <div v-if="showAddRankModal" class="fixed inset-0 z-50 flex items-center justify-center">
                        <!-- 遮罩层 -->
                        <div class="fixed inset-0 bg-black opacity-30"></div>
                        <!-- 对话框内容 -->
                        <div class="bg-white rounded-lg p-6 w-full max-w-md relative z-10">
                            <div class="flex justify-between items-center mb-4">
                                <h3 class="text-lg font-semibold">添加排行榜</h3>
                                <button @click="showAddRankModal = false" class="text-gray-500 hover:text-gray-700">
                                    <i class="fas fa-times text-xl"></i>
                                </button>
                            </div>
                            <div class="space-y-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">排行榜类型</label>
                                    <select v-model="newRank.rankType"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-md">
                                        <option value="click">点击榜</option>
                                        <option value="new">新书榜</option>
                                        <option value="update">更新榜</option>
                                        <option value="comment">评论榜</option>
                                    </select>
                                </div>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">小说ID</label>
                                    <input v-model.number="newRank.bookId" type="number"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-md">
                                </div>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                                    <input v-model.number="newRank.order" type="number" min="0"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-md">
                                </div>
                                <div class="flex justify-end space-x-3 pt-4">
                                    <button @click="showAddRankModal = false"
                                        class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors">取消</button>
                                    <button @click="addRank"
                                        class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-dark transition-colors">添加</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="overflow-x-auto rounded-lg shadow-md">
                        <table class="min-w-full table-fixed divide-y divide-gray-200">
                            <thead class=" bg-gray-100 rounded-t-lg">
                                <tr>
                                    <th
                                        class="px-8 py-4 text-center text-sm font-semibold uppercase tracking-wider w-24">
                                        排序</th>
                                    <th
                                        class="px-8 py-4 text-right text-sm font-semibold uppercase tracking-wider w-28">
                                        小说ID</th>
                                    <th
                                        class="px-8 py-4 text-left text-sm font-semibold uppercase tracking-wider flex-1 min-w-[200px]">
                                        小说名称</th>
                                    <th
                                        class="px-8 py-4 text-right text-sm font-semibold uppercase tracking-wider w-48">
                                        操作</th>
                                </tr>
                            </thead>
                            <draggable
                                v-if="ranks.length > 0"
                                v-model="filteredRanks"
                                item-key="id"
                                handle=".drag-handle"
                                @end="onRankDragEnd"
                                tag="tbody"
                                class="bg-white divide-y divide-gray-100"
                            >
                                <template #item="{ element: rank, index }">
                                    <tr class="hover:bg-blue-50 transition-colors duration-200">
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-gray-700 text-center">
                                            <span class="inline-flex items-center justify-center">
                                                <i class="fas fa-grip-vertical drag-handle cursor-move text-gray-400 hover:text-primary mr-2"></i>
                                                {{ index + 1 }}
                                            </span>
                                        </td>
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-gray-700 text-right font-medium">{{ rank.bookId }}</td>
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-gray-800">{{ rank.book?.title || '未知' }}</td>
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-right">
                                            <div class="flex justify-end flex-wrap gap-3">
                                                <button @click="deleteRank(rank)"
                                                    class="text-red-600 hover:text-red-800 hover:underline text-sm font-medium transition-colors">删除</button>
                                            </div>
                                        </td>
                                    </tr>
                                </template>
                            </draggable>
                            <tbody v-else class="bg-white">
                                <tr>
                                    <td colspan="4" class="px-8 py-12 text-center text-gray-500">
                                        <i class="fas fa-list text-3xl mb-3 text-gray-400"></i>
                                        <p class="text-base">暂无排行榜数据</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 分页导航 -->
                    <div class="mt-4 flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
                        v-if="ranks.length > 0">
                        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p class="text-sm text-gray-700">
                                    显示 <span class="font-medium">{{ (currentRankPage - 1) * rankPageSize + 1 }}</span> 到
                                    <span class="font-medium">{{ Math.min(currentRankPage * rankPageSize, ranks.length)
                                        }}</span> 条，共
                                    <span class="font-medium">{{ ranks.length }}</span> 条结果
                                </p>
                            </div>
                            <div>
                                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                                    aria-label="Pagination">
                                    <button @click="goToPrevRankPage" :disabled="currentRankPage === 1"
                                        class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">上一页</span>
                                        <i class="fas fa-chevron-left"></i>
                                    </button>
                                    <button v-for="page in totalRankPages" :key="page" @click="goToRankPage(page)"
                                        :class="currentRankPage === page
                                            ? 'relative inline-flex items-center z-10 bg-primary border-primary px-4 py-2 text-sm font-medium text-white focus:outline-offset-0'
                                            : 'relative inline-flex items-center border-gray-300 px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'">
                                        {{ page }}
                                    </button>
                                    <button @click="goToNextRankPage" :disabled="currentRankPage === totalRankPages"
                                        class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">下一页</span>
                                        <i class="fas fa-chevron-right"></i>
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 推荐榜管理页面 -->
                <div v-if="activePage === 'rcmd-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-xl font-semibold text-gray-800">推荐榜管理</h2>
                        <div class="flex gap-2">
                            <button @click="updateRcmd()"
                                class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-dark transition-colors">
                                <i class="fas fa-save mr-2"></i>保存
                            </button>
                            <button @click="showAddRcmdModal = true"
                                class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-dark transition-colors">
                                <i class="fas fa-plus mr-2"></i>添加
                            </button>
                        </div>
                    </div>
                    
                    <!-- 推荐榜类型过滤 -->
                    <div class="flex flex-wrap gap-2 mb-4">
                        <button
                            v-for="type in rcmdTypes"
                            :key="type.value"
                            @click="currentRcmdType = type.value"
                            :class="currentRcmdType === type.value ? 'bg-primary text-white' : 'bg-gray-200 text-gray-700 hover:bg-gray-300'"
                            class="px-4 py-2 rounded-md text-sm font-medium transition-colors"
                        >
                            {{ type.label }}
                        </button>
                    </div>

                    <!-- 添加推荐榜模态框 -->
                    <div v-if="showAddRcmdModal" class="fixed inset-0 z-50 flex items-center justify-center">
                        <!-- 遮罩层 -->
                        <div class="fixed inset-0 bg-black opacity-30"></div>
                        <!-- 对话框内容 -->
                        <div class="bg-white rounded-lg p-6 w-full max-w-md relative z-10">
                            <div class="flex justify-between items-center mb-4">
                                <h3 class="text-lg font-semibold">添加推荐榜</h3>
                                <button @click="showAddRcmdModal = false" class="text-gray-500 hover:text-gray-700">
                                    <i class="fas fa-times text-xl"></i>
                                </button>
                            </div>
                            <div class="space-y-4">
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">推荐类型</label>
                                    <select v-model="newRcmd.rcmdType"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-md">
                                        <option value="hot">热门推荐</option>
                                        <option value="top">本周排行</option>
                                        <option value="curated">精品推荐</option>
                                        <option value="featured">特别推荐</option>
                                    </select>
                                </div>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">小说ID</label>
                                    <input v-model.number="newRcmd.bookId" type="number"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-md">
                                </div>
                                <div>
                                    <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                                    <input v-model.number="newRcmd.order" type="number" min="0"
                                        class="w-full px-3 py-2 border border-gray-300 rounded-md">
                                </div>
                                <div class="flex justify-end space-x-3 pt-4">
                                    <button @click="showAddRcmdModal = false"
                                        class="px-4 py-2 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors">取消</button>
                                    <button @click="addRcmd"
                                        class="bg-primary text-white px-4 py-2 rounded-md hover:bg-primary-dark transition-colors">添加</button>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div class="overflow-x-auto rounded-lg shadow-md">
                        <table class="min-w-full table-fixed divide-y divide-gray-200">
                            <thead class="bg-gray-100 rounded-t-lg">
                                <tr>
                                    <th
                                        class="px-8 py-4 text-center text-sm font-semibold uppercase tracking-wider w-24">
                                        排序</th>
                                    <th
                                        class="px-8 py-4 text-right text-sm font-semibold uppercase tracking-wider w-28">
                                        小说ID</th>
                                    <th
                                        class="px-8 py-4 text-left text-sm font-semibold uppercase tracking-wider flex-1 min-w-[200px]">
                                        小说名称</th>
                                    <th
                                        class="px-8 py-4 text-right text-sm font-semibold uppercase tracking-wider w-48">
                                        操作</th>
                                </tr>
                            </thead>
                            <draggable
                                v-if="rcmds.length > 0"
                                v-model="filteredRcmds"
                                item-key="id"
                                handle=".drag-handle"
                                @end="onRcmdDragEnd"
                                tag="tbody"
                                class="bg-white divide-y divide-gray-100"
                            >
                                <template #item="{ element: rcmd, index }">
                                    <tr class="hover:bg-blue-50 transition-colors duration-200">
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-gray-700 text-center">
                                            <span class="inline-flex items-center justify-center">
                                                <i class="fas fa-grip-vertical drag-handle cursor-move text-gray-400 hover:text-primary mr-2"></i>
                                                {{ index + 1 }}
                                            </span>
                                        </td>
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-gray-700 text-right font-medium">{{ rcmd.bookId }}</td>
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-gray-800">{{ rcmd.book?.title || '未知' }}</td>
                                        <td class="px-8 py-5 whitespace-nowrap text-sm text-right">
                                            <div class="flex justify-end flex-wrap gap-3">
                                                <button @click="deleteRcmd(rcmd)"
                                                    class="text-red-600 hover:text-red-800 hover:underline text-sm font-medium transition-colors">删除</button>
                                            </div>
                                        </td>
                                    </tr>
                                </template>
                            </draggable>
                            <tbody v-else class="bg-white">
                                <tr>
                                    <td colspan="4" class="px-8 py-12 text-center text-gray-500">
                                        <i class="fas fa-star text-3xl mb-3 text-gray-400"></i>
                                        <p class="text-base">暂无推荐榜数据</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 分页导航 -->
                    <div class="mt-4 flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
                        v-if="rcmds.length > 0">
                        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p class="text-sm text-gray-700">
                                    显示 <span class="font-medium">{{ (currentRcmdPage - 1) * rcmdPageSize + 1 }}</span> 到
                                    <span class="font-medium">{{ Math.min(currentRcmdPage * rcmdPageSize, rcmds.length)
                                        }}</span> 条，共
                                    <span class="font-medium">{{ rcmds.length }}</span> 条结果
                                </p>
                            </div>
                            <div>
                                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                                    aria-label="Pagination">
                                    <button @click="goToPrevRcmdPage" :disabled="currentRcmdPage === 1"
                                        class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">上一页</span>
                                        <i class="fas fa-chevron-left"></i>
                                    </button>
                                    <button v-for="page in totalRcmdPages" :key="page" @click="goToRcmdPage(page)"
                                        :class="currentRcmdPage === page
                                            ? 'relative inline-flex items-center z-10 bg-primary border-primary px-4 py-2 text-sm font-medium text-white focus:outline-offset-0'
                                            : 'relative inline-flex items-center border-gray-300 px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'">
                                        {{ page }}
                                    </button>
                                    <button @click="goToNextRcmdPage" :disabled="currentRcmdPage === totalRcmdPages"
                                        class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">下一页</span>
                                        <i class="fas fa-chevron-right"></i>
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 用户管理页面 -->
                <div v-if="activePage === 'user-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-xl font-semibold text-gray-800">用户管理</h2>
                        <div class="relative">
                            <input
                                v-model="userSearchKeyword"
                                type="text"
                                placeholder="搜索用户名或邮箱"
                                class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                            >
                            <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
                        </div>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="min-w-full table-fixed divide-y divide-gray-200 border border-gray-200 rounded-lg">
                            <thead class="bg-gray-100 rounded-t-lg">
                                <tr>
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-20">
                                        ID</th>
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-32">
                                        用户名</th>
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-32">
                                        昵称</th>
                                    <th
                                        class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-48">
                                        邮箱</th>
                                    <th
                                        class="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-24">
                                        状态</th>
                                    <th
                                        class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-32">
                                        操作</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200 rounded-b-lg">
                                <tr v-for="user in paginatedUsers" v-if="users?.length > 0" :key="user.id"
                                    class="hover:bg-gray-50 transition-colors">
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.id }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.username
                                        }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.nickname
                                        }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ user.email }}
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-center">
                                        <span :class="user.status === 20 ? 'text-red-500' : 'text-green-500'">
                                            {{ user.status === 20 ? '已禁言' : '正常' }}
                                        </span>
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-right">
                                        <div class="flex justify-end flex-wrap gap-2">
                                            <button @click="toggleUserStatus(user)"
                                                class="text-blue-600 hover:text-blue-800 text-xs md:text-sm">
                                                {{ user.status === 20 ? '解除禁言' : '禁言' }}
                                            </button>
                                        </div>
                                    </td>
                                </tr>
                                <tr v-else>
                                    <td colspan="6" class="px-4 py-8 text-center text-gray-500">
                                        <i class="fas fa-users text-2xl mb-2"></i>
                                        <p>暂无用户数据</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 分页导航 -->
                    <div class="mt-4 flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
                        v-if="users.length > 0">
                        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p class="text-sm text-gray-700">
                                    显示 <span class="font-medium">{{ (currentUserPage - 1) * userPageSize + 1 }}</span> 到
                                    <span class="font-medium">{{ Math.min(currentUserPage * userPageSize, totalUsers) }}</span> 条，共
                                    <span class="font-medium">{{ totalUsers }}</span> 条结果
                                </p>
                            </div>
                            <div>
                                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                                    aria-label="Pagination">
                                    <button @click="goToPrevUserPage" :disabled="currentUserPage === 1"
                                        class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">上一页</span>
                                        <i class="fas fa-chevron-left"></i>
                                    </button>
                                    <button v-for="page in totalUserPages" :key="page" @click="goToUserPage(page)"
                                        :class="currentUserPage === page
                                            ? 'relative inline-flex items-center z-10 bg-primary border-primary px-4 py-2 text-sm font-medium text-white focus:outline-offset-0'
                                            : 'relative inline-flex items-center border-gray-300 px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'">
                                        {{ page }}
                                    </button>
                                    <button @click="goToNextUserPage" :disabled="currentUserPage === totalUserPages"
                                        class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">下一页</span>
                                        <i class="fas fa-chevron-right"></i>
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 评论管理页面 -->
                <div v-if="activePage === 'comment-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-xl font-semibold text-gray-800">评论管理</h2>
                        <div class="relative">
                            <input
                                v-model="commentSearchKeyword"
                                type="text"
                                placeholder="搜索评论内容"
                                class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                            >
                            <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
                        </div>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="min-w-full table-fixed border border-gray-200 rounded-lg">
                            <thead class="bg-gray-100 rounded-t-lg">
                                <tr>
                                    <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-20">
                                        ID</th>
                                    <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-24">
                                        小说名称</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-32">
                                        用户名</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm flex-1 min-w-[300px]">
                                        内容</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-48">
                                        创建时间</th>
                                    <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-24">
                                        操作</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200 rounded-b-lg">
                                <tr v-for="comment in paginatedComments" v-if="comments.length > 0" :key="comment.id"
                                    class="hover:bg-gray-50 transition-colors">
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-right">{{ comment.id }}
                                    </td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-right">{{ comment.bookTitle
                                        }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{
                                        comment.userNickname || '未知' }}</td>
                                    <td class="px-6 py-4 text-sm text-gray-500 max-w-xs truncate">{{ comment.content
                                        }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{
                                        formatDate(comment.createTime) }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-right">
                                        <div class="flex justify-end gap-2">
                                            <button @click="deleteComment(comment.id)"
                                                class="text-red-600 hover:text-red-800 text-xs md:text-sm">删除</button>
                                        </div>
                                    </td>
                                </tr>
                                <tr v-else>
                                    <td colspan="6" class="px-6 py-8 text-center text-gray-500">
                                        <i class="fas fa-comments text-2xl mb-2"></i>
                                        <p>暂无评论数据</p>
                                    </td>
                                </tr>
                            </tbody>
                        </table>
                    </div>

                    <!-- 分页导航 -->
                    <div class="mt-4 flex items-center justify-between border-t border-gray-200 bg-white px-4 py-3 sm:px-6"
                        v-if="comments.length > 0">
                        <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
                            <div>
                                <p class="text-sm text-gray-700">
                                    显示 <span class="font-medium">{{ (currentCommentPage - 1) * commentPageSize + 1 }}</span> 到
                                    <span class="font-medium">{{ Math.min(currentCommentPage * commentPageSize, totalComments) }}</span> 条，共
                                    <span class="font-medium">{{ totalComments }}</span> 条结果
                                </p>
                            </div>
                            <div>
                                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm"
                                    aria-label="Pagination">
                                    <button @click="goToPrevCommentPage" :disabled="currentCommentPage === 1"
                                        class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">上一页</span>
                                        <i class="fas fa-chevron-left"></i>
                                    </button>
                                    <button v-for="page in totalCommentPages" :key="page" @click="goToCommentPage(page)"
                                        :class="currentCommentPage === page
                                            ? 'relative inline-flex items-center z-10 bg-primary border-primary px-4 py-2 text-sm font-medium text-white focus:outline-offset-0'
                                            : 'relative inline-flex items-center border-gray-300 px-4 py-2 text-sm font-medium text-gray-500 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'">
                                        {{ page }}
                                    </button>
                                    <button @click="goToNextCommentPage"
                                        :disabled="currentCommentPage === totalCommentPages"
                                        class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0">
                                        <span class="sr-only">下一页</span>
                                        <i class="fas fa-chevron-right"></i>
                                    </button>
                                </nav>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- 小说管理页面 -->
                <div v-if="activePage === 'book-list'" class="bg-white rounded-lg shadow-sm p-6">
                    <div class="flex justify-between items-center mb-6 pb-4 border-b border-gray-200">
                        <h2 class="text-xl font-semibold text-gray-800">小说管理</h2>
                        <div class="relative">
                            <input
                                v-model="bookSearchKeyword"
                                type="text"
                                placeholder="搜索书名或作者"
                                class="pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-primary"
                            >
                            <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400"></i>
                        </div>
                    </div>

                    <div class="overflow-x-auto">
                        <table class="min-w-full table-fixed border border-gray-200 rounded-lg">
                            <thead class="bg-gray-100 rounded-t-lg">
                                <tr>
                                    <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-20">
                                        ID</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm flex-1 min-w-[200px]">
                                        书名</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-32">
                                        作者</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-24">
                                        分类</th>
                                    <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-24">
                                        字数</th>
                                    <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-48">
                                        更新时间</th>
                                    <th class="px-6 py-4 text-right text-xs font-medium text-gray-500 uppercase tracking-wider md:text-sm w-28">
                                        操作</th>
                                </tr>
                            </thead>
                            <tbody class="bg-white divide-y divide-gray-200 rounded-b-lg">
                                <tr v-for="book in paginatedBooks" v-if="books.length > 0" :key="book.id"
                                    class="hover:bg-gray-50 transition-colors">
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-right">{{ book.id }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ book.title
                                    }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ book.author
                                    }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ book.category
                                        }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-right">{{ book.wordCount
                                        }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{
                                        formatDate(book.updateTime) }}</td>
                                    <td class="px-6 py-4 whitespace-nowrap text-sm text-right">
                                        <div class="flex justify-end gap-2">
                                            <button @click="deleteBook(book.id)"
                                                class="text-red-600 hover:text-red-800 text-xs md:text-sm">删除</button>
                                        </div>
                                    </td>
                                </tr>
                                <tr v-else>
                                    <td colspan="7" class="px-6 py-8 text-center text-gray-500">
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
                                    <span class="font-medium">{{ Math.min(currentBookPage * bookPageSize, totalBooks) }}</span> 条，共
                                    <span class="font-medium">{{ totalBooks }}</span> 条结果
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
            </main>
        </div>
    </div>

    <!-- Notice组件 -->
    <Notice
        :visible="noticeVisible"
        :type="noticeType"
        :title="noticeTitle"
        :message="noticeMessage"
        @close="closeNotice"
    />
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useUserStore } from '@/stores/user.js'
import { adminAPI, rankAPI, rcmdAPI, bookAPI, userAPI } from '@/api/services'
import draggable from 'vuedraggable'
import Notice from '@/components/Notice.vue'

// 侧边栏状态管理
const isSidebarOpen = ref(false)

// 侧边栏菜单项
const menuItems = [
    { page: 'rank-list', text: '排行榜管理', icon: 'fa-list' },
    { page: 'rcmd-list', text: '推荐榜管理', icon: 'fa-star' },
    { page: 'user-list', text: '用户管理', icon: 'fa-users' },
    { page: 'comment-list', text: '评论管理', icon: 'fa-comments' },
    { page: 'book-list', text: '小说管理', icon: 'fa-book-open' }
]

// 当前活跃页面
const activePage = ref('rank-list')

// 获取用户信息
const userStore = useUserStore();
const user = ref({});

// ==================== 排行榜管理 ====================
// 排行榜数据
const ranks = ref([]);
// 当前选中的排行榜类型
const currentRankType = ref('click');
// 排行榜类型选项
const rankTypes = [
    { value: 'click', label: '点击榜' },
    { value: 'new', label: '新书榜' },
    { value: 'update', label: '更新榜' },
    { value: 'comment', label: '评论榜' }
];
// 过滤后的排行榜数据
const filteredRanks = computed(() => {
    return ranks.value.filter(r => r.rankType === currentRankType.value);
});
// 排行榜列表分页相关
const currentRankPage = ref(1);
const rankPageSize = ref(10);
const totalRankPages = computed(() => {
    return Math.ceil(filteredRanks.value.length / rankPageSize.value);
});

// 添加排行榜相关
const showAddRankModal = ref(false);
const newRank = ref({
    rankType: 'click',
    bookId: null,
    order: 0
});

// 排行榜分页控制方法
const goToRankPage = (page) => {
    if (page >= 1 && page <= totalRankPages.value) {
        currentRankPage.value = page;
    }
};

const goToPrevRankPage = () => {
    if (currentRankPage.value > 1) {
        currentRankPage.value--;
    }
};

const goToNextRankPage = () => {
    if (currentRankPage.value < totalRankPages.value) {
        currentRankPage.value++;
    }
};

// 获取排行榜数据
const fetchRanks = async () => {
    try {
        // 获取所有类型的排行榜数据并合并
        const [clickRank, newBookRank, updateRank, commentRank] = await Promise.all([
            rankAPI.getClickRank(),
            rankAPI.getNewBookRank(),
            rankAPI.getUpdateRank(),
            rankAPI.getCommentRank()
        ]);
        ranks.value = [...clickRank.data, ...newBookRank.data, ...updateRank.data, ...commentRank.data];
    } catch (error) {
        console.error('获取排行榜数据失败:', error);
    }
};

// 拖拽结束事件处理函数 - 更新排行榜顺序
const onRankDragEnd = () => {
    // 更新每个项目的order字段，使其与当前位置一致
    filteredRanks.value.forEach((rank, index) => {
        rank.order = index;
    });
};

// 更新排行榜数据（支持批量更新顺序）
const updateRank = async () => {
    try {
        // 使用批量更新接口，传递当前类型的所有排行榜数据
        await rankAPI.updateRanks(currentRankType.value, filteredRanks.value);
        showNotice('success', '成功', '排序更新成功');
    } catch (error) {
        console.error('更新排行榜数据失败:', error);
        showNotice('error', '失败', '更新失败');
    }
};

// 删除排行榜数据
const deleteRank = async (rank) => {
    if (!confirm('确定要删除这条排行榜数据吗？')) return;
    try {
        await rankAPI.deleteRank(rank.rankType, rank.id);
        fetchRanks();
        alert('删除成功');
    } catch (error) {
        console.error('删除排行榜数据失败:', error);
        alert('删除失败');
    }
};

// 添加排行榜数据
const addRank = async () => {
    if (!newRank.value.bookId) {
        alert('请输入小说ID');
        return;
    }
    try {
        // 设置新添加项目的排序为当前类型列表中的最后一位
        newRank.value.order = filteredRanks.value.length;
        await rankAPI.addRank(newRank.value.rankType, newRank.value);
        showAddRankModal.value = false;
        // 重置表单
        newRank.value = {
            rankType: currentRankType.value, // 默认使用当前选中的类型
            bookId: null,
            order: 0
        };
        fetchRanks();
        alert('添加成功');
    } catch (error) {
        console.error('添加排行榜数据失败:', error);
        alert('添加失败');
    }
};

// ==================== 推荐榜管理 ====================
// 推荐榜数据
const rcmds = ref([]);
// 当前选中的推荐榜类型
const currentRcmdType = ref('hot');
// 推荐榜类型选项
const rcmdTypes = [
    { value: 'hot', label: '热门推荐' },
    { value: 'top', label: '本周排行' },
    { value: 'curated', label: '精品推荐' },
    { value: 'featured', label: '特别推荐' }
];
// 过滤后的推荐榜数据
const filteredRcmds = computed(() => {
    return rcmds.value.filter(r => r.rcmdType === currentRcmdType.value);
});
// 推荐榜列表分页相关
const currentRcmdPage = ref(1);
const rcmdPageSize = ref(10);
const totalRcmdPages = computed(() => {
    return Math.ceil(filteredRcmds.value.length / rcmdPageSize.value);
});

// 添加推荐榜相关
const showAddRcmdModal = ref(false);
const newRcmd = ref({
    rcmdType: 'hot',
    bookId: null,
    order: 0
});

// 推荐榜分页控制方法
const goToRcmdPage = (page) => {
    if (page >= 1 && page <= totalRcmdPages.value) {
        currentRcmdPage.value = page;
    }
};

const goToPrevRcmdPage = () => {
    if (currentRcmdPage.value > 1) {
        currentRcmdPage.value--;
    }
};

const goToNextRcmdPage = () => {
    if (currentRcmdPage.value < totalRcmdPages.value) {
        currentRcmdPage.value++;
    }
};

// 获取推荐榜数据
const fetchRcmds = async () => {
    try {
        // 获取所有类型的推荐榜数据并合并
        const [hotBooks, topBooks, curatedBooks, featuredBooks] = await Promise.all([
            rcmdAPI.getHotBooks(),
            rcmdAPI.getTopBooks(),
            rcmdAPI.getCuratedBooks(),
            rcmdAPI.getFeaturedBooks()
        ]);
        rcmds.value = [...hotBooks.data, ...topBooks.data, ...curatedBooks.data, ...featuredBooks.data];
    } catch (error) {
        console.error('获取推荐榜数据失败:', error);
    }
};

// 拖拽结束事件处理函数 - 更新推荐榜顺序
const onRcmdDragEnd = () => {
    // 更新每个项目的order字段，使其与当前位置一致
    filteredRcmds.value.forEach((rcmd, index) => {
        rcmd.order = index;
    });
};

// 更新推荐榜数据（支持批量更新顺序）
const updateRcmd = async () => {
    try {
        // 使用批量更新接口，传递当前类型的所有推荐榜数据
        await rcmdAPI.updateRcmds(currentRcmdType.value, filteredRcmds.value);
        showNotice('success', '成功', '排序更新成功');
    } catch (error) {
        console.error('更新推荐榜数据失败:', error);
        showNotice('error', '失败', '更新失败');
    }
};

// 删除推荐榜数据
const deleteRcmd = async (rcmd) => {
    if (!confirm('确定要删除这条推荐榜数据吗？')) return;
    try {
        await rcmdAPI.deleteRcmd(rcmd.rcmdType, rcmd.id);
        fetchRcmds();
        alert('删除成功');
    } catch (error) {
        console.error('删除推荐榜数据失败:', error);
        alert('删除失败');
    }
};

// 添加推荐榜数据
const addRcmd = async () => {
    if (!newRcmd.value.bookId) {
        alert('请输入小说ID');
        return;
    }
    try {
        // 设置新添加项目的排序为当前类型列表中的最后一位
        newRcmd.value.order = filteredRcmds.value.length;
        await rcmdAPI.addRcmd(newRcmd.value.rcmdType, newRcmd.value);
        showAddRcmdModal.value = false;
        // 重置表单
        newRcmd.value = {
            rcmdType: currentRcmdType.value, // 默认使用当前选中的类型
            bookId: null,
            order: 0
        };
        fetchRcmds();
        alert('添加成功');
    } catch (error) {
        console.error('添加推荐榜数据失败:', error);
        alert('添加失败');
    }
};

// ==================== 用户管理 ====================
// 用户数据
const users = ref([]);
// 用户总数
const totalUsers = ref(0);
// 用户搜索关键词
const userSearchKeyword = ref('');
// 用户列表分页相关
const currentUserPage = ref(1);
const userPageSize = ref(10);
const totalUserPages = computed(() => {
    return Math.ceil(totalUsers.value / userPageSize.value);
});
// 分页后的用户列表 - 由于API已经返回分页数据，直接使用即可
const paginatedUsers = computed(() => {
    return users.value;
});

// 用户分页控制方法
const goToUserPage = (page) => {
    if (page >= 1 && page <= totalUserPages.value) {
        currentUserPage.value = page;
        fetchUsers();
    }
};

const goToPrevUserPage = () => {
    if (currentUserPage.value > 1) {
        currentUserPage.value--;
        fetchUsers();
    }
};

const goToNextUserPage = () => {
    if (currentUserPage.value < totalUserPages.value) {
        currentUserPage.value++;
        fetchUsers();
    }
};

// 获取用户数据
const fetchUsers = async () => {
    try {
        const response = await userAPI.search(userSearchKeyword.value, currentUserPage.value, userPageSize.value);
        users.value = response.data.users;
        totalUsers.value = response.data.total;
    } catch (error) {
        console.error('获取用户数据失败:', error);
    }
};

// 切换用户状态（禁言/解除禁言）
const toggleUserStatus = async (user) => {
    const newStatus = user.status === 20 ? 1 : 20;
    const action = user.status === 20 ? '解除禁言' : '禁言';
    if (!confirm(`确定要${action}用户 ${user.username} 吗？`)) return;
    try {
        await adminAPI.updateUserStatus(user.id, { status: newStatus });
        user.status = newStatus;
        alert(`${action}成功`);
    } catch (error) {
        console.error(`${action}用户失败:`, error);
        alert(`${action}失败`);
    }
};

// ==================== 评论管理 ====================
// 评论数据
const comments = ref([]);
// 评论总数
const totalComments = ref(0);
// 评论搜索关键词
const commentSearchKeyword = ref('');
// 评论列表分页相关
const currentCommentPage = ref(1);
const commentPageSize = ref(10);
const totalCommentPages = computed(() => {
    return Math.ceil(totalComments.value / commentPageSize.value);
});
// 分页后的评论列表 - 由于API已经返回分页数据，直接使用即可
const paginatedComments = computed(() => {
    return comments.value;
});

// 评论分页控制方法
const goToCommentPage = (page) => {
    if (page >= 1 && page <= totalCommentPages.value) {
        currentCommentPage.value = page;
        fetchComments();
    }
};

const goToPrevCommentPage = () => {
    if (currentCommentPage.value > 1) {
        currentCommentPage.value--;
        fetchComments();
    }
};

const goToNextCommentPage = () => {
    if (currentCommentPage.value < totalCommentPages.value) {
        currentCommentPage.value++;
        fetchComments();
    }
};

// 获取评论数据
const fetchComments = async () => {
    try {
        const response = await adminAPI.searchComments(commentSearchKeyword.value, currentCommentPage.value, commentPageSize.value);
        comments.value = response.data.comments;
        totalComments.value = response.data.total;
    } catch (error) {
        console.error('获取评论数据失败:', error);
    }
};

// 删除评论
const deleteComment = async (id) => {
    if (!confirm('确定要删除这条评论吗？')) return;
    try {
        await adminAPI.deleteComment(id);
        comments.value = comments.value.filter(comment => comment.id !== id);
        alert('删除成功');
    } catch (error) {
        console.error('删除评论失败:', error);
        alert('删除失败');
    }
};

// ==================== 小说管理 ====================
// 小说数据
const books = ref([]);
// 小说总数
const totalBooks = ref(0);
// 小说搜索关键词
const bookSearchKeyword = ref('');
// 小说列表分页相关
const currentBookPage = ref(1);
const bookPageSize = ref(10);
const totalBookPages = computed(() => {
    return Math.ceil(totalBooks.value / bookPageSize.value);
});
// 分页后的小说列表 - 由于API已经返回分页数据，直接使用即可
const paginatedBooks = computed(() => {
    return books.value;
});

// 小说分页控制方法
const goToBookPage = (page) => {
    if (page >= 1 && page <= totalBookPages.value) {
        currentBookPage.value = page;
        fetchBooks();
    }
};

const goToPrevBookPage = () => {
    if (currentBookPage.value > 1) {
        currentBookPage.value--;
        fetchBooks();
    }
};

const goToNextBookPage = () => {
    if (currentBookPage.value < totalBookPages.value) {
        currentBookPage.value++;
        fetchBooks();
    }
};

// 获取小说数据
const fetchBooks = async () => {
    try {
        const response = await bookAPI.search(bookSearchKeyword.value, currentBookPage.value, bookPageSize.value);
        books.value = response.data.books;
        totalBooks.value = response.data.total;
    } catch (error) {
        console.error('获取小说数据失败:', error);
    }
};

// 删除小说
const deleteBook = async (id) => {
    if (!confirm('确定要删除这本小说吗？此操作不可恢复！')) return;
    try {
        await adminAPI.deleteBook(id);
        books.value = books.value.filter(book => book.id !== id);
        alert('删除成功');
    } catch (error) {
        console.error('删除小说失败:', error);
        alert('删除失败');
    }
};

// ==================== 通用方法 ====================
// 获取页面标题
const getPageTitle = () => {
    const item = menuItems.find(item => item.page === activePage.value)
    return item ? item.text : ''
}

// 格式化日期
const formatDate = (date) => {
    if (!date) return '';
    const d = new Date(date);
    return d.toLocaleString();
}

// Notice组件状态
const noticeVisible = ref(false)
const noticeType = ref('info')
const noticeTitle = ref('提示')
const noticeMessage = ref('')

// 显示提示
const showNotice = (type, title, message) => {
  noticeType.value = type
  noticeTitle.value = title
  noticeMessage.value = message
  noticeVisible.value = true
}

// 关闭提示
const closeNotice = () => {
  noticeVisible.value = false
}

// 页面切换时自动加载对应数据
const loadData = () => {
    switch (activePage.value) {
        case 'rank-list':
            fetchRanks();
            break;
        case 'rcmd-list':
            fetchRcmds();
            break;
        case 'user-list':
            fetchUsers();
            break;
        case 'comment-list':
            fetchComments();
            break;
        case 'book-list':
            fetchBooks();
            break;
    }
};

// 监听页面变化
const unwatch = watch(activePage, () => {
    loadData();
});

// 监听搜索关键词变化
watch(userSearchKeyword, () => {
    currentUserPage.value = 1;
    fetchUsers();
});

watch(commentSearchKeyword, () => {
    currentCommentPage.value = 1;
    fetchComments();
});

watch(bookSearchKeyword, () => {
    currentBookPage.value = 1;
    fetchBooks();
});

// 页面加载时初始化数据
onMounted(() => {
    user.value = userStore.user.value;
    loadData();
});

// 页面卸载时清理监听器
onUnmounted(() => {
    unwatch();
});
</script>