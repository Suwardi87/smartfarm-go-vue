<template>
  <MarketplaceLayout>
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="mb-8 flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-2">Alamat Pengiriman</h1>
          <p class="text-gray-600 dark:text-gray-400">Kelola alamat pengiriman untuk pesanan Anda</p>
        </div>
        <button
          @click="showForm = !showForm"
          class="bg-brand-600 text-white px-6 py-3 rounded-lg hover:bg-brand-700 transition font-semibold"
        >
          {{ showForm ? 'Batal' : '+ Tambah Alamat' }}
        </button>
      </div>

      <!-- Loading -->
      <LoadingSpinner :isLoading="isLoading" message="Memuat alamat..." />

      <!-- Form -->
      <div v-if="showForm && !isLoading" class="bg-white dark:bg-gray-800 p-8 rounded-xl border border-gray-200 dark:border-gray-700 mb-8">
        <h2 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
          {{ editingAddress ? 'Edit Alamat' : 'Alamat Baru' }}
        </h2>

        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Label -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Label</label>
            <input
              v-model="form.label"
              type="text"
              placeholder="Rumah, Kantor, dll"
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
            />
            <p v-if="errors.label" class="mt-1 text-sm text-red-600">{{ errors.label }}</p>
          </div>

          <!-- Recipient Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nama Penerima</label>
            <input
              v-model="form.recipient_name"
              type="text"
              required
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
            />
            <p v-if="errors.recipient_name" class="mt-1 text-sm text-red-600">{{ errors.recipient_name }}</p>
          </div>

          <!-- Phone -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nomor Telepon</label>
            <input
              v-model="form.phone_number"
              type="tel"
              required
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
            />
            <p v-if="errors.phone_number" class="mt-1 text-sm text-red-600">{{ errors.phone_number }}</p>
          </div>

          <!-- Street -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Alamat Lengkap</label>
            <textarea
              v-model="form.street"
              required
              rows="3"
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
              placeholder="Jalan, nomor rumah, dll"
            ></textarea>
            <p v-if="errors.street" class="mt-1 text-sm text-red-600">{{ errors.street }}</p>
          </div>

          <!-- Province & City -->
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Provinsi</label>
              <input
                v-model="form.province"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
              />
              <p v-if="errors.province" class="mt-1 text-sm text-red-600">{{ errors.province }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Kota</label>
              <input
                v-model="form.city"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
              />
              <p v-if="errors.city" class="mt-1 text-sm text-red-600">{{ errors.city }}</p>
            </div>
          </div>

          <!-- Postal Code -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Kode Pos</label>
            <input
              v-model="form.postal_code"
              type="text"
              required
              class="w-full px-4 py-2 border border-gray-300 dark:border-gray-600 rounded-lg dark:bg-gray-700 dark:text-white focus:ring-2 focus:ring-brand-600 focus:border-transparent"
            />
            <p v-if="errors.postal_code" class="mt-1 text-sm text-red-600">{{ errors.postal_code }}</p>
          </div>

          <!-- Default Address -->
          <div class="flex items-center gap-3">
            <input
              id="is_default"
              v-model="form.is_default"
              type="checkbox"
              class="w-4 h-4 border-gray-300 rounded accent-brand-600"
            />
            <label for="is_default" class="text-sm text-gray-700 dark:text-gray-300">Jadikan alamat utama</label>
          </div>

          <!-- Buttons -->
          <div class="flex gap-4 pt-6 border-t border-gray-200 dark:border-gray-700">
            <button
              type="submit"
              :disabled="isSaving"
              class="flex-1 bg-brand-600 text-white font-semibold py-3 px-6 rounded-lg hover:bg-brand-700 transition disabled:opacity-50"
            >
              {{ isSaving ? 'Menyimpan...' : 'Simpan Alamat' }}
            </button>
            <button
              type="button"
              @click="resetForm"
              class="flex-1 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-300 font-semibold py-3 px-6 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-600 transition"
            >
              Batal
            </button>
          </div>
        </form>
      </div>

      <!-- Addresses List -->
      <div v-if="!isLoading" class="space-y-4">
        <div v-if="addresses.length === 0" class="text-center py-12">
          <div class="text-6xl mb-4">📍</div>
          <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Belum ada alamat</h3>
          <p class="text-gray-500 dark:text-gray-400">Tambahkan alamat pengiriman untuk melanjutkan checkout</p>
        </div>

        <div v-else v-for="address in addresses" :key="address.id" class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-700 hover:border-brand-600 transition">
          <!-- Header -->
          <div class="flex items-start justify-between mb-4">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-2">
                <h3 class="text-lg font-bold text-gray-900 dark:text-white">{{ address.label }}</h3>
                <span v-if="address.is_default" class="px-2 py-1 bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-300 text-xs font-semibold rounded">
                  Alamat Utama
                </span>
              </div>
              <p class="text-gray-600 dark:text-gray-300 mb-2">{{ address.recipient_name }} • {{ address.phone_number }}</p>
            </div>
            <div class="flex gap-2">
              <button
                @click="editAddress(address)"
                class="text-brand-600 hover:text-brand-700 font-medium text-sm"
              >
                Edit
              </button>
              <button
                @click="deleteAddressHandler(address.id)"
                class="text-red-600 hover:text-red-700 font-medium text-sm"
              >
                Hapus
              </button>
            </div>
          </div>

          <!-- Address Details -->
          <div class="text-gray-600 dark:text-gray-400 text-sm space-y-1 mb-4">
            <p>{{ address.street }}</p>
            <p>{{ address.city }}, {{ address.province }} {{ address.postal_code }}</p>
          </div>

          <!-- Set Default Button -->
          <button
            v-if="!address.is_default"
            @click="setDefaultHandler(address.id)"
            class="text-brand-600 hover:text-brand-700 font-medium text-sm"
          >
            Jadikan Alamat Utama
          </button>
        </div>
      </div>
    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { useRouter } from 'vue-router'
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import LoadingSpinner from '@/components/common/LoadingSpinner.vue'
import { getMyAddresses, createAddress, updateAddress, deleteAddress, setDefaultAddress } from '@/services/addressService'
import type { Address, CreateAddressRequest } from '@/dto/address/Address'
import { useToast } from '@/composables/useToast'
import { normalizeAddress } from '@/utils/transformers'

const router = useRouter()
const { showSuccess, showError } = useToast()

const addresses = ref<Address[]>([])
const isLoading = ref(true)
const isSaving = ref(false)
const showForm = ref(false)
const editingAddress = ref<Address | null>(null)

const form = reactive<CreateAddressRequest>({
  label: '',
  recipient_name: '',
  phone_number: '',
  street: '',
  city: '',
  province: '',
  postal_code: '',
  is_default: false
})

const errors = reactive({
  label: '',
  recipient_name: '',
  phone_number: '',
  street: '',
  city: '',
  province: '',
  postal_code: ''
})

onMounted(async () => {
  try {
    isLoading.value = true
    const response = await getMyAddresses()
    const rawData = response.data.data
    addresses.value = (rawData || []).map(normalizeAddress)
  } catch (error: any) {
    console.error('Failed to fetch addresses', error)
    showError('Gagal memuat alamat')
    if (error.response?.status === 401) {
      router.push('/signin')
    }
  } finally {
    isLoading.value = false
  }
})

const validateForm = () => {
  errors.label = ''
  errors.recipient_name = ''
  errors.phone_number = ''
  errors.street = ''
  errors.city = ''
  errors.province = ''
  errors.postal_code = ''

  if (!form.label.trim()) errors.label = 'Label harus diisi'
  if (!form.recipient_name.trim()) errors.recipient_name = 'Nama penerima harus diisi'
  if (!form.phone_number.trim()) errors.phone_number = 'Nomor telepon harus diisi'
  if (!form.street.trim()) errors.street = 'Alamat harus diisi'
  if (!form.city.trim()) errors.city = 'Kota harus diisi'
  if (!form.province.trim()) errors.province = 'Provinsi harus diisi'
  if (!form.postal_code.trim()) errors.postal_code = 'Kode pos harus diisi'

  return !Object.values(errors).some(e => e)
}

const resetForm = () => {
  form.label = ''
  form.recipient_name = ''
  form.phone_number = ''
  form.street = ''
  form.city = ''
  form.province = ''
  form.postal_code = ''
  form.is_default = false
  editingAddress.value = null
  showForm.value = false
}

const handleSubmit = async () => {
  if (!validateForm()) return

  try {
    isSaving.value = true

    if (editingAddress.value) {
      await updateAddress(editingAddress.value.id, form)
      showSuccess('Alamat berhasil diperbarui')
    } else {
      await createAddress(form)
      showSuccess('Alamat berhasil ditambahkan')
    }

    // Reload addresses
    const response = await getMyAddresses()
    const rawData = response.data.data
    addresses.value = (rawData || []).map(normalizeAddress)
    resetForm()
  } catch (error: any) {
    console.error('Failed to save address', error)
    showError(error.response?.data?.error || 'Gagal menyimpan alamat')
  } finally {
    isSaving.value = false
  }
}

const editAddress = (address: Address) => {
  editingAddress.value = address
  form.label = address.label
  form.recipient_name = address.recipient_name
  form.phone_number = address.phone_number
  form.street = address.street
  form.city = address.city
  form.province = address.province
  form.postal_code = address.postal_code
  form.is_default = address.is_default
  showForm.value = true
}

const deleteAddressHandler = async (addressId: number) => {
  if (confirm('Yakin ingin menghapus alamat ini?')) {
    try {
      await deleteAddress(addressId)
      showSuccess('Alamat berhasil dihapus')
      const response = await getMyAddresses()
      const rawData = response.data.data
      addresses.value = (rawData || []).map(normalizeAddress)
    } catch (error: any) {
      console.error('Failed to delete address', error)
      showError('Gagal menghapus alamat')
    }
  }
}

const setDefaultHandler = async (addressId: number) => {
  try {
    await setDefaultAddress(addressId)
    showSuccess('Alamat utama berhasil diubah')
    const response = await getMyAddresses()
    const rawData = response.data.data
    addresses.value = (rawData || []).map(normalizeAddress)
  } catch (error: any) {
    console.error('Failed to set default address', error)
    showError('Gagal mengubah alamat utama')
  }
}
</script>
