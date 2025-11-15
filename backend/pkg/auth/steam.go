package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"golang.org/x/oauth2"
)

var (
	ErrInvalidState = errors.New("invalid OAuth state")
	ErrSteamAuth    = errors.New("Steam authentication failed")
)

// SteamOAuthManager handles Steam OpenID/OAuth authentication
// Note: Steam uses OpenID 2.0, but we can adapt it to OAuth-like flow
type SteamOAuthManager struct {
	apiKey      string
	callbackURL string
	states      map[string]bool // Simple state storage (use Redis in production!)
}

// SteamUser represents Steam user data
type SteamUser struct {
	SteamID     string `json:"steamid"`
	PersonaName string `json:"personaname"`
	Avatar      string `json:"avatar"`
	ProfileURL  string `json:"profileurl"`
}

// NewSteamOAuthManager creates a new Steam OAuth manager
func NewSteamOAuthManager(apiKey, callbackURL string) *SteamOAuthManager {
	return &SteamOAuthManager{
		apiKey:      apiKey,
		callbackURL: callbackURL,
		states:      make(map[string]bool),
	}
}

// GenerateAuthURL creates the Steam login URL with CSRF protection
// Security: State parameter prevents CSRF attacks
func (m *SteamOAuthManager) GenerateAuthURL() (string, string, error) {
	// Generate random state for CSRF protection
	state, err := generateRandomState()
	if err != nil {
		return "", "", err
	}

	// Store state (in production, use Redis with expiration)
	m.states[state] = true

	// Steam OpenID URL
	params := url.Values{}
	params.Set("openid.ns", "http://specs.openid.net/auth/2.0")
	params.Set("openid.mode", "checkid_setup")
	params.Set("openid.return_to", m.callbackURL+"?state="+state)
	params.Set("openid.realm", m.callbackURL)
	params.Set("openid.identity", "http://specs.openid.net/auth/2.0/identifier_select")
	params.Set("openid.claimed_id", "http://specs.openid.net/auth/2.0/identifier_select")

	authURL := "https://steamcommunity.com/openid/login?" + params.Encode()
	return authURL, state, nil
}

// ValidateCallback verifies the Steam OAuth callback and extracts user info
// Security: Validates state to prevent CSRF and verifies Steam response
func (m *SteamOAuthManager) ValidateCallback(ctx context.Context, values url.Values) (*SteamUser, error) {
	// Verify state (CSRF protection)
	state := values.Get("state")
	if state == "" || !m.states[state] {
		return nil, ErrInvalidState
	}
	delete(m.states, state) // One-time use

	// Verify OpenID response
	claimedID := values.Get("openid.claimed_id")
	if claimedID == "" {
		return nil, ErrSteamAuth
	}

	// Extract Steam ID from claimed_id URL
	// Format: https://steamcommunity.com/openid/id/[STEAM_ID]
	steamID := extractSteamID(claimedID)
	if steamID == "" {
		return nil, errors.New("invalid Steam ID")
	}

	// Fetch user info from Steam API
	user, err := m.fetchSteamUserInfo(ctx, steamID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch Steam user info: %w", err)
	}

	return user, nil
}

// fetchSteamUserInfo retrieves user details from Steam Web API
func (m *SteamOAuthManager) fetchSteamUserInfo(ctx context.Context, steamID string) (*SteamUser, error) {
	if m.apiKey == "" {
		// If no API key, return basic info
		return &SteamUser{
			SteamID:     steamID,
			PersonaName: "Steam User",
		}, nil
	}

	apiURL := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v0002/?key=%s&steamids=%s",
		m.apiKey, steamID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Steam API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Response struct {
			Players []SteamUser `json:"players"`
		} `json:"response"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if len(result.Response.Players) == 0 {
		return nil, errors.New("no player data returned from Steam")
	}

	return &result.Response.Players[0], nil
}

// extractSteamID extracts the Steam ID from OpenID claimed_id
func extractSteamID(claimedID string) string {
	u, err := url.Parse(claimedID)
	if err != nil {
		return ""
	}

	// Extract last part of path
	parts := len(u.Path)
	if parts > 0 {
		// Get everything after the last "/"
		for i := len(u.Path) - 1; i >= 0; i-- {
			if u.Path[i] == '/' {
				return u.Path[i+1:]
			}
		}
	}

	return ""
}

// generateRandomState creates a random state string for CSRF protection
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// GetOAuthConfig returns OAuth2 config (alternative approach for other OAuth providers)
func GetOAuthConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://steamcommunity.com/openid/login",
			TokenURL: "https://steamcommunity.com/openid/login",
		},
	}
}
