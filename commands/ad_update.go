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
	handler "opensearch-cli/handler/ad"

	"github.com/spf13/cobra"
)

const (
	updateDetectorsCommandName = "update"
	forceFlagName              = "force"
	startFlagName              = "start"
)

//updateDetectorsCmd updates detectors with configuration from input file
var updateDetectorsCmd = &cobra.Command{
	Use:   updateDetectorsCommandName + " json-file-path ... [flags]",
	Short: "Update detectors based on JSON files",
	Long: "Update detectors based on JSON files.\n" +
		"To begin, use `opensearch-cli ad get detector-name > detector_to_be_updated.json` to download the detector. " +
		"Modify the file, and then use `opensearch-cli ad update file-path` to update the detector.",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool(forceFlagName)
		start, _ := cmd.Flags().GetBool(startFlagName)
		err := updateDetectors(args, force, start)
		if err != nil {
			DisplayError(err, updateDetectorsCommandName)
		}
	},
}

func init() {
	GetADCommand().AddCommand(updateDetectorsCmd)
	updateDetectorsCmd.Flags().BoolP(forceFlagName, "f", false, "Stop detector and update forcefully")
	updateDetectorsCmd.Flags().BoolP(startFlagName, "s", false, "Start detector if update is successful")
	updateDetectorsCmd.Flags().BoolP("help", "h", false, "Help for "+updateDetectorsCommandName)
}

func updateDetectors(fileNames []string, force bool, start bool) error {
	commandHandler, err := GetADHandler()
	if err != nil {
		return err
	}
	for _, name := range fileNames {
		err = handler.UpdateAnomalyDetector(commandHandler, name, force, start)
		if err != nil {
			return err
		}
	}
	return nil
}
