<template>
  <div class="flex min-h-screen items-center justify-center bg-gray-50">
    <div class="w-full max-w-md bg-white rounded-lg shadow-md p-8">
      <div class="text-center mb-8">
        <i class="fas fa-book text-primary text-4xl mb-2"></i>
        <h2 class="text-2xl font-bold text-gray-800">注册</h2>
        <p class="text-gray-500 mt-2">创建新账号，开启阅读之旅</p>
      </div>

      <form @submit.prevent="handleRegister">
        <div class="mb-4">
          <label for="username" class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
          <input type="text" id="username" v-model="form.username" placeholder="请输入用户名"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            required />
        </div>

        <div class="mb-4">
          <label for="password" class="block text-sm font-medium text-gray-700 mb-1">密码</label>
          <input type="password" id="password" v-model="form.password" placeholder="请输入密码"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            required />
        </div>

        <div class="mb-6">
          <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">确认密码</label>
          <input type="password" id="confirmPassword" v-model="confirmPassword" placeholder="请再次输入密码"
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary/50 text-sm"
            required />
          <p v-if="form.password && confirmPassword && form.password !== confirmPassword"
            class="text-red-500 text-xs mt-1">
            两次输入的密码不一致
          </p>
        </div>

        <div class="mb-6">
          <label class="block text-sm font-medium text-gray-700 mb-2">用户身份</label>
          <div class="flex space-x-4">
            <label class="flex items-center">
              <input type="radio" v-model="form.type" value="normal" class="text-primary focus:ring-primary/50" />
              <span class="ml-2 text-gray-700">读者</span>
            </label>
            <label class="flex items-center">
              <input type="radio" v-model="form.type" value="author" class="text-primary focus:ring-primary/50" />
              <span class="ml-2 text-gray-700">作家</span>
            </label>
            <label class="flex items-center">
              <input type="radio" v-model="form.type" value="admin" class="text-primary focus:ring-primary/50" />
              <span class="ml-2 text-gray-700">管理员</span>
            </label>
          </div>
        </div>

        <button type="submit"
          class="w-full py-2 bg-primary text-white rounded-md hover:bg-primary/90 transition-colors text-sm font-medium"
          :disabled="isLoading || (form.password && confirmPassword && form.password !== confirmPassword)">
          <i v-if="isLoading" class="fas fa-spinner fa-spin mr-2"></i>
          {{ isLoading ? '注册中...' : '注册' }}
        </button>

        <div class="text-center mt-4">
          <p class="text-gray-600 text-sm">
            已有账号？
            <router-link to="/login" class="text-primary hover:text-primary/80 font-medium">立即登录</router-link>
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
import Notice from '@/components/Notice.vue'

const router = useRouter()

const form = ref({
  username: '',
  password: '',
  type: 'normal' // 默认普通用户
})

const confirmPassword = ref('')
const isLoading = ref(false)

// Notice状态
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

const handleRegister = async () => {
  if (form.value.password !== confirmPassword.value) {
    showNotice('error', '注册失败', '两次输入的密码不一致')
    return
  }

  try {
    isLoading.value = true
    await userAPI.register(form.value)
    showNotice('success', '注册成功', '请登录')

    // 延迟跳转，让用户看到成功提示
    setTimeout(() => {
      router.push('/login')
    }, 1500)
  } catch (error) {
    console.error('注册失败:', error)
    showNotice('error', '注册失败', error.response?.data?.error || '注册失败，请稍后重试')
  } finally {
    isLoading.value = false
  }
}
</script>