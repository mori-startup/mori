import { createRouter, createWebHistory } from "vue-router";
import Auth from "../components/Auth.vue";
import store from "@/store";

const routes = [
  {
    path: "/",
    name: "auth",
    component: Auth,
  },
  {
    path: "/sign-in",
    name: "sign-in",
    component: () => import("../views/SignInView.vue"),
  },
  {
    path: "/register",
    name: "register",
    component: () => import("../views/RegisterView.vue"),
  },
  {
    path: "/main",
    name: "mainpage",
    components: {
      default: () => import("../views/MainView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue"),
    },
    meta: { requiresAuth: true },
  },
  {
    path: "/messages",
    name: "messages",
    component: () => import("../views/ChatBoxView.vue"),
    props: (route) => ({
      name: route.query.name || "Conversation",
      receiverId: route.query.receiverId || null,
      type: route.query.type || "PERSON",
    }),
    meta: { requiresAuth: true },
  },
  {
    path: "/profile/:id",
    name: "Profile",
    components: {
      default: () => import("../views/ProfileView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue"),
    },
    meta: { requiresAuth: true },
  },
  {
    path: "/group/:id",
    name: "Group",
    components: {
      default: () => import("../views/GroupView.vue"),
      Chat: () => import("@/components/Chat/Chat.vue"),
    },
    meta: { requiresAuth: true },
  },
  {
    path: "/verified",
    name: "Verified",
    component: () => import("@/views/Verified.vue"),
  },
  {
    path: "/forgotpassword",
    name: "ForgotPassword",
    component: () => import("@/views/ForgotPassword.vue"),
  },
  {
    path: "/reset-password",
    name: "ResetPassword",
    component: () => import("@/views/ResetPassword.vue"),
  },
  // Update the Upgrade IA route to load the new DropFiles view
  {
    path: "/upgrade-ia",
    name: "UpgradeIA",
    component: () => import("../views/DropFiles.vue"),
    meta: { requiresAuth: true },
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

router.beforeEach(async (to, from, next) => {
  const isAuthenticated = await store.dispatch("isLoggedIn");

  // Allow unauthenticated access for specific routes
  if (
    !isAuthenticated &&
    to.name !== "sign-in" &&
    to.name !== "register" &&
    to.name !== "Verified" &&
    to.name !== "ForgotPassword" &&
    to.name !== "ResetPassword"
  ) {
    return next({ name: "sign-in" });
  }

  // Ensure WebSocket connection is established
  if (isAuthenticated && !store.state.wsConn) {
    await store.dispatch("createWebSocketConn");
  }

  // If route requires auth and user is not authenticated, redirect to sign-in
  if (to.meta.requiresAuth && !isAuthenticated) {
    return next({ name: "sign-in" });
  }

  next();
});

export default router;
