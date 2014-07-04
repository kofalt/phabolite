package main

import (
	"database/sql"
	. "fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Database schema details from phabricator
const Database   = "phabricator_user"
const Table_SSH  = Database + "." + "user_sshkey"
const Table_User = Database + "." + "user"

// Represents a single public key
type Key struct {
	header, body, comment string
}

//Golang has a curious connection string syntax
func getConnectionString(server, credentials string) string {
	return credentials + "@" + server + "/" + "phabricator_user"
}

//Connect to the database
func connect(server, credentials string) *sql.DB {
	//Creates a database reference. Does not open a connection
	db, err := sql.Open("mysql", getConnectionString(server, credentials))
	if err != nil { fatal("Error setting up a connection pool.", err) }

	//Check if the server is alive and settings are good (opens connection)
	if err = db.Ping() ; err != nil {
		// This will fail later, but we don't exit immediately in case we're looping.
		// If loop is on, phabolite needs to retry even while the database is down.
		Println("Could not communicate with MySQL:", err)
	}

	return db
}

//Get pubkey information from the database
func queryKeys(db *sql.DB) *sql.Rows {
	//Query the database (see query.go for query documentation)
	rows, err := db.Query(Query)
	if err != nil { fatal("Problem retrieving data from MySQL", err) }

	return rows
}


// Load SQL data into users
func loadUsers(rows *sql.Rows) map[string][]Key {
	users := make(map[string][]Key)

	//Iterate over rows
	for rows.Next() {
		var username, keyBody, keyType, keyComment string
		var isAdmin int

		// Read data
		err := rows.Scan(&username, &isAdmin, &keyBody, &keyType, &keyComment)
		if err != nil { fatal("Could not extract data from query row.", err) }

		// Load into struct
		key := Key{keyType, keyBody, keyComment}

		// Load into map
		if users[username] == nil {
			users[username] = []Key{key}
		} else {
			users[username] = append(users[username], key)
		}
	}

	//Check for any problems while looping
	if err := rows.Err(); err != nil { fatal("Problem looping over key results.", err) }

	return users
}
