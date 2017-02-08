package kcpraw

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/gopacket/layers"
)

var NoHTTP bool

var HTTPHost string

var DSCP int

type callback func()

type myMutex struct {
	sync.Mutex
}

func (m *myMutex) run(f callback) {
	m.Lock()
	defer m.Unlock()
	f()
}

// copy from stackoverflow

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var src = rand.NewSource(time.Now().UnixNano())

const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func randStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

var requestFormat string
var responseFromat string

func init() {
	var requestBuffer bytes.Buffer
	strs := []string{
		"POST /%s HTTP/1.1\r\n",
		"Accept: */*\r\n",
		"Accept-Encoding: */*\r\n",
		"Accept-Language: zh-CN\r\n",
		"Connection: keep-alive\r\n",
		"%s",
		"Content-Length:%d\r\n\r\n",
	}
	for _, str := range strs {
		requestBuffer.WriteString(str)
	}
	requestFormat = requestBuffer.String()
	var responseBuffer bytes.Buffer
	strs = []string{
		"HTTP/1.1 200 OK\r\n",
		"Cache-Control: private, no-store, max-age=0, no-cache\r\n",
		"Content-Type: text/html; charset=utf-8\r\n",
		"Content-Encoding: gzip\r\n",
		"Server: openresty/1.11.2\r\n",
		"Connection: keep-alive\r\n",
		"%s",
		"Content-Length: %d\r\n\r\n",
	}
	for _, str := range strs {
		responseBuffer.WriteString(str)
	}
	responseFromat = responseBuffer.String()
}

func buildHTTPRequest(headers string) string {
	return fmt.Sprintf(requestFormat, randStringBytesMaskImprSrc(10), headers, (src.Int63()%65536 + 10485760))
	// return fmt.Sprintf(requestFormat, randStringBytesMaskImprSrc(10), headers, 0)
}

func buildHTTPResponse(headers string) string {
	return fmt.Sprintf(responseFromat, headers, (src.Int63()%65536 + 104857600))
	// return fmt.Sprintf(responseFromat, headers, 0)
}

func fatalErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type timeoutErr struct {
	op string
}

func (t *timeoutErr) Error() string {
	return t.op + " timeout"
}

func (t *timeoutErr) Temporary() bool {
	return true
}

func (t *timeoutErr) Timeout() bool {
	return true
}

// FIXME
type pktLayers struct {
	eth *layers.Ethernet
	ip4 *layers.IPv4
	tcp *layers.TCP
}

const (
	SYNRECEIVED = 0
	WAITHTTPREQ = 1
	HTTPREPSENT = 2
	ESTABLISHED = 3
)

func getTCPOptions() []layers.TCPOption {
	return []layers.TCPOption{
		layers.TCPOption{
			OptionType:   layers.TCPOptionKindSACKPermitted,
			OptionLength: 2,
		},
	}
}

func checkTCPOtions(options []layers.TCPOption) (ok bool) {
	for _, v := range options {
		if v.OptionType == layers.TCPOptionKindSACKPermitted {
			ok = true
			break
		}
	}
	return
}

type connInfo struct {
	state uint32
	layer *pktLayers
	rep   []byte
}
