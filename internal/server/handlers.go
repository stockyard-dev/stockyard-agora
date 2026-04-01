package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-agora/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Poll{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var p store.Poll;json.NewDecoder(r.Body).Decode(&p);if p.Question==""{writeError(w,400,"question required");return};s.db.Create(&p);writeJSON(w,201,p)}
func(s *Server)handleVote(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{Option string `json:"option"`;VoterID string `json:"voter_id"`};json.NewDecoder(r.Body).Decode(&req);if req.Option==""{writeError(w,400,"option required");return};if req.VoterID==""{req.VoterID="anonymous"};if err:=s.db.Vote(id,req.Option,req.VoterID);err!=nil{writeError(w,409,"already voted");return};writeJSON(w,200,map[string]string{"status":"voted"})}
func(s *Server)handleResults(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);res,_:=s.db.Results(id);writeJSON(w,200,res)}
func(s *Server)handleClose(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.ClosePoll(id);writeJSON(w,200,map[string]string{"status":"closed"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
