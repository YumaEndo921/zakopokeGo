# 概要

pokeAPI と Go を使ってみたくて作ったチープなアプリです。
ギャル（chatGPT）と一緒に作りました。

# 起動方法 (How to Run)

Go言語の開発環境が整っている前提です。
ターミナルでこのプロジェクトのルートディレクトリにいる状態で、以下のコマンドを実行してください。

```bash
go run cmd/app/main.go
```

## 初回起動時・エラーが出る場合

もし依存関係（ライブラリ）のエラーが出る場合は、以下のコマンドで整理してください。

```bash
go mod tidy
```

そのあと、再度 `go run cmd/app/main.go` を実行してください。

## 動作確認

起動に成功すると、以下のようなログが表示されます。

```text
2026/01/21 16:45:00 Server starting on :8080
[GIN-debug] Listening and serving HTTP on :8080
```

ブラウザで `http://localhost:8080` にアクセスして遊んでください！

# 機能

ポケモンを捕まえたり、捕まえたポケモンを確認できます。
![スクリーンショット 2025-04-23 22 34 33](https://github.com/user-attachments/assets/d584153d-733c-488f-89c8-e11613d6756c)

# 開発メモ (Architecture)

このプロジェクトは **ドメイン駆動設計 (DDD)** のレイヤードアーキテクチャを採用しています。

- **`cmd/app/main.go`**: アプリケーションのエントリーポイント。
- **`internal/domain`**: ビジネスロジック（Entity）とインターフェース（Repository）。
- **`internal/application`**: ユースケース（Use Case）。
- **`internal/infrastructure`**: データベース接続や外部APIの実装。
- **`internal/ui`**: HTTPハンドラー（Gin）。
