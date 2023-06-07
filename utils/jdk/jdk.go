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

//func getJdkVersions() ([]JdkVersion, error) {
//	jsonContent, err := web.GetRemoteTextFile(config.Originalpath)
//	if err != nil {
//		return nil, err
//	}
//	var versions []JdkVersion
//	err = json.Unmarshal([]byte(jsonContent), &versions)
//	if err != nil {
//		return nil, err
//	}
//	//fmt.Println(versions)
//	adoptiumJdks := strings.Split(adoptium_jdk_go.ApiListReleases(), "\n")
//	for _, adoptiumJdkUrl := range adoptiumJdks {
//		fileSeparatorIndex := strings.LastIndex(adoptiumJdkUrl, "/")
//		fileName := adoptiumJdkUrl[fileSeparatorIndex+1:]
//		fileVersion := strings.TrimSuffix(fileName, ".zip")
//		//fmt.Println(fileVersion)
//		versions = append(versions, JdkVersion{Version: fileVersion, Url: adoptiumJdkUrl})
//	}
//
//	//Azul JDKs
//	azulJdks := jdk.AzulJDKs()
//	for _, azulJdk := range azulJdks {
//		versions = append(versions, JdkVersion{Version: azulJdk.ShortName, Url: azulJdk.DownloadURL})
//	}
//
//	//fmt.Println(versions)
//	return versions, nil
//}
