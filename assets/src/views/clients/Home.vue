<template>
    <section id="section-client-home" class="client-section">
      <RouterLink class="create-btn" :to="{name: 'client/create'}">{{$t("Create")}}</RouterLink>
      <h1>{{$t("Client List")}}</h1>
      <table>
          <thead>
            <tr>
              <th>{{$t("ID")}}</th>
              <th>{{$t("Hostname")}}</th>
              <th>{{$t("Remark")}}</th>
              <th>{{$t("Route")}}</th>
              <th>{{$t("State")}}</th>
              <th>{{$t("Online")}}</th>
              <th>{{$t("Created At")}}</th>
              <th>{{$t("Connect At")}}</th>
              <th>{{$t("Updated At")}}</th>
              <th>{{$t("Expired At")}}</th>
              <th>{{$t("Action")}}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(client, key) in state.data" :key="key" :class="{'client-online': client.online, 'client-offline': !client.online,'client-shell': client.shell, 'client-expired': state.nowUnix > client.expiredAt }">
              <td class="td-id"><span>{{client.id.substr(client.id.length-6)}}</span></td>
              <td class="td-hostname"><span>{{client.hostname}}</span></td>
              <td class="td-remark"><span>{{client.remark}}</span></td>
              <td class="td-route-address"><span>{{client.routeAddress}}</span></td>
              <td class="td-state"><router-link :to="{path: $route.path, query: {state: client.state || 0}}">{{$t(client.state ? "Unavailable" : "Available")}}</router-link></td>
              <td class="td-online"><router-link :to="{path: $route.path, query: {online: client.online || ''}}">{{$t(client.online ? "Online" : "Offline")}}</router-link></td>
              <td class="td-created-at"><time :datetime="$dayjs(client.createdAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(client.createdAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(client.createdAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-connect-at"><time :datetime="$dayjs(client.connectAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(client.connectAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(client.connectAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-updated-at"><time :datetime="$dayjs(client.updatedAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(client.updatedAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(client.updatedAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-expired-at"><time :datetime="$dayjs(client.expiredAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(client.expiredAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(client.expiredAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-action">
                <RouterLink :to="{name: 'client/update', params: {client: client.id}}">{{$t("Edit")}}</RouterLink> , 
                <RouterLink :to="{name: 'client/read', params: {client: client.id}}">{{$t("View")}}</RouterLink> , 
                <RouterLink :to="{name: 'client/shell/home', params: {client: client.id}}">{{$t("Shell")}}</RouterLink> , 
                <a  href="#" @click.prevent="state.delete.client = client">{{$t("Delete")}}</a>
              </td>
            </tr>
          </tbody>
      </table>
      <DialogModal
      :title="$t('Delete')"
      :close="!state.delete.client"
      :submitting="state.delete.submitting"
      @cancel="state.delete.client = undefined"
      @close="state.delete.client = undefined"
      @confirm="onDelete"
      class="client-delete-modal"
    >
      <div v-if="state.delete.client">
        <p><strong>{{$t("Hostname:")}}</strong><span>{{state.delete.client.hostname}}</span></p>
        <p><strong>{{$t("Remark:")}} </strong><span>{{state.delete.client.remark}}</span></p>
        <p><strong>{{$t("Route address:")}}</strong><span>{{state.delete.client.routeAddress}}</span></p>
        <p><strong>{{$t("State:")}}</strong><span>{{$t(state.delete.client.state ? "Unavailable" : "Available")}}</span></p>
        <p><strong>{{$t("Online:")}}</strong><span>{{$t(state.delete.client.online ? "Online" : "Offline")}}</span></p>
      </div>
      <h4>
        {{$t("Are you sure you want to delete the client with the above information?")}}
      </h4>
    </DialogModal>
    </section>
</template>
<style lang="scss">
#section-client-home {
  .create-btn {
    font-size: 1.4em;
    float: right;
    margin: .3em 1em;
  }
}
</style>
<script lang="ts">
import {inject, onMounted, reactive, watch} from "vue"
import DialogModal from "@/components/DialogModal.vue"
import {API_URL} from "@/config"
import { useRoute } from 'vue-router'
import queryString from 'query-string'
export default {
  name: "ClientHome",
  components:{
    DialogModal,
  },
  
    setup() {
      const axios: any = inject('axios')
      const dayjs: any = inject('dayjs')
      const alertMessage: any = inject('alert-message')
      const state = reactive({
        loading: false,
        data:[] as any[],
        nowUnix: dayjs().unix(),
        delete: {
          submitting: false,
          client: undefined as any,
        },
      })
      const route = useRoute()




      async function onDelete() {
        let client = state.delete.client
        state.delete.submitting = true
        try {
          let res = await axios.delete(API_URL + "/api/client/" + client.id)
          if (res.data.error) {
            throw new Error(res.data.error)
          }
          let newData = []
          for (const key in state.data) {
            let item = state.data[key]
            if (item.id !== client.id) {
              newData.push(item)
            }
          }
          state.data = newData
          state.delete.client = undefined
        } catch (e: any) {
            alertMessage(e?.response?.data?.error || e)
        } finally {
          state.delete.submitting = false
        }
      }

      watch(() => route.query, onList)

      async function onList() {
        if (route.name !== "client/home") {
          return
        }
        state.loading = true
        try {
          let res = await axios.get(API_URL + "/api/client?" +  queryString.stringify(route.query))
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
      
      onMounted(onList)
      
    return {
      onList,
      state,
      onDelete,
      API_URL,
    }
  },
}
</script>