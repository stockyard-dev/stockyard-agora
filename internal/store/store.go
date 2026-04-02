package store
import ("database/sql";"fmt";"os";"path/filepath";"time";_ "modernc.org/sqlite")
type DB struct{db *sql.DB}
type Poll struct {
	ID string `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Type string `json:"type"`
	Options string `json:"options"`
	Votes string `json:"votes"`
	Status string `json:"status"`
	ExpiresAt string `json:"expires_at"`
	CreatedAt string `json:"created_at"`
}
func Open(d string)(*DB,error){if err:=os.MkdirAll(d,0755);err!=nil{return nil,err};db,err:=sql.Open("sqlite",filepath.Join(d,"agora.db")+"?_journal_mode=WAL&_busy_timeout=5000");if err!=nil{return nil,err}
db.Exec(`CREATE TABLE IF NOT EXISTS polls(id TEXT PRIMARY KEY,title TEXT NOT NULL,description TEXT DEFAULT '',type TEXT DEFAULT 'single_choice',options TEXT DEFAULT '',votes TEXT DEFAULT '{}',status TEXT DEFAULT 'open',expires_at TEXT DEFAULT '',created_at TEXT DEFAULT(datetime('now')))`)
return &DB{db:db},nil}
func(d *DB)Close()error{return d.db.Close()}
func genID()string{return fmt.Sprintf("%d",time.Now().UnixNano())}
func now()string{return time.Now().UTC().Format(time.RFC3339)}
func(d *DB)Create(e *Poll)error{e.ID=genID();e.CreatedAt=now();_,err:=d.db.Exec(`INSERT INTO polls(id,title,description,type,options,votes,status,expires_at,created_at)VALUES(?,?,?,?,?,?,?,?,?)`,e.ID,e.Title,e.Description,e.Type,e.Options,e.Votes,e.Status,e.ExpiresAt,e.CreatedAt);return err}
func(d *DB)Get(id string)*Poll{var e Poll;if d.db.QueryRow(`SELECT id,title,description,type,options,votes,status,expires_at,created_at FROM polls WHERE id=?`,id).Scan(&e.ID,&e.Title,&e.Description,&e.Type,&e.Options,&e.Votes,&e.Status,&e.ExpiresAt,&e.CreatedAt)!=nil{return nil};return &e}
func(d *DB)List()[]Poll{rows,_:=d.db.Query(`SELECT id,title,description,type,options,votes,status,expires_at,created_at FROM polls ORDER BY created_at DESC`);if rows==nil{return nil};defer rows.Close();var o []Poll;for rows.Next(){var e Poll;rows.Scan(&e.ID,&e.Title,&e.Description,&e.Type,&e.Options,&e.Votes,&e.Status,&e.ExpiresAt,&e.CreatedAt);o=append(o,e)};return o}
func(d *DB)Delete(id string)error{_,err:=d.db.Exec(`DELETE FROM polls WHERE id=?`,id);return err}
func(d *DB)Count()int{var n int;d.db.QueryRow(`SELECT COUNT(*) FROM polls`).Scan(&n);return n}
