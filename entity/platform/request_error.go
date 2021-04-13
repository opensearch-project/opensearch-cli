package platform

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

//RequestError contains more information that can be used by client to provide
//better error message
type RequestError struct {
	statusCode int
	err        error
	response   []byte
}

//NewRequestError builds RequestError
func NewRequestError(statusCode int, body io.ReadCloser, err error) *RequestError {
	return &RequestError{
		statusCode: statusCode,
		err:        err,
		response:   getResponseBody(body),
	}
}

//Error inherits error interface to pass as error
func (r *RequestError) Error() string {
	return r.err.Error()
}

//StatusCode to get response's status code
func (r *RequestError) StatusCode() int {
	return r.statusCode
}

//GetResponse to get error response from OpenSearch
func (r *RequestError) GetResponse() string {
	var data map[string]interface{}
	if err := json.Unmarshal(r.response, &data); err != nil {
		return string(r.response)
	}
	formattedResponse, _ := json.MarshalIndent(data, "", "  ")
	return string(formattedResponse)
}

//getResponseBody to extract response body from OpenSearch server
func getResponseBody(b io.Reader) []byte {
	resBytes, err := ioutil.ReadAll(b)
	if err != nil {
		fmt.Println("failed to read response")
	}
	return resBytes
}
