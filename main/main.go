package main
import (
	"fmt"
	"log"
	"github.com/gobuffalo/packr/v2"
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
	log.Print(box.List())
	s, err := box.FindString("changelog/db.changelog-master.yaml")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)
}