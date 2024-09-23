package testutils

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	RoosterKidProxyListURL = "https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt"
	ProxlifyProxyListURL   = "https://raw.githubusercontent.com/proxifly/free-proxy-list/main/proxies/protocols/http/data.txt" // nolint: revive
)

func NewProxyClient(proxyListUrls ...string) (*http.Client, error) {
	if len(proxyListUrls) == 0 {
		return nil, fmt.Errorf("at least one proxy list url must be passed")
	}

	proxyTransport, err := newProxyTransport(proxyListUrls)
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: proxyTransport}, nil
}

func newProxyTransport(proxyListUrls []string) (http.RoundTripper, error) {
	proxyLists, err := downloadProxyLists(proxyListUrls)
	if err != nil {
		return nil, fmt.Errorf("error downloading proxy list: %w", err)
	}

	proxyLists = getWorkingProxies(proxyLists, 2*time.Second)
	if len(proxyLists) == 0 {
		return nil, errors.New("none of the proxies in the proxy list are working")
	}

	startingIndex := rand.Intn(len(proxyLists))

	return &proxyTransport{
		proxyList: proxyLists,
		index:     startingIndex,
	}, nil
}

type proxyTransport struct {
	proxyList []*url.URL
	index     int
	mu        sync.Mutex
}

func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Get new proxy and increment the index
	proxyURL := t.proxyList[t.index]
	t.index = (t.index + 1) % len(t.proxyList)

	// Create a new transport and us it for the request
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
		TLSClientConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	return transport.RoundTrip(request)
}

func downloadProxyLists(proxyListUrls []string) ([]*url.URL, error) {
	proxyUrls := []*url.URL{}
	for _, proxyListUrl := range proxyListUrls {
		resp, err := http.Get(proxyListUrl)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		proxyListProxyUrls, err := parseProxyList(resp.Body)
		if err != nil {
			return nil, err
		}

		proxyUrls = append(proxyUrls, proxyListProxyUrls...)
	}

	return proxyUrls, nil
}

func parseProxyList(reader io.Reader) ([]*url.URL, error) {
	var proxyUrls []*url.URL
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		proxyAddress := scanner.Text()
		if !strings.HasPrefix(proxyAddress, "http://") {
			proxyAddress = "http://" + proxyAddress
		}

		parsedURL, err := url.Parse(proxyAddress)
		if err != nil {
			fmt.Printf("Skipping invalid proxy URL: %s\n", proxyAddress)
			continue
		}

		proxyUrls = append(proxyUrls, parsedURL)
	}

	err := scanner.Err()
	if err != nil {
		return nil, err
	}

	return proxyUrls, nil
}

func getWorkingProxies(proxyUrls []*url.URL, timeout time.Duration) []*url.URL {
	resultsChan := make(chan *url.URL, len(proxyUrls))

	var wg sync.WaitGroup
	for _, proxyUrl := range proxyUrls {
		wg.Add(1)
		go func(proxyUrl *url.URL) {
			defer wg.Done()
			ok := checkProxy(proxyUrl, timeout)
			if ok {
				resultsChan <- proxyUrl
			}
		}(proxyUrl)
	}

	wg.Wait()
	close(resultsChan)

	results := make([]*url.URL, 0, len(proxyUrls))
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

func checkProxy(proxyUrl *url.URL, timeout time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	request, err := http.NewRequestWithContext(
		ctx, http.MethodGet,
		"https://example.com",
		http.NoBody,
	)
	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)
	if err != nil || response.StatusCode >= 400 {
		return false
	}

	return true
}
