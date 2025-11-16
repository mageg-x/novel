<template>
    <Header />
    <!-- 小屏幕 -->
    <div class="block md:hidden mx-4">
        <!-- 快捷入口 -->
        <div class="grid grid-cols-2 gap-3 mt-2">
            <router-link to="/class" class="bg-white rounded-lg shadow-sm p-3 flex items-center gap-3 cursor-pointer hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-full bg-gradient-to-r from-emerald-500 to-teal-500 flex items-center justify-center text-white">
                    <i class="fas fa-book-open text-lg"></i>
                </div>
                <div>
                    <h4 class="text-sm font-medium">所有作品</h4>
                    <p class="text-xs text-gray-500">浏览全部小说</p>
                </div>
            </router-link>
            <router-link to="/rank" class="bg-white rounded-lg shadow-sm p-3 flex items-center gap-3 cursor-pointer hover:shadow-md transition-shadow">
                <div class="w-12 h-12 rounded-full bg-gradient-to-r from-orange-500 to-amber-500 flex items-center justify-center text-white">
                    <i class="fas fa-chart-line text-lg"></i>
                </div>
                <div>
                    <h4 class="text-sm font-medium">排行榜</h4>
                    <p class="text-xs text-gray-500">热门书籍排行</p>
                </div>
            </router-link>
        </div>
        
        <!-- 精品推荐 - 轮播效果 -->
        <section class="mb-6 mt-4">
        <h3 class="text-lg font-bold mb-3 border-b-2 border-primary pb-1 inline-block">精品推荐</h3>
        <div class="relative">
            <!-- 轮播容器 -->
            <div class="overflow-hidden rounded-lg">
                <div class="flex transition-transform duration-500 ease-in-out" :style="{ transform: `translateX(-${featuredCurrentIndex * 100}%)` }">
                    <div v-for="(book, index) in featuredBooks" :key="`featured-${index}`" class="w-full flex-shrink-0">
                        <div class="bg-white rounded-lg overflow-hidden shadow-sm cursor-pointer" @click="navigateToBook(book.id)">
                            <img :src="book.cover" :alt="book.title" class="w-full h-40 object-cover" />
                            <div class="p-3">
                                <h4 class="text-sm font-medium line-clamp-2">{{ book.title }}</h4>
                                <p class="text-xs text-gray-500 mt-1 line-clamp-3">{{ book.description }}</p>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
            
            <!-- 轮播指示器 -->
            <div class="flex justify-center mt-2 gap-1">
                <span v-for="(book, index) in featuredBooks" :key="`featured-dot-${index}`"
                    @click="featuredCurrentIndex = index"
                    class="w-2 h-2 rounded-full transition-all cursor-pointer"
                    :class="featuredCurrentIndex === index ? 'bg-primary' : 'bg-gray-300'">
                </span>
            </div>
        </div>
        </section>

        <!-- 热门推荐 -->
        <section class="mb-6">
        <h3 class="text-lg font-bold mb-3 border-b-2 border-primary pb-1 inline-block">热门推荐</h3>
        <div class="space-y-3">
            <div v-for="(book, index) in hotBooks" :key="`recommend-${index}`" class="flex gap-2 p-2 bg-white rounded-lg shadow-sm cursor-pointer" @click="navigateToBook(book.id)">
            <img :src="book.cover" :alt="book.title" class="w-16 h-24 object-cover rounded" />
            <div class="flex-1">
                <h4 class="text-sm font-medium line-clamp-2">{{ book.title }}</h4>
                <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ book.description }}</p>
                <p class="text-xs text-gray-400 mt-1">作者：{{ book.author }}</p>
            </div>
            </div>
        </div>
        </section>

        <!-- 最新更新 -->
        <section class="mb-6">
        <div class="flex justify-between items-center mb-3">
            <h3 class="text-lg font-bold">最新更新</h3>
            <a href="#" class="text-primary text-sm hover:underline">更多 ></a>
        </div>
        <div class="bg-white rounded-lg shadow-sm overflow-hidden">
            <ul class="divide-y divide-gray-200">
            <li v-for="(item, index) in latestBooks" :key="index" class="px-3 py-2 flex justify-between items-start cursor-pointer" @click="navigateToBook(item.id)">
                <div class="flex-1">
                <h4 class="text-sm font-medium line-clamp-1">{{ item.title }}</h4>
                <p class="text-xs text-gray-500 mt-1 line-clamp-2">{{ item.latestChapter }}</p>
                <p class="text-xs text-gray-400 mt-1">作者：{{ item.author }}</p>
                </div>
                <span class="text-xs text-orange-500 font-medium">{{ item.updateTime }}</span>
            </li>
            </ul>
        </div>
        </section>
    </div>


    <!-- 大屏幕 -->
    <div class="hidden md:block container mx-auto px-4">
        <!-- 重点推荐区 -->
        <section class="my-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6 md:items-start">
                <!-- 主推荐 -->
                <div v-if="currentBook" class="w-full md:w-2/3 bg-white rounded-xl shadow-sm overflow-hidden card-hover" @click="navigateToBook(currentBook.id)">
                    <div class="flex flex-col md:flex-row">
                        <!-- 左侧：主图 + 小图轮播 -->
                        <div class="w-full md:w-1/3 relative">
                            <!-- 主图 -->
                            <img :src="currentBook.cover" :alt="currentBook.title" class="w-full h-full object-cover cursor-pointer" />
                            <div v-if="currentBook.source"
                                class="absolute top-2 left-2 bg-red-500 text-white text-xs px-2 py-1 rounded">
                                {{ currentBook.source }}
                            </div>

                            <!-- 竖向小图轮播（仅当有多本书时显示） -->
                            <div v-if="mainBooks.length > 1" class="absolute top-0 right-0 flex flex-col">
                                <button v-for="(book, index) in mainBooks" :key="`thumb-${index}`"
                                    @mouseenter="setCurrentIndex(index)"
                                    @click.stop
                                    class="w-12 h-16 border-2 rounded-none overflow-hidden transition-all"
                                    :class="currentIndex === index ? 'border-primary' : 'border-white opacity-70 hover:opacity-100'">
                                    <img :src="book.cover" :alt="book.title" class="w-full h-full object-cover" />
                                </button>
                            </div>
                        </div>

                        <!-- 右侧：书籍信息 -->
                        <div class="w-full md:w-2/3 p-6">
                            <h2 class="text-[clamp(1.2rem,3vw,1.8rem)] font-bold mb-2 text-shadow cursor-pointer">
                                {{ currentBook.title }}
                            </h2>
                            <p class="text-gray-600 mb-4 line-clamp-7">
                                {{ currentBook.description }}
                            </p>
                            <div class="mb-4">
                                <span class="text-sm text-gray-500">简介：</span>
                                <span class="text-sm">{{ currentBook.tagline || '暂无简介' }}</span>
                            </div>
                            <a href="#"
                                class="inline-block bg-primary hover:bg-dark text-white px-4 py-2 rounded-md transition-colors cursor-pointer" @click.stop="navigateToBook(currentBook.id)">
                                开始阅读
                            </a>
                        </div>
                    </div>
                </div>
                <div v-else class="w-full md:w-2/3 bg-white rounded-xl shadow-sm overflow-hidden card-hover p-10 text-center">
                    <p class="text-gray-500">加载中...</p>
                </div>

                <!-- 本周强推 -->
                <div class="w-full md:w-1/3">
                    <RankList title="本周强推" :books="weeklyBooks" />
                </div>
            </div>
        </section>

        <!-- 热门推荐和点击榜单 -->
        <section class="mb-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6">
                <!-- 热门推荐 -->
                <div class="w-full md:w-2/3">
                    <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">热门推荐</h3>
                    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div v-for="(book, index) in hotBooks" :key="index"
                            class="bg-white rounded-xl shadow-sm p-4 card-hover cursor-pointer" @click="navigateToBook(book.id)">
                            <div class="flex gap-4">
                                <img :src="book.cover" :alt="book.title" class="w-20 h-28 object-cover rounded" />
                                <div>
                                    <h4 class="font-medium">{{ book.title }}</h4>
                                    <p class="text-sm text-gray-600 line-clamp-2 mt-1">
                                        {{ book.description }}
                                    </p>
                                    <p class="text-xs text-gray-500 mt-2">作者：{{ book.author }}</p>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
                <!-- 点击榜单 -->
                <div class="w-full md:w-1/3">
                    <RankList title="点击榜单" :books="clickBooks" :show-more="true" more-link="/rank#clickest" />
                </div>
            </div>
        </section>

        <!-- 精品推荐和新书榜单 -->
        <section class="mb-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6">
                <!-- 精品推荐 -->
                <div class="w-full md:w-2/3">
                    <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">精品推荐</h3>
                    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
                        <div v-for="(book, index) in featuredBooks" :key="`featured-${index}`"
                            class="bg-white rounded-xl shadow-sm overflow-hidden card-hover cursor-pointer" @click="navigateToBook(book.id)">
                            <img :src="book.cover" :alt="book.title" class="w-full h-48 object-cover" />
                            <div class="p-3">
                                <h4 class="font-medium line-clamp-2">{{ book.title }}</h4>
                                <p class="text-xs text-gray-500 mt-1 line-clamp-2">
                                    {{ book.description }}
                                </p>
                            </div>
                        </div>
                    </div>
                </div>
                <!-- 新书榜单 -->
                <div class="w-full md:w-1/3">
                    <RankList title="新书榜单" :books="newBooks" :show-more="true" more-link="/rank#newest" />
                </div>
            </div>
        </section>


        <!-- 最新更新和亚榜榜单 -->
        <section class="mb-10 w-5xl mx-auto">
            <div class="flex flex-col md:flex-row gap-6">
            <!-- 最新更新 -->
            <div class="w-full md:w-2/3">
                <h3 class="text-lg font-bold mb-4 border-b-2 border-primary pb-1 inline-block">最新更新</h3>
                <div class="bg-white rounded-xl shadow-sm overflow-hidden">
                <div class="overflow-x-auto">
                    <table class="w-full min-w-max">
                    <thead>
                        <tr class="bg-gray-50">
                        <th class="px-4 py-3 text-left text-sm font-semibold">类别</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">书名</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">最新章节</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">作者</th>
                        <th class="px-4 py-3 text-left text-sm font-semibold">更新时间</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr
                        v-for="(item, index) in latestBooks"
                        :key="index"
                        class=" hover:bg-gray-50 transition-colors cursor-pointer"
                        @click="navigateToBook(item.id)"
                        >
                        <td class="px-4 py-3 text-sm text-gray-600">{{ item.category }}</td>
                        <td class="px-4 py-3 text-sm font-medium">{{ item.title }}</td>
                        <td class="px-4 py-3 text-sm text-gray-600">{{ item.latestChapter }}</td>
                        <td class="px-4 py-3 text-sm text-gray-600">{{ item.author }}</td>
                        <td class="px-4 py-3 text-sm text-gray-500">{{ item.updateTime }}</td>
                        </tr>
                    </tbody>
                    </table>
                </div>
                </div>
            </div>

            <!-- 更新榜单  -->
            <div class="w-full md:w-1/3">
                <RankList title="更新榜单" :books="updateBooks" :show-more="true" more-link="/rank#updatest"/>
            </div>
            </div>
        </section>
    </div>

