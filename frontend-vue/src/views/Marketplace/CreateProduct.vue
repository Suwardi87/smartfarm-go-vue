<template>
   <AdminLayout>
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">Tambah Produk Baru</h1>

      <div class="bg-white dark:bg-gray-800 p-8 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm">
        <form @submit.prevent="handleSubmit" class="space-y-6">
          
          <!-- Name -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Nama Produk</label>
            <input v-model="form.name" type="text" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
          </div>

          <!-- Price & Stock -->
          <div class="grid grid-cols-2 gap-6">
             <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Harga (Rp)</label>
                <input v-model.number="form.price" type="number" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
             </div>
             <div>
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Stok Awal</label>
                <input v-model.number="form.stock" type="number" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
             </div>
          </div>

          <!-- Description -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Deskripsi</label>
            <textarea v-model="form.description" rows="4" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white"></textarea>
          </div>

          <!-- Feature Flags -->
          <div class="space-y-4 border-t border-gray-200 dark:border-gray-700 pt-4">
             <div class="flex items-start gap-3">
                <input v-model="form.isPreOrder" type="checkbox" id="preorder" class="mt-1 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
                <div>
                   <label for="preorder" class="font-medium text-gray-900 dark:text-white">Pre-Order System</label>
                   <p class="text-sm text-gray-500">Aktifkan jika produk belum panen. Pembeli akan membayar di muka.</p>
                </div>
             </div>
             
             <div v-if="form.isPreOrder">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Estimasi Tanggal Panen</label>
                <input v-model="form.harvestDate" type="date" required class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white" />
             </div>

             <div class="flex items-start gap-3">
                <input v-model="form.isSubscription" type="checkbox" id="sub" class="mt-1 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
                <div>
                   <label for="sub" class="font-medium text-gray-900 dark:text-white">Paket Langganan</label>
                   <p class="text-sm text-gray-500">Produk ini akan dikirim rutin (Mingguan/Bulanan).</p>
                </div>
             </div>

             <div v-if="form.isSubscription">
                <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Periode Langganan</label>
                <select v-model="form.subscriptionPeriod" class="w-full rounded-lg border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white">
                   <option value="weekly">Mingguan</option>
                   <option value="monthly">Bulanan</option>
                </select>
             </div>
          </div>

          <!-- Image Upload -->
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">Foto Produk</label>
            <input type="file" @change="handleFileChange" accept="image/*" class="w-full text-sm text-gray-500 file:mr-4 file:py-2 file:px-4 file:rounded-full file:border-0 file:text-sm file:font-semibold file:bg-brand-50 file:text-brand-700 hover:file:bg-brand-100" />
          </div>

          <!-- Submit -->
          <button type="submit" :disabled="loading" class="w-full bg-brand-600 text-white font-bold py-3 rounded-lg hover:bg-brand-700 transition shadow-lg shadow-brand-500/30 disabled:opacity-50">
             {{ loading ? 'Menyimpan...' : 'Jual Produk Sekarang' }}
          </button>
        </form>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { createProduct } from '@/services/productService'
import AdminLayout from '@/components/layout/AdminLayout.vue'

const router = useRouter()
const loading = ref(false)
const file = ref<File | null>(null)

const form = reactive({
  name: '',
  price: 0,
  stock: 10,
  description: '',
  isPreOrder: false,
  harvestDate: '',
  isSubscription: false,
  subscriptionPeriod: 'weekly'
})

const handleFileChange = (e: Event) => {
    const input = e.target as HTMLInputElement
    if (input.files && input.files[0]) {
        file.value = input.files[0]
    }
}

const handleSubmit = async () => {
    loading.value = true
    try {
        const formData = new FormData()
        formData.append('name', form.name)
        formData.append('price', form.price.toString())
        formData.append('stock', form.stock.toString())
        formData.append('description', form.description)
        formData.append('category', 'Vegetables') // Default category
        
        if(file.value) {
            formData.append('image', file.value)
        }

        if (form.isPreOrder) {
            formData.append('is_pre_order', 'true')
            formData.append('harvest_date', new Date(form.harvestDate).toISOString())
        }

        if (form.isSubscription) {
            formData.append('is_subscription', 'true')
            formData.append('subscription_period', form.subscriptionPeriod)
        }

        await createProduct(formData)
        alert('Produk berhasil ditambahkan!')
        router.push('/farmer/dashboard')

    } catch (error: any) {
        console.error("Failed to create product", error)
        alert('Gagal: ' + (error.response?.data?.error || error.message))
    } finally {
        loading.value = false
    }
}
</script>
