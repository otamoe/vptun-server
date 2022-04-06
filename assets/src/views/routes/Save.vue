<template>
    <section id="section-route-save" class="route-section">
      <h1>{{$t(route.params.route ? 'Route Update' : 'Route Create')}}</h1>
      <form action="post"  @submit.prevent="onSubmit">
        <FormGroup v-if="state.id">
          <FormLabel for="id" :block="true">{{$t('ID')}}</FormLabel>
          <FormInput id="id" name="id" :block="true" type="text" :disabled="true" v-model="state.id"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="remark" :block="true">{{$t('Remark')}}</FormLabel>
          <FormInput id="remark" name="remark" :block="true" type="text" v-model="state.remark"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="source-ip" :block="true">{{$t('Source IP')}}</FormLabel>
          <FormInput id="source-ip" name="source-ip" :block="true" type="text" v-model="state.sourceIP" :placeholder="$t('Netmask: 10.0.0.0/16')"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="source-port" :block="true">{{$t('Source port')}}</FormLabel>
          <FormInput id="source-port" name="source-port" :block="true" type="number" v-model.number="state.sourcePort"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="destination-ip" :block="true">{{$t('Destination IP')}}</FormLabel>
          <FormInput id="destination-ip" name="destination-ip" :block="true" type="text" v-model="state.destinationIP"  :placeholder="$t('Netmask: 10.0.0.0/16')"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="destination-port" :block="true">{{$t('Destination port')}}</FormLabel>
          <FormInput id="destination-port" name="destination-port" :block="true" type="number" v-model.number="state.destinationPort"></FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="type" :block="true">{{$t('Type')}}</FormLabel>
          <FormInput id="type" name="type" :block="true" type="select" v-model="state.type">
              <FormOption value="NONE">{{$t("NONE")}}</FormOption>
              <FormOption value="ICMP">{{$t("ICMP")}}</FormOption>
              <FormOption value="UDP">{{$t("UDP")}}</FormOption>
              <FormOption value="TCP">{{$t("TCP")}}</FormOption>
            </FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="action" :block="true">{{$t('Action')}}</FormLabel>
          <FormInput id="action" name="action" :block="true" type="select" v-model="state.action">
              <FormOption value="REJECT">{{$t("Reject")}}</FormOption>
              <FormOption value="ACCEPT">{{$t("Accept")}}</FormOption>
            </FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="state" :block="true">{{$t('State')}}</FormLabel>
          <FormInput id="state" name="state" :block="true" type="select" v-model="state.state">
              <FormOption value="AVAILABLE">{{$t("Available")}}</FormOption>
              <FormOption value="UNAVAILABLE">{{$t("Unavailable")}}</FormOption>
            </FormInput>
        </FormGroup>
        <FormGroup>
          <FormLabel for="level" :block="true">{{$t('Level')}}</FormLabel>
          <FormInput id="level" name="level" :block="true" type="number" v-model.number="state.level"></FormInput>
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
#section-route-save{
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
  name: "RouteSave",
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
      type: "NONE",
      sourceIP: "",
      sourcePort: 0,
      destinationIP: "",
      destinationPort: 0,
      remark: "",
      action: "ACCEPT",
      state: "AVAILABLE",
      level: 0,
      createdAt: "",
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
        switch (data.type || 0) {
          case 1:
            state.type = "ICMP"
            break;
          case 100:
            state.type = "UDP"
            break;
          case 101:
            state.type = "TCP"
            break;
          default:
            state.type = "NONE"
            break;
        }
        state.sourceIP = data.sourceIP || ""
        state.sourcePort = data.sourcePort || 0
        state.destinationIP = data.destinationIP || ""
        state.destinationPort = data.destinationPort || 0
        state.remark = data.remark || ""
        
        if (data.action) {
          state.action = "ACCEPT"
        } else {
          state.action = "REJECT"
        }
        
        if (data.state) {
          state.state = "UNAVAILABLE"
        } else {
          state.state = "AVAILABLE"
        }
        state.level = data.level || 0
        console.log(data.expiredAt)
        state.createdAt = dayjs(data.createdAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
        state.updatedAt = dayjs(data.updatedAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
        state.expiredAt = dayjs(data.expiredAt * 1000).format('YYYY-MM-DDTHH:mm:ss')
      } else {
        state.id = ""
        state.type = "NONE"
        state.sourceIP = ""
        state.sourcePort = 0
        state.destinationIP = ""
        state.destinationPort = 0
        state.remark = ""
        state.action = "ACCEPT"
        state.state = "AVAILABLE"
        state.level = 0
        state.createdAt = ""
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
      if (route.name !== "route/update" && route.name !== "route/create") {
        return
      }

      if (route.params.route === state.id) {
        return
      }

      if (!route.params.route) {
        stateUpdate(null)
        return
      }

      state.loading = true
      state.alert.value = undefined
      try { 
          let res = await axios.get(API_URL + "/api/route/" + route.params.route)
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
          let res = await axios.post(API_URL + "/api/route" + (route.params.route ? "/" + route.params.route : ''), {
            type: state.type,
            sourceIP: state.sourceIP,
            sourcePort: state.sourcePort,
            destinationIP: state.destinationIP,
            destinationPort: state.destinationPort,
            remark: state.remark,
            action: state.action,
            state: state.state,
            level: state.level,
            expiredAt: expiredAt,
          })
          if (res.data.error) {
            throw new Error(res.error)
          }
          let data = res.data.data
          stateUpdate(data)
          router.push({name: "route/update", params:{route: data.id}})
          state.alert = {value: route.params.route ? 'Updated'  : 'Created', type: 'success', close: false}
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