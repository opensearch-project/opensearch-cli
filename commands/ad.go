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
	g, err := adgateway.New(c, profile)
	if err != nil {
		return nil, err
	}
	esg, err := gateway.New(c, profile)
	if err != nil {
		return nil, err
	}
	esc := ctrl.New(esg)
	ctr := adctrl.New(os.Stdin, esc, g)
	return handler.New(ctr), nil
}
