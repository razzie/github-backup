package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"sync"
)

// exists returns whether the given file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func pullRepo(repo, dir string, wg *sync.WaitGroup) {
	defer wg.Done()

	cmd := exec.Command("git", "-C", dir, "pull")
	err := cmd.Start()
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Print(err.Error())
		return
	}

	fmt.Println("pulling", repo, "done")
}

func cloneRepo(repo, dir string, ssh bool, wg *sync.WaitGroup) {
	defer wg.Done()

	var url string
	if ssh {
		url = fmt.Sprintf("git@github.com:%s.git", repo)
	} else {
		url = fmt.Sprintf("https://github.com/%s.git", repo)
	}

	cmd := exec.Command("git", "clone", url, dir)
	err := cmd.Start()
	if err != nil {
		log.Print(err.Error())
		return
	}
	err = cmd.Wait()
	if err != nil {
		log.Print(err.Error())
		return
	}

	fmt.Println("cloning", repo, "->", dir, "done")
}

func backupRepo(repo, workdir string, ssh bool, wg *sync.WaitGroup) {
	fullpath := path.Join(workdir, repo)
	if ok, _ := exists(fullpath); ok {
		pullRepo(repo, fullpath, wg)
	} else {
		cloneRepo(repo, fullpath, ssh, wg)
	}
}
