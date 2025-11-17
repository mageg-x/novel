<template>
  <transition name="notice">
    <div v-if="visible" class="fixed inset-0 flex items-center justify-center z-50 p-4">
      <!-- 半透明浅灰色遮罩 -->
      <div class="fixed inset-0 bg-black  opacity-30 "></div>
      <!-- 提示框内容 -->
      <div
        class="  bg-white rounded-lg shadow-md p-3 sm:p-4 w-full max-w-xs sm:max-w-sm mx-auto relative z-10 transform transition-all duration-300 min-w-[240px] border border-gray-100">
        <!-- 图标 -->
        <div class="flex justify-center mb-2">
          <div :class="[
            'w-8 sm:w-10 h-8 sm:h-10 rounded-full flex items-center justify-center shadow-inner',
            type === 'success' ? 'bg-green-100' : type === 'error' ? 'bg-red-100' : 'bg-blue-100'
          ]">
            <i :class="[
              'fas text-2xl sm:text-4xl',
              type === 'success' ? 'fa-check-circle text-green-500' :
                type === 'error' ? 'fa-times-circle text-red-500' :
                  'fa-info-circle text-blue-500'
            ]"></i>
          </div>
        </div>

        <!-- 标题 -->
        <h3 class="text-sm sm:text-base font-semibold text-center text-gray-800 mb-1 line-clamp-1">{{ title }}</h3>

        <!-- 消息内容 -->
        <p class="text-xs sm:text-sm text-center text-gray-600 mb-3 sm:mb-4 px-1 line-clamp-3 max-h-16 overflow-hidden">
          {{ message }}</p>

        <!-- 确认按钮 -->
        <button @click="close"
          class=" w-fit min-w-24 block m-auto py-1.5 sm:py-2 px-3 bg-[#469b75] text-white rounded-md hover:bg-[#3d8766] transition-all duration-200 font-medium text-xs sm:text-sm shadow-sm hover:shadow-md active:scale-98 focus:outline-none focus:ring-2 focus:ring-[#469b75] focus:ring-opacity-50">
          确认
        </button>
      </div>
    </div>
  </transition>
</template>

<script setup>
import { ref, onMounted } from 'vue'

// Props
const props = defineProps({
  visible: {
    type: Boolean,
    default: false
  },
  type: {
    type: String,
    default: 'info', // success, error, info
    validator: (value) => ['success', 'error', 'info'].includes(value)
  },
  title: {
    type: String,
    default: '提示'
  },
  message: {
    type: String,
    default: ''
  }
})

// Emits
const emit = defineEmits(['close'])

// 关闭提示框
const close = () => {
  emit('close')
}
</script>

<style scoped>
.notice-enter-from {
  opacity: 0;
  transform: scale(0.9);
}

.notice-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>