import { defineStore } from 'pinia'
import axios from 'axios'

export const useCollectionStore = defineStore('collection', {
	state: () => ({
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
         active: false
      },
      fatal: "",
      message: "",
      logos: []
   }),
   getters: {
      isEditing: state => {
         return state.mode == "edit" || state.mode == "submit"
      },
      isAdding: state => {
         return state.mode == "add" || state.mode == "submit"
      },
   },
   actions: {
      clearDetails() {
         this.selectedID = 0
         this.details.id = 0,
         this.details.title = ""
         this.details.description = ""
         this.details.itemLabel = "Issue"
         this.details.startDate = ""
         this.details.endDate = ""
         this.details.filter = ""
         this.details.features.splice(0, this.details.features.length)
         this.details.image = null
         this.details.active = false
      },
      addCollection(data) {
         this.collections.push({id: data.id, title: data.title})
         this.details.features.splice(0, this.details.features.length)
         data.features.forEach( f => this.details.features.push(f))
         this.details.id = data.id
         this.details.active = data.active
         this.details.title = data.title
         this.details.description  = data.description
         this.details.itemLabel = data.item_label
         this.details.startDate = data.start_date
         this.details.endDate = data.end_date
         this.details.filter = data.filter_name
         if (data.image) {
            this.details.image = data.image
         }
         this.selectedID = data.id
      },
      setCollectionDetail(data) {
         this.details.image = null
         if (data.image) {
            this.details.image = data.image
         }
         this.details.features.splice(0, this.details.features.length)
         data.features.forEach( f => this.details.features.push(f))
         this.details.id = data.id
         this.details.active = data.active
         this.details.title = data.title
         this.details.description  = data.description
         this.details.itemLabel = data.item_label
         this.details.startDate = data.start_date
         this.details.endDate = data.end_date
         this.details.filter = data.filter_name
      },
      setEdit() {
         this.mode = "edit"
      },
      setSubmit() {
         this.mode = "submit"
      },
      setDisplay() {
         this.mode = "display"
      },

      getCollections() {
         this.working = true
         axios.get("/api/collections?all=true").then(resp => {
            this.collections =  resp.data
            this.working = false
         }).catch((e) => {
            this.fatal = e
            this.working = false
         })
      },
      getFeatures() {
         axios.get("/api/features").then(resp => {
            this.features =  resp.data
         }).catch((e) => {
            this.fatal = e
         })
      },
      getLogos() {
         axios.get("/api/logos").then(resp => {
            this.logos = resp.data
         }).catch((e) => {
            this.fatal = e
         })
      },
      getCollectionDetail(id) {
         this.selectedID = id
         this.working = true
         axios.get("/api/collections/"+id).then(resp => {
            this.setCollectionDetail(resp.data)
            this.working = false
         }).catch((e) => {
            this.fatal = e
            this.working = false
         })
      },
      deleteSelectedCollection() {
         this.working = true
         this.message = ""
         axios.delete("/api/collections/"+this.selectedID).then( _resp => {
            this.working = false
            let idx = this.collections.findIndex( c => c.id == this.selectedID)
            if (idx > -1) {
               this.collections.splice(idx, 1)
            }
            this.clearDetails()
         }).catch((e) => {
            this.message = e
            this.working = false
         })
      },
      deletePendingImage(filename) {
         axios.delete("/api/collections/"+this.selectedID+"/logo/"+filename)
      },
      submitCollection(collection) {
         this.working = true
         this.message = ""
         let addingNew = (collection.id == 0)
         axios.post("/api/collections", collection).then(resp => {
            if ( addingNew ) {
               this.addCollection(resp.data)
            } else {
               this.setCollectionDetail(resp.data)
            }
            this.working = false
            this.setDisplay()
         }).catch((e) => {
            this.setEdit()
            this.message = e
            this.working = false
         })
      }
   },
})
