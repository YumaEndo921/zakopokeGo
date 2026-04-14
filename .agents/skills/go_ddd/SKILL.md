---
name: go_ddd
description: zakopokeGoプロジェクトにおけるDDD（ドメイン駆動設計）の実装ガイドライン。AIがレイヤードアーキテクチャに則ったコード生成を行うための指示書。
---

# Go DDD Implementation Skill

このスキルは、`zakopokeGo` プロジェクトにおけるDDDレイヤードアーキテクチャの標準的な実装パターンを定義します。新しい機能（モデル、ユースケース、APIエンドポイント）を追加する際は、必ずこのガイドラインに従ってください。

## アーキテクチャ構成

### 1. Domain Layer (`internal/domain`)
- **Model (`model/`)**: 他のレイヤーに依存しない純粋なデータ構造とコンストラクタを配置します。
- **Repository Interface (`repository/`)**: データの永続化や外部連携のためのインターフェースを定義します。具体的な実装（SQLやHTTPリクエスト）は含めません。

### 2. Application Layer (`internal/application/usecase`)
- ドメインオブジェクトを操作し、ビジネスロジック（ユースケース）を実現します。
- リポジトリインターフェースを介してデータにアクセスすることで、具体的なDB実装から絶縁します（依存性の逆転）。

### 3. Infrastructure Layer (`internal/infrastructure`)
- **Persistence (`persistence/`)**: DB（GORMなど）への具体的なアクセスコード。ドメイン層のリポジトリインターフェースを実装します。
- **External (`external/`)**: 外部API（PokeAPIなど）への具体的なリクエストコード。

### 4. UI Layer (`internal/ui/handler`)
- HTTPリクエストの受付とレスポンスの返却を担当します。
- ユースケースを呼び出しますが、ドメインロジックをここに書かないでください。

## 実装のルール

- **ファイルの命名**: 
  - ハンドラー: `xxx_handler.go`
  - ユースケース: `xxx_usecase.go`
  - リポジトリ: `xxx_repository.go`
- **依存性注入 (DI)**: コンストラクタ（`NewXxx`）を使用して依存関係を注入します。
- **エラーハンドリング**: 下位レイヤーのエラーは上位レイヤーに伝播させ、UI層で適切にレンダリングまたはJSON返却します。

## 新機能追加のステップ

1. `internal/domain/model` にエンティティを定義。
2. `internal/domain/repository` に必要なインターフェースを定義。
3. `internal/infrastructure` にインターフェースの具象実装を作成。
4. `internal/application/usecase` にビジネスロジックを実装。
5. `internal/ui/handler` にエンドポイントを実装し、`cmd/app/main.go` でルーティングを登録。
