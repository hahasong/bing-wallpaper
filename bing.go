package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
	"runtime"
	"time"
)

var host = "https://cn.bing.com"
var ver = "1.0"

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func main() {
	fmt.Printf("bing.com wallpaper downloader ver%s powered.by.hahasong\n", ver)
	time.Sleep(500 * time.Millisecond)

	resp, err := http.Get(host)
	if err != nil {
		fmt.Printf("get %s err: %+v\n", host, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("parse response err: %+v\n", err)
		return
	}

	re := regexp.MustCompile(`data-ultra-definition-src="(.*?)"`)
	substr := re.FindStringSubmatch(string(body))
	if len(substr) < 2 {
		fmt.Printf("search img url err\n")
		return
	}
	imgUrl := host + substr[1]

	re = regexp.MustCompile(`w=\d+&h=\d+`)
	imgUrl = re.ReplaceAllString(imgUrl, "w=&h=")
	fmt.Printf("img url: %s\nstart downloading..\n", imgUrl)
	url, err := url.Parse(imgUrl)
	if err != nil {
		fmt.Printf("get img name err: %+v\n", err)
		return
	}
	fileName := url.Query().Get("id")
	if fileName == "" {
		fmt.Printf("get img name empty\n")
		return
	}

	resp, err = http.Get(imgUrl)
	if err != nil {
		fmt.Printf("download img err: %+v\n", err)
		return
	}
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("download img err: %+v\n", err)
		return
	}
	path := path.Join(userHomeDir(), fileName)
	out, _ := os.Create(path)
	io.Copy(out, bytes.NewReader(body))

	fmt.Printf("[%s] download success\n", path)
}
