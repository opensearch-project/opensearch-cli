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
	"fmt"
	"opensearch-cli/handler/ad"

	"github.com/spf13/cobra"
)

const (
	startDetectorsCommandName = "start"
	stopDetectorsCommandName  = "stop"
	idFlagName                = "id"
)

//startDetectorsCmd start detectors based on id, name or name regex pattern.
//default input is name pattern, one can change this format to be id by passing --id flag
var startDetectorsCmd = &cobra.Command{
	Use:   startDetectorsCommandName + " detector_name ..." + " [flags] ",
	Short: "Start detectors based on a list of IDs, names, or name regex patterns",
	Long: "Start detectors based on a list of IDs, names, or name regex patterns.\n" +
		"Wrap regex patterns in quotation marks to prevent the terminal from matching patterns against the files in the current directory.\n" +
		"The default input is detector name. Use the `--id` flag if input is detector ID instead of name",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStatus, _ := cmd.Flags().GetBool(idFlagName)
		action := ad.StartAnomalyDetectorByNamePattern
		if idStatus {
			action = ad.StartAnomalyDetectorByID
		}
		err := execute(action, args)
		DisplayError(err, startDetectorsCommandName)
	},
}

//stopDetectorsCmd stops detectors based on id and name pattern.
//default input is name pattern, one can change this format to be id by passing --id flag
var stopDetectorsCmd = &cobra.Command{
	Use:   stopDetectorsCommandName + " detector_name ..." + " [flags] ",
	Short: "Stop detectors based on a list of IDs, names, or name regex patterns",
	Long: "Stop detectors based on a list of IDs, names, or name regex patterns.\n" +
		"Wrap regex patterns in quotation marks to prevent the terminal from matching patterns against the files in the current directory.\n" +
		"The default input is detector name. Use the `--id` flag if input is detector ID instead of name",
	Run: func(cmd *cobra.Command, args []string) {
		//If no args, display usage
		if len(args) < 1 {
			fmt.Println(cmd.Usage())
			return
		}
		idStatus, _ := cmd.Flags().GetBool(idFlagName)
		action := ad.StopAnomalyDetectorByNamePattern
		if idStatus {
			action = ad.StopAnomalyDetectorByID
		}
		err := execute(action, args)
		DisplayError(err, stopDetectorsCommandName)
	},
}

func init() {
	startDetectorsCmd.Flags().BoolP(idFlagName, "", false, "Input is detector ID")
	startDetectorsCmd.Flags().BoolP("help", "h", false, "Help for "+startDetectorsCommandName)
	GetADCommand().AddCommand(startDetectorsCmd)
	stopDetectorsCmd.Flags().BoolP(idFlagName, "", false, "Input is detector ID")
	stopDetectorsCmd.Flags().BoolP("help", "h", false, "Help for "+stopDetectorsCommandName)
	GetADCommand().AddCommand(stopDetectorsCmd)
}

func execute(f func(*ad.Handler, string) error, detectors []string) error {
	// iterate over the arguments
	// the first return value is index of fileNames, we can omit it using _
	commandHandler, err := GetADHandler()
	if err != nil {
		return err
	}
	for _, detector := range detectors {
		err := f(commandHandler, detector)
		if err != nil {
			return err
		}
	}
	return nil
}
