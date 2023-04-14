package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UpdateAPI = "https://api.github.com/repos/jinyu121/ImageServer/releases/latest"

func CheckUpdate(current string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Failed to check update")
		}
	}()

	if gin.ReleaseMode != gin.Mode() {
		return
	}
	resp, err := http.Get(UpdateAPI)
	if nil != err {
		return
	}
	defer func() { _ = resp.Body.Close() }()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if nil != err {
		return
	}

	if result["prerelease"].(bool) {
		return
	}
	tagName := result["tag_name"].(string)
	releaseURL := result["html_url"].(string)
	if tagName != current {
		fmt.Printf("New version %s is available, please update: %s\n", tagName, releaseURL)
	}
}
