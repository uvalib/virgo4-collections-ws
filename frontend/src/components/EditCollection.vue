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
               <FileUpload name="file" :url="`/api/collections/${collection.details.id}/logo`"
                  accept="image/*" :fileLimit="1"
                  :showCancelButton="false" :multiple="false" @upload="fileUploaded"
                  chooseLabel="Browse images" @removeUploadedFile="fileRemoved" @remove="fileRemoved" @select="fileAdded"
               >
                  <template #empty>
                     <span>Drop new logo here, or click to browse</span>
                  </template>
               </FileUpload>
               <Button @click="pickImageClicked" label="Select an existing logo" />
            </span>
         </div>
      </dd>
   </dl>
   <div class="controls">
      <Button severity="secondary" @click="collection.setDisplay()" label="Cancel"/>
      <Button v-if="collection.selectedID > 0" @click="submitChanges()" label="Submit"/>
      <Button v-else @click="submitChanges()" label="Create"/>
   </div>

   <Dialog v-model:visible="logosOpened" :modal="true" header="Select a Logo" style="width:80%" @hide="dismissLogo()">
      <Carousel :value="collection.logos" :numVisible="5" :numScroll="5">
         <template #item="slotProps">
            <div class="logos">
               <img class="logo" :src="slotProps.data"
                  :class="{selected: slotProps.index == selectedLogoIdx}"
                  @click="logoClicked(slotProps.index)"
               />
            </div>
         </template>
      </Carousel>
      <template #footer>
         <Button @click="dismissLogo"severity="secondary" label="Cancel"/>
         <Button @click="selectLogo" label="Select Logo" :disabled="selectedLogoIdx < 0"/>
      </template>
   </Dialog>

</template>

<script setup>
import { useCollectionStore } from "@/stores/collection"
import { onMounted, ref } from "vue"
import FileUpload from 'primevue/fileupload'
import { useConfirm } from "primevue/useconfirm"
import Dialog from 'primevue/dialog'
import Carousel from 'primevue/carousel'

const confirm = useConfirm()
const collection = useCollectionStore()

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
   selectedLogoIdx.value = -1
   edit.value.imageStatus = "existing"
}

function hasError( val) {
   return errors.value.includes(val)
}

function fileAdded(event) {
   let file = event.files[0]
   let filename = file.name
   edit.value.imageFile = filename
   edit.value.imageStatus = "new"
}

function fileUploaded() {
   edit.value.imageStatus = "uploaded"
}

function fileRemoved(event) {
   if ( edit.value.imageStatus == "uploaded") {
      collection.deletePendingImage(event.file.name)
   }
   edit.value.imageFile = ""
   edit.value.imageStatus = "no_change"
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
      return
   }

   if ( edit.value.imageStatus== "new") {
       confirm.require({
         message: `You have selected a new logo, but not uploaded it. Submitting now will not change the logo. Continue?`,
         header: 'Confirm Submit',
         icon: 'pi pi-question-circle',
         rejectProps: {
            label: 'Cancel',
            severity: 'secondary'
         },
         acceptProps: {
            label: 'Submit without logo'
         },
         accept: () => {
            collection.submitCollection(edit.value)
         },
      })
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

.logos {
   padding: 10px 0;
   img.logo {
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

   img.logo.selected {
      outline: 5px solid $uva-brand-blue-300;
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
      align-items: flex-start;
      gap: 20px;
      .drop-wrap {
          width: 375px;
         display: flex;
         flex-direction: column;
         gap: 10px;
         :deep(.p-fileupload-header .p-button) {
            flex-grow: 1;
         }
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
.controls {
   padding-top: 20px;
   display: flex;
   flex-flow: row nowrap;
   justify-content: flex-end;
   gap: 10px;
   border-top: 2px solid $uva-grey
}
</style>
