<template>
    <section id="section-client-save" class="client-section">
      <h1>{{$t(route.params.client ? 'Client Update' : 'Client Create')}}</h1>
      <form action="post"  @submit.prevent="onSubmit">
        <FormGroup v-if="state.id">
          <FormLabel for="id" :block="true">{{$t('ID')}}</FormLabel>
          <FormInput id="id" name="id" :block="true" type="text" :disabled="true" v-model="state.id"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="key" :block="true">{{$t('Key')}}</FormLabel>
          <FormInput id="key" name="key" :block="true" type="text" v-model="state.key" :placeholder="$t('Empty auto generate')"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="hostname" :block="true">{{$t('Hostname')}}</FormLabel>
          <FormInput id="hostname" name="hostname" :block="true" type="text" :disabled="true" v-model="state.hostname"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="connect-address" :block="true">{{$t('Connect address')}}</FormLabel>
          <FormInput id="connect-address" name="connect-address" :block="true" type="text" :disabled="true" v-model="state.connectAddress"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="route-address" :block="true">{{$t('Route address')}}</FormLabel>
          <FormInput id="route-address" name="route-address" :block="true" type="text" v-model="state.routeAddress" :placeholder="$t('Empty auto generate')"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="remark" :block="true">{{$t('Remark')}}</FormLabel>
          <FormInput id="remark" name="remark" :block="true" type="text"  v-model="state.remark"></FormInput>
        </FormGroup>
        
        <FormGroup>
          <FormLabel for="state" :block="true">{{$t('State')}}</FormLabel>
          <FormInput id="state" name="state" :block="true" type="select" v-model="state.state">
              <FormOption value="AVAILABLE">{{$t("Available")}}</FormOption>
              <FormOption value="UNAVAILABLE">{{$t("Unavailable")}}</FormOption>
            </FormInput>
        </FormGroup>
        <FormGroup v-if="state.id">
          <FormLabel for="created-at" :block="true">{{$t('Created at')}}</FormLabel>
          <FormInput id="created-at" name="created_at" :block="true" type="datetime-local" :disabled="true" v-model="state.createdAt"></FormInput>
        </FormGroup>
        <FormGroup v-if="state.id">
          <FormLabel for="connect-at" :block="true">{{$t('Connect at')}}</FormLabel>
          <FormInput id="connect-at" name="connect_at" :block="true" type="datetime-local" :disabled="true" v-model="state.connectAt"></FormInput>
        </FormGroup>
        <FormGroup v-if="state.id">
          <FormLabel for="updated-at" :block="true">{{$t('Updated at')}}</FormLabel>
          <FormInput id="updated-at" name="updated_at" :block="true" type="datetime-local" :disabled="true" v-model="state.updatedAt"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="expired-at" :block="true">{{$t('Expired at')}}</FormLabel>
          <FormInput id="expired-at" name="expired-at" :block="true" type="datetime-local"  v-model="state.expiredAt"></FormInput>
        </FormGroup>
        <FormGroup>
          <AlertMessage :value="state.alert.value" :type="state.alert.type" :close="state.alert.close" @close="state.alert.close = true"></AlertMessage>
        </FormGroup>
        <FormGroup>
          <FormButton :block="true" styleSize="lg" type="submit" :disabled="state.loading" :submitting="state.submitting">{{$t('Submit')}}</FormButton>
        </FormGroup>
      </form>
    </section>
</template>
<style lang="scss">
#section-client-save{
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
import FormOption from "@/components/forms/Option.vue"
import FormButton from "@/components/forms/Button.vue"


export default {
  name: "ClientSave",
  components: {
    FormInput,
    FormGroup,
    FormLabel,
    FormButton,
    FormOption,
    AlertMessage,
  },
  setup() {
    const router = useRouter()
    const route = useRoute()

    let dayjs: any = inject('dayjs')
    const axios: any = inject('axios')
    let state = reactive({
      id: "",
      key: "",
      hostname: "",
      connectAddress: "",
      routeAddress: "",
      remark: "",
      state: "AVAILABLE",
      createdAt: "",
      connectAt: "",
      updatedAt: "",
      expiredAt: dayjs(Date.UTC(9000, 0, 1)).format('YYYY-MM-DDTHH:mm:ss'),

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
        // www
        state.id = data.id
        state.key = data.key || ""
        state.hostname = data.hostname || ""
        state.connectAddress = data.connectAddress || ""
        state.routeAddress = data.routeAddress || ""
        state.remark = data.remark || ""
        if (data.state) {
          state.state = "UNAVAILABLE"
        } else {
          state.state = "AVAILABLE"
        }
        state.createdAt = dayjs(data.createdAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
        state.connectAt = dayjs(data.connectAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
        state.updatedAt = dayjs(data.updatedAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
        state.expiredAt = dayjs(data.expiredAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
      } else {
        state.id = ""
        state.key = ""
        state.hostname = ""
        state.connectAddress = ""
        state.routeAddress = ""
        state.remark = ""
        state.state = "AVAILABLE"
        state.createdAt = ""
        state.connectAt = ""
        state.updatedAt = ""
        state.expiredAt = dayjs(Date.UTC(9000, 0, 1)).format('YYYY-MM-DDTHH:mm:ss')

        state.loading = false
        state.submitting = false
        state.alert = {
          value: undefined as any,
          close: true,
          type: "error",
        }
      }
    }

    async function onLoad() {
      if (route.name !== "client/update" && route.name !== "client/create") {
        return
      }

      if (route.params.client === state.id) {
        return
      }

      if (!route.params.client) {
        stateUpdate(null)
        return
      }

      state.loading = true
      state.alert.value = undefined
      try { 
          let res = await axios.get(API_URL + "/api/client/" + route.params.client)
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
        let expiredAt = 0
        if (state.expiredAt !== "") {
          let t = dayjs(state.expiredAt)
          expiredAt = t.unix()
        } else {
          expiredAt = Date.UTC(9000, 0, 1) / 1000
        }
          let res = await axios.post(API_URL + "/api/client" + (route.params.client ? "/" + route.params.client : ''), {
            key: state.key,
            routeAddress: state.routeAddress,
            remark: state.remark,
            state: state.state,
            expiredAt: expiredAt,
          })
          if (res.data.error) {
            throw new Error(res.error)
          }
          let data = res.data.data
          stateUpdate(data)
          router.push({name: "client/update", params:{client: data.id}})
          state.alert = {value: route.params.client ? 'Updated'  : 'Created', type: 'success', close: false}
      } catch (e: any) {
          state.alert = {value: e?.response?.data?.error || e, type: 'error', close: false}
      } finally {
        state.submitting = false
      }
    }

    onMounted(onLoad)
    watch(() => route.params, onLoad)



    return {
      state,
      route,
      onLoad,
      onSubmit,
    }
  }
}
</script>