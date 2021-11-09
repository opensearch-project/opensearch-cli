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
