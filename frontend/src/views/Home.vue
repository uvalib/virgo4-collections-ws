<template>
   <div class="home">
      <h1>Collections Management</h1>
      <WaitSpinner v-if="store.working" message="Initializing system..." :overlay="true"/>
      <div class="content">
         <div class="list-wrap">
            <h2>Collection</h2>
            <div class="list-details">
               <div class="search-div">
                  <label>Find:</label>
                  <input type="text" v-model="query" @input="queryTyped" @keyup.enter.prevent.stop="querySelected"/>
               </div>
               <div class="list" id="collection-list">
                  <div class="item" v-for="c in store.collections" :key="c.id" :id="`c${c.id}`"
                     :class="{selected: c.id==store.selectedID}"
                     @click="collectionClicked(c.id)"
                  >
                     {{c.title}}
                  </div>
               </div>
            </div>
            <div class="list-buttons">
               <Button @click="addCollectionClicked" :disabled="store.isEditing" label="Add Collection"/>
            </div>
         </div>
         <div class="detail-wrap">
            <h2>
               <span>Details</span>
               <span class="detail-butons">
                  <template v-if="store.isEditing && store.selectedID > 0">
                     <Button severity="secondary" @click="cancelEdit" label="Cancel"/>
                     <Button @click="submitClicked" label="Submit"/>
                  </template>
                  <template v-else-if="store.isEditing && store.selectedID == 0">
                     <Button severity="secondary" @click="cancelEdit" label="Cancel"/>
                     <Button @click="submitClicked" label="Create"/>
                  </template>
                  <template v-else-if="store.selectedID > 0">
                     <Button @click="deleteClicked" severity="danger" label="Delete"/>
                     <Button @click="editSelected" label="Edit"/>
                  </template>
               </span>
            </h2>
             <div class="details">
               <div v-if="store.selectedID == 0 && !store.isEditing" class="hint">Select a collection from the list on the left to view and edit the collection details.</div>
               <collection-detail v-if="store.selectedID > 0 && !store.isEditing" />
               <edit-collection v-else-if="store.isEditing" />
             </div>
         </div>
      </div>
   </div>
</template>

<script setup>
import { useCollectionStore } from "@/stores/collection"
import CollectionDetail from "@/components/CollectionDetail.vue"
import EditCollection from "@/components/EditCollection.vue"
import WaitSpinner from "@/components/WaitSpinner.vue"
import { useConfirm } from "primevue/useconfirm"
import { onMounted, ref } from "vue"

const store = useCollectionStore()
const query = ref("")
const confirm = useConfirm()

function queryTyped() {
   let val = store.collections.find( c => c.title.toLowerCase().indexOf(query.value)==0)
   if (val) {
      let eles = document.getElementsByClassName("tgt-collection")
      for (let i = 0; i < eles.length; i++) {
         eles[i].classList.remove('tgt-collection')
      }
      let tgt = document.getElementById(`c${val.id}`)
      if (tgt) {
         scrollParentToChild(document.getElementById("collection-list"), tgt)
         tgt.classList.add("tgt-collection")
      }
   }
}
function querySelected() {
   let val = store.collections.find( c => c.title.toLowerCase().indexOf(query.value)==0)
   if ( val ) {
      collectionClicked(val.id)
   }
}
function scrollParentToChild(parent, child) {
   var parentRect = parent.getBoundingClientRect()
   var parentViewableArea = {
      height: parent.clientHeight,
      width: parent.clientWidth
   }
   var childRect = child.getBoundingClientRect()
   var isViewable = (childRect.top >= parentRect.top) && (childRect.bottom <= parentRect.top + parentViewableArea.height)
   if (!isViewable) {
      // Should we scroll using top or bottom? Find the smaller ABS adjustment
      const scrollTop = childRect.top - parentRect.top;
      const scrollBot = childRect.bottom - parentRect.bottom;
      if (Math.abs(scrollTop) < Math.abs(scrollBot)) {
         // we're near the top of the list
         parent.scrollTop += scrollTop;
      } else {
         // we're near the bottom of the list
         parent.scrollTop += scrollBot;
      }
   }
}
const deleteClicked = (() => {
   confirm.require({
      message: `Delete collection '${store.details.title}'? All data will be lost. This cannot be reveresed. `,
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
          store.deleteSelectedCollection()
      },
   })
})

function submitClicked() {
   store.setSubmit()
}
function collectionClicked(id) {
   store.clearDetails()
   store.setDisplay()
   store.getCollectionDetail(id)
}
function editSelected() {
   store.setEdit()
}
function cancelEdit() {
   store.setDisplay()
}
function addCollectionClicked() {
   store.clearDetails()
   store.setEdit()
}

onMounted(()=>{
   store.getCollections()
   store.getFeatures()
})
</script>

<style lang="scss" scoped>
.content {
   display: flex;
   flex-flow: row nowrap;
   text-align: left;

   .list-wrap, .detail-wrap {
      padding: 10px 20px;
      margin-bottom: 20px;
   }
   .list-wrap {
      flex-basis: 25%;
      h2 {
         padding: 5px 10px;
         margin:5px 0 0 0;
         background: $uva-blue-alt-300;
         border: 1px solid $uva-blue-alt;
         border-bottom: 2px solid $uva-blue-alt;
      }
   }
   .detail-butons {
      display: flex;
      flex-flow: row nowrap;
      gap: 0.5rem;
   }
    .detail-wrap {
       flex-basis: 75%;
       h2 {
         display: flex;
         flex-flow: row nowrap;
         justify-content: space-between;
         background: white;
         border: 0;
         border-bottom: 2px solid $uva-grey;
         border-radius: 0;
         padding: 5px 0px 5px 5px;
      }
    }
    .details {
       padding: 20px 20px 0 0;
       .hint {
         font-size: 1.25em;
      }
    }
    .list-details {
       border: 1px solid $uva-grey;
       .search-div {
          padding: 5px;
          display: flex;
          flex-flow: row nowrap;
          align-items: center;
          justify-content: flex-start;
          background: $uva-grey-200;
          label {
             font-weight: bold;
             margin: 0 5px;
          }
          input {
             box-sizing: border-box;
             flex-grow: 1;
          }
          border-bottom: 1px solid  $uva-grey;;
       }
      .list {
         text-align: left;
         padding: 0;
         max-height: 570px;
         width: 100%;
         overflow: scroll;
         margin: 0;
         box-sizing: border-box;
         border-top: 0;

         .item {
            margin: 0;
            cursor: pointer;
            padding: 5px 15px;
            &:hover {
               background: $uva-blue-alt-300;
            }
         }
         .item.selected {
            background: $uva-brand-blue-100;
            color: white;
         }
      }
    }
   .list-buttons {
      margin: 10px 0;
      text-align: right;
   }
}
</style>
