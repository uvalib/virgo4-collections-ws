<template>
   <div id="app">
      <ConfirmDialog position="top" :closable="false"/>
      <Dialog v-model:visible="collection.showMessage" :modal="true" position="top"
         header="System Message" @hide="collection.clearMessage()" style="width:400px;">
         <div class="message-body" id="msgbody">{{ collection.message }}></div>
      </Dialog>
      <div class="header" role="banner" id="uva-header">
         <div class="library-link">
            <a target="_blank" href="https://library.virginia.edu">
               <uva-library-logo />
            </a>
         </div>
         <div class="site-link">
           <router-link to="/">Collections Management</router-link>
         </div>
      </div>
      <div v-if="collection.fatal" class="fatal-err">
         <h1>Internal System Error</h1>
         <div class="err-txt">
            <p>{{collection.fatal}}</p>
            <p>Sorry for the inconvenience! We are aware of the issue and are working to resolve it. Please check back later.</p>
         </div>
      </div>
      <template v-else>
         <router-view />
      </template>
   </div>
</template>

<script setup>
import UvaLibraryLogo from "@/components/UvaLibraryLogo.vue"
import Dialog from 'primevue/dialog'
import { useCollectionStore } from "@/stores/collection"
const collection = useCollectionStore()
</script>

<style lang="scss">
#app {
   font-family: "franklin-gothic-urw", arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   text-align: center;
   color: $uva-text-color-base;
   margin: 0;
   padding: 0;
   background: white;
   font-size: 16px;
   line-height: 1.3;

   h1 {
      color: $uva-brand-blue;
      margin: 25px 0;
      font-weight: bold;
      position: relative;
   }
   a {
      color: $uva-blue-alt-A;
      font-weight: 500;
      text-decoration: none;
      &:hover {
         text-decoration: underline;
      }
   }
   input[type=text], select {
      padding: .5vw .75vw;
      border-radius: 0.3rem;
      border: 1px solid $uva-grey-100;
      box-sizing: border-box;
   }
}

body {
   margin: 0;
   padding: 0;
   background-color: $uva-blue-alt-B;
}

.fatal-err {
   padding-top: 25px;
   min-height: 500px;
   .err-txt {
      font-size: 1.15em;
   }
}

div.header {
   background-color: $uva-brand-blue;
   color: white;
   padding: 1vw 20px;
   text-align: left;
   position: relative;
   box-sizing: border-box;
   display: flex;
   flex-direction: row;
   flex-wrap: nowrap;
   justify-content: space-between;
   align-content: stretch;
   align-items: center;
   div.library-link {
      height: 45px;
      width: 220px;
      order: 0;
      flex: 0 1 auto;
      align-self: flex-start;
   }
   div.site-link {
      order: 0;
      font-size: 1.5em;
      a {
         color: white !important;
         text-decoration: none;
         &:hover {
            text-decoration: underline;
         }
      }
   }
}
</style>
