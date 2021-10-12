<template>
   <div id="app">
      <div class="header" role="banner" id="uva-header">
         <div class="library-link">
            <a target="_blank" href="https://library.virginia.edu">
               <UvaLibraryLogo />
            </a>
         </div>
         <div class="site-link">
           <router-link to="/">Collections Admin</router-link>
         </div>
      </div>
      <router-view />
   </div>
</template>

<script>
import UvaLibraryLogo from "@/components/UvaLibraryLogo"
export default {
   components: {
      UvaLibraryLogo
   }
}
</script>

<style lang="scss">
:root {
   /* OFFICIAL BRAND COLORS */
   --uvalib-brand-blue-lightest: #87b9d9;
   --uvalib-brand-blue-lighter: #3395d4;
   --uvalib-brand-blue-light: #0370b7;
   --uvalib-brand-blue: #232d4b;

   --uvalib-brand-orange-lightest: #ffead6;
   --uvalib-brand-orange: #e57200;
   --uvalib-brand-orange-dark: #b35900;

   --uvalib-blue-alt-light: #bfe7f7;
   --uvalib-blue-alt: #007bac;
   --uvalib-blue-alt-dark: #005679;
   --uvalib-blue-alt-darkest: #141e3c;

   --uvalib-teal-lightest: #c8f2f4;
   --uvalib-teal-light: #5bd7de;
   --uvalib-teal: #25cad3;
   --uvalib-teal-dark: #1da1a8;
   --uvalib-teal-darker: #16777c;

   --uvalib-green-lightest: #89cc74;
   --uvalib-green: #62bb46;
   --uvalib-green-dark: #4e9737;

   --uvalib-red-lightest: #fbcfda;
   --uvalib-red: #ef3f6b;
   --uvalib-red-emergency: #df1e43;
   --uvalib-red-darker: #b30000;
   --uvalib-red-dark: #df1e43;

   --uvalib-yellow-light: #fef6c8;
   --uvalib-yellow: #ecc602;
   --uvalib-yellow-dark: #b99c02;

   --uvalib-beige: #f7efe1;
   --uvalib-beige-dark: #c0b298;

   --uvalib-grey-lightest: #f1f1f1;
   --uvalib-grey-light: #dadada;
   --uvalib-grey: #808080;
   --uvalib-grey-dark: #565656;
   --uvalib-grey-darkest: #2b2b2b;

   --uvalib-text-light: #ffffff;
   --uvalib-text: var(--uvalib-grey-dark);
   --uvalib-text-dark: var(--uvalib-grey-darkest);

   --uvalib-accessibility-highlight: var(--uvalib-brand-blue-light);

   --color-primary-orange: var(--uvalib-brand-orange);
   --color-link: var(--uvalib-blue-alt-dark);
   --color-link-darker: var(--uvalib-blue-alt-dark);
   --color-primary-text: var(--uvalib-grey-dark);
   --color-error: var(--uvalib-red-emergency);
}
#app {
   font-family: "franklin-gothic-urw", arial, sans-serif;
   -webkit-font-smoothing: antialiased;
   -moz-osx-font-smoothing: grayscale;
   text-align: center;
   color: var(--color-primary-text);
   margin: 0;
   padding: 0;
   background: white;
   font-size: 16px;
   line-height: 1.3;

   h1 {
      color: var(--uvalib-brand-orange);
      margin: 25px 0;
      font-weight: bold;
      position: relative;
   }
   a {
      color: var(--color-link);
      font-weight: 500;
      text-decoration: none;
      &:hover {
         text-decoration: underline;
      }
   }
   input[type=text], select {
      padding: 5px;
      border-radius: 5px;
      border: 1px solid var(--uvalib-grey-light);
      box-sizing: border-box;
   }
}

body {
   margin: 0;
   padding: 0;
   background-color: var(--uvalib-blue-alt-darkest);
}

div.header {
   background-color: var(--uvalib-brand-blue);
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
