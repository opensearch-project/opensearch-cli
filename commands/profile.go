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
	"fmt"

	"golang.org/x/term"

	"odfe-cli/controller/config"
	"odfe-cli/controller/profile"
	"odfe-cli/entity"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

const (
	CreateNewProfileCommandName = "create"
	DeleteProfilesCommandName   = "delete"
	FlagProfileVerbose          = "verbose"
	ListProfilesCommandName     = "list"
	ProfileCommandName          = "profile"
	padding                     = 3
	alignLeft                   = 0
)

//GetProfileController gets controller based on config file
func GetProfileController() (profile.Controller, error) {
	cfgFile, err := GetRoot().Flags().GetString(flagConfig)
	if err != nil {
		return nil, err
	}
	return getProfileController(cfgFile)
}

//profileCommand is main command for profile operations like list, create and delete
var profileCommand = &cobra.Command{
	Use:   ProfileCommandName + " sub-command",
	Short: "Manage a collection of settings and credentials that you can apply to an odfe-cli command",
	Long: "A named profile is a collection of settings and credentials that you can apply to an odfe-cli command. " +
		"When you specify a profile for a command (e.g. `odfe-cli <command> --profile <profile_name>`), odfe-cli uses " +
		"the profile's settings and credentials to run the given command.\n" +
		"To configure a default profile for commands, either specify the default profile name in an environment " +
		"variable (`ODFE_PROFILE`) or create a profile named `default`.",
}

//createProfileCmd creates profile interactively by prompting for name (distinct), user, endpoint, password.
var createProfileCmd = &cobra.Command{
	Use:   CreateNewProfileCommandName,
	Short: "Create profile",
	Long:  "Create profile with following fields: name, endpoint, user and password",
	Run: func(cmd *cobra.Command, args []string) {
		profileController, err := GetProfileController()
		if err != nil {
			DisplayError(err, CreateNewProfileCommandName)
			return
		}
		err = CreateProfile(profileController, getNewProfile)
		if err != nil {
			DisplayError(err, CreateNewProfileCommandName)
			return
		}
	},
}

//deleteProfilesCmd deletes profiles by names
var deleteProfilesCmd = &cobra.Command{
	Use:   DeleteProfilesCommandName + " profile_name ...",
	Short: "Delete profiles by names",
	Long:  "Delete profiles by names from the config file permanently.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			fmt.Println(cmd.Usage())
			return
		}
		if err := deleteProfiles(args); err != nil {
			DisplayError(err, DeleteProfilesCommandName)
			return
		}
	},
}

//listProfileCmd lists profiles by names
var listProfileCmd = &cobra.Command{
	Use:   ListProfilesCommandName,
	Short: "List profiles from the config file",
	Long:  "List profiles from the config file",
	Run: func(cmd *cobra.Command, args []string) {
		if err := listProfiles(cmd); err != nil {
			DisplayError(err, ListProfilesCommandName)
			return
		}
	},
}

//deleteProfiles deletes profiles based on names
func deleteProfiles(profiles []string) error {
	profileController, err := GetProfileController()
	if err != nil {
		return err
	}
	return profileController.DeleteProfiles(profiles)
}

// init to register commands to its parent command to create a hierarchy
func init() {
	profileCommand.AddCommand(createProfileCmd)
	profileCommand.AddCommand(deleteProfilesCmd)
	profileCommand.AddCommand(listProfileCmd)
	listProfileCmd.Flags().BoolP(FlagProfileVerbose, "l", false, "Shows information like name, endpoint, user")
	listProfileCmd.Flags().BoolP("help", "h", false, "Help for "+ListProfilesCommandName)
	createProfileCmd.Flags().BoolP("help", "h", false, "Help for "+CreateNewProfileCommandName)
	deleteProfilesCmd.Flags().BoolP("help", "h", false, "Help for "+DeleteProfilesCommandName)
	profileCommand.Flags().BoolP("help", "h", false, "Help for "+ProfileCommandName)
	GetRoot().AddCommand(profileCommand)
}

//getProfileController gets profile controller by wiring config controller with config file
func getProfileController(cfgFlagValue string) (profile.Controller, error) {
	configFilePath, err := GetConfigFilePath(cfgFlagValue)
	if err != nil {
		return nil, fmt.Errorf("failed to get config file due to: %w", err)
	}
	configController := config.New(configFilePath)
	profileController := profile.New(configController)
	return profileController, nil
}

