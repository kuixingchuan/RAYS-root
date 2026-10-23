<template>
 <main class="rays-console" v-loading="busy">
  <header><div><h1>RAYS 管理平台</h1><p>独立账号 · 会话安全 · 服务维护</p></div><el-button @click="refresh">刷新</el-button></header>
  <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" />
  <el-alert title="个人照片和图库数据保存在用户设备内；此平台仅管理账号及服务端运行数据。" type="info" :closable="false" />
  <el-tabs v-model="tab" @tab-change="refresh">
   <el-tab-pane label="运行概览" name="overview">
    <el-row :gutter="16"><el-col v-for="item in stats" :key="item.label" :xs="24" :sm="12" :lg="6"><el-card><p>{{item.label}}</p><strong>{{item.value ?? '—'}}</strong></el-card></el-col></el-row>
    <el-descriptions title="服务状态" border :column="1"><el-descriptions-item label="API">{{overview.api_status === 'RUNNING' ? '运行中' : '—'}}</el-descriptions-item><el-descriptions-item label="服务运行时间">{{overview.uptime_seconds == null ? '—' : overview.uptime_seconds+' 秒'}}</el-descriptions-item><el-descriptions-item label="登录保护">{{overview.login_rate_limit == null ? '—' : '每身份/IP 在 '+overview.rate_limit_window_seconds+' 秒内最多 '+overview.login_rate_limit+' 次'}}</el-descriptions-item><el-descriptions-item label="数据库">{{overview.checked_at ? '连接正常' : '尚未取得状态'}}</el-descriptions-item><el-descriptions-item label="数据库大小">{{ overview.database_bytes == null ? '—' : (Number(overview.database_bytes)/1048576).toFixed(2)+' MB' }}</el-descriptions-item><el-descriptions-item label="检查时间">{{overview.checked_at || '—'}}</el-descriptions-item></el-descriptions>
   </el-tab-pane>
   <el-tab-pane label="用户管理" name="accounts">
    <div class="toolbar"><el-input v-model="query" placeholder="昵称或登录身份" clearable @keyup.enter="refresh"/><el-select v-model="status" clearable placeholder="全部状态"><el-option label="待登录" value="INVITED"/><el-option label="已激活" value="ACTIVE"/><el-option label="已禁用" value="DISABLED"/></el-select><el-button @click="refresh">查询</el-button><el-button type="primary" @click="openCreate">创建账号</el-button></div>
    <el-table :data="accounts" row-key="id" empty-text="没有匹配账号" style="width:100%"><el-table-column prop="displayName" label="昵称" min-width="120"/><el-table-column prop="subjectMasked" label="登录身份" min-width="170"/><el-table-column prop="provider" label="渠道" width="100"/><el-table-column label="状态" width="110"><template #default="{row}"><el-tag :type="row.status==='DISABLED'?'danger':'success'">{{labels[row.status] || row.status}}</el-tag></template></el-table-column><el-table-column label="操作" min-width="330"><template #default="{row}"><el-button link @click="openEdit(row)">编辑</el-button><el-button link @click="toggle(row)">{{row.status==='DISABLED'?'解禁':'禁用'}}</el-button><el-button link @click="revoke(row)">退出全部会话</el-button><el-button link type="danger" @click="remove(row)">删除</el-button></template></el-table-column></el-table>
   </el-tab-pane>
   <el-tab-pane label="操作审计" name="audit"><p>最近 100 条管理操作，仅记录操作者、动作和目标账号，不记录照片、验证码或凭据。</p><el-table :data="audit"><el-table-column prop="created_at" label="时间" min-width="200"/><el-table-column prop="actor" label="操作者"/><el-table-column prop="action" label="动作"/><el-table-column prop="target" label="目标" min-width="260"/></el-table></el-tab-pane>
   <el-tab-pane label="数据维护" name="maintenance">
    <el-card><h2>过期会话维护</h2><p>将已到期会话标记为过期，不删除账号、业务记录或个人设备文件。</p><el-button @click="expire">整理过期会话</el-button></el-card>
    <el-card><h2>数据库备份与恢复</h2><p>由服务器运维入口执行已有脚本。恢复会覆盖目标数据库，必须先在隔离环境验证，且需要明确指定目标库；不通过网页执行任意命令。</p><pre>scripts/deploy/backup_database.sh /path/to/production.env
