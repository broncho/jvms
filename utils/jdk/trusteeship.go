package jdk

import (
	"encoding/json"
	"github.com/ystyle/jvms/utils/web"
)

const DefaultOriginalPath = "https://raw.githubusercontent.com/ystyle/jvms/new/jdkdlindex.json"

type TrusteeshipJdkSource struct {
	vendor    string
	originUrl string
}

func NewTrusteeshipJdkSource(origin string) *TrusteeshipJdkSource {
	var url = DefaultOriginalPath
	if origin != "" {
		url = origin
	}
	return &TrusteeshipJdkSource{
		vendor:    "Github",
		originUrl: url,
	}
}

func (receiver *TrusteeshipJdkSource) OriginName() string {
	return receiver.vendor
}

func (receiver *TrusteeshipJdkSource) OriginDesc() string {
	return receiver.originUrl
}

func (receiver *TrusteeshipJdkSource) JdkVersions() []JdkVersion {
	body, err := web.Call(receiver.originUrl)
	var jdks []JdkVersion
	err = json.Unmarshal(body, &jdks)
	if err != nil {
		return jdks
	}
	for i := 0; i < len(jdks); i++ {
		jdks[i].Url = receiver.OriginName()
	}
	return jdks
}
