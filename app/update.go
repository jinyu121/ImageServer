package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

var UpdateAPI = "https://api.github.com/repos/jinyu121/ImageServer/releases/latest"

func CheckUpdate() {
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
	if tagName != Version {
		fmt.Printf("New version %s is available, please update.\n%s\n", tagName, releaseURL)
	}
}
