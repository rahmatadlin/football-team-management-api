package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"football-team-management-api/models"
	"football-team-management-api/service"
	"football-team-management-api/utils/response"
)

// Respons khusus GET .../teams/:id/players: tanpa team_id (sudah di URL) dan tanpa objek team.
type playerByTeamItem struct {
	ID           uuid.UUID       `json:"id"`
	Name         string          `json:"name"`
	Height       float64         `json:"height"`
	Weight       float64         `json:"weight"`
	Position     models.Position `json:"position"`
	JerseyNumber int             `json:"jersey_number"`
	DeletedAt    gorm.DeletedAt  `json:"deleted_at,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func toPlayerByTeamItems(players []models.Player) []playerByTeamItem {
	out := make([]playerByTeamItem, 0, len(players))
	for _, p := range players {
		out = append(out, playerByTeamItem{
			ID:           p.ID,
			Name:         p.Name,
			Height:       p.Height,
			Weight:       p.Weight,
			Position:     p.Position,
			JerseyNumber: p.JerseyNumber,
			DeletedAt:    p.DeletedAt,
			CreatedAt:    p.CreatedAt,
			UpdatedAt:    p.UpdatedAt,
		})
	}
	return out
}

type PlayerHandler struct {
	svc service.PlayerService
}

func NewPlayerHandler(svc service.PlayerService) *PlayerHandler {
	return &PlayerHandler{svc: svc}
}

type createPlayerRequest struct {
	Name         string          `json:"name" binding:"required,min=1,max=255"`
	Height       float64         `json:"height" binding:"required,gt=0"`
	Weight       float64         `json:"weight" binding:"required,gt=0"`
	Position     models.Position `json:"position" binding:"required,oneof=striker midfielder defender goalkeeper"`
	JerseyNumber int             `json:"jersey_number" binding:"required,min=1,max=99"`
}

type updatePlayerRequest struct {
	Name         *string          `json:"name" binding:"omitempty,min=1,max=255"`
	Height       *float64         `json:"height" binding:"omitempty,gt=0"`
	Weight       *float64         `json:"weight" binding:"omitempty,gt=0"`
	Position     *models.Position `json:"position" binding:"omitempty,oneof=striker midfielder defender goalkeeper"`
	JerseyNumber *int             `json:"jersey_number" binding:"omitempty,min=1,max=99"`
	TeamID       *string          `json:"team_id" binding:"omitempty,uuid"`
}

type teamPlayersURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

type playerIDURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (h *PlayerHandler) Create(c *gin.Context) {
	var turi teamPlayersURI
	if err := c.ShouldBindUri(&turi); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid team id")
		return
	}
	teamID, _ := uuid.Parse(turi.ID)
	var req createPlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	p := &models.Player{
		TeamID:       teamID,
		Name:         req.Name,
		Height:       req.Height,
		Weight:       req.Weight,
		Position:     req.Position,
		JerseyNumber: req.JerseyNumber,
	}
	if err := h.svc.Create(c.Request.Context(), p); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "player created", p)
}

func (h *PlayerHandler) Update(c *gin.Context) {
	var uri playerIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}
	id, _ := uuid.Parse(uri.ID)
	var req updatePlayerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	patch := &models.Player{}
	if req.Name != nil {
		patch.Name = *req.Name
	}
	if req.Height != nil {
		patch.Height = *req.Height
	}
	if req.Weight != nil {
		patch.Weight = *req.Weight
	}
	if req.Position != nil {
		patch.Position = *req.Position
	}
	if req.JerseyNumber != nil {
		patch.JerseyNumber = *req.JerseyNumber
	}
	if req.TeamID != nil {
		tid, err := uuid.Parse(*req.TeamID)
		if err != nil {
			response.ErrorMessage(c, http.StatusBadRequest, "invalid team_id")
			return
		}
		patch.TeamID = tid
	}
	out, err := h.svc.Update(c.Request.Context(), id, patch)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "player updated", out)
}

func (h *PlayerHandler) Delete(c *gin.Context) {
	var uri playerIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}
	id, _ := uuid.Parse(uri.ID)
	if err := h.svc.SoftDelete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "player deleted", nil)
}

func (h *PlayerHandler) ListByTeam(c *gin.Context) {
	var turi teamPlayersURI
	if err := c.ShouldBindUri(&turi); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid team id")
		return
	}
	teamID, _ := uuid.Parse(turi.ID)
	list, err := h.svc.ListByTeam(c.Request.Context(), teamID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "ok", toPlayerByTeamItems(list))
}
