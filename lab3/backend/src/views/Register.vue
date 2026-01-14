<template>
  <main class="container mt-4" style="max-width: 520px">
    <h1>Регистрация</h1>
    <p class="text-muted">Создайте учетную запись для покупок.</p>

    <form class="card shadow-sm p-4" @submit.prevent="handleSubmit">
      <div class="mb-3">
        <label class="form-label" for="register-name">Имя</label>
        <input
          id="register-name"
          v-model.trim="name"
          type="text"
          class="form-control"
          autocomplete="name"
          required
        />
      </div>
      <div class="mb-3">
        <label class="form-label" for="register-email">Email</label>
        <input
          id="register-email"
          v-model.trim="email"
          type="email"
          class="form-control"
          autocomplete="email"
          required
        />
      </div>
      <div class="mb-3">
        <label class="form-label" for="register-password">Пароль</label>
        <input
          id="register-password"
          v-model="password"
          type="password"
          class="form-control"
          autocomplete="new-password"
          required
        />
      </div>

      <p v-if="errorMessage" class="text-danger mb-3">{{ errorMessage }}</p>

      <button class="btn btn-success w-100" type="submit" :disabled="isSubmitting">
        {{ isSubmitting ? "Создаем..." : "Зарегистрироваться" }}
      </button>
    </form>
  </main>
</template>

<script setup>
import { ref } from "vue";
import { useRouter } from "vue-router";

import { register } from "../stores/auth.js";

const router = useRouter();
const name = ref("");
const email = ref("");
const password = ref("");
const errorMessage = ref("");
const isSubmitting = ref(false);

const handleSubmit = async () => {
  errorMessage.value = "";
  isSubmitting.value = true;
  try {
    await register(name.value, email.value, password.value);
    router.push({ name: "products" });
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Не удалось зарегистрироваться.";
  } finally {
    isSubmitting.value = false;
  }
};
</script>
