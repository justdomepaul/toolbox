package entity

type CommonListResponse struct {
	Count int64 `spanner:"Count" json:"Count,omitempty"`
	Row   int64 `spanner:"Row" json:"Row,omitempty"`
	Page  int64 `spanner:"Page" json:"Page,omitempty"`
}
