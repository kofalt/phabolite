package main

import (
	"bytes"
	"errors"
	. "fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	. "polydawn.net/pogo/gosh"
)

const separator = ":"
const adminRepo = "gitolite-admin.git"

const adminKey = "id_rsa"

type Graph struct {
	// Absolute path to the base/working dir for of the graph git repository.
	dir string

	// Cached command template for exec'ing git
	git Command
}

// clone the gitolite graph to a temp folder
func cloneGraph(sshURI string) *Graph {
	// Aquire a unique temporary folder
	dir, err := ioutil.TempDir("", "phabolite-")
	if err != nil { fatal("Could not generate temporary folder", err) }

	g := &Graph{
		dir: dir,
		git: Sh("git")(Opts{Cwd: dir}),
	}

	// Clone to temp dir and return graph
	g.git("clone", sshURI + separator + adminRepo, dir)()
	return g
}

// Remove all keys that are no longer in use.
func (g *Graph) clean(users map[string][]Key) {
	// Get list of current keys
	keys, err := ioutil.ReadDir(filepath.Join(g.dir, "keydir"))
	if err != nil { fatal("Could not list keys:", err)}

	// For each key, check there's a corresponding user
	for _, f := range keys {
		key := f.Name()
		username := strings.TrimSuffix(key, ".pub")

		// Do not remove admin user or users who still have keys
		if username == adminKey || username == ".gitkeep" { continue }

		// Remove any users who do not have keys anymore
		if users[username] == nil {
			Println("Removing", username)
			g.git("rm", filepath.Join(g.dir, "keydir", key))()
		}
	}

	// Make sure git keeps the folder
	Sh("touch")(filepath.Join(g.dir, "keydir", ".gitkeep"))()
}

// Given a set of users, write their pubkeys
func (g *Graph) writePubkeys(users map[string][]Key) {
	for username, keys := range users {
		// Generate key path
		path := filepath.Join(g.dir, "keydir", username + ".pub")

		// Concat keys
		var keyStr = ""
		for _, key := range keys {
			keyStr += key.header + " " + key.body + " " + key.comment + "\n"
		}

		// Write keys
		err := ioutil.WriteFile(path, []byte(keyStr), 0644)
		if err != nil { fatal("Could not write pubkey file:", err) }
	}
}

// Saves changes
func (g *Graph) finish() {
	g.git("add", "--all")()

	// Run commit, accepting non-zero returns & saving all output
	var buf bytes.Buffer
	commit := g.git("commit", "-m", "Phabolite generated update").BakeOpts(Opts{Out: &buf, Err: &buf}).Start()
	code := commit.GetExitCode()
	output := buf.String()

	// Abort if commit for any reason other than empty commit
	if code != 0 {
		// Git is horrible about exit codes. Manually see what the error was.
		if strings.Contains(output, "nothing to commit, working directory clean") {
			Println("Nothing to commit, either change detection is bogus or phabolite has a bug.")
		} else {
			fatal("Git commit failed:", errors.New(output))
		}
	} else {
		g.git("push")(Opts{Out: nil})()
		Println("Config update; commit pushed")
	}
}

// Remove temporary folder
func (g *Graph) cleanUp() {
	Sh("rm")("-rf", g.dir)()
}

// Full workflow
func updateGitolite(sshURI string, users map[string][]Key) {
	graph := cloneGraph(sshURI)
	graph.clean(users)
	graph.writePubkeys(users)
	graph.finish()
	graph.cleanUp()
}
