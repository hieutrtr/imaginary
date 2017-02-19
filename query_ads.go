package main

import "fmt"

const queryTypeAds string = "ads"

// AdsQuery for ads service
type AdsQuery struct {
	queries QueryMap
	ops     OperationMap
}

func (q *AdsQuery) getQuery(ot string, id string) string {
	return fmt.Sprintf(CephQueryFormat, q.queries[ot], queryTypeAds, id)
}

func (q *AdsQuery) getOperation(ot string) Operation {
	return q.ops[ot]
}

func newAdsQuery() ServiceQuery {
	return &AdsQuery{
		QueryMap{
			"full":      "width=600&height=460",
			"thumbnail": "width=100",
		},
		OperationMap{
			"full":      Resize,
			"thumbnail": Thumbnail,
		},
	}
}

func init() {
	ServiceQueryRegister(queryTypeAds, newAdsQuery())
}
