package jdk

import (
	"fmt"
	"github.com/ystyle/jvms/utils/file"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type JdkVersion struct {
	Version string `json:"version"`
	Url     string `json:"url"`
	Origin  string `json:"origin"`
}

type JdkOrigin interface {
	OriginName() string
	OriginDesc() string
	JdkVersions() []JdkVersion
}

var JdkOrigins = []JdkOrigin{
	NewAdoptiumJdkSource(),
	NewAzulJdkSource(),
	NewTrusteeshipJdkSource(""),
}

func GetInstalled(root string) []string {
	list := make([]string, 0)
	files, _ := ioutil.ReadDir(root)
	for i := len(files) - 1; i >= 0; i-- {
		list = append(list, files[i].Name())
		//Windows Symlink
		//if files[i].IsDir() || files[i].Mode() == fs.ModeSymlink {
		//}
	}
	return list
}

func IsVersionInstalled(root string, version string) bool {
	isInstalled := file.Exists(fmt.Sprintf("%s/%s/bin/javac.exe", root, version))
	return isInstalled
}

func RemoteJdkVersions() ([]JdkVersion, error) {
	var versions []JdkVersion
	for _, source := range JdkOrigins {
		versions = append(versions, source.JdkVersions()...)
	}
	return versions, nil
}

func GetJavaHome(jdkDir string) string {
	var javaHome string
	_ = fs.WalkDir(os.DirFS(jdkDir), ".", func(path string, d fs.DirEntry, err error) error {
		if filepath.Base(path) == "javac.exe" {
			temPath := strings.Replace(path, "bin/javac.exe", "", -1)
			javaHome = filepath.Join(jdkDir, temPath)
			return fs.SkipDir
		}
		return nil
	})
	return javaHome
}
