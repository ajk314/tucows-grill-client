package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	
	"tucows-grill-client/internal/models"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
	jwtToken   string
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{},
		baseURL:    "http://localhost:8080", // should use an env var to update this value, and set the env var via aws config or k8 configmap, or a secret
	}
}

func (c *Client) Login(username, password string) (string, error) {
	credentials := map[string]string{
		"username": username,
		"password": password,
	}
	body, err := json.Marshal(credentials)
	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Post(fmt.Sprintf("%s/login", c.baseURL), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("login failed, status: %s", resp.Status)
	}

	var result map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	var ok bool
	c.jwtToken, ok = result["token"]
	if !ok {
		return "", fmt.Errorf("no token found in response")
	}

	return c.jwtToken, nil
}

func (c *Client) GetIngredientByID(id int) (*models.Ingredient, error) {
	url := fmt.Sprintf("%s/ingredients/%d", c.baseURL, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.jwtToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch ingredient, status: %s", resp.Status)
	}

	var ingredient models.Ingredient
	if err := json.NewDecoder(resp.Body).Decode(&ingredient); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &ingredient, nil
}

func (c *Client) PostIngredient(ingredient models.Ingredient) error {
	url := fmt.Sprintf("%s/ingredients", c.baseURL)

	jsonData, err := json.Marshal(ingredient)
	if err != nil {
		return fmt.Errorf("failed to marshal ingredient data: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.jwtToken))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to post ingredient, status: %s", resp.Status)
	}

	return nil
}

// async flag will change the url to hit the async endpoint, which uses go routines
func (c *Client) GetTotalCostForItem(itemID int, asyncBool string) (float64, error) {
	path := "total-cost"
	if asyncBool == "true" {
		path = "total-cost-async"
	}
    req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s?item_id=%d", c.baseURL, path, itemID), nil)
    if err != nil {
        return 0, err
    }
    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.jwtToken))

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return 0, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return 0, fmt.Errorf("failed to get total cost, status: %s", resp.Status)
    }

    var result map[string]float64
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return 0, err
    }

    totalCost, ok := result["total_cost"]
    if !ok {
        return 0, fmt.Errorf("total_cost not found in response")
    }

    return totalCost, nil
}