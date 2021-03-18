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

package signer

import (
	"bytes"
	"errors"
	"odfe-cli/entity"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/hashicorp/go-retryablehttp"
)

func GetV4Signer(credentials *credentials.Credentials) *v4.Signer {
	return v4.NewSigner(credentials)
}
func sign(req *retryablehttp.Request, region *string, serviceName string, signer *v4.Signer) error {
	bodyBytes, err := req.BodyBytes()
	if err != nil {
		return err
	}
	if region == nil || len(*region) == 0 {
		return errors.New("aws region is not found. Either set 'AWS_REGION' or add this information during aws profile creation step")
	}
	// Sign the request
	_, err = signer.Sign(req.Request, bytes.NewReader(bodyBytes), serviceName, *region, time.Now())
	return err
}

//SignRequest signs the request using SigV4
func SignRequest(req *retryablehttp.Request, awsProfile entity.AWSIAM, getSigner func(*credentials.Credentials) *v4.Signer) error {
	awsSession, err := session.NewSessionWithOptions(session.Options{
		Profile:           awsProfile.ProfileName,
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return err
	}
	signer := getSigner(awsSession.Config.Credentials)
	return sign(req, awsSession.Config.Region, awsProfile.ServiceName, signer)
}
