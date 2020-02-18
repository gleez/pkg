// Copyright 2013 The Beego Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package httplib is used as http.Client
// Usage:
//
// import "github.com/gleez/app/pkg/httplib"
//
//	r := httplib.Post("http://gleez.com/")
//	r.Param("username","astaxie")
//	r.Param("password","123456")
//	r.PostFile("uploadfile1", "httplib.pdf")
//	r.PostFile("uploadfile2", "httplib.txt")
//	str, err := r.String()
//	if err != nil {
//		t.Fatal(err)
//	}
//	fmt.Println(str)
//
package httplib

import (
	"bytes"
	"compress/gzip"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

var defaultCookieJar http.CookieJar
var settingMutex sync.Mutex

// createDefaultCookie creates a global cookiejar to store cookies.
func createDefaultCookie() {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultCookieJar, _ = cookiejar.New(nil)
}

// SetDefaultSetting overwrites default settings
func SetDefaultSetting(setting Settings) {
	settingMutex.Lock()
	defer settingMutex.Unlock()
	defaultSetting = setting
}

// newRequest returns *Request with specific method
func newRequest(rawurl, method string) *Request {
	var resp http.Response
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Println("Httplib:", err)
	}
	req := http.Request{
		URL:        u,
		Method:     method,
		Header:     make(http.Header),
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
	}

	return &Request{
		url:     rawurl,
		req:     &req,
		params:  map[string][]string{},
		files:   map[string]string{},
		setting: defaultSetting,
		resp:    &resp,
	}
}

// NewRequest returns *Request with specific method
func NewRequest(url, method string) *Request {
	return newRequest(url, method)
}

// Get returns *Request with GET method.
func Get(url string) *Request {
	return newRequest(url, "GET")
}

// Post returns *Request with POST method.
func Post(url string) *Request {
	return newRequest(url, "POST")
}

// Put returns *Request with PUT method.
func Put(url string) *Request {
	return newRequest(url, "PUT")
}

// Delete returns *Request DELETE method.
func Delete(url string) *Request {
	return newRequest(url, "DELETE")
}

// Head returns *Request with HEAD method.
func Head(url string) *Request {
	return newRequest(url, "HEAD")
}

// Settings is the default settings for http client
type Settings struct {
	ShowDebug        bool
	UserAgent        string
	ConnectTimeout   time.Duration
	ReadWriteTimeout time.Duration
	TLSClientConfig  *tls.Config
	Proxy            func(*http.Request) (*url.URL, error)
	Transport        http.RoundTripper
	CheckRedirect    func(req *http.Request, via []*http.Request) error
	EnableCookie     bool
	Gzip             bool
	DumpBody         bool
	Retries          int // if set to -1 means will retry forever
}

// Request provides more useful methods for requesting one url than http.Request.
type Request struct {
	url     string
	req     *http.Request
	params  map[string][]string
	files   map[string]string
	setting Settings
	resp    *http.Response
	body    []byte
	dump    []byte
}

// GetRequest return the request object
func (r *Request) GetRequest() *http.Request {
	return r.req
}

// Setting changes request settings
func (r *Request) Setting(setting Settings) *Request {
	r.setting = setting
	return r
}

// SetBasicAuth sets the request's Authorization header to use HTTP Basic Authentication with the provided username and password.
func (r *Request) SetBasicAuth(username, password string) *Request {
	r.req.SetBasicAuth(username, password)
	return r
}

// SetEnableCookie sets enable/disable cookiejar
func (r *Request) SetEnableCookie(enable bool) *Request {
	r.setting.EnableCookie = enable
	return r
}

// SetUserAgent sets User-Agent header field
func (r *Request) SetUserAgent(useragent string) *Request {
	r.setting.UserAgent = useragent
	return r
}

// Debug sets show debug or not when executing request.
func (r *Request) Debug(isdebug bool) *Request {
	r.setting.ShowDebug = isdebug
	return r
}

