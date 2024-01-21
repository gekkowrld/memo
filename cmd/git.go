package cmd

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/gekkowrld/go-gitconfig"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func commit(commitMsg string, filename string) {
	isGitEnabled := getKeyValue("Git").(bool)
	if !isGitEnabled {
		return
	}
	// First move to the memoDirectory, then commit
	memoDirectory := getKeyValue("MemoDir").(string)

	// First check if the repo exists, init it if it doesn't exist
	if !checkIfRepoExists() {
		_, err := git.PlainInit(memoDirectory, false)
		if err != nil {
			log.Fatalf("Can't initialize %s", memoDirectory)
		}
	}

	// Panic if the directory isn't created yet
	if !DirectoryExists(memoDirectory) {
		log.Panicf("Can't find %s, commit failed", memoDirectory)
	}

	// Now stage the file
	read, err := git.PlainOpen(memoDirectory)
  if err != nil {
    log.Panicf("%v open", err)
  }
	work, err := read.Worktree()
  if err != nil {
    log.Fatalf("%v read", err)
  }
  _, err = work.Add(filepath.Join(filename))
  if err != nil {
    log.Fatalf("%v add", err)
  }
  
  stats, err := work.Status()
  if err != nil {
    log.Fatalf("%v stat", err)
  }

  fmt.Println(stats)

	// Should now get the user git credentials
	// Since the program can be run from anywhere, specify the starting location
	username, _ := gogitconfig.GetValue("user.name", memoDirectory)
	email, _ := gogitconfig.GetValue("user.email", memoDirectory)

	// Commit the code
	co, err := work.Commit(commitMsg, &git.CommitOptions{
		Author: &object.Signature{
			Name:  username,
			Email: email,
			When:  time.Now(),
		},
	})

	if err != nil {
    log.Panicf("[%v]: Couldn't Commit", err)
	}

	obj, err := read.CommitObject(co)
	if err != nil {
		log.Panic("Couldn't get ay info")
	}

	fmt.Println(obj)

}

func getGitValues(keyType string) string {
	repoDir := getKeyValue("MemoDir").(string)
	if !checkIfRepoExists() {
		return ""
	}

	read, _ := git.PlainOpen(repoDir)
	cfg, _ := read.Config()

	switch keyType {
	case "username":
		return cfg.User.Name
	case "email":
		return cfg.User.Email
	default:
		return ""
	}
}

func checkIfRepoExists() bool {
	repoDir := getKeyValue("MemoDir").(string)
	_, err := git.PlainOpen(repoDir)
	if err != git.ErrRepositoryNotExists {
		return true
	}

	return false
}
