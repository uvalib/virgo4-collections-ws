<template>
   <dl>
      <dt>Active:</dt>
      <dd>
         <input type="checkbox" v-model="edit.active" id="active"/>
      </dd>
      <dt>Title:</dt>
      <dd>
         <input type="text" v-model="edit.title" id="title" aria-required="true" required="required"/>
         <span v-if="hasError('title')" class="error">Collection title is required</span>
      </dd>
      <dt>Description:</dt>
      <dd>
         <textarea rows="4" v-model="edit.description" id="description"></textarea>
      </dd>
      <dt>Item Label:</dt>
      <dd>
         <input type="text" v-model="edit.itemLabel" id="item-label"/>
      </dd>
      <dt>Start Date:</dt>
      <dd>
         <input type="text" v-model="edit.startDate" id="start-date"/>
      </dd>
      <dt>End Date:</dt>
      <dd>
         <input type="text" v-model="edit.endDate" id="end-date"/>
      </dd>
      <dt>Facet Name:</dt>
      <dd>
         <input type="text" v-model="edit.filter" id="filter" aria-required="true" required="required"/>
         <span v-if="hasError('filter')" class="error">Facet Name is required</span>
      </dd>
      <dt>Features:</dt>
      <dd>
         <template v-for="f in collection.features" :key="`f${f.id}`">
            <label class="cb-label" :for="`f${f.id}`">
               <input type="checkbox" :id="`f${f.id}`" :value="f.id" v-model="edit.features" />
               <span>{{ f.name }}</span>
            </label>
         </template>
         <div>
            <span v-if="hasError('features')" class="error">At least one feature is required</span>
         </div>
      </dd>
      <dt>Logo Title:</dt>
      <dd>
         <input type="text" v-model="edit.imageTitle" id="image-title"/>
      </dd>
      <dt>Logo Alt Text:</dt>
      <dd>
         <input type="text" v-model="edit.imageAlt" id="image-alt"/>
      </dd>
      <dt>Logo:</dt>
      <dd>
         <div class="logo-wrap">
            <img v-if="edit.imageURL" class="thumb" :src="edit.imageURL"/>
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
                  <Button @click="pickImageClicked" label="Select an existing logo"/>
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
               <li v-for="(i,idx) in collection.logos" :key="`logo${idx}`" :class="{selected: idx==selectedLogoIdx}" @click="logoClicked(idx)">
                  <img :src="i" />
               </li>
            </ul>
         </div>
         <div class="controls">
            <Button @click="dismissLogo"severity="secondary" label="Cancel"/>
            <Button @click="selectLogo" label="Select Logo"/>
         </div>
      </div>
   </div>
</template>

<script setup>
import { useCollectionStore } from "@/stores/collection"
import { storeToRefs } from "pinia";
import { onMounted, ref, watch, computed } from "vue"
const collection = useCollectionStore()
const { mode } = storeToRefs (collection)

const uploadURL = computed(()=>{
   return `${window.location.href}api/collections/${collection.details.id}/logo`
})

const required = ['title', 'filter', 'itemLabel']
const logosOpened = ref(false)
const selectedLogoIdx = ref(-1)
const errors = ref([])
const edit = ref({
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
})

watch(mode, (newValue, _oldValue) => {
   if (newValue =="submit") {
      submitChanges()
   }
})


function pickImageClicked() {
   logosOpened.value = true
   collection.getLogos()
}
function logoClicked(idx) {
   selectedLogoIdx.value = idx
}
function dismissLogo() {
   logosOpened.value = false
   selectedLogoIdx.value = -1;
}
function selectLogo() {
   logosOpened.value = false
   edit.value.imageURL = collection.logos[selectedLogoIdx.value]
   let bits = edit.value.imageURL.split("/")
   edit.value.imageFile = bits[bits.length-1]
   console.log(edit.value.imageFile)
   selectedLogoIdx.value = -1
   edit.value.imageStatus = "existing"
}
function hasError( val) {
   return errors.value.includes(val)
}
function fileAdded(item) {
   let filename = item.file.name
   edit.value.imageFile = filename
   edit.value.imageStatus = "new"
}
function fileRemoved(item) {
   edit.value.imageFile = ""
   edit.value.imageStatus = "no_change"
   collection.deletePendingImage(item.file.name)
}
function submitChanges() {
   errors.value.splice(0, errors.value.length)
   for (let [key, value] of Object.entries(edit.value)) {
      if ( key == "features") {
         if ( value.length == 0) {
            errors.value.push(key)
         }
      } else if ( required.includes(key) && value == "") {
         errors.value.push(key)
      }
   }

   let focused = false
   if (errors.value.length > 0) {
      let tgtID = errors.value[0]
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
      collection.setEdit()
   } else {
      collection.submitCollection(edit.value)
   }
}

onMounted(()=>{
   if ( collection.selectedID > 0) {
      edit.value.id = collection.details.id
      edit.value.active = collection.details.active
      edit.value.title = collection.details.title
      edit.value.description = collection.details.description
      edit.value.itemLabel = collection.details.itemLabel
      edit.value.startDate = collection.details.startDate
      edit.value.endDate = collection.details.endDate
      edit.value.filter = collection.details.filter
      edit.value.features = []
      collection.details.features.forEach( f => {
         edit.value.features.push(f.id)
      })
      if (collection.details.image != null) {
         edit.value.imageTitle = collection.details.image.title
         edit.value.imageAlt = collection.details.image.alt_text
         edit.value.imageURL = collection.details.image.url
         let bits = edit.value.imageURL.split("/")
         edit.value.imageFile = bits[bits.length-1]
      }
   }
})
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
      width: 75%;
      border-radius: 5px;
      margin: 5% auto;
      .bar {
         padding: 5px;
         font-weight: 500;
         display: flex;
         flex-flow: row nowrap;
         align-items: center;
         justify-content: space-between;
         background-color: v$uva-blue-alt-300;
         border-bottom: 2px solid $uva-blue-alt;
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
         border: 1px solid $uva-grey-100;
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
               border: 1px solid $uva-grey-100;
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
               outline: 5px solid $uva-brand-blue-300;
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
      color: $uva-red-A;
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
      border: 2px dashed $uva-grey-100;
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
         border: 1px solid $uva-grey-100;
         border-radius: 5px;
         font-family: "franklin-gothic-urw", arial, sans-serif;
         -webkit-font-smoothing: antialiased;
         -moz-osx-font-smoothing: grayscale;
         padding: 5px;
      }
      .thumb {
         max-width: 200px;
         border:1px solid $uva-grey-100;
      }
   }
   .na {
      color: #aaa;
      font-style: italic;
   }
}
</style>
