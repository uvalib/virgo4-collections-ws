<template>
   <div class="home">
      <h1>Collections Management</h1>
      <wait-spinner v-if="working" message="Initializing system..." :overlay="true"/>
      <div class="content">
         <div class="list-wrap">
            <h2>Collection</h2>
            <div class="list-details">
               <div class="search-div">
                  <label>Find:</label>
                  <input type="text" v-model="query" @input="queryTyped" @keyup.enter.prevent.stop="querySelected"/>
               </div>
               <div class="list" id="collection-list">
                  <div class="item" v-for="c in collections" :key="c.id" :id="`c${c.id}`"
                     :class="{selected: c.id==selectedID}"
                     @click="collectionClicked(c.id)"
                  >
                     {{c.title}}
                  </div>
               </div>
            </div>
            <div class="list-buttons">
               <uva-button @click="addCollectionClicked" :class="{disabled: isEditing}">Add Collection</uva-button>
            </div>
         </div>
         <div class="detail-wrap">
            <h2>
               <span>Details</span>
               <span class="detail-butons">
                  <template v-if="isEditing && selectedID > 0">
                     <uva-button @click="cancelEdit" class="cancel">Cancel</uva-button>
                     <uva-button @click="submitClicked">Submit</uva-button>
                  </template>
                  <template v-else-if="isEditing && selectedID == 0">
                     <uva-button @click="cancelEdit" class="cancel">Cancel</uva-button>
                     <uva-button @click="submitClicked">Create</uva-button>
                  </template>
                  <template v-else-if="selectedID > 0">
                     <confirm @confirm="deleteCollection"
                        buttonText="Delete"
                        :message="`Delete collection <b>'${selected.title}</b>'?<br/>All data will be removed. This is not reversable.`" />
                     <uva-button class="pad-left" @click="editSelected">Edit</uva-button>
                  </template>
               </span>
            </h2>
             <div class="details">
               <div v-if="selectedID == 0 && !isEditing" class="hint">Select a collection from the list on the left to view and edit the collection details.</div>
               <collection-detail v-if="selectedID > 0 && !isEditing" />
               <edit-collection v-else-if="isEditing" />
             </div>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState, mapGetters } from "vuex"
import CollectionDetail from "@/components/CollectionDetail.vue"
import EditCollection from "@/components/EditCollection.vue"
import Confirm from "@/components/Confirm.vue"
export default {
   name: 'Home',
   components: {
      CollectionDetail,EditCollection,Confirm
   },
   data() {
      return {
         query: ""
      }
   },
   computed: {
      ...mapState({
         working: state => state.working,
         collections: state => state.collections,
         selectedID: state => state.selectedID,
         selected: state => state.details
      }),
      ...mapGetters({
         isEditing: 'isEditing',
      })
   },
   methods: {
      queryTyped() {
         let val = this.collections.find( c => c.title.toLowerCase().indexOf(this.query)==0)
         if (val) {
            let eles = document.getElementsByClassName("tgt-collection")
            for (let i = 0; i < eles.length; i++) {
               eles[i].classList.remove('tgt-collection')
            }
            let tgt = document.getElementById(`c${val.id}`)
            if (tgt) {
               // tgt.scrollIntoView()
               this.scrollParentToChild(document.getElementById("collection-list"), tgt)
               tgt.classList.add("tgt-collection")
            }
         }
      },
      querySelected() {
         let val = this.collections.find( c => c.title.toLowerCase().indexOf(this.query)==0)
         if ( val ) {
            this.collectionClicked(val.id)
         }
      },
      scrollParentToChild(parent, child) {
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
      },
      deleteCollection() {
         this.$store.dispatch("deleteSelectedCollection")
      },
      deleteClicked() {
         this.$store.commit("setShowConfirm", true)
      },
      submitClicked() {
         this.$store.commit("setSubmit")
      },
      collectionClicked(id) {
         this.$store.commit("clearDetails")
         this.$store.commit("setDisplay")
         this.$store.dispatch("getCollectionDetail", id)
      },
      editSelected() {
        this.$store.commit("setEdit")
      },
      cancelEdit() {
         this.$store.commit("setDisplay")
      },
      addCollectionClicked() {
         this.$store.commit("clearDetails")
        this.$store.commit("setEdit")
      }
   },
   created() {
      this.$store.dispatch("getCollections")
      this.$store.dispatch("getFeatures")
   }
}
</script>

<style lang="scss" scoped>
.home {
   min-height: 400px;
}
.content {
   display: flex;
   flex-flow: row nowrap;
   text-align: left;
   h2 {
      margin:5px 0 0 0;
      color: var(--uvalib-text);
      background: var(--uvalib-grey-lightest);
      padding: 5px 10px;
      border: 1px solid var(--uvalib-grey-light);
      border-radius:5px 5px 0 0;
   }
   .pad-left {
      margin-left: 5px;
   }
   .list-wrap, .detail-wrap {
      padding: 10px 20px;
      margin-bottom: 20px;
   }
   .list-wrap {
      flex-basis: 25%;
      h2 {
         background: var(--uvalib-blue-alt-light);
         border: 1px solid var(--uvalib-blue-alt);
         border-bottom: 2px solid var(--uvalib-blue-alt);
      }
   }
    .detail-wrap {
       flex-basis: 75%;
       h2 {
         display: flex;
         flex-flow: row nowrap;
         justify-content: space-between;
         background: white;
         border: 0;
         border-bottom: 2px solid var(--uvalib-grey);
         border-radius: 0;
         padding: 5px 0px 5px 5px;
         .delete  {
            margin-right: 5px;
            background: #a00;
            border-color: #800;
            &:hover {
               background: #c00;
            }
         }
         .cancel {
            margin-right: 5px;
            background-color: var(--uvalib-grey-lightest);
            border: 1px solid var(--uvalib-grey);
            color: black;
            &:hover {
               background-color: var(--uvalib-grey-light);
            }
         }
      }
    }
    .details {
       border-radius: 5px;
       min-height: 600px;
       padding: 20px 20px 0 0;
       .hint {
         font-size: 1.25em;
         font-style: italic;
         color: var(--uvalib-grey);
         margin: 25px;
      }
    }
    .list-details {
       border-radius: 0 0 5px 5px;
       box-shadow:0 2px 4px rgba(0, 0, 0, 0.10), 0 2px 4px rgba(0, 0, 0, 0.12);
       border: 1px solid #aaa;
       .search-div {
          padding: 5px;
          display: flex;
          flex-flow: row nowrap;
          align-items: center;
          justify-content: flex-start;
          background: #f0f0f0;
          label {
             font-weight: bold;
             margin: 0 5px;
          }
          input {
             box-sizing: border-box;
             flex-grow: 1;
          }
          border-bottom: 1px solid #aaa;
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
            border-bottom: 1px solid var(--uvalib-grey-lightest);
            &:hover {
               background: var(--uvalib-teal-lightest);
            }
         }
         .tgt-collection {
            background: var(--uvalib-teal-lightest);
         }
         .item.selected {
            background: var(--uvalib-brand-blue-light);
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
