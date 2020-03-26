package common

import (
	"github.com/adaptiveteam/adaptive/pagination"
)
// InterfaceEffectfulPredicate is a predicate that requires connection
type InterfaceEffectfulPredicate = func(conn DynamoDBConnection) func(i interface{}) (bool, error)

// InterfaceStream is a function that given connection returns a pager that will be fetching data pages.
// This type supports DSL for creating streams with the desired contents.
type InterfaceStream func(conn DynamoDBConnection) pagination.InterfacePager

// AsInterfaceStreamSlice casts each element to InterfaceStream
func AsInterfaceStreamSlice(is pagination.InterfaceSlice) (res []InterfaceStream) {
	for _, i := range is {
		res = append(res, i.(InterfaceStream))
	}
	return
}

// InterfaceStreamPure is an InterfaceStream constructor that uses the given sequence of elements.
func InterfaceStreamPure(slice ... interface{}) InterfaceStream {
	return func(conn DynamoDBConnection) pagination.InterfacePager {
		return pagination.InterfacePagerPure(conn)
	}
}

// InterfaceStreamFromPager constructs a stream given pager
func InterfaceStreamFromPager(pager pagination.InterfacePager) InterfaceStream {
	return func(conn DynamoDBConnection) pagination.InterfacePager {
		return pager
	}
}

// InterfaceStreamConcat concatenates a few streams
func InterfaceStreamConcat(streams ... InterfaceStream) (res InterfaceStream) {
	if len(streams) > 0 {
		head := streams[0]
		tail := streams[1:]
		res = func(conn DynamoDBConnection) pagination.InterfacePager {
			return pagination.InterfacePagerConcat(
				head(conn), 
				InterfaceStreamConcat(tail...)(conn),
			)
		}
	} else {
		res = InterfaceStreamPure()
	}
	return
}
// Run starts the stream and returns the pager.
func (i InterfaceStream)Run(conn DynamoDBConnection) pagination.InterfacePager {
	return i(conn)
}

// Limit creates a stream that will return at most limit elements.
func (i InterfaceStream)Limit(limit int) InterfaceStream {
	return func(conn DynamoDBConnection) pagination.InterfacePager {
		return i.Run(conn).Limit(limit)
	}
}

// Map converts each element using the provided conversion function
func (i InterfaceStream)Map(f func (interface{}) interface{}) InterfaceStream {
	return func(conn DynamoDBConnection) pagination.InterfacePager {
		return i.Run(conn).Map(f)
	}
}

// FilterF returns only such elements that pred is true
func (i InterfaceStream)FilterF(pred InterfaceEffectfulPredicate) InterfaceStream {
	return func(conn DynamoDBConnection) pagination.InterfacePager {
		return i.Run(conn).FilterE(pred(conn))
	}
}

// FlatMapF constructs a stream out of inner streams
func (i InterfaceStream)FlatMapF(f func (interface{}) InterfaceStream) InterfaceStream {
	return func(conn DynamoDBConnection) pagination.InterfacePager {
		return i.Run(conn).FlatMap(
			func (i interface{})pagination.InterfacePager { 
				return f(i).Run(conn)
			},
		)
	}
}
