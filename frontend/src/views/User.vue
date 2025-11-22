<template>
  <div class="hidden md:block">
    <Header />
  </div>
  <!-- 小屏幕 -->
  <div class="block md:hidden">
    <!-- 顶部导航 -->
    <div class="sticky top-0 z-50 bg-gradient-to-r from-emerald-600 to-teal-600 text-white shadow-lg">
      <div class="flex items-center justify-between h-14 px-4">
        <button @click="goBack"
          class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
          <i class="fas fa-arrow-left"></i>
        </button>
        <span class="text-lg font-semibold">个人中心</span>
        <button @click="goHome"
          class="text-white text-xl transition-transform duration-200 hover:scale-110 cursor-pointer">
          <i class="fa-solid fa-house"></i>
        </button>
      </div>
    </div>
  </div>

  <!-- 背景渐变 -->
  <div class="min-h-screen bg-gradient-to-br from-blue-50 via-purple-50 to-pink-50 py-12 relative">
    <!-- 消息提示组件 -->
    <div v-if="message" :class="[
      'fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg transform transition-all duration-500 ease-in-out',
      messageType === 'success' ? 'bg-green-500 text-white' :
        messageType === 'error' ? 'bg-red-500 text-white' :
          'bg-blue-500 text-white'
    ]">
      <i :class="[
        'mr-2',
        messageType === 'success' ? 'fas fa-check-circle' :
          messageType === 'error' ? 'fas fa-exclamation-circle' :
            'fas fa-info-circle'
      ]"></i>
      {{ message }}
    </div>

    <!-- 加载指示器 -->
    <div v-if="isLoading" class="fixed inset-0 bg-black/20 backdrop-blur-sm flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-xl shadow-2xl flex flex-col items-center">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-primary mb-4"></div>
        <p class="text-gray-700">加载中...</p>
      </div>
    </div>

    <div class="container mx-auto px-4">
      <div class="max-w-4xl mx-auto">
        <!-- 用户信息卡片 -->
        <div class="bg-white rounded-2xl shadow-xl overflow-hidden transition-all duration-300 hover:shadow-2xl">
          <!-- 卡片头部背景 -->
          <div class="bg-gradient-to-r from-primary to-purple-600 h-24"></div>

          <!-- 用户头像和基本信息 -->
          <div class="text-center relative -mt-12 mb-8">
            <div class="inline-block relative">
              <img v-if="previewAvatar || userInfo.avatar" :src="previewAvatar || userInfo.avatar" alt="用户头像"
                class="w-28 h-28 rounded-full object-cover border-4 border-white shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-105">
              <div v-else
                class="w-28 h-28 rounded-full bg-gradient-to-br from-primary to-purple-500 flex items-center justify-center text-white text-5xl shadow-lg hover:shadow-xl transition-all duration-300 transform hover:scale-105">
                <i class="fas fa-user-circle"></i>
              </div>
              <!-- 编辑头像提示 -->
              <button v-if="isEditing"
                class="absolute bottom-0 right-0 bg-primary text-white rounded-full p-2 shadow-md hover:bg-primary/90 transition-colors cursor-pointer"
                @click="generateRandomAvatar">
                <i class="fas fa-random"></i>
              </button>
            </div>
            <h1 class="text-3xl font-bold text-gray-800 mt-4">{{ userInfo.nickname || userInfo.username }}</h1>
            <p class="text-gray-600 mt-1 flex items-center justify-center">
              <i class="fas fa-user-tag mr-2"></i>
              {{ userInfo.type === 'author' ? '作家' : userInfo.type === 'admin' ? '管理员' : '读者' }}
            </p>
          </div>

          <div class="px-8 pb-8 space-y-8">
            <!-- 用户基本信息 -->
            <div>
              <div class="flex justify-between items-center mb-6">
                <h2 class="text-xl font-bold text-gray-800 flex items-center">
                  <i class="fas fa-user-circle text-primary mr-2"></i>
                  基本信息
                </h2>
                <div class="flex space-x-3">
                  <!-- 编辑/保存按钮 -->
                  <button @click="handleEditToggle"
                    class="px-5 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition-all duration-300 shadow-md hover:shadow-lg transform hover:-translate-y-1">
                    <i :class="isEditing ? 'fas fa-save' : 'fas fa-edit'" class="mr-2"></i>
                    {{ isEditing ? '保存' : '编辑' }}
                  </button>

                  <!-- 取消修改按钮（仅在编辑模式下显示） -->
                  <button v-if="isEditing" @click="handleCancelEdit"
                    class="px-5 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-all duration-300 shadow-md hover:shadow-lg transform hover:-translate-y-1">
                    <i class="fas fa-times mr-2"></i>
                    取消
                  </button>
                </div>
              </div>

              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- 用户名 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">用户名</label>
                  <p class="text-gray-900 font-medium">{{ userInfo.username }}</p>
                </div>

                <!-- 用户身份 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">用户身份</label>
                  <span
                    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-primary/10 text-primary">
                    {{ userInfo.type === 'author' ? '作家' : userInfo.type === 'admin' ? '管理员' : '读者' }}
                  </span>
                </div>

                <!-- 昵称 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">昵称</label>
                  <input v-if="isEditing" v-model="editForm.nickname" type="text"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all duration-200 bg-white shadow-sm">
                  <p v-else class="text-gray-900 font-medium">{{ userInfo.nickname || '未设置' }}</p>
                </div>

                <!-- 邮箱 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">邮箱</label>
                  <input v-if="isEditing" v-model="editForm.email" type="email"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all duration-200 bg-white shadow-sm">
                  <p v-else class="text-gray-900 font-medium">{{ userInfo.email || '未设置' }}</p>
                </div>

                <!-- 性别 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">性别</label>
                  <select v-if="isEditing" v-model="editForm.sex"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all duration-200 bg-white shadow-sm">
                    <option value="0">保密</option>
                    <option value="1">男</option>
                    <option value="2">女</option>
                  </select>
                  <p v-else class="text-gray-900 font-medium">
                    {{ userInfo.sex === 1 ? '男' : userInfo.sex === 2 ? '女' : '保密' }}
                  </p>
                </div>

                <!-- 所在地 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">所在地</label>
                  <input v-if="isEditing" v-model="editForm.location" type="text"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all duration-200 bg-white shadow-sm">
                  <p v-else class="text-gray-900 font-medium">{{ userInfo.location || '未设置' }}</p>
                </div>

                <!-- 个性签名 -->
                <div class="md:col-span-2 bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">个性签名</label>
                  <textarea v-if="isEditing" v-model="editForm.desc" rows="4"
                    class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all duration-200 bg-white shadow-sm resize-none"></textarea>
                  <p v-else class="text-gray-900 font-medium leading-relaxed">{{ userInfo.desc || '未设置' }}</p>
                </div>

                <!-- 用户等级 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">用户等级</label>
                  <div class="flex items-center">
                    <i class="fas fa-crown text-yellow-500 mr-2 text-xl"></i>
                    <span class="text-gray-900 font-medium text-lg">Lv.{{ userInfo.level }}</span>
                  </div>
                </div>

                <!-- VIP状态 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">VIP状态</label>
                  <span v-if="userInfo.isVip"
                    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-yellow-100 text-yellow-800">
                    <i class="fas fa-star text-yellow-500 mr-1"></i>
                    VIP用户
                  </span>
                  <span v-else
                    class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-gray-100 text-gray-800">
                    普通用户
                  </span>
                </div>

                <!-- 注册时间 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">注册时间</label>
                  <p class="text-gray-700">{{ formatDate(userInfo.createTime) }}</p>
                </div>

                <!-- 最后更新 -->
                <div class="bg-gray-50 p-4 rounded-xl hover:shadow-md transition-shadow">
                  <label class="block text-sm font-semibold text-gray-600 mb-2">最后更新</label>
                  <p class="text-gray-700">{{ formatDate(userInfo.updateTime) }}</p>
                </div>
              </div>
            </div>

            <!-- 操作按钮 -->
            <div class="flex justify-between items-center pt-4 border-t border-gray-100">
              <button v-if="userInfo.type === 'author'" @click="goToAuthorPage"
                class="px-6 py-3 bg-gradient-to-r from-primary to-purple-500 text-white rounded-lg hover:shadow-lg transition-all duration-300 transform hover:-translate-y-1">
                <i class="fas fa-pen-fancy mr-2"></i>
                进入作家专区
              </button>
              <button @click="handleLogout"
                class="px-6 py-3 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-all duration-300 shadow-md hover:shadow-lg transform hover:-translate-y-1">
                <i class="fas fa-sign-out-alt mr-2"></i>
                退出登录
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- 小屏幕 -->
  <div class="block md:hidden">
    <!-- 底部导航 -->
    <!-- <ToolBar :activeTab="activeTab" @tab-change="handleTabChange" /> -->
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { userAPI } from '@/api/services'
import Header from '@/components/Header.vue'
import ToolBar from '@/components/ToolBar.vue';

