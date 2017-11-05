package main

// Wrapper around postgres interactions
type PostgresClient struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

// Add an image to the images table
func (p *PostgresClient) InsertImage() {

}

// Checks the source id to see if we've already stored this image
func (p *PostgresClient) ImageExists(sourceId string) bool {
	return false
}

func NewPostgresClient(pgHost string, pgPort int, pgUser string, pgPassword string, pgDbname string) *PostgresClient {
	p := new(PostgresClient)
	p.Host = pgHost
	p.Port = pgPort
	p.User = pgUser
	p.Password = pgPassword
	p.Dbname = pgDbname
	return p
}
