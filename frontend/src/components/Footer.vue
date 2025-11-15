<template>
  <footer class="bg-white ">
    <div class="container mx-auto px-4 py-6">
      <div class="footer-links flex flex-wrap gap-4 justify-center mb-4 text-sm text-gray-600">
        <a
          v-for="(link, index) in footerLinks"
          :key="index"
          :href="link.url"
          class="hover:text-primary transition-colors"
        >
          {{ link.text }}
        </a>
      </div>
      <p class="copyright text-center text-xs text-gray-500">
        Copyright (C) novel.mageg.cn All rights reserved. 阅林小说网站版权所有
      </p>
    </div>
  </footer>
</template>

<script setup>
// 可维护的链接数据（便于后续国际化或配置化）
const footerLinks = [
  { text: '关闭通道', url: '#' },
  { text: '小说地图屋', url: '#' },
  { text: '网站屋屋', url: '#' },
  { text: '作家实实', url: '#' }
]
</script>

<style scoped>
/* 如果你使用的是全局 CSS，这行可省略 */
/* .container 已假设在全局定义（如 Tailwind 的 container 或自定义） */
</style>