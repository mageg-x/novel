import api from "./axios";
import { useUserStore } from "@/stores/user";

// 分类API
export const categoryAPI = {
  getBooks: (categoryId, page = 1, limit = 20) =>api.get(`/categories/${categoryId}/books?page=${page}&limit=${limit}`),
};

// 书籍API
 export const bookAPI = {
   getAll: (page = 1, limit = 20) =>api.get(`/books?page=${page}&limit=${limit}`),
   getById: (id) => api.get(`/books/${id}`),
   getChapters: (bookId) => api.get(`/books/${bookId}/chapters`),
   getChapter: (bookId, chapterId) =>api.get(`/books/${bookId}/chapters/${chapterId}`),
   getRelatedBooks: (bookId) => api.get(`/books/${bookId}/related`),
   getComments: (bookId) => api.get(`/books/${bookId}/comments`),
   search: (query, page = 1, limit = 20) => api.get(`/books/search?keyword=${query}&page=${page}&limit=${limit}`),
   addBook: (data) => api.post("/books", data),
   updateBook: (id, data) => api.put(`/books/${id}`, data),
   deleteBook: (id) => api.delete(`/books/${id}`),
 };

// 章节API
export const chapterAPI = {
  addChapter: (bookId, chapterId, data) =>api.post(`/books/${bookId}/chapters/${chapterId}`, data),
  updateChapter: (bookId, chapterId, data) =>api.put(`/books/${bookId}/chapters/${chapterId}`, data),
  deleteChapter: (bookId, chapterId) =>api.delete(`/books/${bookId}/chapters/${chapterId}`),
};

// 排行榜API
export const rankAPI = {
  getClickRank: () => api.get("/ranks/click"),
  getNewBookRank: () => api.get("/ranks/new"),
  getUpdateRank: () => api.get("/ranks/update"),
  getCommentRank: () => api.get("/ranks/comment"),
  // 管理接口 - 需要管理员权限
  addRank: (type, data) => api.post(`/ranks/${type}`, data),
  deleteRank: (type, id) => api.delete(`/ranks/${type}/${id}`),
  updateRanks: (type, data) => api.put(`/ranks/${type}`, data),
};

// 推荐API
export const rcmdAPI = {
  getHotBooks: () => api.get("/rcmd/hot"),
  getTopBooks: () => api.get("/rcmd/top"),
  getCuratedBooks: () => api.get("/rcmd/curated"),
  getFeaturedBooks: () => api.get("/rcmd/featured"),
  // 管理接口 - 需要管理员权限
  addRcmd: (type, data) => api.post(`/rcmd/${type}`, data),
  deleteRcmd: (type, id) => api.delete(`/rcmd/${type}/${id}`),
  updateRcmds: (type, data) => api.put(`/rcmd/${type}`, data),
};

// 用户API
 export const userAPI = {
   login: (data) => api.post("/users/login", data),
   register: (data) => api.post("/users/register", data),
   getByID: async (userId) => {
     const userStore = useUserStore();
     const response = await api.get(`/users/${userId}`);
     if (response.data && response.data.id === userStore.userId.value) {
       userStore.setUser(response.data);
     }
     return response;
   },
   getByName: async (username) => {
     const userStore = useUserStore();
     const response = await api.get(`/users/info/${username}`);
     if (response.data && response.data.id === userStore.userId.value) {
       userStore.setUser(response.data);
     }
     return response;
   },
   search: (query, page = 1, limit = 20) => api.get(`/users/search?keyword=${query}&page=${page}&limit=${limit}`),
   updateUser: async (userId, data) => {
     const userStore = useUserStore();
     const response = await api.put(`/users/${userId}`, data);
     if (response.data && response.data.id === userStore.userId.value) {
       userStore.setUser(response.data);
     }
     return response;
   },
 };

// 作者API
export const authorAPI = {
  getBooks: (authorId) => api.get(`/author/books/${authorId}`),
  getInfo: (authorId) => api.get(`/author/info/${authorId}`),
  getStats: (authorId) => api.get(`/author/stats/${authorId}`),
  // 作者操作相关接口 - 需要登录且为作者身份
  addBook: (data) => api.post("/author/books", data),
  updateBook: (bookId, data) => api.put(`/author/books/${bookId}`, data),
  deleteBook: (bookId) => api.delete(`/author/books/${bookId}`),
  addChapter: (bookId, data) =>api.post(`/author/books/${bookId}/chapters`, data),
  updateChapter: (bookId, chapterId, data) =>api.put(`/author/books/${bookId}/chapters/${chapterId}`, data),
  deleteChapter: (bookId, chapterId) =>api.delete(`/author/books/${bookId}/chapters/${chapterId}`),
};

// 书架API
export const shelfAPI = {
  getShelf: (userId) => api.get(`/users/${userId}/shelf`),
  addToShelf: (userId, data) => api.post(`/users/${userId}/shelf`, data),
  removeFromShelf: (userId, bookId) =>api.delete(`/users/${userId}/shelf/${bookId}`),
};

// 阅读历史API
export const historyAPI = {
  getHistory: (userId) => api.get(`/users/${userId}/history`),
  updateReadingProgress: (userId, bookId, data) =>api.post(`/users/${userId}/books/${bookId}/progress`, data),
};

// 管理员API - 只包含评论管理接口
 export const adminAPI = {
   // 评论管理 - 需要管理员权限
   getComments: (bookId = "", search = "", page = 1, limit = 20) =>
     api.get(
       `/comments?bookId=${bookId}&search=${search}&page=${page}&limit=${limit}`
     ),
   searchComments: (query, page = 1, limit = 20) => api.get(`/comments/search?keyword=${query}&page=${page}&limit=${limit}`),
   deleteComment: (id) => api.delete(`/comments/${id}`),
 };
