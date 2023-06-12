package cmd

import (
	"errors"
	"fmt"
	"github.com/tucnak/store"
	"github.com/ystyle/jvms/utils/file"
	"github.com/ystyle/jvms/utils/jdk"
	"github.com/ystyle/jvms/utils/web"
	"path/filepath"
)

const AppName string = "jvms"
const AppConfigName string = "jvms.json"
const JdkVersionCacheName = "versions.json"
const Version = "3.0.0"
const JdkBaseName = "jdk"
const DownloadBaseName = "download"
const StoreBaseName = "store"

type Config struct {
	JvmsHome          string `json:"jvms_home"`
	JavaHome          string `json:"java_home"`
	CurrentJDKVersion string `json:"current_jdk_version"`
	OriginalPath      string `json:"original_path"`
	Proxy             string `json:"proxy"`
	store             string
	download          string
}

var AppConfig Config

func InitConfig() error {
	store.Init(AppName)
	if err := store.Load(AppConfigName, &AppConfig); err != nil {
		return errors.New("failed to load the config:" + err.Error())
	}
	workHome := file.GetCurrentPath()
	if AppConfig.JvmsHome != "" {
		workHome = AppConfig.JvmsHome
	} else {
		AppConfig.JvmsHome = workHome
	}
	if AppConfig.JavaHome == "" {
		AppConfig.JavaHome = filepath.Join(AppConfig.JvmsHome, JdkBaseName)
	}
	AppConfig.store = filepath.Join(workHome, StoreBaseName)
	AppConfig.download = filepath.Join(workHome, DownloadBaseName)
	if AppConfig.OriginalPath == "" {
		AppConfig.OriginalPath = jdk.DefaultOriginalPath
	}
	if AppConfig.Proxy != "" {
		web.SetProxy(AppConfig.Proxy)
	}
	return nil
}

func StoreConfig() error {
	if err := store.Save(AppConfigName, &AppConfig); err != nil {
		return errors.New("failed to save the config:" + err.Error())
	}
	return nil
}

func ShowConfig() string {
	return fmt.Sprintf("%+v", AppConfig)
}

func CachePutVersion(versions []jdk.JdkVersion) {
	store.Init(AppName)
	err := store.Save(JdkVersionCacheName, &versions)
	if err != nil {
		return
	}
}

func CacheGetVersion() []jdk.JdkVersion {
	store.Init(AppName)
	var versions []jdk.JdkVersion
	_ = store.Load(JdkVersionCacheName, &versions)
	return versions
}
