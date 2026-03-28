package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"football-team-management-api/service"
	"football-team-management-api/utils/response"
)

type ReportHandler struct {
	svc service.ReportService
}

func NewReportHandler(svc service.ReportService) *ReportHandler {
	return &ReportHandler{svc: svc}
}

func (h *ReportHandler) MatchReport(c *gin.Context) {
	var uri matchIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		response.ErrorMessage(c, http.StatusBadRequest, "invalid match id")
		return
	}
	matchID, _ := uuid.Parse(uri.ID)
	rep, err := h.svc.MatchReport(c.Request.Context(), matchID)
	if err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, http.StatusOK, "ok", rep)
}
