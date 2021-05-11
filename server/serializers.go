package main

import "encoding/xml"

//WopiDiscovery represents the XML from <WOPI>/hosting/discovery
type WopiDiscovery struct {
	XMLName xml.Name `xml:"wopi-discovery"`
	Text    string   `xml:",chardata"`
	NetZone struct {
		Text string `xml:",chardata"`
		Name string `xml:"name,attr"`
		App  []struct {
			Text   string `xml:",chardata"`
			Name   string `xml:"name,attr"`
			Action []struct {
				Text   string `xml:",chardata"`
				Ext    string `xml:"ext,attr"`
				Name   string `xml:"name,attr"`
				URLSrc string `xml:"urlsrc,attr"`
			} `xml:"action"`
		} `xml:"app"`
	} `xml:"net-zone"`
}

// WOPICheckFileInfo is the required response from http://wopi.readthedocs.io/projects/wopirest/en/latest/files/CheckFileInfo.html#checkfileinfo
type WOPICheckFileInfo struct {
	// The string name of the file, including extension, without a path. Used for display in user interface (UI), and determining the extension of the file.
	BaseFileName string `json:"BaseFileName"`

	// The size of the file in bytes, expressed as a long, a 64-bit signed integer.
	Size int64 `json:"Size"`

	// A string that uniquely identifies the owner of the file.
	OwnerID string `json:"OwnerId"`

	// A string value uniquely identifying the user currently accessing the file.
	UserID string `json:"UserId"`

	// The name visible to other users while editing collaboratively.
	UserFriendlyName string `json:"UserFriendlyName"`

	// User permissions
	UserCanWrite bool `json:"UserCanWrite"`

	// Enables/disables the "Save As" acton in the File menu
	UserCanNotWriteRelative bool `json:"UserCanNotWriteRelative"`
}

//WOPIFileInfo is used top map file extension with the action & url
type WOPIFileInfo struct {
	URL    string //WOPI url to view/edit the file
	Action string //edit or view
}

//CollaboraFileInfo contains file information sent to the client
type CollaboraFileInfo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Extension string `json:"extension"`
	Action    string `json:"action"` //view or edit
}