// Retries sets Retries times.
// default is 0 means no retried.
// -1 means retried forever.
// others means retried times.
func (r *Request) Retries(times int) *Request {
	r.setting.Retries = times
	return r
}

// DumpBody setting whether need to Dump the Body.
func (r *Request) DumpBody(isdump bool) *Request {
	r.setting.DumpBody = isdump
	return r
}

// DumpRequest return the DumpRequest
func (r *Request) DumpRequest() []byte {
	return r.dump
}

// SetTimeout sets connect time out and read-write time out for BeegoRequest.
func (r *Request) SetTimeout(connectTimeout, readWriteTimeout time.Duration) *Request {
	r.setting.ConnectTimeout = connectTimeout
	r.setting.ReadWriteTimeout = readWriteTimeout
	return r
}

// SetTLSClientConfig sets tls connection configurations if visiting https url.
func (r *Request) SetTLSClientConfig(config *tls.Config) *Request {
	r.setting.TLSClientConfig = config
	return r
}

// Header add header item string in request.
func (r *Request) Header(key, value string) *Request {
	r.req.Header.Set(key, value)
	return r
}

// HeaderWithSensitiveCase add header item in request and keep the case of the header key.
func (r *Request) HeaderWithSensitiveCase(key, value string) *Request {
	r.req.Header[key] = []string{value}
	return r
}

// Headers returns headers in request.
func (r *Request) Headers() http.Header {
	return r.req.Header
}

// SetHost set the request host
func (r *Request) SetHost(host string) *Request {
	r.req.Host = host
	return r
}

// SetProtocolVersion sets the protocol version for incoming requests.
// Client requests always use HTTP/1.1.
func (r *Request) SetProtocolVersion(vers string) *Request {
	if len(vers) == 0 {
		vers = "HTTP/1.1"
	}

	major, minor, ok := http.ParseHTTPVersion(vers)
	if ok {
		r.req.Proto = vers
		r.req.ProtoMajor = major
		r.req.ProtoMinor = minor
	}

	return r
}

// SetCookie add cookie into request.
func (r *Request) SetCookie(cookie *http.Cookie) *Request {
	r.req.Header.Add("Cookie", cookie.String())
	return r
}

// SetTransport sets transport to
func (r *Request) SetTransport(transport http.RoundTripper) *Request {
	r.setting.Transport = transport
	return r
}

// SetProxy sets http proxy
// example:
//
//	func(req *http.Request) (*url.URL, error) {
// 		u, _ := url.ParseRequestURI("http://127.0.0.1:8118")
// 		return u, nil
// 	}
func (r *Request) SetProxy(proxy func(*http.Request) (*url.URL, error)) *Request {
	r.setting.Proxy = proxy
	return r
}

// SetCheckRedirect specifies the policy for handling redirects.
//
// If CheckRedirect is nil, the Client uses its default policy,
// which is to stop after 10 consecutive requests.
func (r *Request) SetCheckRedirect(redirect func(req *http.Request, via []*http.Request) error) *Request {
	r.setting.CheckRedirect = redirect
	return r
}

// Param adds query param in to request.
// params build query string as ?key1=value1&key2=value2...
func (r *Request) Param(key, value string) *Request {
	if param, ok := r.params[key]; ok {
		r.params[key] = append(param, value)
	} else {
		r.params[key] = []string{value}
	}
	return r
}

// PostFile uploads file via http
func (r *Request) PostFile(formname, filename string) *Request {
	r.files[formname] = filename
	return r
}

// Body adds request raw body.
// it supports string and []byte.
func (r *Request) Body(data interface{}) *Request {
	switch t := data.(type) {
	case string:
		bf := bytes.NewBufferString(t)
		r.req.Body = ioutil.NopCloser(bf)
		r.req.ContentLength = int64(len(t))
	case []byte:
		bf := bytes.NewBuffer(t)
		r.req.Body = ioutil.NopCloser(bf)
		r.req.ContentLength = int64(len(t))
	}
	return r
}

