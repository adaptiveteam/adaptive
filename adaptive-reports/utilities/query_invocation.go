package utilities

// QueryInvocation is a structure that captures invocation of a query
// with certain parameters.
type QueryInvocation struct {
	Name string
	SQL string
	ArgumentValues []interface{}
}

// InvokeQuery constructs QueryInvocation
func InvokeQuery(name string, sql string, args ... string) QueryInvocation {
	arguments := make([]interface{}, 0, len(args))
	for _, a := range args {
		arguments = append(arguments, a)
	}

	return QueryInvocation{
		Name: name,
		SQL: sql,
		ArgumentValues: arguments,
	}
}
// QueryInvocations is a slice of QueryInvocation s
type QueryInvocations []QueryInvocation

// InvokeQueries constructs a sequence of query invocations
func InvokeQueries(q ... QueryInvocation) QueryInvocations {
	return q
}