package main

import (
	"fmt"
	"bytes"
	"github.com/slack-go/slack"
)

func main() {
	token := "xoxp-4..."
	api := slack.New(token)

	buf := bytes.NewBufferString("hello")
	params := slack.FileUploadParameters{
		Title:           "Simple file",
		Filename:        "file.txt",
		Reader:          buf,
		Channels:        []string{"UJ0SX0G9X"},
		ThreadTimestamp: "",
	}
	var slackFile *slack.File
	var err error
	slackFile, err = api.UploadFile(params)

	if err == nil {
		fmt.Println(slackFile)
	} else {
		fmt.Printf("Error:%+v", err)
	}
}
