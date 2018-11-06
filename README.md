# Go Article API

This is a restful service, built with golang.

These services include create an article, get article details based on passed key and get articles details based on passed tags and date parameter


## Swagger documentation
- I have created a swagger documentation for articles API

  ### Steps to see documentation
  1. open this link in the browser:  https://editor.swagger.io/  
  2. copy article.yaml file and paste in the open editor

## Dependencies

| Package | Justification|
| ---- | ---- |
| `go-chi/chi` | Router & Middleware |
| `logrus` | Log wrapper|

Dependancies
  - go get gopkg.in/mgo.v2
  - go get github.com/go-chi/ch
  - go github.com/sirupsen/logrus

## setting up Mongodb locally
  - Please install monogodb first, here is docs link :
  https://docs.mongodb.com/manual/tutorial/install-mongodb-enterprise-on-os-x/

## Setting Up Local Environment
* Install Golang. Two options, `brew`, or from their website
  * `brew install golang`
  * [golang.org](https://golang.org/)
  * Go workspace setup - slightly different workspace than most. See the [Golang site](https://golang.org/doc/code.html) for more details.
  1. Create the proper directories if you don't already have them
      ```bash
      mkdir -p $GOPATH/src/github.com/
      mkdir -p $GOPATH/src && mkdir -p $GOPATH/bin && mkdir -p $GOPATH/pkg
      cd $GOPATH/src/github.com/
      ```
  2. Clone repo into the directory you `cd`'d into
  3. Build the binary
    ```bash
    cd article
    go get
    go build
    ```
### Testing Code Locally
1. `go get`
2. `go build`
5. `go run main.go`


### Testing with curl

##### Sample curl request

###### Create Article Request
```
curl -X POST \
  http://localhost:8001/articles \
  -H 'Cache-Control: no-cache' \
  -H 'Content-Type: application/json' \
  -H 'id: 1' \
  -d '{
    "id": "1",
    "title":"latest science shows that potato chips are better for you than sugar",
    "Date": "2016-09-22",
    "Body": "some text, potentially containing simple markup about how potato chips are great",
    "Tags":["health", "fitness", "science"]
}'
```


###### Response:
```
201 Created
```

###### Get Article Request
  ```
  curl -X GET \
  http://localhost:8001/articles/1 \
  -H 'Cache-Control: no-cache'
  ```
###### Get Article Response
```
{
    "id": "1",
    "title": "latest science shows that potato chips are better for you than sugar",
    "date": "2016-09-22",
    "body": "some text, potentially containing simple markup about how potato chips are great",
    "tags": [
        "health",
        "fitness",
        "science"
    ]
}
```

###### Get Article based on tagName and Date Request

```
curl -X GET \
  http://localhost:8001/tags/health/20160922 \
  -H 'Cache-Control: no-cache'
  ```
###### Response:
```
{
    "tag": "health",
    "count": 1,
    "articles": [
        "1"
    ],
    "related_tags": [
        "fitness",
        "science"
    ]
}
```
## TODO List
- Now in current test cases, records getting created in database I need to mock it.
- Add Benchmarking test for all endpoints