</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import Header from '@/components/Header.vue'
import RankList from '@/components/RankList.vue'
import { useRouter } from 'vue-router'
import { bookAPI, rankAPI, rcmdAPI } from '@/api/services'

const router = useRouter()

// 当前显示第几本
const currentIndex = ref(0)

// 精品推荐轮播当前索引
const featuredCurrentIndex = ref(0)

// 导航到书籍详情页
const navigateToBook = (bookId) => {
    router.push(`/book/${bookId}`)
}


// 将模拟数据替换为响应式变量
const mainBooks = ref([])
const weeklyBooks = ref([])
const hotBooks = ref([])
const featuredBooks = ref([])
const clickBooks = ref([])
const newBooks = ref([])
const latestBooks = ref([])
const updateBooks = ref([])

// 计算属性：当前书籍
const currentBook = computed(() => mainBooks.value[currentIndex.value])

// 切换方法
function setCurrentIndex(index) {
  if (index >= 0 && index < mainBooks.value.length) {
    console.log(index)
    currentIndex.value = index
  }
}

onMounted(async () => {
  try {
    // 并行获取所有需要的数据
    const [
      curatedBooksResponse,
      topBooksResponse,
      hotBooksResponse,
      featuredBooksResponse,
      clickRankResponse,
      newBookRankResponse,
      updateRankResponse,
      allBooksResponse
    ] = await Promise.all([
      rcmdAPI.getCuratedBooks(),
      rcmdAPI.getTopBooks(),
      rcmdAPI.getHotBooks(),
      rcmdAPI.getFeaturedBooks(),
      rankAPI.getClickRank(),
      rankAPI.getNewBookRank(),
      rankAPI.getUpdateRank(),
      bookAPI.getAll(1, 20)
    ])
    
    // 提取并处理各响应数据
    mainBooks.value = (curatedBooksResponse.data || []).map(item => item.book).slice(0, 6) || []
    weeklyBooks.value = (topBooksResponse.data || []).map(item => item.book).slice(0, 6) || []
    hotBooks.value = (hotBooksResponse.data || []).map(item => item.book).slice(0, 6) || []
    featuredBooks.value = (featuredBooksResponse.data || []).map(item => item.book).slice(0, 6) || []
    clickBooks.value = (clickRankResponse.data || []).map(item => item.book).slice(0, 6) || []
    newBooks.value = (newBookRankResponse.data || []).map(item => item.book).slice(0, 6) || []
    updateBooks.value = (updateRankResponse.data || []).map(item => item.book).slice(0, 6) || []
    latestBooks.value = (allBooksResponse.data.books || []).slice(0, 20)
    
  } catch (error) {
    console.error('Failed to fetch data:', error)
  }
  
  // 精品推荐自动轮播
  if (featuredBooks.value.length > 0) {
    const featuredTimer = setInterval(() => {
      featuredCurrentIndex.value = (featuredCurrentIndex.value + 1) % featuredBooks.value.length
    }, 3000)
  }
  
  console.log('Component mounted.')
})
</script>

<style scoped>
.card-hover:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease;
  cursor: pointer;
}

.text-shadow {
  text-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
}
</style>