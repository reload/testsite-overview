package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/a-h/templ"
	_ "golang.org/x/crypto/x509roots/fallback"
)

// OAuthResponse holds the temporary access token from the auth endpoint
type OAuthResponse struct {
	AccessToken string `json:"access_token"`
}

// Environment represents the necessary fields from an Upsun environment
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

	environments := data()
	// Enable streaming in the handler
	handler := templ.Handler(Page(environments), templ.WithStreaming())

	http.Handle("/", handler)

	log.Printf("Server running on %s\n", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatalf("could not start webserver: %s\n", err)
	}
}

func data() []Environment {
	apiToken := os.Getenv("UPSUN_API_TOKEN")
	projectID := os.Getenv("UPSUN_PROJECT_ID")

	if apiToken == "" || projectID == "" {
		log.Println("Error: Please set UPSUN_API_TOKEN and UPSUN_PROJECT_ID environment variables.")
		os.Exit(1)
	}

	// 1. Exchange the API Token for an OAuth Access Token
	authData := url.Values{}
	authData.Set("grant_type", "api_token")
	authData.Set("api_token", apiToken)

	req, _ := http.NewRequest("POST", "https://auth.upsun.com/oauth2/token", strings.NewReader(authData.Encode()))
	// Upsun requires basic auth using 'platform-api-user' with no password for token exchange
	req.SetBasicAuth("platform-api-user", "")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error requesting access token: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to authenticate. HTTP Status: %d\n", resp.StatusCode)
		os.Exit(1)
	}

	var authResp OAuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		log.Printf("Error parsing auth response: %v\n", err)
		os.Exit(1)
	}

	// 2. Fetch the list of environments
	envURL := fmt.Sprintf("https://api.upsun.com/projects/%s/environments", projectID)
	req, _ = http.NewRequest("GET", envURL, nil)
	req.Header.Add("Authorization", "Bearer "+authResp.AccessToken)

	resp, err = client.Do(req)
	if err != nil {
		log.Printf("Error fetching environments: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			log.Printf("failed to close response body: %v", err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		log.Printf("Failed to get environments. HTTP %d: %s\n", resp.StatusCode, string(bodyBytes))
		os.Exit(1)
	}

	// body, _ := ioutil.ReadAll(resp.Body)
	// log.Printf("envs: %s", body)

	// 3. Parse the JSON array and filter for 'active'
	var environments []Environment
	if err := json.NewDecoder(resp.Body).Decode(&environments); err != nil {
		log.Printf("Error parsing environments response: %v\n", err)
		os.Exit(1)
	}

	var filteredEnvironments []Environment
	for _, env := range environments {
		if env.IsPR {
			filteredEnvironments = append(filteredEnvironments, env)
		}
	}

	return filteredEnvironments
}
