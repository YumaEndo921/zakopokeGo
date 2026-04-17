package usecase

import (
	"fmt"
	"math/rand"
	"time"
	"zakopokeGo/internal/domain/model"
	"zakopokeGo/internal/domain/repository"
	"zakopokeGo/internal/domain/service"
)

type BattleResult struct {
	MyPokemon      *model.Pokemon
	EnemyPokemon   *model.Pokemon
	MyName         string
	EnemyName      string
	MyHPPercent    int
	EnemyHPPercent int
	Logs           []string
	IsFinished     bool
	Winner         string // "player" or "npc"
	LeveledUp      bool
}

type BattleUseCase interface {
	StartBattle(userID uint, myPokemonID uint) (*BattleResult, error)
	ExecuteTurn(userID uint, myPokemonID uint, action string, myName string, enemyName string, enemyNo int, enemyLevel int, enemyHP int, enemyMaxHP int, enemyAtk int, enemyDef int, enemyMoveName string, enemyMoveType string, enemyMovePower int, enemyType1 string, enemyType2 string) (*BattleResult, error)
}

type battleUseCase struct {
	pokemonRepo     repository.PokemonRepository
	metadataRepo    repository.PokemonMetadataRepository
	battleService   service.BattleService
	experienceService service.ExperienceService
}

func NewBattleUseCase(
	pokemonRepo repository.PokemonRepository,
	metadataRepo repository.PokemonMetadataRepository,
	battleService service.BattleService,
	experienceService service.ExperienceService,
) BattleUseCase {
	return &battleUseCase{
		pokemonRepo:     pokemonRepo,
		metadataRepo:    metadataRepo,
		battleService:   battleService,
		experienceService: experienceService,
	}
}

func (u *battleUseCase) StartBattle(userID uint, myPokemonID uint) (*BattleResult, error) {
	myPokemon, err := u.pokemonRepo.FindByID(myPokemonID)
	if err != nil {
		return nil, err
	}
	if myPokemon.UserID != userID {
		return nil, fmt.Errorf("unauthorized")
	}

	if myPokemon.Level < 1 {
		myPokemon.Level = 1
	}

	// 古いデータ向けの初期化ロジック
	if myPokemon.MaxHP == 0 {
		meta, _ := u.metadataRepo.GetPokemonMetadata(myPokemon.PokemonNo)
		if meta != nil {
			myPokemon.MaxHP = meta.Stats["hp"]
			myPokemon.Attack = meta.Stats["attack"]
			myPokemon.Defense = meta.Stats["defense"]
			if len(meta.Moves) > 0 {
				myPokemon.MoveName = meta.Moves[0].Name
				myPokemon.MoveType = meta.Moves[0].Type
				myPokemon.MovePower = meta.Moves[0].Power
			}
			myPokemon.Type1 = meta.Types[0]
			if len(meta.Types) > 1 {
				myPokemon.Type2 = meta.Types[1]
			}
		} else {
			// 完全なフォールバック
			myPokemon.MaxHP = 50
			myPokemon.Attack = 10
			myPokemon.Defense = 10
			myPokemon.MoveName = "たいあたり"
			myPokemon.MovePower = 40
			myPokemon.MoveType = "normal"
		}
	}
	if myPokemon.CurrentHP <= 0 {
		myPokemon.CurrentHP = myPokemon.MaxHP
	}
	u.pokemonRepo.Save(myPokemon)

	myMeta, _ := u.metadataRepo.GetPokemonMetadata(myPokemon.PokemonNo)
	myName := "ポケモン"
	if myMeta != nil {
		myName = myMeta.JapaneseName
	}

	// 敵ポケモンの生成 (ランダム選出)
	rand.Seed(time.Now().UnixNano())
	enemyNo := rand.Intn(151) + 1
	meta, err := u.metadataRepo.GetPokemonMetadata(enemyNo)
	if err != nil {
		return nil, err
	}

	enemy := &model.Pokemon{
		PokemonNo: enemyNo,
		Level:     myPokemon.Level + rand.Intn(3) - 1, // 自分と同じくらいのレベル
		MaxHP:     meta.Stats["hp"],
		CurrentHP: meta.Stats["hp"],
		Attack:    meta.Stats["attack"],
		Defense:   meta.Stats["defense"],
		MoveName:  meta.Moves[0].Name,
		MoveType:  meta.Moves[0].Type,
		MovePower: meta.Moves[0].Power,
		Type1:     meta.Types[0],
	}
	if enemy.Level < 1 {
		enemy.Level = 1
	}
	if len(meta.Types) > 1 {
		enemy.Type2 = meta.Types[1]
	}

	return &BattleResult{
		MyPokemon:      myPokemon,
		EnemyPokemon:   enemy,
		MyName:         myName,
		EnemyName:      meta.JapaneseName,
		MyHPPercent:    myPokemon.CurrentHP * 100 / myPokemon.MaxHP,
		EnemyHPPercent: 100,
		Logs:           []string{fmt.Sprintf("野生の %s が あらわれた！", meta.JapaneseName)},
		IsFinished:     false,
	}, nil
}

