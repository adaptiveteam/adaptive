package pagination

// InterfacePager is a function that returns pages of elements
type InterfacePager func() (InterfaceSlice, InterfacePager, error)

// Page returns a single page and continuation
func (p InterfacePager)Page() (InterfaceSlice, InterfacePager, error) {
	return p()
}
// Drain collects all data from pager to a single slice
func (p InterfacePager)Drain() (sl InterfaceSlice, err error) {
	var ip InterfacePager
	sl, ip, err = p.Page()
	if err == nil && len(sl) > 0 {
		var sl2 InterfaceSlice
		sl2, err = ip.Drain()
		if err == nil {
			sl = append(sl, sl2...)
		}
	}
	return
}

// IsEmpty fetches the first page and makes sure it's empty
func (p InterfacePager)IsEmpty() (fl bool, err error) {
	var sl InterfaceSlice
	sl, _, err = p.Page()
	fl = len(sl) == 0
	return
}
// NonEmpty fetches the first page and makes sure it's non empty
func (p InterfacePager)NonEmpty() (fl bool, err error) {
	fl, err = p.IsEmpty()
	fl = !fl
	return
}
// InterfacePagerPure constructs pager that will emit the given values
func InterfacePagerPure(slice ... interface{}) (ip InterfacePager) {
	return InterfacePagerFromSlice(slice)
}

// InterfacePagerFromSlice constructs pager that will emit the given values
func InterfacePagerFromSlice(slice []interface{}) (ip InterfacePager) {
	if len(slice) > 0 {
		ip = func () (InterfaceSlice, InterfacePager, error) {
			return slice, InterfacePagerPure(), nil
		}
	} else {
		ip = emptyPage
	}
	return
}

// InterfacePagerConcat concatenates a few pagers
func InterfacePagerConcat(pagers ... InterfacePager) (res InterfacePager) {
	if len(pagers) > 0 {
		head := pagers[0]
		tail := pagers[1:]
		res = func() (InterfaceSlice, InterfacePager, error) {
			sl,p,err2 := head()

			if err2 == nil {
				if len(sl) > 0 {
					pagers = append([]InterfacePager{p}, tail...)
					return sl, InterfacePagerConcat(pagers...), err2	
				} else {
					return InterfacePagerConcat(tail...)()
				}
			} else {
				return sl,p,err2
			}
		}
	} else {
		res = emptyPage
	}
	return
}

func emptyPage() (sl InterfaceSlice, ip InterfacePager, err error) {
	ip = emptyPage
	return
}
// Limit creates a pager that will return at most limit elements.
func (p InterfacePager)Limit(limit int) InterfacePager {
	return func() (sl InterfaceSlice, ip InterfacePager, err error) {
		if limit > 0 {
			var slice InterfaceSlice
			slice, ip, err = p()
			if err == nil {
				if len(slice) > limit {
					sl = slice[0:limit]
					ip = emptyPage // ignoring the rest
				} else {
					sl = slice
					ip = ip.Limit(limit - len(slice))
				}
			}
		} else {
			ip = emptyPage
		}
		return
	}
}
// Map converts elements using f
func (p InterfacePager)Map(f func (interface{}) interface{}) InterfacePager {
	return func() (sl InterfaceSlice, ip InterfacePager, err error) {
		var slice InterfaceSlice
		slice, ip, err = p()
		if err == nil {
			sl = make(InterfaceSlice, 0, len(slice))
			for _, i := range slice {
				sl = append(sl, f(i))
			}
		}
		return
	}
}
// MapE converts elements using f. Stops in case of errors
func (p InterfacePager)MapE(f func (interface{}) (interface{}, error)) InterfacePager {
	return func() (sl InterfaceSlice, ip InterfacePager, err error) {
		var slice InterfaceSlice
		slice, ip, err = p()
		if err == nil {
			for _, i := range slice {
				var el interface {}
				el, err = f(i)
				if err != nil {
					return
				}
				sl = append(sl, el)
			}
		}
		return
	}
}
// FilterE leaves only elements for which `pred` is true
func (p InterfacePager)FilterE(pred func(interface{}) (bool, error)) InterfacePager {
	return func() (sl InterfaceSlice, ip InterfacePager, err error) {
		var slice InterfaceSlice
		slice, ip, err = p()
		if err == nil && len(slice) > 0 {
			for _, i := range slice {
				var fl bool
				fl, err = pred(i)
				if err != nil {
					return
				}
				if fl {
					sl = append(sl, i)
				}
			}
			if len(sl) == 0 {
				return ip.FilterE(pred)()
			}
		}
		return
	}
}
// FlatMap constructs a pager out of inner pagers
func (p InterfacePager)FlatMap(f func (interface{}) InterfacePager) InterfacePager {
	pp := p.Map(func (i interface{})interface{} { return f(i)})
	return func() (sl InterfaceSlice, ip InterfacePager, err error) {
		sl, ip, err = pp()
		if err == nil && len(sl) > 0 {
			pagers := sl.AsInterfacePagerSlice()
			p0 := InterfacePagerConcat(pagers...)
			pager := InterfacePagerConcat(p0, ip)
			sl, ip, err = pager()
		}
		return
	}
}
