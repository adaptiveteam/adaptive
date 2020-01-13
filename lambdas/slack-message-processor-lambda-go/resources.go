package lambda

const (
	HolidaysMenuKey                     = "holidays"
	HolidaysCreateNewHolidayKey         = "holidays-create-new"
	HolidaysSimpleListKey               = "holidays-simple-list"
	HolidaysListHolidaysKey             = "holidays-list"
	AdaptiveValuesMenuKey               = "adaptive-values"
	AdaptiveValuesCreateNewHolidayKey   = "adaptive-values-create-new"
	AdaptiveValuesSimpleListKey         = "adaptive-values-simple-list"
	AdaptiveValuesListAdaptiveValuesKey = "adaptive-values-list"
)

const (
	holidaysMenuTemplate                           = "Holidays"
	holidaysCreateNewHolidayTemplate               = "Create a new holiday"
	holidaysSimpleListHolidaysTemplate             = "List holidays"
	holidaysListHolidaysTemplate                   = "Update or delete existing holidays" // Includes delete a holiday
	adaptiveValuesMenuTemplate                     = "Adaptive values"
	adaptiveValuesCreateNewHolidayTemplate         = "Create a new value"
	adaptiveValuesSimpleListAdaptiveValuesTemplate = "List adaptive values"
	adaptiveValuesListAdaptiveValuesTemplate       = "Update or delete existing adaptive values" // Includes delete a holiday
)

var (
	templates = map[string]string{
		HolidaysMenuKey:                     holidaysMenuTemplate,
		HolidaysCreateNewHolidayKey:         holidaysCreateNewHolidayTemplate,
		HolidaysSimpleListKey:               holidaysSimpleListHolidaysTemplate,
		HolidaysListHolidaysKey:             holidaysListHolidaysTemplate,
		AdaptiveValuesMenuKey:               adaptiveValuesMenuTemplate,
		AdaptiveValuesCreateNewHolidayKey:   adaptiveValuesCreateNewHolidayTemplate,
		AdaptiveValuesSimpleListKey:         adaptiveValuesSimpleListAdaptiveValuesTemplate,
		AdaptiveValuesListAdaptiveValuesKey: adaptiveValuesListAdaptiveValuesTemplate,
	}
)

// RetrieveTemplate returns a dialog message template for
// a given context and subject
func RetrieveTemplate(key string) string {
	return templates[key]
}
