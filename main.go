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

var tpl *template.Template
var r = httprouter.New()
var err error


func init(){
	tpl = template.Must(template.ParseGlob("html/*.html"))
}


func main() {
	prodCat1 := make(map[int]string)
	// Create the database handle, confirm driver is present
	db, err = sql.Open("mysql", "root:Bambie69@/Transaction")
	if err != nil {
		log.Fatalf("Error on initializing database connection: %s", err.Error())
	}
	fmt.Println(prodCat1)
	fmt.Println("db opened at root:****@/Transaction")
	db.SetMaxIdleConns(100)
	defer db.Close()

	// make sure connection is available
	err = db.Ping()
	if err != nil {
	log.Fatalf("Error on opening database connection: %s", err.Error())
	}else {fmt.Println("verified db is open")}

	r.GET("/", HomeHandler)
	r.GET("/index", indexPageHandler)
	r.POST("/login", loginHandler)
	r.POST("/logout", logoutHandler)
	r.GET("/internal", internalPageHandler)
	r.GET("/products", ProductsHandler)
	r.GET("/articles", SitesHandler)
	r.GET("/submit" , submit)

	//r.GET("/internal/prodCat1", getProdCat1)
	// Create room for static files serving
	//r.ServeFiles("/node_modules/", http.StripPrefix("/node_modules", http.FileServer(http.Dir("./node_modules"))))
	r.ServeFiles("/html/*filepath", http.Dir("html"))
	r.ServeFiles("/js/*filepath", http.Dir("js"))
	r.ServeFiles("/ts/*filepath", http.Dir("ts"))
	r.ServeFiles("/css/*filepath", http.Dir("css"))

	fmt.Println("router open for business on port 8080")
	http.ListenAndServe(":9080", r)
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
			redirectTarget = "/index"
			http.Redirect(w, r, redirectTarget, 302)

		}else {
			setSession(name, w)
			redirectTarget = "/internal"
			http.Redirect(w, r, redirectTarget, 302)
		}
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
	err := db.QueryRow("SELECT prodDesc FROM products WHERE prodID=?", "1").Scan(&msg)
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

	var cat = [12]map[string]string{}
	//var store = [6]map[string]string{}
	//var dim = make(map[int][]string)

	//var categories = ([]prod, []store,)
	cat[0] = getCat("pCat1")
	cat[1] = getCat("pCat2")
	cat[2] = getCat("pCat3")
	cat[3] = getCat("pCat4")
	cat[4] = getCat("pCat5")
	cat[5] = getCat("pCat6")
	cat[6] = getCat("sCat1")
	cat[7] = getCat("sCat2")
	cat[8] = getCat("sCat3")
	cat[9] = getCat("sCat4")
	cat[10] = getCat("sCat5")
	cat[11] = getCat("sCat6")

	if userName != "" {
		err = tpl.ExecuteTemplate(w, "bdr.html", cat)
		if err != nil{log.Fatalln(err)}
	} else {
		http.Redirect(w, r, "/", 302)
	}
}


const indexPage = `
<h1 align="center">Big Data Rebel</h1><br>

<h1>Login</h1>
<form method="post" action="/login">
     <label for="name">User name</label>
     <input type="text" id="name" name="name">
     <label for="password">Password</label>
     <input type="password" id="password" name="password">
     <button type="submit">Login</button>
</form>
`
const reportPage =`
<h1> You may pick up your report at ??? </h1>

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
func submit(w http.ResponseWriter, r *http.Request, _ httprouter.Params)  {
	fmt.Println(r)

	seasonality()
}