// CreateProfile creates a new named profile
func CreateProfile(profileController profile.Controller, getNewProfile func(map[string]entity.Profile) entity.Profile) error {
	profiles, err := profileController.GetProfilesMap()
	if err != nil {
		return fmt.Errorf("failed to get profile names due to: %w", err)
	}
	newProfile := getNewProfile(profiles)
	if err = profileController.CreateProfile(newProfile); err != nil {
		return fmt.Errorf("failed to create profile %v due to: %w", newProfile, err)
	}
	return nil
}

// getNewProfile gets new profile information from user using command line
func getNewProfile(profileMap map[string]entity.Profile) entity.Profile {
	var name string
	for {
		fmt.Printf("Enter profile's name: ")
		name = getUserInputAsText(checkInputIsNotEmpty)
		if _, ok := profileMap[name]; !ok {
			break
		}
		fmt.Println("profile", name, "already exists.")
	}
	fmt.Printf("Elasticsearch Endpoint: ")
	endpoint := getUserInputAsText(checkInputIsNotEmpty)
	fmt.Printf("Is security plugin enabled? Y/N ")
	isSecured := getUserInputAsBoolean(checkInputIsNotEmpty)
	var user, password string
	if isSecured {
		fmt.Printf("User Name: ")
		user = getUserInputAsText(checkInputIsNotEmpty)
		fmt.Printf("Password: ")
		password = getUserInputAsMaskedText(checkInputIsNotEmpty)
	}
	return entity.Profile{
		Name:     name,
		Endpoint: endpoint,
		UserName: user,
		Password: password,
	}
}

func getUserInputAsBoolean(isValid func(input string) bool) bool {
	response := getUserInputAsText(isValid)
	switch strings.ToLower(response) {
	case "y", "yes":
		return true
	case "n", "no":
		return false
	default:
		fmt.Printf("please type (y)es or (n)o: ")
		return getUserInputAsBoolean(isValid)
	}
}

// getUserInputAsText get value from user as text
func getUserInputAsText(isValid func(string) bool) string {
	var response string
	//Ignore return value since validation is applied below
	_, _ = fmt.Scanln(&response)
	if !isValid(response) {
		return getUserInputAsText(isValid)
	}
	return strings.TrimSpace(response)
}

// checkInputIsNotEmpty checks whether input is empty or not
func checkInputIsNotEmpty(input string) bool {
	if len(input) < 1 {
		fmt.Print("Value cannot be empty, please enter non-empty value: ")
		return false
	}
	return true
}

// getUserInputAsMaskedText get value from user as masked text, since credentials like password
// should not be displayed on console for security reasons
func getUserInputAsMaskedText(isValid func(string) bool) string {
	maskedValue, err := term.ReadPassword(0)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	value := fmt.Sprintf("%s", maskedValue)
	if !isValid(value) {
		return getUserInputAsMaskedText(isValid)
	}
	fmt.Println()
	return value
}

//listProfiles list profiles from the config file
func listProfiles(cmd *cobra.Command) error {
	ok, err := cmd.Flags().GetBool(FlagProfileVerbose)
	if err != nil {
		return err
	}
	profileController, err := GetProfileController()
	if err != nil {
		return err
	}
	if !ok {
		return displayProfileNames(profileController)
	}
	return displayCompleteProfiles(profileController)
}

//displayCompleteProfiles lists complete profile information as below
/*
Name       UserName     Endpoint-url
----       --------     ------------
default    admin      	https://localhost:9200
dev        test      	https://127.0.0.1:9200
*/
func displayCompleteProfiles(p profile.Controller) (err error) {
	var profiles []entity.Profile
	if profiles, err = p.GetProfiles(); err != nil {
		return
	}
	if len(profiles) < 1 {
		return fmt.Errorf("no profiles found")
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, padding, ' ', alignLeft)
	defer func() {
		err = w.Flush()
	}()
	_, err = fmt.Fprintln(w, "Name\t\tUserName\t\tEndpoint-url\t")
	_, err = fmt.Fprintf(w, "%s\t\t%s\t\t%s\t\n", "----", "--------", "------------")
	for _, p := range profiles {
		_, err = fmt.Fprintf(w, "%s\t\t%s\t\t%s\t\n", p.Name, p.UserName, p.Endpoint)
	}
	return
}

//displayProfileNames lists only profile names
func displayProfileNames(p profile.Controller) (err error) {

	var names []string
	if names, err = p.GetProfileNames(); err != nil {
		return
	}
	if len(names) < 1 {
		return fmt.Errorf("no profiles found")
	}
	for _, name := range names {
		fmt.Println(name)
	}
	return nil
}
