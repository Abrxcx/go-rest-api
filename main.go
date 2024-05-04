package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
const maxRetries = 3

func main() {
	var err error

	db, err = sql.Open("mysql", "<username>:<password>@tcp(localhost:3306)/emaildata")
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()
	
	_, err = db.Exec("DROP TABLE IF EXISTS email_data")
	if err != nil {
		log.Fatal(err)
	}
	
	_, err = db.Exec("CREATE TABLE email_data (id INT PRIMARY KEY, name VARCHAR(100), email VARCHAR(100))")
	if err != nil {
		log.Fatal(err)
	}
	
	r := gin.Default()
	r.POST("/api/emaildata", handleRequest)
	r.GET("/api/emaildata/:id/binary", handleBinaryRequest)
	r.GET("/api/moveXML", moveXMLFile)
	r.Run(":8080")
}
