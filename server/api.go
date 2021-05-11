package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
)

const (
	HeaderMattermostUserID = "Mattermost-User-Id"
)

// InitAPI initializes the REST API
func (p *Plugin) InitAPI() *mux.Router {
	r := mux.NewRouter()
	r.Use(p.withRecovery)

	p.handleStaticFiles(r)
	s := r.PathPrefix("/api/v1").Subrouter()

	// Add the custom plugin routes here
	s.HandleFunc("/fileInfo", handleAuthRequired(p.parseFileIDs)).Methods(http.MethodGet)
	s.HandleFunc("/wopiFileList", handleAuthRequired(p.returnWopiFileList)).Methods(http.MethodGet)
	s.HandleFunc("/collaboraURL", handleAuthRequired(p.returnCollaboraOnlineFileURL)).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}", p.returnWopiFileInfo).Methods(http.MethodGet)
	s.HandleFunc("/wopi/files/{fileID:[a-z0-9]+}/contents", p.parseWopiRequests).Methods(http.MethodGet, http.MethodPost)

	// 404 handler
	r.Handle("{anything:.*}", http.NotFoundHandler())
	return r
}

func (p *Plugin) getBaseAPIURL() string {
	return *p.API.GetConfig().ServiceSettings.SiteURL + "/plugins/" + manifest.Id + "/api/v1"
}

// withRecovery allows recovery from panics
func (p *Plugin) withRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if x := recover(); x != nil {
				p.API.LogError("Recovered from a panic",
					"url", r.URL.String(),
					"error", x,
					"stack", string(debug.Stack()))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// handleStaticFiles handles the static files under the assets directory.
func (p *Plugin) handleStaticFiles(r *mux.Router) {
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		p.API.LogWarn("Failed to get bundle path.", "Error", err.Error())
		return
	}

	// This will serve static files from the 'assets' directory under '/static/<filename>'
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(filepath.Join(bundlePath, "assets")))))
}

// handleAuthRequired verifies if provided request is performed by a logged-in Mattermost user.
func handleAuthRequired(handleFunc func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Header.Get(HeaderMattermostUserID)
		if userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		handleFunc(w, r)
	}
}

