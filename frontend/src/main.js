import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import axios from 'axios'

// 设置axios的默认基础URL
// 上线部署需要改动axios.defaults.baseURL = 'http://你的服务器公网IP:8080'
axios.defaults.baseURL = 'http://localhost:8080'

createApp(App)
  .use(router)
  .mount('#app')