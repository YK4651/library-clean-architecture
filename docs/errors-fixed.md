# 修正したエラーの記録（文章形式）

ドメイン（User / UserID）の API を変更したあと、アプリ層・テスト・他ドメインサービスが古い API を前提にしていたため、複数箇所でコンパイルエラーが出ました。ここでは、どこでどんなエラーがあり、どう直したかを文章でまとめます。

---

## 全体の背景

User エンティティは「貸出数（currentLoanCount）をフィールドとして持たない」設計になっています。貸出数は Loan テーブルから導出する値なので、ユースケースやドメインサービスが引数で渡し、`User.CanBorrow(currentLoanCount int)` で判定する形に変わっています。また UserID は 8 桁の値となり、`NewUserID(value string)` は検証付きで `(*UserID, error)` を返し、ランダム生成には `GenerateUserID()` を使う API になっています。これらの変更に呼び出し側が追いついていなかったことが、今回のエラーの共通原因です。

---

## 1. CreateUserUseCase（アプリケーション層）

**ファイル:** `internal/application/command/createuser/create_user_usecase.go`

ユースケースの出力 DTO に `CurrentBorrowCount` を設定する際、存在しないメソッド `u.CurrentBorrowCount()` を呼んでいました。User は貸出数を持たないため、新規作成時は常に貸出数 0 です。そこで、`CurrentBorrowCount` にはリテラルの `0` を代入するように修正しました。

---

## 2. User ドメインのテスト（前半：以前からあったテスト）

**ファイル:** `internal/domain/userdm/user_test.go`（`TestCanCreateValidUser` や `TestCannotBorrowMoreWhenMaxLoansReached` など、別ファイルにあったテスト群）

ここでは、User にない `CurrentLoanCount()` や `CanBorrowMore()` を呼んでいました。また `ReconstructUser` に「貸出数」と「延滞料」の二つの数値を渡して 7 引数で呼んでいましたが、実際の `ReconstructUser` は 6 引数で、第 5 引数は `overdueFees` のみです。

修正内容は次のとおりです。`CurrentLoanCount()` のアサーションは削除しました。貸出可否のテストでは `CanBorrowMore()` の代わりに `CanBorrow(currentLoanCount)` を使い、貸出数は `CanBorrow(5)` や `CanBorrow(2)` の引数で渡すようにしました。`ReconstructUser` の呼び出しは 6 引数に揃え、第 5 引数は延滞料のみ（`overdueFees`）を渡す形にしました。

---

## 3. User ドメインのテスト（後半：TestUser と Immutability）

**ファイル:** `internal/domain/userdm/user_test.go`（`TestUser` とそのサブテスト）

`TestUser` 内では、引数なしの `u.CanBorrow()`、未公開フィールドの `u.currentLoanCount`、存在しない `u.BorrowBook()` と `u.ReturnBook()` を参照していました。さらに `ReconstructUser` を 7 引数で呼んでいる箇所がありました。

まず「Immutability」サブテストについてです。`u.ReturnBook()` は User に存在しないため、不変性の検証対象を「返却」から「延滞料の支払い」に変更しました。`ReconstructUser` で延滞料 10 の User を作り、`u.PayOverdueFee(3.0)` で新しい User を取得します。新しいインスタンスの延滞料が 7 になっていること、元の User の延滞料が 10 のままであることを確認する形にしました。

そのほかのサブテストは次のように直しました。`u.CanBorrow()` は `u.CanBorrow(0)` や `u.CanBorrow(MaxLoans)` のように、貸出数を引数で渡す呼び出しに変更しました。`u.currentLoanCount` への参照は削除しました。`u.BorrowBook()` を使っていたサブテストは、User が貸出数を持たない設計のため、「貸出数が上限のとき借りられない／上限未満なら借りられる」を検証する「CanBorrowAtLimit」に置き換え、`CanBorrow(MaxLoans)` と `CanBorrow(MaxLoans-1)` の結果で確認するようにしました。`ReconstructUser` の呼び出しは、第 5 引数が `overdueFees` のみになるよう 6 引数に統一しました。

---

## 4. UserID ドメインのテスト

**ファイル:** `internal/domain/userdm/user_id_test.go`

`NewUserID()` を引数なしで呼び、戻り値を 1 つの変数で受けていました。実際の `NewUserID(value string)` は 2 つ返すため、コンパイルエラー（引数不足と代入の個数不一致）になっていました。ランダムな ID が欲しいテストでは、引数なしで 1 つだけ返す `GenerateUserID()` に差し替えました。あわせて、ID 長の期待値が ULID 仕様の 26 のままだった箇所を、現在の 8 桁仕様に合わせて 8 に変更しました。

---

## 5. LoanEligibility のテスト

**ファイル:** `internal/domain/loandm/loan_eligibility_service_test.go`

`userdm.ReconstructUser` を 7 引数（貸出数と延滞料の二つの数値を含む）で呼んでいました。`ReconstructUser` は 6 引数で、貸出数は受け取らないため、余分な数値引数を削除し、第 5 引数は `overdueFees` のみを渡すようにしました。あわせて、サービス側のシグネチャ変更に合わせ、`CanBorrow(u, b)` を `CanBorrow(u, 0, b)` および `CanBorrow(u, 5, b)` のように、現在の貸出数（currentLoanCount）を第 2 引数で渡す形に修正しました。「貸出上限に達したユーザーは借りられない」テストでは、User を Active のままにして、`CanBorrow(u, 5, b)` で貸出数 5 のとき借りられないことを検証するようにしました。

---

## 6. LoanEligibilityService（ドメインサービス）

**ファイル:** `internal/domain/loandm/loan_eligibility_service.go`

サービス内で、User に存在しない `CanBorrowMore()`、`HasOverdueBooks()`、`GetMaxLoans()` を呼んでいました。

`CanBorrow(u, b)` を `CanBorrow(u, currentLoanCount int, b)` に変更し、呼び出し元が現在の貸出数を渡すようにしました。ルール 1（貸出上限）とルール 2（延滞）の判定は、User の `CanBorrow(currentLoanCount)` に任せる形にし、`CanBorrowMore()` と `HasOverdueBooks()` の呼び出しは削除しました。`IneligibilityReason` も同様に `currentLoanCount` を引数に追加し、`u.CanBorrow(currentLoanCount)` が false の理由を、`userdm.MaxLoans` と `u.OverdueFees()` を使ってメッセージに組み立てるようにしました。`GetMaxLoans()` の代わりには、定数 `userdm.MaxLoans` を参照するようにしました。

---

## まとめ

出たエラーは、いずれも「ドメインの API と、それを使う側の前提が食い違っていた」ことが原因です。具体的には、(1) User に存在しないメソッドやフィールド（CurrentBorrowCount、CurrentLoanCount、CanBorrowMore、HasOverdueBooks、GetMaxLoans、BorrowBook、ReturnBook）を呼んでいたこと、(2) 関数のシグネチャと引数・戻り値の個数が合っていなかったこと（NewUserID の引数と戻り値、ReconstructUser の引数の数）、の二つに整理できます。

対応方針は、ドメインの現在の API に合わせて呼び出し側を修正することです。貸出数は User のフィールドではなく、ユースケースやサービスが引数で渡し、`CanBorrow(currentLoanCount)` で判定する形に統一しました。ID の生成・検証は、ランダムなら `GenerateUserID()`、指定値の検証付きなら `NewUserID(value)` を使うようにし、`ReconstructUser` は 6 引数（id, name, email, status, overdueFees, createdAt）で呼ぶようにしました。
