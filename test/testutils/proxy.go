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

func NewProxyClient(proxyListUrl string) (*http.Client, error) {
	proxyTransport, err := newProxyTransport(proxyListUrl)
	if err != nil {
		return nil, err
	}

	return &http.Client{Transport: proxyTransport}, nil
}

func newProxyTransport(proxyListUrl string) (http.RoundTripper, error) {
	proxyList, err := downloadProxyList(proxyListUrl)
	if err != nil {
		return nil, fmt.Errorf("error downloading proxy list: %w", err)
	}

	proxyList = getWorkingProxies(proxyList, 2*time.Second)
	if len(proxyList) == 0 {
		return nil, errors.New("none of the proxies in the proxy list are working")
	}

	startingIndex := rand.Intn(len(proxyList))

	return &proxyTransport{
		proxyList: proxyList,
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

func downloadProxyList(urlString string) ([]*url.URL, error) {
	resp, err := http.Get(urlString)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return parseProxyList(resp.Body)
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
