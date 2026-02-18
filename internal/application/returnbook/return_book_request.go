package returnbook

import "errors"

type ReturnBookRequest struct {
	LoanID string
}

func NewReturnBookRequest(loanID string) (*ReturnBookRequest, error) {
	if loanID == "" {
		return nil, errors.New("loanID cannot be empty")
	}
	return &ReturnBookRequest{
		LoanID: loanID,
	}, nil
}
