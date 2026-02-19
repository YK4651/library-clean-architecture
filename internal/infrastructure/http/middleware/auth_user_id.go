package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ContextKeyUserID はコンテキストに格納するユーザーIDのキー
const ContextKeyUserID = "userID"

// AuthUserID は認証済み userId をコンテキストに設定するミドルウェア。
// 本番では JWT トークン/セッションから userId を取得すること。ここでは開発用に X-User-ID ヘッダーを使用。
func AuthUserID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-ID")
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization required"})
			c.Abort()
			return
		}
		c.Set(ContextKeyUserID, userID)
		c.Next()
	}
}
