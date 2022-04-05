<template>
<div :class="['dialog-modal', close ? 'hide' : 'open']" tabindex="-1" @keyup.esc="onClose2" >
  <div :class="['dialog-modal-container', 'dialog-modal-size-' + size]">
    <component :is="onSubmit ? 'form': 'tag'" class="dialog-modal-content" @submit.prevent="onSubmit">
      <header class="dialog-modal-header">
        <slot name="header">
          <FormButton v-if="onClose" :title="$t('Close')" style-type="text" class="btn-close" @click.prevent="onClose2">
            <span class="hide">
              {{$t('Close')}}
            </span>
            <svg class="icon icon-close" aria-hidden="true">
              <use xlink:href="#icon-close"></use>
            </svg>
          </FormButton>
          <h5 class="dialog-modal-title" v-if="title">{{title}}</h5>
        </slot>
      </header>
      <div class="dialog-modal-body">
        <slot></slot>
      </div>
      <footer class="dialog-modal-footer">
        <slot name="footer">
          <slot name="footer-before"></slot>
          <FormButton v-if="onCancel" style-type="secondary" @click.prevent="onCancel2">
            {{$t('Cancel')}}
          </FormButton>
          <FormButton v-if="onSubmit" type="submit" :submitting="submitting" style-type="primary">
            {{$t('Submit')}}
          </FormButton>
          <FormButton v-if="onConfirm" style-type="primary" @click.prevent="onConfirm2" :submitting="submitting" >
            {{$t('Confirm')}}
          </FormButton>
          <slot name="footer-after"></slot>
        </slot>
      </footer>
    </component>
  </div>
  <div class="dialog-modal-backdrop" @click.prevent="onClose2"></div>
</div>
</template>
<style lang="scss">
.dialog-modal {
  position: fixed;
  z-index: 991;
  left: 0;
  right: 0;
  bottom: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow-x: hidden;
  overflow-y: auto;
  outline: 0;
  visibility: hidden;
  opacity: 0;
  &.open {
    visibility: visible;
    opacity: 1;
  }
}
.dialog-modal-header {
  padding: .6em 1em;
  border-bottom: 1px solid #dee2e6;
  &.hide, .hide {
    display: none;
  }
  & .btn-close {
    float: right;
    display: block;
    padding: .125em .25em;
    background: transparent;
    border: 0;
    margin: 0;
  }
}
.dialog-modal-title {
  margin: 0;
  font-size:1.25em;
}


.dialog-modal-container {
  z-index: 991;
  margin: 1.5em auto 3em auto;
  display: flex;
  align-items: center;
  min-height: calc(100% - 4.5em);
}
.dialog-modal-content {
  border-radius: .3em;
  z-index: 992;
  position: relative;
  width: 100%;
  background: #fff;
}
.dialog-modal-size-xs {
  max-width: 20em;
}
.dialog-modal-size-sm {
  max-width: 30em;
}
.dialog-modal-size-md {
  max-width: 40em;
}
.dialog-modal-size-lg {
  max-width: 60em;
}
.dialog-modal-size-xl {
  max-width: 80em;
  min-height: 100%;
}
.dialog-modal-size-xxl {
  max-width: 110em;
}
.dialog-modal-size-fl {
  max-width: none;
  margin: 1.75em;
  @media (min-width: 40em) {
    margin: 0;
    min-height: 100%;
    .dialog-modal-content {
      min-height: 100%;   
    }
  }
}

.dialog-modal-body {
  padding: .6em 1em;
}
.dialog-modal-footer {
  padding: .6em 1em;
  text-align: right;
  border-top: 1px solid #dee2e6;
  .btn {
    margin: 0 0 0 .8em;
  }
}
.dialog-modal-backdrop {
  background-color: #000;
  opacity: .2;
  position: fixed;
  z-index: 990;
  top:0;
  left: 0;
  right: 0;
  bottom: 0;
}
</style>
<script lang="ts">
import {SetupContext, defineComponent, watch, onUnmounted, onMounted, getCurrentInstance, nextTick} from "vue"
import FormButton from "@/components/forms/Button.vue"
export default defineComponent({
  components: {
    FormButton,
  },
  name: "DialogModal",
  props: {
    title: {
      type: String,
    },
    size: {
      type: String,
      default: 'sm',
      // xs sm md lg xl fl
    },
    close: {
      type: Boolean,
      default: false
    },
    submitting: {
      type: Boolean,
      default: false
    },
    onClose: {
    },
    onConfirm: {
    },
    onCancel: {
    },
    onSubmit: {
    },
  },
  setup(props: any, context: SetupContext) {
    const internalInstance = getCurrentInstance() // works
    const onClose2 = () => {
      context.emit("close")
    }
    const onConfirm2 = () => {
      context.emit("confirm")
    }
    const onCancel2 = () => {
      context.emit("cancel")
    }
    watch(()=> props.close, function() {
      
      nextTick(function() {
        // @ts-ignore
        internalInstance.ctx.$el.focus()
      })
    })
    
    
    return {
      onConfirm2,
      onCancel2,
      onClose2,
    }
  }
})
</script>
