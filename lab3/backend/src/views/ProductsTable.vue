<template>
  <main class="container mt-4">
    <h1>Список товаров</h1>

    <div class="mb-3">
      <RouterLink class="btn btn-outline-secondary" to="/products">
        <i class="bi bi-card-list me-1"></i>
        Показать карточками
      </RouterLink>
    </div>

    <p class="text-muted">Табличное представление ассортимента товаров магазина автозапчастей.</p>

    <p v-if="isLoading" class="text-muted">Загрузка списка товаров...</p>
    <p v-else-if="loadError" class="text-danger">{{ loadError }}</p>

    <div v-else class="table-responsive">
      <table class="table table-bordered table-striped table-hover align-middle">
        <thead class="table-light">
          <tr>
            <th>#</th>
            <th>Наименование</th>
            <th>Описание</th>
            <th>Цена</th>
            <th>В корзину</th>
          </tr>
        </thead>

        <tbody>
          <tr v-for="(product, index) in products" :key="product.id">
            <td>{{ index + 1 }}</td>
            <td>{{ product.name }}</td>
            <td>{{ product.description }}</td>
            <td>{{ formatPrice(product.price) }} руб.</td>
            <td>
              <button class="btn btn-warning btn-sm" @click="addToCart(product.id)">
                <i class="bi bi-cart"></i>
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <nav v-if="totalPages > 1" class="mt-4">
      <ul class="pagination justify-content-center">
        <li class="page-item" :class="{ disabled: currentPage === 1 }">
          <button class="page-link" :disabled="currentPage === 1" @click="changePage(currentPage - 1)">
            Назад
          </button>
        </li>
        <li v-for="page in totalPages" :key="page" class="page-item" :class="{ active: page === currentPage }">
          <button class="page-link" @click="changePage(page)">{{ page }}</button>
        </li>
        <li class="page-item" :class="{ disabled: currentPage === totalPages }">
          <button
            class="page-link"
            :disabled="currentPage === totalPages"
            @click="changePage(currentPage + 1)"
          >
            Далее
          </button>
        </li>
      </ul>
    </nav>
  </main>
</template>

<script setup>
import { onMounted, ref } from "vue";

const products = ref([]);
const isLoading = ref(true);
const loadError = ref("");
const currentPage = ref(1);
const totalPages = ref(1);
const pageSize = 3;

const addToCart = async (productId) => {
  try {
    const response = await fetch("/api/cart", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ productId, quantity: 1 }),
    });
    if (response.status === 401) {
      window.alert("Нужна авторизация.");
      return;
    }
    if (!response.ok) {
      throw new Error("Не удалось добавить товар.");
    }
    window.alert("Товар добавлен в корзину");
  } catch (error) {
    window.alert(error instanceof Error ? error.message : "Ошибка добавления товара.");
  }
};

const formatPrice = (price) => new Intl.NumberFormat("ru-RU").format(price);

const loadProducts = async (page = 1) => {
  isLoading.value = true;
  try {
    const response = await fetch(`/api/products?page=${page}&limit=${pageSize}`);
    if (!response.ok) {
      throw new Error("Не удалось загрузить список товаров.");
    }
    const data = await response.json();
    products.value = data.items ?? [];
    currentPage.value = data.page ?? page;
    totalPages.value = data.totalPages ?? 1;
  } catch (error) {
    loadError.value = error instanceof Error ? error.message : "Ошибка загрузки товаров.";
  } finally {
    isLoading.value = false;
  }
};

onMounted(() => {
  loadProducts();
});

const changePage = (page) => {
  if (page < 1 || page > totalPages.value || page === currentPage.value) {
    return;
  }
  loadProducts(page);
};
</script>
