package sys

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var paths = make(map[string]string)

// GetPath gets the full path of a cli cmd.
func GetPath(cmd string) (folder string, err error) {
	if _, ok := paths[cmd]; ok {
		return paths[cmd], nil
	}

	path, err := exec.LookPath(cmd)
	if err != nil {
		log.Fatalf("please install %s", cmd)
		return "", err
	}
	log.Printf("%s is available at %s\n", cmd, path)
	paths[cmd] = path
	return paths[cmd], nil

}

// RunCmd runs a cmd.
func RunCmd(cmd *exec.Cmd) (result string, err error) {
	log.Printf("running cmd: %s", cmd.Args)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Println(err)
	}

	go func() {
		defer stdin.Close()
		io.WriteString(stdin, "values written to stdin are passed to cmd's standard input")
	}()

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(fmt.Sprint(err) + ": " + fmt.Sprintf("%s\n", out))
		return "", err
	}
	return fmt.Sprintf("%s", out), nil
}

// WriteFile writes the content in a file at filename.
func WriteFile(content string, filename string) (filepath string, err error) {
	filepath = filename
	file, err := os.Create(filepath)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer file.Close()

	fmt.Fprintf(file, content)
	return filepath, nil
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letters = []rune("abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandSeq - produces a random sequence of strings.
func RandSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
