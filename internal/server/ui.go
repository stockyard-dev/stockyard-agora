package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1.0">
<title>Agora</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--orange:#d4843a;--blue:#5b8dd9;--purple:#9d6bb8;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}
body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5;font-size:13px}
.hdr{padding:.8rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center;gap:1rem;flex-wrap:wrap}
.hdr h1{font-size:.9rem;letter-spacing:2px}
.hdr h1 span{color:var(--rust)}
.main{padding:1.2rem 1.5rem;max-width:980px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(4,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.7rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700;color:var(--gold)}
.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.2rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;flex-wrap:wrap;align-items:center}
.search{flex:1;min-width:180px;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.filter-sel{padding:.4rem .5rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.65rem}
.list{display:flex;flex-direction:column;gap:.7rem}
.poll{background:var(--bg2);border:1px solid var(--bg3);padding:1rem 1.1rem;display:flex;flex-direction:column;gap:.5rem;transition:border-color .15s}
.poll:hover{border-color:var(--leather)}
.poll.closed{opacity:.6}
.poll-hdr{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.poll-title{font-size:.95rem;font-weight:700;color:var(--cream);flex:1}
.poll-actions{display:flex;gap:.3rem;flex-shrink:0}
.poll-desc{font-size:.7rem;color:var(--cd);line-height:1.5}
.poll-meta{display:flex;gap:.5rem;flex-wrap:wrap;align-items:center;font-size:.55rem;color:var(--cm)}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid var(--bg3);color:var(--cm);font-weight:700}
.badge.open{border-color:var(--green);color:var(--green)}
.badge.closed{border-color:var(--cm);color:var(--cm)}
.badge.draft{border-color:var(--orange);color:var(--orange)}
.badge.expired{border-color:var(--red);color:var(--red)}
.badge.type{border-color:var(--blue);color:var(--blue)}
.poll-options{display:flex;flex-direction:column;gap:.35rem;margin-top:.4rem}
.opt{position:relative;padding:.5rem .7rem;background:var(--bg);border:1px solid var(--bg3);cursor:pointer;overflow:hidden;transition:.15s}
.opt:hover{border-color:var(--rust)}
.opt.disabled{cursor:default;opacity:.7}
.opt.disabled:hover{border-color:var(--bg3)}
.opt-bar{position:absolute;top:0;left:0;bottom:0;background:var(--bg3);z-index:0;transition:width .3s}
.opt-content{position:relative;z-index:1;display:flex;justify-content:space-between;align-items:center;gap:.5rem;font-size:.75rem}
.opt-name{color:var(--cream);font-weight:500}
.opt-count{font-family:var(--mono);font-size:.65rem;color:var(--gold);font-weight:700}
.opt-pct{font-family:var(--mono);font-size:.55rem;color:var(--cm);margin-left:.5rem}
.poll-extra{font-size:.55rem;color:var(--cd);margin-top:.4rem;padding-top:.4rem;border-top:1px dashed var(--bg3);display:flex;flex-direction:column;gap:.15rem}
.poll-extra-row{display:flex;gap:.4rem}
.poll-extra-label{color:var(--cm);text-transform:uppercase;letter-spacing:.5px;min-width:90px}
.poll-extra-val{color:var(--cream)}
.btn{font-family:var(--mono);font-size:.6rem;padding:.3rem .55rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:.15s}
.btn:hover{border-color:var(--leather);color:var(--cream)}
.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-p:hover{opacity:.85;color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.btn-del{color:var(--red);border-color:#3a1a1a}
.btn-del:hover{border-color:var(--red);color:var(--red)}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}
.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:520px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}
.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.fr-section{margin-top:1rem;padding-top:.8rem;border-top:1px solid var(--bg3)}
.fr-section-label{font-size:.55rem;color:var(--rust);text-transform:uppercase;letter-spacing:1px;margin-bottom:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.acts .btn-del{margin-right:auto}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.85rem}
@media(max-width:600px){.stats{grid-template-columns:repeat(2,1fr)}}
</style>
</head>
<body>

<div class="hdr">
<h1 id="dash-title"><span>&#9670;</span> AGORA</h1>
<button class="btn btn-p" onclick="openNew()">+ New Poll</button>
</div>

<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar">
<input class="search" id="search" placeholder="Search polls..." oninput="debouncedRender()">
<select class="filter-sel" id="status-filter" onchange="render()">
<option value="">All Statuses</option>
<option value="open">Open</option>
<option value="closed">Closed</option>
<option value="draft">Draft</option>
</select>
<select class="filter-sel" id="type-filter" onchange="render()">
<option value="">All Types</option>
<option value="single_choice">Single Choice</option>
<option value="multi_choice">Multi Choice</option>
</select>
</div>
<div id="list" class="list"></div>
</div>

<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()">
<div class="modal" id="mdl"></div>
</div>

<script>
var A='/api';
var RESOURCE='polls';
var polls=[],pollExtras={},editId=null,searchTimer=null;
var customFields=[];

function fmtDate(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return s;
return d.toLocaleDateString('en-US',{month:'short',day:'numeric',year:'numeric'});
}catch(e){return s}
}

