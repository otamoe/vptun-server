/*! For license information please see 300.chunk.968c76.js.LICENSE.txt */
"use strict";(self.webpackChunkvptun_server_assets=self.webpackChunkvptun_server_assets||[]).push([[300],{9300:(t,e,a)=>{a.r(e),a.d(e,{default:()=>f});var r=a(3613),n=a(6739);const s={id:"section-route-read",class:"route-section"};var o=a(3097),u=a(7724),i=a(5464),d=a(9594),l=a(561);const c={name:"RouteRead",components:{HighlightJson:i.Z},setup(){const t=(0,r.f3)("axios"),e=(0,r.f3)("alert-message"),a=(0,o.qj)({loading:!1,data:{},delete:{submitting:!1,route:void 0}}),n=(0,d.yj)();async function s(){if("route/read"===n.name){a.loading=!0;try{let e=await t.get(u.T5+"/api/route/"+n.params.route+"?"+l.stringify(n.query));if(e.data.error)throw new Error(e.data.error);a.data=e.data.data}catch(t){e(t)}finally{a.loading=!1}}}return(0,r.YP)((()=>n.query),s),(0,r.bv)(s),{ontRead:s,state:a,API_URL:u.T5}}};const f=(0,a(4485).Z)(c,[["render",function(t,e,a,o,u,i){const d=(0,r.up)("HighlightJson");return(0,r.wg)(),(0,r.iD)("section",s,[(0,r._)("h1",null,(0,n.zw)(t.$t("Route Read")),1),(0,r.Wm)(d,null,{default:(0,r.w5)((()=>[(0,r.Uk)((0,n.zw)(o.state.data),1)])),_:1})])}]])}}]);
//# sourceMappingURL=300.chunk.968c76.js.map