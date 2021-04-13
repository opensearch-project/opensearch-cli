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

package commands

import (
	"opensearch-cli/client"
	adctrl "opensearch-cli/controller/ad"
	ctrl "opensearch-cli/controller/platform"
	adgateway "opensearch-cli/gateway/ad"
	gateway "opensearch-cli/gateway/platform"
	handler "opensearch-cli/handler/ad"
	"os"

	"github.com/spf13/cobra"
)

const (
	adCommandName = "ad"
)

//adCommand is base command for Anomaly Detection plugin.
var adCommand = &cobra.Command{
	Use:   adCommandName,
	Short: "Manage the Anomaly Detection plugin",
	Long:  "Use the Anomaly Detection commands to create, configure, and manage detectors.",
}

func init() {
	adCommand.Flags().BoolP("help", "h", false, "Help for Anomaly Detection")
	GetRoot().AddCommand(adCommand)
}

//GetADCommand returns AD base command, since this will be needed for subcommands
//to add as parent later
func GetADCommand() *cobra.Command {
	return adCommand
}

//GetADHandler returns handler by wiring the dependency manually
func GetADHandler() (*handler.Handler, error) {
	c, err := client.New(nil)
	if err != nil {
		return nil, err
	}
	profile, err := GetProfile()
	if err != nil {
		return nil, err
	}
	g := adgateway.New(c, profile)
	esg := gateway.New(c, profile)
	esc := ctrl.New(esg)
	ctr := adctrl.New(os.Stdin, esc, g)
	return handler.New(ctr), nil
}
