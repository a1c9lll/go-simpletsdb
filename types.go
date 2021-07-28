package simpletsdb

import "net/http"

type SimpleTSDB struct {
	// The host used to connect to SimpleTSDB
	Host string
	// The port used to connect to SimpleTSDB
	Port   int
	client *http.Client
}

type InsertPointRequest struct {
	// The metric name of the inserted point
	Metric string
	// The tags associated with the inserted point
	Tags map[string]string
	// The point to insert
	Point *Point
}

type AggregatorQuery struct {
	// Name of aggregator
	Name string `json:"name"`
	// Options for aggregator.
	Options map[string]interface{} `json:"options"`
}

type QueryPointsRequest struct {
	// The metric to query.
	Metric string `json:"metric"`
	// The timestamp in nanoseconds for the start of the query.
	Start int64 `json:"start"`
	// optional: The timestamp in nanoseconds for the end of the query.
	End int64 `json:"end"`
	// optional: The max number of points to return.
	N int64 `json:"n"`
	// optional: Key/value pairs to add criteria to the query.
	Tags map[string]string `json:"tags"`
	// optional: Windowing options.
	Window map[string]interface{} `json:"window"`
	// optional: An array of aggregators.
	Aggregators []*AggregatorQuery `json:"aggregators"`
}

type DeletePointsRequest struct {
	// The metric to query.
	Metric string `json:"metric"`
	// The timestamp in nanoseconds for the start of the query.
	Start int64 `json:"start"`
	// The timestamp in nanoseconds for the end of the query.
	End int64 `json:"end"`
	// optional: Key/value pairs to add criteria to the query.
	Tags map[string]string `json:"tags"`
}

type Point struct {
	Value     interface{} `json:"value"`
	Timestamp int64       `json:"timestamp"`
	// If `window` was set in the query this will be set (unless a windowed aggregator was also applied)
	Window int64 `json:"window,omitempty"`
}

type pointReader struct {
	reqs []*InsertPointRequest
	i    int
}
