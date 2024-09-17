package plugin

import (
	"bytes"
	"encoding/json"
	"net/http"
)

func (p *Plugin) RunCurlCommand() {
	configuration := p.API.GetPluginConfig()
	var url string
	for key, value := range configuration {
		webhookURL, ok := value.(string)
		if key == "webhook" && webhookURL != "" {
			if ok {
				url = webhookURL
			}
		}
	}
	// url := "https://mattermost.vivasoftltd.com/hooks/uweuwoha7iy9intjxq6d5njt8r"
	// url := "http://localhost:8065/hooks/agpyd15ptjgk58scamzh7e9epw"
	message := map[string]interface{}{
		"text": "**🌅 Good Morning @all!** \n\n**How do you feel today?** \nHope your **Work From Home** status has been updated properly.\n*Have a great day ahead!*\n\n**Thanks!**",
		// "text": "## 🌅 **Good Morning** @all!\n\n---\n\n### **How do you feel today?** 😊\n\nWe hope your *work from home* status has been updated properly.\n\n---\n\n**Thanks!** 🙏\n\nHave a great day ahead!",
		// "text": "\n\n**Good Morning** @all! ☀️\n\n---\n\n**How do you feel today?** 😊  \nHope your *work from home* status has been updated properly.\n\n---\n\n**Thanks!** 🙏",
	}
	jsonData, err := json.Marshal(message)
	if err != nil {
		p.API.LogError("Failed to marshal JSON data", "error", err)
		return
	}
	var req *http.Request
	if url != "" {
		req, err = http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
		if err != nil {
			p.API.LogError("Failed to create new request", "error", err)
			return
		}
	}

	req.Header.Set("Content-Type", "application/json")
	p.API.LogError("test command", "error", err)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		p.API.LogError("Failed to execute POST request", "error", err)
		return
	}
	defer resp.Body.Close()

	p.API.LogInfo("POST request sent successfully", "status", resp.Status)
}