func (u *battleUseCase) ExecuteTurn(userID uint, myPokemonID uint, action string, myName string, enemyName string, enemyNo int, enemyLevel int, enemyHP int, enemyMaxHP int, enemyAtk int, enemyDef int, enemyMoveName string, enemyMoveType string, enemyMovePower int, enemyType1 string, enemyType2 string) (*BattleResult, error) {
	myPokemon, err := u.pokemonRepo.FindByID(myPokemonID)
	if err != nil {
		return nil, err
	}
	
	enemy := &model.Pokemon{
		PokemonNo: enemyNo,
		Level:     enemyLevel,
		CurrentHP: enemyHP,
		MaxHP:     enemyMaxHP,
		Attack:    enemyAtk,
		Defense:   enemyDef,
		MoveName:  enemyMoveName,
		MoveType:  enemyMoveType,
		MovePower: enemyMovePower,
		Type1:     enemyType1,
		Type2:     enemyType2,
	}

	logs := []string{}
	result := &BattleResult{
		MyPokemon:    myPokemon,
		EnemyPokemon: enemy,
		MyName:       myName,
		EnemyName:    enemyName,
	}

	// 1. プレイヤーのターン
	if action == "attack" {
		damage, _, msg := u.battleService.CalculateDamage(myPokemon, enemy)
		enemy.CurrentHP -= damage
		if enemy.CurrentHP < 0 {
			enemy.CurrentHP = 0
		}
		logs = append(logs, fmt.Sprintf("%s の %s！", myName, myPokemon.MoveName))
		if msg != "" {
			logs = append(logs, msg)
		}
		logs = append(logs, fmt.Sprintf("%s に %d のダメージ！", enemyName, damage))
	}

	if enemy.CurrentHP <= 0 {
		logs = append(logs, fmt.Sprintf("%s は たおれた！", enemyName))
		result.IsFinished = true
		result.Winner = "player"
		
		// 经验値獲得
		leveledUp, growthMsgs := u.experienceService.AddExperience(myPokemon, 15) // 固定値
		result.LeveledUp = leveledUp
		logs = append(logs, fmt.Sprintf("%s は 15 の経験値を獲得した！", myName))
		logs = append(logs, growthMsgs...)

		u.pokemonRepo.Save(myPokemon)
	} else {
		// 2. 敵のターン
		damage, _, msg := u.battleService.CalculateDamage(enemy, myPokemon)
		myPokemon.CurrentHP -= damage
		if myPokemon.CurrentHP < 0 {
			myPokemon.CurrentHP = 0
		}
		logs = append(logs, fmt.Sprintf("%s の %s！", enemyName, enemy.MoveName))
		if msg != "" {
			logs = append(logs, msg)
		}
		logs = append(logs, fmt.Sprintf("%s は %d のダメージをうけた！", myName, damage))

		if myPokemon.CurrentHP <= 0 {
			logs = append(logs, fmt.Sprintf("%s は たおれた...", myName))
			result.IsFinished = true
			result.Winner = "npc"
		}
	}

	// ターンの終わりに自分のポケモンの状態を保存（HPなど）
	u.pokemonRepo.Save(myPokemon)
	
	result.MyHPPercent = myPokemon.CurrentHP * 100 / myPokemon.MaxHP
	result.EnemyHPPercent = enemy.CurrentHP * 100 / enemy.MaxHP
	result.Logs = logs
	return result, nil
}
