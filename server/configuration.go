package main

import (
	"encoding/xml"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/pkg/errors"
)

type configuration struct {
	WOPIAddress string
}

var (
	//WOPIData contains the XML from <WOPI>/hosting/discovery
	WOPIData WopiDiscovery

	//WOPIFiles maps file extension with file action & url
	WOPIFiles map[string]WOPIFileInfo
)

// Clone deep copies the configuration
func (c *configuration) Clone() *configuration {
	return &configuration{WOPIAddress: c.WOPIAddress}
}

// ProcessConfiguration processes the config.
func (c *configuration) ProcessConfiguration() error {
	// trim trailing slash or spaces from the WOPI address, if needed
	c.WOPIAddress = strings.TrimSpace(c.WOPIAddress)
	c.WOPIAddress = strings.Trim(c.WOPIAddress, "/")

	return nil
}

// IsValid checks if all needed fields are set.
func (c *configuration) IsValid() error {
	if c.WOPIAddress == "" {
		return errors.New("please provide the WOPIAddress")
	}

	return nil
}

// OnConfigurationChange is called when plugin's configuration changes
func (p *Plugin) OnConfigurationChange() error {
	var configuration = new(configuration)

	// Load the public configuration fields from the Mattermost server configuration.
	if loadConfigErr := p.API.LoadPluginConfiguration(configuration); loadConfigErr != nil {
		return errors.Wrap(loadConfigErr, "failed to load plugin configuration")
	}

	if err := configuration.ProcessConfiguration(); err != nil {
		p.API.LogError("Error in ProcessConfiguration.", "Error", err.Error())
		return err
	}

	if err := configuration.IsValid(); err != nil {
		p.API.LogError("Error in Validating Configuration.", "Error", err.Error())
		return err
	}

	p.setConfiguration(configuration)
	p.loadWopiFileInfo(configuration.WOPIAddress)

	return nil
}

// setConfiguration sets the new configuration
func (p *Plugin) setConfiguration(configuration *configuration) {
	p.configurationLock.Lock()
	defer p.configurationLock.Unlock()

	if configuration != nil && p.configuration == configuration {
		// Ignore assignment if the configuration struct is empty. Go will optimize the
		// allocation for same to point at the same memory address, breaking the check
		// above.
		if reflect.ValueOf(*configuration).NumField() == 0 {
			return
		}

		panic("setConfiguration called with the existing configuration")
	}

	p.configuration = configuration
}

//loadWopiFileInfo loads the WOPI file data
func (p *Plugin) loadWopiFileInfo(wopiAddress string) {
	client := getHTTPClient()
	resp, err := client.Get(wopiAddress + "/hosting/discovery")
	if err != nil {
		p.API.LogError("WOPI request error. Please check the WOPI address.", err.Error())
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.API.LogError("WOPI request error. Failed to read WOPI request body. Please check the WOPI address.", err.Error())
		return
	}

	if err := xml.Unmarshal(body, &WOPIData); err != nil {
		p.API.LogError("WOPI request error. Failed to unmarshal WOPI XML. Please check the WOPI address.", err.Error())
		return
	}

	WOPIFiles = make(map[string]WOPIFileInfo)
	for i := 0; i < len(WOPIData.NetZone.App); i++ {
		for j := 0; j < len(WOPIData.NetZone.App[i].Action); j++ {
			ext := strings.ToLower(WOPIData.NetZone.App[i].Action[j].Ext)
			if ext == "" || ext == "png" || ext == "jpg" || ext == "jpeg" || ext == "gif" {
				continue
			}
			WOPIFiles[strings.ToLower(ext)] = WOPIFileInfo{WOPIData.NetZone.App[i].Action[j].URLSrc, WOPIData.NetZone.App[i].Action[j].Name}
		}
	}
	p.API.LogInfo("WOPI file info loaded successfully!")
}
