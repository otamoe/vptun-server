<template>
  <div :class="{'form-flex': true, 'form-inline-flex': inline}" v-bind="$attrs">
    <slot></slot>
  </div>
</template>
<style lang="scss">
.form-flex {
  margin-top: .25em;
  margin-bottom: .25em;
  position: relative;
  display: flex;
  vertical-align: middle;
  align-items: stretch;
  & .btn,
  & .input {
    &:focus,
    &:focus {
      z-index: 1;
      position: relative;
    }
    &:not(:last-child) {
      border-top-right-radius: 0;
      border-bottom-right-radius: 0;
      border-right: 0;
    }
    &:not(:first-child) {
      border-top-left-radius: 0;
      border-bottom-left-radius: 0;
      border-left: 0;
    }
  }
}

.form-inline-flex {
  display: inline-flex;
}
</style>
<script lang="ts">
export default {
  props: {
    inline: Boolean
  }
}
</script>
