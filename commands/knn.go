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
	"opensearch-cli/client"
	ctrl "opensearch-cli/controller/knn"
	gateway "opensearch-cli/gateway/knn"
	handler "opensearch-cli/handler/knn"

	"github.com/spf13/cobra"
)

const (
	knnCommandName        = "knn"
	knnStatsCommandName   = "stats"
	knnWarmupCommandName  = "warmup"
	knnStatsNodesFlagName = "nodes"
	knnStatsNamesFlagName = "stat-names"
)

//knnCommand is base command for k-NN plugin.
var knnCommand = &cobra.Command{
	Use:   knnCommandName,
	Short: "Manage the k-NN plugin",
	Long:  "Use the k-NN commands to perform operations like stats, warmup.",
}

//knnStatsCommandName provide stats command for k-NN plugin.
var knnStatsCommand = &cobra.Command{
	Use:   knnStatsCommandName,
	Short: "Display current status of the k-NN Plugin",
	Long:  "Display current status of the k-NN Plugin.",
	Run: func(cmd *cobra.Command, args []string) {
		h, err := GetKNNHandler()
		if err != nil {
			DisplayError(err, knnStatsCommandName)
			return
		}
		nodes, err := cmd.Flags().GetString(knnStatsNodesFlagName)
		if err != nil {
			DisplayError(err, knnStatsCommandName)
			return
		}
		names, err := cmd.Flags().GetString(knnStatsNamesFlagName)
		if err != nil {
			DisplayError(err, knnStatsCommandName)
			return
		}
		err = getStatistics(h, nodes, names)
		DisplayError(err, knnStatsCommandName)
	},
}

//knnWarmupCommand warmups shards
var knnWarmupCommand = &cobra.Command{
	Use:   knnWarmupCommandName + " index ..." + " [flags] ",
	Args:  cobra.MinimumNArgs(1),
	Short: "Warmup shards for given indices",
	Long: "Warmup command loads all graphs for all of the shards (primaries and replicas) " +
		"for given indices into native memory.\nThis is an asynchronous operation. If the command times out, " +
		"the operation will still be going on in the cluster.\nTo monitor this, use the OpenSearch _tasks API. " +
		"Use `opensearch-cli knn stats` command to verify whether indices are successfully loaded into memory.",
	Run: func(cmd *cobra.Command, args []string) {
		h, err := GetKNNHandler()
		if err != nil {
			DisplayError(err, knnWarmupCommandName)
			return
		}
		err = warmupIndices(h, args)
		DisplayError(err, knnWarmupCommandName)
	},
}

func GetKNNCommand() *cobra.Command {
	return knnCommand
}

func GetKNNStatsCommand() *cobra.Command {
	return knnStatsCommand
}

func GetKNNWarmupCommand() *cobra.Command {
	return knnWarmupCommand
}

func init() {
	//knn base command
	knnCommand.Flags().BoolP("help", "h", false, "Help for k-NN plugin")
	GetRoot().AddCommand(knnCommand)
	//knn stats command
	knnStatsCommand.Flags().BoolP("help", "h", false, "Help for k-NN plugin stats command")
	knnStatsCommand.Flags().StringP(knnStatsNodesFlagName, "n", "", "Input is list of node Ids, separated by ','")
	knnStatsCommand.Flags().StringP(knnStatsNamesFlagName, "s", "", "Input is list of stats names, separated by ','")
	knnCommand.AddCommand(knnStatsCommand)
	//knn warmup command
	knnWarmupCommand.Flags().BoolP("help", "h", false, "Help for k-NN plugin warmup command")
	knnCommand.AddCommand(knnWarmupCommand)
}

func getStatistics(h *handler.Handler, nodes string, names string) error {
	stats, err := handler.GetStatistics(h, nodes, names)
	if err != nil {
		return err
	}
	fmt.Println(string(stats))
	return nil
}

func warmupIndices(h *handler.Handler, index []string) error {
	shards, err := handler.WarmupIndices(h, index)
	if err != nil {
		return err
	}
	if shards.Failed > 0 {
		return fmt.Errorf("%d/%d shards were failed to load into memory", shards.Failed, shards.Total)
	}
	fmt.Printf("successfully loaded %d shards into memory\n", shards.Total)
	return nil
}

//GetKNNHandler returns handler by wiring the dependency manually
func GetKNNHandler() (*handler.Handler, error) {
	c, err := client.New(nil)
	if err != nil {
		return nil, err
	}
	profile, err := GetProfile()
	if err != nil {
		return nil, err
	}
	g, err := gateway.New(c, profile)
	if err != nil {
		return nil, err
	}
	ctr := ctrl.New(g)
	return handler.New(ctr), nil
}
