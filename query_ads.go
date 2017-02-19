package main

import "fmt"

const queryTypeAds string = "ads"

type AdsQuery struct {
	queries QueryMap
	ops     OperationMap
}

func (q *AdsQuery) getQuery(ot string) string {
	return fmt.Sprintf("%s&cpool=%s", q.queries[ot], queryTypeAds)
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
