import type { User } from "~/types/auth";

export const useUserStore = defineStore("user", () => {
  const user = useState<User | null>("user", () => null);
  const avatarUrl = ref<string>("");

  return { user, avatarUrl };
});
