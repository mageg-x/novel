<!-- WordHeader.vue -->
<template>
  <div class="word-header" data-v-3fda83c4 data-v-0cc19a20>
    <div class="word-header-title" data-v-3fda83c4></div>
    <div class="word-header-top" data-v-3fda83c4>
      <div class="word-header-top-l" data-v-3fda83c4></div>
      <div class="word-header-top-r" data-v-3fda83c4></div>
    </div>
    <div class="word-header-bottom" data-v-3fda83c4>
      <div class="word-header-bottom-l" data-v-3fda83c4></div>
      <div class="word-header-bottom-m" data-v-3fda83c4></div>
      <div class="word-header-bottom-r" data-v-3fda83c4>
        <div class="exit" data-v-3fda83c4 @click="onExit">退出摸鱼</div>
      </div>
    </div>
    <div class="close" data-v-3fda83c4 @click="onClose"></div>
  </div>
</template>

<script setup>
// 如果需要传递事件给父组件，可以使用 defineEmits
const emit = defineEmits(['exit', 'close'])

function onExit() {
  emit('exit')
}

function onClose() {
  emit('close')
}
</script>

<style scoped>
/* 保留原始样式（可选，如果全局已引入则无需重复） */
.word-header[data-v-3fda83c4] {
  position: fixed;
  top: 0;
  width: 100%;
  z-index: 9999;
}
.word-header-title[data-v-3fda83c4] {
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/top2.95b18e3.png) no-repeat;
  background-size: contain;
  height: 32px;
  left: 50%;
  position: fixed;
  top: 0;
  transform: translateX(-50%);
  width: 168px;
}
.word-header-top[data-v-3fda83c4] {
  align-items: start;
  background-color: #2a579a;
  display: flex;
  height: 71px;
  justify-content: space-between;
}
.word-header-top-l[data-v-3fda83c4] {
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/top1.ea88a28.jpg) no-repeat;
  background-size: auto 100%;
  height: 100%;
  width: 738px;
}
.word-header-top-r[data-v-3fda83c4] {
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/top3.41cd289.jpg) no-repeat;
  background-size: auto 100%;
  height: 100%;
  width: 213px;
}
.word-header-bottom[data-v-3fda83c4] {
  align-items: center;
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/bottom3.92c1038.png);
  display: flex;
  height: 114px;
  justify-content: space-between;
}
.word-header-bottom-l[data-v-3fda83c4] {
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/bottom1.c397253.jpg) no-repeat;
  background-size: auto 100%;
  height: 100%;
  width: 801px;
}
.word-header-bottom-m[data-v-3fda83c4] {
  flex: 1;
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/bottom2.fdaa13b.jpg) no-repeat;
  background-size: auto 100%;
  height: 100%;
}
.word-header-bottom-r[data-v-3fda83c4] {
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/bottom4.2f5a94a.jpg) no-repeat;
  background-size: auto 100%;
  height: 100%;
  position: relative;
  width: 287px;
}
.word-header .exit[data-v-3fda83c4] {
  background: url(//qqreader-site-1252317822.file.myqcloud.com/qq-pc/static/img/exit.ba0d295.png) no-repeat;
  background-size: 29px 29px;
  color: #444;
  cursor: pointer;
  font-size: 10px;
  left: 140px;
  padding-top: 32px;
  position: absolute;
  top: 64px;
}
.word-header .close[data-v-3fda83c4] {
  cursor: pointer;
  height: 40px;
  position: fixed;
  right: 10px;
  top: 0;
  width: 40px;
}
</style>