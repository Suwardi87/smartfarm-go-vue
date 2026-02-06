import { reactive, computed, watch } from 'vue'
import type { Product } from '@/dto/product/Product'
import { normalizeProduct } from '@/utils/transformers'

export interface CartItem {
    product: Product
    quantity: number
}

const state = reactive({
    items: [] as CartItem[]
})

// Load from localStorage
const stored = localStorage.getItem('smartfarm_cart')
if (stored) {
    try {
        const loaded = JSON.parse(stored)
        // Normalize product data when loading from storage
        state.items = loaded.map((item: any) => ({
            product: normalizeProduct(item.product),
            quantity: item.quantity
        }))
    } catch (e) {
        console.error('Failed to load cart', e)
    }
}

// Watch and save
watch(() => state.items, (newItems) => {
    localStorage.setItem('smartfarm_cart', JSON.stringify(newItems))
}, { deep: true })

export const useCart = () => {
    const addItem = (product: Product, quantity = 1) => {
        // Normalize product to ensure consistent field naming
        const normalizedProduct = normalizeProduct(product)
        const existing = state.items.find(i => i.product.id === normalizedProduct.id)
        if (existing) {
            existing.quantity += quantity
        } else {
            state.items.push({ product: normalizedProduct, quantity })
        }
    }

    const removeItem = (productId: number) => {
        const index = state.items.findIndex(i => i.product.id === productId)
        if (index > -1) {
            state.items.splice(index, 1)
        }
    }

    const updateQuantity = (productId: number, quantity: number) => {
        const item = state.items.find(i => i.product.id === productId)
        if (item) {
            if (quantity <= 0) removeItem(productId)
            else item.quantity = quantity
        }
    }

    const clearCart = () => {
        state.items = []
    }

    const totalItems = computed(() => state.items.reduce((sum, item) => sum + item.quantity, 0))

    const totalPrice = computed(() => state.items.reduce((sum, item) => sum + (item.product.price * item.quantity), 0))

    return {
        state,
        addItem,
        removeItem,
        updateQuantity,
        clearCart,
        totalItems,
        totalPrice
    }
}
