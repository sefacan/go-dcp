package couchbase

import (
	"encoding/base64"
	"fmt"

	"github.com/Trendyol/go-dcp/config"

	jsoniter "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type PoolsResult struct {
	ImplementationVersion string `json:"implementationVersion"`
}

type BucketResult struct {
	BucketType     string `json:"bucketType"`
	StorageBackend string `json:"storageBackend"`
}

type HTTPClient interface {
	Connect() error
	GetVersion() (string, error)
	GetBucketInformation() (*BucketResult, error)
}

type httpClient struct {
	config     *config.Dcp
	httpClient *fasthttp.Client
	client     Client
	baseURL    string
}

func (h *httpClient) Connect() error {
	pingResult, err := h.client.Ping()
	if err != nil {
		return err
	}

	h.baseURL = pingResult.MgmtEndpoint

	return nil
}

func (h *httpClient) doRequest(req *fasthttp.Request, v interface{}) error {
	req.Header.Set(
		"Authorization",
		"Basic "+base64.StdEncoding.EncodeToString([]byte(h.config.Username+":"+h.config.Password)),
	)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	err := h.httpClient.Do(req, res)
	if err != nil {
		return err
	}

	err = jsoniter.Unmarshal(res.Body(), v)
	if err != nil {
		return err
	}

	return nil
}

func (h *httpClient) GetVersion() (string, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%v/pools", h.baseURL))
	req.Header.SetMethod("GET")

	var result PoolsResult
	err := h.doRequest(req, &result)
	if err != nil {
		return "", err
	}

	return result.ImplementationVersion, nil
}

func (h *httpClient) GetBucketInformation() (*BucketResult, error) {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%v/pools/default/buckets/%v", h.baseURL, h.config.BucketName))
	req.Header.SetMethod("GET")

	var result BucketResult
	err := h.doRequest(req, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func NewHTTPClient(config *config.Dcp, client Client) HTTPClient {
	return &httpClient{
		config:     config,
		httpClient: &fasthttp.Client{},
		client:     client,
	}
}
