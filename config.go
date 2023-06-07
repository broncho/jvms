package main

import (
	"errors"
	"github.com/tucnak/store"
	"github.com/ystyle/jvms/utils/file"
	"github.com/ystyle/jvms/utils/jdk"
	"github.com/ystyle/jvms/utils/web"
	"path/filepath"
)

const AppName string = "jvm"
const ConfigName string = "jvm.json"
const version = "3.0.0"

type Config struct {
	JavaHome          string `json:"java_home"`
	CurrentJDKVersion string `json:"current_jdk_version"`
	Originalpath      string `json:"original_path"`
	Proxy             string `json:"proxy"`
	store             string `json:"store"`
	download          string `json:"download"`
}

var config Config

func initConfig() error {
	store.Init(AppName)
	if err := store.Load(ConfigName, &config); err != nil {
		return errors.New("failed to load the config:" + err.Error())
	}
	s := file.GetCurrentPath()
	config.store = filepath.Join(s, "store")
	config.download = filepath.Join(s, "download")
	if config.Originalpath == "" {
		config.Originalpath = jdk.DefaultOriginalPath
	}
	if config.Proxy != "" {
		web.SetProxy(config.Proxy)
	}
	return nil
}

func storeConfig() error {
	if err := store.Save(ConfigName, &config); err != nil {
		return errors.New("failed to save the config:" + err.Error())
	}
	return nil
}
