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
	"es-cli/odfe-cli/handler/ad"
	"fmt"

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
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(cmd.Usage())
			return
		}
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
	startDetectorsCmd.Flags().StringP(flagProfileName, "p", "", "Use a specific profile from your configuration file")
	startDetectorsCmd.Flags().BoolP("help", "h", false, "Help for "+startDetectorsCommandName)
	GetADCommand().AddCommand(startDetectorsCmd)
	stopDetectorsCmd.Flags().BoolP(idFlagName, "", false, "Input is detector ID")
	stopDetectorsCmd.Flags().StringP(flagProfileName, "p", "", "Use a specific profile from your configuration file")
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
