<template>
  <div class="bg-white container md:w-5xl mx-auto">
    <Header :showSearchBar="false" />
    <!-- Search Bar -->
    <div class="bg-white py-8">
      <div class="max-w-[1200px] mx-auto px-4">
        <div class="flex justify-center">
          <div class="relative w-full max-w-2xl">
            <input
              type="text"
              v-model="searchQuery"
              @keyup.enter="performSearch"
              class="w-full h-12 pl-12 pr-32 border border-gray-300 rounded-full text-gray-700 focus:outline-none focus:border-blue-400"
              placeholder="请输入搜索关键词"
            >
            <svg class="absolute left-4 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
            <button 
              @click="performSearch"
              class="absolute right-1 top-1/2 transform -translate-y-1/2 bg-[#21b2f0] text-white px-8 py-2 rounded-full hover:bg-[#1a9dd8]"
            >
              搜索
            </button>
            <button 
              @click="clearSearch"
              class="absolute right-28 top-1/2 transform -translate-y-1/2 text-gray-400 hover:text-gray-600"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-[1200px] mx-auto px-4 py-6">
      <div class="flex gap-8">
        <!-- Left Column - Search Results -->
        <div class="flex-1">
          <!-- Loading State -->
          <div v-if="loading" class="text-center py-12">
            搜索中...
          </div>
          
          <!-- Error State -->
          <div v-else-if="error" class="text-center py-12 text-red-500">
            搜索中出错，请稍后再试
          </div>
          
          <!-- No Results -->
          <div v-else-if="books.length === 0 && searchQuery" class="text-center py-12 text-gray-500">
            未找到相关书籍
          </div>
          
          <!-- Book Result List -->
          <div v-else>
            <div 
              v-for="(book, index) in books" 
              :key="index"
              class="flex gap-4 mb-6 pb-6 border-b border-gray-100"
            >
              <div class="relative flex-shrink-0">
                <img 
                  :src="book.cover || 'https://via.placeholder.com/120x160?text=No+Cover'" 
                  :alt="book.title"
                  class="w-[120px] h-[160px] object-cover rounded"
                >
                <div 
                  v-if="book.status === '完结' || book.status === '全本'"
                  class="absolute top-0 left-0 bg-red-600 text-white text-xs px-2 py-0.5 rounded-tl"
                >
                  {{ book.status }}
                </div>
              </div>
              <div class="flex-1 flex flex-col justify-evenly">
                <h3 class="text-lg mb-2">
                  <a 
                    href="#" 
                    class="hover:underline"
                    @click.prevent="viewBook(book)"
                  >
                    <span 
                      v-for="(part, i) in highlightText(book.title)"
                      :key="i"
                      :class="{'text-[#21b2f0] font-medium': part.isMatch}"
                    >
                      {{ part.text }}
                    </span>
                  </a>
                </h3>
                <p class="text-gray-600 text-sm mb-2 line-clamp-2">
                  <template v-if="book.description">
                    <span 
                      v-for="(part, i) in highlightText(book.description)"
                      :key="i"
                      :class="{'text-[#21b2f0] font-medium': part.isMatch}"
                    >
                      {{ part.text }}
                    </span>
                  </template>
                  <template v-else>暂无描述</template>
                </p>
                <p class="text-gray-400 text-xs">{{ book.author }} · {{ book.category }} · {{ book.status }}    {{ book.wordCount }}</p>
              </div>
            </div>
          </div>

          <!-- Pagination -->
          <div v-if="totalPages > 1" class="flex items-center justify-center gap-2 mt-8 mb-12">
            <button 
              @click="changePage(1)"
              :disabled="currentPage === 1"
              class="px-4 py-2 rounded text-gray-600 hover:bg-gray-100 disabled:opacity-50"
            >
              ‹
            </button>
            
            <button 
              v-for="page in pagination.pages" 
              :key="page"
              @click="changePage(page)"
              :class="[
                'px-4 py-2 rounded',
                currentPage === page 
                  ? 'bg-[#e8f6fd] text-[#21b2f0]' 
                  : 'text-gray-600 hover:bg-gray-100'
              ]"
            >
              {{ page }}
            </button>
            
            <span class="px-2 text-gray-400" v-if="pagination.hasMore">...</span>
            
            <button v-if="pagination.hasMore" @click="changePage(pagination.totalPages)"
              class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded"
            >
              {{ pagination.totalPages }}
            </button>
            
            <button @click="nextPage" :disabled="currentPage === totalPages"
              class="px-4 py-2 text-gray-600 hover:bg-gray-100 rounded disabled:opacity-50"
            >
              ›
            </button>
            
            <div class="ml-4 flex items-center gap-2">
              <span class="text-sm text-gray-600">跳转至</span>
              <input 
                v-model="jumpPage"
                type="text" 
                class="w-12 h-8 border border-gray-300 rounded text-center text-sm"
                @keyup.enter="jumpToPage"
              >
              <span class="text-sm text-gray-600">页</span>
            </div>
          </div>
        </div>
      </div>
    </div>
</div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import Header from '@/components/Header.vue'
import { bookAPI } from '@/api/services.js'

