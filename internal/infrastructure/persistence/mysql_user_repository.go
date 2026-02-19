package persistence

import (
	"context"
	"database/sql"
	"time"

	"github.com/YK4651/library-clean-architecture/internal/domain/userdm"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewMySQLUserRepository(db *sql.DB) userdm.IUserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) FindByID(ctx context.Context, id *userdm.UserID) (*userdm.User, error) {
	query := "SELECT id, name, email, status, created_at FROM users WHERE id = ?"

	var idStr, name, email string
	var status uint8
	var createdAt time.Time
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(&idStr, &name, &email, &status, &createdAt)
	_ = idStr // id already known from parameter
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Use the same id we queried with (caller's *UserID); DB may store ULID while domain validates 8-digit
	// status: 1=ACTIVE, 2=SUSPENDED, 3=INACTIVE (schema)
	userStatus := userdm.UserStatusActive
	if status == 2 {
		userStatus = userdm.UserStatusSuspended
	}

	return userdm.ReconstructUser(
		id,
		name,
		email,
		userStatus,
		0, // overdue_fees not in users table
		createdAt,
	), nil
}
