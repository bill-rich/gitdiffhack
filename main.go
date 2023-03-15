package main

import (
	"bufio"
	"fmt"
	"os/exec"
)

var repo = "/home/hrich/go/src/github.com/hashicorp/terraform"

// var repo = "/home/hrich/go/src/github.com/duckduckgo/tracker-radar"

//var repo = "/home/hrich/go/src/github.com/trufflesecurity/trufflehog"

type diff struct {
	commit1 string
	commit2 string
}

func main() {
	commits := listCommits()
	commitChan := make(chan diff)
	for i := 0; i < len(commits)-1; i++ {
		runDiff(commits[i], commits[i+1])
	}
	close(commitChan)

}

func runDiff(commit1, commit2 string) {
	fmt.Println("Running diff for commits", commit1, commit2)

	cmd := exec.Command("git", "-C", repo, "diff", "--name-only", commit1, commit2)

	//cmd := exec.Command("git", "-C", repo, "log", "-p", "-U5", "--full-history", "--date=format:%a %b %d %H:%M:%S %Y %z")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(stdout)
	// f, err := os.Create(fmt.Sprintf("/home/hrich/tmp/%d.txt", time.Now().UnixNano()))
	// if err != nil {
	// 	panic(err)
	// }
	// defer f.Close()
	for scanner.Scan() {
		filename := scanner.Text()
		cmd := exec.Command("git", "-C", repo, "diff", commit1, commit2, "--", filename)

		_, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err := cmd.Start(); err != nil {
			panic(err)
		}

		// f, err := os.Create(fmt.Sprintf("/home/hrich/tmp/%s-%s-%d.txt", commit1, commit2, time.Now().UnixNano()))
		// if err != nil {
		// 	fmt.Println(err)
		// 	continue
		// }
		// defer f.Close()
		//somebytes := make([]byte, 1000000)
		//for {
		// _, err := stdout.Read(somebytes)
		//if errors.Is(err, io.EOF) {
		//	break
		//}
		//if bytes.Contains(somebytes, []byte("password")) {
		//	fmt.Println("Found password in", filename)
		//}
		//}

	}

}

func listCommits() []string {
	commits := []string{}
	cmd := exec.Command("git", "-C", repo, "log")
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
			fmt.Println(commit)
		}
	}
	cmd.Wait()
	return commits
}
