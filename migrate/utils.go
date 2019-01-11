package migrate

import "fmt"

func CheckIfError(err error, message string) bool {
	if err == nil {
		return false
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m \n", fmt.Sprintf("error: %s, %s", err, message))
	return  true
}
