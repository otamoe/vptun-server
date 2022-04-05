<template>
<pre class="highlight highlight-json" v-html="syntaxHighlight(state.data)"></pre>
</template>
<style lang="scss">
.highlight-json {
  outline: 1px solid #ccc; 
  padding: .5em 1em;
  .string { color: green; }
  .number { color: darkorange; }
  .boolean { color: blue; }
  .null { color: magenta; }
  .key { color: red; }
}

</style>
<script lang="ts">
import {onBeforeMount, onMounted, watch, reactive} from "vue"

export default {
    setup(_:any, {slots}: any) {
    //   function syntaxHighlight(obj: any) {
    //           let text = ''
    //         if (typeof vnode.children === 'string') {
    //         text = vnode.children
    //         } else {
    //         for (let i = 0; i < vnode.children.length; i++) {
    //             text += vnode.children[i].text || vnode.children[i].children || ''
    //         }
    //         }

    //   }

    const state = reactive({
        data: "",
    })
    function ontUpdate() {
        let data = ""
        let slotsDefault =  slots.default()
        for (let index = 0; index < slotsDefault.length; index++) {
            const item = slotsDefault[index];
            data += item.children
        }
        state.data = data
    }
    onBeforeMount(ontUpdate)
    onMounted(ontUpdate)
    watch(slots.default, ontUpdate)


      
      function syntaxHighlight(obj: any) {
        let json: string
        if (obj && (typeof obj === "object")) {
          json = JSON.stringify(obj, undefined, 2);
        } else {
          json = obj || {}
        }

        json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
        return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
              let cls = 'number';
              if (/^"/.test(match)) {
                  if (/:$/.test(match)) {
                      cls = 'key';
                  } else {
                      cls = 'string';
                  }
              } else if (/true|false/.test(match)) {
                  cls = 'boolean';
              } else if (/null/.test(match)) {
                  cls = 'null';
              }
              return '<span class="' + cls + '">' + match + '</span>';
          });
      }

      return {
          state,
          syntaxHighlight,
          ontUpdate,
      }
    }
}
</script>