<template>
   <dl>
      <dt>Title:</dt>
      <dd>

      </dd>
      <dt>Description:</dt>
      <dd>

      </dd>
      <dt>Item Label:</dt>
      <dd>

      </dd>
      <dt>Start Date:</dt>
      <dd>

      </dd>
      <dt>End Date:</dt>
      <dd>

      </dd>
      <dt>Facet Name:</dt>
      <dd>

      </dd>
      <dt>Features:</dt>
      <dd>{{details.features.join(", ")}}</dd>
      <template v-if="details.image">
         <dt>Logo Title:</dt>
         <dd>
            <template v-if="details.image.title">{{details.image.title}}</template>
            <span v-else class="na">N/A</span>
         </dd>
         <dt>Logo Alt Text:</dt>
         <dd>
            <template v-if="details.image.alt_text">{{details.image.alt_text}}</template>
            <span v-else class="na">N/A</span>
         </dd>
         <dt>Logo:</dt>
         <dd>
            <img class="thumb" :src="details.image.url"/>
         </dd>
      </template>
   </dl>
</template>

<script>
import { mapState } from "vuex"
export default {
   computed: {
      ...mapState({
         selectedID: state => state.selectedID,
         details: state => state.details,
      })
   },
   methods: {
   },
}
</script>

<style lang="scss" scoped>
dl {
   margin-left: 25px;
   display: inline-grid;
   grid-template-columns: max-content 2fr;
   grid-column-gap: 15px;
   dt {
      font-weight: bold;
      text-align: right;
      margin: 0 0 10px 0;
   }
   dd {
      margin: 0 0 10px 0;
      word-break: break-word;
      -webkit-hyphens: auto;
      -moz-hyphens: auto;
      hyphens: auto;
      .thumb {
         max-width: 200px;
         border:1px solid var(--uvalib-grey-light);
      }
   }
   .na {
      color: #aaa;
      font-style: italic;
   }
}
</style>