// XMLBody adds request raw body encoding by XML.
func (r *Request) XMLBody(obj interface{}) (*Request, error) {
	if r.req.Body == nil && obj != nil {
		byts, err := xml.Marshal(obj)
		if err != nil {
			return r, err
		}
		r.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		r.req.ContentLength = int64(len(byts))
		r.req.Header.Set("Content-Type", "application/xml")
	}
	return r, nil
}

// YAMLBody adds request raw body encoding by YAML.
func (r *Request) YAMLBody(obj interface{}) (*Request, error) {
	if r.req.Body == nil && obj != nil {
		byts, err := yaml.Marshal(obj)
		if err != nil {
			return r, err
		}
		r.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		r.req.ContentLength = int64(len(byts))
		r.req.Header.Set("Content-Type", "application/x+yaml")
	}
	return r, nil
}

// JSONBody adds request raw body encoding by JSON.
func (r *Request) JSONBody(obj interface{}) (*Request, error) {
	if r.req.Body == nil && obj != nil {
		byts, err := json.Marshal(obj)
		if err != nil {
			return r, err
		}
		r.req.Body = ioutil.NopCloser(bytes.NewReader(byts))
		r.req.ContentLength = int64(len(byts))
		r.req.Header.Set("Content-Type", "application/json")
	}
	return r, nil
}

func (r *Request) buildURL(paramBody string) {
	// build GET url with query string
	if r.req.Method == "GET" && len(paramBody) > 0 {
		if strings.Contains(r.url, "?") {
			r.url += "&" + paramBody
		} else {
			r.url = r.url + "?" + paramBody
		}
		return
	}

	// build POST/PUT/PATCH url and body
	if (r.req.Method == "POST" || r.req.Method == "PUT" || r.req.Method == "PATCH" || r.req.Method == "DELETE") && r.req.Body == nil {
		// with files
		if len(r.files) > 0 {
			pr, pw := io.Pipe()
			bodyWriter := multipart.NewWriter(pw)
			go func() {
				for formname, filename := range r.files {
					fileWriter, err := bodyWriter.CreateFormFile(formname, filename)
					if err != nil {
						log.Println("Httplib:", err)
					}
					fh, err := os.Open(filename)
					if err != nil {
						log.Println("Httplib:", err)
					}
					//iocopy
					_, err = io.Copy(fileWriter, fh)
					fh.Close()
					if err != nil {
						log.Println("Httplib:", err)
					}
				}
				for k, v := range r.params {
					for _, vv := range v {
						bodyWriter.WriteField(k, vv)
					}
				}
				bodyWriter.Close()
				pw.Close()
			}()
			r.Header("Content-Type", bodyWriter.FormDataContentType())
			r.req.Body = ioutil.NopCloser(pr)
			return
		}

		// with params
		if len(paramBody) > 0 {
			r.Header("Content-Type", "application/x-www-form-urlencoded")
			r.Body(paramBody)
		}
	}
}

func (r *Request) getResponse() (*http.Response, error) {
	if r.resp.StatusCode != 0 {
		return r.resp, nil
	}
	resp, err := r.DoRequest()
	if err != nil {
		return nil, err
	}
	r.resp = resp
	return resp, nil
}

