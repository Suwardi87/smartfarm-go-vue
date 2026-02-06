<template>
  <div v-if="isLoading" class="flex items-center justify-center py-12">
    <div class="flex flex-col items-center gap-4">
      <div class="w-12 h-12 border-4 border-gray-200 border-t-brand-600 rounded-full animate-spin"></div>
      <p class="text-gray-500 dark:text-gray-400">{{ message }}</p>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  isLoading: boolean
  message?: string
}

withDefaults(defineProps<Props>(), {
  message: 'Memuat...'
})
</script>
