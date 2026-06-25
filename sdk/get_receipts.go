package expo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type getReceiptsRequest struct {
	IDs []string `json:"ids"`
}

// GetReceipts fetches delivery receipts for previously sent push notification tickets.
func (c *PushClient) GetReceipts(ids []string) (map[string]PushReceipt, error) {
	if len(ids) == 0 {
		return nil, errors.New("No receipt ids")
	}

	url := fmt.Sprintf("%s%s/push/getReceipts", c.host, c.apiURL)
	jsonBytes, err := json.Marshal(getReceiptsRequest{IDs: ids})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	if c.accessToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.accessToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = checkStatus(resp); err != nil {
		return nil, err
	}

	var response ReceiptsResponse
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	if response.Errors != nil {
		return nil, NewReceiptsServerError("Invalid server response", resp, &response, response.Errors)
	}
	if response.Data == nil {
		return nil, NewReceiptsServerError("Invalid server response", resp, &response, nil)
	}

	return response.Data, nil
}
