import { createApp, markRaw } from 'vue'
import App from './App.vue'
import router from './router'

import { createPinia } from 'pinia'
const pinia = createPinia()
pinia.use(({ store }) => {
   // all stores can access router with this.router
   store.router = markRaw (router)
})

const app = createApp(App)
app.use(pinia)
app.use(router)


import DropZone from 'dropzone-vue'
import 'dropzone-vue/dist/dropzone-vue.common.css'
app.use(DropZone)

import WaitSpinner from "@/components/WaitSpinner.vue"
app.component('WaitSpinner', WaitSpinner)

import UvaButton from "@/components/UvaButton.vue"
app.component('UvaButton', UvaButton)

import '@fortawesome/fontawesome-free/css/all.css'

app.mount('#app')
