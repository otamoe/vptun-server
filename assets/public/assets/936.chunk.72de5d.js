/*! For license information please see 936.chunk.72de5d.js.LICENSE.txt */
"use strict";(self.webpackChunkvptun_server_assets=self.webpackChunkvptun_server_assets||[]).push([[936],{4342:(t,e,l)=>{l.r(e),l.d(e,{default:()=>d});var n=l(9739),a=l.n(n),s=l(4419),i=l.n(s)()(a());i.push([t.id,"#section-client-home .create-btn{font-size:1.4em;float:right;margin:.3em 1em}","",{version:3,sources:["webpack://./src/views/clients/Home.vue"],names:[],mappings:"AAEE,iCACE,eAAA,CACA,WAAA,CACA,eAAA",sourcesContent:["\n#section-client-home {\n  .create-btn {\n    font-size: 1.4em;\n    float: right;\n    margin: .3em 1em;\n  }\n}\n"],sourceRoot:""}]);const d=i},9439:(t,e,l)=>{var n=l(3379),a=l.n(n),s=l(7795),i=l.n(s),d=l(569),o=l.n(d),r=l(3565),u=l.n(r),c=l(9216),m=l.n(c),f=l(4589),_=l.n(f),w=l(4342),p={};p.styleTagTransform=_(),p.setAttributes=u(),p.insert=o().bind(null,"head"),p.domAPI=i(),p.insertStyleElement=m();var h=a()(w.default,p);if(!w.default.locals||t.hot.invalidate){var Y=!w.default.locals,A=Y?w:w.default.locals;t.hot.accept(4342,(e=>{w=l(4342),function(t,e,l){if(!t&&e||t&&!e)return!1;var n;for(n in t)if((!l||"default"!==n)&&t[n]!==e[n])return!1;for(n in e)if(!(l&&"default"===n||t[n]))return!1;return!0}(A,Y?w:w.default.locals,Y)?(A=Y?w:w.default.locals,h(w.default)):t.hot.invalidate()}))}t.hot.dispose((function(){h()}));w.default&&w.default.locals&&w.default.locals},7936:(t,e,l)=>{l.r(e),l.d(e,{default:()=>x});var n=l(3613),a=l(6739),s=l(1961);const i={id:"section-client-home",class:"client-section"},d={class:"td-id"},o={class:"td-hostname"},r={class:"td-remark"},u={class:"td-route-address"},c={class:"td-state"},m={class:"td-online"},f={class:"td-created-at"},_=["datetime","title"],w={class:"td-connect-at"},p=["datetime","title"],h={class:"td-updated-at"},Y=["datetime","title"],A={class:"td-expired-at"},$=["datetime","title"],z={class:"td-action"},D=(0,n.Uk)(" , "),y=(0,n.Uk)(" , "),g=(0,n.Uk)(" , "),v=["onClick"],M={key:0};var k=l(3097),b=l(6126),C=l(7724),H=l(9594),j=l(561);const U={name:"ClientHome",components:{DialogModal:b.Z},setup(){const t=(0,n.f3)("axios"),e=(0,n.f3)("dayjs"),l=(0,n.f3)("alert-message"),a=(0,k.qj)({loading:!1,data:[],nowUnix:e().unix(),delete:{submitting:!1,client:void 0}}),s=(0,H.yj)();async function i(){if("client/home"===s.name){a.loading=!0;try{let e=await t.get(C.T5+"/api/client?"+j.stringify(s.query));if(e.data.error)throw new Error(e.data.error);a.data=e.data.data}catch(t){l(t)}finally{a.loading=!1}}}return(0,n.YP)((()=>s.query),i),(0,n.bv)(i),{onList:i,state:a,onDelete:async function(){let e=a.delete.client;a.delete.submitting=!0;try{let l=await t.delete(C.T5+"/api/client/"+e.id);if(l.data.error)throw new Error(l.data.error);let n=[];for(const t in a.data){let l=a.data[t];l.id!==e.id&&n.push(l)}a.data=n,a.delete.client=void 0}catch(t){l(t?.response?.data?.error||t)}finally{a.delete.submitting=!1}},API_URL:C.T5}}};l(9439);const x=(0,l(4485).Z)(U,[["render",function(t,e,l,k,b,C){const H=(0,n.up)("RouterLink"),j=(0,n.up)("router-link"),U=(0,n.up)("DialogModal");return(0,n.wg)(),(0,n.iD)("section",i,[(0,n.Wm)(H,{class:"create-btn",to:{name:"client/create"}},{default:(0,n.w5)((()=>[(0,n.Uk)((0,a.zw)(t.$t("Create")),1)])),_:1}),(0,n._)("h1",null,(0,a.zw)(t.$t("Client List")),1),(0,n._)("table",null,[(0,n._)("thead",null,[(0,n._)("tr",null,[(0,n._)("th",null,(0,a.zw)(t.$t("ID")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Hostname")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Remark")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Route")),1),(0,n._)("th",null,(0,a.zw)(t.$t("State")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Online")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Created At")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Connect At")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Updated At")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Expired At")),1),(0,n._)("th",null,(0,a.zw)(t.$t("Action")),1)])]),(0,n._)("tbody",null,[((0,n.wg)(!0),(0,n.iD)(n.HY,null,(0,n.Ko)(k.state.data,((e,l)=>((0,n.wg)(),(0,n.iD)("tr",{key:l,class:(0,a.C_)({"client-online":e.online,"client-offline":!e.online,"client-shell":e.shell,"client-expired":k.state.nowUnix>e.expiredAt})},[(0,n._)("td",d,[(0,n._)("span",null,(0,a.zw)(e.id.substr(e.id.length-6)),1)]),(0,n._)("td",o,[(0,n._)("span",null,(0,a.zw)(e.hostname),1)]),(0,n._)("td",r,[(0,n._)("span",null,(0,a.zw)(e.remark),1)]),(0,n._)("td",u,[(0,n._)("span",null,(0,a.zw)(e.routeAddress),1)]),(0,n._)("td",c,[(0,n.Wm)(j,{to:{path:t.$route.path,query:{state:e.state||0}}},{default:(0,n.w5)((()=>[(0,n.Uk)((0,a.zw)(t.$t(e.state?"Unavailable":"Available")),1)])),_:2},1032,["to"])]),(0,n._)("td",m,[(0,n.Wm)(j,{to:{path:t.$route.path,query:{online:e.online||""}}},{default:(0,n.w5)((()=>[(0,n.Uk)((0,a.zw)(t.$t(e.online?"Online":"Offline")),1)])),_:2},1032,["to"])]),(0,n._)("td",f,[(0,n._)("time",{datetime:t.$dayjs(1e3*e.createdAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.createdAt).format("YYYY-MM-DD HH:mm:ss")},(0,a.zw)(t.$dayjs(1e3*e.createdAt).format("YYYY-MM-DD")),9,_)]),(0,n._)("td",w,[(0,n._)("time",{datetime:t.$dayjs(1e3*e.connectAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.connectAt).format("YYYY-MM-DD HH:mm:ss")},(0,a.zw)(t.$dayjs(1e3*e.connectAt).format("YYYY-MM-DD")),9,p)]),(0,n._)("td",h,[(0,n._)("time",{datetime:t.$dayjs(1e3*e.updatedAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.updatedAt).format("YYYY-MM-DD HH:mm:ss")},(0,a.zw)(t.$dayjs(1e3*e.updatedAt).format("YYYY-MM-DD")),9,Y)]),(0,n._)("td",A,[(0,n._)("time",{datetime:t.$dayjs(1e3*e.expiredAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.expiredAt).format("YYYY-MM-DD HH:mm:ss")},(0,a.zw)(t.$dayjs(1e3*e.expiredAt).format("YYYY-MM-DD")),9,$)]),(0,n._)("td",z,[(0,n.Wm)(H,{to:{name:"client/update",params:{client:e.id}}},{default:(0,n.w5)((()=>[(0,n.Uk)((0,a.zw)(t.$t("Edit")),1)])),_:2},1032,["to"]),D,(0,n.Wm)(H,{to:{name:"client/read",params:{client:e.id}}},{default:(0,n.w5)((()=>[(0,n.Uk)((0,a.zw)(t.$t("View")),1)])),_:2},1032,["to"]),y,(0,n.Wm)(H,{to:{name:"client/shell/home",params:{client:e.id}}},{default:(0,n.w5)((()=>[(0,n.Uk)((0,a.zw)(t.$t("Shell")),1)])),_:2},1032,["to"]),g,(0,n._)("a",{href:"#",onClick:(0,s.iM)((t=>k.state.delete.client=e),["prevent"])},(0,a.zw)(t.$t("Delete")),9,v)])],2)))),128))])]),(0,n.Wm)(U,{title:t.$t("Delete"),close:!k.state.delete.client,submitting:k.state.delete.submitting,onCancel:e[0]||(e[0]=t=>k.state.delete.client=void 0),onClose:e[1]||(e[1]=t=>k.state.delete.client=void 0),onConfirm:k.onDelete,class:"client-delete-modal"},{default:(0,n.w5)((()=>[k.state.delete.client?((0,n.wg)(),(0,n.iD)("div",M,[(0,n._)("p",null,[(0,n._)("strong",null,(0,a.zw)(t.$t("Hostname:")),1),(0,n._)("span",null,(0,a.zw)(k.state.delete.client.hostname),1)]),(0,n._)("p",null,[(0,n._)("strong",null,(0,a.zw)(t.$t("Remark:")),1),(0,n._)("span",null,(0,a.zw)(k.state.delete.client.remark),1)]),(0,n._)("p",null,[(0,n._)("strong",null,(0,a.zw)(t.$t("Route address:")),1),(0,n._)("span",null,(0,a.zw)(k.state.delete.client.routeAddress),1)]),(0,n._)("p",null,[(0,n._)("strong",null,(0,a.zw)(t.$t("State:")),1),(0,n._)("span",null,(0,a.zw)(t.$t(k.state.delete.client.state?"Unavailable":"Available")),1)]),(0,n._)("p",null,[(0,n._)("strong",null,(0,a.zw)(t.$t("Online:")),1),(0,n._)("span",null,(0,a.zw)(t.$t(k.state.delete.client.online?"Online":"Offline")),1)])])):(0,n.kq)("",!0),(0,n._)("h4",null,(0,a.zw)(t.$t("Are you sure you want to delete the client with the above information?")),1)])),_:1},8,["title","close","submitting","onConfirm"])])}]])}}]);
//# sourceMappingURL=936.chunk.72de5d.js.map