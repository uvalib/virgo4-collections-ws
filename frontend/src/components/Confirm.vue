<template>
   <uva-button @click="triggerClicked" :class="{disabled: isOpen}" >{{buttonText}}</uva-button>
   <div v-if="isOpen" class="messsage-box">
      <div class="message" role="dialog" aria-modal="true"
         aria-labelledby="msgtitle" aria-describedby="msgbody"
         @keyup.esc="dismiss"
      >
         <div class="bar">
            <span id="msgtitle" class="title">Confirm Action</span>
         </div>
         <div class="message-body" id="msgbody" v-html="message"></div>
         <div class="message-body">Continue?</div>
         <div class="controls">
            <button id="cancelbtn" class="pad" @esc="dismiss" @click="dismiss" >
               Cancel
            </button>
            <button id="okbtn" @esc="dismiss" @click="ok" >
               OK
            </button>
         </div>
      </div>
   </div>
</template>

<script>
export default {
   emits: ['cancel', 'confirm' ],
   props: {
      message: {
         type: String,
         reqired: true
      },
      buttonText: {
         type: String,
         reqired: true
      },
   },
   data: function()  {
      return {
         isOpen: false
      }
   },
   methods: {
      triggerClicked() {
         this.isOpen = true
      },
      ok() {
         this.isOpen = false
         this.$emit('confirm')
      },
      dismiss() {
         this.isOpen = false
         this.$emit('cancel')
      },
   },
};
</script>

<style lang="scss" scoped>

div.messsage-box {
   position: fixed;
   left: 0;
   top: 0;
   width: 100%;
   height: 100%;
   z-index: 1000;
   background: rgba(0, 0, 0, 0.2);

   .message {
      box-shadow:  0 3px 6px rgba(0, 0, 0, 0.16), 0 3px 6px rgba(0, 0, 0, 0.23);
      font-size: 16px;
      position: fixed;
      height: auto;
      z-index: 8000;
      background: white;
      top: 15%;
      left: 50%;
      transform: translate(-50%, -15%);
      border-radius: 5px;
      min-width: 300px;
      word-break: break-word;

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

      .message-body {
         max-height: 55vh;
         overflow-y: auto;
         text-align: left;
         padding: 20px 30px 0 30px;
         font-weight: normal;
         opacity: 1;
         visibility: visible;
         text-align: left;
         word-break: break-word;
         -webkit-hyphens: auto;
         -moz-hyphens: auto;
         hyphens: auto;
         color: var(--uvalib-primary-text);
      }

      .controls {
         padding: 15px 10px 10px 0;
         font-size: 0.9em;
         text-align: right;
         .pad {
            margin-right: 5px;
         }
      }
   }
}

</style>
