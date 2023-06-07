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

func TestAzulApiEndpoint(t *testing.T) {
	url := AzulApiEndpoint()
	fmt.Println(url)
}
