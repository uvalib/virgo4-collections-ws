<template>
   <button tabindex="0" class="uva-button" :class="{disabled: disabled}"
      @keydown.exact.tab="tabNext"
      @keydown.shift.tab="tabBack"
      @click.prevent.stop="clicked" @keydown.prevent.stop.enter="clicked" @keydown.space.prevent.stop="clicked" @keyup.stop.esc="escClicked">
      <slot></slot>
   </button>
</template>

<script setup>
const props = defineProps({
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
})
const emit = defineEmits( ['click', 'esc', 'tabback', 'tabnext' ] )

function escClicked() {
   if (!props.disabled) {
      emit('esc')
   }
}
function clicked() {
   if (!props.disabled) {
      emit('click')
   }
}
function tabBack(event) {
   if (props.focusBackOverride ) {
      event.stopPropagation()
      event.preventDefault()
      emit('tabback')
   }
}
function tabNext(event) {
   if (props.focusNextOverride ) {
      event.stopPropagation()
      event.preventDefault()
      emit('tabnext')
   }
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