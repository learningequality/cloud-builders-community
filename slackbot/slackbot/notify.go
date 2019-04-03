package slackbot

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	cloudbuild "google.golang.org/api/cloudbuild/v1"
)

// Notify posts a notification to Slack that the build is complete.
func Notify(b *cloudbuild.Build, webhook string) {
	url := fmt.Sprintf("https://console.cloud.google.com/cloud-build/builds/%s", b.Id)
	server := fmt.Sprintf("https://%s.studio.cd.learningequality.org", b.Substitutions["_RELEASE_NAME"])
	var i string
	switch b.Status {
	case "SUCCESS":
		i = fmt.Sprintf(":white_check_mark: %s\n<%s> is going to be turned down.", b.Status, server)
	case "FAILURE", "CANCELLED":
		i = fmt.Sprintf(":x: %s\nPlease try again later.", b.Status)
	case "STATUS_UNKNOWN", "INTERNAL_ERROR":
		i = fmt.Sprintf(":interrobang: %s", b.Status)
	default:
		i = fmt.Sprintf(":question: %s", b.Status)
	}
	j := fmt.Sprintf(
		`{"text": "Cloud Build %s complete: %s",
		    "attachments": [
				{
					"fallback": "Open build details at %s",
					"actions": [
						{
							"type": "button",
							"text": "Open details",
							"url": "%s"
						}
					]
				}
			]}`, b.Id, i, url, url)

	r := strings.NewReader(j)
	resp, err := http.Post(webhook, "application/json", r)
	if err != nil {
		log.Fatalf("Failed to post to Slack: %v", err)
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Posted message to Slack: [%v], got response [%s]", j, body)
}
