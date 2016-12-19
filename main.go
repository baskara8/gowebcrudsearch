package main 

import "fmt"
import "sso/config"
import "net"
import "net/http"
import "text/template"
import "flag"
import "strings"
import "crypto/sha1"
import "github.com/kataras/go-sessions"
import "github.com/fatih/color"
import "io"
import "os"
import "time"
import "github.com/tealeg/xlsx"
import "database/sql"
import "math/rand"
import "strconv"
import "encoding/json"

type sysuser struct{
	Userid	 string
	Username string 
	Email 	 string 
	Jekel 	 string 
	Foto 	 string 
	Oldpass  string 
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
		hasspass  := hassdata(password)
		file, header, _ := r.FormFile("file")
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
	sql 			:= "SELECT id,nama,email,jekel,foto FROM "+table+" "+where
	var col,err 	= db.Query(sql)
	config.CheckError(err)
	var each  = sysuser{}
	var res     = []sysuser{}
	for col.Next(){
		var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel,&each.Foto)
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
	sql   := "SELECT id,nama,email,jekel,foto,password FROM "+table +" WHERE id =?"
	var col,err = db.Query(sql,getid)
	config.CheckError(err)
	var each  = sysuser{}
	for col.Next(){
		var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel,&each.Foto,&each.Oldpass)
		config.CheckError(err)
	}
	tmpl.ExecuteTemplate(w,"edit",each)
}
func doedit(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	if r.Method=="POST"{
		var hasspass string
		userid 	 := r.FormValue("userid")
		username := r.FormValue("username")
		email 	 := r.FormValue("email")
		jekel 	 := r.FormValue("jekel")
		password := r.FormValue("password")
		if password !=""{
		hasspass  = hassdata(password)
		}else{
		hasspass  = r.FormValue("oldpassword")
		}
		if userid==""{
			fmt.Fprintln(w,"Isi userid")
		}else if username==""{
			fmt.Fprintln(w,"Isi username")
		}else if email==""{
			fmt.Fprintln(w,"Isi email")
		}else if jekel==""{
			fmt.Fprintln(w,"Pilih jenis kelamin")
		}else{
		files, _, _ := r.FormFile("file")
		if files!=nil{
		var data = sysuser{}
		sql := "SELECT foto FROM siswa WHERE email =?"
		var errs = db.QueryRow(sql,email).Scan(&data.Foto)
		config.CheckError(errs)
		if data.Foto!="" && data.Foto!="default.jpg"{
		_ = os.Remove("./upload/" + data.Foto)
		}
		file, header, _ := r.FormFile("file")
		var statement,err = db.Prepare("UPDATE siswa SET nama =?,email=?,jekel=?,password=?,foto=? WHERE id =?")
		config.CheckError(err)
		statement.Exec(username,email,jekel,hasspass,header.Filename,userid)
		out, _ := os.Create("./upload/" + header.Filename)
		_, _ = io.Copy(out, file)
		file.Close()
		out.Close()
		}else{
		var statement,err = db.Prepare("UPDATE siswa SET nama =?,email=?,jekel=?,password=? WHERE id =?")
		config.CheckError(err)
		statement.Exec(username,email,jekel,hasspass,userid)
		}
		http.Redirect(w,r,"/",301)
		}
	}
}
func delete(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	var getid  = r.URL.Query().Get("id")
	var data = sysuser{}
		sql := "SELECT foto FROM siswa WHERE id =?"
		var errs = db.QueryRow(sql,getid).Scan(&data.Foto)
		config.CheckError(errs)
		if data.Foto!="" && data.Foto!="default.jpg"{
		_ = os.Remove("./upload/" + data.Foto)
	}
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
func downloadexcel(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	var path string = "/upload/siswa.xlsx"
	table := "siswa"
	sql   := "SELECT id AS ID,nama AS NAME,email AS EMAIL,jekel AS GENDER,foto AS AVA FROM "+table
	var col,err = db.Query(sql)
	config.CheckError(err)
	err = generateXLSXFromRows(col, "./"+path,"siswa")
	if err != nil {
		fmt.Print(err)
	}
	url := "http://"+r.Host+path

	timeout := time.Duration(5) * time.Second
	transport := &http.Transport{
		ResponseHeaderTimeout: timeout,
		Dial: func(network, addr string) (net.Conn, error) {
			return net.DialTimeout(network, addr, timeout)
		},
		DisableKeepAlives: true,
	}
	client := &http.Client{
		Transport: transport,
	}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}

	//copy the relevant headers. If you want to preserve the downloaded file name, extract it with go's url parser.
	w.Header().Set("Content-Disposition", "attachment; filename=siswa.xlsx")
	w.Header().Set("Content-Type", r.Header.Get("Content-Type"))
	w.Header().Set("Content-Length", r.Header.Get("Content-Length"))

	//stream the body to the client without fully loading it into memory
	io.Copy(w, resp.Body)
	resp.Body.Close()
	_ = os.Remove("./"+path)
	// http.Redirect(w, r, "/", 301)
}
func uploadexcel(w http.ResponseWriter,r *http.Request){
	tmpl.ExecuteTemplate(w,"upload",nil)
}
func doupload(w http.ResponseWriter,r *http.Request){
	var db = config.Connect()
	defer db.Close()
	rand.Seed(time.Now().Unix())
	file, header, _ := r.FormFile("file")
	var randomdata  = randoms(0,10)
	var files   = randomdata+header.Filename
	out, _ := os.Create("./upload/"+files)
	_, _ = io.Copy(out, file)
	file.Close()
	out.Close()
	excelFileName := "./upload/"+files
    xlFile, err := xlsx.OpenFile(excelFileName)
    if err != nil {
        fmt.Print(err.Error())
    }
    hasspass := hassdata("123456")
    for _, sheet := range xlFile.Sheets {
    	for a, _  := range sheet.Rows {
    	if a>0{
    	sql := "INSERT INTO siswa (id,nama,email,jekel,foto,password) VALUES(?,?,?,?,?,?)"
    	var statement,_ = db.Prepare(sql)
    	statement.Exec(sheet.Cell(a,0).Value,sheet.Cell(a,1).Value,sheet.Cell(a,2).Value,sheet.Cell(a,3).Value,sheet.Cell(a,4).Value,hasspass)
    	}
    
    	}
    }
    _ = os.Remove("./upload/"+files)
    http.Redirect(w, r, "/", 301)
}
func generateXLSXFromRows(rows *sql.Rows, outf string,sheet string) error {
	
	var err error
	// Get column names from query result
	colNames, err := rows.Columns()
	if err != nil {
		return fmt.Errorf("error fetching column names, %s\n", err)
	}
	length := len(colNames)

	// Create a interface slice filled with pointers to interface{}'s
	pointers := make([]interface{}, length)
	container := make([]interface{}, length)
	for i := range pointers {
		pointers[i] = &container[i]
	}

	// Create output xlsx workbook
	xfile := xlsx.NewFile()
	xsheet, err := xfile.AddSheet(sheet)
	if err != nil {
		return fmt.Errorf("error adding sheet to xlsx file, %s\n", err)
	}

	// Write Headers to 1st row
	xrow := xsheet.AddRow()
	xrow.WriteSlice(&colNames, -1)

	// Process sql rows
	for rows.Next() {

		// Scan the sql rows into the interface{} slice
		err = rows.Scan(pointers...)
		if err != nil {
			return fmt.Errorf("error scanning sql row, %s\n", err)
		}

		xrow = xsheet.AddRow()

		// Here we range over our container and look at each column
		// and set some different options depending on the column type.
		for _, v := range container {
			xcell := xrow.AddCell()
			switch v := v.(type) {
			case string:
				xcell.SetString(v)
			case []byte:
				xcell.SetString(string(v))
			case int64:
				xcell.SetInt64(v)
			case float64:
				xcell.SetFloat(v)
			case bool:
				xcell.SetBool(v)
			case time.Time:
				xcell.SetDateTime(v)
			default:
				xcell.SetValue(v)
			}

		}

	}

	// Save the excel file to the provided output file
	err = xfile.Save(outf)
	if err != nil {
		return fmt.Errorf("error writing to output file %s, %s\n", outf, err)
	}

	return nil
}
func downloadfromurl(urllocation,locationslashfilename string){
	img, _ := os.Create(locationslashfilename)
    defer img.Close()

    resp, _ := http.Get(urllocation)
    defer resp.Body.Close()

    b, _ := io.Copy(img, resp.Body)
    fmt.Println("File size: ", b)
}
func randoms(min,max int)string{
	var value = rand.Int() % (max-min+1)+min
	var val   = strconv.Itoa(value)
	return val
}
func users(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	if r.Method == "POST"{
	var db = config.Connect()
	defer db.Close()
	sql 			:= "SELECT id,nama,email,jekel,foto,password FROM siswa"
	var col,err 	= db.Query(sql)
	config.CheckError(err)
	var each  = sysuser{}
	var res     = []sysuser{}
	for col.Next(){
		var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel,&each.Foto,&each.Oldpass)
		config.CheckError(err)
		res 			= append(res,each)
	}
	var result,error = json.Marshal(res)
	if error!=nil{
		http.Error(w,err.Error(),http.StatusInternalServerError)
		return 
	}
	w.Write(result)
	return
	}
	http.Error(w,"",http.StatusBadRequest)
}
func getuser(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application-json")
	if r.Method =="POST"{
		var userid =r.FormValue("userid")
		if userid !=""{
		var db = config.Connect()
		defer db.Close()
		var each  = sysuser{}
		var res     = []sysuser{}
		sql 			:= "SELECT id,nama,email,jekel,foto,password FROM siswa WHERE id = ?"
		var col,err 	= db.Query(sql,userid)
		if err !=nil{
		for col.Next(){
			var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel,&each.Foto,&each.Oldpass)
			config.CheckError(err)
			res 			= append(res,each)
		}
		var result,error = json.Marshal(res)
		if error!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return 
		}
		w.Write(result)
		return
		}else{
		http.Error(w,"Siswa not found",http.StatusBadRequest)
		}
		}else{
		http.Error(w,"Parameter userid kosong",http.StatusBadRequest)
		}
	}
}
func insertuser(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application-json")
	if r.Method =="POST"{
		userid 	  := r.FormValue("userid")
		username  := r.FormValue("username")
		jekel     := r.FormValue("jekel")
		email     := r.FormValue("email")
		password  := r.FormValue("password")
		hasspass  := hassdata(password)
		if userid ==""{
		http.Error(w,"Mohon isi paramater userid",http.StatusBadRequest)
		}else if username ==""{
		http.Error(w,"Mohon isi paramater username",http.StatusBadRequest)
		}else if jekel ==""{
		http.Error(w,"Mohon isi paramater jekel",http.StatusBadRequest)
		}else if userid ==""{
		http.Error(w,"Mohon isi paramater email",http.StatusBadRequest)
		}else if userid ==""{
		http.Error(w,"Mohon isi paramater password",http.StatusBadRequest)
		}else{
		var db = config.Connect()
		defer db.Close()
		var each  = sysuser{}
		var res     = []sysuser{}
		var statement,error1 = db.Prepare("INSERT INTO siswa (id,nama,email,jekel,password) VALUES(?,?,?,?,?)")
		config.CheckError(error1)
		statement.Exec(strings.ToLower(userid),username,email,jekel,hasspass)
		sql 			:= "SELECT id,nama,email,jekel,foto FROM siswa WHERE id = ?"
		var col,err 	= db.Query(sql,userid)
		config.CheckError(err)
		for col.Next(){
			var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel,&each.Foto)
			config.CheckError(err)
			res 			= append(res,each)
		}
		var result,error = json.Marshal(res)
		if error!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return 
		}
		w.Write(result)
		return
		http.Error(w,"Siswa not found",http.StatusBadRequest)
		}
	}
	http.Error(w,"",http.StatusInternalServerError)
}
func deleteuser(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application-json")
	if r.Method =="POST"{
		userid 	  := r.FormValue("userid")
		if userid !=""{
		var db = config.Connect()
		defer db.Close()
		var each  = sysuser{}
		var res     = []sysuser{}
		sql := "SELECT foto FROM siswa WHERE id =?"
		var errs = db.QueryRow(sql,userid).Scan(&each.Foto)
		if errs!=nil{
		fmt.Println("Tidak ada foto")
		}
		if each.Foto!="" && each.Foto!="default.jpg"{
		_ = os.Remove("./upload/" + each.Foto)
		}
		var statement,error1 = db.Prepare("DELETE FROM siswa WHERE id = ?")
		config.CheckError(error1)
		statement.Exec(strings.ToLower(userid))	
		sql1 			:= "SELECT id,nama,email,jekel,password FROM siswa WHERE id = ?"
		var col,err 	= db.Query(sql1,userid)
		if err!=nil{
		for col.Next(){
			var err = col.Scan(&each.Userid,&each.Username,&each.Email,&each.Jekel,&each.Oldpass)
			config.CheckError(err)
			res 			= append(res,each)
		}
		var result,error = json.Marshal(res)
		if error!=nil{
			http.Error(w,err.Error(),http.StatusInternalServerError)
			return 
		}
		w.Write(result)
		return
		}else{
		http.Error(w,"Siswa not found",http.StatusBadRequest)
		}
		}else{
		http.Error(w,"Parameter userid kosong atau salah kunci",http.StatusBadRequest)
		}
	}
		
}
func centang(w http.ResponseWriter, r *http.Request){
	if r.Method=="POST"{
	pilihan := r.FormValue("pilihan")
	if pilihan=="hapus"{
	r.ParseForm()
	data := r.Form["data"]
	for _,d := range data{
		var db = config.Connect()
		defer db.Close()
		var each  = sysuser{}
		sql := "SELECT foto FROM siswa WHERE id =?"
		var errs = db.QueryRow(sql,d).Scan(&each.Foto)
		if errs!=nil{
		fmt.Println("Tidak ada foto")
		}
		if each.Foto!="" && each.Foto!="default.jpg"{
		_ = os.Remove("./upload/" + each.Foto)
		}
		var statement,error1 = db.Prepare("DELETE FROM siswa WHERE id = ?")
		config.CheckError(error1)
		statement.Exec(strings.ToLower(d))	
	}
		http.Redirect(w,r,"home",301)
	}else{
		http.Redirect(w,r,"home?"+pilihan,301)
	}
	}

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
	http.HandleFunc("/download",downloadexcel)
	http.HandleFunc("/upload",uploadexcel)
	http.HandleFunc("/doupload",doupload)
	http.HandleFunc("/users",users)
	http.HandleFunc("/getuser",getuser)
	http.HandleFunc("/insertuser",insertuser)
	http.HandleFunc("/deleteuser",deleteuser)
	http.HandleFunc("/centang",centang)
	//Makes a folder upload to be public
	http.Handle("/upload/", http.StripPrefix("/upload/", http.FileServer(http.Dir("upload"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
}
func main() {
	route()
	var port = flag.String("port","85","isi port")
	flag.Parse()
	var timenow = time.Now()
	c := color.New(color.FgGreen).Add(color.Underline)
	c.Printf("%s %s %s",timenow,"Berjalan di port",*port)
	http.ListenAndServe(":"+*port,nil)
}