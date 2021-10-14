import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import store from './store'

const app = createApp(App)
app.use(store)
app.use(router)

store.router = router

import DropZone from 'dropzone-vue'
app.use(DropZone)

import WaitSpinner from "@/components/WaitSpinner"
app.component('WaitSpinner', WaitSpinner)

import UvaButton from "@/components/UvaButton"
app.component('UvaButton', UvaButton)

import '@fortawesome/fontawesome-free/css/all.css'

app.mount('#app')
