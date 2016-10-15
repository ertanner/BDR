
package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"

	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/securecookie"
)

var db *sql.DB
var cookieHandler = securecookie.New(
	securecookie.GenerateRandomKey(64),
	securecookie.GenerateRandomKey(32))

type store struct {
	StoreID int
	Address string
	City    string
	Zip     string
}
type prod struct {
	ProdID   int
	ProdDesc string
}

var tpl *template.Template
var r = httprouter.New()
var err error

func init(){
	tpl = template.Must(template.ParseGlob("html/*.html"))
}


func main() {

	// Create the database handle, confirm driver is present
	db, err = sql.Open("mysql", "root:Bambie69@/test")
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
	fmt.Println("db opened at dbabis:dbabis11@/test2")
	db.SetMaxIdleConns(100)
	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
	log.Fatalf("Error on opening database connection: %s", err.Error())
	}else {fmt.Println("verified db is open")}

	//r.GET("/", indexPageHandler)
	r.GET("/", HomeHandler)
	r.GET("/index", indexPageHandler)
	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)
	r.GET("/internal", internalPageHandler)
	r.GET("/products", ProductsHandler)
	r.GET("/articles", SitesHandler)

	// Create room for static files serving
	//r.ServeFiles("/node_modules/", http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules"))))
	r.ServeFiles("/html/*filepath", http.Dir("html"))
	r.ServeFiles("/js/*filepath", http.Dir("js"))
	r.ServeFiles("/ts/*filepath", http.Dir("ts"))
	r.ServeFiles("/css/*filepath", http.Dir("css"))
	//http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("."))))
	fmt.Println("router open for business on port 8080")
	http.ListenAndServe(":8080", r)
}

func loginHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	name := r.FormValue("name")
	pass := r.FormValue("password")
	redirectTarget := "/"
	if name != "" && pass != "" {
		// .. check credentials ..
		var msg string
		err := db.QueryRow("SELECT userName FROM user WHERE userName=? and userPwd=?", name, pass).Scan(&msg)
		if err != nil {
			log.Println(err)
			redirectTarget = "/login"
			http.Redirect(w, r, redirectTarget, 302)
			//fmt.Fprintf(w, "Database Error!")
		} //else {
		//	fmt.Fprintf(w, msg)
		//}
		setSession(name, w)
		redirectTarget = "/internal"
		http.Redirect(w, r, redirectTarget, 302)
	}
}
func logoutHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	clearSession(w)
	http.Redirect(w, r, "/", 302)
}

func indexPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
fmt.Fprintf(w, indexPage)
}

func HomeHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
var msg string
err := db.QueryRow("SELECT prodDesc FROM prod WHERE prodID=?", "1").Scan(&msg)
if err != nil {
log.Println(err)
fmt.Fprintf(w, "Database Error!")
} else {
fmt.Fprintf(w, msg)
}
}
func ProductsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
func SitesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
func internalPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
userName := getUserName(r)
if userName != "" {
err = tpl.Execute(w, []string {"a_1", "a&2", "a$3"} )
if err != nil{
log.Fatalln(err)
}
//fmt.Fprintf(w, internalPage, userName)
} else {
http.Redirect(w, r, "/", 302)
}

}

const indexPage = `
<h1>Login</h1>
<form method="post" action="/login">
     <label for="name">User name</label>
     <input type="text" id="name" name="name">
     <label for="password">Password</label>
     <input type="password" id="password" name="password">
     <button type="submit">Login</button>
</form>
`


func setSession(userName string, response http.ResponseWriter) {
value := map[string]string{
"name": userName,
}
if encoded, err := cookieHandler.Encode("session", value); err == nil {
cookie := &http.Cookie{
Name:  "session",
Value: encoded,
Path:  "/",
}
http.SetCookie(response, cookie)
}
}

func getUserName(request *http.Request) (userName string) {
if cookie, err := request.Cookie("session"); err == nil {
cookieValue := make(map[string]string)
if err = cookieHandler.Decode("session", cookie.Value, &cookieValue); err == nil {
userName = cookieValue["name"]
}
}
return userName
}

func clearSession(response http.ResponseWriter) {
cookie := &http.Cookie{
Name:   "session",
Value:  "",
Path:   "/",
MaxAge: -1,
}
http.SetCookie(response, cookie)
}


