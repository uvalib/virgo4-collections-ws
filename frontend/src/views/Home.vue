<template>
   <div class="home">
      <h1>Virgo Collections Management</h1>
      <wait-spinner v-if="working && selectedID == -1" message="Initializing system..."/>
      <div v-else class="content">
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
         </div>
         <div class="detail-wrap">
             <h2>Details</h2>
             <div class="details">
                <div v-if="selectedID == -1" class="hint">Select a collection from the list on the left to view and edit the collection details.</div>
                <template v-else>
                   <collection-detail v-if="!editing" />
                </template>
             </div>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from "vuex"
import CollectionDetail from '../components/CollectionDetail.vue'
export default {
   name: 'Home',
   components: {
      CollectionDetail
   },
   computed: {
      ...mapState({
         working: state => state.working,
         collections: state => state.collections,
         selectedID: state => state.selectedID,
         editing: state => state.editing,
      })
   },
   methods: {
      collectionClicked(id) {
         this.$store.dispatch("getCollectionDetail", id)
         // let ele = document.getElementById(`c${id}`)
         // ele.scrollIntoView()
         // console.log("scrolled")
      }
   },
   created() {
      this.$store.dispatch("getCollections")
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
         background: white;
         border: 0;
         border-bottom: 2px solid var(--uvalib-grey);
         border-radius: 0;
      }
    }
    .details {
       border-radius: 5px;
       min-height: 600px;
       max-height: 600px;
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
      max-height: 600px;
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
      .item.selected {
         background: var(--uvalib-brand-blue-light);
         color: white;
      }
   }
}
</style>
