import { createApp } from 'vue'
import App from './App.vue'
import './assets/main.css'
import router from './router'
import { useUserStore } from './stores/user'
import 'vue-sonner/style.css'

// Create Vue app instance
const app = createApp(App)
const { checkUserStatus } = useUserStore()
checkUserStatus()

// Mount app to the DOM
app.use(router)
app.mount('#app')
