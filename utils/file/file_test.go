package file

import (
	"testing"
)

func TestUnzip(t *testing.T) {
	src := "D:\\__data\\jvmst\\download\\zulu10.3.5-ca-jdk10.0.2.zip"
	dest := "D:\\__data\\jvmst\\download\\zulu10.3.5-ca-jdk10.0.2_tmp"
	err := Unzip(src, dest)
	if err != nil {
		t.Fatal(err)
	}
}
