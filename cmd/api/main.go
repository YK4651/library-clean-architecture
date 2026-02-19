package main

import (
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"

	"github.com/YK4651/library-clean-architecture/internal/application/extendloan"
	"github.com/YK4651/library-clean-architecture/internal/application/returnbook"
	"github.com/YK4651/library-clean-architecture/internal/infrastructure/http/handler"
	"github.com/YK4651/library-clean-architecture/internal/infrastructure/http/middleware"
	"github.com/YK4651/library-clean-architecture/internal/infrastructure/persistence"
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
	loanRepo := persistence.NewMySQLLoanRepository(db)
	userRepo := persistence.NewMySQLUserRepository(db)
	returnBookUC := returnbook.NewReturnBookUseCase(loanRepo)
	r.POST("/api/loans/:loanID/returns", handler.ReturnBook(returnBookUC))

	extendLoanUC := extendloan.NewExtendLoanUseCase(userRepo, loanRepo)
	r.POST("/api/books/:bookID/loans/:loanID/extend", middleware.AuthUserID(), handler.ExtendLoan(extendLoanUC))

	r.Run(":8080")
}
