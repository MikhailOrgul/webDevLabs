<script setup>
import { computed, onMounted, ref, watch } from "vue";
import Swal from "sweetalert2";

const STORAGE_KEY = "todolist_notes_v1";

const makeId = () => {
  return Math.random().toString(16).slice(2) + Date.now().toString(16);
};

const defaultNotes = () => {
  return [
    {
      id: makeId(),
      title: "Покупки",
      todos: [
        { id: makeId(), text: "Масляный фильтр", done: false },
        { id: makeId(), text: "Тормозная жидкость", done: true },
      ],
    },
    {
      id: makeId(),
      title: "Учёба",
      todos: [{ id: makeId(), text: "Лаба 2 по Vue", done: false }],
    },
  ];
};

const notes = ref([]);
const selectedNoteId = ref(null);

const selectedNote = computed(() => {
  return notes.value.find((n) => n.id === selectedNoteId.value) ?? null;
});

const selectedNoteIndex = computed(() => {
  return notes.value.findIndex((n) => n.id === selectedNoteId.value);
});

const newNoteTitle = ref("");
const newTodoText = ref("");

const selectNote = (id) => {
  selectedNoteId.value = id;
};

const loadFromStorage = () => {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) {
      return defaultNotes();
    }
    const parsed = JSON.parse(raw);
    if (!Array.isArray(parsed)) {
      return defaultNotes();
    }
    return parsed;
  } catch (e) {
    return defaultNotes();
  }
};

const saveToStorage = () => {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(notes.value));
};

onMounted(() => {
  notes.value = loadFromStorage();
  selectedNoteId.value = notes.value.at(0)?.id ?? null;
});

watch(notes, saveToStorage, { deep: true });

const addNote = async () => {
  const title = newNoteTitle.value.trim();
  if (!title) {
    await Swal.fire({
      icon: "warning",
      title: "Введите название заметки",
    });
    return;
  }

  const note = {
    id: makeId(),
    title,
    todos: [],
  };

  notes.value.unshift(note);
  selectedNoteId.value = note.id;
  newNoteTitle.value = "";
};

const renameNote = async () => {
  if (!selectedNote.value) {
    return;
  }

  const { value } = await Swal.fire({
    title: "Изменить название",
    input: "text",
    inputValue: selectedNote.value.title,
    showCancelButton: true,
    confirmButtonText: "Сохранить",
    cancelButtonText: "Отмена",
  });

  if (value === undefined) {
    return;
  }

  const title = String(value).trim();
  if (!title) {
    await Swal.fire({
      icon: "warning",
      title: "Название не может быть пустым",
    });
    return;
  }

  selectedNote.value.title = title;
};

const deleteNote = async () => {
  if (!selectedNote.value) {
    return;
  }

  const result = await Swal.fire({
    icon: "warning",
    title: "Удалить заметку?",
    text: `«${selectedNote.value.title}» будет удалена без возможности восстановления.`,
    showCancelButton: true,
    confirmButtonText: "Удалить",
    cancelButtonText: "Отмена",
  });

  if (!result.isConfirmed) {
    return;
  }

  const idx = selectedNoteIndex.value;
  if (idx < 0) {
    return;
  }

  notes.value.splice(idx, 1);

  selectedNoteId.value = notes.value.at(0)?.id ?? null;
};

const addTodo = async () => {
  if (!selectedNote.value) {
    return;
  }

  const text = newTodoText.value.trim();
  if (!text) {
    await Swal.fire({
      icon: "warning",
      title: "Введите текст задачи",
    });
    return;
  }

  selectedNote.value.todos.push({
    id: makeId(),
    text,
    done: false,
  });

  newTodoText.value = "";
};

const deleteTodo = async (todoId) => {
  if (!selectedNote.value) {
    return;
  }

  const idx = selectedNote.value.todos.findIndex((t) => t.id === todoId);
  if (idx < 0) {
    return;
  }

  selectedNote.value.todos.splice(idx, 1);
};

