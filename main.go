package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sync"
)

func main() {
	user := flag.String("u", "", "Github username")
	key := flag.String("key", "", "Github access key")
	dir := flag.String("o", "", "Backup directory")
	forks := flag.Bool("forks", false, "Enable backup of forked repos")
	ssh := flag.Bool("ssh", false, "Using SSH cloning instead of HTTPS")
	flag.Parse()

	if len(*user) == 0 {
		fmt.Println("Github username not specified")
		os.Exit(1)
	}

	if len(*dir) == 0 {
		fmt.Println("Backup directory not specified")
		os.Exit(1)
	}

	git := exec.Command("git")
	err := git.Start()
	if err != nil {
		panic(err)
	}
	git.Wait()

	repos, err := getRepos(*user, *key, *forks)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, repo := range repos {
		wg.Add(1)
		go backupRepo(repo, *dir, *ssh, &wg)
	}

	wg.Wait()
}
