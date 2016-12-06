package main 

import "sso/config"
import "strconv"
import "fmt"
import "crypto/sha1"

func hassdata(s string)string{
	var sha = sha1.New()
	sha.Write([]byte(s))
	var encrypted = sha.Sum(nil)
	var encryptedstring = fmt.Sprintf("%x",encrypted)
	return encryptedstring
}
func main(){
var db = config.Connect()
defer db.Close()
var N,M int
fmt.Print("Dari :")
fmt.Scan(&M)
fmt.Print("Sampai :")
fmt.Scan(&N)
for i:=M;i<N;i++{
var statement,error1 = db.Prepare("INSERT INTO siswa (id,nama,email,jekel,password) VALUES(?,?,?,?,?)")
config.CheckError(error1)
var stringangka string = strconv.Itoa(i)
var hasspass = hassdata(stringangka)
statement.Exec("dana"+stringangka,"Dana No "+stringangka,"dana"+stringangka+"@gmail.com","L",hasspass)
}
}