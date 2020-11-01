package main

import (
	"github.com/Kong/go-pdk"
)

type Config struct {
	Uri string
	In  string
	Out string
}

func New() interface{} {
	return &Config{}
}

func (conf Config) Access(kong *pdk.PDK) {
	kong.Log.Info(conf.Uri)

	r := kong.Router
	s, _ := r.GetService()
	kong.Log.Info(s.Path)

	headers := make(map[string][]string)
	kong.Response.Exit(200, s.Path, headers)
}
