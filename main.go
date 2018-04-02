/**
    * @project SEOCrawler
    * @date 02.04.2018 19:29
    * @author Nikita Zaytsev (exluap) <nickzaytsew@gmail.com>
    * @twitter https://twitter.com/exluap
    * @keybase https://keybase.io/exluap
*/

package main

import (
	"github.com/exluap/SEOCrawler/utils"
	"fmt"
	"os"
	"bufio"
	"strings"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"log"
	"github.com/blang/semver"
)

const version = "1.0.0"

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

	fmt.Print("Do you want to update to", latest.Version, "? (y/n): ")
	input, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil || (input != "y\n" && input != "n\n") {
		log.Println("Invalid input")
		return
	}
	if input == "n\n" {
		return
	}

	if err := selfupdate.UpdateTo(latest.AssetURL, os.Args[0]); err != nil {
		log.Println("Error occurred while updating binary:", err)
		return
	}
	log.Println("Successfully updated to version", latest.Version)
}

func crawl() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Write URL (http://example.com/): ")
	url, _ := reader.ReadString('\n')
	utils.StartHere(strings.TrimSuffix(url,"\n"))
}
