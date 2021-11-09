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
	deleteDetectorsCommandName    = "delete"
	deleteDetectorIDFlagName      = "id"
	detectorForceDeletionFlagName = "force"
)

//deleteDetectorsCmd deletes detectors based on id, name or name regex pattern.
//default input is name pattern, one can change this format to be id by passing --id flag
var deleteDetectorsCmd = &cobra.Command{
	Use:   deleteDetectorsCommandName + " detector_name ..." + " [flags] ",
	Short: "Delete detectors based on a list of IDs, names, or name regex patterns",
	Long: "Delete detectors based on list of IDs, names, or name regex patterns.\n" +
		"Wrap regex patterns in quotation marks to prevent the terminal from matching patterns against the files in the current directory.\nThe default input is detector name. Use the `--id` flag if input is detector ID instead of name",

	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool(detectorForceDeletionFlagName)
		detectorID, _ := cmd.Flags().GetBool(deleteDetectorIDFlagName)
		action := handler.DeleteAnomalyDetectorByNamePattern
		if detectorID {
			action = handler.DeleteAnomalyDetectorByID
		}
		err := deleteDetectors(args, force, action)
		DisplayError(err, deleteDetectorsCommandName)
	},
}

func init() {
	GetADCommand().AddCommand(deleteDetectorsCmd)
	deleteDetectorsCmd.Flags().BoolP(detectorForceDeletionFlagName, "f", false, "Delete the detector even if it is running")
	deleteDetectorsCmd.Flags().BoolP(deleteDetectorIDFlagName, "", false, "Input is detector ID")
	deleteDetectorsCmd.Flags().BoolP("help", "h", false, "Help for "+deleteDetectorsCommandName)
}

//deleteDetectors deletes detectors with force by calling delete method provided
func deleteDetectors(detectors []string, force bool, f func(*handler.Handler, string, bool) error) error {
	commandHandler, err := GetADHandler()
	if err != nil {
		return err
	}
	for _, detector := range detectors {
		err = f(commandHandler, detector, force)
		if err != nil {
			return err
		}
	}
	return nil
}
