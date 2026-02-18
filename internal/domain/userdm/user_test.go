package userdm

import (
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	t.Run("CanBorrow", func(t *testing.T) {
		u := NewUser("田中太郎", "tanaka@example.com")

		if !u.CanBorrow(0) {
			t.Error("新規ユーザーは本を借りられるべきです（貸出数0）")
		}

		if u.name != "田中太郎" {
			t.Errorf("名前は '田中太郎' であるべきですが、'%s' でした", u.name)
		}
		if u.email != "tanaka@example.com" {
			t.Errorf("メールは 'tanaka@example.com' であるべきですが、'%s' でした", u.email)
		}
		if u.status != UserStatusActive {
			t.Errorf("ステータスは 'active' であるべきですが、'%s' でした", u.status)
		}
	})

	t.Run("CanBorrowAtLimit", func(t *testing.T) {
		u := NewUser("田中太郎", "tanaka@example.com")
		// 貸出数が MaxLoans なら借りられない
		if u.CanBorrow(MaxLoans) {
			t.Error("貸出数が上限のユーザーは本を借りられないべきです")
		}
		if !u.CanBorrow(MaxLoans - 1) {
			t.Error("貸出数が上限未満のユーザーは借りられるべきです")
		}
	})

	t.Run("CannotBorrowWhenSuspended", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}

		// 停止中のユーザーを再構築
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusSuspended,
			0, // overdueFees
			time.Now(),
		)

		if u.CanBorrow(0) {
			t.Error("停止中のユーザーは本を借りられないべきです")
		}
	})

	t.Run("CannotBorrowAtMaxLimit", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}

		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusActive,
			0, // overdueFees
			time.Now(),
		)
		// 貸出数が MaxLoans のとき借りられない
		if u.CanBorrow(MaxLoans) {
			t.Error("最大貸出数に達したユーザーは本を借りられないべきです")
		}
	})

	t.Run("CannotBorrowWithOverdueFees", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}

		// 延滞料金があるユーザーを再構築
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusActive,
			10.50, // overdueFees
			time.Now(),
		)

		if u.CanBorrow(0) {
			t.Error("延滞料金があるユーザーは本を借りられないべきです")
		}
	})

	t.Run("FeeManagement", func(t *testing.T) {
		u := NewUser("田中太郎", "tanaka@example.com")

		u2, err := u.AddOverdueFee(5.00)
		if err != nil {
			t.Fatalf("延滞料金の追加に失敗しました: %v", err)
		}

		if u2.overdueFees != 5.00 {
			t.Errorf("延滞料金は 5.00 であるべきですが、%.2f でした", u2.overdueFees)
		}

		// 負の料金をテスト
		_, err = u.AddOverdueFee(-1.00)
		if err == nil {
			t.Error("負の料金に対してエラーが期待されました")
		}
	})

	t.Run("Immutability", func(t *testing.T) {
		userID, err := NewUserID("12345678")
		if err != nil {
			t.Fatal(err)
		}
		u := ReconstructUser(
			userID,
			"田中太郎",
			"tanaka@example.com",
			UserStatusActive,
			10.0, // overdueFees
			time.Now(),
		)

		// PayOverdueFee の不変性をテスト
		u2, err := u.PayOverdueFee(3.0)
		if err != nil {
			t.Fatalf("延滞料金の支払いに失敗しました: %v", err)
		}

		if u2.overdueFees != 7.0 {
			t.Errorf("支払い後の延滞料金は 7.0 であるべきですが、%.2f でした", u2.overdueFees)
		}

		// 元のユーザーは変更されないべき
		if u.overdueFees != 10.0 {
			t.Error("元のユーザーは変更されないべきです")
		}
	})
}
