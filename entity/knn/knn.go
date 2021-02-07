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
