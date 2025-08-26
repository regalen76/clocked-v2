import { useAuth } from "~/composables/useAuth";

export default defineNuxtRouteMiddleware((to, _from) => {
  const { isLoggedIn } = useAuth();

  if (!isLoggedIn.value && to.path !== "/login") {
    return navigateTo("/login");
  }
});
