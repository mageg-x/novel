import api from './axios';

// 分类API
export const categoryAPI = {
  getBooks: (categoryId, page = 1, limit = 20) => api.get(`/categories/${categoryId}/books?page=${page}&limit=${limit}`)
};

// 书籍API
export const bookAPI = {
  getAll: (page = 1, limit = 20) => api.get(`/books?page=${page}&limit=${limit}`),
  getById: (id) => api.get(`/books/${id}`),
  getChapters: (bookId) => api.get(`/books/${bookId}/chapters`),
  getChapter: (bookId, chapterId) => api.get(`/books/${bookId}/chapters/${chapterId}`),
  getRelatedBooks: (bookId) => api.get(`/books/${bookId}/related`),
  getComments: (bookId) => api.get(`/books/${bookId}/comments`),
  addBook: (data) => api.post('/books', data),
  updateBook: (id, data) => api.put(`/books/${id}`, data),
  deleteBook: (id) => api.delete(`/books/${id}`)
};

// 章节API
export const chapterAPI = {
  addChapter: (bookId, chapterId, data) => api.post(`/books/${bookId}/chapters/${chapterId}`, data),
  updateChapter: (bookId, chapterId, data) => api.put(`/books/${bookId}/chapters/${chapterId}`, data),
  deleteChapter: (bookId, chapterId) => api.delete(`/books/${bookId}/chapters/${chapterId}`)
};

// 排行榜API
export const rankAPI = {
  getClickRank: () => api.get('/ranks/click'),
  getNewBookRank: () => api.get('/ranks/new'),
  getUpdateRank: () => api.get('/ranks/update'),
  getCommentRank: () => api.get('/ranks/comment')
};

// 推荐API
export const rcmdAPI = {
  getHotBooks: () => api.get('/rcmd/hot'),
  getTopBooks: () => api.get('/rcmd/top'),
  getCuratedBooks: () => api.get('/rcmd/curated'),
  getFeaturedBooks: () => api.get('/rcmd/featured')
};

// 用户API
export const userAPI = {
  login: (data) => api.post('/users/login', data),
  register: (data) => api.post('/users/register', data),
  getByID: (userId) => api.get(`/users/${userId}`),
  getByName: (username) => api.get(`/users/info/${username}`),
};

// 作者API
export const authorAPI = {
  getBooks: (authorId) => api.get(`/author/books/${authorId}`),
  getInfo: (authorId) => api.get(`/author/info/${authorId}`),
  getStats: (authorId) => api.get(`/author/stats/${authorId}`),
};

// 书架API
export const shelfAPI = {
  getShelf: (userId) => api.get(`/users/${userId}/shelf`),
  addToShelf: (userId, data) => api.post(`/users/${userId}/shelf`, data),
  removeFromShelf: (userId, bookId) => api.delete(`/users/${userId}/shelf/${bookId}`)
};

// 阅读历史API
export const historyAPI = {
  getHistory: (userId) => api.get(`/users/${userId}/history`),
  updateReadingProgress: (userId, bookId, data) => api.post(`/users/${userId}/books/${bookId}/progress`, data)
};
