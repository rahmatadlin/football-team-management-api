package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Match is a scheduled match between two teams.
type Match struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	MatchDate   time.Time `json:"match_date" gorm:"type:date;not null;index"`
	MatchTime   string    `json:"match_time" gorm:"size:8;not null"` // "HH:MM:SS" or "15:04:05"
	HomeTeamID  uuid.UUID `json:"home_team_id" gorm:"type:char(36);not null;index"`
	AwayTeamID  uuid.UUID `json:"away_team_id" gorm:"type:char(36);not null;index"`
	HomeTeam    Team      `json:"home_team,omitempty" gorm:"foreignKey:HomeTeamID"`
	AwayTeam    Team      `json:"away_team,omitempty" gorm:"foreignKey:AwayTeamID"`
	Result      *MatchResult `json:"result,omitempty" gorm:"foreignKey:MatchID"`
	Goals       []Goal       `json:"goals,omitempty" gorm:"foreignKey:MatchID"`
	SoftDeleteModel
	Timestamps
}

func (m *Match) BeforeCreate(tx *gorm.DB) error {
	return BeforeCreateUUID(&m.ID)
}

func (Match) TableName() string {
	return "matches"
}
