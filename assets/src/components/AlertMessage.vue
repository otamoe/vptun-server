<template>
  <div role="alert" :class="['alert-message', 'alert-message-' + type, forValue.length && !close ? 'show' : 'hide']">
    <FormButton v-if="onClose" :title="$t('Close')" style-type="text" @click.prevent="onClose2">
      <span class="hide">
        {{$t('Close')}}
      </span>
      <svg class="icon icon-close" aria-hidden="true">
        <use xlink:href="#icon-close"></use>
      </svg>
    </FormButton>
    <slot>
      <p v-for="(text, key) in forValue" :key="key">
        {{ translate(text) }}
      </p>
    </slot>
  </div>
</template>
<style lang="scss"> 
@import "@/scss/variable.scss";
@mixin alert-message-type-mixin($value) {
  color: map-get($value, color);
  background: map-get($value, background);
  box-shadow: shade-color(map-get($value, background), 20);
  & .btn-text {
    color: map-get($value, color);
    &:focus{
      box-shadow: 0 0 0 0.2em rgba(map-get($value, color), .35);
    }
  }
}

@mixin alert-messages-mixin() {  
  @each $name, $value in $alert-messages {
    @if $name == "default" {
      .alert-message {
        @include alert-message-type-mixin($value);
      }  
    } @else {
      .alert-message-#{$name} {
        @include alert-message-type-mixin($value);
      }
    }
  }
}

.alert-message{
  border: 1px solid transparent;
  border-radius: .25em;
  padding: 1em;
  margin-bottom: 1em;

  &.hide, .hide {
    display: none;
  }
  & p {
    margin: 0 0 .75em 0;
    &:last-child{
      margin-bottom: 0
    }
  }
  & .btn-text {
    float: right;
    display: block;
    padding: 0 .25em;
    background: transparent;
    border: 0;
    margin: 0;
  }
}
@include alert-messages-mixin();
</style>

<script lang="ts">
import {SetupContext, computed} from "vue"
import FormButton from "@/components/forms/Button.vue"

export default {
  components: {
    FormButton,
  },
  name: "AlertMessage",
  props: {
    type: {
      type: String,
      default: 'default',
    },
    value: {},
    
    close: {
      type: Boolean,
      default: false
    },
    onClose: {
    },
  },
  setup(props: any, context: SetupContext) {
    const translate = function(message: any) {
      let value
      if (!message || typeof message !== 'object') {
        value = message
      } else if (message.message) {
        value = message.message
      } else if (message.error) {
        value = message.error
      } else {
        value = message
      }
      return value
    }

    const onClose2 = () => {
      context.emit("close")
    }

    const forValue = computed((): any => {
      let value = props.value
      if (!value) {
        return []
      }

      if (value?.graphQLErrors?.length) {
        return value.graphQLErrors
      }
      if (value?.networkError?.result?.errors?.length) {
        return value.networkError.result.errors
      }

      if (value?.errors?.length) {
        return value.errors
      }

      // array
      if (value instanceof Array) {
        let messages: any[] = []
        for (let i = 0; i < value.length; i++) {
          const message = value[i];
          if (!message || typeof message !== 'object') {
            messages.push({message})
          } else {
            messages.push(message)
          }
        }
        return messages
      }

    if (value !== 'object') {
        return [{message: value}]
      }
      return [value]
    })

    return {translate, forValue, onClose2}
  }
}
</script>
