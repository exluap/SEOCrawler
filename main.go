/**
 * @project SEOCrawler
 * @date 02.04.2018 19:29
 * @author Nikita Zaytsev (exluap) <nickzaytsew@gmail.com>
 * @twitter https://twitter.com/exluap
 * @keybase https://keybase.io/exluap
 */

package main

import (
	"bufio"
	"fmt"
	"github.com/blang/semver"
	"github.com/exluap/SEOCrawler/utils"
	"github.com/getsentry/raven-go"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"log"
	"os"
	"runtime"
	"strings"
)

var version = "1.0.1"

var DSN = ""

func init() {
	raven.SetDSN(DSN)
}

func main() {

	confirmAndSelfUpdate()

	crawl()

}

func confirmAndSelfUpdate() {
	latest, found, err := selfupdate.DetectLatest("exluap/SEOCrawler")
	if err != nil {
		log.Println("Error occurred while detecting version:", err)
		return
	}

	v := semver.MustParse(version)
	if !found || latest.Version.Equals(v) {
		log.Println("Current version is the latest")
		return
	}

	fmt.Print("Detected new version: ", latest.Version)
	fmt.Print("Do you want to update to ", latest.Version, "? (y/n): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input != "y\n" && input != "n\n") || (input != "y\r" && input != "y\r") {
		log.Println("Invalid input")
		return
	}
	if input == "n\n" || input == "n\r" {
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, os.Args[0]); err != nil {
		log.Println("Error occurred while updating binary:", err)
		raven.CaptureErrorAndWait(err, nil)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
}

func crawl() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Write URL (http://example.com/): ")
	url, _ := reader.ReadString('\n')
	if runtime.GOOS == "windows" {
		url = strings.Replace(url, "\r\n", "", -1)
	} else {
		url = strings.Replace(url, "\n", "", -1)
	}

	utils.StartHere(url)
}
