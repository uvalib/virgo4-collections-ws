<template>
   <div class="home">
      <h1>Virgo Collections Management</h1>
      <WaitSpinner v-if="working" message="Initializing system..." />
      <div v-else class="content">
         <div class="list-wrap">
            <h2>Collection</h2>
            <div class="list">
               <div class="item" v-for="c in collections" :key="c.id">{{c.title}}</div>
            </div>
         </div>
         <div class="detail-wrap">
             <h2>Detail</h2>
             <div class="details">
                <div v-if="selectedID == -1" class="hint">Select a collection from the list on the left to view and edit the collection details.</div>
                <template v-else>
                </template>
             </div>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from "vuex"
export default {
   name: 'Home',
   components: {
   },
   computed: {
      ...mapState({
         working: state => state.working,
         collections: state => state.collections,
         selectedID: state => state.selectedID
      })
   },
   created() {
      this.$store.dispatch("getCollections")
   }
}
</script>

<style lang="scss">
.home {
   min-height: 400px;
}
.content {
   display: flex;
   flex-flow: row nowrap;
   text-align: left;
   h2 {
      margin:5px 0 10px 0;
      color: var(--uvalib-text);
   }
   .list-wrap, .detail-wrap {
      padding: 10px 20px;
      margin-bottom: 20px;
   }
   .list {
      text-align: left;
      border: 1px solid var(--uvalib-grey-light);
      padding: 0;
      border-radius: 5px;
      max-height: 450px;
      max-width: 400px;
      overflow: scroll;
      .item {
         margin: 0;
         cursor: pointer;
         padding: 5px 15px;
         &:hover {
            background: var(--uvalib-teal-lightest);
         }
      }
   }
   .details {
      .hint {
         font-size: 1.25em;
         font-style: italic;
         color: var(--uvalib-grey);
         margin-left: 25px;
      }
   }
}
</style>
