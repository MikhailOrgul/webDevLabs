<script setup>
import { RouterLink, RouterView } from "vue-router";
</script>

<template>
  <header>
    <nav class="navbar navbar-expand-lg navbar-light bg-light">
      <div class="container-fluid">
        <a class="navbar-brand" href="#">Редактор заметок</a>

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
              <RouterLink class="nav-link" to="/">
                <i class="bi bi-house"></i>
                Главная
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink class="nav-link" to="/rate">
                <i class="bi bi-star"></i>
                Оценить
              </RouterLink>
            </li>
            <li class="nav-item">
              <RouterLink class="nav-link" to="/about">
                <i class="bi bi-person-circle"></i>
                Об авторе
              </RouterLink>
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </header>

  <main class="container py-4">
    <RouterView />
  </main>

  <footer class="border-top py-3">
    <div class="container text-muted small">Лабораторная 2 — Vue SPA (заметки и todo-листы)</div>
  </footer>
</template>

<style scoped></style>
