package main

type RequestMethod string
type BodyType string

type Env map[string]string

type Response struct {
	StatusCode int      `yaml:"statusCode"`
	BodyType   BodyType `yaml:"bodyType"`
	Body       any      `yaml:"body"`
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

const RMGet RequestMethod = "GET"
const RMPost RequestMethod = "POST"
const RMPatch RequestMethod = "PATCH"
const RMPut RequestMethod = "PUT"
const RMDelete RequestMethod = "DELETE"

const RTPlain BodyType = "plain"
const RMJson BodyType = "map"
