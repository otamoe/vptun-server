<template>
    <section id="section-client-shell-save" class="client-section">
      <h1 v-if="route.params.shell">{{$t("Shell $routeAddress - $hostname Update", {$routeAddress: state.client.routeAddress || "", $hostname: state.client.hostname || "Unknown"})}}</h1>
      <h1 v-else>{{$t("Shell $routeAddress - $hostname Create", {$routeAddress: state.client.routeAddress || "", $hostname: state.client.hostname || "Unknown"})}}</h1>
      <form action="post"  @submit.prevent="onSubmit">
        <FormGroup v-if="state.id">
          <FormLabel for="id" :block="true">{{$t('ID')}}</FormLabel>
          <FormInput id="id" name="id" :block="true" type="text" :disabled="true" v-model="state.id"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="client-id" :block="true">{{$t('Client ID')}}</FormLabel>
          <FormInput id="client-id" name="client-id" :block="true" type="text" :disabled="true" v-model="state.clientId"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="remark" :block="true">{{$t('Remark')}}</FormLabel>
          <FormInput id="remark" name="remark" :block="true" type="text" v-model="state.remark"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="input" :block="true">{{$t('Input')}}</FormLabel>
          <FormInput id="input" name="input" :block="true" rows="6" type="textarea" :disabled="state.id !== ''" v-model="state.input"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="output" :block="true">{{$t('Output')}}</FormLabel>
          <FormInput id="output" name="output" :block="true" rows="10" type="textarea" :disabled="state.id !== ''" v-model="state.output"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="status" :block="true">{{$t('Status')}}</FormLabel>
          <FormInput id="status" name="status" :block="true" type="text" v-model="state.status" :disabled="true"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="timeout" :block="true">{{$t('Timeout')}}</FormLabel>
          <FormInput id="timeout" name="timeout" :block="true" type="number" :disabled="state.id !== ''" v-model.number="state.timeout"></FormInput>
        </FormGroup>
        <FormGroup v-if="state.id">
          <FormLabel for="created-at" :block="true">{{$t('Created at')}}</FormLabel>
          <FormInput id="created-at" name="created_at" :block="true" type="datetime-local" :disabled="true" v-model="state.createdAt"></FormInput>
        </FormGroup>
        <FormGroup v-if="state.id">
          <FormLabel for="updated-at" :block="true">{{$t('Updated at')}}</FormLabel>
          <FormInput id="updated-at" name="updated_at" :block="true" type="datetime-local" :disabled="true" v-model="state.updatedAt"></FormInput>
        </FormGroup>
        <FormGroup>
          <AlertMessage :value="state.alert.value" :type="state.alert.type" :close="state.alert.close" @close="state.alert.close = true"></AlertMessage>
        </FormGroup>
        <FormGroup>
          <FormButton :block="true" styleSize="lg" type="submit" :disabled="state.loading || state.clientLoading" :submitting="state.submitting">{{$t('Submit')}}</FormButton>
        </FormGroup>
      </form>
    </section>
</template>
<style lang="scss">
#section-client-shell-save{
  form {
    padding: 1em 2em;
  }
}



</style>
<script lang="ts">
import {reactive, inject, onMounted, watch} from "vue"
import { useRouter, useRoute } from 'vue-router'
import {API_URL} from "@/config"

import AlertMessage from "@/components/AlertMessage.vue"
import FormInput from "@/components/forms/Input.vue"
import FormLabel from "@/components/forms/Label.vue"
import FormGroup from "@/components/forms/Group.vue"
import FormButton from "@/components/forms/Button.vue"


export default {
  name: "ClientSave",
  components: {
    FormInput,
    FormGroup,
    FormLabel,
    FormButton,
    AlertMessage,
  },
  setup() {
    const router = useRouter()
    const route = useRoute()
    const alertMessage: any = inject('alert-message')
    let dayjs: any = inject('dayjs')
    const axios: any = inject('axios')
    let state = reactive({
      id: "",
      clientId: "",
      remark: "",
      input: "",
      output: "",
      timeout: 600,
      status: -1,
      createdAt: "",
      updatedAt: "",


      clientLoading: false,
      client: {} as any,

      loading: false,
      submitting: false,

      alert: {
        value: undefined as any,
        close: true,
        type: "error",
      }
    })


    function stateUpdate(data : any) {
      if (data) {
        state.id = data.id
        state.clientId = data.clientId
        state.remark = data.remark || ""
        state.input = data.input || ""
        state.output = data.output || ""
        state.timeout = data.timeout || 600
        state.status = data.status || 0
        
        state.createdAt = dayjs(data.createdAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
        state.updatedAt = dayjs(data.updatedAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
      } else {
        state.id = ""
        state.clientId = route.params.client as string
        state.remark = ""
        state.input = ""
        state.output = ""
        state.timeout = 600
        state.status = -1
        state.createdAt = ""
        state.updatedAt = ""
        
        state.loading = false
        state.submitting = false
        state.alert = {
          value: undefined as any,
          close: true,
          type: "error",
        }
      }
    }


    async function onReadClient() {
        if (route.name !== "client/shell/update" && route.name !== "client/shell/create") {
          return
        }
        state.clientLoading = true
        try {
          let res = await axios.get(API_URL + "/api/client/" + route.params.client)
          if (res.data.error) {
            throw new Error(res.data.error)
          }
          state.client = res.data.data
        } catch (e) {
          alertMessage(e)
          state.client = {}
        } finally {
          state.clientLoading = false
        }
      }

    async function onLoad() {
      if (route.name !== "client/shell/update" && route.name !== "client/shell/create") {
        return
      }

      if (route.params.client === state.clientId && route.params.shell === state.id) {
        return
      }

      if (!route.params.shell) {
        stateUpdate(null)
        return
      }

      state.loading = true
      state.alert.value = undefined
      try { 
          let res = await axios.get(API_URL + "/api/client/" + route.params.client + "/shell/" + route.params.shell)
          if (res.data.error) {
            throw new Error(res.error)
          }
          let data = res.data.data
          stateUpdate(data)
      } catch (e: any) {
          state.alert = {value: e?.response?.data?.error || e, type: 'error', close: false}
      } finally {
        state.loading = false
      }
    }


    async function onSubmit() {
      state.submitting = true
      try {
          let post: any = {remark: state.remark}
          if (!route.params.shell) {
            post.input = state.input
            post.timeout = state.timeout
          }
          let res = await axios.post(API_URL + "/api/client/" + route.params.client  + "/shell" + (route.params.shell ? "/" + route.params.shell : ''), post)
          if (res.data.error) {
            throw new Error(res.error)
          }
          let data = res.data.data
          stateUpdate(data)
          router.push({name: "client/shell/update", params:{client: data.clientId, shell: data.id}})
          state.alert = {value: route.params.client ? 'Updated'  : 'Created', type: 'success', close: false}
      } catch (e: any) {
          state.alert = {value: e?.response?.data?.error || e, type: 'error', close: false}
      } finally {
        state.submitting = false
      }
    }

    onMounted(onLoad)
    onMounted(onReadClient)
    watch(() => route.params, onLoad)
    watch(() => route.params, onReadClient)



    return {
      state,
      route,
      onLoad,
      onSubmit,
    }
  }
}
</script>