package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
	"football-team-management-api/repository"
	"football-team-management-api/utils/apperror"
)

type MatchOutcome string

const (
	OutcomeHomeWin MatchOutcome = "HOME_WIN"
	OutcomeAwayWin MatchOutcome = "AWAY_WIN"
	OutcomeDraw    MatchOutcome = "DRAW"
)

type TopScorer struct {
	PlayerID   uuid.UUID `json:"player_id"`
	PlayerName string    `json:"player_name"`
	Goals      int       `json:"goals"`
}

type MatchReport struct {
	Schedule struct {
		MatchDate string `json:"match_date"`
		MatchTime string `json:"match_time"`
	} `json:"schedule"`
	HomeTeam          models.Team    `json:"home_team"`
	AwayTeam          models.Team    `json:"away_team"`
	HomeScore         int            `json:"home_score"`
	AwayScore         int            `json:"away_score"`
	Result            MatchOutcome   `json:"match_result"`
	TopScorer         *TopScorer     `json:"top_scorer"`
	HomeTeamWinsUntil int            `json:"home_team_wins_until_match"`
	AwayTeamWinsUntil int            `json:"away_team_wins_until_match"`
}

type ReportService interface {
	MatchReport(ctx context.Context, matchID uuid.UUID) (*MatchReport, error)
}

type reportService struct {
	matches repository.MatchRepository
	goals   repository.GoalRepository
}

func NewReportService(matches repository.MatchRepository, goals repository.GoalRepository) ReportService {
	return &reportService{matches: matches, goals: goals}
}

func outcome(home, away int) MatchOutcome {
	if home > away {
		return OutcomeHomeWin
	}
	if away > home {
		return OutcomeAwayWin
	}
	return OutcomeDraw
}

func countWinsForTeam(teamID uuid.UUID, past []models.Match) int {
	wins := 0
	for _, m := range past {
		if m.Result == nil {
			continue
		}
		if m.HomeTeamID == teamID && m.Result.HomeScore > m.Result.AwayScore {
			wins++
		}
		if m.AwayTeamID == teamID && m.Result.AwayScore > m.Result.HomeScore {
			wins++
		}
	}
	return wins
}

func (s *reportService) MatchReport(ctx context.Context, matchID uuid.UUID) (*MatchReport, error) {
	m, err := s.matches.FindByID(ctx, matchID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.NotFound("match not found")
		}
		return nil, err
	}
	if m.Result == nil {
		return nil, apperror.BadRequest("match has no result yet")
	}
	goals, err := s.goals.GoalsByMatch(ctx, matchID)
	if err != nil {
		return nil, err
	}
	counts := map[uuid.UUID]int{}
	names := map[uuid.UUID]string{}
	for _, g := range goals {
		counts[g.PlayerID]++
		if g.Player.ID != uuid.Nil {
			names[g.PlayerID] = g.Player.Name
		}
	}
	var top *TopScorer
	for pid, c := range counts {
		name := names[pid]
		if top == nil {
			top = &TopScorer{PlayerID: pid, PlayerName: name, Goals: c}
			continue
		}
		if c > top.Goals {
			top = &TopScorer{PlayerID: pid, PlayerName: name, Goals: c}
		} else if c == top.Goals && pid.String() < top.PlayerID.String() {
			top = &TopScorer{PlayerID: pid, PlayerName: name, Goals: c}
		}
	}
	if len(counts) == 0 {
		top = nil
	}

	dateStr := m.MatchDate.Format("2006-01-02")
	prior, err := s.matches.FindCompletedBefore(ctx, dateStr, m.MatchTime, m.ID)
	if err != nil {
		return nil, err
	}

	rep := &MatchReport{
		HomeTeam:          m.HomeTeam,
		AwayTeam:          m.AwayTeam,
		HomeScore:         m.Result.HomeScore,
		AwayScore:         m.Result.AwayScore,
		Result:            outcome(m.Result.HomeScore, m.Result.AwayScore),
		TopScorer:         top,
		HomeTeamWinsUntil: countWinsForTeam(m.HomeTeamID, prior),
		AwayTeamWinsUntil: countWinsForTeam(m.AwayTeamID, prior),
	}
	rep.Schedule.MatchDate = dateStr
	rep.Schedule.MatchTime = m.MatchTime

	return rep, nil
}
