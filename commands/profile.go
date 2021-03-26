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
	"errors"
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
	FlagProfileCreateName       = "name"
	FlagProfileCreateEndpoint   = "endpoint"
	FlagProfileCreateAuthType   = "auth-type"
	FlagProfileMaxRetry         = "max-retry"
	FlagProfileTimeout          = "timeout"
	FlagProfileHelp             = "help"
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
	Long:  "Create named profile to save settings and credentials that you can apply to an odfe-cli command.",
	Run: func(cmd *cobra.Command, args []string) {
		profileController, err := GetProfileController()
		if err != nil {
			DisplayError(err, CreateNewProfileCommandName)
			return
		}
		name, err := getProfileName(cmd, profileController)
		if err != nil {
			DisplayError(err, CreateNewProfileCommandName)
			return
		}
		endpoint, _ := cmd.Flags().GetString(FlagProfileCreateEndpoint)
		maxAttempt, _ := cmd.Flags().GetInt(FlagProfileMaxRetry)
		timeout, _ := cmd.Flags().GetInt64(FlagProfileTimeout)
		newProfile := entity.Profile{
			Name:     name,
			Endpoint: endpoint,
			MaxRetry: &maxAttempt,
			Timeout:  &timeout,
		}
		switch authType, _ := cmd.Flags().GetString(FlagProfileCreateAuthType); authType {
		case "disabled":
			break
		case "basic":
			getBasicAuthDetails(&newProfile)
		case "aws-iam":
			getAWSIAMAuthDetails(&newProfile)
		default:
			DisplayError(errors.New("invalid value for auth-type. Use --help -h command to see permitted values"), CreateNewProfileCommandName)
			return
		}
		err = CreateProfile(profileController, newProfile)
		if err != nil {
			DisplayError(err, CreateNewProfileCommandName)
			return
		}
		fmt.Println("Profile created successfully.")
	},
}

func getProfileName(cmd *cobra.Command, controller profile.Controller) (string, error) {
	name, _ := cmd.Flags().GetString(FlagProfileCreateName)
	if err := validateProfileName(name, controller); err != nil {
		return "", err
	}
	return name, nil
}

//deleteProfilesCmd deletes profiles by names
var deleteProfilesCmd = &cobra.Command{
	Use:   DeleteProfilesCommandName + " profile_name ...",
	Short: "Delete profiles by names",
	Long:  "Delete profiles by names from the config file permanently.",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := deleteProfiles(args); err != nil {
			DisplayError(err, DeleteProfilesCommandName)
			return
		}
		fmt.Println("Profile deleted successfully.")
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

	//profile flags
	profileCommand.Flags().BoolP(FlagProfileHelp, "h", false, "Help for "+ProfileCommandName)

	//profile list flags
	listProfileCmd.Flags().BoolP(FlagProfileVerbose, "l", false, "Shows information like name, endpoint, user")
	listProfileCmd.Flags().BoolP(FlagProfileHelp, "h", false, "Help for "+ListProfilesCommandName)

	//profile create flags
	createProfileCmd.Flags().StringP(FlagProfileCreateName, "n", "", "Create profile with this name")
	_ = createProfileCmd.MarkFlagRequired(FlagProfileCreateName)
	createProfileCmd.Flags().StringP(FlagProfileCreateEndpoint, "e", "", "Create profile with this endpoint or host")
	_ = createProfileCmd.MarkFlagRequired(FlagProfileCreateEndpoint)
	createProfileCmd.Flags().StringP(FlagProfileCreateAuthType, "a", "", "Authentication type. Options are disabled, basic and aws-iam."+
		"\nIf security is disabled, provide --auth-type='disabled'.\nIf security uses HTTP basic authentication, provide --auth-type='basic'.\n"+
		"If security uses AWS IAM ARNs as users, provide --auth-type='aws-iam'.\nodfe-cli asks for additional information based on your choice of authentication type.")
	_ = createProfileCmd.MarkFlagRequired(FlagProfileCreateAuthType)
	createProfileCmd.Flags().IntP(FlagProfileMaxRetry, "m", 3, "Maximum retry attempts allowed if transient problems occur.\n"+
		"You can override this value by using the ODFE_MAX_RETRY environment variable.")
	createProfileCmd.Flags().Int64P(FlagProfileTimeout, "t", 10, "Maximum time allowed for connection in seconds.\n"+
		"You can override this value by using the ODFE_TIMEOUT environment variable.")
	createProfileCmd.Flags().BoolP(FlagProfileHelp, "h", false, "Help for "+CreateNewProfileCommandName)

	//profile delete flags
	deleteProfilesCmd.Flags().BoolP(FlagProfileHelp, "h", false, "Help for "+DeleteProfilesCommandName)

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
func CreateProfile(profileController profile.Controller, newProfile entity.Profile) error {
	if err := profileController.CreateProfile(newProfile); err != nil {
		return fmt.Errorf("failed to create profile %v due to: %w", newProfile, err)
	}
	return nil
}

func validateProfileName(name string, controller profile.Controller) error {
	profileMap, err := controller.GetProfilesMap()
	if err != nil {
		return err
	}
	if _, ok := profileMap[name]; !ok {
		return nil
	}
	return fmt.Errorf("profile %s already exists", name)
}

// getBasicAuthDetails gets new basic HTTP Auth profile information from user using command line
func getBasicAuthDetails(newProfile *entity.Profile) {
	fmt.Printf("Username: ")
	newProfile.UserName = getUserInputAsText(checkInputIsNotEmpty)
	fmt.Printf("Password: ")
	newProfile.Password = getUserInputAsMaskedText(checkInputIsNotEmpty)
}

// getAWSIAMAuthDetails gets new AWS IAM Auth profile information from user using command line
func getAWSIAMAuthDetails(newProfile *entity.Profile) {
	fmt.Printf("AWS profile name (leave blank if you want to provide credentials using environment variables): ")
	awsIAM := &entity.AWSIAM{}
	awsIAM.ProfileName = getUserInputAsText(nil)
	fmt.Printf("AWS service name where your cluster is deployed (for Amazon Elasticsearch Service, use 'es'. For EC2, use 'ec2'): ")
	awsIAM.ServiceName = getUserInputAsText(checkInputIsNotEmpty)
	newProfile.AWS = awsIAM
}

// getUserInputAsText get value from user as text
func getUserInputAsText(isValid func(string) bool) string {
	var response string
	//Ignore return value since validation is applied below
	_, _ = fmt.Scanln(&response)
	if isValid != nil && !isValid(response) {
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