const editTodo = async (todoId) => {
  if (!selectedNote.value) {
    return;
  }

  const todo = selectedNote.value.todos.find((t) => t.id === todoId);
  if (!todo) {
    return;
  }

  const { value } = await Swal.fire({
    title: "Редактировать задачу",
    input: "text",
    inputValue: todo.text,
    showCancelButton: true,
    confirmButtonText: "Сохранить",
    cancelButtonText: "Отмена",
  });

  if (value === undefined) {
    return;
  }

  const text = String(value).trim();
  if (!text) {
    await Swal.fire({
      icon: "warning",
      title: "Текст не может быть пустым",
    });
    return;
  }

  todo.text = text;
};

const doneCount = computed(() => {
  if (!selectedNote.value) {
    return 0;
  }
  return selectedNote.value.todos.filter((t) => t.done).length;
});
</script>

<template>
  <div class="row g-3">
    <div class="col-12 col-lg-4">
      <div class="d-flex gap-2 mb-2">
        <input
          v-model="newNoteTitle"
          type="text"
          class="form-control"
          placeholder="Новая заметка..."
          @keydown.enter="addNote"
        />
        <button class="btn btn-primary" type="button" @click="addNote" title="Создать">
          <i class="bi bi-plus-lg"></i>
        </button>
      </div>

      <div class="list-group">
        <button
          v-for="note in notes"
          :key="note.id"
          type="button"
          class="list-group-item list-group-item-action note-list-item"
          :class="{ active: note.id === selectedNoteId }"
          @click="selectNote(note.id)"
        >
          <div class="d-flex justify-content-between align-items-center">
            <div class="fw-semibold text-truncate">
              {{ note.title }}
            </div>
            <span class="badge bg-secondary ms-2">
              {{ note.todos.length }}
            </span>
          </div>
        </button>

        <div v-if="notes.length === 0" class="text-muted small p-3 border rounded">
          Заметок пока нет. Создай первую слева сверху.
        </div>
      </div>
    </div>

    <div class="col-12 col-lg-8">
      <div v-if="!selectedNote" class="alert alert-info">
        Выберите заметку слева или создайте новую.
      </div>

      <div v-else class="card shadow-sm">
        <div class="card-body">
          <div class="d-flex justify-content-between align-items-start gap-2 mb-3">
            <div>
              <h2 class="h4 mb-1">{{ selectedNote.title }}</h2>
              <div class="text-muted small">
                Выполнено: {{ doneCount }} / {{ selectedNote.todos.length }}
              </div>
            </div>

            <div class="d-flex gap-2">
              <button class="btn btn-outline-secondary btn-sm" type="button" @click="renameNote">
                <i class="bi bi-pencil"></i>
                <span class="d-none d-md-inline">Переименовать</span>
              </button>
              <button class="btn btn-outline-danger btn-sm" type="button" @click="deleteNote">
                <i class="bi bi-trash"></i>
                <span class="d-none d-md-inline">Удалить</span>
              </button>
            </div>
          </div>

          <div class="d-flex gap-2 mb-3">
            <input
              v-model="newTodoText"
              type="text"
              class="form-control"
              placeholder="Новая задача..."
              @keydown.enter="addTodo"
            />
            <button class="btn btn-success" type="button" @click="addTodo" title="Добавить">
              <i class="bi bi-plus-lg"></i>
            </button>
          </div>

          <div v-if="selectedNote.todos.length === 0" class="text-muted">
            В этой заметке пока нет задач. Добавь первую.
          </div>

          <ul v-else class="list-group">
            <li
              v-for="todo in selectedNote.todos"
              :key="todo.id"
              class="list-group-item d-flex justify-content-between align-items-center"
            >
              <div class="form-check d-flex align-items-center gap-2">
                <input class="form-check-input" type="checkbox" :id="todo.id" v-model="todo.done" />
                <label class="form-check-label" :for="todo.id" :class="{ 'todo-done': todo.done }">
                  {{ todo.text }}
                </label>
              </div>

              <div class="d-flex gap-2">
                <button
                  class="btn btn-outline-secondary btn-sm"
                  type="button"
                  title="Редактировать"
                  @click="editTodo(todo.id)"
                >
                  <i class="bi bi-pencil"></i>
                </button>

                <button
                  class="btn btn-outline-danger btn-sm"
                  type="button"
                  title="Удалить"
                  @click="deleteTodo(todo.id)"
                >
                  <i class="bi bi-x-lg"></i>
                </button>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {};
</script>
