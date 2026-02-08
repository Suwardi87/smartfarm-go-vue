<template>
  <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 md:grid-cols-4 md:gap-6">
    <!-- Revenue -->
    <div class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] md:p-6">
      <div class="flex items-center justify-center w-12 h-12 bg-gray-100 rounded-xl dark:bg-gray-800 text-brand-500">
        <svg fill="currentColor" width="24" height="24" viewBox="0 0 24 24"><path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 14H11V15H12V13H11V12H13V15H14V16H13zM12 11C11.45 11 11 10.55 11 10s.45-1 1-1 1 .45 1 1-.45 1-1 1z"/></svg>
      </div>
      <div class="flex items-end justify-between mt-5">
        <div>
          <span class="text-sm text-gray-500 dark:text-gray-400">Total Revenue</span>
          <h4 class="mt-2 font-bold text-gray-800 text-title-sm dark:text-white/90">
            {{ formatRupiah(stats?.total_revenue || 0) }}
          </h4>
        </div>
        <span :class="[(stats?.revenue_growth ?? 0) >= 0 ? 'bg-success-50 text-success-600' : 'bg-error-50 text-error-600', 'flex items-center gap-1 rounded-full py-0.5 pl-2 pr-2.5 text-sm font-medium']">
          {{ (stats?.revenue_growth ?? 0) >= 0 ? '↑' : '↓' }} {{ Math.abs(stats?.revenue_growth ?? 0) }}%
        </span>
      </div>
    </div>

    <!-- Orders -->
    <div class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] md:p-6">
      <div class="flex items-center justify-center w-12 h-12 bg-gray-100 rounded-xl dark:bg-gray-800 text-brand-500">
        <svg fill="currentColor" width="24" height="24" viewBox="0 0 24 24"><path d="M7 18c-1.1 0-1.99.9-1.99 2S5.9 22 7 22s2-.9 2-2-.9-2-2-2zM1 2v2h2l3.6 7.59-1.35 2.45c-.16.28-.25.61-.25.96 0 1.1.9 2 2 2h12v-2H7.42c-.14 0-.25-.11-.25-.25l.03-.12.9-1.63h7.45c.75 0 1.41-.41 1.75-1.03l3.58-6.49c.08-.14.12-.31.12-.48 0-.55-.45-1-1-1H5.21l-.94-2H1zm16 16c-1.1 0-1.99.9-1.99 2s.89 2 1.99 2 2-.9 2-2-.9-2-2-2z"/></svg>
      </div>
      <div class="flex items-end justify-between mt-5">
        <div>
          <span class="text-sm text-gray-500 dark:text-gray-400">Total Orders</span>
          <h4 class="mt-2 font-bold text-gray-800 text-title-sm dark:text-white/90">{{ stats?.total_orders || 0 }}</h4>
        </div>
        <span :class="[(stats?.orders_growth ?? 0) >= 0 ? 'bg-success-50 text-success-600' : 'bg-error-50 text-error-600', 'flex items-center gap-1 rounded-full py-0.5 pl-2 pr-2.5 text-sm font-medium']">
          {{ (stats?.orders_growth ?? 0) >= 0 ? '↑' : '↓' }} {{ Math.abs(stats?.orders_growth ?? 0) }}%
        </span>
      </div>
    </div>

    <!-- Customers -->
    <div class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] md:p-6">
      <div class="flex items-center justify-center w-12 h-12 bg-gray-100 rounded-xl dark:bg-gray-800 text-brand-500">
        <svg fill="currentColor" width="24" height="24" viewBox="0 0 24 24"><path d="M16 11c1.66 0 2.99-1.34 2.99-3S17.66 5 16 5s-3 1.34-3 3 1.34 3 3 3zm-8 0c1.66 0 2.99-1.34 2.99-3S9.66 5 8 5 5 6.34 5 8s1.34 3 3 3zm0 2c-2.33 0-7 1.17-7 3.5V19h14v-2.5c0-2.33-4.67-3.5-7-3.5zm8 0c-.29 0-.62.02-.97.05 1.16.84 1.97 1.97 1.97 3.45V19h6v-2.5c0-2.33-4.67-3.5-7-3.5z"/></svg>
      </div>
      <div class="flex items-end justify-between mt-5">
        <div>
          <span class="text-sm text-gray-500 dark:text-gray-400">Customers</span>
          <h4 class="mt-2 font-bold text-gray-800 text-title-sm dark:text-white/90">{{ stats?.total_customers || 0 }}</h4>
        </div>
        <span :class="[(stats?.customers_growth ?? 0) >= 0 ? 'bg-success-50 text-success-600' : 'bg-error-50 text-error-600', 'flex items-center gap-1 rounded-full py-0.5 pl-2 pr-2.5 text-sm font-medium']">
          {{ (stats?.customers_growth ?? 0) >= 0 ? '↑' : '↓' }} {{ Math.abs(stats?.customers_growth ?? 0) }}%
        </span>
      </div>
    </div>

    <!-- Products -->
    <div class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] md:p-6">
      <div class="flex items-center justify-center w-12 h-12 bg-gray-100 rounded-xl dark:bg-gray-800 text-brand-500">
        <svg fill="currentColor" width="24" height="24" viewBox="0 0 24 24"><path d="M12 3L2 12h3v8h6v-6h2v6h6v-8h3L12 3z"/></svg>
      </div>
      <div class="flex items-end justify-between mt-5">
        <div>
          <span class="text-sm text-gray-500 dark:text-gray-400">Total Products</span>
          <h4 class="mt-2 font-bold text-gray-800 text-title-sm dark:text-white/90">{{ stats?.total_products || 0 }}</h4>
        </div>
        <span :class="[(stats?.products_growth ?? 0) >= 0 ? 'bg-success-50 text-success-600' : 'bg-error-50 text-error-600', 'flex items-center gap-1 rounded-full py-0.5 pl-2 pr-2.5 text-sm font-medium']">
          {{ (stats?.products_growth ?? 0) >= 0 ? '↑' : '↓' }} {{ Math.abs(stats?.products_growth ?? 0) }}%
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { FarmerStats } from '@/services/analyticsService'

defineProps<{
  stats: FarmerStats | null
}>()

const formatRupiah = (value: number) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0
  }).format(value)
}
</script>
