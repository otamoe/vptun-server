/*! For license information please see 760.chunk.c65814.js.LICENSE.txt */
"use strict";(self.webpackChunkvptun_server_assets=self.webpackChunkvptun_server_assets||[]).push([[760],{7894:(e,t,a)=>{a.r(t),a.d(t,{default:()=>r});var l=a(9739),d=a.n(l),o=a(4419),n=a.n(o)()(d());n.push([e.id,"#section-client-save form{padding:1em 2em}","",{version:3,sources:["webpack://./src/views/clients/Save.vue"],names:[],mappings:"AAEE,0BACE,eAAA",sourcesContent:["\n#section-client-save{\n  form {\n    padding: 1em 2em;\n  }\n}\n\n\n\n"],sourceRoot:""}]);const r=n},9217:(e,t,a)=>{var l=a(3379),d=a.n(l),o=a(7795),n=a.n(o),r=a(569),s=a.n(r),u=a(3565),m=a.n(u),i=a(9216),c=a.n(i),p=a(4589),f=a.n(p),k=a(7894),A={};A.styleTagTransform=f(),A.setAttributes=m(),A.insert=s().bind(null,"head"),A.domAPI=n(),A.insertStyleElement=c();var w=d()(k.default,A);if(!k.default.locals||e.hot.invalidate){var b=!k.default.locals,y=b?k:k.default.locals;e.hot.accept(7894,(t=>{k=a(7894),function(e,t,a){if(!e&&t||e&&!t)return!1;var l;for(l in e)if((!a||"default"!==l)&&e[l]!==t[l])return!1;for(l in t)if(!(a&&"default"===l||e[l]))return!1;return!0}(y,b?k:k.default.locals,b)?(y=b?k:k.default.locals,w(k.default)):e.hot.invalidate()}))}e.hot.dispose((function(){w()}));k.default&&k.default.locals&&k.default.locals},9730:(e,t,a)=>{a.d(t,{Z:()=>o});var l=a(3613);const d={name:"FormOption"};const o=(0,a(4485).Z)(d,[["render",function(e,t,a,d,o,n){return(0,l.wg)(),(0,l.iD)("option",null,[(0,l.WI)(e.$slots,"default")])}]])},8760:(e,t,a)=>{a.r(t),a.d(t,{default:()=>w});var l=a(3613),d=a(6739),o=a(1961);const n={id:"section-client-save",class:"client-section"};var r=a(3097),s=a(9594),u=a(7724),m=a(7329),i=a(7787),c=a(1384),p=a(8319),f=a(9730),k=a(762);const A={name:"ClientSave",components:{FormInput:i.Z,FormGroup:p.Z,FormLabel:c.Z,FormButton:k.Z,FormOption:f.Z,AlertMessage:m.Z},setup(){const e=(0,s.tv)(),t=(0,s.yj)();let a=(0,l.f3)("dayjs");const d=(0,l.f3)("axios");let o=(0,r.qj)({id:"",key:"",hostname:"",connectAddress:"",routeAddress:"",remark:"",state:"AVAILABLE",createdAt:"",connectAt:"",updatedAt:"",expiredAt:a(Date.UTC(9e3,0,1)).format("YYYY-MM-DDTHH:mm:ss"),loading:!1,submitting:!1,alert:{value:void 0,close:!0,type:"error"}});function n(e){e?(o.id=e.id,o.key=e.key||"",o.hostname=e.hostname||"",o.connectAddress=e.connectAddress||"",o.routeAddress=e.routeAddress||"",o.remark=e.remark||"",e.state?o.state="UNAVAILABLE":o.state="AVAILABLE",o.createdAt=a(1e3*e.createdAt).format("YYYY-MM-DDTHH:mm:ss"),o.connectAt=a(1e3*e.connectAt).format("YYYY-MM-DDTHH:mm:ss"),o.updatedAt=a(1e3*e.updatedAt).format("YYYY-MM-DDTHH:mm:ss"),o.expiredAt=a(1e3*e.expiredAt).format("YYYY-MM-DDTHH:mm:ss")):(o.id="",o.key="",o.hostname="",o.connectAddress="",o.routeAddress="",o.remark="",o.state="AVAILABLE",o.createdAt="",o.connectAt="",o.updatedAt="",o.expiredAt=a(Date.UTC(9e3,0,1)).format("YYYY-MM-DDTHH:mm:ss"),o.loading=!1,o.submitting=!1,o.alert={value:void 0,close:!0,type:"error"})}async function m(){if(("client/update"===t.name||"client/create"===t.name)&&t.params.client!==o.id)if(t.params.client){o.loading=!0,o.alert.value=void 0;try{let e=await d.get(u.T5+"/api/client/"+t.params.client);if(e.data.error)throw new Error(e.error);n(e.data.data)}catch(e){o.alert={value:e?.response?.data?.error||e,type:"error",close:!1}}finally{o.loading=!1}}else n(null)}return(0,l.bv)(m),(0,l.YP)((()=>t.params),m),{state:o,route:t,onLoad:m,onSubmit:async function(){o.submitting=!0;try{let l=0;if(""!==o.expiredAt){l=a(o.expiredAt).unix()}else l=Date.UTC(9e3,0,1)/1e3;let r=await d.post(u.T5+"/api/client"+(t.params.client?"/"+t.params.client:""),{key:o.key,routeAddress:o.routeAddress,remark:o.remark,state:o.state,expiredAt:l});if(r.data.error)throw new Error(r.error);let s=r.data.data;n(s),e.push({name:"client/update",params:{client:s.id}}),o.alert={value:t.params.client?"Updated":"Created",type:"success",close:!1}}catch(e){o.alert={value:e?.response?.data?.error||e,type:"error",close:!1}}finally{o.submitting=!1}}}}};a(9217);const w=(0,a(4485).Z)(A,[["render",function(e,t,a,r,s,u){const m=(0,l.up)("FormLabel"),i=(0,l.up)("FormInput"),c=(0,l.up)("FormGroup"),p=(0,l.up)("FormOption"),f=(0,l.up)("AlertMessage"),k=(0,l.up)("FormButton");return(0,l.wg)(),(0,l.iD)("section",n,[(0,l._)("h1",null,(0,d.zw)(e.$t(r.route.params.client?"Client Update":"Client Create")),1),(0,l._)("form",{action:"post",onSubmit:t[12]||(t[12]=(0,o.iM)(((...e)=>r.onSubmit&&r.onSubmit(...e)),["prevent"]))},[r.state.id?((0,l.wg)(),(0,l.j4)(c,{key:0},{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"id",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("ID")),1)])),_:1}),(0,l.Wm)(i,{id:"id",name:"id",block:!0,type:"text",disabled:!0,modelValue:r.state.id,"onUpdate:modelValue":t[0]||(t[0]=e=>r.state.id=e)},null,8,["modelValue"])])),_:1})):(0,l.kq)("",!0),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"key",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Key")),1)])),_:1}),(0,l.Wm)(i,{id:"key",name:"key",block:!0,type:"text",modelValue:r.state.key,"onUpdate:modelValue":t[1]||(t[1]=e=>r.state.key=e),placeholder:e.$t("Empty auto generate")},null,8,["modelValue","placeholder"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"hostname",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Hostname")),1)])),_:1}),(0,l.Wm)(i,{id:"hostname",name:"hostname",block:!0,type:"text",disabled:!0,modelValue:r.state.hostname,"onUpdate:modelValue":t[2]||(t[2]=e=>r.state.hostname=e)},null,8,["modelValue"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"connect-address",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Connect address")),1)])),_:1}),(0,l.Wm)(i,{id:"connect-address",name:"connect-address",block:!0,type:"text",disabled:!0,modelValue:r.state.connectAddress,"onUpdate:modelValue":t[3]||(t[3]=e=>r.state.connectAddress=e)},null,8,["modelValue"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"route-address",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Route address")),1)])),_:1}),(0,l.Wm)(i,{id:"route-address",name:"route-address",block:!0,type:"text",modelValue:r.state.routeAddress,"onUpdate:modelValue":t[4]||(t[4]=e=>r.state.routeAddress=e),placeholder:e.$t("Empty auto generate")},null,8,["modelValue","placeholder"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"remark",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Remark")),1)])),_:1}),(0,l.Wm)(i,{id:"remark",name:"remark",block:!0,type:"text",modelValue:r.state.remark,"onUpdate:modelValue":t[5]||(t[5]=e=>r.state.remark=e)},null,8,["modelValue"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"state",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("State")),1)])),_:1}),(0,l.Wm)(i,{id:"state",name:"state",block:!0,type:"select",modelValue:r.state.state,"onUpdate:modelValue":t[6]||(t[6]=e=>r.state.state=e)},{default:(0,l.w5)((()=>[(0,l.Wm)(p,{value:"AVAILABLE"},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Available")),1)])),_:1}),(0,l.Wm)(p,{value:"UNAVAILABLE"},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Unavailable")),1)])),_:1})])),_:1},8,["modelValue"])])),_:1}),r.state.id?((0,l.wg)(),(0,l.j4)(c,{key:1},{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"created-at",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Created at")),1)])),_:1}),(0,l.Wm)(i,{id:"created-at",name:"created_at",block:!0,type:"datetime-local",disabled:!0,modelValue:r.state.createdAt,"onUpdate:modelValue":t[7]||(t[7]=e=>r.state.createdAt=e)},null,8,["modelValue"])])),_:1})):(0,l.kq)("",!0),r.state.id?((0,l.wg)(),(0,l.j4)(c,{key:2},{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"connect-at",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Connect at")),1)])),_:1}),(0,l.Wm)(i,{id:"connect-at",name:"connect_at",block:!0,type:"datetime-local",disabled:!0,modelValue:r.state.connectAt,"onUpdate:modelValue":t[8]||(t[8]=e=>r.state.connectAt=e)},null,8,["modelValue"])])),_:1})):(0,l.kq)("",!0),r.state.id?((0,l.wg)(),(0,l.j4)(c,{key:3},{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"updated-at",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Updated at")),1)])),_:1}),(0,l.Wm)(i,{id:"updated-at",name:"updated_at",block:!0,type:"datetime-local",disabled:!0,modelValue:r.state.updatedAt,"onUpdate:modelValue":t[9]||(t[9]=e=>r.state.updatedAt=e)},null,8,["modelValue"])])),_:1})):(0,l.kq)("",!0),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(m,{for:"expired-at",block:!0},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Expired at")),1)])),_:1}),(0,l.Wm)(i,{id:"expired-at",name:"expired-at",block:!0,type:"datetime-local",modelValue:r.state.expiredAt,"onUpdate:modelValue":t[10]||(t[10]=e=>r.state.expiredAt=e)},null,8,["modelValue"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(f,{value:r.state.alert.value,type:r.state.alert.type,close:r.state.alert.close,onClose:t[11]||(t[11]=e=>r.state.alert.close=!0)},null,8,["value","type","close"])])),_:1}),(0,l.Wm)(c,null,{default:(0,l.w5)((()=>[(0,l.Wm)(k,{block:!0,styleSize:"lg",type:"submit",disabled:r.state.loading,submitting:r.state.submitting},{default:(0,l.w5)((()=>[(0,l.Uk)((0,d.zw)(e.$t("Submit")),1)])),_:1},8,["disabled","submitting"])])),_:1})],32)])}]])}}]);
//# sourceMappingURL=760.chunk.c65814.js.map