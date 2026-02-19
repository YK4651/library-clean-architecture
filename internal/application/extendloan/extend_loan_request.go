package extendloan

import "errors"

// ExtendLoanRequest は延長リクエスト。UserIDは認証（JWT/セッション）からハンドラが設定する。
type ExtendLoanRequest struct {
	BookID string
	LoanID string
	UserID string
}

// NewExtendLoanRequest はパスパラメータと認証済みUserIDからリクエストを生成する。
func NewExtendLoanRequest(bookID, loanID, userID string) (*ExtendLoanRequest, error) {
	if bookID == "" {
		return nil, errors.New("bookID cannot be empty")
	}
	if loanID == "" {
		return nil, errors.New("loanID cannot be empty")
	}
	if userID == "" {
		return nil, errors.New("userID cannot be empty")
	}
	return &ExtendLoanRequest{
		BookID: bookID,
		LoanID: loanID,
		UserID: userID,
	}, nil
}
