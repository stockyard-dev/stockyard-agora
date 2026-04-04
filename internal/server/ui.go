package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Agora</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:800px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.poll{background:var(--bg2);border:1px solid var(--bg3);padding:1rem;margin-bottom:.6rem;transition:border-color .2s}
.poll:hover{border-color:var(--leather)}
.poll-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.poll-title{font-size:.88rem;font-weight:700}
.poll-desc{font-size:.7rem;color:var(--cd);margin-top:.2rem}
.poll-options{margin-top:.5rem}
.opt{display:flex;align-items:center;gap:.5rem;margin-bottom:.3rem;font-size:.7rem}
.opt-bar-bg{flex:1;height:18px;background:var(--bg);border:1px solid var(--bg3);position:relative;cursor:pointer}
.opt-bar{height:100%;background:var(--rust);transition:width .3s}
.opt-label{position:absolute;left:.4rem;top:1px;font-size:.55rem;color:var(--cream);z-index:1}
.opt-count{font-size:.6rem;color:var(--cm);min-width:30px;text-align:right}
.poll-meta{font-size:.55rem;color:var(--cm);margin-top:.4rem;display:flex;gap:.6rem;flex-wrap:wrap;align-items:center}
.poll-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.active{border-color:var(--green);color:var(--green)}.badge.closed{border-color:var(--cm);color:var(--cm)}.badge.draft{border-color:var(--gold);color:var(--gold)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> AGORA</h1><button class="btn btn-p" onclick="openForm()">+ New Poll</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search polls..." oninput="render()"></div>
<div id="polls"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',polls=[],editId=null;

async function load(){var r=await fetch(A+'/polls').then(function(r){return r.json()});polls=r.polls||[];renderStats();render();}

function renderStats(){
var total=polls.length;
var active=polls.filter(function(p){return p.status==='active'}).length;
var totalVotes=0;polls.forEach(function(p){try{var v=JSON.parse(p.votes||'{}');Object.values(v).forEach(function(n){totalVotes+=n})}catch(e){}});
document.getElementById('stats').innerHTML=[
{l:'Polls',v:total},{l:'Active',v:active,c:'var(--green)'},{l:'Total Votes',v:totalVotes}
].map(function(x){return '<div class="st"><div class="st-v" style="'+(x.c?'color:'+x.c:'')+'">'+x.v+'</div><div class="st-l">'+x.l+'</div></div>'}).join('');
}

function render(){
var q=(document.getElementById('search').value||'').toLowerCase();
var f=polls;
if(q)f=f.filter(function(p){return(p.title||'').toLowerCase().includes(q)||(p.description||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('polls').innerHTML='<div class="empty">No polls yet. Create one to start gathering opinions.</div>';return;}
var h='';f.forEach(function(p){
var opts=[];try{opts=JSON.parse(p.options||'[]')}catch(e){}
var votes={};try{votes=JSON.parse(p.votes||'{}')}catch(e){}
var totalV=0;Object.values(votes).forEach(function(n){totalV+=n});
h+='<div class="poll"><div class="poll-top"><div style="flex:1">';
h+='<div class="poll-title">'+esc(p.title)+'</div>';
if(p.description)h+='<div class="poll-desc">'+esc(p.description)+'</div>';
h+='</div><div class="poll-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+p.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+p.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
if(opts.length){h+='<div class="poll-options">';
opts.forEach(function(o){
var count=votes[o]||0;
var pct=totalV>0?Math.round(count/totalV*100):0;
h+='<div class="opt"><div class="opt-bar-bg" onclick="vote(''+p.id+'',''+esc(o)+'')">';
h+='<div class="opt-bar" style="width:'+pct+'%"></div>';
h+='<span class="opt-label">'+esc(o)+'</span>';
h+='</div><span class="opt-count">'+count+'</span></div>';
});
h+='</div>';}
h+='<div class="poll-meta">';
h+='<span class="badge '+(p.status||'active')+'">'+esc(p.status||'active')+'</span>';
if(totalV)h+='<span>'+totalV+' votes</span>';
if(p.expires_at)h+='<span>Expires: '+ft(p.expires_at)+'</span>';
h+='<span>'+ft(p.created_at)+'</span>';
h+='</div></div>';
});
document.getElementById('polls').innerHTML=h;
}

async function vote(id,option){
var poll=null;for(var j=0;j<polls.length;j++){if(polls[j].id===id){poll=polls[j];break;}}
if(!poll||poll.status==='closed')return;
var votes={};try{votes=JSON.parse(poll.votes||'{}')}catch(e){}
votes[option]=(votes[option]||0)+1;
await fetch(A+'/polls/'+id,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify({votes:JSON.stringify(votes)})});
load();
}

async function del(id){if(!confirm('Delete this poll?'))return;await fetch(A+'/polls/'+id,{method:'DELETE'});load();}

function formHTML(poll){
var i=poll||{title:'',description:'',type:'single',options:'[]',status:'active',expires_at:''};
var opts=[];try{opts=JSON.parse(i.options||'[]')}catch(e){}
var isEdit=!!poll;
var h='<h2>'+(isEdit?'EDIT POLL':'NEW POLL')+'</h2>';
h+='<div class="fr"><label>Title *</label><input id="f-title" value="'+esc(i.title)+'" placeholder="What should we decide?"></div>';
h+='<div class="fr"><label>Description</label><input id="f-desc" value="'+esc(i.description)+'"></div>';
h+='<div class="fr"><label>Options (one per line)</label><textarea id="f-opts" rows="4" placeholder="Option A
Option B
Option C">'+opts.join('\n')+'</textarea></div>';
h+='<div class="row2"><div class="fr"><label>Status</label><select id="f-status">';
['active','draft','closed'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s.charAt(0).toUpperCase()+s.slice(1)+'</option>';});
h+='</select></div><div class="fr"><label>Expires</label><input id="f-exp" type="date" value="'+esc(i.expires_at)+'"></div></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;
}

function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');document.getElementById('f-title').focus();}
function openEdit(id){var p=null;for(var j=0;j<polls.length;j++){if(polls[j].id===id){p=polls[j];break;}}if(!p)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(p);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}

async function submit(){
var title=document.getElementById('f-title').value.trim();
if(!title){alert('Title is required');return;}
var optsText=document.getElementById('f-opts').value.trim();
var opts=optsText.split('\n').map(function(o){return o.trim()}).filter(function(o){return o});
var body={title:title,description:document.getElementById('f-desc').value.trim(),options:JSON.stringify(opts),status:document.getElementById('f-status').value,expires_at:document.getElementById('f-exp').value};
if(editId){await fetch(A+'/polls/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{body.votes='{}';await fetch(A+'/polls',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();
}

function ft(t){if(!t)return'';try{return new Date(t).toLocaleDateString('en-US',{month:'short',day:'numeric'})}catch(e){return t;}}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
