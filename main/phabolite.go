package main

import (
	. "fmt"
	"os"
	"time"
)

// Sometimes git is weird about error codes.
// Phabolite wants to continue in the face of panics.
func panicHandler() {
	if err := recover(); err != nil {
		Println(err)
		Println("Phabolite encountered a panic. Probably just git being bad at exit codes.")
	}
}

func main() {
	//Load config, connect to the database and get pubkeys
	config := loadConf()

	// Closure to run update-check, ignoring any panics from failed git commands
	runOnce := func() {
		defer panicHandler()

		// Connect to the database
		db := connect(config.Server, config.Credentials)
		defer db.Close()

		//Load data
		rows := queryKeys(db)
		users := loadUsers(rows)

		// If data has changed, update settings
		updateGitolite(config.Ssh, users)
	}

	// Loop forever
	for {
		runOnce()

		//Exit if user wanted a single run, else sleep
		if config.Loop == false { break }
		time.Sleep(config.WaitSeconds)
	}

	os.Exit(0)
}
