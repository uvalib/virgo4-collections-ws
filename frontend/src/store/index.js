import { createStore } from 'vuex'
import axios from 'axios'

export default createStore({
   state: {
      working: false,
      collections: [],
      selectedID: 0,
      mode: "display",
      features: [],
      details: {
         id: 0,
         title: "",
         description: "",
         itemLabel: "Issue",
         startDate: "",
         endDate: "",
         filter: "",
         features: [],
         image: null,
      },
      fatal: "",
      message: "",
      logos: []
   },
   getters: {
      isEditing: state => {
         return state.mode == "edit" || state.mode == "submit"
      },
      isAdding: state => {
         return state.mode == "add" || state.mode == "submit"
      },
   },
   mutations: {
      setCollections(state, data) {
         state.collections.splice(0, state.collections.length)
         data.forEach(c => state.collections.push(c))
      },
      setLogos(state, data) {
         state.logos.splice(0, state.logos.length)
         data.forEach(c => state.logos.push(c))
      },
      clearDetails(state) {
         state.selectedID = 0
         state.details.id = 0,
         state.details.title = ""
         state.details.description = ""
         state.details.itemLabel = "Issue"
         state.details.startDate = ""
         state.details.endDate = ""
         state.details.filter = ""
         state.details.features.splice(0, state.details.features.length)
         state.details.image = null
      },
      deleteSelectedCollection(state) {
         let idx = state.collections.findIndex( c => c.id == state.selectedID)
         if (idx > -1) {
            state.collections.splice(idx, 1)
         }
      },
      addCollection(state, data) {
         state.collections.push({id: data.id, title: data.title})
         state.details.features.splice(0, state.details.features.length)
         data.features.forEach( f => state.details.features.push(f))
         state.details.id = data.id
         state.details.title = data.title
         state.details.description  = data.description
         state.details.itemLabel = data.item_label
         state.details.startDate = data.start_date
         state.details.endDate = data.end_date
         state.details.filter = data.filter_name
         state.selectedID = data.id
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
      setEdit(state) {
         state.mode = "edit"
      },
      setSubmit(state) {
         state.mode = "submit"
      },
      setDisplay(state) {
         state.mode = "display"
      },
      setFeatures(state, data) {
         state.features.splice(0, state.features.length)
         data.forEach( f => state.features.push(f))
      },
      setFatalError(state, err) {
         state.fatal = err
      },
      setMessage(state, err) {
         state.message = err
      },
      clearMessage(state) {
         state.message = ""
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
      getLogos(ctx) {
         axios.get("/api/logos").then(resp => {
            ctx.commit("setLogos", resp.data)
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
      },
      deleteSelectedCollection(ctx) {
         ctx.commit("setWorking", true)
         ctx.commit("clearMessage")
         axios.delete("/api/collections/"+ctx.state.selectedID).then( _resp => {
            ctx.commit("setWorking", false)
            ctx.commit("deleteSelectedCollection")
            ctx.commit("clearDetails")
         }).catch((e) => {
            ctx.commit("setMessage", e)
            ctx.commit("setWorking", false)
         })
      },
      deletePendingImage(ctx, filename) {
         axios.delete("/api/collections/"+ctx.state.selectedID+"/logo/"+filename)
      },
      submitCollection(ctx, collection) {
         ctx.commit("setWorking", true)
         ctx.commit("clearMessage")
         let addingNew = (collection.id == 0)
         axios.post("/api/collections", collection).then(resp => {
            if ( addingNew ) {
               ctx.commit("addCollection", resp.data)
            } else {
               ctx.commit("setCollectionDetail", resp.data)
            }
            ctx.commit("setWorking", false)
            ctx.commit("setDisplay")
         }).catch((e) => {
            ctx.commit("setEdit")
            ctx.commit("setMessage", e)
            ctx.commit("setWorking", false)
         })
      }
   },
})
