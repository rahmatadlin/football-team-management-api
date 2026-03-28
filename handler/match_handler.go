package handler

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"football-team-management-api/models"
	"football-team-management-api/service"
	"football-team-management-api/utils/response"
)

// go-playground/validator v10 tidak menyediakan tag "matches"; validasi manual.
var matchTimeHHMMSS = regexp.MustCompile(`^([0-1][0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)

type MatchHandler struct {
	svc service.MatchService
}

func NewMatchHandler(svc service.MatchService) *MatchHandler {
	return &MatchHandler{svc: svc}
}

type createScheduleRequest struct {
	MatchDate   string `json:"match_date" binding:"required"` // YYYY-MM-DD
	MatchTime   string `json:"match_time" binding:"required"`
	HomeTeamID  string `json:"home_team_id" binding:"required,uuid"`
	AwayTeamID  string `json:"away_team_id" binding:"required,uuid"`
}

type matchIDURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type submitResultRequest struct {
	HomeScore int                 `json:"home_score" binding:"required,min=0"`
	AwayScore int                 `json:"away_score" binding:"required,min=0"`
	Goals     []service.GoalInput `json:"goals" binding:"dive"`
}

func parseDate(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02", s, time.Local)
}

func (h *MatchHandler) CreateSchedule(c *gin.Context) {
	var req createScheduleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	d, err := parseDate(req.MatchDate)
	if err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "match_date must be YYYY-MM-DD")
		return
	}
	if !matchTimeHHMMSS.MatchString(req.MatchTime) {
		response.ErrorMessage(c, http.StatusBadRequest, "match_time must be HH:MM:SS (24h)")
		return
	}
	homeID, _ := uuid.Parse(req.HomeTeamID)
	awayID, _ := uuid.Parse(req.AwayTeamID)
	m := &models.Match{
		MatchDate:  d,
		MatchTime:  req.MatchTime,
		HomeTeamID: homeID,
		AwayTeamID: awayID,
	}
	if err := h.svc.CreateSchedule(c.Request.Context(), m); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "match scheduled", m)
}

func (h *MatchHandler) ListSchedules(c *gin.Context) {
	list, err := h.svc.ListSchedules(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "ok", list)
}

func (h *MatchHandler) SubmitResult(c *gin.Context) {
	var uri matchIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid match id")
		return
	}
	matchID, _ := uuid.Parse(uri.ID)
	var req submitResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	if err := h.svc.SubmitResult(c.Request.Context(), matchID, req.HomeScore, req.AwayScore, req.Goals); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "match result saved", gin.H{
		"match_id":    matchID,
		"home_score":  req.HomeScore,
		"away_score":  req.AwayScore,
		"goals_count": len(req.Goals),
	})
}
