package http

import (
	"fmt"
	"github.com/mozilla-services/heka/message"
	. "github.com/mozilla-services/heka/pipeline"
	"strconv"
	"net/http"
	"github.com/mozilla-services/heka/plugins/tcp"
	"net/url"
	"regexp"
	"errors"
	"time"
	"bytes"
	"io/ioutil"
)

// Heka Filter plugin that can send a http request
// if a given regex matches in the response body then
// a "success" message will be created
// otherwise a "failure" message will be created
type HttpFilter struct {
	*HttpFilterConfig
	url          *url.URL
	client       *http.Client
	useBasicAuth bool
	Match        string
}

type HttpFilterConfig struct {
	HttpTimeout uint32 `toml:"http_timeout"`
	Address     string
	Method	    string
	Headers     http.Header
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	Tls         *tcp.TlsConfig
	MatchRegex string `toml:"match_regex"`
}

func (hf *HttpFilter) ConfigStruct() interface{} {
	return &HttpFilterConfig{
		HttpTimeout: 0,
		Headers:     make(http.Header),
	}
}

func (hf *HttpFilter) Init(config interface{}) (err error) {
	hf.HttpFilterConfig = config.(*HttpFilterConfig)
	
	hf.Match = hf.HttpFilterConfig.MatchRegex
	//if hf.Match, err = regexp.Compile(hf.HttpFilterConfig.MatchRegex); err != nil {
	//	err = fmt.Errorf("HttpFilter: %s", err)
	//	return
	//}
	
	if hf.url, err = url.Parse(hf.Address); err != nil {
		return fmt.Errorf("Can't parse URL '%s': %s", hf.Address, err.Error())
	}
	if hf.url.Scheme != "http" && hf.url.Scheme != "https" {
		return errors.New("`address` must contain an absolute http or https URL.")
	}
	hf.Method = "POST"
	hf.client = new(http.Client)
	if hf.HttpTimeout > 0 {
		hf.client.Timeout = time.Duration(hf.HttpTimeout) * time.Millisecond
	}
	if hf.Username != "" || hf.Password != "" {
		hf.useBasicAuth = true
	}
	if hf.url.Scheme == "https" && hf.Tls != nil {
		transport := &http.Transport{}
		if transport.TLSClientConfig, err = tcp.CreateGoTlsConfig(hf.Tls); err != nil {
			return fmt.Errorf("TLS init error: %s", err.Error())
		}
		hf.client.Transport = transport
	}
	return
}

func (hf *HttpFilter) Run(fr FilterRunner, h PluginHelper) (err error) {
	var (
		success  bool
		values = make(map[string]string)
		val      string
	)

	inChan := fr.InChan()

	for pack := range inChan {
		values["Payload"] = pack.Message.GetPayload()
		
		for _, field := range pack.Message.Fields {
			// It's painful to be converting these numeric values to strings,
			// but for now it's the only way to get numeric data into the stat
			// accumulator.
			if field.GetValueType() == message.Field_STRING && len(field.ValueString) > 0 {
				val = field.ValueString[0]
			} else if field.GetValueType() == message.Field_DOUBLE {
				val = strconv.FormatFloat(field.ValueDouble[0], 'f', -1, 64)
			} else if field.GetValueType() == message.Field_INTEGER {
				val = strconv.FormatInt(field.ValueInteger[0], 10)
			}
			values[field.GetName()] = val
		}
		
		if success = hf.request(fr, hf.Match, []byte(values["Payload"])); success {
			// change message to success
			pack.Message.SetType("http.success")
		} else{
			// change message to failure
			pack.Message.SetType("http.failure")
		}
                
                fr.Inject(pack) //after changing the type, inject the message back to the heka pipeline
	}

	return
}

func (hf *HttpFilter) request(fr FilterRunner, re string, outBytes []byte) (matched bool) {
	var(
		resp       *http.Response
		err        error
	)

	req := &http.Request{
		Method: hf.Method,
		URL:    hf.url,
		Header: hf.Headers,
	}
	
	if hf.useBasicAuth {
		req.SetBasicAuth(hf.Username, hf.Password)
	}
	
	req.Body = ioutil.NopCloser(bytes.NewReader(outBytes))
	req.ContentLength = int64(len(outBytes))

	if resp, err = hf.client.Do(req); err != nil {
		return false
	}
	defer resp.Body.Close()

	var body []byte
		if resp.ContentLength > 0 {
			body = make([]byte, resp.ContentLength)
			resp.Body.Read(body)
		}
	
	if resp.StatusCode >= 400 {
		return false
	}
	
	matched, err = regexp.MatchString(re, string(body))
     
	return matched
}

func init() {
	RegisterPlugin("HttpFilter", func() interface{} {
		return new(HttpFilter)
	})
}
