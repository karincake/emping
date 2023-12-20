package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
	fmt.Println("Testing request for:", string(r.Method), r.Url)
	resp, err := client.Do(req)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		return
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		os.Stderr.Write([]byte(err.Error()))
		return
	}

	if resp.StatusCode != r.Want.StatusCode {
		os.Stderr.Write([]byte(fmt.Sprintf("Failed: wrong status, want: %v, get: %v\n%v\n", r.Want.StatusCode, resp.StatusCode, string(respBody))))
		return
	}

	errorCount := 0
	if r.Want.Body != nil {
		if r.Want.BodyType == "map" {
			if wantBody, ok := r.Want.Body.(map[string]any); ok {
				keyVal := map[string]any{}
				json.Unmarshal(respBody, &keyVal)
				for kWantBody, vWantBody := range wantBody {
					vRespBody, ok := keyVal[kWantBody]
					if !ok {
						errorCount++
						os.Stderr.Write([]byte(fmt.Sprintf("Failed: misisng key '%v' on response body \n", kWantBody)))
						break
					}
					if vWantBody != vRespBody {
						errorCount++
						os.Stderr.Write([]byte(fmt.Sprintf("Failed: different value for key '%v' \n", kWantBody)))
						break
					}
				}
			} else {
				errorCount++
				os.Stderr.Write([]byte(fmt.Sprintf("Failed: wanted body specs doesn't fit 'bodyType'\n")))
			}
		} else if r.Want.Body != string(respBody) {
			errorCount++
			os.Stderr.Write([]byte(fmt.Sprintf("Failed: wanted body specs doesn't fit 'bodyType'\n")))
		}
	}

	if errorCount > 0 {
		os.Exit(1)
	}
	fmt.Printf("Succeed\n\n")

	defer resp.Body.Close()
}

func main() {
	args := os.Args[1:]
	if len(os.Args) == 1 {
		os.Stderr.Write([]byte("error: please provide a valid yaml file as the first argument\n\n"))
		os.Exit(1)
		return
	}

	f, err := os.ReadFile(args[0])
	if err != nil {
		os.Stderr.Write([]byte("error: " + err.Error() + "\n\n"))
		os.Exit(1)
		return
	}

	myJob := Job{
		Env:     Env{},
		ReqList: []Request{},
	}
	if err := yaml.Unmarshal(f, &myJob); err != nil {
		os.Stderr.Write([]byte("error: " + err.Error() + "\n\n"))
		os.Exit(1)
		return
	}

	for _, req := range myJob.ReqList {
		execute(req, myJob.Env)
	}
}
