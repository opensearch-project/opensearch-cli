/*
 * Copyright 2021 Amazon.com, Inc. or its affiliates. All Rights Reserved.
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
	"io/ioutil"
	"net/http"
	"opensearch-cli/entity/platform"
	"path/filepath"
	"reflect"
	"testing"
)

func helperLoadBytes(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return contents
}

func TestCommandToCurlRequestParameter(t *testing.T) {
	type args struct {
		request platform.CurlCommandRequest
	}
	tests := []struct {
		name       string
		args       args
		wantResult platform.CurlRequest
		wantErr    bool
	}{
		{
			"success: with data from file",
			args{
				request: platform.CurlCommandRequest{
					Action:      "post",
					Path:        "sample-path/two",
					QueryParams: "a=b&c=d",
					Headers:     "ct:value;h:23",
					Data:        "@testdata/index.json",
					Pretty:      false,
				},
			},
			platform.CurlRequest{
				Action:      http.MethodPost,
				Path:        "sample-path/two",
				QueryParams: "a=b&c=d",
				Headers: map[string]string{
					"ct": "value",
					"h":  "23",
				},
				Data: helperLoadBytes(t, "index.json"),
			},
			false,
		},
		{
			"success: with data from stdin",
			args{
				request: platform.CurlCommandRequest{
					Action:      "post",
					Path:        "sample-path/two",
					QueryParams: "a=b&c=d",
					Headers:     "ct:value;h:23",
					Data:        string(helperLoadBytes(t, "index.json")),
					Pretty:      true,
				},
			},
			platform.CurlRequest{
				Action:      http.MethodPost,
				Path:        "sample-path/two",
				QueryParams: "a=b&c=d&pretty=true",
				Headers: map[string]string{
					"ct": "value",
					"h":  "23",
				},
				Data: helperLoadBytes(t, "index.json"),
			},
			false,
		},
		{
			"success: with basic data",
			args{
				request: platform.CurlCommandRequest{
					Action:       "post",
					Path:         "",
					QueryParams:  "",
					Headers:      "",
					Data:         "",
					Pretty:       true,
					OutputFormat: "yaml",
				},
			},
			platform.CurlRequest{
				Action:      http.MethodPost,
				Path:        "",
				QueryParams: "&pretty=true&format=yaml",
				Headers:     nil,
				Data:        nil,
			},
			false,
		},
		{
			"fail: invalid action",
			args{
				request: platform.CurlCommandRequest{
					Action:      "test",
					Path:        "sample-path/two",
					QueryParams: "a=b&c=d",
					Headers:     "ct:value;h:23",
					Data:        "@testdata/index.json",
					Pretty:      false,
				},
			},
			platform.CurlRequest{},
			true,
		},
		{
			"fail: empty action",
			args{
				request: platform.CurlCommandRequest{
					Action:      "",
					Path:        "sample-path/two",
					QueryParams: "a=b&c=d",
					Headers:     "ct:value;h:23",
					Data:        "@testdata/index.json",
					Pretty:      false,
				},
			},
			platform.CurlRequest{},
			true,
		},
		{
			"fail: invalid header",
			args{
				request: platform.CurlCommandRequest{
					Action:      "post",
					Path:        "sample-path/two",
					QueryParams: "a=b&c=d",
					Headers:     "ct:value:invalid;h:23",
					Data:        "@testdata/index.json",
					Pretty:      false,
				},
			},
			platform.CurlRequest{},
			true,
		},
		{
			"success:  empty header",
			args{
				request: platform.CurlCommandRequest{
					Action:      "Get",
					Path:        "  ",
					QueryParams: "",
					Headers:     "  ;  ",
					Data:        "{}",
					Pretty:      true,
				},
			},
			platform.CurlRequest{
				Action:      http.MethodGet,
				QueryParams: "&pretty=true",
				Headers:     map[string]string{},
				Data:        []byte(`{}`),
			},
			false,
		},
		{
			"fail: invalid data",
			args{
				request: platform.CurlCommandRequest{
					Action:      "post",
					Path:        "",
					QueryParams: "",
					Headers:     "",
					Data:        "this is not a json data",
					Pretty:      false,
				},
			},
			platform.CurlRequest{},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, err := CommandToCurlRequestParameter(tt.args.request)
			if (err != nil) != tt.wantErr {
				t.Errorf("CommandToCurlRequestParameter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotResult, tt.wantResult) {
				t.Errorf("CommandToCurlRequestParameter() gotResult = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}
