<template>
  <AdminLayout>
    <div class="space-y-6">
      <!-- Header -->
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-800 dark:text-white/90">Dashboard Petani</h1>
          <p class="text-sm text-gray-500 dark:text-gray-400 text-theme-sm">Selamat datang kembali! Berikut ringkasan performa pertanian Anda.</p>
        </div>
        <router-link
          to="/products/create"
          class="inline-flex items-center justify-center gap-2 rounded-lg bg-brand-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-brand-700 transition-colors"
        >
          <svg fill="none" stroke="currentColor" viewBox="0 0 24 24" class="w-5 h-5"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          Tambah Produk
        </router-link>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex items-center justify-center py-20">
        <div class="h-10 w-10 animate-spin rounded-full border-4 border-gray-200 border-t-brand-600"></div>
      </div>

      <!-- Dashboard Content -->
      <div v-else class="space-y-6">
        <!-- Metrics -->
        <EcommerceMetrics :stats="dashboardData?.stats || null" />

        <div class="grid grid-cols-12 gap-4 md:gap-6">
          <!-- Main Chart -->
          <div class="col-span-12 xl:col-span-8">
            <StatisticsChart />
          </div>

          <!-- Secondary Widget (Smart Prediction) -->
          <div class="col-span-12 xl:col-span-4">
             <div class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] h-full">
                <h3 class="text-lg font-semibold text-gray-800 dark:text-white/90 mb-4">Smart Prediction</h3>
                <p class="text-sm text-gray-500 mb-6">Komoditas dengan tren pencarian tertinggi minggu ini.</p>
                <div class="space-y-4">
                   <div v-for="item in trendingProducts" :key="item.id" class="flex items-center justify-between p-3 rounded-xl bg-gray-50 dark:bg-gray-800/50">
                      <span class="font-medium text-gray-700 dark:text-gray-300">{{ item.name }}</span>
                      <span class="text-brand-600 font-bold">{{ item.views }} views</span>
                   </div>
                </div>
             </div>
          </div>

          <!-- Recent Orders -->
          <div class="col-span-12">
            <RecentOrders :orders="dashboardData?.recent_orders || []" />
          </div>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import EcommerceMetrics from '@/components/ecommerce/EcommerceMetrics.vue'
import StatisticsChart from '@/components/ecommerce/StatisticsChart.vue'
import RecentOrders from '@/components/ecommerce/RecentOrders.vue'
import { getFarmerDashboardData, getRecommendations } from '@/services/analyticsService'
import type { FarmerDashboardData, CommodityTrend } from '@/services/analyticsService'

const loading = ref(true)
const dashboardData = ref<FarmerDashboardData | null>(null)
const trendingProducts = ref<CommodityTrend[]>([])

onMounted(async () => {
  try {
    const [dashRes, trendRes] = await Promise.all([
      getFarmerDashboardData(),
      getRecommendations()
    ])
    
    dashboardData.value = dashRes.data.data
    trendingProducts.value = trendRes.data.data.slice(0, 5) // Top 5
  } catch (error) {
    console.error('Failed to fetch dashboard data', error)
  } finally {
    loading.value = false
  }
})
</script>

