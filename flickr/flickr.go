package flickr

import (
    "encoding/json"
    "io/ioutil"
    "net/http"
    "strings"
)

const (
    endpoint = "https://api.flickr.com/services/rest/"
)

type Flickr struct {
    ApiKey string
}

type Request struct {
    ApiKey string
    Method string
    Arguments map[string]string
}

type Response struct {
    Code int
    Body string
    Message Message
}

type Message struct {
    Photosets Photosets
    Photoset Photoset
    Sizes Sizes
    Stat string
}

type Photo struct {
    Id string
    Title string
}

type Photoset struct {
    Id string
    Photo []Photo
    Title string
}

type PhotosetLite struct {
    Id string
    Title Title
}

type Photosets struct {
    Page uint
    Pages uint
    PerPage uint
    Total uint
    Photoset []PhotosetLite
}

type Title struct {
    Content string `json:"_content"`
}

type Size struct {
    Label string
    Source string
}

type Sizes struct {
    Size []Size
}

func (flickr *Flickr) Request(method string, arguments map[string]string) (*Response, error)  {
    r := Request{flickr.ApiKey, method, arguments}
    return r.Exectue()
}

func (request *Request) URL() (string) {
    args := request.Arguments
    args["api_key"] = request.ApiKey
    args["method"] = request.Method
    args["format"] = "json"

    return endpoint + "?" + request.queryString()
}

func (request *Request) Exectue() (response *Response, ret error) {
    var body []byte

    s := request.URL()
    response = new(Response)

    res, err := http.Get(s)
    if err != nil {
        return response, err
    }

    body, _ = ioutil.ReadAll(res.Body)
    body = stripJsonP(body)

    response.Code = res.StatusCode
    response.Body = string(body)
    err = json.Unmarshal(body, &response.Message)
    if err != nil {
        return nil, err
    }

    return
}

func (request *Request) queryString() (string) {
    nArguments := len(request.Arguments)
    arguments := make([]string, nArguments)
    i := 0

    for argument, value := range request.Arguments {
        arguments[i] = argument + "=" + value
        i++
    }

    return strings.Join(arguments, "&")
}

func stripJsonP(body []byte) ([]byte) {
    return body[14:len(body) - 1]
}
