package jdk

import (
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
	"github.com/ystyle/jvms/utils/web"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

type AdoptiumJdkSource struct {
	vendor     string
	vendorHome string
	apiDoc     string
}

func NewAdoptiumJdkSource() *AdoptiumJdkSource {
	return &AdoptiumJdkSource{
		vendor:     "Eclipse",
		vendorHome: "https://adoptium.net/", //
		apiDoc:     "https://api.adoptium.net/q/swagger-ui",
	}
}

func (receiver *AdoptiumJdkSource) OriginName() string {
	return receiver.vendor
}

func (receiver *AdoptiumJdkSource) OriginDesc() string {
	return fmt.Sprintf("%s (%s)", receiver.vendorHome, receiver.apiDoc)
}

func (receiver *AdoptiumJdkSource) JdkVersions() []JdkVersion {
	var versions []JdkVersion
	query := AdoptiumQuery{
		OS:   runtime.GOOS,
		ARCH: runtime.GOARCH,
	}
	adoptiumJdks, err := QueryAdoptiumVersions(query)
	if err != nil {
		return versions
	}
	for _, adoptiumJdkUrl := range adoptiumJdks {
		fileSeparatorIndex := strings.LastIndex(adoptiumJdkUrl, "/")
		fileName := adoptiumJdkUrl[fileSeparatorIndex+1:]
		fileVersion := strings.TrimSuffix(fileName, ".zip")
		versions = append(versions, JdkVersion{Version: fileVersion, Url: adoptiumJdkUrl, Origin: receiver.OriginName()})
	}
	return versions
}

func buildAdoptiumQueryUrl(release string) string {
	const api = "https://api.adoptium.net/v3/assets/latest/$RELEASE/hotspot?vendor=eclipse"
	return strings.Replace(api, "$RELEASE", release, 1)
}

func QueryAdoptiumVersions(query AdoptiumQuery) ([]string, error) {
	const jqQuery = `.[] | .binary | select(.image_type == "jdk") | select(.architecture == "$ARCH") | select(.os == "$OS") | .package.link`
	filter := strings.Replace(jqQuery, "$ARCH", fixedArch(query.ARCH), 1)
	filter = strings.Replace(filter, "$OS", query.OS, 1)
	release := QueryAdoptiumRelease()

	var wx sync.WaitGroup
	wx.Add(len(release))
	var downloadUrls []string
	for _, v := range release {
		go func(v string, f string) {
			urls := queryAdoptiumVersionAndFilter(v, filter)
			downloadUrls = append(downloadUrls, urls...)
			defer wx.Done()
		}(v, filter)
	}
	wx.Wait()
	return downloadUrls, nil
}

func fixedArch(arch string) string {
	if arch == "amd64" {
		return "x64"
	} else {
		return arch
	}
}

func queryAdoptiumVersionAndFilter(release string, jqQuery string) []string {
	var downloadUrl []string
	url := buildAdoptiumQueryUrl(release)
	body, err := web.Call(url)
	if err != nil {
		return downloadUrl
	}
	var unmarshalledJson interface{}
	e := json.Unmarshal(body, &unmarshalledJson)
	if e != nil {
		return downloadUrl
	}
	queryRunner, err := gojq.Parse(jqQuery)
	if err != nil {
		return downloadUrl
	}
	iter := queryRunner.Run(unmarshalledJson)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err != nil {
			}
		}
		if num, ok := v.(float64); ok {
			downloadUrl = append(downloadUrl, fmt.Sprintf("%f", num))
		}
		if s, ok := v.(string); ok {
			downloadUrl = append(downloadUrl, s)
		}
	}
	return downloadUrl
}

func QueryAdoptiumRelease() []string {
	var release []string
	var url = "https://api.adoptium.net/v3/info/available_releases"
	body, err := web.Call(url)
	if err != nil {
		return release
	}
	var unmarshalledJson interface{}
	e := json.Unmarshal(body, &unmarshalledJson)
	if e != nil {
		return release
	}

	queryRunner, err := gojq.Parse(".available_releases")
	if err != nil {
		return release
	}
	iter := queryRunner.Run(unmarshalledJson)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if array, ok := v.([]interface{}); ok {
			for _, item := range array {
				r := strconv.FormatFloat(item.(float64), 'f', 0, 64)
				release = append(release, r)
			}
		}
	}
	return release
}

type AdoptiumQuery struct {
	OS   string `json:"os"`
	ARCH string `json:"arch"`
}
