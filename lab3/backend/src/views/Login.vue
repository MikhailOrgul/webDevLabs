<template>
  <main class="container mt-4" style="max-width: 520px">
    <h1>Вход</h1>
    <p class="text-muted">Введите почту и пароль, чтобы войти.</p>

    <form class="card shadow-sm p-4" @submit.prevent="handleSubmit">
      <div class="mb-3">
        <label class="form-label" for="login-email">Email</label>
        <input
          id="login-email"
          v-model.trim="email"
          type="email"
          class="form-control"
          autocomplete="username"
          required
        />
      </div>
      <div class="mb-3">
        <label class="form-label" for="login-password">Пароль</label>
        <input
          id="login-password"
          v-model="password"
          type="password"
          class="form-control"
          autocomplete="current-password"
          required
        />
      </div>

      <p v-if="errorMessage" class="text-danger mb-3">{{ errorMessage }}</p>

      <button class="btn btn-primary w-100" type="submit" :disabled="isSubmitting">
        {{ isSubmitting ? "Входим..." : "Войти" }}
      </button>
    </form>
  </main>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

import { login } from "../stores/auth.js";

const router = useRouter();
const email = ref("");
const password = ref("");
const errorMessage = ref("");
const isSubmitting = ref(false);

const handleSubmit = async () => {
  errorMessage.value = "";
  isSubmitting.value = true;
  try {
    await login(email.value, password.value);
    router.push({ name: "products" });
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Не удалось войти.";
  } finally {
    isSubmitting.value = false;
  }
};
</script>
