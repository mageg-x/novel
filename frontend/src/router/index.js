import { createRouter, createWebHistory } from "vue-router";
import { useUserStore } from "@/stores/user";
import Home from "@/views/Home.vue";
import Cover from "@/views/Cover.vue";
import Book from "@/views/Book.vue";
import Toc from "@/views/Toc.vue";
import Class from "@/views/Class.vue";
import Rank from "@/views/Rank.vue";
import History from "@/views/History.vue";
import Author from "@/views/Author.vue";
import Login from "@/views/Login.vue";
import Register from "@/views/Register.vue";
import User from "@/views/User.vue";

const routes = [
  { path: "/", redirect: "/home" },
  { path: "/admin", redirect: "/home" },
  {
    path: "/home",
    name: "home",
    component: Home,
  },
  {
    path: "/login",
    name: "login",
    component: Login,
  },
  {
    path: "/register",
    name: "register",
    component: Register,
  },
  {
    path: "/user",
    name: "user",
    component: User,
    meta: {
      requiresAuth: true // 需要登录才能访问
    }
  },
  {
    path: "/class",
    name: "class",
    component: Class,
  },
  {
    path: "/rank",
    name: "rank",
    component: Rank,
  },
  {
    path: "/history",
    name: "history",
    component: History,
    meta: {
      requiresAuth: true // 需要登录才能访问
    }
  },
  {
    path: "/book/:id",
    name: "cover",
    component: Cover,
  },
  {
    path: "/book/:bookId/:chapterId",
    name: "book",
    component: Book,
  },
  {
    path: "/book/:bookId/toc",
    name: "toc",
    component: Toc,
  },
  {
    path: "/author",
    name: "author",
    component: Author,
    meta: {
      requiresAuth: true, // 需要登录才能访问
      requiresAuthor: true // 需要作者身份才能访问
    }
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

// 路由守卫
router.beforeEach((to, from, next) => {
  const userStore = useUserStore();

  // 检查是否需要登录
  if (to.meta.requiresAuth) {
    if (!userStore.isLoggedIn.value) {
      // 未登录，跳转到登录页面
      next({ path: "/login" });
      return;
    }

    // 检查是否需要作者身份
    if (to.meta.requiresAuthor && !userStore.isAuthor.value) {
      // 不是作者，跳转到首页
      next({ path: "/home" });
      return;
    }
  }

  next();
});

export default router;
