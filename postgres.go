package main

import "database/sql"
import "fmt"
import _ "github.com/lib/pq"

// Wrapper around postgres interactions
type PostgresClient struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

// Get a DB connection
func (p *PostgresClient) GetDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Dbname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

// Add an image to the images table
func (p *PostgresClient) InsertImage(filename string, original_url string,
	source string, source_id string) {
	db := p.GetDB()
	defer db.Close()
	sqlStatement := `  
  INSERT INTO images (filename, original_url, source, source_id, classified)
  VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(sqlStatement, filename, original_url, source, source_id, false)
	if err != nil {
		panic(err)
	}
}

// Checks the source id to see if we've already stored this image
// N.B. This method is only useful for images that come from
// social media sources that provide an internally unique ID
// connected to the image
func (p *PostgresClient) ImageExists(sourceId string) bool {
	// TODO: Look up by source id
	return false
}

func NewPostgresClient(pgHost string, pgPort int, pgUser string,
	pgPassword string, pgDbname string) *PostgresClient {

	p := new(PostgresClient)
	p.Host = pgHost
	p.Port = pgPort
	p.User = pgUser
	p.Password = pgPassword
	p.Dbname = pgDbname
	return p
}
