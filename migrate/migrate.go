package migrate

import (
	"log"
	"os"
	//"io"
	"io/ioutil"
	"github.com/gobuffalo/packr/v2"
	"strings"
	"bytes"
	"fmt"
	"encoding/json"
	"net/http"
)

type DbType string

type DataSource struct {
	DbType DbType `json:"dbType"`
	Host string `json:"host"`
	Port int `json:"port"`
	Database string `json:"database"`
	Parameters string `json:"parameters"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Migrate struct {
	Changelog string `json:"changeLog"`
	Contexts string `json:"contexts"`
	DataSource DataSource `json:"dataSource"`
}

const PackrChangLogBox = "changeLogs"

// package the database sql migrate files
// use the first parameter as the changelog path that relative to the main entry
// if ignore parameter the path will be set to ../db as default
func Init(path ... string) *packr.Box {
	changelogPath := "../db"
	if len(path) > 0 {
		changelogPath = path[0]
	}
	return packr.New(PackrChangLogBox, changelogPath)
}

func MigrateChangeLogsWithServer(programName string, serverAddr string, migrate Migrate)  {

	migrateData, err := json.Marshal(migrate)
	if CheckIfError(err, "migrate to json.") {
		return
	}

	log.Printf(string(migrateData))

	httpClient := &http.Client{}
	migrateRequest, err := http.NewRequest("POST", serverAddr + "migrates", bytes.NewBuffer(migrateData))
	if CheckIfError(err, "do migrate.") {
		return
	}
	migrateRequest.Header.Set("Accept", "application/json")
	migrateRequest.Header.Set("Content-Type", "application/json")

	if err == nil {
		createRepoResponse, _ := httpClient.Do(migrateRequest)
		body, _ := ioutil.ReadAll(createRepoResponse.Body)
		fmt.Println(string(body))
	}
}

func ExtractChangeLogs(changelogDir string)  {
	box := packr.New(PackrChangLogBox, "./db")
	files := box.List()
	pathSeparator := string(os.PathSeparator)
	changelogParentDir := changelogDir + "db" + pathSeparator;
	for _, file := range files {
		log.Print(file)
		lastIndexPathSeparator := strings.LastIndex(file, string(os.PathSeparator))
		fileDir := changelogParentDir + file[0:lastIndexPathSeparator]
		os.MkdirAll(fileDir, 0666)
		content, err := box.Find(file)
		if err != nil {
			log.Fatal(err)
		} else {
			err2 := ioutil.WriteFile(changelogParentDir + file, content, 0666)
			if err2 != nil {
				log.Fatal(err2)
			}
		}
	}
}