func (r *Request) DoRequest() (resp *http.Response, err error) {
	var paramBody string
	if len(r.params) > 0 {
		var buf bytes.Buffer
		for k, v := range r.params {
			for _, vv := range v {
				buf.WriteString(url.QueryEscape(k))
				buf.WriteByte('=')
				buf.WriteString(url.QueryEscape(vv))
				buf.WriteByte('&')
			}
		}
		paramBody = buf.String()
		paramBody = paramBody[0 : len(paramBody)-1]
	}

	r.buildURL(paramBody)
	urlParsed, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}

	r.req.URL = urlParsed

	trans := r.setting.Transport

	if trans == nil {
		// create default transport
		trans = &http.Transport{
			TLSClientConfig:     r.setting.TLSClientConfig,
			Proxy:               r.setting.Proxy,
			Dial:                TimeoutDialer(r.setting.ConnectTimeout, r.setting.ReadWriteTimeout),
			MaxIdleConnsPerHost: 100,
		}
	} else {
		// if r.transport is *http.Transport then set the settings.
		if t, ok := trans.(*http.Transport); ok {
			if t.TLSClientConfig == nil {
				t.TLSClientConfig = r.setting.TLSClientConfig
			}
			if t.Proxy == nil {
				t.Proxy = r.setting.Proxy
			}
			if t.Dial == nil {
				t.Dial = TimeoutDialer(r.setting.ConnectTimeout, r.setting.ReadWriteTimeout)
			}
		}
	}

	var jar http.CookieJar
	if r.setting.EnableCookie {
		if defaultCookieJar == nil {
			createDefaultCookie()
		}

		jar = defaultCookieJar
	} else {
		jar = nil
	}

	client := &http.Client{
		Transport: trans,
		Jar:       jar,
	}

	if len(r.setting.UserAgent) > 0 && len(r.req.Header.Get("User-Agent")) == 0 {
		r.req.Header.Set("User-Agent", r.setting.UserAgent)
	}

	if r.setting.CheckRedirect != nil {
		client.CheckRedirect = r.setting.CheckRedirect
	}

	if r.setting.ShowDebug {
		dump, err := httputil.DumpRequest(r.req, r.setting.DumpBody)
		if err != nil {
			log.Println(err.Error())
		}

		r.dump = dump
	}

	// retries default value is 0, it will run once.
	// retries equal to -1, it will run forever until success
	// retries is setted, it will retries fixed times.
	for i := 0; r.setting.Retries == -1 || i <= r.setting.Retries; i++ {
		resp, err = client.Do(r.req)
		if err == nil {
			break
		}
	}

	return resp, err
}

// String returns the body string in response.
// it calls Response inner.
func (r *Request) String() (string, error) {
	data, err := r.Bytes()
	if err != nil {
		return "", err
	}

	return string(data), nil
}

// Bytes returns the body []byte in response.
// it calls Response inner.
func (r *Request) Bytes() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}
	resp, err := r.getResponse()
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, nil
	}

	defer resp.Body.Close()
	if r.setting.Gzip && resp.Header.Get("Content-Encoding") == "gzip" {
		reader, err := gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
		r.body, err = ioutil.ReadAll(reader)
		return r.body, err
	}

	r.body, err = ioutil.ReadAll(resp.Body)
	return r.body, err
}

// ToFile saves the body data in response to one file.
// it calls Response inner.
func (r *Request) ToFile(filename string) error {
	resp, err := r.getResponse()
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return nil
	}
	defer resp.Body.Close()
	err = pathExistAndMkdir(filename)
	if err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	return err
}

//Check that the file directory exists, there is no automatically created
func pathExistAndMkdir(filename string) (err error) {
	filename = path.Dir(filename)
	_, err = os.Stat(filename)
	if err == nil {
		return nil
	}
	if os.IsNotExist(err) {
		err = os.MkdirAll(filename, os.ModePerm)
		if err == nil {
			return nil
		}
	}
	return err
}

// ToJSON returns the map that marshals from the body bytes as json in response .
// it calls Response inner.
func (r *Request) ToJSON(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}

	err = json.Unmarshal(data, v)
	return err
}

// ToXML returns the map that marshals from the body bytes as xml in response .
// it calls Response inner.
func (r *Request) ToXML(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	err = xml.Unmarshal(data, v)
	return err
}

// ToYAML returns the map that marshals from the body bytes as yaml in response .
// it calls Response inner.
func (r *Request) ToYAML(v interface{}) error {
	data, err := r.Bytes()
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, v)
}

// Response executes request client gets response manually.
func (r *Request) Response() (*http.Response, error) {
	return r.getResponse()
}

// TimeoutDialer returns functions of connection dialer with timeout settings for http.Transport Dial field.
func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		conn, err := net.DialTimeout(netw, addr, cTimeout)
		if err != nil {
			return nil, err
		}
		return conn, conn.SetDeadline(time.Now().Add(rwTimeout))
	}
}
