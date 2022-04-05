import 'es6-promise/auto'
import 'core-js/features/object/assign'
import '@/iconfont'
import App from '@/views/App.vue'
import { createApp, h } from 'vue'
import { createLocale } from '@/locale'
import { createRouter } from '@/router'
import dayjs from '@/plugins/dayjs'
import axios from 'axios'
import VueAxios from 'vue-axios'


export default async function init() {
    const locale = createLocale()
    const router = createRouter()    
    locale.load()

    
    const app = createApp({
        render() {
            return h(App)
        },
    })

    app.use(locale)
    app.use(router)
    if (process.env.NODE_ENV === 'development') {
        app.use(VueAxios, axios.create({
            headers: {
                'Authorization': 'Basic YWRtaW46YWRtaW4='
            },
        }))
    } else {
        app.use(VueAxios, axios.create({
            withCredentials: true,
        }))
    }

    app.use(dayjs)

    app.provide('axios', app.config.globalProperties.axios)

 
    
    // 设置当前 url 信息
    router.push(window.location.pathname+ window.location.search) 

    // 载入
    router.isReady().then(function() {
        app.mount('#app')
    })
}

init()