const router = useRouter()
const userStore = useUserStore()

// 用户信息和编辑状态
const userInfo = ref({})
const isEditing = ref(false)
const editForm = ref({})
// 头像预览
const previewAvatar = ref('')
// 加载状态
const isLoading = ref(false)
// 提示信息
const message = ref('')
const messageType = ref('')
// 底部导航激活状态
const activeTab = ref('user')
// 处理底部导航切换
const handleTabChange = (tab) => {
  activeTab.value = tab
}

// 生成随机头像
const generateRandomAvatar = () => {
  // 使用随机字符串作为seed
  const seed = Math.random().toString(36).substring(2, 15) + Math.random().toString(36).substring(2, 15)
  // 生成DiceBear头像URL
  const avatarUrl = `https://api.dicebear.com/9.x/avataaars/svg?seed=${seed}`
  // 更新预览头像
  previewAvatar.value = avatarUrl
  // 更新编辑表单中的头像
  editForm.value.avatar = avatarUrl
  showMessage('头像已更新', 'success')
}

// 回退和首页跳转方法
const goBack = () => {
  router.back()
}

const goHome = () => {
  router.push('/')
}

// 格式化日期
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// 显示提示信息
const showMessage = (text, type = 'info', duration = 3000) => {
  message.value = text
  messageType.value = type
  setTimeout(() => {
    message.value = ''
    messageType.value = ''
  }, duration)
}

