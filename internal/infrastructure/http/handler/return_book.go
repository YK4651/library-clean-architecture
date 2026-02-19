package handler

import (
	"errors"
	"net/http"

	"github.com/YK4651/library-clean-architecture/internal/application/returnbook"
	"github.com/gin-gonic/gin"
)

// ReturnBook は POST /loans/:loanID/returns のハンドラを返す。
func ReturnBook(uc *returnbook.ReturnBookUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		loanID := c.Param("loanID")
		if loanID == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "loanID is required"})
			return
		}

		req, err := returnbook.NewReturnBookRequest(loanID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ctx := c.Request.Context()
		resp, err := uc.Execute(ctx, req)
		if err != nil {
			switch {
			case errors.As(err, new(*returnbook.LoanNotFoundError)):
				c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			case errors.As(err, new(*returnbook.LoanAlreadyReturnedError)):
				c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.JSON(http.StatusOK, resp)
	}
}
