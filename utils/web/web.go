package web

import (
	"errors"
	"fmt"
	"github.com/ystyle/jvms/utils/file"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	pb "gopkg.in/cheggaaa/pb.v1"
)

var client = &http.Client{}

func SetProxy(p string) {
	if p != "" && p != "none" {
		proxyUrl, _ := url.Parse(p)
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	} else {
		client = &http.Client{}
	}
}

func Download(url string, target string) bool {
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false
	}
	if response.StatusCode != 200 {
		fmt.Println("Error status while downloading", url, "-", response.StatusCode)
		return false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(response.Body)

	output, err := os.Create(target)
	if err != nil {
		fmt.Println("Error while creating", target, "-", err)
		return false
	}
	defer func(output *os.File) {
		err := output.Close()
		if err != nil {

		}
	}(output)

	// 创建一个进度条
	bar := pb.New(int(response.ContentLength)).SetUnits(pb.U_BYTES_DEC).SetRefreshRate(time.Millisecond * 10)
	// 显示下载速度
	bar.ShowSpeed = true

	// 显示剩余时间
	bar.ShowTimeLeft = true

	// 显示完成时间
	bar.ShowFinalTime = true

	bar.SetWidth(80)

	bar.Start()
	writer := io.MultiWriter(output, bar)
	_, err = io.Copy(writer, response.Body)
	if err != nil {
		fmt.Println("Error while downloading", url, "-", err)
		return false
	}
	bar.Finish()

	return true
}

func GetJDK(download string, version string, url string) (string, bool) {
	fileName := filepath.Join(download, fmt.Sprintf("%s.zip", version))
	if file.Exists(fileName) {
		fmt.Printf("JDK already download %v", fileName)
		return fileName, true
	}
	//os.Remove(fileName)
	if url == "" {
		//No url should mean this version/arch isn't available
		fmt.Printf("JDK %s isn't available right now.", version)
	} else {
		fmt.Printf("Downloading jdk version %s...\n", version)
		if Download(url, fileName) {
			fmt.Println("Complete")
			return fileName, true
		} else {
			return "", false
		}
	}
	return "", false
}

func Call(url string) ([]byte, error) {
	res, err := client.Get(url)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Call %s error %v", url, err))
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Call %s read response error %v", url, err))
	}
	return body, nil
}
