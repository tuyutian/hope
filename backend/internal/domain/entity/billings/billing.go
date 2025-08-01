package billings

import "github.com/shopspring/decimal"

type CommissionListResponse struct {
	List  []*CommissionBill `json:"list"`
	Total int64             `json:"total"`
}
type BillingSummaryResponse struct {
	List  []*BillingPeriodSummary `json:"list"`
	Total int64                   `json:"total"`
}
type CurrentPeriodResponse struct {
	PeriodStart int64           `json:"period_start"`
	PeriodEnd   int64           `json:"period_end"`
	Amount      decimal.Decimal `json:"amount"`
}
