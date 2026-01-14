import { createRouter, createWebHistory } from "vue-router";

import HomeView from "../views/Home.vue";
import ProductsView from "../views/Products.vue";
import ProductsTableView from "../views/ProductsTable.vue";
import ProductsEditorView from "../views/ProductsEditor.vue";
import AboutView from "../views/About.vue";
import LoginView from "../views/Login.vue";
import RegisterView from "../views/Register.vue";
import CartView from "../views/Cart.vue";
import DirectorView from "../views/EmployeeDirector.vue";
import ManagerView from "../views/EmployeeManager.vue";
import SupportView from "../views/EmployeeSupport.vue";

import { authState, fetchMe } from "../stores/auth.js";

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "home",
      component: HomeView,
    },
    {
      path: "/products",
      name: "products",
      component: ProductsView,
    },
    {
      path: "/products-table",
      name: "products-table",
      component: ProductsTableView,
    },
    {
      path: "/products-editor",
      name: "products-editor",
      component: ProductsEditorView,
      meta: { requiresAuth: true, requiresRole: "admin" },
    },
    {
      path: "/cart",
      name: "cart",
      component: CartView,
      meta: { requiresAuth: true },
    },
    {
      path: "/login",
      name: "login",
      component: LoginView,
    },
    {
      path: "/register",
      name: "register",
      component: RegisterView,
    },
    {
      path: "/about",
      name: "about",
      component: AboutView,
    },
    {
      path: "/employees/director",
      name: "employee-director",
      component: DirectorView,
    },
    {
      path: "/employees/manager",
      name: "employee-manager",
      component: ManagerView,
    },
    {
      path: "/employees/support",
      name: "employee-support",
      component: SupportView,
    },
  ],
});

router.beforeEach(async (to) => {
  const requiresAuth = to.meta?.requiresAuth;
  const requiresRole = to.meta?.requiresRole;

  if (requiresAuth || requiresRole) {
    if (!authState.user) {
      await fetchMe();
    }
    if (!authState.user) {
      return { name: "login" };
    }
    if (requiresRole && authState.user.role !== requiresRole) {
      return { name: "products" };
    }
  }
});

export default router;
