<template>
   <dl>
      <dt>Active:</dt>
      <dd>
         <input type="checkbox" v-model="collection.active" id="active"/>
      </dd>
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
            <label class="cb-label" :for="`f${f.id}`">
               <input type="checkbox" :id="`f${f.id}`" :value="f.id" v-model="collection.features" />
               <span>{{ f.name }}</span>
            </label>
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
      <dt>Logo:</dt>
      <dd>
         <div class="logo-wrap">
            <img v-if="collection.imageURL" class="thumb" :src="collection.imageURL"/>
            <span class="drop-wrap">
               <DropZone
                  :maxFiles="1"
                  :clickable="true"
                  :uploadOnDrop="true"
                  :acceptedFiles="['png', 'jpg']"
                  :multipleUpload="false"
                  dropzoneClassName="logo-drop"
                  :url="uploadURL"
                  @addedFile="fileAdded"
                  @removedFile="fileRemoved">
               >
                  <template v-slot:message>
                     <span>Drop new logo here, or click to browse.</span>
                     <span class="note"><b>Note:</b> A newly uploaded logo will replace the current logo upon submission.</span>
                  </template>
               </DropZone>
               <div class="other-opts">
                  <p>OR</p>
                  <uva-button @click="pickImageClicked">Select an existing logo</uva-button>
               </div>
            </span>
         </div>
      </dd>
   </dl>
   <div class="picker-dimmer" v-if="logosOpened">
      <div class="message" role="dialog" aria-modal="true"
         aria-labelledby="picketitle" aria-describedby="msgbody"
         @keyup.esc="dismissLogo"
      >
         <div class="bar">
            <span id="picketitle" class="title">Select a logo</span>
         </div>
         <div class="content">
            <ul class="images">
               <li v-for="(i,idx) in logos" :key="`logo${idx}`" :class="{selected: idx==selectedLogoIdx}" @click="logoClicked(idx)">
                  <img :src="i" />
               </li>
            </ul>
         </div>
         <div class="controls">
            <button id="cancel-logo" @esc="dismissLogo" @click="selectLogo" >
               Cancel
            </button>
            <button id="okl-ogo" @esc="dismissLogo" @click="selectLogo" >
               Select Logo
            </button>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from "vuex"
