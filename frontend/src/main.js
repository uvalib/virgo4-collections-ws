import { createApp, markRaw } from 'vue'
import App from './App.vue'
import router from './router'

// primevue setup
import PrimeVue from 'primevue/config'
import 'primeicons/primeicons.css'
import Collections from './assets/theme/collections'
import Button from 'primevue/button'
import ConfirmationService from 'primevue/confirmationservice'
import ConfirmDialog from 'primevue/confirmdialog'

import { createPinia } from 'pinia'
const pinia = createPinia()
pinia.use(({ store }) => {
   // all stores can access router with this.router
   store.router = markRaw (router)
})

const app = createApp(App)
app.use(pinia)
app.use(router)

app.use(PrimeVue, {
   theme: {
      preset: Collections,
      options: {
         prefix: 'p',
         darkModeSelector: '.dpg-dark'
      }
   }
})

app.component('Button', Button)
app.component("ConfirmDialog", ConfirmDialog)
app.use(ConfirmationService)


import DropZone from 'dropzone-vue'
import 'dropzone-vue/dist/dropzone-vue.common.css'
app.use(DropZone)

app.mount('#app')
