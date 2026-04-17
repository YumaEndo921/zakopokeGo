package service

import (
	"fmt"
	"zakopokeGo/internal/domain/model"
)

type ExperienceService interface {
	AddExperience(p *model.Pokemon, exp int) (bool, []string)
}

type experienceService struct{}

func NewExperienceService() ExperienceService {
	return &experienceService{}
}

func (s *experienceService) AddExperience(p *model.Pokemon, exp int) (bool, []string) {
	p.Exp += exp
	messages := []string{}
	leveledUp := false

	// シンプルなレベルアップ必要経験値: level * 10
	for p.Exp >= p.Level*10 {
		p.Exp -= p.Level * 10
		p.Level++
		leveledUp = true

		// ステータス上昇
		hpGain := 2
		atkGain := 1
		defGain := 1
		
		p.MaxHP += hpGain
		p.CurrentHP += hpGain
		p.Attack += atkGain
		p.Defense += defGain

		messages = append(messages, fmt.Sprintf("レベル %d にレベルアップした！", p.Level))
	}

	return leveledUp, messages
}
