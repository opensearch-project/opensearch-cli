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

package gateway

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"opensearch-cli/client"
	"opensearch-cli/entity"
	"opensearch-cli/entity/platform"
	"opensearch-cli/environment"
	"opensearch-cli/gateway/aws/signer"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// HTTPGateway type for gateway client
type HTTPGateway struct {
	Client  *client.Client
	Profile *entity.Profile
}

// GetDefaultHeaders returns common headers
func GetDefaultHeaders() map[string]string {
	return map[string]string{
		"content-type": "application/json",
	}
}
func GetTLSConfig(trust *entity.Trust) (*tls.Config, error) {
	config := &tls.Config{}
	if trust.ClientCertificateFilePath != nil && trust.ClientKeyFilePath != nil {
		cert, err := tls.LoadX509KeyPair(*trust.ClientCertificateFilePath, *trust.ClientKeyFilePath)
		if err != nil {
			return nil, fmt.Errorf(
				"error creating x509 keypair from client cert file %s and client key file %s",
				*trust.ClientCertificateFilePath, *trust.ClientKeyFilePath)
		}
		config.Certificates = []tls.Certificate{cert}
	}
	caCertPool := x509.NewCertPool()
	if trust.CAFilePath != nil {
		caCert, err := ioutil.ReadFile(*trust.CAFilePath)
		if err != nil {
			return nil, fmt.Errorf("error opening certificate file %s, error: %s", *trust.CAFilePath, err)
		}
		caCertPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caCertPool
	}

	return config, nil
}

// NewHTTPGateway creates new HTTPGateway instance
func NewHTTPGateway(c *client.Client, p *entity.Profile) (*HTTPGateway, error) {

	if p.Certificate != nil {
		tlsConfig, err := GetTLSConfig(p.Certificate)
		if err != nil {
			return nil, err
		}
		c.HTTPClient.HTTPClient.Transport = &http.Transport{
			TLSClientConfig: tlsConfig,
		}
	}

	// set max retry if provided by command
	if p.MaxRetry != nil {
		c.HTTPClient.RetryMax = *p.MaxRetry
	}
	//override with environment variable if exists
	if val, ok := overrideValue(p, environment.OPENSEARCH_MAX_RETRY); ok {
		c.HTTPClient.RetryMax = *val
	}

	// set connection timeout if provided by command
	if p.Timeout != nil {
		c.HTTPClient.HTTPClient.Timeout = time.Duration(*p.Timeout) * time.Second
	}
	//override with environment variable if exists
	if duration, ok := overrideValue(p, environment.OPENSEARCH_TIMEOUT); ok {
		c.HTTPClient.HTTPClient.Timeout = time.Duration(*duration) * time.Second
	}

	return &HTTPGateway{
		Client:  c,
		Profile: p,
	}, nil
}

func overrideValue(p *entity.Profile, envVariable string) (*int, bool) {
	if val, ok := os.LookupEnv(envVariable); ok {
		//ignore error from non positive number
		if attempt, err := strconv.Atoi(val); err == nil {
			return &attempt, true
		}
	}
	return nil, false
}

// isValidResponse checks whether the response is valid or not by checking the status code
func (g *HTTPGateway) isValidResponse(response *http.Response) error {
	if response == nil {
		return errors.New("response is nil")
	}
	// client error if 400 <= status code < 500
	if response.StatusCode >= http.StatusBadRequest && response.StatusCode < http.StatusInternalServerError {

		return platform.NewRequestError(
			response.StatusCode,
			response.Body,
			fmt.Errorf("%d Client Error: %s for url: %s", response.StatusCode, response.Status, response.Request.URL))
	}
	// server error if status code >= 500
	if response.StatusCode >= http.StatusInternalServerError {

		return platform.NewRequestError(
			response.StatusCode,
			response.Body,
			fmt.Errorf("%d Server Error: %s for url: %s", response.StatusCode, response.Status, response.Request.URL))
	}
	return nil
}

// Execute calls request using http and check if status code is ok or not
func (g *HTTPGateway) Execute(req *retryablehttp.Request) ([]byte, error) {
	if g.Profile.AWS != nil {
		//sign request
		if err := signer.SignRequest(req, *g.Profile.AWS, signer.GetV4Signer); err != nil {
			return nil, err
		}
	}
	response, err := g.Client.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := response.Body.Close()
		if err != nil {
			return
		}
	}()
	if err = g.isValidResponse(response); err != nil {
		return nil, err
	}
	return ioutil.ReadAll(response.Body)
}

// Call calls request using http and return error if status code is not expected
func (g *HTTPGateway) Call(req *retryablehttp.Request, statusCode int) ([]byte, error) {
	resBytes, err := g.Execute(req)
	if err == nil {
		return resBytes, nil
	}
	r, ok := err.(*platform.RequestError)
	if !ok {
		return nil, err
	}
	if r.StatusCode() != statusCode {
		return nil, fmt.Errorf(r.GetResponse())
	}
	return nil, err

}

// BuildRequest builds request based on method and appends payload for given url with headers
// TODO: Deprecate this method by replace this with BuildCurlRequest
func (g *HTTPGateway) BuildRequest(ctx context.Context, method string, payload interface{}, url string, headers map[string]string) (*retryablehttp.Request, error) {
	if payload == nil {
		return g.BuildCurlRequest(ctx, method, nil, url, headers)
	}
	reqBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return g.BuildCurlRequest(ctx, method, reqBytes, url, headers)
}

// BuildCurlRequest builds request based on method and add payload (in byte)
func (g *HTTPGateway) BuildCurlRequest(ctx context.Context, method string, payload []byte, url string, headers map[string]string) (*retryablehttp.Request, error) {
	r, err := retryablehttp.NewRequest(method, url, payload)
	if err != nil {
		return nil, err
	}
	req := r.WithContext(ctx)
	if len(g.Profile.UserName) != 0 {
		req.SetBasicAuth(g.Profile.UserName, g.Profile.Password)
	}
	if len(headers) == 0 {
		return req, nil
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil
}

// BuildCurlMultipartFormRequest builds multipart file-upload request based on method and add payload (in byte)
func (g *HTTPGateway) BuildCurlMultipartFormRequest(ctx context.Context, method string, filePath string, url string, headers map[string]string) (*retryablehttp.Request, error) {
	file, _ := os.Open(filePath)
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filePath)
	io.Copy(part, file)
	writer.Close()

	r, err := retryablehttp.NewRequest(method, url, body)

	r.Header.Add("Content-Type", writer.FormDataContentType())

	if err != nil {
		return nil, err
	}
	req := r.WithContext(ctx)
	if len(g.Profile.UserName) != 0 {
		req.SetBasicAuth(g.Profile.UserName, g.Profile.Password)
	}
	if len(headers) == 0 {
		return req, nil
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return req, nil

}

// GetValidEndpoint get url based on user config
func GetValidEndpoint(profile *entity.Profile) (*url.URL, error) {
	u, err := url.ParseRequestURI(profile.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint: %v due to %v", profile.Endpoint, err)
	}
	return u, nil
}
