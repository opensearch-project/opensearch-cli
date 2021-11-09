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

package signer

import (
	"bytes"
	"errors"
	"opensearch-cli/entity"
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
