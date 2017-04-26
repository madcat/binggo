package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type meta struct {
	StartDate string `json:"startdate"`
	URL       string `json:"url"`
	CopyRight string `json:"copyright"`
}

type result struct {
	Images []meta `json:"images"`
}

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

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func download(filepath string, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func title(filepath string, label string) error {
	str := "gm convert -font /Users/lingfei/Library/Fonts/msyh.ttc -fill white -pointsize 24 -gravity South -draw"
	parts := strings.Fields(str)
	head := parts[0]
	parts = parts[1:]
	parts = append(parts, fmt.Sprintf("text 10,20 '%s'", label))
	parts = append(parts, filepath)
	parts = append(parts, filepath)
	out, err := exec.Command(head, parts...).CombinedOutput()
	if err != nil {
		log.Println(string(out))
		return err
	}
	return nil
}

func main() {
	var flagDir = flag.String("dir", userHomeDir(), "directory downloaded images are saved")
	var flagCopy = flag.Bool("copy", false, "whether to print copyright info at bottom. (requires graphicsmagick)")

	flag.Parse()

	resp, err := http.Get("http://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=1&mkt=en-US")
	checkError(err)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	checkError(err)
	var r result
	err = json.Unmarshal(body, &r)
	checkError(err)

	img := r.Images[0]
	filename := fmt.Sprintf("%s/%s.jpg", *flagDir, img.StartDate)

	imgURL := fmt.Sprintf("%s%s", "http://www.bing.com/", img.URL)
	checkError(download(filename, imgURL))
	if *flagCopy {
		checkError(title(filename, img.CopyRight))
	}
}
