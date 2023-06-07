package jdk

import (
	"fmt"
	"testing"
)

func TestAdoptiumJdkSource_JdkVersions(t *testing.T) {
	source := NewAdoptiumJdkSource()
	for _, value := range source.JdkVersions() {
		fmt.Printf("%+v\n", value)
	}
}
