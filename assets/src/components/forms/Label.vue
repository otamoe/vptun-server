<template>
  <label :class="{'form-label': true, block: block}" v-bind="$attrs">
    <slot></slot>
  </label>
</template>
<style lang="scss">
.form-label {
  margin: .25rem 0;
  &.block {
    display: block
  }
  .form-check-input {
    vertical-align: -0.06em;
  }
}
</style>
<script lang="ts">
export default {
  name: "FormLabel",
  props: {
    block: {
      type: Boolean,
    },
  }
}
</script>
