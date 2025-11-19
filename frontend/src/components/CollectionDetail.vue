<template>
   <dl>
      <dt>Active:</dt>
      <dd>
         <span v-if="collection.details.active">Yes</span>
         <span v-else>No</span>
      </dd>
      <dt>Title:</dt>
      <dd>
         <template v-if="collection.details.title">{{collection.details.title}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Description:</dt>
      <dd>
         <span v-if="collection.details.description" v-html="collection.details.description"></span>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Item Label:</dt>
      <dd>
         <template v-if="collection.details.itemLabel">{{collection.details.itemLabel}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Start Date:</dt>
      <dd>
         <template v-if="collection.details.startDate">{{collection.details.startDate}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>End Date:</dt>
      <dd>
         <template v-if="collection.details.endDate">{{collection.details.endDate}}</template>
         <span v-else class="na">N/A</span>
      </dd>
      <dt>Facet Name:</dt>
      <dd>
         {{collection.details.filter}}
      </dd>
      <dt>Features:</dt>
      <dd>
         <template v-for="f in collection.details.features" :key="`f${f.id}`"><span class="feature">{{f.name}}</span></template>
      </dd>
      <template v-if="collection.details.image">
         <dt>Logo Title:</dt>
         <dd>
            <template v-if="collection.details.image.title">{{collection.details.image.title}}</template>
            <span v-else class="na">N/A</span>
         </dd>
         <dt>Logo Alt Text:</dt>
         <dd>
            <template v-if="collection.details.image.alt_text">{{collection.details.image.alt_text}}</template>
            <span v-else class="na">N/A</span>
         </dd>
         <dt>Logo:</dt>
         <dd>
            <img class="thumb" :src="collection.details.image.url"/>
         </dd>
      </template>
   </dl>
   <div class="controls">
      <Button @click="deleteClicked" severity="danger" label="Delete"/>
      <Button @click="collection.setEdit()" label="Edit"/>
   </div>
</template>

<script setup>
import { useCollectionStore } from "@/stores/collection"
import { useConfirm } from "primevue/useconfirm"

const collection = useCollectionStore()
const confirm = useConfirm()

const deleteClicked = (() => {
   confirm.require({
      message: `Delete collection '${collection.details.title}'? All data will be lost. This cannot be reveresed. `,
      header: 'Confirm Collection Delete',
      icon: 'pi pi-question-circle',
      rejectProps: {
         label: 'Cancel',
         severity: 'secondary'
      },
      acceptProps: {
         label: 'Delete'
      },
      accept: () => {
         collection.deleteSelectedCollection()
      },
   })
})


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
         border:1px solid $uva-grey-100;
      }
      .feature {
         display: inline-block;
         margin-right: 10px;
         border:1px solid $uva-grey-100;
         background: $uva-grey-200;
         padding: 2px 15px;
         border-radius: 15px;
      }
   }
   .na {
      color: #aaa;
      font-style: italic;
   }
}
.controls {
   padding-top: 20px;
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   gap: 10px;
   border-top: 2px solid $uva-grey
}
</style>
