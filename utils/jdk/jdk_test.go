package jdk

import (
	"fmt"
	"testing"
)

func TestGetJavaHome_With_JRE(t *testing.T) {
	dest := "D:\\__data\\jvmst\\download\\zulu10.3.5-ca-jdk10.0.2_tmp"
	javaHome := GetJavaHome(dest)
	fmt.Println(javaHome)
}

func TestGetJavaHome(t *testing.T) {
	dest := "D:\\__data\\jvmst\\download\\zulu7.5.0.1-ca-jdk7.0.60_tmp"
	javaHome := GetJavaHome(dest)
	fmt.Println(javaHome)
}

func TestGetInstalled(t *testing.T) {
	dest := "D:\\__data\\jvmst\\store"
	installs := GetInstalled(dest)
	fmt.Println(installs)
}
