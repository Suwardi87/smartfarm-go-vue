<template>
  <div class="flex flex-col min-h-screen bg-gray-50 dark:bg-gray-900 text-gray-800 dark:text-white">
    <!-- Navbar -->
    <header class="sticky top-0 z-40 w-full bg-white border-b border-gray-200 dark:bg-gray-900 dark:border-gray-800">
      <div class="container mx-auto px-4 h-16 flex items-center justify-between">
        <!-- Logo -->
        <router-link to="/" class="flex items-center gap-2 font-bold text-xl text-brand-600 dark:text-brand-400">
          <span class="text-2xl">🌱</span> SmartFarm
        </router-link>

        <!-- Search (Hidden on mobile for now or simple icon) -->
        <div class="hidden md:flex flex-1 mx-8 max-w-lg">
          <div class="relative w-full">
            <input
              type="text"
              placeholder="Cari sayuran segar..."
              class="w-full pl-4 pr-10 py-2 rounded-lg border border-gray-300 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 focus:ring-2 focus:ring-brand-500"
            />
            <button class="absolute right-2 top-1/2 -translate-y-1/2 text-gray-500">
              🔍
            </button>
          </div>
        </div>

        <!-- Right Actions -->
        <div class="flex items-center gap-4">
          <!-- Cart -->
          <router-link to="/cart" class="relative p-2 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-full">
            🛒
            <span v-if="totalItems > 0" class="absolute top-0 right-0 bg-red-500 text-white text-xs w-4 h-4 rounded-full flex items-center justify-center">
              {{ totalItems }}
            </span>
          </router-link>

          <!-- Auth Buttons / User Menu -->
          <div v-if="user" class="flex items-center gap-4">
             <span class="text-sm font-medium">Hi, {{ user.name }}</span>
             <router-link v-if="user.role === 'petani'" to="/farmer/dashboard" class="text-sm text-green-600 font-bold hover:underline">
               Dashboard
             </router-link>
             <button @click="handleLogout" class="text-sm text-red-500 hover:text-red-700 font-medium">
               Keluar
             </button>
          </div>
          <div v-else class="flex items-center gap-4">
            <router-link to="/signin" class="text-sm font-medium hover:text-brand-600">
              Masuk
            </router-link>
            <router-link to="/signup" class="px-4 py-2 bg-brand-600 hover:bg-brand-700 text-white rounded-lg text-sm font-medium transition">
              Daftar
            </router-link>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="flex-1 container mx-auto px-4 py-8">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 dark:bg-gray-900 dark:border-gray-800 py-8">
      <div class="container mx-auto px-4 text-center text-gray-500 text-sm">
        &copy; {{ new Date().getFullYear() }} SmartFarm Marketplace. Melayani Petani & Pembeli.
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { useCart } from '@/stores/cart'
import { useUser } from '@/stores/user'
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const { totalItems } = useCart()
const userStore = useUser()
const user = computed(() => userStore.state.user)
const router = useRouter()

onMounted(() => {
  if (!userStore.state.isAuthenticated) {
    userStore.fetchUser()
  }
})

const handleLogout = async () => {
  await userStore.logout()
  router.push('/signin')
}
</script>