function isExpired(p){
if(!p.expires_at)return false;
return new Date(p.expires_at) < new Date();
}

function debouncedRender(){
clearTimeout(searchTimer);
searchTimer=setTimeout(render,200);
}

function parseOptions(p){
try{return JSON.parse(p.options||'[]')||[]}catch(e){return[]}
}

function parseVotes(p){
try{return JSON.parse(p.votes||'{}')||{}}catch(e){return{}}
}

// ─── Loading ──────────────────────────────────────────────────────

async function load(){
try{
var resps=await Promise.all([
fetch(A+'/polls').then(function(r){return r.json()}),
fetch(A+'/stats').then(function(r){return r.json()})
]);
polls=resps[0].polls||[];
renderStats(resps[1]||{});

try{
var ex=await fetch(A+'/extras/'+RESOURCE).then(function(r){return r.json()});
pollExtras=ex||{};
polls.forEach(function(p){
var x=pollExtras[p.id];
if(!x)return;
Object.keys(x).forEach(function(k){if(p[k]===undefined)p[k]=x[k]});
});
}catch(e){pollExtras={}}
}catch(e){
console.error('load failed',e);
polls=[];
}
render();
}

function renderStats(s){
var total=s.total||0;
var byStatus=s.by_status||{};
var open=byStatus.open||0;
var closed=byStatus.closed||0;
var totalVotes=s.total_votes||0;
document.getElementById('stats').innerHTML=
'<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Polls</div></div>'+
'<div class="st"><div class="st-v">'+open+'</div><div class="st-l">Open</div></div>'+
'<div class="st"><div class="st-v">'+closed+'</div><div class="st-l">Closed</div></div>'+
'<div class="st"><div class="st-v">'+totalVotes+'</div><div class="st-l">Total Votes</div></div>';
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var sf=document.getElementById('status-filter').value;
var tf=document.getElementById('type-filter').value;

var f=polls.slice();
if(q)f=f.filter(function(p){
return(p.title||'').toLowerCase().includes(q)||(p.description||'').toLowerCase().includes(q);
});
if(sf)f=f.filter(function(p){return p.status===sf});
if(tf)f=f.filter(function(p){return p.type===tf});

if(!f.length){
var msg=window._emptyMsg||'No polls yet. Create your first poll.';
document.getElementById('list').innerHTML='<div class="empty">'+esc(msg)+'</div>';
return;
}

var h='';
f.forEach(function(p){h+=pollHTML(p)});
document.getElementById('list').innerHTML=h;
}

