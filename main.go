package main

import (
	"bufio"
	"os/exec"
	"sync"
	"time"
)

//var repo = "/home/hrich/go/src/github.com/hashicorp/terraform"

var repo = "/home/hrich/go/src/github.com/duckduckgo/tracker-radar"

//var repo = "/home/hrich/go/src/github.com/trufflesecurity/trufflehog"

type diff struct {
	commit1 string
	commit2 string
}

func main() {
	commits := listCommits()
	commitChan := make(chan diff)
	wg := sync.WaitGroup{}
	go func() {
		for i := 0; i < len(commits)-1; i++ {
			commitChan <- diff{commits[i], commits[i+1]}
		}
		close(commitChan)
	}()
	for i := 0; i < 12; i++ {
		wg.Add(1)
		go func() {
			for _ = range commitChan {
				//runDiff(cdiff.commit1, cdiff.commit2)
			}
			time.Sleep(1 * time.Second)
			wg.Done()
		}()
	}
	wg.Wait()

}

func runDiff(commit1, commit2 string) {
	cmd := exec.Command("git", "-C", repo, "diff", commit1, commit2)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		_ = scanner.Text()
		//fmt.Println(str)
	}

}

func listCommits() []string {
	commits := []string{}
	cmd := exec.Command("git", "-C", repo, "log", "-p", "-U5", "--full-history", "--date=format:%a %b %d %H:%M:%S %Y %z")
	//cmd := exec.Command("git", "-C", repo, "log")
	stdOut, _ := cmd.StdoutPipe()
	err := cmd.Start()
	if err != nil {
		panic(err)
	}
	logReader := bufio.NewReader(stdOut)
	for {
		line, err := logReader.ReadBytes('\n')
		if err != nil {
			break
		}
		if len(line) > 7 && string(line[:6]) == "commit" {
			commit := string(line[7 : len(line)-1])
			commits = append(commits, commit)
		}
	}
	return commits
}
