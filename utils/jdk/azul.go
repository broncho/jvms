package jdk

import (
	"encoding/json"
	"fmt"
	"github.com/ystyle/jvms/utils/web"
	"runtime"
	"strings"
)

type AzulJdkSource struct {
	vendor     string
	vendorHome string
	apiDoc     string
}

func NewAzulJdkSource() *AzulJdkSource {
	return &AzulJdkSource{
		vendor:     "Azul",
		vendorHome: "https://www.azul.com/",
		apiDoc:     "https://api.azul.com/metadata/v1/docs/swagger",
	}
}

func (receiver *AzulJdkSource) OriginName() string {
	return receiver.vendor
}

func (receiver *AzulJdkSource) OriginDesc() string {
	return fmt.Sprintf("%s (%s)", receiver.vendorHome, receiver.apiDoc)
}

func (receiver *AzulJdkSource) JdkVersions() []JdkVersion {
	var versions []JdkVersion
	query := AzulQuery{
		OS:       runtime.GOOS,
		ARCH:     runtime.GOARCH,
		Latest:   true,
		Page:     1,
		PageSize: 100,
	}
	azulJdks, err := QueryAzulJdkVersions(query)
	if err != nil {
		return versions
	}
	for _, value := range azulJdks {
		versions = append(versions, JdkVersion{Version: value.ShortName, Url: value.DownloadURL, Origin: receiver.OriginName()})
	}
	return versions
}

func QueryAzulJdkVersions(query AzulQuery) ([]AzulJDK, error) {
	url := buildAzulQueryUrl(query)
	body, err := web.Call(url)
	if err != nil {
		return nil, err
	}
	var jdks []AzulJDK
	err = json.Unmarshal(body, &jdks)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(jdks); i++ {
		lastIndex := strings.LastIndex(jdks[i].Name, "-")
		jdks[i].ShortName = jdks[i].Name[0:lastIndex]
	}
	return jdks, nil
}

func buildAzulQueryUrl(query AzulQuery) string {
	//https://api.azul.com/metadata/v1/docs/swagger
	var queryParams = map[string]interface{}{
		"archive_type":       "zip",
		"java_package_type":  "jdk",
		"javafx_bundled":     false,
		"latest":             query.Latest,
		"release_status":     "ga",
		"availability_types": "CA",
		"certifications":     "tck",
		"page":               1,
		"page_size":          100,
		"os":                 query.OS,
		"arch":               query.ARCH,
	}
	var queryBuilder strings.Builder
	queryBuilder.WriteString("https://api.azul.com/metadata/v1/zulu/packages")
	queryBuilder.WriteString("?")
	for k, v := range queryParams {
		queryBuilder.WriteString(k)
		queryBuilder.WriteString("=")
		queryBuilder.WriteString(fmt.Sprintf("%v", v))
		queryBuilder.WriteString("&")
	}
	return strings.TrimRight(queryBuilder.String(), "&")
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
