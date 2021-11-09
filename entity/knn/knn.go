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

package knn

//Shards represents number of shards succeeded or failed to warmup
type Shards struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Failed     int `json:"failed"`
}

//WarmupAPIResponse warmup api response structure
type WarmupAPIResponse struct {
	Shards Shards `json:"_shards"`
}

//RootCause gives information about type and reason
type RootCause struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

//Error contains root cause
type Error struct {
	RootCause []RootCause `json:"root_cause"`
}

//ErrorResponse knn request failure error response
type ErrorResponse struct {
	KNNError Error `json:"error"`
	Status   int   `json:"status"`
}
