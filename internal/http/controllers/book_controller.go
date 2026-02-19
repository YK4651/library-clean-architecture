package controllers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/application/query/getbook"
	"github.com/YK4651/library-clean-architecture/internal/application/query/listbooks"
	"github.com/gorilla/mux"
)

// BookController - 書籍APIコントローラー（プレゼンテーション層）
type BookController struct {
	getBookUseCase   *getbook.GetBookUseCase
	listBooksUseCase *listbooks.ListBooksUseCase
}

// NewBookController - 新しいBookControllerを作成
func NewBookController(getBookUseCase *getbook.GetBookUseCase, listBooksUseCase *listbooks.ListBooksUseCase) *BookController {
	return &BookController{
		getBookUseCase:   getBookUseCase,
		listBooksUseCase: listBooksUseCase,
	}
}

// GetBook - GET /api/books/{bookID} を処理
func (c *BookController) GetBook(w http.ResponseWriter, r *http.Request) {
	// 1. URLパスからbookIdを抽出
	vars := mux.Vars(r)
	bookID := vars["bookID"]

	// 2. リクエストDTOを作成
	request, err := getbook.NewGetBookRequest(bookID)
	if err != nil {
		// リクエストDTOのバリデーション失敗
		c.sendJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, http.StatusBadRequest) // 400 Bad Request
		return
	}

	// 3. ユースケースを実行
	response := c.getBookUseCase.Execute(r.Context(), request)

	// 4. 失敗レスポンスを処理
	if !response.Success {
		c.handleError(w, *response.ErrorMessage)
		return
	}

	// 5. 成功レスポンスを返す
	c.sendJSON(w, map[string]interface{}{
		"success": true,
		"data": map[string]interface{}{
			"id":          response.Data.ID,
			"isbn":        response.Data.ISBN,
			"title":       response.Data.Title,
			"author":      response.Data.Author,
			"isAvailable": response.Data.IsAvailable,
			"currentLoan": response.Data.CurrentLoan,
		},
	}, http.StatusOK) // 200 OK
}

// ListBooks - GET /books?limit=20&offset=0 を処理（ページネーション、空の場合は200で空配列）
func (c *BookController) ListBooks(w http.ResponseWriter, r *http.Request) {
	limitParam := r.URL.Query().Get("limit")
	offsetParam := r.URL.Query().Get("offset")

	request, err := listbooks.NewListBooksRequest(limitParam, offsetParam)
	if err != nil {
		c.sendJSON(w, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	response := c.listBooksUseCase.Execute(r.Context(), request)
	if !response.Success {
		c.handleError(w, *response.ErrorMessage)
		return
	}

	// API仕様: books, total, limit, offset（日付は YYYY-MM-DD）
	bookItems := make([]map[string]interface{}, 0, len(response.Data.Books))
	for _, b := range response.Data.Books {
		item := map[string]interface{}{
			"id":          b.ID,
			"isbn":        b.ISBN,
			"title":       b.Title,
			"author":      b.Author,
			"isAvailable": b.IsAvailable,
			"currentLoan": nil,
		}
		if b.CurrentLoan != nil {
			item["currentLoan"] = map[string]interface{}{
				"loanId":       b.CurrentLoan.LoanID,
				"userId":       b.CurrentLoan.UserID,
				"borrowedDate": formatDateOnly(b.CurrentLoan.BorrowedDate),
				"dueDate":      formatDateOnly(b.CurrentLoan.DueDate),
			}
		}
		bookItems = append(bookItems, item)
	}

	c.sendJSON(w, map[string]interface{}{
		"books":  bookItems,
		"total":  response.Data.Total,
		"limit":  response.Data.Limit,
		"offset": response.Data.Offset,
	}, http.StatusOK)
}

// formatDateOnly - "2006-01-02 15:04:05" または "2006-01-02" を "2006-01-02" に変換
func formatDateOnly(s string) string {
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t.Format("2006-01-02")
	}
	if t, err := time.Parse("2006-01-02", s); err == nil {
		return t.Format("2006-01-02")
	}
	return s
}

// handleError - エラーメッセージを適切なHTTPステータスコードにマッピング
func (c *BookController) handleError(w http.ResponseWriter, errorMessage string) {
	statusCode := http.StatusInternalServerError // デフォルト: 500

	if strings.Contains(strings.ToLower(errorMessage), "not found") {
		statusCode = http.StatusNotFound // 404 Not Found
	} else if strings.Contains(errorMessage, "Invalid") {
		statusCode = http.StatusBadRequest // 400 Bad Request
	}

	c.sendJSON(w, map[string]interface{}{
		"success": false,
		"error":   errorMessage,
	}, statusCode)
}

// sendJSON - JSONレスポンスを送信するヘルパー関数
func (c *BookController) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(data)
}
