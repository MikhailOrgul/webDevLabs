<template>
  <main class="container mt-4">
    <h1>Корзина</h1>

    <p v-if="isLoading" class="text-muted">Загрузка корзины...</p>
    <p v-else-if="errorMessage" class="text-danger">{{ errorMessage }}</p>

    <div v-else>
      <div v-if="cartItems.length === 0" class="alert alert-light">
        Корзина пуста.
      </div>

      <div v-else class="table-responsive">
        <table class="table table-bordered align-middle">
          <thead class="table-light">
            <tr>
              <th>Товар</th>
              <th>Количество</th>
              <th>Стоимость</th>
              <th></th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="item in cartItems" :key="item.product.id">
              <td>{{ item.product.name }}</td>
              <td>{{ item.quantity }}</td>
              <td>{{ formatPrice(item.lineTotal) }} руб.</td>
              <td class="text-end">
                <button class="btn btn-outline-danger btn-sm" @click="removeItem(item.product.id)">
                  Удалить
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <div v-if="cartItems.length" class="d-flex justify-content-end align-items-center gap-3">
        <div class="fs-5">Итого: <strong>{{ formatPrice(total) }} руб.</strong></div>
        <button class="btn btn-outline-secondary" @click="clearCart">Очистить</button>
      </div>
    </div>
  </main>
</template>

<script setup>
import { onMounted, ref } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const cartItems = ref([]);
const total = ref(0);
const isLoading = ref(true);
const errorMessage = ref("");

const formatPrice = (price) => new Intl.NumberFormat("ru-RU").format(price);

const loadCart = async () => {
  isLoading.value = true;
  errorMessage.value = "";
  try {
    const response = await fetch("/api/cart");
    if (response.status === 401) {
      router.push({ name: "login" });
      return;
    }
    if (!response.ok) {
      throw new Error("Не удалось загрузить корзину.");
    }
    const data = await response.json();
    cartItems.value = data.items ?? [];
    total.value = data.total ?? 0;
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Ошибка загрузки корзины.";
  } finally {
    isLoading.value = false;
  }
};

const removeItem = async (productId) => {
  await fetch(`/api/cart/${productId}`, { method: "DELETE" });
  loadCart();
};

const clearCart = async () => {
  await fetch("/api/cart/clear", { method: "POST" });
  loadCart();
};

onMounted(() => {
  loadCart();
});
</script>
