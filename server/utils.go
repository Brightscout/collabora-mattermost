package main

import (
	"crypto/tls"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

func (p *Plugin) getHTTPClient() *http.Client {
	config := p.getConfiguration()
	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	if config.DisableCertificateVerification {
		customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	client := &http.Client{Transport: customTransport}
	return client
}

// getAccessTokenFromURI extracts the access_token from the URI
// We need to do this manually as Mattermost removes the access_token before it reaches the plugin HTTP request parser
func getAccessTokenFromURI(uri string) (string, error) {
	parsedURL, err := url.Parse(uri)
	if err != nil {
		return "", errors.Wrap(err, "failed to parse uri")
	}
	urlValues, parseErr := url.ParseQuery(parsedURL.RawQuery)
	if parseErr != nil {
		return "", errors.Wrap(parseErr, "failed to parse raw query")
	}
	return urlValues.Get("access_token"), nil
}
