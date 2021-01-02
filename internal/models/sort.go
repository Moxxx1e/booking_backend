package models

type Sort struct {
	OrderBy string `query:"order_by"`
	Desc    bool   `query:"desc"`
}
