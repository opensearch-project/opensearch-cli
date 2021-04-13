// +build integration

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
	"encoding/json"
	"fmt"
	"net/http"
	"opensearch-cli/client"
	adctrl "opensearch-cli/controller/ad"
	"opensearch-cli/controller/platform"
	"opensearch-cli/entity"
	adentity "opensearch-cli/entity/ad"
	"opensearch-cli/environment"
	adgateway "opensearch-cli/gateway/ad"
	esg "opensearch-cli/gateway/platform"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const (
	EcommerceIndexName     = "ecommerce"
	EcommerceIndexFileName = "ecommerce"
)

//ADTestSuite suite specific to AD plugin
type ADTestSuite struct {
	CLISuite
	DetectorRequest adentity.CreateDetectorRequest
	Detector        adentity.CreateDetector
	DetectorId      string
	ADGateway       adgateway.Gateway
	ESController    platform.Controller
}

func getRawFeatureAggregation() []byte {
	return []byte(`
	{
		"sum_value": {
			"sum": {
				"field": "total_quantity"
			}
		}
	}`)
}

//SetupSuite runs once for every test suite
func (a *ADTestSuite) SetupSuite() {
	var err error
	a.Client, err = client.New(nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.Profile = &entity.Profile{
		Name:     "test",
		Endpoint: os.Getenv(environment.OPENSEARCH_ENDPOINT),
		UserName: os.Getenv(environment.OPENSEARCH_USER),
		Password: os.Getenv(environment.OPENSEARCH_PASSWORD),
	}
	if err = a.ValidateProfile(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	a.CreateIndex(EcommerceIndexFileName, "")
	g := esg.New(a.Client, a.Profile)
	a.ESController = platform.New(g)
	a.ADGateway = adgateway.New(a.Client, a.Profile)
	a.DetectorRequest = getCreateDetectorRequest()
	a.Detector = adentity.CreateDetector{
		Name:        "setup-detector-it1",
		Description: a.DetectorRequest.Description,
		TimeField:   a.DetectorRequest.TimeField,
		Index:       a.DetectorRequest.Index,
		Features: []adentity.Feature{
			{
				Name:             "sum_value",
				Enabled:          true,
				AggregationQuery: getRawFeatureAggregation(),
			},
		},
		Interval: adentity.Interval{
			Period: adentity.Period{
				Duration: 1,
				Unit:     "Minutes",
			},
		},
		Delay: adentity.Interval{
			Period: adentity.Period{
				Duration: 1,
				Unit:     "Minutes",
			},
		},
	}
}

func (a *ADTestSuite) TearDownSuite() {
	a.DeleteIndex(EcommerceIndexName)
}

// This will run right before the test starts
// and receives the suite and test names as input
func (a *ADTestSuite) BeforeTest(suiteName, testName string) {
	// We don't need to create detector for create use case
	if testName != "TestCreateDetectors" {
		a.CreateDetectorUsingRESTAPI(a.T())
	}
}

// This will run after test finishes
// and receives the suite and test names as input
func (a *ADTestSuite) AfterTest(suiteName, testName string) {
	if testName != "TestCreateDetectors" || a.DetectorId != "" {
		a.StopDetectorUsingRESTAPI(a.T(), a.DetectorId)
		a.DeleteDetectorUsingRESTAPI(a.T(), a.DetectorId)
	}
}

//DeleteDetectorUsingRESTAPI helper to delete detector using rest api
func (a *ADTestSuite) DeleteDetectorUsingRESTAPI(t *testing.T, ID string) {
	indexURL := fmt.Sprintf("%s/_opendistro/_anomaly_detection/detectors/%s", a.Profile.Endpoint, ID)
	_, err := a.callRequest(http.MethodDelete, []byte(""), indexURL)
	if err != nil {
		t.Fatal(err)
	}
}

//StartDetectorUsingRESTAPI helper to start detector using rest api
func (a *ADTestSuite) StartDetectorUsingRESTAPI(t *testing.T, ID string) {
	if ID == "" {
		t.Fatal("Detector ID cannot be empty")
	}
	indexURL := fmt.Sprintf("%s/_opendistro/_anomaly_detection/detectors/%s/_start", a.Profile.Endpoint, ID)
	_, err := a.callRequest(http.MethodPost, []byte(""), indexURL)
	if err != nil {
		t.Fatal(err)
	}
}

//StopDetectorUsingRESTAPI helper to stop detector using rest api
func (a *ADTestSuite) StopDetectorUsingRESTAPI(t *testing.T, ID string) {
	if ID == "" {
		t.Fatal("Detector ID cannot be empty")
	}
	indexURL := fmt.Sprintf("%s/_opendistro/_anomaly_detection/detectors/%s/_stop", a.Profile.Endpoint, ID)
	_, err := a.callRequest(http.MethodPost, []byte(""), indexURL)
	if err != nil {
		t.Fatal(err)
	}
}

//CreateDetectorUsingRESTAPI helper to create detector using rest api
func (a *ADTestSuite) CreateDetectorUsingRESTAPI(t *testing.T) {
	indexURL := fmt.Sprintf("%s/_opendistro/_anomaly_detection/detectors", a.Profile.Endpoint)
	reqBytes, err := json.Marshal(a.Detector)
	if err != nil {
		t.Fatal(err)
	}
	response, err := a.callRequest(http.MethodPost, reqBytes, indexURL)
	if err != nil {
		t.Fatal(err)
	}
	var data map[string]interface{}
	_ = json.Unmarshal(response, &data)
	if val, ok := data["_id"]; ok {
		a.DetectorId = fmt.Sprintf("%s", val)
		return
	}
	t.Fatal(data)
}

func getRawFilter() []byte {
	return []byte(`
	{
		"bool":{
			"filter": {
				"term": {
					"currency": "EUR"
				}
			}
		}
	}`)
}

func getCreateDetectorRequest() adentity.CreateDetectorRequest {
	return adentity.CreateDetectorRequest{
		Name:        "testdata-detector",
		Description: "Test detector",
		TimeField:   "utc_time",
		Index:       []string{EcommerceIndexName},
		Features: []adentity.FeatureRequest{{
			AggregationType: []string{"sum"},
			Enabled:         true,
			Field:           []string{"total_quantity"},
		}},
		Filter:         getRawFilter(),
		Interval:       "1m",
		Delay:          "1m",
		Start:          false,
		PartitionField: nil,
	}
}

func (a *ADTestSuite) TestCreateDetectors() {
	a.T().Run("create success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := adctrl.New(os.Stdin, a.ESController, a.ADGateway)
		response, err := ctrl.CreateAnomalyDetector(ctx, a.DetectorRequest)
		assert.NoError(t, err, "failed to create detectors")
		assert.NotNil(t, response)
		a.DeleteDetectorUsingRESTAPI(t, *response)
	})
}

func (a *ADTestSuite) TestStopDetectors() {
	a.T().Run("stop success", func(t *testing.T) {
		a.StartDetectorUsingRESTAPI(t, a.DetectorId)
		ctx := context.Background()
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl := adctrl.New(&stdin, a.ESController, a.ADGateway)
		err := ctrl.StopDetectorByName(ctx, a.Detector.Name, false)
		assert.NoError(t, err, "failed to stop detectors")
	})
}

func (a *ADTestSuite) TestStartDetectors() {
	a.T().Run("start success", func(t *testing.T) {
		a.StopDetectorUsingRESTAPI(t, a.DetectorId)
		ctx := context.Background()
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl := adctrl.New(&stdin, a.ESController, a.ADGateway)
		err := ctrl.StartDetectorByName(ctx, a.Detector.Name, false)
		assert.NoError(t, err, "failed to start detectors")
	})
}
func (a *ADTestSuite) TestDeleteDetectorsForce() {
	a.T().Run("delete force success", func(t *testing.T) {
		a.StartDetectorUsingRESTAPI(t, a.DetectorId)
		ctx := context.Background()
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl := adctrl.New(&stdin, a.ESController, a.ADGateway)
		err := ctrl.DeleteDetectorByName(ctx, a.Detector.Name, true, false)
		assert.NoError(t, err, "failed to delete detectors")
	})
}

func (a *ADTestSuite) TestDeleteDetectors() {
	a.T().Run("delete stopped success", func(t *testing.T) {
		ctx := context.Background()
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl := adctrl.New(&stdin, a.ESController, a.ADGateway)
		err := ctrl.DeleteDetectorByName(ctx, a.Detector.Name, false, false)
		assert.NoError(t, err, "failed to delete detectors")
	})
}

func (a *ADTestSuite) TestGetDetectors() {
	a.T().Run("get detector success", func(t *testing.T) {
		ctx := context.Background()
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl := adctrl.New(&stdin, a.ESController, a.ADGateway)
		output, err := ctrl.GetDetectorsByName(ctx, a.Detector.Name, false)
		assert.NoError(t, err, "failed to get detectors")
		assert.EqualValues(t, 1, len(output))
		assert.EqualValues(t, a.DetectorId, output[0].ID)
	})
}

func (a *ADTestSuite) TestUpdateDetectorsForce() {
	a.T().Run("update detector success", func(t *testing.T) {
		a.StartDetectorUsingRESTAPI(t, a.DetectorId)
		ctx := context.Background()
		ctrl := adctrl.New(os.Stdin, a.ESController, a.ADGateway)
		output, err := ctrl.GetDetector(ctx, a.DetectorId)
		assert.NoError(t, err, "failed to get detector")
		updatedDetector := adentity.UpdateDetectorUserInput{
			ID:            output.ID,
			Name:          output.Name,
			Description:   output.Description,
			TimeField:     output.TimeField,
			Index:         output.Index,
			Features:      output.Features,
			Filter:        output.Filter,
			Interval:      output.Interval,
			Delay:         "5m",
			LastUpdatedAt: output.LastUpdatedAt,
			SchemaVersion: output.SchemaVersion,
		}
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl = adctrl.New(&stdin, a.ESController, a.ADGateway)
		err = ctrl.UpdateDetector(ctx, updatedDetector, true, false)
		assert.NoError(t, err, "failed to update detector")
		output, err = ctrl.GetDetector(ctx, a.DetectorId)
		assert.NoError(t, err, "failed to get detector")
		assert.EqualValues(t, "5m", output.Delay)
	})

}
func (a *ADTestSuite) TestUpdateDetectors() {
	a.T().Run("update detector success", func(t *testing.T) {
		ctx := context.Background()
		ctrl := adctrl.New(os.Stdin, a.ESController, a.ADGateway)
		output, err := ctrl.GetDetector(ctx, a.DetectorId)
		assert.NoError(t, err, "failed to get detector")
		updatedDetector := adentity.UpdateDetectorUserInput{
			ID:            output.ID,
			Name:          output.Name,
			Description:   output.Description,
			TimeField:     output.TimeField,
			Index:         output.Index,
			Features:      output.Features,
			Filter:        output.Filter,
			Interval:      output.Interval,
			Delay:         "5m",
			LastUpdatedAt: output.LastUpdatedAt,
			SchemaVersion: output.SchemaVersion,
		}
		var stdin bytes.Buffer
		stdin.Write([]byte("yes\n"))
		ctrl = adctrl.New(&stdin, a.ESController, a.ADGateway)
		err = ctrl.UpdateDetector(ctx, updatedDetector, false, false)
		assert.NoError(t, err, "failed to update detector")
		output, err = ctrl.GetDetector(ctx, a.DetectorId)
		assert.NoError(t, err, "failed to get detector")
		assert.EqualValues(t, "5m", output.Delay)
	})
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestADSuite(t *testing.T) {
	suite.Run(t, new(ADTestSuite))
}
