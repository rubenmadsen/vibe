package nodes

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"costner/pkg/types"
)

type RequestNode struct {
	types.BaseNode
}

func NewRequestNode(id string) *RequestNode {
	node := &RequestNode{
		BaseNode: types.BaseNode{
			NodeID:   id,
			NodeType: "request",
			NodeName: "HTTP Request",
			Inputs: []types.NodeInput{
				{Name: "url", Type: "string", Required: true, Description: "Request URL"},
				{Name: "method", Type: "string", Required: false, Description: "HTTP method", Value: "GET"},
				{Name: "headers", Type: "map", Required: false, Description: "Request headers"},
				{Name: "body", Type: "string", Required: false, Description: "Request body"},
				{Name: "timeout", Type: "int", Required: false, Description: "Timeout in seconds", Value: 30},
			},
			Outputs: []types.NodeOutput{
				{Name: "status_code", Type: "int", Description: "HTTP status code"},
				{Name: "headers", Type: "map", Description: "Response headers"},
				{Name: "body", Type: "string", Description: "Response body"},
				{Name: "duration", Type: "duration", Description: "Request duration"},
			},
			Config: make(map[string]interface{}),
		},
	}
	return node
}

func (n *RequestNode) Execute(ctx context.Context, inputs map[string]interface{}) (map[string]interface{}, error) {
	url, ok := inputs["url"].(string)
	if !ok || url == "" {
		return nil, fmt.Errorf("url is required")
	}

	method := "GET"
	if m, ok := inputs["method"].(string); ok && m != "" {
		method = strings.ToUpper(m)
	}

	timeout := 30
	if t, ok := inputs["timeout"].(int); ok && t > 0 {
		timeout = t
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	// Prepare request body
	var bodyReader io.Reader
	if body, ok := inputs["body"].(string); ok && body != "" {
		bodyReader = strings.NewReader(body)
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	if headers, ok := inputs["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}

	// Execute request
	start := time.Now()
	resp, err := client.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Convert response headers to map
	responseHeaders := make(map[string]interface{})
	for key, values := range resp.Header {
		if len(values) == 1 {
			responseHeaders[key] = values[0]
		} else {
			responseHeaders[key] = values
		}
	}

	result := map[string]interface{}{
		"status_code": resp.StatusCode,
		"headers":     responseHeaders,
		"body":        string(bodyBytes),
		"duration":    duration,
	}

	// Update output values
	for i := range n.Outputs {
		if value, exists := result[n.Outputs[i].Name]; exists {
			n.Outputs[i].Value = value
		}
	}

	return result, nil
}

func (n *RequestNode) Clone() types.Node {
	clone := NewRequestNode(n.NodeID)
	clone.NodeName = n.NodeName
	clone.Position = n.Position
	clone.Config = make(map[string]interface{})
	for k, v := range n.Config {
		clone.Config[k] = v
	}
	return clone
}

func (n *RequestNode) Serialize() ([]byte, error) {
	return json.Marshal(n.BaseNode)
}

func (n *RequestNode) Deserialize(data []byte) error {
	return json.Unmarshal(data, &n.BaseNode)
}