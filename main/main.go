package main
import (
	"fmt"
	"log"
	"os"
	//"io"
	"io/ioutil"
	"github.com/gobuffalo/packr/v2"
	"strings"
	"gopkg.in/src-d/go-git.v4"
	"net/http"
	//"net/url"
	"bytes"
)
var(
	version string
	gitcommit string
	buildstamp string
)

func main() {
	fmt.Printf("version: %s\n", version)
	fmt.Printf("gitcommit: %s\n", gitcommit)
	fmt.Printf("buildstamp: %s\n", buildstamp)

	box := packr.New("myBox", "./db")

	log.Print(box.ResolutionDir)


	//s, err := box.FindString("changelog/db.changelog-master.yaml")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(s)

	//The traditional argv[0] in C is available in os.Args[0] in Go. The flags package simply processes the slice os.Args[1:]
	programName := strings.Replace(os.Args[0], ".", "-", -1)

	log.Print(git.AllTags)

	gitBaseUrl := "http://localhost:8090/git-repos/"

	httpClient := &http.Client{}

	request, err := http.NewRequest("GET", gitBaseUrl + programName , nil)
	request.Header.Set("Accept", "application/json")
	response, _ := httpClient.Do(request)
	body, _ := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	status := response.StatusCode
	fmt.Println(status)

	if len(body) == 0 {
		log.Printf("The repo %s is not exists.", programName)
		//var clusterinfo = url.Values{}
		////var clusterinfo = map[string]string{}
		//clusterinfo.Add("name", programName)
		//data := clusterinfo.Encode()

		data := bytes.NewBuffer([]byte(`{"name": "` + programName + `"}`))
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
	repoUrl := "http://localhost:8090/git/" + programName + ".git"
	//repoUrl = "http://localhost:8090/git/test"
	pathSeparator := string(os.PathSeparator)
	changelogDir := os.TempDir() + pathSeparator + "sql-changelogs" + pathSeparator + programName + pathSeparator
	log.Printf("git clone %s %s --recursive", repoUrl, changelogDir)

	repo, err := git.PlainClone(changelogDir, false, &git.CloneOptions{
		URL:               repoUrl,
		//RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	//git.PlainClone(changelogDir, false, &git.CloneOptions{
	//	URL: "https://github.com/git-fixtures/basic.git",
	//})

	if ! CheckIfError(err) {
		log.Printf("git clone %s successful", repoUrl)
		// ... retrieving the branch being pointed by HEAD
		ref, err := repo.Head()
		if ! CheckIfError(err) {
			log.Print("get git ref")
			// ... retrieving the commit object
			commit, err := repo.CommitObject(ref.Hash())
			if ! CheckIfError(err) {
				log.Print("get git commit")
				fmt.Println(commit)
			}
		}
	}
	//git.

	log.Printf(changelogDir)
	os.MkdirAll(changelogDir, 0666)
	files := box.List()
	for _, file := range files {
		log.Print(file)
		//lastIndexPathSeparator := strings.LastIndex(file, string(os.PathSeparator))
		//fileDir := changelogDir + file[0:lastIndexPathSeparator]
		//os.MkdirAll(fileDir, 0666)
		//content, err := box.Find(file)
		//if(err != nil) {
		//	log.Fatal(err)
		//} else {
		//	err2 := ioutil.WriteFile(changelogDir + file, content, 0666)
		//	if(err2 != nil) {
		//		log.Fatal(err2)
		//	}
		//}
	}
}

func CheckIfError(err error) bool {
	if err == nil {
		return false
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	return  true
}
