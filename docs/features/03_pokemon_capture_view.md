# ポケモン捕獲・一覧確認機能 (Pokemon Capture & View)

## 機能概要
認証済みユーザーがポケモンを捕まえ、自身が過去に捕獲したポケモンの一覧を確認できるコア機能です。

## 詳細機能
- ランダムな遭遇、または特定の条件でのポケモン捕獲処理
- 捕獲したポケモンデータをユーザー情報と紐づけてデータベースに保存
- ユーザーごとの「マイポケモン一覧」の取得および画面表示

## 関連ファイル
- **Domain (Model)**: `internal/domain/model/pokemon.go`
- **Application (Usecase)**: `internal/application/usecase/pokemon_usecase.go`
- **UI (Handler)**: `internal/ui/handler/pokemon_handler.go`
