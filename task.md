# バトル機能実装タスクリスト

## Phase 1: データモデルとインフラの拡張
- [x] `internal/domain/model/pokemon.go` のフィールド拡張
- [x] `internal/infrastructure/persistence/db.go` のDBモデル拡張と変換ロジック更新
- [x] `internal/infrastructure/external/pokeapi_repository.go` でのステータス・技取得の実装
- [x] 動作確認（新規捕獲時にデータが保存されること）

## Phase 2: バトルロジックの実装
- [x] `internal/domain/service/battle_service.go` の作成（ダメージ計算・タイプ相性）
- [x] `internal/domain/service/experience_service.go` の作成（レベルアップ・成長）
- [x] `internal/application/usecase/battle_usecase.go` の作成（バトルフロー管理）

## Phase 3: UI/UXの構築
- [x] `internal/ui/handler/battle_handler.go` の作成
- [x] `templates/battle.html` の新規作成
- [x] CSSパーツ（HPバー、アニメーション）の実装
- [x] 既存画面（ボックス画面など）にバトル開始ボタンを追加

## Phase 4: ブラッシュアップ・テスト
- [x] レベルアップ時の演出追加（JS/CSS）
- [x] 各種エッジケース（HP 0、にげる、レベル最大など）の動作確認
- [x] 単体テストの作成
