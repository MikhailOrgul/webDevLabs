import { reactive } from "vue";

//password admin123

export const authState = reactive({
  user: null,
  isLoading: false,
});

export const fetchMe = async () => {
  if (authState.isLoading) {
    return authState.user;
  }
  authState.isLoading = true;
  try {
    const response = await fetch("/api/me");
    if (!response.ok) {
      authState.user = null;
      return null;
    }
    authState.user = await response.json();
    return authState.user;
  } finally {
    authState.isLoading = false;
  }
};

export const login = async (email, password) => {
  const response = await fetch("/api/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  });
  if (!response.ok) {
    const data = await response.json().catch(() => null);
    const message = data?.error ?? "Не удалось войти.";
    throw new Error(message);
  }
  authState.user = await response.json();
  return authState.user;
};

export const register = async (name, email, password) => {
  const response = await fetch("/api/register", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email, password }),
  });
  if (!response.ok) {
    const data = await response.json().catch(() => null);
    const message = data?.error ?? "Не удалось зарегистрироваться.";
    throw new Error(message);
  }
  authState.user = await response.json();
  return authState.user;
};

export const logout = async () => {
  await fetch("/api/logout", { method: "POST" });
  authState.user = null;
};