// 响应式数据
const searchQuery = ref('')
const currentPage = ref(1)
const jumpPage = ref('')
const route = useRoute()
const router = useRouter()

// 书籍数据 - 替换模拟数据
const books = ref([])
const loading = ref(false)
const error = ref('')
const totalCount = ref(0)
const pageSize = ref(20)

// 分页数据
const pagination = reactive({
  pages: [],
  totalPages: 1,
  hasMore: false
})


// 计算总页数
const totalPages = computed(() => {
  return Math.ceil(totalCount.value / pageSize.value)
})

// 更新分页显示
const updatePagination = () => {
  const total = totalPages.value
  const current = currentPage.value
  const pages = []
  
  // 生成显示的页码
  if (total <= 8) {
    // 总页数小于等于8，显示所有页码
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
    pagination.hasMore = false
  } else {
    // 总页数大于8，根据当前页生成页码显示
    if (current <= 5) {
      // 当前页在前5页，显示1-8
      for (let i = 1; i <= 8; i++) {
        pages.push(i)
      }
      pagination.hasMore = true
    } else if (current >= total - 4) {
      // 当前页在后5页，显示total-7到total
      for (let i = total - 7; i <= total; i++) {
        pages.push(i)
      }
      pagination.hasMore = false
    } else {
      // 当前页在中间，显示current-3到current+3
      for (let i = current - 3; i <= current + 3; i++) {
        pages.push(i)
      }
      pagination.hasMore = true
    }
  }
  
  pagination.pages = pages
  pagination.totalPages = total
}

// 高亮文本
  const highlightText = (text) => {
    if (!searchQuery.value) {
      return [{ text: text, isMatch: false }]
    }
    
    const query = searchQuery.value.toLowerCase()
    const parts = []
    let lastIndex = 0
    let searchIndex = text.toLowerCase().indexOf(query)
    
    while (searchIndex !== -1) {
      // 添加匹配前的文本
      if (searchIndex > lastIndex) {
        parts.push({
          text: text.substring(lastIndex, searchIndex),
          isMatch: false
        })
      }
      
      // 添加匹配的文本
      parts.push({
        text: text.substring(searchIndex, searchIndex + query.length),
        isMatch: true
      })
      
      lastIndex = searchIndex + query.length
      searchIndex = text.toLowerCase().indexOf(query, lastIndex)
    }
    
    // 添加剩余文本
    if (lastIndex < text.length) {
      parts.push({
        text: text.substring(lastIndex),
        isMatch: false
      })
    }
    
    return parts
  }

// 执行搜索
const performSearch = async () => {
  if (!searchQuery.value.trim()) {
    return
  }
  
  currentPage.value = 1
  await fetchSearchResults(searchQuery.value, 1)
  
  // 更新URL参数
  router.push({
    path: '/search',
    query: { q: searchQuery.value }
  })
}

// 获取搜索结果
const fetchSearchResults = async (query, page) => {
  if (!query) return
  
  loading.value = true
  error.value = null
  
  try {
    const response = await bookAPI.search(query, page, pageSize.value)
    
    // 数据验证和格式化
    if (!response || typeof response !== 'object') {
      throw new Error('API返回格式错误')
    }
    
    // 确保books是数组格式
    books.value = Array.isArray(response.data.books) ? response.data.books : []
    totalCount.value = Number(response.data.total) || 0
    console.log('Search results:', books.value)
  } catch (err) {
    console.error('Search error:', err)
    error.value = err.message || '搜索失败，请稍后重试'
    books.value = []
    totalCount.value = 0
  } finally {
    loading.value = false
  }
  
  // 更新分页
  updatePagination()
}

// 清空搜索
const clearSearch = () => {
  searchQuery.value = ''
  books.value = []
  totalCount.value = 0
  currentPage.value = 1
  updatePagination()
  
  // 移除URL参数
  router.push({
    path: '/search',
    query: {}
  })
}

// 切换页面
const changePage = async (page) => {
  if (page !== currentPage.value && page >= 1 && page <= totalPages.value) {
    currentPage.value = page
    await fetchSearchResults(searchQuery.value, page)
  }
}

// 下一页
const nextPage = async () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    await fetchSearchResults(searchQuery.value, currentPage.value)
  }
}

// 跳转到指定页面
const jumpToPage = async () => {
  const page = parseInt(jumpPage.value)
  if (page && page >= 1 && page <= totalPages.value) {
    await changePage(page)
    jumpPage.value = ''
  }
}

// 查看书籍详情
const viewBook = (book) => {
  // 假设书籍有id字段
  router.push(`/book/${book.id}`)
}

// 搜索热门关键词
const searchHot = (keyword) => {
  searchQuery.value = keyword
  performSearch()
}

onMounted(() => {
  // 从URL参数中获取搜索关键词
  const searchParam = route.query.q
  if (searchParam) {
    searchQuery.value = searchParam
    // 如果有搜索参数，自动执行搜索
    fetchSearchResults(searchQuery.value, 1)
  }
})
</script>

<style scoped>
/* 添加一些自定义样式 */
.book-title-highlight {
  color: #21b2f0;
  font-weight: 500;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
