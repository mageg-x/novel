<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-8">
      <div class="text-center mb-8">
        <i class="fas fa-book text-primary text-4xl mb-2"></i>
        <h2 class="text-2xl font-bold text-gray-800">登录</h2>
        <p class="text-gray-500 mt-2">欢迎回来，请登录您的账号</p>
      </div>

      <form @submit.prevent="handleLogin">
        <div class="mb-4">
          <label for="username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
          <input type="text" id="username" v-model="form.username" placeholder="请输入用户名"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            required />
        </div>

        <div class="mb-6">
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">密码</label>
          <input type="password" id="password" v-model="form.password" placeholder="请输入密码"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            required />
        </div>

        <button type="submit"
          class="w-full py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm font-medium"
          :disabled="isLoading">
          <i v-if="isLoading" class="fas fa-spinner fa-spin mr-2"></i>
          {{ isLoading ? '登录中...' : '登录' }}
        </button>

        <div class="text-center mt-4">
          <p class="text-gray-600 text-sm">
            还没有账号？
            <router-link to="/register" class="text-primary hover:text-primary/80 font-medium">立即注册</router-link>
          </p>
        </div>
      </form>
    </div>
  </div>

  <!-- 通知组件 -->
  <Notice :visible="notice.visible" :type="notice.type" :title="notice.title" :message="notice.message"
    @close="closeNotice" />
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { userAPI } from '@/api/services'
import { useUserStore } from '@/stores/user'
import Notice from '@/components/Notice.vue'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: ''
})

const isLoading = ref(false)

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

const handleLogin = async () => {
  try {
    isLoading.value = true
    const response = await userAPI.login(form.value)
    const { user, token } = response.data

    // 保存用户信息和token到本地存储
    localStorage.setItem('user', JSON.stringify(user))
    localStorage.setItem('token', token)

    // 更新用户状态
    userStore.setUser(user)

    // 跳转到首页
    router.push('/')
  } catch (error) {
    console.error('登录失败:', error)
    showNotice('error', '登录失败', error.response?.data?.error || '未知错误')
  } finally {
    isLoading.value = false
  }
}
</script>