import UvaButton from './UvaButton.vue'
export default {
   components: { UvaButton },
   computed: {
      ...mapState({
         selectedID: state => state.selectedID,
         details: state => state.details,
         features: state => state.features,
         mode: state => state.mode,
         logos: state => state.logos
      }),
      uploadURL() {
         return `${window.location.href}/api/collections/${this.details.id}/logo`
      }
   },
   data: function()  {
      return {
         logosOpened: false,
         selectedLogoIdx: -1,
         error: "",
         errors: [],
         required: ['title', 'filter', 'itemLabel'],
         collection: {
            id: 0,
            active: false,
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
            imageFile: "",
            imageStatus: "no_change"
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
      pickImageClicked() {
         this.logosOpened = true
         this.$store.dispatch("getLogos")
      },
      logoClicked(idx) {
         this.selectedLogoIdx = idx
      },
      dismissLogo() {
         this.logosOpened = false
         this.selectedLogoIdx = -1;
      },
      selectLogo() {
         this.logosOpened = false
         this.collection.imageURL = this.logos[this.selectedLogoIdx]
         let bits = this.collection.imageURL.split("/")
         this.collection.imageFile = bits[bits.length-1]
         this.selectedLogoIdx = -1
         this.collection.imageStatus = "existing"
      },
      hasError( val) {
         return this.errors.includes(val)
      },
      fileAdded(item) {
         let filename = item.file.name
         this.collection.imageFile = filename
         this.collection.imageStatus = "new"
      },
      fileRemoved(item) {
         this.collection.imageFile = ""
         this.collection.imageStatus = "no_change"
         this.$store.dispatch("deletePendingImage", item.file.name)
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
         this.collection.active = this.details.active
         this.collection.title = this.details.title
         this.collection.description = this.details.description
         this.collection.itemLabel = this.details.itemLabel
         this.collection.startDate = this.details.startDate
         this.collection.endDate = this.details.endDate
         this.collection.filter = this.details.filter
         this.collection.features = []
         this.details.features.forEach( f => {
            this.collection.features.push(f.id)
         })
         if (this.details.image != null) {
            this.collection.imageTitle = this.details.image.title
            this.collection.imageAlt = this.details.image.alt_text
            this.collection.imageURL = this.details.image.url
            let bits = this.collection.imageURL.split("/")
            this.collection.imageFile = bits[bits.length-1]
         }
      }
   }
}
</script>

<style lang="scss" scoped>
.picker-dimmer {
   position: fixed;
   left: 0;
   width: 100%;
   top:0;
   height: 100%;
   background: rgba(0, 0, 0, 0.2);
   .message {
      display: block;
      text-align: left;
      background: white;
      padding: 0px;
      box-shadow:  0 3px 6px rgba(0, 0, 0, 0.16), 0 3px 6px rgba(0, 0, 0, 0.23);
      width: 75%;
      border-radius: 5px;
      margin: 5% auto;
      .bar {
         padding: 5px;
         color: var(--uvalib-text-dark);
         font-weight: 500;
         display: flex;
         flex-flow: row nowrap;
         align-items: center;
         justify-content: space-between;
         background-color: var(--uvalib-blue-alt-light);
         border-bottom: 2px solid var(--uvalib-blue-alt);
         border-radius: 5px 5px 0 0;
         font-size: 1.1em;
         padding: 10px;
      }
      .content {
         margin: 10px 20px 10px 0;
         padding: 0;
      }
      .images {
         box-sizing: border-box;
         margin: 10px;
         padding: 0;
         white-space: nowrap;
         width: 100%;
         overflow-x: auto;
         border: 1px solid var(--uvalib-grey-light);
         display: flex;
         flex-flow: row nowrap;
         align-items: flex-start;
         li {
            display: inline-block;
            width: 200px;
            margin: 10px;

            img {
               max-width: 200px;
               margin: 0;
               padding: 0;
               display: inline-block;
               border: 1px solid var(--uvalib-grey-light);
               box-shadow: 0 1px 3px rgba(0,0,0,.06), 0 1px 2px rgba(0,0,0,.12);
               position: relative;
               cursor: pointer;
               &:hover, &:focus-within, &:focus {
                  top: -2px;
                  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.5), 0 1px 2px rgba(0, 0, 0,1);
               }
            }
         }
         li.selected {
            img {
               outline: 5px solid var(--uvalib-brand-blue-lightest);
            }
         }
      }
      .controls {
         padding: 15px 10px 10px 0;
         text-align: right;
         button {
            margin-left: 5px;
         }
      }
   }
}
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
   .logo-wrap {
      display: flex;
      flex-flow: row nowrap;
      justify-content: flex-start;
   }
   .drop-wrap {
      width: 250px;
      height: 250px;
      margin-left: 10px;
   }
   .other-opts {
      text-align: center;
   }
   .logo-drop {
      border: 2px dashed var(--uvalib-grey-light);
      border-radius: 5px;
      padding: 25px;
      .note {
         display: block;
         margin-top: 10px;
         font-size: 0.85em;
         font-style: italic;
      }
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
         width: 20px;
         height: 20px;
      }
      .cb-label {
         margin: 0 0 10px 0;
         display: flex;
         flex-flow: row nowrap;
         justify-content: flex-start;
         align-items: center;
         input[type=checkbox] {
            margin-right: 10px;
         }
      }
      textarea {
         width: 100%;
         box-sizing: border-box;
         border: 1px solid var(--uvalib-grey-light);
         border-radius: 5px;
         font-family: "franklin-gothic-urw", arial, sans-serif;
         -webkit-font-smoothing: antialiased;
         -moz-osx-font-smoothing: grayscale;
         padding: 5px;
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
