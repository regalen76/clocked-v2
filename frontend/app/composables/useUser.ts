import { useAuth } from "./useAuth";
import type { UserResponse } from "~/types/auth";
import { useUserStore } from "./useUserStore";

export const useUser = () => {
  const { token } = useAuth();
  const userStore = useUserStore();

  async function fetchUser() {
    if (!token.value) return;

    const { data, error } = await useFetch<UserResponse>(
      "http://localhost:8000/api/user",
      {
        headers: {
          Authorization: `Bearer ${token.value}`,
        },
      },
    );

    if (error.value) {
      console.error(error.value);
      return;
    }

    if (data.value) {
      userStore.$patch({ user: data.value.data });
      fetchAvatar();
    }
  }

  async function fetchAvatar() {
    if (!userStore.user) return;

    const res = await fetch(
      `http://localhost:8000/api/user/${userStore.user.ID}/avatar`,
      {
        headers: {
          Authorization: `Bearer ${useAuth().token.value}`,
        },
      },
    );

    if (res.ok) {
      const blob = await res.blob();
      userStore.$patch({ avatarUrl: URL.createObjectURL(blob) });
    } else {
      userStore.$patch({ avatarUrl: "" });
    }
  }

  return { fetchUser, fetchAvatar };
};
