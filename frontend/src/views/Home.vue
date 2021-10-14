<template>
   <div class="home">
      <h1>Collections Management</h1>
      <wait-spinner v-if="working" message="Initializing system..." :overlay="true"/>
      <div class="content">
         <div class="list-wrap">
            <h2>Collection</h2>
            <div class="list">
               <div class="item" v-for="c in collections" :key="c.id" :id="`c${c.id}`"
                  :class="{selected: c.id==selectedID}"
                  @click="collectionClicked(c.id)"
               >
                  {{c.title}}
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
                     <uva-button class="delete">Delete</uva-button>
                     <uva-button @click="editSelected">Edit</uva-button>
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
import CollectionDetail from '../components/CollectionDetail.vue'
import EditCollection from '../components/EditCollection.vue'
export default {
   name: 'Home',
   components: {
      CollectionDetail,EditCollection
   },
   computed: {
      ...mapState({
         working: state => state.working,
         collections: state => state.collections,
         selectedID: state => state.selectedID,
      }),
      ...mapGetters({
         isEditing: 'isEditing',
      })
   },
   methods: {
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
   .list {
      text-align: left;
      border: 1px solid #aaa;
      padding: 0;
      border-radius: 0 0 5px 5px;
      max-height: 570px;
      width: 100%;
      overflow: scroll;
      margin: 0;
      box-sizing: border-box;
      border-top: 0;
      box-shadow:0 2px 4px rgba(0, 0, 0, 0.10), 0 2px 4px rgba(0, 0, 0, 0.12);
      .item {
         margin: 0;
         cursor: pointer;
         padding: 5px 15px;
         border-bottom: 1px solid var(--uvalib-grey-lightest);
         &:hover {
            background: var(--uvalib-teal-lightest);
         }
      }
      .item.selected {
         background: var(--uvalib-brand-blue-light);
         color: white;
      }
   }
   .list-buttons {
      margin: 10px 0;
      text-align: right;
   }
}
</style>
