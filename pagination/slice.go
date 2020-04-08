package pagination

// InterfaceSlice is a generic slice
type InterfaceSlice []interface{}

// AsStringSlice casts each element to string
func (is InterfaceSlice)AsStringSlice() (res []string) {
	for _, i := range is {
		res = append(res, i.(string))
	}
	return
}

// AsIntSlice casts each element to int
func (is InterfaceSlice)AsIntSlice() (res []int) {
	for _, i := range is {
		res = append(res, i.(int))
	}
	return
}

// AsInterfacePagerSlice casts each element to InterfacePager
func (is InterfaceSlice)AsInterfacePagerSlice() (res []InterfacePager) {
	for _, i := range is {
		res = append(res, i.(InterfacePager))
	}
	return
}
// AsInterfaceSlice casts each element to interface{}
func (is InterfaceSlice)AsInterfaceSlice() (res []interface{}) {
	return is
}

// StringsToInterfaceSlice - []string to []interface{}
func StringsToInterfaceSlice(strs []string) (res InterfaceSlice) {
	for _, s := range strs {
		res = append(res, s)
	}
	return
}
