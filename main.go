package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/a-h/templ"
	_ "golang.org/x/crypto/x509roots/fallback"
)

// OAuthResponse holds the temporary access token from the auth endpoint.
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
}

// Environment represents the necessary fields from an Upsun environment.
type Environment struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Status          string          `json:"status"`
	Title           string          `json:"title"`
	IsPR            bool            `json:"is_pr"`
	URL             string          `json:"edge_hostname"`
	DeploymentState DeploymentState `json:"deployment_state"`
	Links           Links           `json:"_links"`
}

type DeploymentState struct {
	LastStateUpdateSuccessful bool `json:"last_state_update_successful"`
}

type Links struct {
	PfRoutes []Link `json:"pf:routes,omitempty"`
}

type Link struct {
	Href string `json:"href"`
}

//go:generate go tool templ generate
func main() {
	port := ":80"
	if value, ok := os.LookupEnv("PORT"); ok {
		port = value
	}

	ctx := context.Background()

	environments, err := data(ctx)
	if err != nil {
		log.Fatalf("could not get data: %v", err)
	}

	// Enable streaming in the handler
	handler := templ.Handler(Page(environments), templ.WithStreaming())

	log.Printf("Server running on %s\n", port)

	//nolint:mnd
	server := &http.Server{
		Addr:              port,
		ReadHeaderTimeout: 3 * time.Second,
	}

	http.Handle("/", handler)

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("could not start webserver: %s\n", err)
	}
}

//nolint:cyclop,funlen
func data(ctx context.Context) ([]Environment, error) {
	apiToken := os.Getenv("UPSUN_API_TOKEN")
	projectID := os.Getenv("UPSUN_PROJECT_ID")

	if apiToken == "" || projectID == "" {
		log.Fatalln("please set UPSUN_API_TOKEN and UPSUN_PROJECT_ID environment variables.")
	}

	// 1. Exchange the API Token for an OAuth Access Token
	authData := url.Values{}
	authData.Set("grant_type", "api_token")
	authData.Set("api_token", apiToken)

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://auth.upsun.com/oauth2/token",
		strings.NewReader(authData.Encode()))
	if err != nil {
		log.Fatalf("could not create request: %v\n", err)
	}

	// Upsun requires basic auth using 'platform-api-user' with no password for token exchange
	req.SetBasicAuth("platform-api-user", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error requesting access token: %w", err)
	}

	if resp == nil {
		//nolint:err113
		return nil, errors.New("error requesting access token: response is nil")
	}

	defer resp.Body.Close()

	//nolint:err113
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to authenticate, HTTP status: %d", resp.StatusCode)
	}

	var authResp OAuthResponse

	err = json.NewDecoder(resp.Body).Decode(&authResp)
	if err != nil {
		return nil, fmt.Errorf("error parsing auth response: %w", err)
	}

	// 2. Fetch the list of environments
	envURL := fmt.Sprintf("https://api.upsun.com/projects/%s/environments", projectID)

	//nolint:gosec
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, envURL, nil)
	if err != nil {
		return nil, fmt.Errorf("could not create request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+authResp.AccessToken)

	//nolint:gosec
	resp, err = client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error fetching environments: %w", err)
	}

	if resp == nil {
		//nolint:err113
		return nil, errors.New("error fetching environments: response is nil")
	}

	defer resp.Body.Close()

	//nolint:err113
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get environments, HTTP status: %d", resp.StatusCode)
	}

	// 3. Parse the JSON array and filter for 'active'
	var environments []Environment

	err = json.NewDecoder(resp.Body).Decode(&environments)
	if err != nil {
		return nil, fmt.Errorf("error parsing environments response: %w", err)
	}

	var filteredEnvironments []Environment

	for _, env := range environments {
		if env.IsPR {
			filteredEnvironments = append(filteredEnvironments, env)
		}
	}

	return filteredEnvironments, nil
}
