import axios from 'axios';
import { genToken } from '@/utils/tiny';

// 创建axios实例
const api = axios.create({
  baseURL: 'http://localhost:3002/api',
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
});

// 请求拦截器
api.interceptors.request.use(
  async (config) => { // ← 加上 async
    // 添加认证信息，如token
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // 添加一次性请求token，防止爬虫
    const t = await genToken(); // ← 现在可以 await 了
    config.headers['X-Request-Token'] = t;
   
    return config;
  },
  error => {
    return Promise.reject(error);
  }
);

// 响应拦截器
api.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    console.error('API Error:', error);
    return Promise.reject(error);
  }
);

export default api;