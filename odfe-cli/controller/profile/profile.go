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

package profile

import (
	"es-cli/odfe-cli/controller/config"
	"es-cli/odfe-cli/entity"
	"fmt"
	"os"
	"strings"
)

const (
	odfeProfileEnvVarName  = "ODFE_PROFILE"
	odfeDefaultProfileName = "default"
)

//go:generate go run -mod=mod github.com/golang/mock/mockgen -destination=mocks/mock_profile.go -package=mocks . Controller
type Controller interface {
	CreateProfile(profile entity.Profile) error
	DeleteProfiles(names []string) error
	GetProfiles() ([]entity.Profile, error)
	GetProfileNames() ([]string, error)
	GetProfilesMap() (map[string]entity.Profile, error)
	GetProfileForExecution(name string) (entity.Profile, bool, error)
}

type controller struct {
	configCtrl config.Controller
}

//New returns new config controller instance
func New(c config.Controller) Controller {
	return &controller{
		configCtrl: c,
	}
}

//GetProfiles gets list of profiles fom config file
func (c controller) GetProfiles() ([]entity.Profile, error) {
	data, err := c.configCtrl.Read()
	if err != nil {
		return nil, err
	}
	return data.Profiles, nil
}

//GetProfileNames gets list of profile names
func (c controller) GetProfileNames() ([]string, error) {
	profiles, err := c.GetProfiles()
	if err != nil {
		return nil, err
	}
	var names []string
	for _, profile := range profiles {
		names = append(names, profile.Name)
	}
	return names, nil
}

//GetProfilesMap returns a map view of the profiles contained in config
func (c controller) GetProfilesMap() (map[string]entity.Profile, error) {
	profiles, err := c.GetProfiles()
	if err != nil {
		return nil, err
	}
	result := make(map[string]entity.Profile)
	for _, p := range profiles {
		result[p.Name] = p
	}
	return result, nil
}

//CreateProfile creates profile by gets list of existing profiles, append new profile to list
//and saves it in config file
func (c controller) CreateProfile(p entity.Profile) error {
	data, err := c.configCtrl.Read()
	if err != nil {
		return err
	}
	data.Profiles = append(data.Profiles, p)
	return c.configCtrl.Write(data)
}

//DeleteProfiles loads all profile, deletes selected profiles, and saves rest in config file
func (c controller) DeleteProfiles(names []string) error {
	profilesMap, err := c.GetProfilesMap()
	if err != nil {
		return err
	}
	var invalidProfileNames []string
	for _, name := range names {
		if _, ok := profilesMap[name]; !ok {
			invalidProfileNames = append(invalidProfileNames, name)
			continue
		}
		delete(profilesMap, name)
	}

	//load config
	data, err := c.configCtrl.Read()
	if err != nil {
		return err
	}

	//empty existing profile
	data.Profiles = nil
	for _, p := range profilesMap {
		// add existing profiles to the list
		data.Profiles = append(data.Profiles, p)
	}

	//save config
	err = c.configCtrl.Write(data)
	if err != nil {
		return err
	}

	// if found any invalid profiles
	if len(invalidProfileNames) > 0 {
		return fmt.Errorf("no profiles found for: %s", strings.Join(invalidProfileNames, ", "))
	}
	return nil
}

// GetProfileForExecution returns profile information for current command execution
// if profile name is provided as an argument, will return the profile,
// if profile name is not provided as argument, we will check for environment variable
// in session, then will check for profile named `default`
// bool determines whether profile is valid or not
func (c controller) GetProfileForExecution(name string) (value entity.Profile, ok bool, err error) {
	profiles, err := c.GetProfilesMap()
	if err != nil {
		return
	}
	if name != "" {
		value, ok = profiles[name]
		return
	}
	if envProfileName, exists := os.LookupEnv(odfeProfileEnvVarName); exists {
		value, ok = profiles[envProfileName]
		return
	}
	value, ok = profiles[odfeDefaultProfileName]
	return
}
