package events

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/infracost/infracost/internal/config"
	log "github.com/sirupsen/logrus"
)

func SendReport(key string, data interface{}) {
	if config.Config.PricingAPIEndpoint != config.Config.DefaultPricingAPIEndpoint {
		// skip for non-default pricing API endpoints
		return
	}

	url := fmt.Sprintf("%s/report", config.Config.PricingAPIEndpoint)

	j := make(map[string]interface{})
	j[key] = data
	j["environment"] = config.Environment

	body, err := json.Marshal(j)
	if err != nil {
		log.Debugf("Unable to generate event: %v", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Debugf("Unable to generate event: %v", err)
		return
	}

	config.AddAuthHeaders(req)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("Unable to send event: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Debugf("Unexpected response sending event: %d", resp.StatusCode)
	}
}
