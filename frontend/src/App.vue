<template>
  <div id="app">
    <!-- 路由视图 -->
    <router-view />
    <div class="hidden md:block">
      <Footer />
    </div>
    
  </div>
</template>

<script setup>
import Footer from '@/components/Footer.vue'
</script>

<style scoped>

</style>
