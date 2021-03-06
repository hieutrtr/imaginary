package main

// CephQueryFormat format query string by pool and oid
const CephQueryFormat string = "%s&cpool=%s&coid=%s"

// QueryMap mapping query with operation
type QueryMap map[string]string
type OperationMap map[string]Operation

// ServiceQueryMap mapping service name and ServiceQuery
var ServiceQueryMap = make(map[string]ServiceQuery)

// ServiceQuery interface getting query by service's operation
type ServiceQuery interface {
	getQuery(string, string) string
	getOperation(string) Operation
}

// ServiceQueryRegister register service for handling query
func ServiceQueryRegister(t string, q ServiceQuery) {
	ServiceQueryMap[t] = q
}
