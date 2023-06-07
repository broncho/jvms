package jdk

import (
	"fmt"
	"testing"
)

func TestAzulJDKs(t *testing.T) {
	jdks := AzulJDKs()
	for _, value := range jdks {
		fmt.Printf("%+v \n", value)
	}
}

func TestAzulJdkSource_JdkVersions(t *testing.T) {
	source := NewAzulJdkSource()
	for _, value := range source.JdkVersions() {
		fmt.Printf("%+v\n", value)
	}
}

func TestAzulApiEndpoint(t *testing.T) {
	url := AzulApiEndpoint()
	fmt.Println(url)
}