scripts/deploy/restore_database.sh /path/to/production.env /path/to/rays.dump</pre><p>备份包含账号、会话及服务端配置数据，不包含用户设备内的个人图库。恢复脚本会验证备份校验和。</p></el-card>
   </el-tab-pane>
  </el-tabs>
  <el-dialog v-model="dialog" :title="editing ? '编辑账号' : '创建邀请账号'" width="min(520px, 92vw)" :close-on-click-modal="false">
   <el-form label-position="top" @submit.prevent="save"><el-form-item label="昵称"><el-input v-model="form.displayName" maxlength="120"/></el-form-item>
    <template v-if="!editing"><el-form-item label="登录渠道"><el-select v-model="form.provider"><el-option label="邮箱" value="EMAIL"/><el-option label="手机号" value="PHONE"/><el-option label="微信" value="WECHAT"/></el-select></el-form-item><el-form-item label="登录身份"><el-input v-model="form.subject" maxlength="320"/></el-form-item><el-form-item v-if="form.provider==='WECHAT'" label="微信 issuer"><el-input v-model="form.issuer" maxlength="120"/></el-form-item></template>
    <el-form-item label="语言"><el-select v-model="form.locale"><el-option label="简体中文" value="zh-CN"/><el-option label="English" value="en-US"/></el-select></el-form-item>
   </el-form><template #footer><el-button @click="dialog=false">取消</el-button><el-button type="primary" :loading="busy" @click="save">保存</el-button></template>
  </el-dialog>
 </main>
</template>
<script setup lang="ts">
import {ref, reactive, computed, onMounted} from 'vue';
import {ElMessage, ElMessageBox} from 'element-plus';
import request from '@/common/request';
const tab=ref('overview'), busy=ref(false), error=ref(''), dialog=ref(false), query=ref(''),status=ref('');
const accounts=ref<any[]>([]),audit=ref<any[]>([]),overview=ref<any>({}), editing=ref<any>(null);
const form=reactive({displayName:'',provider:'EMAIL',subject:'',issuer:'',locale:'zh-CN'});
const labels:Record<string,string>={INVITED:'待首次登录',ACTIVE:'已激活',DISABLED:'已禁用'};
const stats=computed(()=>[{label:'账号总数',value:overview.value.accounts},{label:'禁用账号',value:overview.value.disabled},{label:'有效会话',value:overview.value.sessions},{label:'审计记录',value:overview.value.audit_events}]);
async function call(method:string,path:string,data?:any){return request.xhrReq(method,'/rays'+path,data);}
async function action(fn:()=>Promise<void>){if(busy.value)return;busy.value=true;error.value='';try{await fn();}catch(e:any){error.value=e?.message || '操作失败，请检查管理连接和权限';}finally{busy.value=false;}}
async function load(){if(tab.value==='accounts')accounts.value=await call('get','/accounts',{q:query.value,status:status.value || undefined});else if(tab.value==='audit')audit.value=await call('get','/audit');else overview.value=await call('get','/overview');}
function refresh(){return action(load);}
function openCreate(){editing.value=null;Object.assign(form,{displayName:'',provider:'EMAIL',subject:'',issuer:'',locale:'zh-CN'});dialog.value=true;}
function openEdit(row:any){editing.value=row;Object.assign(form,row);dialog.value=true;}
async function save(){if(!form.displayName.trim() || (!editing.value && !form.subject.trim())){ElMessage.warning('请填写昵称和登录身份');return;}await action(async()=>{if(editing.value)await call('put','/accounts/'+editing.value.id,{displayName:form.displayName,locale:form.locale,disabled:editing.value.status==='DISABLED'});else await call('post','/accounts',{...form,displayName:form.displayName.trim()});dialog.value=false;await load();ElMessage.success('已保存');});}
async function confirmed(message:string,fn:()=>Promise<void>){try{await ElMessageBox.confirm(message,'确认操作',{type:'warning',confirmButtonText:'确认',cancelButtonText:'取消'});}catch{return;}await action(fn);}
async function toggle(row:any){await confirmed(row.status==='DISABLED'?'解禁此账号？':'禁用此账号并使现有登录失效？',async()=>{await call('put','/accounts/'+row.id,{displayName:row.displayName,locale:row.locale,disabled:row.status!=='DISABLED'});await load();});}
async function remove(row:any){await confirmed('删除此账号的登录身份和服务端凭据并撤销会话；审计保留，个人设备文件不会删除。继续？',async()=>{await call('delete','/accounts/'+row.id);await load();});}
async function revoke(row:any){await confirmed('使此账号的全部现有会话退出？',async()=>{await call('post','/accounts/'+row.id+'/revoke-sessions',{});ElMessage.success('会话已撤销');});}
async function expire(){await confirmed('将已到期会话标记为过期？',async()=>{const result=await call('post','/maintenance/expire-sessions',{});ElMessage.success('已整理 '+result.expired+' 条会话');});}
onMounted(refresh);
</script>
<style scoped>
.rays-console{padding:24px;max-width:1440px;margin:auto}header{display:flex;justify-content:space-between;align-items:center}h1{font-size:26px}header p{color:var(--el-text-color-secondary)}.el-alert,.el-card,.el-descriptions{margin-bottom:18px}.toolbar{display:flex;gap:12px;margin:16px 0;flex-wrap:wrap}.toolbar .el-input{width:260px}.toolbar .el-select{width:150px}strong{font-size:32px}pre{white-space:pre-wrap;overflow-wrap:anywhere}h2{font-size:18px}
</style>
