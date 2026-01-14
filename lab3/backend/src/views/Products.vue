<template>
  <main class="container mt-4">
    <h1>Товары</h1>

    <div class="mb-3">
      <RouterLink class="btn btn-outline-secondary" to="/products-table">
        <i class="bi bi-table me-1"></i>
        Показать таблицей
      </RouterLink>
    </div>

    <p v-if="isLoading" class="text-muted">Загрузка списка товаров...</p>
    <p v-else-if="loadError" class="text-danger">{{ loadError }}</p>

    <div v-else>
      <div v-for="product in products" :key="product.id" class="card shadow rounded my-3">
        <div class="card-body">
          <div class="row">
            <div class="col-md-2">
              <img :src="product.image" class="product-image" :alt="product.name" />
            </div>

            <div class="col-md-10 d-flex flex-column">
              <div class="product-name">{{ product.name }}</div>

              <div class="product-description">
                {{ product.description }}
              </div>

              <div class="flex-grow-1"></div>

              <div class="d-flex justify-content-end align-items-center">
                <div class="product-price">
                  Цена: <strong>{{ formatPrice(product.price) }} руб.</strong>
                </div>

                <button class="btn btn-warning" title="В корзину" @click="addToCart(product.id)">
                  <i class="bi bi-cart"></i>
                  <span class="d-none d-md-inline">В корзину</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
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
