import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'

// 从localStorage加载用户信息
const loadUserFromStorage = () => {
  const userStr = localStorage.getItem('user')
  if (userStr) {
    try {
      return JSON.parse(userStr)
    } catch (error) {
      console.error('解析用户信息失败:', error)
      localStorage.removeItem('user')
      return null
    }
  }
  return null
}

// 用户状态管理
export const useUserStore = () => {
  const router = useRouter()
  const user = ref(loadUserFromStorage())

  // 设置用户信息
  const setUser = (userData) => {
    user.value = userData
    if (userData) {
      localStorage.setItem('user', JSON.stringify(userData))
    } else {
      localStorage.removeItem('user')
    }
  }

  // 退出登录
  const logout = () => {
    setUser(null)
    router.push('/')
  }

  // 检查用户是否登录
  const isLoggedIn = computed(() => !!user.value)

  // 检查用户是否为作者
  const isAuthor = computed(() => user.value?.type === 'author')

  // 检查用户是否为管理员
  const isAdmin = computed(() => user.value?.type === 'admin')

  // 获取用户ID
  const userId = computed(() => user.value?.id)

  return {
    user,
    setUser,
    logout,
    isLoggedIn,
    isAuthor,
    isAdmin,
    userId
  }
}