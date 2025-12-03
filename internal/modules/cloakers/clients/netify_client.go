package clients

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/italoservio/serviosoftware_ads/pkg/env"
)

type NetifyClient struct {
	env env.Env
}

type NetifySuccessResponseData struct {
	Address     string  `json:"address"`
	Version     string  `json:"version"`
	SharedScore float64 `json:"shared_score"`
	Network     struct {
		ID    int    `json:"id"`
		Tag   string `json:"tag"`
		Label string `json:"label"`
	} `json:"network"`
	AsnRoute string `json:"asn_route"`
	Asn      struct {
		ID    int    `json:"id"`
		Tag   string `json:"tag"`
		Label string `json:"label"`
	} `json:"asn"`
	Geolocation struct {
		Continent struct {
			Code  string `json:"code"`
			Label string `json:"label"`
		} `json:"continent"`
		Country struct {
			Code  string `json:"code"`
			Label string `json:"label"`
		} `json:"country"`
	} `json:"geolocation"`
	ApplicationList []struct {
		ID    int    `json:"id"`
		Tag   string `json:"tag"`
		Label string `json:"label"`
	} `json:"application_list"`
}

type NetifySuccessResponse struct {
	StatusCode    int                       `json:"status_code"`
	StatusMessage string                    `json:"status_message"`
	Data          NetifySuccessResponseData `json:"data"`
}

type NetifyErrorResponse struct {
	StatusCode    string `json:"status_code"`
	StatusMessage string `json:"status_message"`
	StatusFields  []struct {
		Key string `json:"key"`
	} `json:"status_fields"`
}

func (e *NetifyErrorResponse) Error() string {
	msg := []string{e.StatusCode, e.StatusMessage}

	for _, field := range e.StatusFields {
		msg = append(msg, field.Key)
	}

	return strings.Join(msg, "; ")
}

func NewNetifyClient(env env.Env) *NetifyClient {
	return &NetifyClient{
		env: env,
	}
}

func (c *NetifyClient) GetIPMetadata(IP string) (*NetifySuccessResponse, error) {
	url := fmt.Sprintf("%s/api/v2/ips/%s", c.env.NETIFY_BASE_URL, IP)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.env.NETIFY_DATA_FEED_API_KEY)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var netifyError NetifyErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&netifyError); err != nil {
			return nil, fmt.Errorf("falha ao deserializar resposta: %w", err)
		}

		return nil, &netifyError
	}

	var netifyResponse NetifySuccessResponse
	if err := json.NewDecoder(resp.Body).Decode(&netifyResponse); err != nil {
		return nil, fmt.Errorf("failed to decode success response: %w", err)
	}

	return &netifyResponse, nil
}
