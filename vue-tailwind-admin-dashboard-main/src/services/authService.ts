import http from "@/lib/http"
import type { LoginRequest } from "@/dto/auth/LoginRequest"
import type { LoginResponse } from "@/dto/auth/LoginResponse"

export function login(payload: LoginRequest) {
  return http.post<LoginResponse>("/signin", payload)
}

export function logout() {
  return http.post("/logout")
}

export interface RegisterRequest {
  name: string
  email: string
  password: string
  role: string
}

export function register(payload: RegisterRequest) {
  return http.post("/signup", payload)
}

export function getMe() {
  return http.get("/me")
}

export interface UpdateProfileRequest {
  name: string
  email: string
}

export function updateProfile(payload: UpdateProfileRequest) {
  return http.put("/me", payload)
}
