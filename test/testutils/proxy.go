package testutils

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

const ProxyListURL = "https://raw.githubusercontent.com/roosterkid/openproxylist/main/HTTPS_RAW.txt"

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
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
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
