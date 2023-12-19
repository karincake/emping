package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

var client = &http.Client{}

func applyEnv(r *Request, e Env) {
	regex, _ := regexp.Compile(`{{([a-zA-Z][\w]*)}}`)

	res := regex.FindAllStringSubmatch(r.Url, 2)
	if len(res) > 0 {
		for idx, _ := range res { // using index is faster
			if _, ok := e[res[idx][1]]; ok {
				r.Url = strings.Replace(r.Url, res[idx][0], e[res[idx][1]], -1)
			}
		}
	}
}

func execute(r Request, e Env) {
	applyEnv(&r, e)

	// create new http req
	req, err := http.NewRequest(string(r.Method), r.Url, bytes.NewBuffer([]byte(r.Body)))
	if r.Header != nil {
		for _, v := range r.Header {
			req.Header.Set(v[0], v[1])
		}
	}

	// send the req
	resp, err := client.Do(req)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}

	_, err = io.ReadAll(resp.Body)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
	}

	fmt.Println("Testing request for:", string(r.Method), r.Url)

	if resp.StatusCode != r.Want.StatusCode {
		os.Stderr.Write([]byte(fmt.Sprintf("Failed: wrong status, want: %v, get: %v\n", r.Want.StatusCode, resp.StatusCode)))
	}

	fmt.Println("Succeed")

	defer resp.Body.Close()
}

func main() {
	args := os.Args[1:]
	if args[0] == "" {
		log.Fatal("please provide a valid yaml file")
	}

	f, err := os.ReadFile(args[0])
	if err != nil {
		log.Fatal(err)
	}

	myJob := Job{
		Env:     Env{},
		ReqList: []Request{},
	}
	if err := yaml.Unmarshal(f, &myJob); err != nil {
		log.Fatal(err)
	}

	for _, req := range myJob.ReqList {
		execute(req, myJob.Env)
	}
}
