import { defineStore } from 'pinia'

export const useUiStore = defineStore('ui', {
  state: () => ({
    // 控制Footer是否显示的状态，默认为true
    showFooter: true
  }),
  
  actions: {
    // 显示Footer
    displayFooter() {
      this.showFooter = true
    },
    
    // 隐藏Footer
    hideFooter() {
      this.showFooter = false
    },
    
    // 切换Footer显示状态
    toggleFooter() {
      this.showFooter = !this.showFooter
    },
    
    // 设置Footer显示状态
    setFooterVisibility(visible) {
      this.showFooter = visible
    }
  },
  
  getters: {
    // 获取Footer显示状态
    isFooterVisible: (state) => state.showFooter
  }
})