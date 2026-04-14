# 外部ポケモンデータ取得 (PokeAPI Integration)

## 機能概要
外部の公式Pokemon API（PokeAPI）と連携して、アプリケーション内で利用するポケモンの各種情報やメタデータを取得・変換する機能です。

## 詳細機能
- PokeAPI (`https://pokeapi.co/`) へのHTTPリクエスト送信
- 取得したJSONデータのアプリケーション内ドメインエンティティへの変換・マッピング機能

## 関連ファイル
- **Infrastructure (Repository / External)**: `internal/infrastructure/external/pokemon_metadata_repository.go` (※または該当パスに格納されているレポジトリ実装)
- **Domain (Repository interface)**: `internal/domain/repository/pokemon_metadata_repository.go`
