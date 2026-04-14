---
name: requirements_definition
description: This skill should be used when the user wants to "create a new feature", "implement a change", "define requirements", or asks to "start a new task". It provides a structured process for Creating a Product Requirements Document (PRD) before any code is written.
---

# Requirements Definition Skill

新規機能や大幅な改修を行う前に、ユーザーと「何を作るか（What）」および「なぜ作るか（Why）」を明確化せよ。コード（How）を書き始める前に、必ず以下のプロセスを実行すること。

## 要件引き出しのステップ

以下の手順でユーザーとディスカッションを行い、要件を定義せよ。

1. **目的の確認 (Objective)**: 機能が必要な理由、解決したい課題、期待するユーザー体験を確認する。
2. **ユーザーストーリーの記述**: 「誰が」「何を」「どうしたいのか」を列挙し、機能要件を抽出する。
3. **境界条件・エッジケースの検討**: 例外ケースやエラーハンドリング、予期せぬ入力への対応についてユーザーに質問する。
4. **スコープの定義 (In/Out of Scope)**: 実装に含める内容と、意図的に「やらないこと」を明確にし、スコープクリープを防ぐ。

## 出力：PRD (Product Requirements Document) アーティファクト

要件が合意に至ったら、`prd.md` という名称でアーティファクトを生成せよ。記載すべき項目については、後述の `examples/template_prd.md` を参照すること。生成後、必ずユーザーの承認（Approve）を得ること。

## 行動指針

- 必ずユーザーへ質問し、独断で要件を決定しない。
- PRDの承認を得るまでは、実装計画（`implementation_plan`）やコードの実装に着手しない。
- プロジェクト独自の「可愛さ（Aesthetic）」や世界観が、UI/UXにどう反映されるべきかを確認する。

## 関連リソース

### Examples
- **`examples/template_prd.md`** - 標準的なPRDのフォーマットと記載例。
