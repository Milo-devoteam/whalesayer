package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var logger *log.Logger = log.New(os.Stdout, "[whalesayer]", 1)
var COW_PATH string = os.Getenv("COWPATH")
var PORT string = ":8080"

func main() {
	rand.Seed(time.Now().Unix())
	http.HandleFunc("/", cowsay)

	fmt.Println("Whale listening on", PORT)
	http.ListenAndServe(PORT, nil)
}

func cowsay(w http.ResponseWriter, r *http.Request) {
	buf := new(bytes.Buffer)
	n, _ := buf.ReadFrom(r.Body)
	msg := buf.String()
	if n == 0 {
		msg = "say something!"
	}
	logger.Println(r.Method, "request to:", r.URL)

	animal := select_animal()
	response := retry_say(animal, msg)
	fmt.Fprintln(w, response)

	// Replicate cowsay message to stdout
	logger.Println("responded with:\n", response)
}

func retry_say(animal, msg string) string {
	cmd := exec.Command("cowsay", "-f", animal, msg)
	stdout, err := cmd.Output()
	if err != nil {
		// Retry with regular cow
		cmd = exec.Command("cowsay", msg)
		stdout, _ = cmd.Output()
	}
	return string(stdout)
}

func select_animal() string {
	files, err := ioutil.ReadDir(COW_PATH)
	arr := make([]string, len(files))
	if err != nil {
		panic(err.Error())
	}

	for i, file := range files {
		arr[i] = strings.Replace(file.Name(), ".cow", "", 1)
	}

	return arr[rand.Intn(len(arr)-1)]
}
