<template>
    <section id="section-client-read" class="client-shell-section">
      <h1>{{$t("Client Shell Read")}}</h1>
      <HighlightJson>{{state.data}}</HighlightJson>
      <h2>{{$t('Input')}}</h2>
      <pre>{{state.input}}</pre>
      <h2>{{$t('Output')}}</h2>
      <pre>{{state.output}}</pre>
    </section>
</template>
<style lang="scss">
#section-client-read {
  pre {
    outline: 1px solid #ccc;
    padding: .5em 1em;
  }
}
</style>
<script lang="ts">
import {inject, onMounted, reactive, watch} from "vue"
import {API_URL} from "@/config"
import HighlightJson from "@/components/highlight/Json.vue"
import { useRoute } from 'vue-router'
import queryString from 'query-string'

export default {
  name: "ClientShellRead",
  components:{
    HighlightJson,
  },
  


    setup() {
      const axios: any = inject('axios')
      const alertMessage: any = inject('alert-message')
      const state = reactive({
        loading: false,
        data: {} as any,
        input: "",
        output: "",
        delete: {
          submitting: false,
          client: undefined as any,
        },
      })
      const route = useRoute()



      

      watch(() => route.query, ontRead)
      async function ontRead() {
        if (route.name !== "client/shell/read") {
          return
        }
        state.loading = true
        try {
          let res = await axios.get(API_URL + "/api/client/" + route.params.client + "/shell/" + route.params.shell + "?"  +  queryString.stringify(route.query))
          if (res.data.error) {
            throw new Error(res.data.error)
          }
          state.data = res.data.data
          state.input = res.data.data.input || ""
          state.output = res.data.data.output || ""
          state.data.input = undefined
          state.data.output = undefined
        } catch (e) {
          alertMessage(e)
        } finally {
          state.loading = false
        }
      };
      
      onMounted(ontRead)
      
    return {
      ontRead,
      state,
      API_URL,
    }
  },



}
</script>