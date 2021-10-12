<template>
   <button tabindex="0" class="uva-button" :class="{disabled: disabled}"
      @keydown.exact.tab="tabNext"
      @keydown.shift.tab="tabBack"
      @click.prevent.stop="clicked" @keydown.prevent.stop.enter="clicked" @keydown.space.prevent.stop="clicked" @keyup.stop.esc="escClicked">
      <slot></slot>
   </button>
</template>

<script>
export default {
   props: {
      focusNextOverride: {
         type: Boolean,
         default: false
      },
      focusBackOverride: {
         type: Boolean,
         default: false
      },
      disabled: {
         type: Boolean,
         default: false
      }
   },
   emits: ['click', 'esc', 'tabback', 'tabnext' ],
   methods: {
      escClicked() {
         if (!this.disabled) {
            this.$emit('esc')
         }
      },
      clicked() {
         if (!this.disabled) {
            this.$emit('click')
         }
      },
      tabBack(event) {
         if (this.focusBackOverride ) {
            event.stopPropagation()
            event.preventDefault()
            this.$emit('tabback')
         }
      },
      tabNext(event) {
         if (this.focusNextOverride ) {
            event.stopPropagation()
            event.preventDefault()
            this.$emit('tabnext')
         }
      }
   },
}
</script>

<style lang="scss" scoped>
button.uva-button {
   font-weight: normal;
   padding: 8px 20px;
   border-radius: 5px;
   cursor: pointer;
   background-color: var(--uvalib-brand-blue-light);
   border: 1px solid var(--uvalib-brand-blue-light);
   color: white;
   &:hover {
      background-color: var(--uvalib-brand-blue-lighter);
   }
}
.uva-button.disabled {
   cursor: default;
   opacity: 0.25;
}
</style>