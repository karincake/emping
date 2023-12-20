# Emping
Simple E2E tester. It uses yaml file to configure environment and target api.

## Installation and Usage
Install the app
`go install github.com/karincake/emping`

Running the app uses the first parameter as yaml file name for the source.
`emping [yaml_file]`

## Source Specification
Conventional for the yaml formated source. Key or val of `any` follows yaml standards.
- [some_key] - optional key of `any`
- [*:if-someKey=x] - optional key, if key `someKey` value is value `x`
- [some_val] - optional val
- <some_val> - mandatory value, still depends on the key existence
- <some_*_val> - mandatory value, with external specifaction of *

The yaml formatted source
```
env:
  [some_key]: <some_val>

reqList:
  - url: <some_val>
    method: GET|POST|PUT|PATCH|DELETE
    header:
      - [<val>, <val>]
    body: <some_json_formated_val>
    want:
      statusCode: <some_http_status_code_val>
      [bodyType]: <plain|map>
      [body:if-bodyType=plain]: <some_val>
      [body:if-bodyType=map]:
        [some_key]: <some_val>
```
