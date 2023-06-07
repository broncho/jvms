package jdk

import (
	"github.com/baneeishaque/adoptium_jdk_go"
	"strings"
)

const vendor = "Eclipse"
const url = "https://api.adoptium.net/"

type AdoptiumJdkSource struct {
	vendor string
	url    string
}

func NewAdoptiumJdkSource() *AdoptiumJdkSource {
	return &AdoptiumJdkSource{
		vendor: "Eclipse",
		url:    "https://api.adoptium.net/",
	}
}

func (receiver *AdoptiumJdkSource) SourceName() string {
	return receiver.vendor
}

func (receiver *AdoptiumJdkSource) SourceUrl() string {
	return receiver.url
}

func (receiver *AdoptiumJdkSource) JdkVersions() []JdkVersion {
	var versions []JdkVersion
	adoptiumJdks := strings.Split(adoptium_jdk_go.ApiListReleases(), "\n")
	for _, adoptiumJdkUrl := range adoptiumJdks {
		fileSeparatorIndex := strings.LastIndex(adoptiumJdkUrl, "/")
		fileName := adoptiumJdkUrl[fileSeparatorIndex+1:]
		fileVersion := strings.TrimSuffix(fileName, ".zip")
		versions = append(versions, JdkVersion{Version: fileVersion, Url: adoptiumJdkUrl, Source: receiver.SourceName()})
	}
	return versions
}
