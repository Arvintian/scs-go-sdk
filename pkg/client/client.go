package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//Client scs http client
type Client struct {
	endpoint  string
	accesskey string
	secretkey string
	hc        *http.Client
}

//Request scs http request
type Request struct {
	Method  string
	Bucket  string
	Path    string
	Params  url.Values
	Headers http.Header
	Body    io.Reader
	//internal
	baseuri  string
	signpath string
	prepared bool
}

var sigleParams = map[string]bool{
	"acl":       true,
	"meta":      true,
	"multipart": true,
	"relax":     true,
}

// HTTPTimeout defines HTTP timeout.
type HTTPTimeout struct {
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	HeaderTimeout    time.Duration
	LongTimeout      time.Duration
	IdleConnTimeout  time.Duration
}

// HTTPMaxConns defines max idle connections and max idle connections per host
type HTTPMaxConns struct {
	MaxIdleConns        int
	MaxIdleConnsPerHost int
}

//NewClient return a client point
func NewClient(accesskey, secretkey, endpoint string) *Client {

	httpTimeOut := HTTPTimeout{}
	httpTimeOut.ConnectTimeout = time.Second * 30   // 30s
	httpTimeOut.ReadWriteTimeout = time.Second * 60 // 60s
	httpTimeOut.HeaderTimeout = time.Second * 60    // 60s
	httpTimeOut.LongTimeout = time.Second * 300     // 300s
	httpTimeOut.IdleConnTimeout = time.Second * 50  // 50s
	httpMaxConns := HTTPMaxConns{}
	httpMaxConns.MaxIdleConns = 100
	httpMaxConns.MaxIdleConnsPerHost = 100

	return &Client{
		accesskey: accesskey,
		secretkey: secretkey,
		endpoint:  endpoint,
		hc: &http.Client{
			Transport: &http.Transport{
				Dial: func(netw, addr string) (net.Conn, error) {
					d := net.Dialer{
						Timeout:   httpTimeOut.ConnectTimeout,
						KeepAlive: 30 * time.Second,
					}
					conn, err := d.Dial(netw, addr)
					if err != nil {
						return nil, err
					}
					return newTimeoutConn(conn, httpTimeOut.ReadWriteTimeout, httpTimeOut.LongTimeout), nil
				},
				MaxIdleConns:          httpMaxConns.MaxIdleConns,
				MaxIdleConnsPerHost:   httpMaxConns.MaxIdleConnsPerHost,
				IdleConnTimeout:       httpTimeOut.IdleConnTimeout,
				ResponseHeaderTimeout: httpTimeOut.HeaderTimeout,
			},
		},
	}
}

//Query scs http query
func (c *Client) Query(req *Request) (http.Header, io.ReadCloser, error) {
	err := c.prepare(req)
	if err != nil {
		return nil, ioutil.NopCloser(bytes.NewBuffer([]byte{})), err
	}
	hresp, err := c.run(req)
	if err != nil || hresp == nil {
		return nil, ioutil.NopCloser(bytes.NewBuffer([]byte{})), err
	}
	return hresp.Header, hresp.Body, nil
}

func (c *Client) prepare(req *Request) error {
	if !req.prepared {
		req.prepared = true
		if req.Method == "" {
			req.Method = "GET"
		}
		// Copy so they can be mutated without affecting on retries.

		params := make(url.Values)
		headers := make(http.Header)
		for k, v := range req.Params {
			params[k] = v
		}
		for k, v := range req.Headers {
			headers[k] = v
		}
		req.Params = params
		req.Headers = headers
		if !strings.HasPrefix(req.Path, "/") {
			req.Path = "/" + req.Path
		}

		if req.Bucket != "" {
			if strings.IndexAny(req.Bucket, "/:@") >= 0 {
				return fmt.Errorf("bad S3 bucket: %q", req.Bucket)
			}
			req.signpath = "/" + req.Bucket + (&url.URL{Path: req.Path}).RequestURI()
		} else {
			req.signpath = (&url.URL{Path: req.Path}).RequestURI()
		}
		req.baseuri = c.endpoint
		req.baseuri = strings.Replace(req.baseuri, "$", req.Bucket, -1)
	}
	u, err := url.Parse(req.baseuri)
	if err != nil {
		return fmt.Errorf("bad S3 endpoint URL %q: %v", req.baseuri, err)
	}
	req.Headers["Host"] = []string{u.Host}
	req.Headers["Date"] = []string{time.Now().In(time.UTC).Format(time.RFC1123)}
	req.Headers["User-Agent"] = []string{"s3gosdk-1.0"}
	sign(*c, req.Method, req.signpath, req.Params, req.Headers)
	return nil
}

func (c *Client) run(req *Request) (hresp *http.Response, err error) {
	// fmt.Println(req.Headers)
	// fmt.Println(req.urlencode())
	u, err := req.urlencode()
	if err != nil {
		return nil, err
	}
	hreq := http.Request{
		URL:    u,
		Method: req.Method,
		Header: req.Headers,
		//Close:  true,
	}
	if v, ok := req.Headers["Content-Length"]; ok {
		hreq.ContentLength, _ = strconv.ParseInt(v[0], 10, 64)
		delete(req.Headers, "Content-Length")
	}
	htCli := c.hc
	if req.Body != nil {
		hreq.Body = ioutil.NopCloser(req.Body)
	}
	hresp, err = htCli.Do(&hreq)
	if err != nil {
		return nil, err
	}
	if hresp.StatusCode != 200 && hresp.StatusCode != 204 && hresp.StatusCode != 206 {
		return nil, buildError(hresp)
	}
	return hresp, nil
}

func (req *Request) urlencode() (*url.URL, error) {
	var sigleArray []string
	var value = url.Values{}
	u, err := url.Parse(req.baseuri)
	if err != nil {
		return nil, fmt.Errorf("bad S3 endpoint URL %q: %v", req.baseuri, err)
	}
	for k, v := range req.Params {
		if sigleParams[k] {
			sigleArray = append(sigleArray, k)
		} else {
			value.Add(k, v[0])
		}
	}
	switch {
	case len(sigleArray) > 0 && len(value) > 0:
		u.RawQuery = strings.Join(sigleArray, "&") + "&" + value.Encode()
	case len(sigleArray) <= 0:
		u.RawQuery = value.Encode()
	default:
		u.RawQuery = strings.Join(sigleArray, "&")
	}
	re := regexp.MustCompile(req.Bucket)
	if re.MatchString(u.Host) {
		u.Path = req.Path
	} else {
		u.Path = "/" + req.Bucket + req.Path
	}
	return u, nil
}

//Error scs client error
type Error struct {
	StatusCode int
	RequestID  string
	ErrorCode  string
	Date       string
}

func (e *Error) Error() string {
	return e.ErrorCode
}

func buildError(r *http.Response) error {
	var err Error
	err.StatusCode = r.StatusCode
	err.RequestID = r.Header["X-Requestid"][0]
	if ErrCode, ok := r.Header["X-Error-Code"]; ok {
		err.ErrorCode = ErrCode[0]
	} else {
		err.ErrorCode = strconv.FormatInt(int64(r.StatusCode), 10)
	}
	err.Date = r.Header["Date"][0]
	return &err
}
