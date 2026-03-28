package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"football-team-management-api/models"
	"football-team-management-api/service"
	"football-team-management-api/utils/response"
)

type TeamHandler struct {
	svc service.TeamService
}

func NewTeamHandler(svc service.TeamService) *TeamHandler {
	return &TeamHandler{svc: svc}
}

type createTeamRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	LogoURL     string `json:"logo_url" binding:"max=512"`
	FoundedYear int    `json:"founded_year" binding:"required,min=1800,max=2100"`
	Address     string `json:"address" binding:"max=512"`
	City        string `json:"city" binding:"max=128"`
}

type updateTeamRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=255"`
	LogoURL     *string `json:"logo_url" binding:"omitempty,max=512"`
	FoundedYear *int    `json:"founded_year" binding:"omitempty,min=1800,max=2100"`
	Address     *string `json:"address" binding:"omitempty,max=512"`
	City        *string `json:"city" binding:"omitempty,max=128"`
}

type idURI struct {
	ID string `uri:"id" binding:"required,uuid"`
}

func (h *TeamHandler) Create(c *gin.Context) {
	var req createTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	t := &models.Team{
		Name:        req.Name,
		LogoURL:     req.LogoURL,
		FoundedYear: req.FoundedYear,
		Address:     req.Address,
		City:        req.City,
	}
	if err := h.svc.Create(c.Request.Context(), t); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusCreated, "team created", t)
}

func (h *TeamHandler) Update(c *gin.Context) {
	var uri idURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}
	id, _ := uuid.Parse(uri.ID)
	var req updateTeamRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}
	patch := &service.TeamUpdateInput{
		Name:        req.Name,
		LogoURL:     req.LogoURL,
		FoundedYear: req.FoundedYear,
		Address:     req.Address,
		City:        req.City,
	}
	t, err := h.svc.Update(c.Request.Context(), id, patch)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "team updated", t)
}

func (h *TeamHandler) Delete(c *gin.Context) {
	var uri idURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}
	id, _ := uuid.Parse(uri.ID)
	if err := h.svc.SoftDelete(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "team deleted", nil)
}

func (h *TeamHandler) Get(c *gin.Context) {
	var uri idURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid id")
		return
	}
	id, _ := uuid.Parse(uri.ID)
	t, err := h.svc.Get(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "ok", t)
}

func (h *TeamHandler) List(c *gin.Context) {
	list, err := h.svc.List(c.Request.Context())
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "ok", list)
}
