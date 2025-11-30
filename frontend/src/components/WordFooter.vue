<!-- WordFooter.vue -->
<template>
  <div class="word-footer" data-v-53a0721a data-v-0cc19a20>
    <div class="word-footer-l" data-v-53a0721a></div>
    <div class="word-footer-m" data-v-53a0721a>
      <!-- 上一章按钮 -->
      <button 
        @click="handlePrevClick"
        :class="['chapter-btn', 'prev-btn', { 'disabled': !hasPreviousChapter }]"
        :disabled="!hasPreviousChapter"
        data-v-53a0721a
      >
        <span data-v-53a0721a>&lt; 上一章</span>
      </button>
      
      <!-- 章节间距 -->
      <span v-if="hasPreviousChapter && hasNextChapter" class="chapter-spacer" data-v-53a0721a></span>
      
      <!-- 下一章按钮 -->
      <button 
        @click="handleNextClick"
        :class="['chapter-btn', 'next-btn', { 'disabled': !hasNextChapter }]"
        :disabled="!hasNextChapter"
        data-v-53a0721a
      >
        <span data-v-53a0721a>下一章 &gt;</span>
      </button>
    </div>
    <div class="word-footer-r" data-v-53a0721a></div>
  </div>
</template>

<script setup>
// 定义组件属性
const props = defineProps({
  // 是否有上一章
  hasPreviousChapter: {
    type: Boolean,
    default: false
  },
  // 是否有下一章
  hasNextChapter: {
    type: Boolean,
    default: false
  }
})

// 定义事件
const emit = defineEmits(['onpre', 'onnext'])

// 处理上一章点击事件
function handlePrevClick() {
  // 如果上一章不可用，不执行任何操作
  if (!props.hasPreviousChapter) {
    return
  }
  
  // 发送上一章点击事件
  emit('onpre')
}

// 处理下一章点击事件
function handleNextClick() {
  // 如果下一章不可用，不执行任何操作
  if (!props.hasNextChapter) {
    return
  }
  
  // 发送下一章点击事件
  emit('onnext')
}
</script>

<style scoped>
/* 保留原始样式并添加新样式 */
.word-footer[data-v-53a0721a] {
  justify-content: space-between;
  background: #efefef url(https://fastly.jsdelivr.net/gh/mageg-x/novel/frontend/src/assets/footer3.webp);
  bottom: 0;
  display: flex;
  height: 27px;
  position: fixed;
  width: 100%;
  z-index: 9999;
}
.word-footer-l[data-v-53a0721a] {
  background: url(https://fastly.jsdelivr.net/gh/mageg-x/novel/frontend/src/assets/footer1.webp) no-repeat;
  background-size: auto 100%;
  height: 100%;
  width: 322px;
}
.word-footer-m[data-v-53a0721a] {
  color: #acacac;
  cursor: pointer;
  line-height: 27px;
  padding: 0 20px;
  display: flex;
  align-items: center;
}
.chapter-btn[data-v-53a0721a] {
  color: #acacac;
  text-decoration: none;
  padding: 0 20px;
  transition: color 0.3s;
}
.chapter-btn:hover:not(.disabled)[data-v-53a0721a] {
  color: #000;
}
.chapter-btn.disabled[data-v-53a0721a] {
  color: #d0d0d0;
  cursor: not-allowed;
  opacity: 0.5;
}
.chapter-spacer[data-v-53a0721a] {
  width: 40px;
}
.word-footer-r[data-v-53a0721a] {
  background: url(https://fastly.jsdelivr.net/gh/mageg-x/novel/frontend/src/assets/footer2.webp) no-repeat;
  height: 100%;
  width: 344px;
}
</style>