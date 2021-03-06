To deal with Relational DB's, Go has to install the ORM Library package so we can do the operations in the databaseand also connect to the DB

So to use the GORM library,

main.go
 import (
 “log”
 "github.com/jinzhu/gorm"
 _ "github.com/jinzhu/gorm/dialects/mysql"
)
 
func main() {
 db, err := gorm.Open("mysql", "root:root@tcp(127.0.0.1:3306)/ormdemo?charset=utf8&parseTime=True")
 // to connect to the DB
defer db.Close()
 //Always close the database when it is not in use using defer db.Close()
 if err!=nil{
 log.Println(“Connection Failed to Open”)
 } 
 log.Println(“Connection Established”)
}
