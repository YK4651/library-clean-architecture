package handler

import (
	"errors"
	"net/http"

	"github.com/YK4651/library-clean-architecture/internal/application/extendloan"
	"github.com/YK4651/library-clean-architecture/internal/infrastructure/http/middleware"
	"github.com/gin-gonic/gin"
)

// ExtendLoan は POST /books/:bookID/loans/:loanID/extend のハンドラ。
// userId は認証（JWT/セッション）から取得。開発用に middleware.AuthUserID で X-User-ID を参照。
func ExtendLoan(uc *extendloan.ExtendLoanUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		bookID := c.Param("bookID")
		loanID := c.Param("loanID")
		userIDVal, ok := c.Get(middleware.ContextKeyUserID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			return
		}
		userID := userIDVal.(string)

		req, err := extendloan.NewExtendLoanRequest(bookID, loanID, userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		resp, err := uc.Execute(ctx, req)
		if err != nil {
			switch {
			case errors.As(err, new(*extendloan.UserNotFoundError)):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.As(err, new(*extendloan.UserSuspendedError)):
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			case errors.As(err, new(*extendloan.UserHasOverdueLoansError)):
				c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			case errors.As(err, new(*extendloan.LoanNotFoundOrForbiddenError)):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.As(err, new(*extendloan.BookLoanMismatchError)):
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			case errors.As(err, new(*extendloan.LoanAlreadyReturnedError)):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
