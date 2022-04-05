<template>
  <RouterLink
    v-if="!disabled && type === 'route'"
    :disabled="disabled || submitting"
    :to="to || href"
    :class="className"
    v-bind="$attrs"
  >
    <span v-if="submitting" class="btn-submitting">
      <svg class="icon icon-loading" aria-hidden="true">
        <use xlink:href="#icon-loading"></use>
      </svg>
    </span>
    <slot></slot>
  </RouterLink>
  <a
    v-else-if="type === 'route'|| type === 'link' || type == 'a'"
    :disabled="disabled || submitting"
    :href="disabled || submitting ? '#' : href"
    :class="className"
    v-bind="$attrs"
    @click="onClick"
  >
    <span v-if="submitting" class="btn-submitting">
      <svg class="icon icon-loading" aria-hidden="true">
        <use xlink:href="#icon-loading"></use>
      </svg>
    </span>
    <slot></slot>
  </a>
  <label
    v-else-if="type === 'label'"
    :type="type"
    :disabled="disabled || submitting"
    :class="className"
    v-bind="$attrs"
  >
    <span v-if="submitting" class="btn-submitting">
      <svg class="icon icon-loading" aria-hidden="true">
        <use xlink:href="#icon-loading"></use>
      </svg>
    </span>
    <slot></slot>
  </label>
  <button
    v-else
    :type="type"
    :disabled="disabled || submitting"
    :class="className"
    v-bind="$attrs"
  >
    <span v-if="submitting" class="btn-submitting">
      <svg class="icon icon-loading" aria-hidden="true">
        <use xlink:href="#icon-loading"></use>
      </svg>
    </span>
    <slot></slot>
  </button>
</template>
<style lang="scss">
@import "@/scss/variable.scss";

@mixin btn-mixin-value($value) {
  $color: map-get($value, color);
  $background: map-get($value, background);

  color: $color;
  background: $background;

  @if $background != transparent {
    border-color: shade-color($background, 7);
  } @else {
    border-color: transparent;
  }
  &:focus,
  &.focus {
    @if $background != transparent {
      box-shadow: 0 0 0 0.2em rgba($background, .5);
    } @else {
      box-shadow: 0;
    }
  }
  &:hover,
  &.hover {
    @if $background != transparent {
      background: shade-color($background, 7);
      border-color: shade-color($background, 10);
    } @else {
      background: transparent;
      border-color: transparent;
    }
  }
  &:disabled,
  &.disabled {
    &:hover,
    &.hover {
      color: $color;
      background: $background;
    }
  }
  &:active,
  &.active {
    @if $background != transparent {
      background: shade-color($background, 14);
      border-color: shade-color($background, 20);
    } @else {
      background: transparent;
      border-color: transparent;
    }
  }
  &.outline {
    @if $background != transparent {
      color: $background;
      &:focus,
      &.focus {
        color: $background;
        box-shadow: 0 0 0 0.15em rgba($background, .3);
      }
      &:hover,
      &.hover,
      &:active,
      &.active {
        background: $background;
        color: $color;
      }
      &:disabled,
      &.disabled {
        &:hover,
        &.hover {
          color: $background;
          background: transparent;
        }
      }
    }
  }
}

@mixin btns-mixin() {  
  @each $name, $value in $btns {
    @if $name == "default" {
      .btn {
        @include btn-mixin-value($value);
      }  
    } @else {
      .btn-#{$name} {
        @include btn-mixin-value($value);
      }
    }
  }
}



.btn {
  border-radius: .2em;
  text-decoration: none;
  border: 1px solid transparent;
  padding: .375em .75em;
  position: relative;
  outline: 0;
  font-size: 1em;
  line-height: 1.5;
  cursor: pointer;
  display: inline-block;
  font-weight: 400;
  text-align: center;
  vertical-align: middle;
  user-select: none;
  transition: color .15s ease-in-out,background-color .15s ease-in-out,border-color .15s ease-in-out,box-shadow .15s ease-in-out;
  margin-top: .25em;
  margin-bottom: .25em;
  white-space: nowrap;
  &:disabled,
  &.disabled {
    cursor: default;
    opacity: .65;
  }
  &:hover,
  &.hover,
  &:active,
  &.active {
    text-decoration: none;
  }
  &.outline {
    background: transparent;
  }
}
@include btns-mixin();

.btn-submitting {
  position: absolute;
  z-index: 1;
  background: inherit;
  text-align: center;
  left: 0;
  right: 0;
  svg {
    font-size:1.5em;
    animation: btn-submitting 1s linear infinite;
  }
}

@keyframes btn-submitting {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}

.btn-block {
  width: 100%;
  display: block;
}


.btn-xs {
  font-size: 0.75em;
}
.btn-sm {
  font-size: 0.875em;
}
.btn-lg {
  font-size: 1.25em;
}
.btn-xl {
  font-size: 1.5em;
}
</style>

<script>
export default {
  props: {
    type: {
      type: String,
      default: 'button',
    },
    to: {
    },
    href: {
    },
    block: {
      type: Boolean,
    },
    styleSize: {
      type: String,
      default: 'md',
    },
    disabled: {
      type: Boolean,
    },
    submitting: {
      type: Boolean,
    },
    outline: {
      type: Boolean,
    },
    styleType: {
      type: String,
      default: 'primary',
    },
  },
  computed: {
    className() {
      return [
        'btn',
        'btn-' + this.styleSize,
        'btn-' + this.styleType,
        this.outline ? 'outline': undefined,
        this.block ? 'btn-block' : undefined,
        this.submitting ? 'submitting' : undefined,
        this.disabled || this.submitting ? 'disabled' : undefined,
      ]
    }
  },
  methods: {
    onClick(e) {
      if (this.disabled || this.submitting) {
        e.preventDefault()
      }
    }
  }
}
</script>
