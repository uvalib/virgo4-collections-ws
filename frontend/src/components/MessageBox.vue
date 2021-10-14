<template>
   <div v-if="message" class="messsage-box">
      <div class="message" role="dialog" aria-modal="true"
         aria-labelledby="msgtitle" aria-describedby="msgbody"
         @keyup.esc="dismiss"
      >
         <div class="bar">
            <span id="msgtitle" class="title">System Message</span>
         </div>
         <div class="message-body" id="msgbody" v-html="message"></div>
         <div class="controls">
            <button id="okbtn" @esc="dismiss" @click="dismiss" >
               OK
            </button>
         </div>
      </div>
   </div>
</template>

<script>
import { mapState } from "vuex"
export default {
   computed: {
      ...mapState({
         message: state => state.message,
      }),
   },
   methods: {
      dismiss() {
         this.$store.commit("clearMessage")
      },
   },
   created() {
      let ele = document.getElementById("okbtn")
      if (ele ) {
               ele.focus()
            }
   },
};
</script>

<style lang="scss" scoped>

div.messsage-box {
   position: fixed;
   left: 0;
   right: 0;
   z-index: 9999;
   top: 20%;

   .details {
      text-align: left;
      padding: 0 30px 20px 30px;
   }

   .message {
      display: inline-block;
      text-align: left;
      background: white;
      padding: 0px;
      box-shadow:  0 3px 6px rgba(0, 0, 0, 0.16), 0 3px 6px rgba(0, 0, 0, 0.23);
      min-width: 20%;
      max-width: 80%;
      border-radius: 5px;
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
      }
   }
}

</style>
