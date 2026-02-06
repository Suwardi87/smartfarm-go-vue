import { reactive } from 'vue'
import { getMe, logout as apiLogout } from '@/services/authService'

interface User {
    id: number
    name: string
    email: string
    role: string
}

const state = reactive({
    user: null as User | null,
    isAuthenticated: false
})

export const useUser = () => {
    const fetchUser = async () => {
        try {
            const response = await getMe()
            if (response.data && response.data.data) {
                state.user = response.data.data
                state.isAuthenticated = true
                return state.user
            }
        } catch (e) {
            state.user = null
            state.isAuthenticated = false
            throw e
        }
    }

    const logout = async () => {
        try {
            await apiLogout()
        } finally {
            state.user = null
            state.isAuthenticated = false
        }
    }

    return {
        state,
        fetchUser,
        logout
    }
}
