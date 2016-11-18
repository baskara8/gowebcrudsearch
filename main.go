package main 

import "fmt"
import "sso/config"
import "net/http"
import "text/template"
import "flag"
import "strings"
import "crypto/sha1"
import "github.com/kataras/go-sessions"
import "io"
import "os"

type sysuser struct{
	Userid	 string
	Username string 
	Email string 
	Jekel string 
}
var tmpl = template.Must(template.ParseGlob("template/*"))
func hassdata(s string)string{
	var sha = sha1.New()
	sha.Write([]byte(s))
	var encrypted = sha.Sum(nil)
	var encryptedstring = fmt.Sprintf("%x",encrypted)
	return encryptedstring
}
func home(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()

	defer db.Close()
		var bagianWhere,where string
		submit 	  := r.FormValue("submit")
		userid 	  := r.FormValue("userid")
		username  := r.FormValue("username")
		jekel     := r.FormValue("Jekel")
		email     := r.FormValue("email")
		password  := r.FormValue("password")
		file, header, _ := r.FormFile("file")
		hasspass  := hassdata(password)
		if submit =="Tambah Data"{
		if userid==""{
			fmt.Fprintln(w,"Isi userid")
		}else if username==""{
			fmt.Fprintln(w,"Isi username")
		}else if email==""{
			fmt.Fprintln(w,"Isi email")
		}else if jekel==""{
			fmt.Fprintln(w,"Pilih jenis kelamin")
		}else{
		var statement,err = db.Prepare("INSERT INTO siswa (id,nama,email,jekel,password,foto) VALUES(?,?,?,?,?,?)")
		config.CheckError(err)
		statement.Exec(strings.ToLower(userid),username,email,jekel,hasspass,header.Filename)
		out, _ := os.Create("./upload/" + header.Filename)
		_, _ = io.Copy(out, file)
		file.Close()
		out.Close()
		http.Redirect(w,r,"/",301)
		}
		}else{
		bagianWhere = ""
		if userid!=""{
			if bagianWhere ==""{
				bagianWhere = "id ='"+userid+"'"
			}
		}
		if username!=""{
			if bagianWhere ==""{
				bagianWhere += "nama LIKE '%"+username+"%'"
			}
			if bagianWhere !=""{
				bagianWhere += "AND nama LIKE '%"+username+"%'"
			}
		}
		if email!=""{
			if bagianWhere ==""{
				bagianWhere += "email LIKE '%"+email+"%'"
			}
			if bagianWhere !=""{
				bagianWhere += "AND email LIKE '%"+email+"%'"
			}
		}
		if jekel!=""{
			if bagianWhere ==""{
				bagianWhere += "jekel = '"+jekel+"'"
			}
			if bagianWhere !=""{
				bagianWhere += "AND jekel = '"+jekel+"'"
			}
		}
		if bagianWhere==""{
			where = ""
		}else{
			where = "WHERE "+bagianWhere
		}
	table 			:= "siswa"
	sql 			:= "SELECT id,nama,email,jekel FROM "+table+" "+where
	var col,err 	= db.Query(sql)
	config.CheckError(err)
	var each  = sysuser{}
	var res     = []sysuser{}
	for col.Next(){
		var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel)
		config.CheckError(err)
		res 			= append(res,each)
	}
	tmpl.ExecuteTemplate(w, "home", res)
	}
}
func edit(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	var getid  = r.URL.Query().Get("id")
	table := "siswa"
	sql   := "SELECT id,nama,email,jekel FROM "+table +" WHERE id =?"
	var col,err = db.Query(sql,getid)
	config.CheckError(err)
	var each  = sysuser{}
	for col.Next(){
		var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel)
		config.CheckError(err)
	}
	tmpl.ExecuteTemplate(w,"edit",each)
}
func doedit(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	if r.Method=="POST"{
		userid 	 := r.FormValue("userid")
		username := r.FormValue("username")
		email 	 := r.FormValue("email")
		jekel 	 := r.FormValue("jekel")
		password := r.FormValue("password")
		hasspass := hassdata(password)
		if userid==""{
			fmt.Fprintln(w,"Isi userid")
		}else if username==""{
			fmt.Fprintln(w,"Isi username")
		}else if email==""{
			fmt.Fprintln(w,"Isi email")
		}else if jekel==""{
			fmt.Fprintln(w,"Pilih jenis kelamin")
		}else if password==""{
			fmt.Fprintln(w,"Isi Password")
		}else{
		var statement,err = db.Prepare("UPDATE siswa SET nama =?,email=?,jekel=?,password=? WHERE id =?")
		config.CheckError(err)
		statement.Exec(username,email,jekel,hasspass,userid)
		http.Redirect(w,r,"/",301)
		}
	}
}
func delete(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	var getid  = r.URL.Query().Get("id")
	var statement,err = db.Prepare("DELETE FROM siswa WHERE id =?")
	config.CheckError(err)
	statement.Exec(getid)
	http.Redirect(w,r,"/",301)
}
func user(w http.ResponseWriter,r *http.Request){
    session := sessions.Start(w, r)
	var suserid = session.GetString("suserid")
	var data    = make(map[string]string)
	data["suserid"] = suserid
	data["err"] = r.URL.Query().Get("err")
	if suserid!=""{
	tmpl.ExecuteTemplate(w,"user",data)
	}else{
	http.Redirect(w,r,"login?err=Harap login terlebih dahulu",301)
	}
}
func login(w http.ResponseWriter,r *http.Request){
	var data    = make(map[string]string)
	data["err"] = r.URL.Query().Get("err")
	tmpl.ExecuteTemplate(w,"login",data)
}
func dologin(w http.ResponseWriter, r*http.Request){
	var db = config.Connect()
	defer db.Close()
	if r.Method=="POST"{
	email 	 := r.FormValue("email")
	password := r.FormValue("password")
	hasspass := hassdata(password)
	sql := "SELECT id FROM siswa WHERE email =? AND password =?"
	var data = sysuser{}
	var err = db.QueryRow(sql,email,hasspass).Scan(&data.Email)
	if err!=nil{
	http.Redirect(w,r,"user?err=Email dan Password anda Salah",301)
	}else{
	session := sessions.Start(w, r)
	session.Set("suserid", data.Email)
	http.Redirect(w,r,"user",301)
	}
	}
}
func dologout(w http.ResponseWriter, r*http.Request){
	session := sessions.Start(w,r)
    session.Clear()
    sessions.Destroy(w,r)
	http.Redirect(w, r, "/", 302)
}
func about(w http.ResponseWriter,r *http.Request){
	tmpl.ExecuteTemplate(w,"about","")
}
func route(){
	http.HandleFunc("/",home)
	http.HandleFunc("/edit",edit)
	http.HandleFunc("/doedit",doedit)
	http.HandleFunc("/delete",delete)
	http.HandleFunc("/login",login)
	http.HandleFunc("/dologin",dologin)
	http.HandleFunc("/logout",dologout)
	http.HandleFunc("/user",user)
	http.HandleFunc("/about",about)
}
func main() {
	route()
	var port = flag.String("port","85","isi port")
	flag.Parse()
	fmt.Print("Berjalan di port ",*port)
	http.ListenAndServe(":"+*port,nil)
}