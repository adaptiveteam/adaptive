package models

import (
	"github.com/pkg/errors"
	"strconv"
	core "github.com/adaptiveteam/adaptive/core-utils-go"
	"strings"
	"net/url"
)

const (
	PathSeparator = "/"
	PathQuerySeparator = "?"
)
// Path is a list of strings - parts of path.
type Path []string

// ActionPath represents part of URL after host - path and query.
// It might be used for routing and passing data to Slack callbacks.
type ActionPath struct {
	Path Path
	Values url.Values
}

// RelActionPath is a relative path around some global prefix.
type RelActionPath struct {
	Path Path
	Values url.Values
}

// NewPath constructs path from sequence of strings
func NewPath(p ...string) Path {
	return Path(p)
}

// ParsePath constructs Path from string
func ParsePath(p string) Path {
	s := strings.Split(p, PathSeparator)
	if strings.HasPrefix(p, PathSeparator) {
		s = s[1:]
	}
	return Path(s)
}

// NewActionPath returns ActionPath constructed from path and
func NewActionPath(path Path, values ...Pair) ActionPath {
	return ActionPath {
		Path: path,
		Values: Values(values...),
	}
}

// ParseActionPath parses path and query - /path/?key=value
// NB! It's unsafe!!!
func ParseActionPath(a string) ActionPath {
	s := strings.SplitN(a, PathQuerySeparator, 2)
	p := ParsePath(s[0])
	switch len(s) {
	case 1:
		return NewActionPath(p)
	case 2:
		v, err := url.ParseQuery(s[1])
		core.ErrorHandler(err, "ActionPath", "ParseActionPath")
		return ActionPath{Path: p, Values: v}
	default:
		panic(errors.New("Couldn't parse ActionPath " + a))
	}
}

// Pair captures key and value
type Pair struct {
	Key string
	Value string
}

// P construct a pair from key&value
func P(key, value string) Pair {
	return Pair{
		Key: key,
		Value: value,
	}
}

// Values construct url.Values from sequence of pairs
func Values(values ...Pair) url.Values {
	v := make(url.Values)
	for _, p := range values {
		v.Set(p.Key, p.Value)
	}
	return v
}

// ValuesIsEmpty checks if url.Values is empty
func ValuesIsEmpty(v url.Values) bool {
	return len(map[string][]string(v)) == 0
}

// Encode converts path to string
func (p Path) Encode() string {
	return PathSeparator + strings.Join(p, PathSeparator)
}

// Encode converts AcitonPath to string
func (a ActionPath) Encode() string {
	p := a.Path.Encode()
	if ValuesIsEmpty(a.Values) {
		return p
	}
	return p + PathQuerySeparator + a.Values.Encode()
}

// Head returns the head of path
func (p Path) Head() string {
	if len(p) == 0 { return "" }
	return []string(p)[0]
}

// Tail returns the tail of path
func (p Path) Tail() Path {
	if len(p) == 0 { return Path{} }
	return []string(p)[1:]
}

// Append adds a few more parts at the end of path
func (p Path)Append(parts ... string) Path {
	return append([]string(p), parts...)
}

// Prepend adds a new part in the beginning of the path.
func (p Path)Prepend(part string) Path {
	return append([]string{part}, []string(p)...)
}

// Tail returns action path with stripped out head
func (a ActionPath) Tail() RelActionPath {
	return RelActionPath{Path: a.Path.Tail(), Values: a.Values}
}

// HeadTail divides action path into head (for routing) and rest of action path
func (a ActionPath) HeadTail() (string, RelActionPath) {
	return a.Path.Head(), a.Tail()
}

// ToRelActionPath converts to RelActionPath
func (a ActionPath) ToRelActionPath() RelActionPath {
	return RelActionPath{Path: a.Path, Values: a.Values}
}
// Param gets parameter value from Values
func (a ActionPath) Param(key string) string {
	return a.Values.Get(key)
}

// SetParam overrides one of the existing Values
func (a ActionPath) SetParam(key string, value string) (res ActionPath) {
	res = a
	res.Values.Set(key, value)
	return 
}

// CurrentQuarterHash mixes year and month into a single string for hashing purposes.
func CurrentQuarterHash() string {
	year, month := core.CurrentYearMonth()
	hash := strconv.Itoa(year) + "-" + strconv.Itoa(int(month))
	return hash
}

// CurrentQuarterHashPair constructs `quarter_hash=2019-06`
func CurrentQuarterHashPair() Pair {
	return P("quarter_hash", CurrentQuarterHash())
}

// Prepend prepends a single step to this relative action path.
func (r RelActionPath)Prepend(part string) RelActionPath {
	return RelActionPath{
		Path: r.Path.Prepend(part),
		Values: r.Values,
	}
}

// ToActionPath converts to a global ActionPath
func (r RelActionPath)ToActionPath(prefix Path) ActionPath {
	return ActionPath{
		Path: prefix.Append(r.Path...),
		Values: r.Values,
	}
}

// Tail returns action path with stripped out head
func (r RelActionPath) Tail() RelActionPath {
	return RelActionPath{Path: r.Path.Tail(), Values: r.Values}
}

// HeadTail divides action path into head (for routing) and rest of action path
func (r RelActionPath) HeadTail() (string, RelActionPath) {
	return r.Path.Head(), r.Tail()
}

// Param gets parameter value from Values
func (r RelActionPath) Param(key string) string {
	return r.Values.Get(key)
}

// SetParam overrides one of the existing Values
func (r RelActionPath) SetParam(key string, value string) (res RelActionPath) {
	res = r
	res.Values.Set(key, value)
	return 
}

// Encode converts RelActionPath to string
func (r RelActionPath) Encode() (res string) {
	p := strings.Join(r.Path, PathSeparator)
	if ValuesIsEmpty(r.Values) {
		res = p
	} else {
		res = p + PathQuerySeparator + r.Values.Encode()
	}
	return
}
