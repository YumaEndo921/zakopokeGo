package service

import (
	"zakopokeGo/internal/domain/model"
)

type BattleService interface {
	CalculateDamage(attacker *model.Pokemon, defender *model.Pokemon) (int, float64, string)
	GetTypeEffectiveness(attackType string, targetTypes []string) (float64, string)
}

type battleService struct{}

func NewBattleService() BattleService {
	return &battleService{}
}

func (s *battleService) CalculateDamage(attacker *model.Pokemon, defender *model.Pokemon) (int, float64, string) {
	defenderTypes := []string{defender.Type1}
	if defender.Type2 != "" {
		defenderTypes = append(defenderTypes, defender.Type2)
	}
	effectiveness, message := s.GetTypeEffectiveness(attacker.MoveType, defenderTypes)

	// ダメージ計算式
	baseDamage := (float64(attacker.Attack) * float64(attacker.MovePower) / float64(defender.Defense))
	if baseDamage < 1 {
		baseDamage = 1
	}
	
	finalDamage := int(baseDamage * effectiveness)
	if finalDamage < 1 {
		finalDamage = 1
	}

	return finalDamage, effectiveness, message
}

func (s *battleService) GetTypeEffectiveness(attackType string, targetTypes []string) (float64, string) {
	// 簡易相性表
	multipliers := map[string]map[string]float64{
		"fire": {
			"grass": 2.0,
			"water": 0.5,
			"fire":  0.5,
		},
		"water": {
			"fire":  2.0,
			"grass": 0.5,
			"water": 0.5,
		},
		"grass": {
			"water": 2.0,
			"fire":  0.5,
			"grass": 0.5,
		},
		"electric": {
			"water": 2.0,
			"grass": 0.5,
		},
	}

	totalEff := 1.0
	for _, t := range targetTypes {
		if m, ok := multipliers[attackType][t]; ok {
			totalEff *= m
		}
	}

	message := ""
	if totalEff > 1.0 {
		message = "こうかは ばつぐんだ！"
	} else if totalEff < 1.0 {
		message = "こうかは いまいちの ようだ..."
	}

	return totalEff, message
}
