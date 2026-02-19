package extendloan

// ExtendLoanResponse は延長結果
type ExtendLoanResponse struct {
	NewDueDate   string `json:"newDueDate"`
	ExtendedDays int    `json:"extendedDays"`
}
