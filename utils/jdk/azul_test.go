package jdk

import (
	"fmt"
	"testing"
)

func TestAzulJdkSource_JdkVersions(t *testing.T) {
	source := NewAzulJdkSource()
	for _, value := range source.JdkVersions() {
		fmt.Printf("%+v\n", value)
	}
}

func TestQueryAzulJdkVersions(t *testing.T) {
	query := AzulQuery{PageSize: 10, Page: 1, OS: "linux", ARCH: "amd64", Latest: true}
	versions, err := QueryAzulJdkVersions(query)
	if err != nil {
		t.Fail()
	}
	for _, value := range versions {
		fmt.Printf("%+v\n", value)
	}
}
