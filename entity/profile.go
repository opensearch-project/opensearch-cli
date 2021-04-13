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

package entity

type AWSIAM struct {
	ProfileName string `yaml:"profile"`
	ServiceName string `yaml:"service"`
}

//Trust contains file path for certificate and private key locations
type Trust struct {
	CAFilePath                *string
	ClientCertificateFilePath *string
	ClientKeyFilePath         *string
}

type Profile struct {
	Name        string  `yaml:"name"`
	Endpoint    string  `yaml:"endpoint"`
	UserName    string  `yaml:"user,omitempty"`
	Password    string  `yaml:"password,omitempty"`
	AWS         *AWSIAM `yaml:"aws_iam,omitempty"`
	Certificate *Trust  `yaml:"certificate,omitempty"`
	MaxRetry    *int    `yaml:"max_retry,omitempty"`
	Timeout     *int64  `yaml:"timeout,omitempty"`
}
