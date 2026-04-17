# Battle System Design Spec

## 1. 概要
`zakopokeGo` にNPCトレーナーとのコマンド式バトル機能を追加する。
捕まえたポケモンには固有の「わざ」と「ステータス」を持たせ、勝利することでレベルアップする成長要素を導入する。

## 2. 目標 (Success Criteria)
- 自分のポケモン1体 vs NPCのポケモン1体のターン制バトルができる。
- 捕獲時にPokeAPIからランダムで「わざ」を1つ取得し、保存できる。
- タイプ相性（例：水は炎に強い）がダメージ計算に反映される。
- バトルに勝利するとレベルが上がり、ステータスが上昇する。

## 3. ユーザー体験 (UX)
1. ポケモンを捕まえた際、「〇〇は [わざ名] を覚えている！」というメッセージが表示される。
2. ボックス画面などから「バトル」を選択すると、ランダムなNPCトレーナーが現れる。
3. バトル画面では、交互にコマンドを選択して相手のHPを0にする。
4. 勝利すると経験値を獲得し、一定量でレベルアップ。ステータス上昇の演出が出る。

## 4. システムアーキテクチャ (DDD)

### Domain Layer (`internal/domain`)
- **Model**: `Pokemon` 構造体に `Level`, `Exp`, `BaseHP`, `BaseAttack`, `MoveName`, `MoveType`, `MovePower` などのフィールドを追加。
- **Service**: `BattleService` を作成。ダメージ計算、タイプ相性、勝敗判定、経験値計算のロジックを持つ。
- **Repository**: ポケモンの新規項目を保存できるように DB スキーマ（GORM）を更新。

### Application Layer (`internal/application`)
- **UseCase**: `BattleUseCase` を作成。バトルの開始、ターンの進行、勝敗後の報酬処理などをオーケストレートする。

### Infrastructure Layer (`internal/infrastructure`)
- **External**: `PokeAPIRepository` を拡張し、ポケモンの基本ステータス（Stats）と覚える技（Moves）のリストを取得する機能を追加。
- **Persistence**: `Pokemon` テーブルに新しいカラムを追加。

### UI Layer (`internal/ui`)
- **Handler**: `/battle` エンドポイントの実装。バトル画面のレンダリング。
- **Template**: `templates/battle.html` の新規作成。

## 5. データモデル (Pokemon)
```go
type Pokemon struct {
    ID          uint
    UserID      uint
    PokemonNo   int
    Level       int
    Exp         int
    MaxHP       int
    CurrentHP   int
    Attack      int
    Defense     int
    MoveName    string
    MoveType    string
    MovePower   int
}
```

## 6. タイプ相性
簡易的なタイプ相性テーブル（Map）をコード内に定義する。
例: `{"water": {"fire": 2.0, "grass": 0.5}}`

## 7. テスト計画
- ダメージ計算ロジックの単体テスト。
- タイプ相性が正しく反映されるかのテスト。
- レベルアップ時のステータス上昇計算のテスト。

## 8. WOW要素
バトルのログ（「〇〇の こうげき！ こうかは ばつぐんだ！」）に動的なアニメーションや、勝利時のファンファーレ風 UI を追加し、チープながらも「遊んでいる感」を演出する。
