package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"text/template"

	"github.com/Kong/go-pdk"
	"github.com/Masterminds/sprig"
)

type Config struct {
	Uri    string
	Method string
	In     string
	Out    string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	sr := kong.Router
	s, err := sr.GetService()
	kongError(err, kong)

	reqPath, err := kong.Request.GetPath()
	kongError(err, kong)

	path, err := BuildPath(conf.Uri, s.Path, reqPath)
	kongError(err, kong)
	url := s.Protocol + "://" + s.Host + path

	rawBody, err := kong.Request.GetRawBody()
	kongError(err, kong)
	reqBody, err := BuildTmpl([]byte(rawBody), conf.In)
	kongError(err, kong)

	data, err := HttpReq(conf.Method, url, reqBody)
	kongError(err, kong)

	respBody, err := BuildTmpl(data, conf.Out)
	kongError(err, kong)

	headers := make(map[string][]string)
	kong.Response.Exit(200, respBody, headers)
}

func kongError(err error, kong *pdk.PDK) {
	if err != nil {
		kong.Log.Err(err.Error())
		kong.Response.ExitStatus(500)
	}
}

func sprigTmpl(n string) *template.Template {
	return template.New(n).Funcs(sprig.TxtFuncMap())
}

func BuildTmpl(data []byte, tmpl string) (compiled string, err error) {
	if tmpl == " " {
		return string(data), nil
	}
	var t *template.Template
	t, err = sprigTmpl("midnight-tmpl-engine").Parse(tmpl)
	if err != nil {
		return
	}
	params := make(map[string]interface{})
	err = json.Unmarshal(data, &params)
	if err != nil {
		return
	}
	var buf strings.Builder
	err = t.Execute(&buf, params)
	if err != nil {
		return
	}
	compiled = buf.String()
	return
}

func BuildPath(cfgUri, svcPath, reqPath string) (path string, err error) {
	if cfgUri == " " {
		return svcPath, nil
	}
	mapUri := make(map[string]int)
	err = json.Unmarshal([]byte(cfgUri), &mapUri)
	if err != nil {
		return
	}
	reqPaths := strings.Split(reqPath, "/")

	for k, v := range mapUri {
		reqValue := reqPaths[v+1]
		svcPath = strings.ReplaceAll(svcPath, ":"+k, reqValue)
	}
	path = svcPath

	return
}

func HttpReq(method, url, body string) (data []byte, err error) {
	client := &http.Client{}
	client.Transport = &http.Transport{
		DisableKeepAlives: true,
	}

	var req *http.Request
	if len(body) > 0 {
		req, err = http.NewRequest(method, url, strings.NewReader(body))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, resp.Body)

	if err != nil {
		return
	}
	data = buf.Bytes()

	return
}
