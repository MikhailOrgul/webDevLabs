<script setup>
import { computed, ref } from "vue";

const rating = ref(0);

const title = computed(() => {
  switch (rating.value) {
    case 1:
      return "Плохо";
    case 2:
      return "Так себе";
    case 3:
      return "Удовлетворительно";
    case 4:
      return "Хорошо";
    case 5:
      return "Отлично";
    default:
      return "Выберите оценку";
  }
});

const setRating = (value) => {
  rating.value = value;
};
</script>

<template>
  <div class="card shadow-sm">
    <div class="card-body text-center">
      <h1 class="h4 mb-4">Оценка приложения</h1>

      <div class="d-flex justify-content-center gap-2 mb-3 star-btn">
        <span
          v-for="n in 5"
          :key="n"
          role="button"
          tabindex="0"
          @click="setRating(n)"
          @keydown.enter="setRating(n)"
        >
          <i
            class="bi"
            :class="rating >= n ? 'bi-star-fill text-warning' : 'bi-star text-warning'"
          ></i>
        </span>
      </div>

      <div class="star-hint text-muted">
        {{ title }}
      </div>
    </div>
  </div>
</template>

<script>
export default {};
</script>
