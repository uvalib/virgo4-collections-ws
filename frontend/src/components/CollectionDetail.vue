<template>
   <dl>
      <dt>Title:</dt>
      <dd>
         <template v-if="details.title">{{details.title}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Description:</dt>
      <dd>
         <span v-if="details.description" v-html="details.description"></span>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Item Label:</dt>
      <dd>
         <template v-if="details.itemLabel">{{details.itemLabel}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Start Date:</dt>
      <dd>
         <template v-if="details.startDate">{{details.startDate}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>End Date:</dt>
      <dd>
         <template v-if="details.endDate">{{details.endDate}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Facet Name:</dt>
      <dd>
         {{details.filter}}
      </dd>
      <dt>Features:</dt>
      <dd>{{details.features.join(", ")}}</dd>
      <template v-if="details.image">
         <dt>Logo:</dt>
         <dd>
            <img class="thumb" :src="details.image.url"/>
         </dd>
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
      </template>
   </dl>
</template>

itemLabel: "Issue",
         startDate: "",
         endDate: "",
         filter: "",
         features: [],
         images: []

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
