package main

import (
	. "fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

const ConfigFileName = "phabolite.toml"

//A container's settings
type PhaboliteConfig struct {
	//MySQL connection,  example: "unix(/var/run/mysqld/mysqld.sock)"
	Server                 string  `toml:server`

	//MySQL credentials, example: "phabolite:terriblePassword"
	Credentials            string  `toml:credentials`

	//SSH access string, example: "git@127.0.0.1"
	Ssh                    string  `toml:ssh`

	//Run continuously?
	Loop                   bool    `toml:"loop"`

	//How long to wait between loops, if enabled?
	WaitSeconds            time.Duration     `toml:"waitseconds"`
}

// Load and parse a config
func loadConf() *PhaboliteConfig {
	buf, err := ioutil.ReadFile(ConfigFileName)
	if err != nil { fatal("Could not decode conf file:", err) }

	return parseConf(string(buf))
}

// Given a string, parse the config
func parseConf(data string) *PhaboliteConfig {
	//Set defaults
	config := PhaboliteConfig {
		Server:                "unix(/var/run/mysqld/mysqld.sock)",
		Credentials:           "phabolite:terriblePassword",
		Ssh:                   "git@127.0.0.1",
		Loop:                  false,
		WaitSeconds:           30,
	}

	_, err := toml.Decode(data, &config)
	if err != nil { fatal("Could not decode conf file:", err) }

	//Toml don't care, toml does what it wants!
	config.WaitSeconds *= time.Second

	return &config
}

// Sugar function for failures
func fatal(message string, err interface{}) {
	Println(message)
	Println(err)
	os.Exit(1)
}
