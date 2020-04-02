package lambda

import (
	"fmt"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)


var (
	// NoReportTemplate is a message about absent report
	NoReportTemplate ui.RichText = "_Sorry, unable to locate the report :disappointed:. It could be that there isn't any feedback yet._"
)

// TitleTemplate template for report title
func TitleTemplate(userID string) ui.PlainText {
	return ui.PlainText(fmt.Sprintf("Performance Report for <@%s>", userID))
}
