package manifest

import (
	"encoding/json"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
)

var Manifest *model.Manifest

const manifestStr = `
{
  "id": "wfh-plugin",
  "name": "WFH-Plugin",
  "description": "This plugin serves as a starting point for writing a Mattermost plugin.",
  "homepage_url": "https://github.com/mattermost/mattermost-plugin-starter-template",
  "support_url": "https://github.com/mattermost/mattermost-plugin-starter-template/issues",
  "icon_path": "assets/starter-template-icon.svg",
  "version": "0.0.0+8867637",
  "min_server_version": "6.2.1",
  "server": {
    "executables": {
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "darwin-arm64": "server/dist/plugin-darwin-arm64",
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "linux-arm64": "server/dist/plugin-linux-arm64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  },
  "webapp": {
    "bundle_path": "webapp/dist/main.js"
  },
  "settings_schema": {
    "header": "Header: Configure your plugin settings below.",
    "footer": "",
    "settings": [
      {
        "key": "ChannelName",
        "display_name": "Channel Name:",
        "type": "text",
        "help_text": "The channel to use as part of the plugin, created for each team automatically if it does not exist.",
        "placeholder": "channel name (small case)",
        "default": "",
        "hosting": ""
      },
      {
        "key": "MessageKey",
        "display_name": "Input Message:",
        "type": "text",
        "help_text": "Select the input message  for the system user.",
        "placeholder": "input message key",
        "default": null,
        "hosting": ""
      }
    ]
  },
  "props": {
    "support_packet": "Plugin support packet"
  }
}
`

func init() {
	_ = json.NewDecoder(strings.NewReader(manifestStr)).Decode(&Manifest)
}
