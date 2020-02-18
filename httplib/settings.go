package httplib

import (
	"net"
	"net/http"
	"time"
)

var defaultSetting = Settings{
	UserAgent:        "GleezServer",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,
	Gzip:             true,
	DumpBody:         true,
}

var CustomSetting = Settings{
	UserAgent:        "GleezServer",
	ConnectTimeout:   60 * time.Second,
	ReadWriteTimeout: 60 * time.Second,

	Transport: &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 40 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		TLSHandshakeTimeout:   20 * time.Second,
		IdleConnTimeout:       20 * time.Second,
		ExpectContinueTimeout: 20 * time.Second,
	},

	Gzip:     true,
	DumpBody: true,
}
