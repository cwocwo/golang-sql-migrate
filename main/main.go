package main
import (
	"fmt"
	"log"
	"os"
	"strings"
	m "../migrate"
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

	box := m.Init("db/")

	log.Print(box.ResolutionDir)

	// The traditional argv[0] in C is available in os.Args[0] in Go. The flags package simply processes the slice os.Args[1:]
	programName := strings.Replace(os.Args[0], ".", "-", -1)
	//pathSeparator := string(os.PathSeparator)
	//changelogDir := os.TempDir() + pathSeparator + "sql-changelogs" + pathSeparator + programName + pathSeparator
	//
	//gitServerAddr := "http://localhost:8090/"
	//m.CloneRepo(programName, changelogDir, gitServerAddr)
	//m.ExtractChangeLogs(changelogDir)
	//m.CommitChangeLogs(changelogDir)
	//
	dataSource := m.DataSource{"POSTGRESQL", "localhost", 5432, "test", "", "test", "123456a?"}
	migrate := m.Migrate{programName, "", dataSource}
	//m.DoMigrateWithServer(programName, gitServerAddr, migrate)

	m.DoMigrate("http://localhost:8090/", migrate)

}
