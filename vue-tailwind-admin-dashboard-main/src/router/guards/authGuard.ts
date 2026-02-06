import { getMe } from "@/services/authService"

export async function authGuard() {
  try {
    await getMe()
    return true
  } catch {
    return "/signin"
  }
}
