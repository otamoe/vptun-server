<template>
  <div id="wrapper">
    <header id="header">
      <div class="container">
        <div id="header-nav-toggle">
            <button type="button" class="btn btn btn-primary">
                <svg class="icon icon-menu" aria-hidden="true">
                    <use xlink:href="#icon-menu"></use>
                </svg>
                <span>{{$t("Menu")}}</span>
            </button>
        </div>
        <component id="logo" :is="$route.to === '/' ? 'h1' : 'h2'">
          <RouterLink to="/" title="VPTun Server" rel="home" class="logo-link">
            <span class="logo-image">
                <svg class="icon icon-logo" aria-hidden="true">
                    <use xlink:href="#icon-logo"></use>
                </svg>
            </span>
            <span class="logo-text">
              <span class="logo-name">{{$t("VPTun")}}</span>
              <span class="logo-separator">-</span>
              <span class="logo-description">{{$t("VPTun Server")}}</span>
            </span>
          </RouterLink>
        </component>
        <nav id="header-nav">
          <ul>
            <li><RouterLink :to="{name: 'client/home'}">{{$t("Client List")}}</RouterLink></li>
            <li><RouterLink :to="{name: 'route/home'}">{{$t("Route List")}}</RouterLink></li>
          </ul>
        </nav>
      </div>
    </header>
    <RouterView></RouterView>
    <footer id="footer">
    </footer>
    <div class="popup-alert-message"  @dblclick.prevent="alertMessageClose">
        <AlertMessage :type="alertMessage.type"  :value="alertMessage.value" :close="alertMessage.close"></AlertMessage>
    </div>
  </div>
</template>
<style lang="scss">
@import "normalize.css";
@import "@/scss/variable.scss";

*, *::before, *::after {
    box-sizing: border-box;
    font-smoothing: antialiased;
    word-break: break-word;
    word-wrap: break-word;

}

body, html, #app, #wrapper {
    margin: 0;
    padding: 0;
    color: $body-color;
    background: $body-background;
    font-size: $body-font-size;
    line-height: $body-line-height;
    font-family: $body-font-family;
    font-weight: $body-font-weight;
}

img {
    border: none;
    outline: 0;
}

[tabindex]{
    outline: 0;
}
table {
	border-collapse: collapse;
	border-spacing: 0;
}
label {
    cursor: pointer;
    display: inline-block;
}
[type=button],
[type=reset],
[type=submit],
button {
    -webkit-appearance: button;
}
button, select {
    text-transform: none
}

a {
    background:transparent;
    color: $link-color;
    text-decoration: none;
    &:hover {
        text-decoration: underline;
        color: $link-hover-color;
    }
}

button, [type="button"], [type="submit"] {
    cursor: pointer;
}

.small, small {
    font-size: 80%;
}

input[type=search] {
    appearance: none;
}

h1 {
    font-size: 2em;
}
h2 {
    font-size: 1.5em;
}
h3 {
    font-size: 1.25em;
}
h4 {
    font-size: 1em;
}
h5 {
    font-size: .9em;
}
h6 {
    font-size: .8em;
}

#app, #wrapper {
    min-height: 100vh;
}

.ratio-item {
    width: 100%;
    height: 100%;
    position: absolute;
    overflow: hidden;
}

.ratio-21by9 {
    width: 100%;
    padding-bottom: 42.857143%;
    height: 0;
    position: relative;
}
.ratio-16by9 {
    width: 100%;
    padding-bottom: 56.25%;
    height: 0;
    position: relative;
}
.ratio-16by10 {
    width: 100%;
    padding-bottom: 62.5%;
    height: 0;
    position: relative;
}
.ratio-5by4 {
    width: 100%;
    padding-bottom: 80%;
    height: 0;
    position: relative;
}
.ratio-4by3 {
    width: 100%;
    padding-bottom: 75%;
    height: 0;
    position: relative;
}
.ratio-1by1 {
    width: 100%;
    padding-bottom: 100%;
    height: 0;
    position: relative;
}

.ratio-link {
    display: block;
    height: 100%;
    width: 100%;
    position: relative;
    overflow: hidden;
}

.ratio-img {
    position: absolute;
    top: 0;
    bottom: 0;
    left: 0;
    right: 0;
    width: 100%;
    margin: auto;
    height: auto; 
}



.ratio-media {
    display: block;
    width: 100%;
    height: 100%;   
}


.container {
    width: $container-width-md;
    margin: 0 auto;
}



.clearfix {
    @include clearfix;
} 

.content-none {
    @include content-none;
}

.text-none {
    @include text-none;
}

.icon {
    stroke-width: 0;
    stroke: currentColor;
    fill: currentColor;
    overflow: hidden;
    width: 1em;
    height: 1em;
    vertical-align: -0.15em;
    fill: currentColor;
    display: inline-block;
    color: inherit;
}


@include media-lt-md {
    .container {
        width: auto;
    }
}

@include media-lt-sm {
    body, html, #app, #wrapper {
        font-size: ($body-font-size * 0.875);
    }
}

@include media-lt-xxs {
    body, html, #app, #wrapper {
        font-size: ($body-font-size * 0.8125);
    }
}

@include media-gt-lg {
    .container {
        width: $container-width-lg;
    }
}



