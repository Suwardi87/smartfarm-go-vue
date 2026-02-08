<template>
  <div class="group bg-white dark:bg-gray-800 rounded-xl overflow-hidden border border-gray-100 dark:border-gray-700 hover:shadow-lg transition flex flex-col h-full">
    <router-link :to="`/products/${product.id}`" class="block aspect-square bg-gray-200 dark:bg-gray-700 relative overflow-hidden">
       <!-- Image -->
       <img
        v-if="product.image_url && !imageError"
         :src="getImageUrl(product.image_url)"
        :alt="product.name"
        class="w-full h-full object-cover group-hover:scale-110 transition duration-500"
        @error="handleImageError"
        loading="lazy"
       />

       <div v-else class="absolute inset-0 flex items-center justify-center text-gray-400 bg-gradient-to-br from-green-50 to-green-100 dark:from-gray-700 dark:to-gray-800">
         <span class="text-6xl">🥬</span>
       </div>
       <!-- Label -->
       <span
        v-if="label"
        class="absolute top-2 left-2 text-white text-xs font-bold px-2 py-1 rounded shadow-sm"
        :class="labelColor"
       >
         {{ label }}
       </span>
    </router-link>

    <div class="p-4 flex flex-col flex-1">
      <router-link :to="`/products/${product.id}`">
        <h3 class="font-bold text-lg mb-1 group-hover:text-brand-600 transition line-clamp-2">
          {{ product.name }}
        </h3>
      </router-link>

      <p class="text-sm text-gray-500 mb-3" v-if="product.farmer_name">
        👨‍🌾 {{ product.farmer_name }}
      </p>

      <div class="mt-auto flex items-center justify-between">
        <div class="flex flex-col">
          <span class="font-bold text-lg text-brand-600">
            {{ formatRupiah(product.price) }}
          </span>
          <span class="text-xs text-gray-400" v-if="product.subscription_period">
             / {{ product.subscription_period }}
          </span>
        </div>

        <button
          @click.prevent="$emit('addToCart', product)"
          class="w-10 h-10 rounded-full bg-gray-100 dark:bg-gray-700 flex items-center justify-center text-gray-600 dark:text-gray-300 hover:bg-brand-600 hover:text-white transition shadow-sm"
        >
          <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="9" cy="21" r="1"/><circle cx="20" cy="21" r="1"/><path d="M1 1h4l2.68 13.39a2 2 0 0 0 2 1.61h9.72a2 2 0 0 0 2-1.61L23 6H6"/></svg>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { getImageUrl } from '@/utils/image'
import type { Product } from '@/dto/product/Product'

interface Props {
  product: Product
  label?: string
  labelColor?: string
}

const props = withDefaults(defineProps<Props>(), {
  labelColor: 'bg-green-500'
})

defineEmits(['addToCart'])

const imageError = ref(false)

const handleImageError = () => {
  console.warn('Failed to load image:', props.product.image_url)
  imageError.value = true
}

const formatRupiah = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}
</script>
