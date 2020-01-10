package utilities

import excel "github.com/360EntSecGroup-Skylar/excelize/v2"

func CreateDocumentProperties(
	category string,
	description string,
	keywords []string,
	subject string,
	title string,
	) (rv *excel.DocProperties) {

	keywordsString := "Adaptive"
	if len(keywords) > 1 {
		keywordsString = keywordsString+", "
		for i, v := range keywords {
			keywordsString = keywordsString + v
			if i < len(keywords) -1 {
				keywordsString= keywordsString+", "
			}
		}
	}

	rv = &excel.DocProperties{
		Category:       category,
		ContentStatus:  "Generated",
		Creator:        "Adaptive",
		Description:    description,
		Identifier:     "xlsx",
		Keywords:       keywordsString,
		LastModifiedBy: "Adaptive",
		Revision:       "0",
		Subject:        subject,
		Title:          title,
		Language:       "en-US",
		Version:        "1.0.0",
	}

	return rv
}