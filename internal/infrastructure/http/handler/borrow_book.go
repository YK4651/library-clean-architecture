package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BorrowBook は POST /api/books/borrow のハンドラです。
// リクエストボディ: { "user_id": "...", "book_id": "..." }
// ctx.Request.Context().Value("db") で *sql.DB または *sql.Tx を取得可能
func BorrowBook(c *gin.Context) {
	_ = c.Request.Context().Value("db") // トランザクションはミドルウェアが設定
	c.JSON(http.StatusNotImplemented, gin.H{"error": "not implemented"})
}
