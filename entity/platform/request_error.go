/*
 * SPDX-License-Identifier: Apache-2.0
 *
 * The OpenSearch Contributors require contributions made to
 * this file be licensed under the Apache-2.0 license or a
 * compatible open source license.
 *
 * Modifications Copyright OpenSearch Contributors. See
 * GitHub history for details.
 */
/*
 * Copyright 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License").
 * You may not use this file except in compliance with the License.
 * A copy of the License is located at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 * or in the "license" file accompanying this file. This file is distributed
 * on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
 * express or implied. See the License for the specific language governing
 * permissions and limitations under the License.
 */
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
