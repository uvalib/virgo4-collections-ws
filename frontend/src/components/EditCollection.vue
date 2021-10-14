<template>
   <dl>
      <dt>Title:</dt>
      <dd>
         <input type="text" v-model="collection.title" id="title" aria-required="true" required="required"/>
         <span v-if="hasError('title')" class="error">Collection title is required</span>
      </dd>
      <dt>Description:</dt>
      <dd>
         <textarea rows="4" v-model="collection.description" id="description"></textarea>
      </dd>
      <dt>Item Label:</dt>
      <dd>
         <input type="text" v-model="collection.itemLabel" id="item-label"/>
      </dd>
      <dt>Start Date:</dt>
      <dd>
         <input type="text" v-model="collection.startDate" id="start-date"/>
      </dd>
      <dt>End Date:</dt>
      <dd>
         <input type="text" v-model="collection.endDate" id="end-date"/>
      </dd>
      <dt>Facet Name:</dt>
      <dd>
         <input type="text" v-model="collection.filter" id="filter" aria-required="true" required="required"/>
         <span v-if="hasError('filter')" class="error">Facet Name is required</span>
      </dd>
      <dt>Features:</dt>
      <dd>
         <template v-for="f in features" :key="`f${f.id}`">
            <input type="checkbox" :id="`f${f.id}`" :value="f.id" v-model="collection.features" />
            <label class="cb-label" :for="`f${f.id}`">{{ f.name }}</label>
         </template>
         <div>
            <span v-if="hasError('features')" class="error">At least one feature is required</span>
         </div>
      </dd>
      <dt>Logo Title:</dt>
      <dd>
            <input type="text" v-model="collection.imageTitle" id="image-title"/>
      </dd>
      <dt>Logo Alt Text:</dt>
      <dd>
         <input type="text" v-model="collection.imageAlt" id="image-alt"/>
      </dd>
      <template v-if="collection.imageURL">
         <dt>Logo:</dt>
         <dd>
            <img class="thumb" :src="collection.imageURL"/>
         </dd>
      </template>
   </dl>
</template>

<script>
import { mapState } from "vuex"
export default {
   computed: {
      ...mapState({
         selectedID: state => state.selectedID,
         details: state => state.details,
         features: state => state.features,
         mode: state => state.mode
      })
   },
   data: function()  {
      return {
         error: "",
         errors: [],
         required: ['title', 'filter', 'itemLabel'],
         collection: {
            id: 0,
            title: "",
            description: "",
            itemLabel: "Issue",
            startDate: "",
            endDate: "",
            filter: "",
            features: [],
            imageTitle: "",
            imageAlt: "",
            imageURL: "",
            imageFile: ""
         }
      }
   },
   watch: {
      mode(newVal, _oldVal)  {
         if (newVal=="submit") {
            this.submitChanges()
         }
      }
   },
   methods: {
      hasError( val) {
         return this.errors.includes(val)
      },
      submitChanges() {
         this.errors.splice(0, this.errors.length)
         for (let [key, value] of Object.entries(this.collection)) {
            if ( key == "features") {
               if ( value.length == 0) {
                  this.errors.push(key)
               }
            } else if ( this.required.includes(key) && value == "") {
               this.errors.push(key)
            }
         }

         let focused = false
         if (this.errors.length > 0) {
            let tgtID = this.errors[0]
            if (tgtID == "itemLabel") {
               tgtID = "item-label"
            }
            if (!focused) {
               let first = document.getElementById(tgtID)
               if ( first ) {
                  first.focus()
                  focused = true
               }
            }
            this.$store.commit("setEdit")
         } else {
            this.$store.dispatch("submitCollection", this.collection)
         }
      }
   },
   mounted() {
      if ( this.selectedID > 0) {
         this.collection.id = this.details.id
         this.collection.title = this.details.title
         this.collection.description = this.details.description
         this.collection.itemLabel = this.details.itemLabel
         this.collection.startDate = this.details.startDate
         this.collection.endDate = this.details.endDate
         this.collection.filter = this.details.filter
         this.collection.features = []
         this.details.features.forEach( f => {
            let data = this.features.find( sf => sf.name == f)
            if (data) {
               this.collection.features.push(data.id)
            }
         })
         if (this.details.image) {
            this.collection.imageTitle = this.details.image.title
            this.collection.imageAlt = this.details.image.alt_text
            this.collection.imageURL = this.details.image.url
         }
      }
   }
}
</script>

<style lang="scss" scoped>
dl {
   margin-left: 25px;
   display: inline-grid;
   grid-template-columns: max-content 2fr;
   grid-column-gap: 15px;
   width: 95%;
   span.error {
      color: var(--uvalib-red-emergency);
      font-style: italic;
      display: inline-block;
      margin-top: 5px;
   }
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
      input[type=text] {
         width: 100%;
      }
      input[type=checkbox] {
         margin: 0 5px 0 0;
      }
      .cb-label {
         margin-right: 20px;
      }
      textarea {
         width: 100%;
         box-sizing: border-box;
         border: 1px solid var(--uvalib-grey-light);
         border-radius: 5px;
      }
      .thumb {
         max-width: 200px;
         border:1px solid var(--uvalib-grey-light);
      }
   }
   .na {
      color: #aaa;
      font-style: italic;
   }
}
</style>
