package lambda

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"regexp"
	"testing"
)

func TestUserRegex(t *testing.T) {
	var re = regexp.MustCompile(`(?m)<@([a-zA-Z0-9]+)> ([a-zA-Z\s\d]+)`)
	var str = `<@UEFF123> test report`

	fmt.Println(re.FindStringSubmatch(str))
}

func DisabledTestListUsers(t *testing.T) {
	api := slack.New("xoxb-1234")
	users, _ := api.GetUsers()
	for _, each := range users {
		byt, _ := json.Marshal(each)
		fmt.Println(string(byt))
	}
}
