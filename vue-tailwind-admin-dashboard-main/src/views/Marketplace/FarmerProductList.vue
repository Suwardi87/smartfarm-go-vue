<template>
  <MarketplaceLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <div class="flex justify-between items-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 dark:text-white">Daftar Produk Saya</h1>
        <router-link to="/products/create" class="px-6 py-3 bg-brand-600 text-white rounded-lg font-bold hover:bg-brand-700 transition">
          + Tambah Produk
        </router-link>
      </div>

      <div v-if="loading" class="flex justify-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-brand-600"></div>
      </div>

      <div v-else-if="products.length === 0" class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 p-20 text-center">
        <div class="text-6xl mb-4">📦</div>
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-2">Belum ada produk</h3>
        <p class="text-gray-500 mb-6">Mulai jual hasil tani Anda sekarang!</p>
        <router-link to="/products/create" class="inline-block px-6 py-3 bg-brand-600 text-white rounded-lg font-bold hover:bg-brand-700 transition">
          Jual Produk Pertama
        </router-link>
      </div>

      <div v-else class="bg-white dark:bg-gray-800 rounded-xl border border-gray-200 dark:border-gray-700 overflow-hidden shadow-sm">
        <div class="overflow-x-auto">
          <table class="w-full text-left border-collapse">
            <thead class="bg-gray-50 dark:bg-gray-700/50">
              <tr>
                <th class="px-6 py-4 text-sm font-semibold text-gray-900 dark:text-white">Produk</th>
                <th class="px-6 py-4 text-sm font-semibold text-gray-900 dark:text-white">Kategori</th>
                <th class="px-6 py-4 text-sm font-semibold text-gray-900 dark:text-white">Harga</th>
                <th class="px-6 py-4 text-sm font-semibold text-gray-900 dark:text-white">Stok</th>
                <th class="px-6 py-4 text-sm font-semibold text-gray-900 dark:text-white text-right">Aksi</th>
              </tr>
            </thead>
            <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
              <tr v-for="product in products" :key="product.id" class="hover:bg-gray-50 dark:hover:bg-gray-700/30 transition">
                <td class="px-6 py-4">
                  <div class="flex items-center gap-4">
                    <div class="w-12 h-12 rounded-lg bg-gray-100 dark:bg-gray-700 overflow-hidden flex-shrink-0">
                      <img v-if="product.image_url" :src="getImageUrl(product.image_url)" :alt="product.name" class="w-full h-full object-cover">
                      <div v-else class="w-full h-full flex items-center justify-center text-xl">🥦</div>
                    </div>
                    <div>
                      <div class="font-bold text-gray-900 dark:text-white">{{ product.name }}</div>
                      <div class="text-xs text-gray-500 line-clamp-1 max-w-[200px]">{{ product.description }}</div>
                    </div>
                  </div>
                </td>
                <td class="px-6 py-4">
                  <span class="px-2 py-1 text-xs font-medium rounded-full bg-brand-50 text-brand-700 dark:bg-brand-900/30 dark:text-brand-400 capitalize">
                    {{ product.category }}
                  </span>
                </td>
                <td class="px-6 py-4 text-gray-900 dark:text-white font-medium">
                  {{ formatRupiah(product.price) }}
                </td>
                <td class="px-6 py-4">
                  <div class="flex flex-col">
                    <span :class="[
                      'text-sm font-bold',
                      product.stock > 0 ? 'text-gray-900 dark:text-white' : 'text-red-500'
                    ]">
                      {{ product.stock }}
                    </span>
                    <span v-if="product.is_pre_order" class="text-[10px] text-orange-500 font-bold uppercase">Pre-Order</span>
                  </div>
                </td>
                <td class="px-6 py-4 text-right space-x-2">
                  <router-link :to="`/products/edit/${product.id}`" class="p-2 text-gray-400 hover:text-brand-600 transition inline-block" title="Edit">
                    ✏️
                  </router-link>
                  <button @click="handleDelete(product.id)" class="p-2 text-gray-400 hover:text-red-500 transition" title="Hapus">
                    🗑️
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import { ref, onMounted } from 'vue'
import { getFarmerProducts, deleteProduct } from '@/services/productService'
import type { Product } from '@/dto/product/Product'
import { formatRupiah } from '@/utils/formatter'
import { getImageUrl } from '@/utils/image'

const products = ref<Product[]>([])
const loading = ref(true)

const fetchProducts = async () => {
  try {
    const response = await getFarmerProducts()
    products.value = response.data.data ?? []
  } catch (error) {
    console.error('Failed to fetch farmer products', error)
  } finally {
    loading.value = false
  }
}

const handleDelete = async (id: number) => {
    if (!confirm('Apakah Anda yakin ingin menghapus produk ini?')) return

    try {
        await deleteProduct(id)
        alert('Produk berhasil dihapus')
        fetchProducts() // Refresh list
    } catch (error: any) {
        console.error('Failed to delete product', error)
        alert('Gagal menghapus produk: ' + (error.response?.data?.error || error.message))
    }
}

onMounted(fetchProducts)
</script>
