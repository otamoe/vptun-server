<template>
    <section id="section-client-read" class="client-section">
      <h1>{{$t("Client Read")}}</h1>
      <HighlightJson>{{state.data}}</HighlightJson>
    </section>
</template>

<script lang="ts">
import {inject, onMounted, reactive, watch} from "vue"
import {API_URL} from "@/config"
import HighlightJson from "@/components/highlight/Json.vue"
import { useRoute } from 'vue-router'
import queryString from 'query-string'

export default {
  name: "ClientRead",
  components:{
    HighlightJson,
  },
  


    setup() {
      const axios: any = inject('axios')
      const alertMessage: any = inject('alert-message')
      const state = reactive({
        loading: false,
        data: {} as any,
        delete: {
          submitting: false,
          client: undefined as any,
        },
      })
      const route = useRoute()



      

      watch(() => route.query, ontRead)
      async function ontRead() {
        if (route.name !== "client/read") {
          return
        }
        state.loading = true
        try {
          let res = await axios.get(API_URL + "/api/client/" + route.params.client + "?"  +  queryString.stringify(route.query))
          if (res.data.error) {
            throw new Error(res.data.error)
          }
          state.data = res.data.data
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