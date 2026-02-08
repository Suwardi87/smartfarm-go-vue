import { createRouter, createWebHistory } from 'vue-router'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  scrollBehavior(to, from, savedPosition) {
    return savedPosition || { left: 0, top: 0 }
  },
  routes: [
    {
      path: '/',
      name: 'Marketplace',
      component: () => import('../views/Marketplace/Home.vue'),
      meta: {
        title: 'SmartFarm Marketplace',
      },
    },
    {
      path: '/marketplace',
      name: 'AllProducts',
      component: () => import('../views/Marketplace/Marketplace.vue'),
      meta: {
        title: 'Katalog Produk',
      },
    },
    {
      path: '/farmer/dashboard',
      name: 'FarmerDashboard',
      component: () => import('../views/Marketplace/FarmerDashboard.vue'),
      meta: {
        title: 'Dashboard Petani',
      },
    },
    {
      path: '/farmer/products',
      name: 'FarmerProducts',
      component: () => import('../views/Marketplace/FarmerProductList.vue'),
    },
    {
      path: '/products/create',
      name: 'CreateProduct',
      component: () => import('../views/Marketplace/CreateProduct.vue'),
      meta: {
        title: 'Jual Produk',
      },
    },
    {
      path: '/products/edit/:id',
      name: 'EditProduct',
      component: () => import('../views/Marketplace/EditProduct.vue'),
    },
    {
      path: '/cart',
      name: 'Cart',
      component: () => import('../views/Marketplace/Cart.vue'),
      meta: {
        title: 'Keranjang Belanja',
      },
    },
    {
      path: '/checkout',
      name: 'Checkout',
      component: () => import('../views/Marketplace/Checkout.vue'),
      meta: {
        title: 'Checkout',
        requiresAuth: true,
      },
    },
    {
      path: '/products/:id',
      name: 'ProductDetail',
      component: () => import('../views/Marketplace/ProductDetail.vue'),
      meta: {
        title: 'Detail Produk',
      },
    },
    {
      path: '/orders',
      name: 'Orders',
      component: () => import('../views/Marketplace/Orders.vue'),
      meta: {
        title: 'Riwayat Pesanan',
        requiresAuth: true,
      },
    },
    {
      path: '/profile',
      name: 'ProfileEdit',
      component: () => import('../views/Marketplace/Profile.vue'),
      meta: {
        title: 'Profil Saya',
        requiresAuth: true,
      },
    },
    {
      path: '/addresses',
      name: 'Addresses',
      component: () => import('../views/Marketplace/Addresses.vue'),
      meta: {
        title: 'Alamat Pengiriman',
        requiresAuth: true,
      },
    },
    {
      path: '/dashboard',
      name: 'Dashboard',
      component: () => import('../views/Ecommerce.vue'),
      meta: {
        title: 'Dashboard',
        requiresAuth: true,
      },
    },
    {
      path: '/calendar',
      name: 'Calendar',
      component: () => import('../views/Others/Calendar.vue'),
      meta: {
        title: 'Calendar',
        requiresAuth: true,
      },
    },
    {
      path: '/profile',
      name: 'Profile',
      component: () => import('../views/Others/UserProfile.vue'),
      meta: {
        title: 'Profile',
      },
    },
    {
      path: '/form-elements',
      name: 'Form Elements',
      component: () => import('../views/Forms/FormElements.vue'),
      meta: {
        title: 'Form Elements',
      },
    },
    {
      path: '/basic-tables',
      name: 'Basic Tables',
      component: () => import('../views/Tables/BasicTables.vue'),
      meta: {
        title: 'Basic Tables',
      },
    },
    {
      path: '/line-chart',
      name: 'Line Chart',
      component: () => import('../views/Chart/LineChart/LineChart.vue'),
    },
    {
      path: '/bar-chart',
      name: 'Bar Chart',
      component: () => import('../views/Chart/BarChart/BarChart.vue'),
    },
    {
      path: '/alerts',
      name: 'Alerts',
      component: () => import('../views/UiElements/Alerts.vue'),
      meta: {
        title: 'Alerts',
      },
    },
    {
      path: '/avatars',
      name: 'Avatars',
      component: () => import('../views/UiElements/Avatars.vue'),
      meta: {
        title: 'Avatars',
      },
    },
    {
      path: '/badge',
      name: 'Badge',
      component: () => import('../views/UiElements/Badges.vue'),
      meta: {
        title: 'Badge',
      },
    },

    {
      path: '/buttons',
      name: 'Buttons',
      component: () => import('../views/UiElements/Buttons.vue'),
      meta: {
        title: 'Buttons',
      },
    },

    {
      path: '/images',
      name: 'Images',
      component: () => import('../views/UiElements/Images.vue'),
      meta: {
        title: 'Images',
      },
    },
    {
      path: '/videos',
      name: 'Videos',
      component: () => import('../views/UiElements/Videos.vue'),
      meta: {
        title: 'Videos',
      },
    },
    {
      path: '/blank',
      name: 'Blank',
      component: () => import('../views/Pages/BlankPage.vue'),
      meta: {
        title: 'Blank',
      },
    },

    {
      path: '/error-404',
      name: '404 Error',
      component: () => import('../views/Errors/FourZeroFour.vue'),
      meta: {
        title: '404 Error',
      },
    },

    {
      path: '/signin',
      name: 'Signin',
      component: () => import('../views/Auth/Signin.vue'),
      meta: {
        title: 'Signin',
        guestOnly: true,
      },
    },
    {
      path: '/signup',
      name: 'Signup',
      component: () => import('../views/Auth/Signup.vue'),
      meta: {
        title: 'Signup',
      },
    },
  ],
})

export default router

import { useUser } from "@/stores/user"

router.beforeEach(async (to, from, next) => {
  document.title = `Vue.js ${to.meta.title} | TailAdmin - Vue.js Tailwind CSS Dashboard Template`

  const userStore = useUser()
  let isLoggedIn = false

  try {
    await userStore.fetchUser()
    isLoggedIn = userStore.state.isAuthenticated
  } catch {
    isLoggedIn = false
  }

  // ✅ sudah login tapi ke signin / signup
  if (to.meta.guestOnly && isLoggedIn) {
    return next('/dashboard')
  }

  // ❌ belum login tapi ke dashboard
  if (to.meta.requiresAuth && !isLoggedIn) {
    return next('/signin')
  }

  next()
})
