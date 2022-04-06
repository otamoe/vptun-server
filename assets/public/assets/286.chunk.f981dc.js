/*! For license information please see 286.chunk.f981dc.js.LICENSE.txt */
"use strict";(self.webpackChunkvptun_server_assets=self.webpackChunkvptun_server_assets||[]).push([[286],{5352:(t,e,a)=>{a.r(e),a.d(e,{default:()=>s});var l=a(9739),n=a.n(l),o=a(4419),r=a.n(o)()(n());r.push([t.id,"#section-route-home .create-btn{font-size:1.4em;float:right;margin:.3em 1em}","",{version:3,sources:["webpack://./src/views/routes/Home.vue"],names:[],mappings:"AAEE,gCACE,eAAA,CACA,WAAA,CACA,eAAA",sourcesContent:["\n#section-route-home {\n  .create-btn {\n    font-size: 1.4em;\n    float: right;\n    margin: .3em 1em;\n  }\n}\n"],sourceRoot:""}]);const s=r},1786:(t,e,a)=>{var l=a(3379),n=a.n(l),o=a(7795),r=a.n(o),s=a(569),u=a.n(s),d=a(3565),i=a.n(d),c=a(9216),m=a.n(c),p=a(4589),_=a.n(p),w=a(5352),f={};f.styleTagTransform=_(),f.setAttributes=i(),f.insert=u().bind(null,"head"),f.domAPI=r(),f.insertStyleElement=m();var h=n()(w.default,f);if(!w.default.locals||t.hot.invalidate){var z=!w.default.locals,$=z?w:w.default.locals;t.hot.accept(5352,(e=>{w=a(5352),function(t,e,a){if(!t&&e||t&&!e)return!1;var l;for(l in t)if((!a||"default"!==l)&&t[l]!==e[l])return!1;for(l in e)if(!(a&&"default"===l||t[l]))return!1;return!0}($,z?w:w.default.locals,z)?($=z?w:w.default.locals,h(w.default)):t.hot.invalidate()}))}t.hot.dispose((function(){h()}));w.default&&w.default.locals&&w.default.locals},3286:(t,e,a)=>{a.r(e),a.d(e,{default:()=>j});var l=a(3613),n=a(6739),o=a(1961);const r={id:"section-route-home",class:"route-section"},s={class:"td-id"},u={class:"td-type"},d={class:"td-remark"},i={class:"td-source"},c={class:"td-destination"},m={class:"td-route-action"},p={class:"td-state"},_={class:"td-level"},w={class:"td-created-at"},f=["datetime","title"],h={class:"td-updated-at"},z=["datetime","title"],$={class:"td-expired-at"},y=["datetime","title"],A={class:"td-action"},Y=(0,l.Uk)(" , "),g=(0,l.Uk)(" , "),D=["onClick"],v={key:0};var k=a(3097),b=a(6126),M=a(7724),C=a(9594),H=a(561);const U={name:"RouteHome",components:{DialogModal:b.Z},setup(){const t=(0,l.f3)("axios"),e=(0,l.f3)("dayjs"),a=(0,l.f3)("alert-message"),n=(0,k.qj)({loading:!1,data:[],nowUnix:e().unix(),delete:{submitting:!1,route:void 0}}),o=(0,C.yj)();async function r(){if("route/home"===o.name){n.loading=!0;try{let e=await t.get(M.T5+"/api/route?"+H.stringify(o.query));if(e.data.error)throw new Error(e.data.error);n.data=e.data.data}catch(t){a(t)}finally{n.loading=!1}}}return(0,l.YP)((()=>o.query),r),(0,l.bv)(r),{onList:r,state:n,onDelete:async function(){let e=n.delete.route;n.delete.submitting=!0;try{let a=await t.delete(M.T5+"/api/route/"+e.id);if(a.data.error)throw new Error(a.data.error);let l=[];for(const t in n.data){let a=n.data[t];a.id!==e.id&&l.push(a)}n.data=l,n.delete.route=void 0}catch(t){a(t?.response?.data?.error||t)}finally{n.delete.submitting=!1}},typeString:function(t){switch(t||0){case 1:return"ICMP";case 100:return"UDP";case 101:return"TCP";default:return"NONE"}},API_URL:M.T5}}};a(1786);const j=(0,a(4485).Z)(U,[["render",function(t,e,a,k,b,M){const C=(0,l.up)("RouterLink"),H=(0,l.up)("router-link"),U=(0,l.up)("DialogModal");return(0,l.wg)(),(0,l.iD)("section",r,[(0,l.Wm)(C,{class:"create-btn",to:{name:"route/create"}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(t.$t("Create")),1)])),_:1}),(0,l._)("h1",null,(0,n.zw)(t.$t("Route List")),1),(0,l._)("table",null,[(0,l._)("thead",null,[(0,l._)("tr",null,[(0,l._)("th",null,(0,n.zw)(t.$t("ID")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Type")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Remark")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Source")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Destination")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Route action")),1),(0,l._)("th",null,(0,n.zw)(t.$t("State")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Level")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Created At")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Updated At")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Expired At")),1),(0,l._)("th",null,(0,n.zw)(t.$t("Action")),1)])]),(0,l._)("tbody",null,[((0,l.wg)(!0),(0,l.iD)(l.HY,null,(0,l.Ko)(k.state.data,((e,a)=>((0,l.wg)(),(0,l.iD)("tr",{key:a,class:(0,n.C_)({"route-expired":k.state.nowUnix>e.expiredAt})},[(0,l._)("td",s,[(0,l._)("span",null,(0,n.zw)(e.id.substr(e.id.length-6)),1)]),(0,l._)("td",u,[(0,l.Wm)(H,{to:{path:t.$route.path,query:{type:e.type||0}}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(k.typeString(e.type||0)),1)])),_:2},1032,["to"])]),(0,l._)("td",d,[(0,l._)("span",null,(0,n.zw)(e.remark),1)]),(0,l._)("td",i,[(0,l._)("span",null,(0,n.zw)(e.sourceIP+" - "+(e.sourcePort||0)),1)]),(0,l._)("td",c,[(0,l._)("span",null,(0,n.zw)(e.destinationIP+" - "+(e.destinationPort||0)),1)]),(0,l._)("td",m,[(0,l.Wm)(H,{to:{path:t.$route.path,query:{action:e.action||0}}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(t.$t(e.action?"Accept":"Reject")),1)])),_:2},1032,["to"])]),(0,l._)("td",p,[(0,l.Wm)(H,{to:{path:t.$route.path,query:{state:e.state||0}}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(t.$t(e.state?"Unavailable":"Available")),1)])),_:2},1032,["to"])]),(0,l._)("td",_,(0,n.zw)(e.level||0),1),(0,l._)("td",w,[(0,l._)("time",{datetime:t.$dayjs(1e3*e.createdAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.createdAt).format("YYYY-MM-DD HH:mm:ss")},(0,n.zw)(t.$dayjs(1e3*e.createdAt).format("YYYY-MM-DD")),9,f)]),(0,l._)("td",h,[(0,l._)("time",{datetime:t.$dayjs(1e3*e.updatedAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.updatedAt).format("YYYY-MM-DD HH:mm:ss")},(0,n.zw)(t.$dayjs(1e3*e.updatedAt).format("YYYY-MM-DD")),9,z)]),(0,l._)("td",$,[(0,l._)("time",{datetime:t.$dayjs(1e3*e.expiredAt).format("YYYY-MM-DD HH:mm:ss"),title:t.$dayjs(1e3*e.expiredAt).format("YYYY-MM-DD HH:mm:ss")},(0,n.zw)(t.$dayjs(1e3*e.expiredAt).format("YYYY-MM-DD")),9,y)]),(0,l._)("td",A,[(0,l.Wm)(C,{to:{name:"route/update",params:{route:e.id}}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(t.$t("Edit")),1)])),_:2},1032,["to"]),Y,(0,l.Wm)(C,{to:{name:"route/read",params:{route:e.id}}},{default:(0,l.w5)((()=>[(0,l.Uk)((0,n.zw)(t.$t("View")),1)])),_:2},1032,["to"]),g,(0,l._)("a",{href:"#",onClick:(0,o.iM)((t=>k.state.delete.route=e),["prevent"])},(0,n.zw)(t.$t("Delete")),9,D)])],2)))),128))])]),(0,l.Wm)(U,{title:t.$t("Delete"),close:!k.state.delete.route,submitting:k.state.delete.submitting,onCancel:e[0]||(e[0]=t=>k.state.delete.route=void 0),onClose:e[1]||(e[1]=t=>k.state.delete.route=void 0),onConfirm:k.onDelete,class:"route-delete-modal"},{default:(0,l.w5)((()=>[k.state.delete.route?((0,l.wg)(),(0,l.iD)("div",v,[(0,l._)("p",null,[(0,l._)("strong",null,(0,n.zw)(t.$t("Type:")),1),(0,l._)("span",null,(0,n.zw)(k.typeString(k.state.delete.route.type)),1)]),(0,l._)("p",null,[(0,l._)("strong",null,(0,n.zw)(t.$t("Remark:")),1),(0,l._)("span",null,(0,n.zw)(k.state.delete.route.remark),1)]),(0,l._)("p",null,[(0,l._)("strong",null,(0,n.zw)(t.$t("Source:")),1),(0,l._)("span",null,(0,n.zw)(k.state.delete.route.sourceIP+" - "+(k.state.delete.route.sourcePort||0)),1)]),(0,l._)("p",null,[(0,l._)("strong",null,(0,n.zw)(t.$t("Destination:")),1),(0,l._)("span",null,(0,n.zw)(k.state.delete.route.destinationIP+" - "+(k.state.delete.route.destinationPort||0)),1)]),(0,l._)("p",null,[(0,l._)("strong",null,(0,n.zw)(t.$t("State:")),1),(0,l._)("span",null,(0,n.zw)(t.$t(k.state.delete.route.state?"Unavailable":"Available")),1)]),(0,l._)("p",null,[(0,l._)("strong",null,(0,n.zw)(t.$t("Route action:")),1),(0,l._)("span",null,(0,n.zw)(t.$t(k.state.delete.route.action?"Accept":"Reject")),1)])])):(0,l.kq)("",!0),(0,l._)("h4",null,(0,n.zw)(t.$t("Are you sure you want to delete the route with the above information?")),1)])),_:1},8,["title","close","submitting","onConfirm"])])}]])}}]);
//# sourceMappingURL=286.chunk.f981dc.js.map