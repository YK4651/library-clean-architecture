package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/YK4651/library-clean-architecture/internal/infrastructure/http/handler"
	"github.com/YK4651/library-clean-architecture/internal/infrastructure/http/middleware"
)

func main() {
	// データベース接続の初期化
	db, err := sql.Open("mysql", "user:password@tcp(localhost:3306)/library_db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ginルーターのセットアップ
	r := gin.Default()

	// トランザクションミドルウェアをグローバルに適用
	r.Use(middleware.TransactionMiddleware(db))

	// ルーティング
	r.POST("/api/books/borrow", handler.BorrowBook)
	// ハンドラー内で ctx.Request.Context().Value("db") からDB/txを取得

	r.Run(":8080")
}
