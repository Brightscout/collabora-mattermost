package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/model"
	"github.com/mattermost/mattermost-server/plugin"
)

// MessageWillBePosted will set a post type for each post that contains at least one file
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	//change the post type only if it contains any files that can be viewed/edited with Collabora Online
	changePostType := false
	for i := 0; i < len(post.FileIds); i++ {
		fileInfo, fileInfoError := p.API.GetFileInfo(post.FileIds[i])
		if fileInfoError != nil {
			p.API.LogError("Could not retrieve file info on message post")
			continue
		}
		_, ok := WOPIFiles[strings.ToLower(fileInfo.Extension)]
		if ok {
			changePostType = true
		}
	}

	if changePostType {
		post.Type = "custom_post_with_file"
	}

	return post, ""
}
