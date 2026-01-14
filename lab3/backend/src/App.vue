<template>
  <header>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="container-fluid">
        <RouterLink class="navbar-brand d-flex align-items-center" to="/">
          <img
            src="/assets/img/logo.png"
            alt="Авто-Оргул"
            width="30"
            height="30"
            class="d-inline-block align-text-top me-2"
          />
          Запчасти ВАЗ
        </RouterLink>

        <button
          class="navbar-toggler"
          type="button"
          data-bs-toggle="collapse"
          data-bs-target="#navbarNav"
          aria-controls="navbarNav"
          aria-expanded="false"
          aria-label="Toggle navigation"
        >
          <span class="navbar-toggler-icon"></span>
        </button>

        <div class="collapse navbar-collapse" id="navbarNav">
          <ul class="navbar-nav">
            <li class="nav-item">
              <RouterLink class="nav-link" active-class="active" to="/">
                <i class="bi bi-house-door me-1"></i>
                Главная
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink class="nav-link" active-class="active" to="/products">
                <i class="bi bi-basket2-fill me-1"></i>
                Каталог
              </RouterLink>
            </li>
            <li class="nav-item" v-if="authState.user">
              <RouterLink class="nav-link" active-class="active" to="/cart">
                <i class="bi bi-cart3 me-1"></i>
                Корзина
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink class="nav-link" active-class="active" to="/about">
                <i class="bi bi-info-circle me-1"></i>
                О нас
              </RouterLink>
            </li>
            <li class="nav-item" v-if="isAdmin">
              <RouterLink class="nav-link" active-class="active" to="/products-editor">
                <i class="bi bi-pencil-square me-1"></i>
                Редактор
              </RouterLink>
            </li>
          </ul>

          <div class="ms-auto d-flex align-items-center gap-2">
            <template v-if="authState.user">
              <span class="text-muted small">Привет, {{ authState.user.name }}</span>
              <button class="btn btn-outline-secondary btn-sm" type="button" @click="handleLogout">
                Выйти
              </button>
            </template>
            <template v-else>
              <RouterLink class="btn btn-outline-primary btn-sm" to="/login">Вход</RouterLink>
              <RouterLink class="btn btn-primary btn-sm" to="/register">Регистрация</RouterLink>
            </template>
          </div>
        </div>
      </div>
    </nav>
  </header>

  <RouterView />

  <footer></footer>
</template>

<script setup>
import { computed, onMounted } from "vue";
import { RouterLink, RouterView } from "vue-router";

import { authState, fetchMe, logout } from "./stores/auth.js";

const isAdmin = computed(() => authState.user?.role === "admin");

const handleLogout = async () => {
  await logout();
};

onMounted(() => {
  fetchMe();
});
</script>
