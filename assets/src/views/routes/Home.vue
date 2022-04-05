<template>
    <section id="section-route-home" class="route-section">
      <RouterLink class="create-btn" :to="{name: 'route/create'}">{{$t("Create")}}</RouterLink>
      <h1>{{$t("Route List")}}</h1>
      <table>
          <thead>
            <tr>
              <th>{{$t("ID")}}</th>
              <th>{{$t("Type")}}</th>
              <th>{{$t("Remark")}}</th>
              <th>{{$t("Source")}}</th>
              <th>{{$t("Destination")}}</th>
              <th>{{$t("Route action")}}</th>
              <th>{{$t("State")}}</th>
              <th>{{$t("Level")}}</th>
              <th>{{$t("Created At")}}</th>
              <th>{{$t("Updated At")}}</th>
              <th>{{$t("Expired At")}}</th>
              <th>{{$t("Action")}}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(route, key) in state.data" :key="key" :class="{'route-expired': state.nowUnix > route.expiredAt }">
              <td class="td-id"><span>{{route.id.substr(route.id.length-6)}}</span></td>
              <td class="td-type"><router-link :to="{path: $route.path, query: {type: route.type || 0}}">{{typeString(route.type || 0)}}</router-link></td>
              <td class="td-remark"><span>{{route.remark}}</span></td>
              <td class="td-source"><span>{{route.sourceIP + " - " + (route.sourcePort || 0)}}</span></td>
              <td class="td-destination"><span>{{route.destinationIP + " - " + (route.destinationPort || 0)}}</span></td>
              <td class="td-route-action"><router-link :to="{path: $route.path, query: {action: route.action || 0}}">{{$t(route.action ? "Accept" : "Reject")}}</router-link></td>
              <td class="td-state"><router-link :to="{path: $route.path, query: {state: route.state || 0}}">{{$t(route.state ? "Unavailable" : "Available")}}</router-link></td>
              <td class="td-level">{{route.level || 0}}</td>
              <td class="td-created-at"><time :datetime="$dayjs(route.createdAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(route.createdAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(route.createdAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-updated-at"><time :datetime="$dayjs(route.updatedAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(route.updatedAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(route.updatedAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-expired-at"><time :datetime="$dayjs(route.expiredAt * 1000).format('YYYY-MM-DD HH:mm:ss')" :title="$dayjs(route.expiredAt * 1000).format('YYYY-MM-DD HH:mm:ss')">{{$dayjs(route.expiredAt * 1000).format("YYYY-MM-DD")}}</time></td>
              <td class="td-action">
                <RouterLink :to="{name: 'route/update', params: {route: route.id}}">{{$t("Edit")}}</RouterLink> , 
                <RouterLink :to="{name: 'route/read', params: {route: route.id}}">{{$t("View")}}</RouterLink> , 
                <a  href="#" @click.prevent="state.delete.route = route">{{$t("Delete")}}</a>
              </td>
            </tr>
          </tbody>
      </table>
      <DialogModal
      :title="$t('Delete')"
      :close="!state.delete.route"
      :submitting="state.delete.submitting"
      @cancel="state.delete.route = undefined"
      @close="state.delete.route = undefined"
      @confirm="onDelete"
      class="route-delete-modal"
    >
      <div v-if="state.delete.route">
        <p><strong>{{$t("Type:")}}</strong><span>{{typeString(state.delete.route.type)}}</span></p>
        <p><strong>{{$t("Remark:")}} </strong><span>{{state.delete.route.remark}}</span></p>
        <p><strong>{{$t("Source:")}}</strong><span>{{state.delete.route.sourceIP + " - " + (state.delete.route.sourcePort || 0)}}</span></p>
        <p><strong>{{$t("Destination:")}}</strong><span>{{state.delete.route.destinationIP + " - " + (state.delete.route.destinationPort || 0)}}</span></p>
        <p><strong>{{$t("State:")}}</strong><span>{{$t(state.delete.route.state ? "Unavailable" : "Available")}}</span></p>
        <p><strong>{{$t("Route action:")}}</strong><span>{{$t(state.delete.route.action ? "Accept" : "Reject")}}</span></p>
      </div>
      <h4>
        {{$t("Are you sure you want to delete the route with the above information?")}}
      </h4>
    </DialogModal>
    </section>
</template>
<style lang="scss">
#section-route-home {
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
  name: "RouteHome",
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
          route: undefined as any,
        },
      })
      const route = useRoute()




      function typeString(t: number) {
        switch (t || 0) {
          case 1:
            return "ICMP"
          case 100:
            return "UDP"
          case 101:
            return "TCP"
          default:
            return "NONE"
        }
      }

      async function onDelete() {
        let route = state.delete.route
        state.delete.submitting = true
        try {
          let res = await axios.delete(API_URL + "/api/route/" + route.id)
          if (res.data.error) {
            throw new Error(res.data.error)
          }
          let newData = []
          for (const key in state.data) {
            let item = state.data[key]
            if (item.id !== route.id) {
              newData.push(item)
            }
          }
          state.data = newData
          state.delete.route = undefined
        } catch (e: any) {
            alertMessage(e?.response?.data?.error || e)
        } finally {
          state.delete.submitting = false
        }
      }

      watch(() => route.query, onList)

      async function onList() {
        if (route.name !== "route/home") {
          return
        }
        state.loading = true
        try {
          let res = await axios.get(API_URL + "/api/route?" +  queryString.stringify(route.query))
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
      typeString,
      API_URL,
    }
  },
}
</script>