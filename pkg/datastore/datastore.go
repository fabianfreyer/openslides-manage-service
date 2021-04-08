package datastore

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Get gets a fqField from the datastore.
func Get(ctx context.Context, addr string, key string, value interface{}) error {
	keyParts := strings.Split(key, "/")
	if len(keyParts) != 3 {
		return fmt.Errorf("invalid key %s, expected two `/`", key)
	}

	reqBody := fmt.Sprintf(
		`{
			"fqid": "%s/%s",
			"mapped_fields": ["%s"]
		}`,
		keyParts[0], keyParts[1], keyParts[2],
	)

	req, err := http.NewRequestWithContext(ctx, "POST", addr+"/internal/datastore/reader/get", strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("creating request to datastore: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request to datastore: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("[can not read body]")
		}
		return fmt.Errorf("got response `%s`: %s", resp.Status, body)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	var respData map[string]json.RawMessage
	if err := json.Unmarshal(respBody, &respData); err != nil {
		return fmt.Errorf("decoding response body `%s`: %w", respBody, err)
	}

	if err := json.Unmarshal(respData[keyParts[2]], value); err != nil {
		return fmt.Errorf("decoding response field: %w", err)
	}

	return nil
}

// Set sets a fqField at the datastore. Value has to be json.
func Set(ctx context.Context, addr, key string, value json.RawMessage) error {
	parts := strings.Split(key, "/")
	if len(parts) != 3 {
		return fmt.Errorf("invalid key %s, expected two `/`", key)
	}

	reqBody := fmt.Sprintf(
		`{
			"user_id": 0,
			"information": {},
			"locked_fields":{}, "events":[
				{"type":"update","fqid":"%s/%s","fields":{"%s":%s}}
			]
		}`,
		parts[0], parts[1], parts[2], value,
	)

	req, err := http.NewRequestWithContext(ctx, "POST", addr+"/internal/datastore/writer/write", strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("creating request to %s: %w", addr, err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending request to %s: %w", addr, err)
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte("[can not read body]")
		}
		return fmt.Errorf("got response `%s`: %s", resp.Status, body)
	}

	return nil
}