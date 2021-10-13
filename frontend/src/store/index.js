import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
   state: {
      working: false,
      collections: [],
      selectedID: -1,
      editing: false,
      details: {
         id: -1,
         title: "",
         description: "",
         itemLabel: "Issue",
         startDate: "",
         endDate: "",
         filter: "",
         features: [],
         images: [],
      },
      fatal: ""
   },
   mutations: {
      setCollections(state, data) {
         state.collections.splice(0, state.collections.length)
         data.forEach(c => state.collections.push(c))
      },
      setCollectionDetail(state, data) {
         state.details.images.splice(0, state.details.images.length)
         state.details.features.splice(0, state.details.features.length)
         data.images.forEach( i => state.details.images.push(i))
         data.features.forEach( f => state.details.features.push(f))
         state.details.id = data.id
         state.details.title = data.title
         state.details.description  = data.description
         state.details.itemLabel = data.item_label
         state.details.startDate = data.start_date
         state.details.endDate = data.end_date
         state.details.filter = data.filter_name
      },
      setEditing(state, flag) {
         state.editing = flag
      },
      setFatalError(state, err) {
         state.fatal = err
      },
      setSelectedCollectionID(state, id) {
         state.selectedID = id
      },
      setWorking(state, flag) {
         state.working = flag
      }
   },
   actions: {
      getCollections(ctx) {
         ctx.commit("setWorking", true)
         axios.get("/api/collections").then(resp => {
            ctx.commit("setCollections", resp.data)
            ctx.commit("setWorking", false)
         }).catch((e) => {
            ctx.commit("setFatalError", e)
            ctx.commit("setWorking", false)
         })
      },
      getCollectionDetail(ctx, id) {
         ctx.commit("setSelectedCollectionID", id)
         ctx.commit("setWorking", true)
         axios.get("/api/collections/"+id).then(resp => {
            ctx.commit("setCollectionDetail", resp.data)
            ctx.commit("setWorking", false)
         }).catch((e) => {
            ctx.commit("setFatalError", e)
            ctx.commit("setWorking", false)
         })
      }
   },
})
