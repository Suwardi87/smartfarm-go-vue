<template>
  <div
    class="overflow-hidden rounded-2xl border border-gray-200 bg-white px-4 pb-3 pt-4 dark:border-gray-800 dark:bg-white/[0.03] sm:px-6"
  >
    <div class="flex flex-col gap-2 mb-4 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h3 class="text-lg font-semibold text-gray-800 dark:text-white/90">Recent Orders</h3>
      </div>
      <router-link
        to="/orders"
        class="inline-flex items-center gap-2 rounded-lg border border-gray-300 bg-white px-4 py-2 text-theme-sm font-medium text-gray-700 shadow-theme-xs hover:bg-gray-50 hover:text-gray-800 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-400 dark:hover:bg-white/[0.03] dark:hover:text-gray-200"
      >
        See all
      </router-link>
    </div>

    <div class="max-w-full overflow-x-auto custom-scrollbar">
      <table class="min-w-full">
        <thead>
          <tr class="border-t border-gray-100 dark:border-gray-800">
            <th class="py-3 text-left">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Products</p>
            </th>
            <th class="py-3 text-left">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Category</p>
            </th>
            <th class="py-3 text-left">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Price</p>
            </th>
            <th class="py-3 text-left">
              <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">Status</p>
            </th>
          </tr>
        </thead>
        <tbody>
          <tr
            v-for="(order, index) in orders"
            :key="index"
            class="border-t border-gray-100 dark:border-gray-800"
          >
            <td class="py-3 whitespace-nowrap">
              <div class="flex items-center gap-3">
                <div class="h-[50px] w-[50px] overflow-hidden rounded-md bg-gray-100">
                   <img :src="getImageUrl(order.image_url)" :alt="order.product_name" class="h-full w-full object-cover" />
                </div>
                <div>
                  <p class="font-medium text-gray-800 text-theme-sm dark:text-white/90">
                    {{ order.product_name }}
                  </p>
                  <span class="text-gray-500 text-theme-xs dark:text-gray-400">Order #{{ order.id }}</span>
                </div>
              </div>
            </td>
            <td class="py-3 whitespace-nowrap">
              <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ order.category }}</p>
            </td>
            <td class="py-3 whitespace-nowrap">
              <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ formatRupiah(order.price) }}</p>
            </td>
            <td class="py-3 whitespace-nowrap">
              <span
                :class="{
                  'rounded-full px-2 py-0.5 text-theme-xs font-medium': true,
                  'bg-success-50 text-success-600 dark:bg-success-500/15 dark:text-success-500':
                    order.status.toLowerCase() === 'paid' || order.status.toLowerCase() === 'delivered',
                  'bg-warning-50 text-warning-600 dark:bg-warning-500/15 dark:text-orange-400':
                    order.status.toLowerCase() === 'pending',
                }"
              >
                {{ order.status }}
              </span>
            </td>
          </tr>
          <tr v-if="orders.length === 0">
            <td colspan="4" class="py-10 text-center text-gray-500">Belum ada pesanan terbaru.</td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FarmerRecentOrder } from '@/services/analyticsService'
import { getImageUrl } from '@/utils/image'

defineProps<{
  orders: FarmerRecentOrder[]
}>()

const formatRupiah = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}
</script>
