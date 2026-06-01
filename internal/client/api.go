package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/LemuriiL/GophKeeperPassword/internal/dto"
	"github.com/LemuriiL/GophKeeperPassword/internal/secure"
)

type APIClient struct {
	baseURL string
	client  *http.Client
}

func NewAPIClient(baseURL string) *APIClient {
	return &APIClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{},
	}
}

func (c *APIClient) Register(login string, password string) (string, string, error) {
	body, err := json.Marshal(dto.RegisterRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", "", err
	}

	resp, err := c.client.Post(c.baseURL+"/api/register", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		b, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("register failed: %s", string(b))
	}

	var out dto.LoginResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", "", err
	}

	salt, err := secure.NewSalt()
	if err != nil {
		return "", "", err
	}

	return out.Token, salt, nil
}

func (c *APIClient) Login(login string, password string) (string, error) {
	body, err := json.Marshal(dto.LoginRequest{
		Login:    login,
		Password: password,
	})
	if err != nil {
		return "", err
	}

	resp, err := c.client.Post(c.baseURL+"/api/login", "application/json", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("login failed: %s", string(b))
	}

	var out dto.LoginResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}

	return out.Token, nil
}

func (c *APIClient) SaveItem(token string, req dto.UpsertItemRequest) (dto.ItemResponse, error) {
	var out dto.ItemResponse

	body, err := json.Marshal(req)
	if err != nil {
		return out, err
	}

	httpReq, err := http.NewRequest(http.MethodPost, c.baseURL+"/api/items", bytes.NewReader(body))
	if err != nil {
		return out, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return out, fmt.Errorf("save failed: %s", string(b))
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func (c *APIClient) UpdateItem(token string, id string, req dto.UpsertItemRequest) (dto.ItemResponse, error) {
	var out dto.ItemResponse

	req.ID = id

	body, err := json.Marshal(req)
	if err != nil {
		return out, err
	}

	httpReq, err := http.NewRequest(http.MethodPut, c.baseURL+"/api/items/"+id, bytes.NewReader(body))
	if err != nil {
		return out, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return out, fmt.Errorf("update failed: %s", string(b))
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func (c *APIClient) GetItem(token string, id string) (dto.ItemResponse, error) {
	var out dto.ItemResponse

	httpReq, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/items/"+id, nil)
	if err != nil {
		return out, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return out, fmt.Errorf("get failed: %s", string(b))
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	return out, err
}

func (c *APIClient) ListItems(token string) ([]dto.ItemResponse, error) {
	httpReq, err := http.NewRequest(http.MethodGet, c.baseURL+"/api/items", nil)
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("list failed: %s", string(b))
	}

	var out []dto.ItemResponse
	if err = json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}

	return out, nil
}

func (c *APIClient) DeleteItem(token string, id string) error {
	httpReq, err := http.NewRequest(http.MethodDelete, c.baseURL+"/api/items/"+id, nil)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("delete failed: %s", string(b))
	}

	return nil
}
