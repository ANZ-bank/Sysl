package syslutil

import (
	"sort"

	sysl "github.com/anz-bank/sysl/pkg/sysl"
)

type endpointSourceOrder []*sysl.Endpoint

// Len is the number of elements in the collection.
func (t endpointSourceOrder) Len() int {
	return len(t)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (t endpointSourceOrder) Less(i, j int) bool {
	a := t[i].SourceContext
	b := t[j].SourceContext
	if a == nil {
		return b != nil
	}
	if b == nil {
		return false
	}
	return a.Start.Line < b.Start.Line
}

// Swap swaps the elements with indexes i and j.
func (t endpointSourceOrder) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

// EndpointsInSourceOrder an endpoint map's names in the order the endpoint
// appear in the source file.
func EndpointsInSourceOrder(endpoint []*sysl.Endpoint) []*sysl.Endpoint {
	tmp := make([]*sysl.Endpoint, len(endpoint))
	copy(tmp, endpoint)
	sort.Sort(endpointSourceOrder(tmp))
	return tmp
}

type endpointNamesSourceOrder struct {
	endpoint map[string]*sysl.Endpoint
	names    []string
}

// Len is the number of elements in the collection.
func (t *endpointNamesSourceOrder) Len() int {
	return len(t.names)
}

// Less reports whether the element with index i should sort before the element
// with index j.
func (t *endpointNamesSourceOrder) Less(i, j int) bool {
	a := t.endpoint[t.names[i]].SourceContext.Start.Line
	b := t.endpoint[t.names[j]].SourceContext.Start.Line
	return a < b
}

// Swap swaps the elements with indexes i and j.
func (t *endpointNamesSourceOrder) Swap(i, j int) {
	t.names[i], t.names[j] = t.names[j], t.names[i]
}

// EndpointNamesInSourceOrder an endpoint map's names in the order the endpoint
// appear in the source file.
func EndpointNamesInSourceOrder(endpoint map[string]*sysl.Endpoint) []string {
	names := make([]string, 0, len(endpoint))
	for name := range endpoint {
		names = append(names, name)
	}
	sort.Sort(&endpointNamesSourceOrder{endpoint, names})
	return names
}

// EndpointNamesInAlphabeticalOrder returns an endpoint map's names in
// alphabetical order.
func EndpointNamesInAlphabeticalOrder(endpoint map[string]*sysl.Endpoint) []string {
	names := make([]string, 0, len(endpoint))
	for name := range endpoint {
		names = append(names, name)
	}
	sort.Sort(&endpointNamesSourceOrder{endpoint, names})
	return names
}
