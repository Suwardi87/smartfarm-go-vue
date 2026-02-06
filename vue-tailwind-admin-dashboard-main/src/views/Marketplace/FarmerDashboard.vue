<template>
  <MarketplaceLayout>
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-10">
      <h1 class="text-3xl font-bold text-gray-900 dark:text-white mb-8">Dashboard Petani</h1>

      <!-- Smart Prediction Section -->
      <div class="bg-white dark:bg-gray-800 p-6 rounded-xl border border-gray-200 dark:border-gray-700 shadow-sm mb-8">
        <h2 class="text-xl font-bold mb-4 flex items-center gap-2">
           <span>📊</span> Smart Prediction: Rekomendasi Tanam
        </h2>
        <p class="text-gray-500 mb-6">Berdasarkan data pencarian dan pre-order pembeli, komoditas berikut memiliki permintaan tertinggi saat ini.</p>

        <div v-if="loading" class="flex justify-center py-10">
           <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-brand-600"></div>
        </div>
        <div v-else>
           <apexchart type="bar" height="350" :options="chartOptions" :series="series"></apexchart>
        </div>
      </div>

      <!-- Action -->
      <div class="flex justify-end">
        <router-link to="/products/create" class="px-6 py-3 bg-brand-600 text-white rounded-lg font-bold hover:bg-brand-700 transition">
          + Tambah Produk Baru
        </router-link>
      </div>

    </div>
  </MarketplaceLayout>
</template>

<script setup lang="ts">
import MarketplaceLayout from '@/components/layout/MarketplaceLayout.vue'
import { ref, onMounted } from 'vue'
import { getRecommendations } from '@/services/analyticsService'
import type { CommodityTrend } from '@/services/analyticsService'

const loading = ref(true)

const series = ref([
  {
    name: 'Skor Permintaan',
    data: [] as number[]
  }
])

const chartOptions = ref({
  chart: {
    type: 'bar',
    height: 350
  },
  plotOptions: {
    bar: {
      horizontal: true
    }
  },
  dataLabels: {
    enabled: false
  },
  xaxis: {
    categories: [] as string[]
  },
  colors: ['#10B981']
})

onMounted(async () => {
  try {
    const response = await getRecommendations()

    // 🔑 AMBIL ARRAY YANG BENAR
    const rawData = response.data.data ?? []
    
    // 🔑 MAP KE FORMAT YANG DIINGINKAN CHART
    const data = rawData.map(item => ({
      name: item.name,
      score: item.views ?? 0
    }))

    // 🔑 SORT BERDASARKAN SCORE
    data.sort((a, b) => b.score - a.score)

    // 🔑 UPDATE SERIES DAN OPTIONS BERSAMAAN
    series.value = [{
      name: 'Skor Permintaan',
      data: data.map(item => item.score)
    }]

    chartOptions.value = {
      ...chartOptions.value,
      xaxis: {
        categories: data.map(item => item.name)
      }
    }

  } catch (error) {
    console.error('Failed to fetch analytics', error)
  } finally {
    loading.value = false
  }
})
</script>

