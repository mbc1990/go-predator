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
	Db       *sql.DB
}

func (p *PostgresClient) GetDB() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		p.Host, p.Port, p.User, p.Password, p.Dbname)
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

// Get a slice of source_ids that already exist in a database.
// This is necessary to avoid doing thousands of simultaneous
// db queries to check if tweets/images have already been processed.
func (p *PostgresClient) GetExistingImages() map[string]bool {
	sqlStatement := `
    SELECT source_id FROM images`
	rows, err := p.Db.Query(sqlStatement)
	defer rows.Close()
	if err != nil {
		panic(err)
	}

	var existing []string
	for rows.Next() {
		var source_id string
		if err := rows.Scan(&source_id); err != nil {
			panic(err)
		}
		existing = append(existing, source_id)
	}
	var res map[string]bool
	res = make(map[string]bool)
	for _, el := range existing {
		res[el] = true
	}
	return res
}

// Add an image to the images table
func (p *PostgresClient) InsertImage(filename string, original_url string,
	source string, source_id string) {
	sqlStatement := `  
  INSERT INTO images (filename, original_url, source, source_id, classified)
  VALUES ($1, $2, $3, $4, $5)`
	_, err := p.Db.Exec(sqlStatement, filename, original_url, source, source_id, false)
	if err != nil {
		panic(err)
	}
}

// Checks the source id to see if we've already stored this image
// N.B. This method is only useful for images that come from
// social media sources that provide an internally unique ID
// connected to the image
func (p *PostgresClient) ImageExists(sourceId string) bool {
	sqlStatement := `
    SELECT COUNT(*) FROM images WHERE source_id IN ($1)`
	rows, err := p.Db.Query(sqlStatement, sourceId)
	defer rows.Close()
	if err != nil {
		panic(err)
	}
	rows.Next()
	var count int
	if err := rows.Scan(&count); err != nil {
		panic(err)
	}
	return count > 0
}

func NewPostgresClient(pgHost string, pgPort int, pgUser string,
	pgPassword string, pgDbname string) *PostgresClient {
	p := new(PostgresClient)
	p.Host = pgHost
	p.Port = pgPort
	p.User = pgUser
	p.Password = pgPassword
	p.Dbname = pgDbname
	p.Db = p.GetDB()
	p.Db.SetMaxOpenConns(50)
	return p
}
