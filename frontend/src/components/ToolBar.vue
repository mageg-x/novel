<template>
    <div :class="[
        'fixed bottom-0 left-0 right-0 bg-[#469b75] text-white shadow-md transition-transform duration-300',
        showControls ? 'translate-y-0' : 'translate-y-full'
    ]">
        <div class="flex justify-around items-center h-14">
            <a href="/" class="text-center flex-1 py-2 transition-all duration-200"
                :class="{ 'opacity-100': activeTab === '', 'opacity-60': activeTab !== '' }">
                <div class="text-xl mb-1" :class="{ 'opacity-100': activeTab === '', 'opacity-60': activeTab !== '' }">
                    <i class="fas fa-home"></i>
                </div>
                <div class="text-xs" :class="{ 'opacity-100': activeTab === '', 'opacity-60': activeTab !== '' }">
                    首页
                </div>
            </a>
            <a href="/history?tab=recent" class="text-center flex-1 py-2 px-1 transition-all duration-200"
                :class="{ 'bg-gradient-to-t from-white to-emerald-50 rounded-lg shadow-sm': activeTab === 'recent' }">
                <div class="text-xl mb-1"
                    :class="{ 'text-emerald-600': activeTab === 'recent', 'opacity-60': activeTab !== 'recent' }">
                    <i class="fas fa-book"></i>
                </div>
                <div class="text-xs"
                    :class="{ 'text-emerald-700 font-medium': activeTab === 'recent', 'opacity-60': activeTab !== 'recent' }">
                    阅读记录
                </div>
            </a>
            <a href="/history?tab=shelf" class="text-center flex-1 py-2 px-1 transition-all duration-200"
                :class="{ 'bg-gradient-to-t from-white to-emerald-50 rounded-lg shadow-sm': activeTab === 'shelf' }">
                <div class="text-xl mb-1"
                    :class="{ 'text-emerald-600': activeTab === 'shelf', 'opacity-60': activeTab !== 'shelf' }">
                    <i class="fa-solid fa-swatchbook"></i>
                </div>
                <div class="text-xs"
                    :class="{ 'text-emerald-700 font-medium': activeTab === 'shelf', 'opacity-60': activeTab !== 'shelf' }">
                    书架
                </div>
            </a>
            <a href="/user" class="text-center flex-1 py-2 transition-all duration-200"
                :class="{ 'opacity-100': activeTab === 'user', 'opacity-60': activeTab !== 'user' }">
                <div class="text-xl mb-1"
                    :class="{ 'opacity-100': activeTab === 'user', 'opacity-60': activeTab !== 'user' }">
                    <i class="fas fa-user"></i>
                </div>
                <div class="text-xs" :class="{ 'opacity-100': activeTab === 'user', 'opacity-60': activeTab !== 'user' }">
                    我的
                </div>
            </a>
        </div>
    </div>
</template>

<script setup>
// 定义组件的属性
const props = defineProps({
    showControls: {
        type: Boolean,
        default: true
    },
    activeTab: {
        type: String,
        default: ''
    }
})

// 定义组件的事件
const emit = defineEmits(['tab-change'])
</script>

<style scoped>
/* 可以添加组件特定的样式 */
</style>