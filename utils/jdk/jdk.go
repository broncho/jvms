package jdk

import (
	"fmt"
	"github.com/ystyle/jvms/utils/file"
	"io/ioutil"
)

type JdkVersion struct {
	Version string `json:"version"`
	Url     string `json:"url"`
	Source  string `json:"source"`
}

type JdkSource interface {
	SourceName() string
	SourceUrl() string
	JdkVersions() []JdkVersion
}

func GetInstalled(root string) []string {
	list := make([]string, 0)
	files, _ := ioutil.ReadDir(root)
	for i := len(files) - 1; i >= 0; i-- {
		if files[i].IsDir() {
			list = append(list, files[i].Name())
		}
	}
	return list
}

func IsVersionInstalled(root string, version string) bool {
	isInstalled := file.Exists(fmt.Sprintf("%s/%s/bin/javac.exe", root, version))
	return isInstalled
}

func RemoteJdkVersions() ([]JdkVersion, error) {
	var versions []JdkVersion
	var jdkSources = [...]JdkSource{
		NewAdoptiumJdkSource(),
		NewAzulJdkSource(),
	}
	for _, source := range jdkSources {
		versions = append(versions, source.JdkVersions()...)
	}
	return versions, nil
}
