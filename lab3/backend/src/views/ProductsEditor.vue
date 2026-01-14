<template>
  <main class="container mt-4">
    <h1>Редактор товаров</h1>
    <p class="text-muted">Доступно только администраторам.</p>

    <form class="card shadow-sm p-4 mb-4" @submit.prevent="handleSubmit">
      <div class="row g-3">
        <div class="col-md-6">
          <label class="form-label" for="product-name">Наименование</label>
          <input id="product-name" v-model.trim="form.name" class="form-control" required />
        </div>
        <div class="col-md-6">
          <label class="form-label" for="product-price">Цена</label>
          <input
            id="product-price"
            v-model.number="form.price"
            type="number"
            min="1"
            class="form-control"
            required
          />
        </div>
        <div class="col-12">
          <label class="form-label" for="product-description">Описание</label>
          <textarea
            id="product-description"
            v-model.trim="form.description"
            class="form-control"
            rows="3"
            required
          ></textarea>
        </div>
        <div class="col-12">
          <label class="form-label" for="product-image">Ссылка на изображение</label>
          <input id="product-image" v-model.trim="form.image" class="form-control" />
        </div>
      </div>

      <div class="mt-4 d-flex gap-2">
        <button class="btn btn-primary" type="submit" :disabled="isSubmitting">
          {{ editingId ? "Сохранить" : "Добавить" }}
        </button>
        <button v-if="editingId" class="btn btn-outline-secondary" type="button" @click="resetForm">
          Отмена
        </button>
      </div>

      <p v-if="errorMessage" class="text-danger mt-3">{{ errorMessage }}</p>
    </form>

    <div class="table-responsive">
      <table class="table table-striped align-middle">
        <thead class="table-light">
          <tr>
            <th>#</th>
            <th>Наименование</th>
            <th>Цена</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="product in products" :key="product.id">
            <td>{{ product.id }}</td>
            <td>{{ product.name }}</td>
            <td>{{ formatPrice(product.price) }} руб.</td>
            <td class="text-end">
              <button class="btn btn-outline-primary btn-sm me-2" @click="startEdit(product)">
                Редактировать
              </button>
              <button class="btn btn-outline-danger btn-sm" @click="removeProduct(product.id)">
                Удалить
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </main>
</template>

<script setup>
import { onMounted, ref } from "vue";

const products = ref([]);
const editingId = ref(null);
const isSubmitting = ref(false);
const errorMessage = ref("");

const form = ref({
  name: "",
  description: "",
  price: 0,
  image: "",
});

const formatPrice = (price) => new Intl.NumberFormat("ru-RU").format(price);

const loadProducts = async () => {
  const response = await fetch("/api/products?page=1&limit=100");
  if (!response.ok) {
    return;
  }
  const data = await response.json();
  products.value = data.items ?? [];
};

const resetForm = () => {
  editingId.value = null;
  form.value = { name: "", description: "", price: 0, image: "" };
  errorMessage.value = "";
};

const startEdit = (product) => {
  editingId.value = product.id;
  form.value = {
    name: product.name,
    description: product.description,
    price: product.price,
    image: product.image,
  };
};

const handleSubmit = async () => {
  errorMessage.value = "";
  isSubmitting.value = true;
  try {
    const method = editingId.value ? "PUT" : "POST";
    const url = editingId.value ? `/api/products/${editingId.value}` : "/api/products";
    const response = await fetch(url, {
      method,
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(form.value),
    });
    if (!response.ok) {
      const data = await response.json().catch(() => null);
      const message = data?.error ?? "Не удалось сохранить товар.";
      throw new Error(message);
    }
    await loadProducts();
    resetForm();
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : "Ошибка сохранения.";
  } finally {
    isSubmitting.value = false;
  }
};

const removeProduct = async (id) => {
  await fetch(`/api/products/${id}`, { method: "DELETE" });
  loadProducts();
};

onMounted(() => {
  loadProducts();
});
</script>
