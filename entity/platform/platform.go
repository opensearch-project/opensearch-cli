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

package platform

//Terms contains fields
type Terms struct {
	Field string `json:"field"`
}

//DistinctGroups contains terms
type DistinctGroups struct {
	Term Terms `json:"terms"`
}

//Aggregate contains list of items
type Aggregate struct {
	Group DistinctGroups `json:"items"`
}

//SearchRequest structure for request
type SearchRequest struct {
	Agg  Aggregate `json:"aggs"`
	Size int32     `json:"size"`
}

//Bucket represents bucket used by ES for aggregations
type Bucket struct {
	Key      interface{} `json:"key"`
	DocCount int64       `json:"doc_count"`
}

//Items contains buckets defined by response
type Items struct {
	Buckets []Bucket `json:"buckets"`
}

//Aggregations contains items defined by response
type Aggregations struct {
	Items Items `json:"items"`
}

//Response response defined by response
type Response struct {
	Aggregations Aggregations `json:"aggregations"`
}

//CurlRequest contains parameter to execute REST Action
type CurlRequest struct {
	Action      string
	Path        string
	QueryParams string
	Headers     map[string]string
	Data        []byte
}

//CurlCommandRequest contains parameter from command
type CurlCommandRequest struct {
	Action           string
	Path             string
	QueryParams      string
	Headers          string
	Data             string
	Pretty           bool
	OutputFormat     string
	OutputFilterPath string
}
