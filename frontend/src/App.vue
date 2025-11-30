<template>
  <div id="app">
    <!-- 使用keep-alive缓存Home组件，避免重复加载数据 -->
    <router-view v-slot="{ Component }">
      <keep-alive include="Home">
        <component :is="Component" />
      </keep-alive>
    </router-view>
    <!-- 根据ui store中的showFooter状态来控制Footer的显示 -->
    <Footer v-if="uiStore.showFooter" />
  </div>
</template>

<script setup>
import Footer from '@/components/Footer.vue'
import { useUiStore } from '@/stores/ui.js'

// 创建ui store实例
const uiStore = useUiStore()
</script>

<style scoped></style>
