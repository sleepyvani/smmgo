package smm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Client struct {
	apiURL string
	apiKey string
}

func New(apiURL, apiKey string) *Client {
	return &Client{
		apiURL: apiURL,
		apiKey: apiKey,
	}
}

func (c *Client) request(action string, data map[string]interface{}) (map[string]interface{}, error) {
	formData := url.Values{}
	formData.Set("key", c.apiKey)
	formData.Set("action", action)

	for key, value := range data {
		if value != nil {
			formData.Set(key, fmt.Sprintf("%v", value))
		}
	}

	req, err := http.NewRequest("POST", c.apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 5.01; Windows NT 5.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP Error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		// Some panels return HTML error pages or plain text
		return nil, fmt.Errorf("Invalid JSON response: %s", string(body))
	}

	if errStr, ok := result["error"].(string); ok && errStr != "" {
		return nil, fmt.Errorf("%s", errStr)
	}

	return result, nil
}

func (c *Client) requestList(action string, data map[string]interface{}) ([]interface{}, error) {
	// Similar to request but expects a list response
	// This is a simplified version, ideally refactor request to handle both or return any
	formData := url.Values{}
	formData.Set("key", c.apiKey)
	formData.Set("action", action)

	for key, value := range data {
		if value != nil {
			formData.Set(key, fmt.Sprintf("%v", value))
		}
	}

	req, err := http.NewRequest("POST", c.apiURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Mozilla/4.0 (compatible; MSIE 5.01; Windows NT 5.0)")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP Error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("Invalid JSON response: %s", string(body))
	}

	// Check for error object
	if resMap, ok := result.(map[string]interface{}); ok {
		if errStr, ok := resMap["error"].(string); ok && errStr != "" {
			return nil, fmt.Errorf("%s", errStr)
		}
	}

	if list, ok := result.([]interface{}); ok {
		return list, nil
	}
	
	// Fallback if it returns object but we expected list (should likely error out or wrap)
	return nil, fmt.Errorf("Expected list response, got %T", result)
}


// GetServices returns the list of services
func (c *Client) GetServices() ([]interface{}, error) {
	return c.requestList("services", nil)
}

// GetStatus returns the status of an order
func (c *Client) GetStatus(orderID interface{}) (map[string]interface{}, error) {
	return c.request("status", map[string]interface{}{"order": orderID})
}

// GetMultiStatus returns status for multiple orders
func (c *Client) GetMultiStatus(orderIDs []interface{}) (map[string]interface{}, error) {
	ids := make([]string, len(orderIDs))
	for i, v := range orderIDs {
		ids[i] = fmt.Sprintf("%v", v)
	}
	return c.request("status", map[string]interface{}{"orders": strings.Join(ids, ",")})
}

// CreateRefill creates a refill for an order
func (c *Client) CreateRefill(orderID interface{}) (map[string]interface{}, error) {
	return c.request("refill", map[string]interface{}{"order": orderID})
}

// CreateMultiRefill creates refills for multiple orders
func (c *Client) CreateMultiRefill(orderIDs []interface{}) ([]interface{}, error) {
	ids := make([]string, len(orderIDs))
	for i, v := range orderIDs {
		ids[i] = fmt.Sprintf("%v", v)
	}
	return c.requestList("refill", map[string]interface{}{"orders": strings.Join(ids, ",")})
}

// GetRefillStatus returns the status of a refill
func (c *Client) GetRefillStatus(refillID interface{}) (map[string]interface{}, error) {
	return c.request("refill_status", map[string]interface{}{"refill": refillID})
}

// GetMultiRefillStatus returns status for multiple refills
func (c *Client) GetMultiRefillStatus(refillIDs []interface{}) ([]interface{}, error) {
	ids := make([]string, len(refillIDs))
	for i, v := range refillIDs {
		ids[i] = fmt.Sprintf("%v", v)
	}
	return c.requestList("refill_status", map[string]interface{}{"refills": strings.Join(ids, ",")})
}

// CancelOrders cancels multiple orders
func (c *Client) CancelOrders(orderIDs []interface{}) ([]interface{}, error) {
	ids := make([]string, len(orderIDs))
	for i, v := range orderIDs {
		ids[i] = fmt.Sprintf("%v", v)
	}
	return c.requestList("cancel", map[string]interface{}{"orders": strings.Join(ids, ",")})
}

// GetBalance returns user balance
func (c *Client) GetBalance() (map[string]interface{}, error) {
	return c.request("balance", nil)
}

// AddOrderParams parameters for adding an order
type AddOrderParams struct {
	Service   interface{}
	Link      string
	Quantity  int
	Comments  string
	Runs      int
	Interval  int
}

// AddOrder adds a new order
func (c *Client) AddOrder(params AddOrderParams) (map[string]interface{}, error) {
	data := map[string]interface{}{
		"service":  params.Service,
		"link":     params.Link,
		"quantity": params.Quantity,
	}
	if params.Comments != "" {
		data["comments"] = params.Comments
	}
	if params.Runs > 0 {
		data["runs"] = params.Runs
	}
	if params.Interval > 0 {
		data["interval"] = params.Interval
	}
	return c.request("add", data)
}