function pollHTML(p){
var options=parseOptions(p);
var votes=parseVotes(p);
var totalVotes=0;
options.forEach(function(o){totalVotes+=(votes[o]||0)});
var expired=isExpired(p);
var canVote=p.status==='open'&&!expired;

var cls='poll '+(p.status||'open');

var h='<div class="'+cls+'">';
h+='<div class="poll-hdr">';
h+='<div class="poll-title">'+esc(p.title)+'</div>';
h+='<div class="poll-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(\''+esc(p.id)+'\')">Edit</button>';
h+='</div>';
h+='</div>';

if(p.description)h+='<div class="poll-desc">'+esc(p.description)+'</div>';

h+='<div class="poll-meta">';
if(expired)h+='<span class="badge expired">Expired</span>';
else if(p.status)h+='<span class="badge '+esc(p.status)+'">'+esc(p.status)+'</span>';
if(p.type)h+='<span class="badge type">'+esc(p.type.replace(/_/g,' '))+'</span>';
if(p.expires_at&&!expired)h+='<span>expires '+esc(fmtDate(p.expires_at))+'</span>';
h+='<span>'+totalVotes+' votes</span>';
h+='</div>';

if(options.length){
h+='<div class="poll-options">';
options.forEach(function(opt){
var c=votes[opt]||0;
var pct=totalVotes>0?Math.round((c/totalVotes)*100):0;
var oc=canVote?'opt':'opt disabled';
var click=canVote?(' onclick="vote(\''+esc(p.id)+'\',\''+esc(opt)+'\')"'):'';
h+='<div class="'+oc+'"'+click+'>';
h+='<div class="opt-bar" style="width:'+pct+'%"></div>';
h+='<div class="opt-content">';
h+='<span class="opt-name">'+esc(opt)+'</span>';
h+='<span><span class="opt-count">'+c+'</span><span class="opt-pct">'+pct+'%</span></span>';
h+='</div>';
h+='</div>';
});
h+='</div>';
}else{
h+='<div class="poll-meta" style="color:var(--cm);font-style:italic">No options. Edit poll to add some.</div>';
}

// Custom field display
var customRows='';
customFields.forEach(function(f){
var v=p[f.name];
if(v===undefined||v===null||v==='')return;
customRows+='<div class="poll-extra-row">';
customRows+='<span class="poll-extra-label">'+esc(f.label)+'</span>';
customRows+='<span class="poll-extra-val">'+esc(String(v))+'</span>';
customRows+='</div>';
});
if(customRows)h+='<div class="poll-extra">'+customRows+'</div>';

h+='</div>';
return h;
}

// ─── Voting ──────────────────────────────────────────────────────

async function vote(pollId,option){
try{
var r=await fetch(A+'/polls/'+pollId+'/vote',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify({option:option})});
if(!r.ok){
var e=await r.json().catch(function(){return{}});
alert(e.error||'Vote failed');
return;
}
var d=await r.json();
// Update in-memory poll's votes
for(var i=0;i<polls.length;i++){
if(polls[i].id===pollId){
polls[i].votes=JSON.stringify(d.votes||{});
break;
}
}
load(); // refresh stats
}catch(e){
alert('Network error: '+e.message);
}
}

// ─── Modal ────────────────────────────────────────────────────────

function customFieldHTML(f,value){
var v=value;
if(v===undefined||v===null)v='';
var h='<div class="fr"><label>'+esc(f.label)+'</label>';
if(f.type==='textarea'){
h+='<textarea id="cf-'+f.name+'" rows="2">'+esc(String(v))+'</textarea>';
}else if(f.type==='select'){
h+='<select id="cf-'+f.name+'"><option value="">Select...</option>';
(f.options||[]).forEach(function(o){
var sel=String(v)===String(o)?' selected':'';
h+='<option value="'+esc(String(o))+'"'+sel+'>'+esc(String(o))+'</option>';
});
h+='</select>';
}else if(f.type==='number'){
h+='<input type="number" id="cf-'+f.name+'" value="'+esc(String(v))+'">';
}else{
h+='<input type="text" id="cf-'+f.name+'" value="'+esc(String(v))+'">';
}
h+='</div>';
return h;
}

function formHTML(poll){
var p=poll||{};
var isEdit=!!poll;
var options=isEdit?parseOptions(p):[];

var h='<h2>'+(isEdit?'EDIT POLL':'NEW POLL')+'</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(p.title||'')+'"></div>';
h+='<div class="fr"><label>Description</label><textarea id="f-description" rows="2">'+esc(p.description||'')+'</textarea></div>';
h+='<div class="fr"><label>Options (one per line)</label><textarea id="f-options" rows="4" placeholder="Option A&#10;Option B&#10;Option C">'+esc(options.join('\n'))+'</textarea></div>';
h+='<div class="row2">';
h+='<div class="fr"><label>Type</label><select id="f-type">';
['single_choice','multi_choice'].forEach(function(t){
var sel=(p.type||'single_choice')===t?' selected':'';
h+='<option value="'+t+'"'+sel+'>'+t.replace(/_/g,' ')+'</option>';
});
h+='</select></div>';
h+='<div class="fr"><label>Status</label><select id="f-status">';
['open','closed','draft'].forEach(function(st){
var sel=(p.status||'open')===st?' selected':'';
h+='<option value="'+st+'"'+sel+'>'+st+'</option>';
});
h+='</select></div>';
h+='</div>';
h+='<div class="fr"><label>Expires At (optional)</label><input type="datetime-local" id="f-expires" value="'+esc(toLocalDT(p.expires_at))+'"></div>';

if(customFields.length){
var label=window._customSectionLabel||'Additional Details';
h+='<div class="fr-section"><div class="fr-section-label">'+esc(label)+'</div>';
customFields.forEach(function(f){h+=customFieldHTML(f,p[f.name])});
h+='</div>';
}

h+='<div class="acts">';
if(isEdit){
h+='<button class="btn btn-del" onclick="delPoll()">Delete</button>';
h+='<button class="btn" onclick="resetVotesConfirm()">Reset Votes</button>';
}
h+='<button class="btn" onclick="closeModal()">Cancel</button>';
h+='<button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button>';
h+='</div>';
return h;
}

