package jdk

import (
	"fmt"
	"testing"
)

func TestTrusteeshipJdkSource_JdkVersions(t *testing.T) {
	source := NewTrusteeshipJdkSource("")
	versions := source.JdkVersions()
	for _, value := range versions {
		fmt.Printf("%+v", value)
	}
}