// 加载用户信息
const loadUserInfo = async () => {
  if (!userStore.user.value?.id) {
    showMessage('用户信息不存在', 'error')
    return
  }

  isLoading.value = true
  try {
    const response = await userAPI.getByID(userStore.user.value.id)

    // 确保响应数据结构正确（由于axios响应拦截器，response已经是response.data）
    if (response && response.data) {
      userInfo.value = response.data
      // 清除头像预览
      previewAvatar.value = ''

      // 初始化编辑表单
      editForm.value = {
        nickname: userInfo.value.nickname || '',
        avatar: userInfo.value.avatar || '',
        desc: userInfo.value.desc || '',
        sex: userInfo.value.sex || 0,
        location: userInfo.value.location || '',
        email: userInfo.value.email || ''
      }
    }
  } catch (error) {
    console.error('加载用户信息失败:', error)
    showMessage('加载用户信息失败，请稍后重试', 'error')
  } finally {
    isLoading.value = false
  }
}

// 保存用户信息
const saveUserInfo = async () => {
  if (!userStore.user.value?.id) return

  isLoading.value = true
  try {
    const response = await userAPI.updateUser(userStore.user.value.id, editForm.value)
    userInfo.value = response.data // 由于axios响应拦截器，response已经是response.data
    // 更新本地存储的用户信息
    userStore.setUser(userInfo.value)
    isEditing.value = false
    showMessage('用户信息保存成功', 'success')
  } catch (error) {
    console.error('保存用户信息失败:', error)
    showMessage('保存用户信息失败，请稍后重试', 'error')
  } finally {
    isLoading.value = false
  }
}

// 进入作家专区
const goToAuthorPage = () => {
  router.push(`/author/${userStore.user.value.id}`)
}

// 退出登录
const handleLogout = () => {
  // 清除本地存储
  localStorage.removeItem('user')

  // 更新用户状态
  userStore.setUser(null)

  // 跳转到首页
  router.push('/')
}

// 监听编辑状态变化
const handleEditToggle = () => {
  if (isEditing.value) {
    // 保存修改
    saveUserInfo()
  } else {
    // 准备编辑
    editForm.value = {
      nickname: userInfo.value.nickname || '',
      avatar: userInfo.value.avatar || '',
      desc: userInfo.value.desc || '',
      sex: userInfo.value.sex || 0,
      location: userInfo.value.location || '',
      email: userInfo.value.email || ''
    }
    isEditing.value = true
  }
}

// 取消修改
const handleCancelEdit = () => {
  // 恢复原始数据
  editForm.value = {
    nickname: userInfo.value.nickname || '',
    avatar: userInfo.value.avatar || '',
    desc: userInfo.value.desc || '',
    sex: userInfo.value.sex || 0,
    location: userInfo.value.location || '',
    email: userInfo.value.email || ''
  }
  // 清除头像预览
  previewAvatar.value = ''
  // 退出编辑模式
  isEditing.value = false
}

// 在组件挂载时加载用户信息
onMounted(() => {
  // 检查用户是否已登录
  if (!userStore.isLoggedIn.value) {
    router.push('/login')
    return
  }
  loadUserInfo()
})
</script>