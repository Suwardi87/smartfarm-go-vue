const API_URL = import.meta.env.VITE_API_URL

/**
 * Get full image URL from filename or absolute URL.
 * Handles:
 * 1. Absolute URLs (starting with http)
 * 2. Relative paths (prepends API_URL/uploads/)
 */
export const getImageUrl = (path: string | null | undefined): string => {
    if (!path) return ''

    if (path.startsWith('http')) {
        return path
    }

    // Clean path to avoid double slashes if needed, though typically not an issue with simple concat
    const cleanPath = path.startsWith('/') ? path.substring(1) : path
    return `${API_URL}/uploads/${cleanPath}`
}
