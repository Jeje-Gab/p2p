package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Client struct {
	clientID     string
	clientSecret string
	http         *http.Client
	token        string
	expiry       time.Time
	Market       string // opcional: "BR"
}

func New(clientID, clientSecret string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		http:         &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) ensureToken(ctx context.Context) error {
	if c.token != "" && time.Now().Before(c.expiry.Add(-30*time.Second)) {
		return nil
	}
	form := url.Values{}
	form.Set("grant_type", "client_credentials")

	req, _ := http.NewRequestWithContext(ctx, "POST", "https://accounts.spotify.com/api/token", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	basic := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+basic)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("spotify token status %d: %s", resp.StatusCode, string(b))
	}
	var tk struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tk); err != nil {
		return err
	}
	c.token = tk.AccessToken
	c.expiry = time.Now().Add(time.Duration(tk.ExpiresIn) * time.Second)
	return nil
}

// ---------------------------------------------------------
//
//	A) Busca "parcial" (MANTIDA): SearchPlaylistByName
//
// ---------------------------------------------------------
func (c *Client) SearchPlaylistByName(ctx context.Context, q string) (id, name string, err error) {
	if err = c.ensureToken(ctx); err != nil {
		return "", "", err
	}
	u := "https://api.spotify.com/v1/search?type=playlist&limit=1&q=" + url.QueryEscape(q)
	if c.Market != "" {
		u += "&market=" + url.QueryEscape(c.Market)
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", "", fmt.Errorf("spotify search status %d", resp.StatusCode)
	}
	var sr struct {
		Playlists struct {
			Items []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"items"`
		} `json:"playlists"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", "", err
	}
	if len(sr.Playlists.Items) == 0 {
		return "", "", errors.New("no playlist found")
	}
	return sr.Playlists.Items[0].ID, sr.Playlists.Items[0].Name, nil
}

// ---------------------------------------------------------
//
//	B) Busca "exata": SearchPlaylistByNameExact
//	   - usa aspas na query
//	   - normaliza nome (case/acentos/pontuação) para igualdade
//
// ---------------------------------------------------------
func (c *Client) SearchPlaylistByNameExact(ctx context.Context, name string) (id, display string, err error) {
	if err = c.ensureToken(ctx); err != nil {
		return "", "", err
	}

	q := `"` + name + `"` // aspas para frase exata
	v := url.Values{}
	v.Set("type", "playlist")
	v.Set("limit", "50")
	v.Set("q", q)
	if c.Market != "" {
		v.Set("market", c.Market)
	}
	u := "https://api.spotify.com/v1/search?" + v.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return "", "", fmt.Errorf("spotify search status %d: %s", resp.StatusCode, string(b))
	}

	var sr struct {
		Playlists struct {
			Items []struct {
				ID   string `json:"id"`
				Name string `json:"name"`
			} `json:"items"`
		} `json:"playlists"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&sr); err != nil {
		return "", "", err
	}
	if len(sr.Playlists.Items) == 0 {
		return "", "", errors.New("no playlist found")
	}

	target := slug(name)
	for _, it := range sr.Playlists.Items {
		if slug(it.Name) == target {
			return it.ID, it.Name, nil
		}
	}
	return "", "", errors.New("no exact playlist match")
}

// (opcional) wrapper para alternar por flag
func (c *Client) SearchPlaylist(ctx context.Context, name string, exact bool) (id, display string, err error) {
	if exact {
		return c.SearchPlaylistByNameExact(ctx, name)
	}
	return c.SearchPlaylistByName(ctx, name)
}

// ---------- helpers ----------
func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	ns, _, _ := transform.String(t, s)
	var b strings.Builder
	for _, r := range ns {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// Mantém GetPlaylistTracks como você já tinha:
type Track struct {
	Name   string
	Artist string
	Link   string
}

func (c *Client) GetPlaylistTracks(ctx context.Context, playlistID string, limit int) ([]Track, error) {
	if err := c.ensureToken(ctx); err != nil {
		return nil, err
	}
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	u := fmt.Sprintf("https://api.spotify.com/v1/playlists/%s/tracks?fields=items(track(name,external_urls,artists(name,external_urls)))&limit=%d", url.PathEscape(playlistID), limit)
	if c.Market != "" {
		u += "&market=" + url.QueryEscape(c.Market)
	}
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("spotify tracks status %d: %s", resp.StatusCode, string(b))
	}

	var pr struct {
		Items []struct {
			Track struct {
				Name         string            `json:"name"`
				ExternalURLs map[string]string `json:"external_urls"`
				Artists      []struct {
					Name         string            `json:"name"`
					ExternalURLs map[string]string `json:"external_urls"`
				} `json:"artists"`
			} `json:"track"`
		} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}
	out := make([]Track, 0, len(pr.Items))
	for _, it := range pr.Items {
		if it.Track.Name == "" || len(it.Track.Artists) == 0 {
			continue
		}
		artist := it.Track.Artists[0]
		link := artist.ExternalURLs["spotify"]
		if link == "" {
			link = it.Track.ExternalURLs["spotify"]
		}
		out = append(out, Track{
			Name:   it.Track.Name,
			Artist: artist.Name,
			Link:   link,
		})
	}
	return out, nil
}
