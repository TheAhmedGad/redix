package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"redix/pkg/server"

	_ "github.com/go-sql-driver/mysql" // Register MySQL driver

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	// Define command line flags
	mysqlHost := flag.String("mysql-host", "localhost", "MySQL host address")
	mysqlPort := flag.String("mysql-port", "3306", "MySQL port")
	mysqlUser := flag.String("mysql-user", "root", "MySQL username")
	mysqlPass := flag.String("mysql-pass", "root", "MySQL password")
	mysqlDB := flag.String("mysql-db", "redix", "MySQL database name")
	redixPort := flag.String("port", ":6379", "Redix server port")

	flag.Parse()

	// Build DSN from flags or fallback to environment variable
	var dsn string
	if *mysqlHost != "" {
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			*mysqlUser, *mysqlPass, *mysqlHost, *mysqlPort, *mysqlDB)
	} else {
		dsn = os.Getenv("MYSQL_DSN")
	}

	if dsn == "" {
		log.Fatal("MySQL connection information not provided. Use command line flags or MYSQL_DSN environment variable")
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("MySQL connect error: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("MySQL ping failed: %v", err)
	}
	defer db.Close()

	srv := server.New(db)
	log.Printf("🚀 Redix server running on %s", *redixPort)
	if err := srv.Listen(*redixPort); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
