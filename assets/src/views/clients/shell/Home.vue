<template>
    <section id="section-client-shell-home" class="client-shell-section">
      <RouterLink class="create-btn" :to="{name: 'client/shell/create', params: {client: route.params.client}}">{{$t("Create")}}</RouterLink>
      <h1>{{$t("Shell $routeAddress - $hostname", {$routeAddress: state.client.routeAddress || "", $hostname: state.client.hostname || "Unknown"})}}</h1>
      <table>
          <thead>
            <tr>
              <th>{{$t("ID")}}</th>
              <th>{{$t("Remark")}}</th>
              <th>{{$t("Input")}}</th>
              <th>{{$t("Timeout")}}</th>
              <th>{{$t("Status")}}</th>
              <th>{{$t("Created At")}}</th>
              <th>{{$t("Updated At")}}</th>
              <th>{{$t("Action")}}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(shell, key) in state.data" :key="key" :class="{'shell-status-': (shell.status || 0)}">
              <td class="td-id"><span>{{shell.id.substr(shell.id.length-6)}}</span></td>
              <td class="td-input"><span>{{(shell.input && shell.input.length > 32 ? shell.input.substr(32) + "..." : shell.input)}}</span></td>
              <td class="td-remark"><span>{{shell.remark}}</span></td>
              <td class="td-timeout"><span>{{shell.timeout}}s</span></td>
              <td class="td-status"><span>{{shell.status || 0}}</span></td>
              <td class="td-created-at"><time :datetime="$dayjs(shell.createdAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(shell.createdAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(shell.createdAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-updated-at"><time :datetime="$dayjs(shell.updatedAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(shell.updatedAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(shell.updatedAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-action">
                <RouterLink :to="{name: 'client/shell/update', params: {client: shell.clientId, shell: shell.id}}">{{$t("Edit")}}</RouterLink> , 
                <RouterLink :to="{name: 'client/shell/read', params: {client: shell.clientId, shell: shell.id}}">{{$t("View")}}</RouterLink> , 
                <a  href="#" @click.prevent="state.delete.shell = shell">{{$t("Delete")}}</a>
              </td>
            </tr>
          </tbody>
      </table>
      <DialogModal
        :title="$t('Delete')"
        :close="!state.delete.shell"
        :submitting="state.delete.submitting"
        @cancel="state.delete.shell = undefined"
        @close="state.delete.shell = undefined"
        @confirm="onDelete"
        class="shell-delete-modal"
      >
        <div v-if="state.delete.shell">
          <p><strong>{{$t("Remark:")}} </strong><span>{{state.delete.shell.remark}}</span></p>
          <p><strong>{{$t("Input:")}} </strong><span>{{state.delete.shell.input}}</span></p>
        </div>
        <h4>
          {{$t("Are you sure you want to delete the shell with the above information?")}}
        </h4>
      </DialogModal>
    </section>
</template>
<style lang="scss">
#section-client-shell-home {
  .create-btn {
    font-size: 1.4em;
    float: right;
    margin: .3em 1em;
  }
}
</style>
<script lang="ts">
import {inject, onMounted, reactive, watch} from "vue"
import {API_URL} from "@/config"
import { useRoute } from 'vue-router'
import queryString from 'query-string'
import DialogModal from "@/components/DialogModal.vue"

export default {
  name: "ClientShellList",
  components:{
    DialogModal,
  },

    setup() {
      const axios: any = inject('axios')
      const alertMessage: any = inject('alert-message')
      const state = reactive({
        loading: false,
        clientLoading: false,
        client: {} as any,
        data: {} as any,
        delete: {
          shell: undefined as any,
          submitting: false,
        },
      })
      const route = useRoute()

  
    

      watch(() => route.query, onList)
      
      async function onReadClient() {
        if (route.name !== "client/shell/home") {
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
        } finally {
          state.clientLoading = false
        }
      }

      async function onList() {
        if (route.name !== "client/shell/home") {
          return
        }
        state.loading = true
        try {
          let res = await axios.get(API_URL + "/api/client/" + route.params.client + "/shell?limit=1000&"  +  queryString.stringify(route.query))
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

        async function onDelete() {
        let shell = state.delete.shell
        state.delete.submitting = true
        try {
          let res = await axios.delete(API_URL + "/api/client/" + shell.clientId + "/shell/" + shell.id)
          if (res.data.error) {
            throw new Error(res.data.error)
          }
          let newData = []
          for (const key in state.data) {
            let item = state.data[key]
            if (item.id !== shell.id) {
              newData.push(item)
            }
          }
          state.data = newData
          state.delete.shell = undefined
        } catch (e: any) {
            alertMessage(e?.response?.data?.error || e)
        } finally {
          state.delete.submitting = false
        }
      }

      
      onMounted(onList)
      onMounted(onReadClient)
    return {
      onList,
      onReadClient,
      onDelete,
      state,
      route,
      API_URL,
    }
  },



}
</script>