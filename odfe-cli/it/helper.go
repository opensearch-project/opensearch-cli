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

package it

import (
	"bytes"
	"context"
	"es-cli/odfe-cli/client"
	"es-cli/odfe-cli/entity"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/suite"
)

type ODFECLISuite struct {
	suite.Suite
	Client  *client.Client
	Profile *entity.Profile
}

//HelperLoadBytes loads file from testdata and stream contents
func HelperLoadBytes(name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return contents
}

// DeleteEcommerceIndex deletes test index created for integration tests
func (a *ODFECLISuite) DeleteIndex(indexName string) {
	_, err := a.callRequest(http.MethodPut, HelperLoadBytes("ecommerce.json"), fmt.Sprintf("%s/%s/_doc", a.Profile.Endpoint, indexName))

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (a *ODFECLISuite) ValidateProfile() error {
	if a.Profile.Endpoint == "" {
		return fmt.Errorf("odfe endpoint cannot be empty. set env ODFE_ENDPOINT")
	}
	if a.Profile.UserName == "" {
		return fmt.Errorf("odfe user name cannot be empty. set env ODFE_USER")
	}
	if a.Profile.Password == "" {
		return fmt.Errorf("odfe endpoint cannot be empty. set env ODFE_PASSWORD")
	}
	return nil
}

//CreateIndex creates test data for plugin processing
func (a *ODFECLISuite) CreateIndex(indexName string) {
	res, err := a.callRequest(http.MethodPut, HelperLoadBytes("ecommerce.json"), fmt.Sprintf("%s/%s/_doc/1?refresh", a.Profile.Endpoint, indexName))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(string(res))
}

func (a *ODFECLISuite) callRequest(method string, reqBytes []byte, url string) ([]byte, error) {
	var reqReader *bytes.Reader
	if reqBytes != nil {
		reqReader = bytes.NewReader(reqBytes)
	}
	r, err := retryablehttp.NewRequest(method, url, reqReader)
	if err != nil {
		return nil, err
	}
	req := r.WithContext(context.Background())
	req.SetBasicAuth(a.Profile.UserName, a.Profile.Password)
	req.Header.Set("Content-Type", "application/json")
	response, err := a.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			return
		}
	}()
	return ioutil.ReadAll(response.Body)
}