#header {
    font-size: 1.2em;
    padding: .25em 0;
    & .container {
        display: flex;
    }
}


#logo {
    & .logo-separator, 
    & .logo-description {
        @include content-none
    }
    & .logo-text{
        vertical-align: middle;
        display: block;
        line-height: 1;
    }
    & .logo-name {
        display: block;
        line-height: 1;
    }
    & .logo-image {
        vertical-align: middle;
        height: 1em;
        width: 1em;
        display: block;
    }
    & a {
        font-weight: normal;
        width: 100%;
        height: 100%;
        line-height: 1;
        display: flex;
        &:hover {
            text-decoration: none;
        }
        text-transform: uppercase;
    }
    font-size: 3em;
    margin: 0;
    padding: 0;
    height: 1em;
    width: auto;
    white-space: nowrap;
    flex-basis: 2%;
}

#header-nav-toggle{
    flex-basis: 48%;
    flex-grow: 1;
    height: 3em;
    padding: .3em 0;
    & button {
        display: none;
        background: #eee;
        border-radius: .25em;
        color: #444;
        padding: 0;
        margin: 0;
        border: 0;
        width: 2.4em;
        height: 2.4em;
        & span {
            @include content-none
        }
        & .icon{
            width: 80%;
            height: 80%;
        }
    }
}

#header-nav{
    flex-basis: 48%;
    flex-grow: 1;
    height: 3em;
    padding: .3em 0;
    & ul {
        padding: 0;
        margin: 0;
        list-style: none;
        float: right;
        height: 2.4em;
        & li {
            height: 100%;
            display: inline-block;
        }
        & a {
            display: inline-block;
            line-height: 2.4;
            padding: 0 .8em;
        }
    }
}


#main {
    min-height: 95vh;
    min-height: calc(100vh - 7em - 4em);
}

#footer {
    background: #444;
    color: #bbb;
    & a {
        color: #fff;
        &:hover {
            color: #bbb;
        }
    }
}
#footer-nav {
    & ul {
        margin: 0;
        padding: 0;
        display: flex;
        & li {
            text-align: center;
            list-style: none;
            flex-grow: 1;
            & a {
                padding: .8em;
                display: inline-block;
            }
        }
    }
    border-bottom: 1px solid rgb(200, 200, 200, .2);
}
#footer-copyright {
    text-align: center;
    padding: .1em 0;
}

@include media-lt-md {
    #header {
        margin-left: .25em;
        margin-right: .25em;
    }
    #header-nav-toggle {
        padding-left: .3em;
    }
}

.popup-alert-message {
    max-width: 20em;
    z-index: 1030;
    position: fixed;
    top: 40%;
    left: 50%;
    transform: translate(-50%, -50%);
    .alert-message {
        text-align: center;
        opacity: .9;
        border: 0;
        margin: 0;
        .btn {
            display: none;
        }
    }
}

@keyframes loading {
  from {
    transform: rotate(0deg);
  }
  to {
    transform: rotate(360deg);
  }
}
.loading, .load-more {
  margin: 1em 0;
  color: #2196f3;
  a {
    color: #2196f3;
  }
  text-align: center;
  svg {
    font-size: 2.6em;
  }
}
.loading  svg {
  animation: loading 1s linear infinite;
}
.breadcrumb{
  ol {
    display: flex;
    flex-wrap: wrap;
    padding: .75rem 1rem;
    margin-bottom: 1rem;
    list-style: none;
    background-color: #e9ecef;
    border-radius: .25rem;
    li {
      display: flex;
      max-width: 25em;
      overflow: hidden;
      white-space: nowrap;
      text-overflow: ellipsis;
      &.active {
        color: #6c757d;
      }
      span {
        display: inline-block;
        max-width: 100%;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
      }
      a {
        display: inline-block;
        max-width: 100%;
        overflow: hidden;
        white-space: nowrap;
        text-overflow: ellipsis;
        svg {
          vertical-align: -0.1em;
        }
        &:hover {
          text-decoration: none;
        }
      }
    }
    li+li {
      &:before {
        display: inline-block;
        padding-right: .5rem;
        color: #6c757d;
        content: "/";
      }
      padding-left: .5rem;
    }
  }
}
</style>
<script lang="ts">
import {provide, reactive} from "vue"
import AlertMessage from '@/components/AlertMessage.vue'

export default {
  components: { AlertMessage },
  name: "App",
  setup() {
    const alertMessage = reactive({
        type: 'error',
        value: [] as any,
        close: true,
        timeer: null as null | number
    })

    function alertMessageClose() {
        if (alertMessage.timeer !== null) {
            window.clearTimeout(alertMessage.timeer)
        }
        alertMessage.close = true
        alertMessage.timeer = null
    }

    provide("alert-message", function(value: any, type: string, timeout?: number) {
        alertMessage.type = type || 'error'
        alertMessage.value = value
        alertMessage.close = false
        if (alertMessage.timeer !== null) {
            clearTimeout(alertMessage.timeer)
        }
        alertMessage.timeer = window.setTimeout(function() {
            alertMessage.close = true
            alertMessage.timeer = null
        }, timeout || 4000)
    })
    provide("confirm", async function(message: string): Promise<boolean> {
        return false
    })


    
    return {
      alertMessage,
      alertMessageClose,
    };
  }
}
</script>