function toLocalDT(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return'';
var pad=function(n){return n<10?'0'+n:''+n};
return d.getFullYear()+'-'+pad(d.getMonth()+1)+'-'+pad(d.getDate())+'T'+pad(d.getHours())+':'+pad(d.getMinutes());
}catch(e){return''}
}

function fromLocalDT(s){
if(!s)return'';
try{
var d=new Date(s);
if(isNaN(d.getTime()))return'';
return d.toISOString();
}catch(e){return''}
}

function openNew(){
editId=null;
document.getElementById('mdl').innerHTML=formHTML();
document.getElementById('mbg').classList.add('open');
var t=document.getElementById('f-title');
if(t)t.focus();
}

function openEdit(id){
var p=null;
for(var i=0;i<polls.length;i++){if(polls[i].id===id){p=polls[i];break}}
if(!p)return;
editId=id;
document.getElementById('mdl').innerHTML=formHTML(p);
document.getElementById('mbg').classList.add('open');
}

function closeModal(){
document.getElementById('mbg').classList.remove('open');
editId=null;
}

async function submit(){
var titleEl=document.getElementById('f-title');
if(!titleEl||!titleEl.value.trim()){alert('Title is required');return}

var optionsRaw=document.getElementById('f-options').value;
var optionsList=optionsRaw.split('\n').map(function(o){return o.trim()}).filter(function(o){return o});

var body={
title:titleEl.value.trim(),
description:document.getElementById('f-description').value,
type:document.getElementById('f-type').value,
status:document.getElementById('f-status').value,
options:JSON.stringify(optionsList),
expires_at:fromLocalDT(document.getElementById('f-expires').value)
};

var extras={};
customFields.forEach(function(f){
var el=document.getElementById('cf-'+f.name);
if(!el)return;
extras[f.name]=f.type==='number'?(parseFloat(el.value)||0):el.value.trim();
});

var savedId=editId;
try{
if(editId){
var r1=await fetch(A+'/polls/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r1.ok){var e1=await r1.json().catch(function(){return{}});alert(e1.error||'Save failed');return}
}else{
var r2=await fetch(A+'/polls',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});
if(!r2.ok){var e2=await r2.json().catch(function(){return{}});alert(e2.error||'Create failed');return}
var created=await r2.json();
savedId=created.id;
}
if(savedId&&Object.keys(extras).length){
await fetch(A+'/extras/'+RESOURCE+'/'+savedId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(extras)}).catch(function(){});
}
}catch(e){
alert('Network error: '+e.message);
return;
}
closeModal();
load();
}

async function delPoll(){
if(!editId)return;
if(!confirm('Delete this poll?'))return;
await fetch(A+'/polls/'+editId,{method:'DELETE'});
closeModal();
load();
}

async function resetVotesConfirm(){
if(!editId)return;
if(!confirm('Reset all vote counts on this poll?'))return;
await fetch(A+'/polls/'+editId+'/reset',{method:'POST'});
closeModal();
load();
}

function esc(s){
if(s===undefined||s===null)return'';
var d=document.createElement('div');
d.textContent=String(s);
return d.innerHTML;
}

document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal()});

// ─── Personalization ──────────────────────────────────────────────

(function loadPersonalization(){
fetch('/api/config').then(function(r){return r.json()}).then(function(cfg){
if(!cfg||typeof cfg!=='object')return;

if(cfg.dashboard_title){
var h1=document.getElementById('dash-title');
if(h1)h1.innerHTML='<span>&#9670;</span> '+esc(cfg.dashboard_title);
document.title=cfg.dashboard_title;
}

if(cfg.empty_state_message)window._emptyMsg=cfg.empty_state_message;
if(cfg.primary_label)window._customSectionLabel=cfg.primary_label+' Details';

if(Array.isArray(cfg.custom_fields)){
customFields=cfg.custom_fields.filter(function(f){return f&&f.name&&f.label});
}
}).catch(function(){
}).finally(function(){
load();
});
})();
</script>
</body>
</html>`
