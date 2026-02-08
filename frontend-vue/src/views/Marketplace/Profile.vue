<template>
  <AdminLayout>
    <div class="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Profil Saya</h1>
        <p class="text-gray-600 dark:text-gray-400">Kelola informasi akun dan pengaturan Anda</p>
      </div>

      <!-- Loading -->
      <LoadingSpinner :isLoading="isLoading" message="Memuat profil..." />

      <!-- Profile Form -->
      <div v-if="!isLoading" class="bg-white dark:bg-gray-800 p-8 rounded-xl border border-gray-200 dark:border-gray-700">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Name -->
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Nama Lengkap
            </label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
              placeholder="Nama lengkap Anda"
            />
            <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
          </div>

          <!-- Email -->
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Email
            </label>
            <input
              id="email"
              v-model="form.email"
              type="email"
              required
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
              placeholder="Email Anda"
            />
            <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
          </div>

          <!-- Role (Read-only) -->
          <div>
            <label for="role" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
              Tipe Akun
            </label>
            <input
              id="role"
              :value="user?.role"
              type="text"
              disabled
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-gray-400 bg-gray-100 cursor-not-allowed"
            />
            <p class="mt-1 text-xs text-gray-500">Hubungi admin untuk mengubah tipe akun</p>
          </div>

          <!-- Buttons -->
          <div class="flex gap-4 pt-6 border-t border-gray-200 dark:border-gray-700">
            <button
              type="submit"
              :disabled="isSaving"
              class="flex-1 bg-brand-600 text-white font-semibold py-3 px-6 rounded-lg hover:bg-brand-700 transition disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {{ isSaving ? 'Menyimpan...' : 'Simpan Perubahan' }}
            </button>
            <button
              type="button"
              @click="handleCancel"
              class="flex-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 font-semibold py-3 px-6 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition"
            >
              Batal
            </button>
          </div>
        </form>
      </div>

      <!-- Other Settings -->
      <div class="mt-12 space-y-6">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white">Pengaturan Lainnya</h2>

        <!-- Change Password (Placeholder) -->
        <div class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Ubah Kata Sandi</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-4">Perbarui kata sandi akun Anda secara berkala untuk keamanan</p>
          <button class="text-brand-600 hover:text-brand-700 font-medium">
            Ubah Kata Sandi →
          </button>
        </div>

        <!-- Addresses (Placeholder) -->
        <div class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-2">Alamat Pengiriman</h3>
          <p class="text-gray-600 dark:text-gray-400 mb-4">Kelola alamat pengiriman untuk pesanan Anda</p>
          <button class="text-brand-600 hover:text-brand-700 font-medium">
            Kelola Alamat →
          </button>
        </div>

        <!-- Logout -->
        <div class="bg-red-50 dark:bg-red-900/20 p-6 rounded-xl border border-red-200 dark:border-red-800">
          <h3 class="text-lg font-semibold text-red-900 dark:text-red-300 mb-2">Keluar dari Akun</h3>
          <p class="text-red-800 dark:text-red-200 mb-4">Anda akan keluar dari semua perangkat</p>
          <button
            @click="handleLogout"
            class="text-red-600 dark:text-red-400 hover:text-red-700 dark:hover:text-red-300 font-medium"
          >
            Keluar Sekarang →
          </button>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { getMe, updateProfile } from '@/services/authService'
import type { User } from '@/dto/auth/User'
import { useToast } from '@/composables/useToast'
import AdminLayout from '@/components/layout/AdminLayout.vue'

const router = useRouter()
const { showSuccess, showError } = useToast()

const user = ref<User | null>(null)
const isLoading = ref(true)
const isSaving = ref(false)

const form = reactive({
  name: '',
  email: ''
})

const errors = reactive({
  name: '',
  email: ''
})

onMounted(async () => {
  try {
    isLoading.value = true
    const response = await getMe()
    user.value = response.data.data
    if (user.value) {
      form.name = user.value.name
      form.email = user.value.email
    }
  } catch (error: any) {
    console.error('Failed to fetch profile', error)
    showError('Gagal memuat profil. Silakan coba lagi.')
    if (error.response?.status === 401) {
      router.push('/signin')
    }
  } finally {
    isLoading.value = false
  }
})

const validateForm = () => {
  errors.name = ''
  errors.email = ''

  if (!form.name.trim()) {
    errors.name = 'Nama harus diisi'
  }
  if (!form.email.trim()) {
    errors.email = 'Email harus diisi'
  } else if (!form.email.includes('@')) {
    errors.email = 'Format email tidak valid'
  }

  return !errors.name && !errors.email
}

const handleSubmit = async () => {
  if (!validateForm()) return

  try {
    isSaving.value = true
    await updateProfile({
      name: form.name,
      email: form.email
    })
    showSuccess('Profil berhasil diperbarui!')
    // Reload user data
    const response = await getMe()
    user.value = response.data.data
  } catch (error: any) {
    console.error('Failed to update profile', error)
    showError(error.response?.data?.error || 'Gagal menyimpan profil')
  } finally {
    isSaving.value = false
  }
}

const handleCancel = () => {
  if (user.value) {
    form.name = user.value.name
    form.email = user.value.email
    errors.name = ''
    errors.email = ''
  }
}

const handleLogout = async () => {
  if (confirm('Apakah Anda yakin ingin keluar?')) {
    try {
      await useRouter().push('/signin')
    } catch (error) {
      console.error('Logout failed', error)
    }
  }
}
</script>
