package main

type RequestMethod string
type ResponseType string

type Env map[string]string

type Response struct {
	StatusCode int          `yaml:"statusCode"`
	Type       ResponseType `yaml:"type"`
	Value      any          `yaml:"value"`
}

type Request struct {
	Method RequestMethod `yaml:"method"`
	Url    string        `yaml:"url"`
	Header [][2]string   `yaml:"header"`
	Body   string        `yaml:"body"`
	Want   Response      `yaml:"want"`
}

type Job struct {
	Env     Env
	ReqList []Request `yaml:"reqList"`
}

const RMGet = "GET"
const RMPost = "POST"
const RMPatch = "PATCH"
const RMPut = "PUT"
const RMDelete = "DELETE"

const RTPlain = "plain"
const RMJson = "json"
