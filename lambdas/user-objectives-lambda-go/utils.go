package lambda


import (
	ebm "github.com/adaptiveteam/adaptive/engagement-builder/model"
	"github.com/adaptiveteam/adaptive/engagement-builder/ui"
)


func attachmentField(label ui.PlainText, value ui.PlainText) ebm.AttachmentField {
	return ebm.AttachmentField{
		Title: string(label), 
		Value: string(value),
	}
}