// parseFileIDs sends the file info to the client (name, extension and id) for each file
// body contains an array with file ids in JSON format
func (p *Plugin) parseFileIDs(w http.ResponseWriter, r *http.Request) {
	//extract fileIDs array from body
	body, bodyReadError := ioutil.ReadAll(r.Body)
	if bodyReadError != nil {
		p.API.LogError("Error when reading body: ", bodyReadError.Error())
		return
	}
	var fileIDs []string
	_ = json.Unmarshal(body, &fileIDs)

	//create an array with more detailed file info for each file
	files := make([]ClientFileInfo, 0, len(fileIDs))
	for _, fileID := range fileIDs {
		fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
		if fileInfoError != nil {
			p.API.LogError("Error when retrieving file info: ", fileInfoError.Error())
			continue
		}
		if value, ok := WopiFiles[strings.ToLower(fileInfo.Extension)]; ok {
			file := ClientFileInfo{
				fileInfo.Id,
				fileInfo.Name,
				fileInfo.Extension,
				value.Action,
			}
			files = append(files, file)
		}
	}

	responseJSON, _ := json.Marshal(files)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

// returnWopiFileList returns the list with file extensions and actions associated with these files
func (p *Plugin) returnWopiFileList(w http.ResponseWriter, r *http.Request) {
	responseJSON, _ := json.Marshal(WopiFiles)
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

// returnCollaboraOnlineFileURL returns the URL and token that the client will use to
// load Collabora Online in the iframe
func (p *Plugin) returnCollaboraOnlineFileURL(w http.ResponseWriter, r *http.Request) {
	//retrieve fileID and file info
	fileID := r.URL.Query().Get("file_id")
	if fileID == "" {
		p.API.LogError("file_id query parameter missing!")
		http.Error(w, "missing file_id parameter", http.StatusBadRequest)
		return
	}

	file, fileError := p.API.GetFileInfo(fileID)
	if fileError != nil {
		p.API.LogError("Failed to retrieve file. Error: ", fileError.Error())
		http.Error(w, "Invalid fileID. Error: " + fileError.Error(), http.StatusBadRequest)
		return
	}

	wopiURL := WopiFiles[strings.ToLower(file.Extension)].URL + "WOPISrc=" + p.getBaseAPIURL() + "/wopi/files/" + fileID
	wopiToken := p.EncodeToken(r.Header.Get(HeaderMattermostUserID), fileID)

	response := struct {
		URL         string `json:"url"`
		AccessToken string `json:"access_token"` //client will pass this token as a POST parameter to Collabora Online when loading the iframe
	}{wopiURL, wopiToken}

	responseJSON, _ := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}

// parseWopiRequests is used by Collabora Online server to get/save the contents of a file
func (p *Plugin) parseWopiRequests(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	token, tokenErr := getAccessTokenFromURI(r.RequestURI)
	if tokenErr != nil {
		p.API.LogError("Error retrieving token from URI: "+r.RequestURI, "Error", tokenErr.Error())
		return
	}

	wopiToken, isValid := p.DecodeToken(token)
	if !isValid || wopiToken.FileID != fileID {
		p.API.LogError("Invalid token.")
		return
	}

	fileContent, err := p.API.GetFile(fileID)
	if err != nil {
		p.API.LogError("Error retrieving file info, fileID: " + fileID)
		return
	}

	fileInfo, fileInfoError := p.API.GetFileInfo(fileID)
	if fileInfoError != nil {
		p.API.LogError("Error occurred when retrieving file info: " + fileInfoError.Error())
		return
	}

	postInfo, postInfoError := p.API.GetPost(fileInfo.PostId)
	if postInfoError != nil {
		p.API.LogError("Error occurred when retrieving post info for file: " + postInfoError.Error())
		return
	}

	//check if user has access to the channel where the file was sent
	//p.API.HasPermissionToChannel(userID,channelID) was returning false for some reason...
	members, channelMembersError := p.API.GetChannelMembersByIds(postInfo.ChannelId, []string{wopiToken.UserID})
	if channelMembersError != nil {
		p.API.LogError("Error occurred when retrieving channel members: " + channelMembersError.Error())
	}
	if members == nil {
		p.API.LogError("User doesn't have access to the channel where the file was sent")
		return
	}

	//send file to Collabora Online
	if r.Method == http.MethodGet {
		if _, err := w.Write(fileContent); err != nil {
			p.API.LogError("failed to write status", "err", err.Error())
		}
	}

	//save file received from Collabora Online
	if r.Method == http.MethodPost {
		f, fileCreateError := os.Create("./data/" + fileInfo.Path)
		if fileCreateError != nil {
			p.API.LogError("Error occurred when creating new file: ", fileCreateError.Error())
			return
		}

		body, bodyReadError := ioutil.ReadAll(r.Body)
		if bodyReadError != nil {
			p.API.LogError("Error occurred when reading body:", bodyReadError.Error())
			return
		}

		_, fileSaveError := f.Write(body)
		if fileSaveError != nil {
			p.API.LogError("Error occurred when writing contents to file: " + fileSaveError.Error())
			f.Close()
			return
		}

		fileCloseError := f.Close()
		if fileCloseError != nil {
			p.API.LogError("Error occurred when closing the file: " + fileCloseError.Error())
			return
		}
	}
}

// returnWopiFileInfo returns the file information, used by Collabora Online
// see: http://wopi.readthedocs.io/projects/wopirest/en/latest/files/CheckFileInfo.html#checkfileinfo
func (p *Plugin) returnWopiFileInfo(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fileID := params["fileID"]

	token, tokenErr := getAccessTokenFromURI(r.RequestURI)
	if tokenErr != nil {
		p.API.LogError("Error retrieving token from URI:"+r.RequestURI, "Error", tokenErr.Error())
		return
	}

	wopiToken, isValid := p.DecodeToken(token)
	if !isValid || wopiToken.FileID != fileID {
		p.API.LogError("Collabora Online called the plugin with an invalid token.")
		return
	}

	user, userErr := p.API.GetUser(wopiToken.UserID)
	if userErr != nil {
		p.API.LogError("Error retrieving user. Token UserID is corrupted or the user doesn't exist.")
		return
	}

	fileInfo, err := p.API.GetFileInfo(fileID)
	if err != nil {
		p.API.LogError("Error retrieving file info, fileID: " + fileID)
		return
	}

	post, postErr := p.API.GetPost(fileInfo.PostId)
	if postErr != nil {
		p.API.LogError("Error retrieving file's post, postId: " + fileInfo.PostId)
		return
	}

	wopiFileInfo := WopiCheckFileInfo{
		BaseFileName:            fileInfo.Name,
		Size:                    fileInfo.Size,
		OwnerID:                 post.UserId,
		UserID:                  user.Id,
		UserFriendlyName:        user.GetFullName(),
		UserCanWrite:            true,
		UserCanNotWriteRelative: true,
	}

	responseJSON, _ := json.Marshal(wopiFileInfo)

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(responseJSON); err != nil {
		p.API.LogError("failed to write status", "err", err.Error())
	}
}
