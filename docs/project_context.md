# ZakopokeGo - プロジェクト概要 & AI駆動開発ガイド

## プロジェクト概要 (Overview)
**ZakopokeGo** は、Go言語とPokeAPIを利用して作成されたWebアプリケーションです。
主な機能として、「ポケモンの取得（PokeAPI利用）」「ポケモンの捕獲」「捕獲したポケモンの確認」を提供します。ドメイン駆動設計（DDD）に基づいたレイヤードアーキテクチャを採用しており、関心の分離を意識した堅牢な構造を持っています。

## 技術スタック (Tech Stack)
- **言語**: Go (1.24)
- **Webフレームワーク**: Gin (`github.com/gin-gonic/gin`)
- **ORM・データベース**: GORM (`gorm.io/gorm`), SQLite (`pokeapp.db`)
- **セッション・認証**: `github.com/gin-contrib/sessions`, `golang.org/x/crypto` (パスワードハッシュ化)
- **外部API**: [PokeAPI](https://pokeapi.co/)
- **フロントエンド**: Go `html/template` (templatesディレクトリ) + 静的ファイル (staticディレクトリ)

---

## 機能一覧 (Features)

現在実装されている各機能の詳細は `docs/features/` ディレクトリ配下のドキュメントを参照してください。
今後機能が追加された場合は、同ディレクトリに新しいマークダウンファイルを追加し管理します。

- [ユーザー認証 (User & Auth)](./features/01_user_auth.md)
- [外部ポケモンデータ取得 (PokeAPI Integration)](./features/02_pokeapi_integration.md)
- [ポケモン捕獲・一覧確認機能 (Pokemon Capture & View)](./features/03_pokemon_capture_view.md)

---

## アーキテクチャ・ディレクトリ構成 (Architecture)

DDD（ドメイン駆動設計）のレイヤードアーキテクチャを採用しています。各層の責務を分けて実装されています。

```text
zakopokeGo/
├── cmd/
│   └── app/
│       └── main.go                  # エントリーポイント。サーバー起動と各レイヤーのDI（依存注入）を設定します。
├── internal/
│   ├── domain/                      # 【ドメイン層】ビジネスルール、Entity、操作のインターフェース
│   │   ├── model/                   # データベースやAPIの構造に依存しない純粋なデータの形 (Entity: user.go, pokemon.go)
│   │   └── repository/              # データ永続化・外部連携のインターフェース (repository.go 等)
│   ├── application/                 # 【アプリケーション層】ユースケースの実現
│   │   └── usecase/                 # フロントエンド(UI)からの要求を受け取り、ドメイン層の機能を組み合わせて処理 (auth_usecase.go, pokemon_usecase.go)
│   ├── infrastructure/              # 【インフラストラクチャ層】技術的詳細の実装
│   │   ├── persistence/             # GORM(SQLite)を用いたデータベースアクセスの実装
│   │   └── external/                # PokeAPIへのHTTPリクエストなどの外部連携実装
│   └── ui/                          # 【プレゼンテーション層】ユーザーとのインターフェース
│       └── handler/                 # Ginを使ったHTTPリクエストの受付・ルーティング処理、レスポンスの生成 (auth_handler.go 等)
├── templates/                       # HTMLテンプレートファイル
├── static/                          # CSS・JS・画像などの静的ファイル
├── pokeapp.db                       # SQLiteデータベースファイル（ローカル永続化用）
├── go.mod / go.sum                  # モジュール管理ファイル
└── README.md                        # 一般的なプロジェクト導入ドキュメント
```

---

## AI Agent (AI駆動開発) へのガイドライン

今後AIエージェントに機能追加や改修を依頼する場合は、以下のルールを遵守して開発を進めてください：

1. **レイヤーの依存関係を厳守する**
   - **依存の方向**: `ui` → `application` → `domain` ← `infrastructure`
   - `domain`層は外部のフレームワークやDB技術（GORMやGinなど）に一切依存してはいけません。
   - `infrastructure`層で実装するDB（SQLite）や外部APIの知識を、他の層に漏らさないようインターフェース（`domain/repository`）経由で操作してください。

2. **機能追加の手順**
   - まず `domain/model` で必要なEntityやStructを定義する。
   - 必要に応じて `domain/repository` にインターフェースを定義する。
   - `application/usecase` にビジネスロジックを実装する。
   - `infrastructure/persistence` や `infrastructure/external` で具体的な技術的処理（SQLクエリ、HTTP呼び出し）を実装する。
   - 最後に `ui/handler` を追加・修正し、`templates` で画面に反映する。

3. **ローカル環境起動とライブラリ管理**
   - 新しいパッケージを入れた場合は速やかに `go mod tidy` を実行する。
   - 動作確認は `go run cmd/app/main.go` にてサーバーを起動し、ブラウザ上で確認を行うこと。

このドキュメントはアプリケーションの成長に合わせて随時アップデートしてください。
