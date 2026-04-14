# ユーザー認証 (User & Auth)

## 機能概要
ユーザーの登録、ログイン処理、およびセッション情報などを管理する機能です。
セッションベースの認証を提供し、各保護されたエンドポイントへのアクセスを制御します。

## 詳細機能
- 新規ユーザー登録（パスワードのハッシュ化を含む）
- ログイン・ログアウト
- セッションベースの認証・状態管理

## 関連ファイル
- **Domain (Model)**: `internal/domain/model/user.go`
- **Application (Usecase)**: `internal/application/usecase/auth_usecase.go`
- **UI (Handler)**: `internal/ui/handler/auth_handler.go`
