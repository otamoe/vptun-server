<template>
  <div class="form-group" v-bind="$attrs">
    <slot></slot>
  </div>
</template>
<style lang="scss">
.form-group {
  margin-bottom: 1rem
}
</style>
<script lang="ts">
export default {
  name: "FormGroup",
  props: {
    label: String
  }
}
</script>
