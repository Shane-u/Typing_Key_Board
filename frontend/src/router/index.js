import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Materials from '../views/Materials.vue'
import History from '../views/History.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/materials',
    name: 'Materials',
    component: Materials
  },
  {
    path: '/history',
    name: 'History',
    component: History
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router 