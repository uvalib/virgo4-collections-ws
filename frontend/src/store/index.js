import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
   state: {
      working: false,
      collections: [],
      selectedID: -1,
      editing: false,
      adding: false,
      features: [],
      details: {
         id: -1,
         title: "",
         description: "",
         itemLabel: "Issue",
         startDate: "",
         endDate: "",
         filter: "",
         features: [],
         image: null,
      },
      fatal: ""
   },
   mutations: {
      setCollections(state, data) {
         state.collections.splice(0, state.collections.length)
         data.forEach(c => state.collections.push(c))
      },
      clearDetails(state) {
         state.selectedID = -1
         state.details.id = -1,
         state.details.title = ""
         state.details.description = ""
         state.details.itemLabel = "Issue"
         state.details.startDate = ""
         state.details.endDate = ""
         state.details.filter = ""
         state.details.features.splice(0, state.details.features.length)
         state.details.image = null
      },
      setCollectionDetail(state, data) {
         state.details.image = null
         if (data.image) {
            state.details.image = data.image
         }
         state.details.features.splice(0, state.details.features.length)
         data.features.forEach( f => state.details.features.push(f))
         state.details.id = data.id
         state.details.title = data.title
         state.details.description  = data.description
         state.details.itemLabel = data.item_label
         state.details.startDate = data.start_date
         state.details.endDate = data.end_date
         state.details.filter = data.filter_name
      },
      setAdding(state, flag) {
         state.adding = flag
      },
      setEditing(state, flag) {
         state.editing = flag
      },
      setFeatures(state, data) {
         state.features.splice(0, state.features.length)
         data.forEach( f => state.features.push(f))
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
      getFeatures(ctx) {
         axios.get("/api/features").then(resp => {
            ctx.commit("setFeatures", resp.data)
         }).catch((e) => {
            ctx.commit("setFatalError", e)
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
