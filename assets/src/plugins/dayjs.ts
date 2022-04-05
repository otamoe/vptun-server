import { App } from 'vue'

import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import utc from 'dayjs/plugin/utc'
import timezone from 'dayjs/plugin/timezone'


declare module '@vue/runtime-core' {
  export interface ComponentCustomProperties {
    $dayjs: (date?: dayjs.ConfigType, format?: dayjs.OptionType, locale?: string, strict?: boolean) => dayjs.Dayjs;
  }
} 

export default {
  install(app: App) {
    dayjs.extend(relativeTime)
    dayjs.extend(utc)
    dayjs.extend(timezone)
    app.provide("dayjs", dayjs)
    app.config.globalProperties.$dayjs = dayjs
  }
}







