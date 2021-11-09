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

package profile

import (
	"fmt"
	"opensearch-cli/controller/config"
	"opensearch-cli/entity"
	"opensearch-cli/environment"
	"os"
	"strings"
)

const (
	DefaultProfileName = "default"
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
		if value, ok = profiles[name]; ok {
			return
		}
		return value, ok, fmt.Errorf("profile '%s' does not exist", name)
	}
	if envProfileName, exists := os.LookupEnv(environment.OPENSEARCH_PROFILE); exists {
		if value, ok = profiles[envProfileName]; ok {
			return
		}
		return value, ok, fmt.Errorf("profile '%s' does not exist", envProfileName)
	}
	value, ok = profiles[DefaultProfileName]
	return
}
