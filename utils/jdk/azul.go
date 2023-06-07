package jdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
)

type AzulJdkSource struct {
	vendor string
	url    string
}

func NewAzulJdkSource() *AzulJdkSource {
	return &AzulJdkSource{
		vendor: "Azul",
		url:    "https://api.azul.com/metadata/v1/zulu/packages",
	}
}

func (receiver *AzulJdkSource) SourceName() string {
	return receiver.vendor
}

func (receiver *AzulJdkSource) SourceUrl() string {
	return receiver.url
}

func (receiver *AzulJdkSource) JdkVersions() []JdkVersion {
	azulJdks := AzulJDKs()
	var versions []JdkVersion
	for _, value := range azulJdks {
		versions = append(versions, JdkVersion{Version: value.ShortName, Url: value.DownloadURL, Source: receiver.SourceName()})
	}
	return versions
}

func AzulJDKs() []AzulJDK {
	url := AzulApiEndpoint()
	body := call(url)
	var jdks []AzulJDK
	err := json.Unmarshal(body, &jdks)
	if err != nil {
		fmt.Printf("error %v \n", err)
	}
	for i := 0; i < len(jdks); i++ {
		lastIndex := strings.LastIndex(jdks[i].Name, "-")
		jdks[i].ShortName = jdks[i].Name[0:lastIndex]
	}
	return jdks
}

type AzulJDK struct {
	PackageUUID        string `json:"package_uuid"`
	Name               string `json:"name"`
	JavaVersion        []int  `json:"java_version"`
	OpenjdkBuildNumber int    `json:"openjdk_build_number"`
	Latest             bool   `json:"latest"`
	DownloadURL        string `json:"download_url"`
	Product            string `json:"product"`
	DistroVersion      []int  `json:"distro_version"`
	AvailabilityType   string `json:"availability_type"`
	ShortName          string
}

type AzulQuery struct {
	OS       string `json:"os"`
	ARCH     string `json:"arch"`
	Latest   bool   `json:"latest"`
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
}

func AzulDefaultQuery() {

}

func AzulApiEndpoint2(query AzulQuery) string {
	//https://api.azul.com/metadata/v1/docs/swagger
	var api = AzulApi() + "?os=$OS&arch=$ARCH&archive_type=zip&java_package_type=jdk&javafx_bundled=false&latest=true&release_status=ga&availability_types=CA&certifications=tck&page=1&page_size=100"
	api = strings.Replace(api, "$OS", query.OS, 1)
	api = strings.Replace(api, "$ARCH", query.ARCH, 1)

	return api
}

func AzulApiEndpoint() string {
	//https://api.azul.com/metadata/v1/docs/swagger
	var api = AzulApi() + "?os=$OS&arch=$ARCH&archive_type=zip&java_package_type=jdk&javafx_bundled=false&latest=true&release_status=ga&availability_types=CA&certifications=tck&page=1&page_size=100"
	api = strings.Replace(api, "$OS", runtime.GOOS, 1)
	api = strings.Replace(api, "$ARCH", runtime.GOARCH, 1)
	return api
}

func AzulApi() string {
	return "https://api.azul.com/metadata/v1/zulu/packages"
}

func call(url string) []byte {
	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	return body
}
