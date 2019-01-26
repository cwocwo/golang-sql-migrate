package migrate

import (
	"fmt"
	"bytes"
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"time"
)

func CloneRepo(repoName string, changelogDir string, gitServerAddr string) {
	gitBaseUrl := gitServerAddr + "git-repos/"

	httpClient := &http.Client{}

	request, err := http.NewRequest("GET", gitBaseUrl + repoName , nil)
	request.Header.Set("Accept", "application/json")
	response, err := httpClient.Do(request)

	if CheckIfError(err, "Get repo info.") {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if CheckIfError(err, "Read response body of repo info.") {
		return
	}
	log.Printf("Get repo %s: %s", repoName, string(body))
	status := response.StatusCode
	fmt.Println(status)

	if len(body) == 0 {
		log.Printf("The repo %s is not exists.", repoName)
		//var clusterinfo = url.Values{}
		////var clusterinfo = map[string]string{}
		//clusterinfo.Add("name", programName)
		//data := clusterinfo.Encode()

		data := bytes.NewBuffer([]byte(`{"name": "` + repoName + `"}`))
		createRepoRequest, err := http.NewRequest("POST", gitBaseUrl, data)
		createRepoRequest.Header.Set("Accept", "application/json")
		createRepoRequest.Header.Set("Content-Type", "application/json")

		if err == nil {
			createRepoResponse, _ := httpClient.Do(createRepoRequest)
			body, _ := ioutil.ReadAll(createRepoResponse.Body)
			fmt.Println(string(body))
		}
	}

	//clone repo
	repoUrl := gitServerAddr +"git/" + repoName + ".git"
	//repoUrl = "http://localhost:8090/git/test"

	log.Printf("git clone %s %s --recursive", repoUrl, changelogDir)

	repo, err := git.PlainClone(changelogDir, false, &git.CloneOptions{
		URL:               repoUrl,
		//RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	//git.PlainClone(changelogDir, false, &git.CloneOptions{
	//	URL: "https://github.com/git-fixtures/basic.git",
	//})

	if ! CheckIfError(err, "cloning repo.") {
		log.Printf("git clone %s successful", repoUrl)
		// ... retrieving the branch being pointed by HEAD
		ref, err := repo.Head()
		if ! CheckIfError(err, "getting head.") {
			log.Print("get git ref")
			// ... retrieving the commit object
			commit, err := repo.CommitObject(ref.Hash())
			if ! CheckIfError(err, "get git commits.") {
				log.Print("get git commit")
				fmt.Println(commit)
			}
		}
	}
	//git.

	log.Printf(changelogDir)
	os.MkdirAll(changelogDir, 0666)
}


func CommitChangeLogs(changelogDir string)  {
	repo, err := git.PlainOpen(changelogDir)
	if CheckIfError(err, "Openning repo.") {
		return
	}
	worktree, err := repo.Worktree()
	if CheckIfError(err, "Getting worktree.") {
		return
	}
	addErr := worktree.AddGlob("./")
	if CheckIfError(addErr, "git add -A.") {
		return
	}
	status, err := worktree.Status()
	log.Print(status)

	commit, err := worktree.Commit("commit sql change logs by go-git", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "go-git",
			Email: "go-git@inspur.com",
			When:  time.Now(),
		},
	})
	if CheckIfError(err, "git commit.") {
		return
	}

	obj, err := repo.CommitObject(commit)
	if CheckIfError(err, "git show -s") {
		return
	}
	log.Print(obj)

	err = repo.Push(&git.PushOptions{})
	if CheckIfError(err, "git push") {
		return
	}
}