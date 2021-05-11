package main

import (
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/pkg/errors"
)

const (
	kvEncryptionPasswordKey = "encryptionPassword"
)

var (
	//if the plugin fails to save password to KV, this fallback password will be used
	fallbackPassword = ""
)

//EnsureEncryptionPassword generates a password for encrypting the tokens, if it does not exist
//This method is called from plugin.go, and will generate a password only the first time when the plugin is loaded
func (p *Plugin) EnsureEncryptionPassword() {
	password := GenerateEncryptionPassword()
	if _, err := p.KVEnsure(kvEncryptionPasswordKey, []byte(password)); err != nil {
		p.API.LogError("Failed to set an encryption password for the plugin, fallback password will be used.", "Error", err.Error())
		fallbackPassword = password
		return
	}
}

func (p *Plugin) getEncryptionPassword() []byte {
	//if the fallbackPassword is set this means the plugin cannot read from KV pair
	if fallbackPassword != "" {
		return []byte(fallbackPassword)
	}

	tokenSignPasswordByte, _ := p.API.KVGet(kvEncryptionPasswordKey)
	return tokenSignPasswordByte
}

//EncodeToken creates a token for WOPI
func (p *Plugin) EncodeToken(userID string, fileID string) string {
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &WopiToken{
		UserID: userID,
		FileID: fileID,
	})
	signedString, err := token.SignedString(p.getEncryptionPassword())
	if err != nil {
		p.API.LogError("Failed to encode WOPI token", "Error", err.Error())
		return ""
	}
	return signedString
}

//DecodeToken decodes a token string an returns WopiToken and isValid
func (p *Plugin) DecodeToken(tokenString string) (WopiToken, bool) {
	wopiToken := WopiToken{}
	_, err := jwt.ParseWithClaims(tokenString, &wopiToken, func(token *jwt.Token) (interface{}, error) {
		return p.getEncryptionPassword(), nil
	})

	if err != nil {
		p.API.LogError("Failed to decode WOPI token", "Error", err.Error())
		return WopiToken{}, false
	}

	return wopiToken, true
}

// GetWopiTokenFromURI decodes a token string from the URI
// returns WopiToken and error
func (p *Plugin) GetWopiTokenFromURI(uri string) (WopiToken, error) {
	token, tokenErr := getAccessTokenFromURI(uri)
	if tokenErr != nil {
		return WopiToken{}, errors.Wrap(tokenErr, "failed to retrieve token from URI: "+uri)
	}

	wopiToken, isValid := p.DecodeToken(token)
	if !isValid {
		return WopiToken{}, errors.New("collaboraOnline called the plugin with an invalid token")
	}

	return wopiToken, nil
}
