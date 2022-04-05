<template>
  <select
    v-if="type === 'select'"
    :class="className"
    :value="modelValue"
    :disabled="disabled"
    v-bind="$attrs"
    v-on:input="onValue"
    v-on:change="onValue"
  >
    <slot></slot>
  </select>
  <textarea
    v-else-if="type === 'textarea'"
    :class="className"
    :value="modelValue"
    :disabled="disabled"
    v-bind="$attrs"
    v-on:input="onValue"
    v-on:change="onValue"
  ></textarea>
  <input
    v-else-if="isCheck"
    :type="type"
    :class="className"
    :disabled="disabled"
    v-bind="$attrs"
    v-on:input="onValue"
    v-on:change="onValue"
    :value="check"
    :checked="checked"
  >
  <RouterLink
    v-else-if="type === 'route'"
    :class="className"
    :disabled="disabled"
    v-bind="$attrs"
    :checked="checked"
    :to="'/'"
  >
    <slot></slot>
  </RouterLink>
  <a
    v-else-if="type === 'route' || type === 'link' || type === 'a'"
    :class="className"
    :disabled="disabled"
    v-bind="$attrs"
    :checked="checked"
    :to="'/'"
  >
    <slot></slot>
  </a>
  <input
    v-else
    :type="type"
    :class="className"
    :value="modelValue"
    :disabled="disabled"
    v-bind="$attrs"
    v-on:input="onValue"
    v-on:change="onValue"
  >
</template>
<style lang="scss">
@import "@/scss/variable.scss";


@mixin input-mixin-value($value) {
  @if map-has-key($value, color) {
    color: map-get($value, color);
  }
  @if map-has-key($value, background) {
    background: map-get($value, background);
  }
  @if map-has-key($value, border-color) {
    border-color: map-get($value, border-color);
  }

  &:focus,
  &.focus {
    @if map-has-key($value, border-color) {
      box-shadow: 0 0 0 0.2em rgba(map-get($value, border-color), .5);
    }
  }
  &:hover,
  &.hover {
    @if map-has-key($value, border-color) {
      border-color: shade-color(map-get($value, border-color), 10);
    }
  }
  &:active,
  &.active {
    @if map-has-key($value, border-color) {
      border-color: shade-color(map-get($value, border-color), 15);
    }
  }
}

@mixin inputs-mixin() {  
  @each $name, $value in $inputs {
    @if $name == "default" {
      .input {
        @include input-mixin-value($value);
      }  
    } @else if $name == "placeholder" {
      .input::placeholder {
        @include input-mixin-value($value);
      }
    } @else if $name == "valid" {
      .input.valid, .input.validity:valid:not(:placeholder-shown) {
        @include input-mixin-value($value);
      }
    } @else if $name == "invalid" {
      .input.invalid, .input.validity:invalid:not(:placeholder-shown) {
        @include input-mixin-value($value);
      }
    } @else if $name == "readonly" {
      .input.readonly, .input[readonly] {
        @include input-mixin-value($value);
      }
    } @else {
      .input.#{$name}, .input:#{$name} {
        @include input-mixin-value($value);
      }
    }
  }
}

.input{
  border-radius: .2em;
  text-decoration: none;
  padding: .375em .75em;
  position: relative;
  outline: 0;
  font-size: 1em;
  line-height: 1.5;
  display: inline-block;
  font-weight: 400;
  vertical-align: middle;
  transition: border-color .15s ease-in-out,box-shadow .15s ease-in-out;
  margin-top: .25em;
  margin-bottom: .25em;
  border: 1px solid rgba(0,0,0, .1);
  &:focus,
  &.focus{
    outline: 0;
  }
  &.invalid, &.validity:invalid:not(:placeholder-shown) {
    &~.form-invalid {
      display: block;
    }
  }
  &[readonly],
  &.readonly {
    opacity: .85;
  }
  &:disabled,
  &.disabled {
    opacity: .85;
  }
}


@include inputs-mixin();


select.input {
    height: calc(2.25em + 2px);
}
select[multiple].input, select[size].input {
    height: auto;
}

.input-block {
  width: 100%;
  display: block;
}

.input-xs {
  font-size: 0.75em;
}
.input-sm {
  font-size: 0.875em;
}
.input-lg {
  font-size: 1.25em;
}
.input-xl {
  font-size: 1.5em;
}
</style>
<script lang="ts">
import {computed, SetupContext} from "vue"
export default {
  name: "FormInput",
  props: {
    block: {
      type: Boolean,
    },
    type: {
      type: String,
      default: 'text',
    },
    styleSize: {
      type: String,
      default: 'md',
    },
    check: {
    },
    
    strict: {
      type: Boolean,
    },
    
    modelValue: {
    },
    disabled: {
    },
  },
  setup(props: any, context: SetupContext): any {
    const isCheck = computed(() => {
      return props.type === 'radio' || props.type === 'checkbox'
    })

    const checked = computed(() => {
      switch (props.type) {
        case 'radio':
          return props.check === props.modelValue
        case 'checkbox':
          if (props.check === undefined) {
            return Boolean(props.modelValue)
          }
          if (!Array.isArray(props.modelValue) || props.strict) {
            return props.modelValue === props.check
          }
          return props.modelValue.indexOf(props.check) !== -1
        default:
          return false
      }
    })

    const className = computed((): any => {
      return [
        isCheck.value ? 'form-check-input' : 'input',
        isCheck.value ? 'form-check-' + props.styleSize : 'input-' + props.styleSize,
        props.disabled ? 'disabled' : undefined,
        checked.value ? 'checked' : undefined,
        props.block ? 'input-block' : undefined,
      ]
    })

    const onValue = (e: any)  => {
      let el = e.target
      let value
      switch (props.type) {
        case 'checkbox':
          if (props.check === undefined) {
            value = el.checked
          } else if (!Array.isArray(props.modelValue)) {
            value = el.checked ? props.check : undefined
          } else if (props.strict) {
            value = el.checked ? props.check : []
          } else {
            value = props.modelValue
            let index = value.indexOf(props.check)
            if (el.checked) {
              if (index === -1) {
                value = value.concat([props.check])
              }
            } else if (index !== -1) {
              value = value.concat([])
              value.splice(index, 1)
            }
          }
          break
        default:
          value = el.value
      }
      context.emit('update:modelValue', value)
    }

    return {
      isCheck,
      checked,
      className,
      onValue,
    };
  }
}
</script>
