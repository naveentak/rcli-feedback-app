package feedback

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rcli/feedback/internal/auth"
)

type SignedSubmitRequest struct {
	DeviceID    string `json:"deviceId" binding:"required"`
	Timestamp   string `json:"timestamp" binding:"required"`
	Signature   string `json:"signature" binding:"required"`
	AppID       string `json:"appId" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Email       string `json:"email,omitempty"`
	AppVersion  string `json:"appVersion,omitempty"`
	OSVersion   string `json:"osVersion,omitempty"`
}

func (h *Handler) SignedSubmit(hmacSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if hmacSecret == "" {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": "signed feedback not configured"})
			return
		}

		var req SignedSubmitRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := auth.ValidateFeedbackHMAC(req.DeviceID, req.Timestamp, req.Signature, hmacSecret); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		ticketType := req.Type
		if ticketType == "feature" {
			ticketType = string(TypeFeatureRequest)
		}
		if !IsValidType(ticketType) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid type: %s", req.Type)})
			return
		}
		if !IsValidApp(req.AppID) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("invalid app: %s", req.AppID)})
			return
		}

		body := req.Description
		if req.AppVersion != "" || req.OSVersion != "" || req.DeviceID != "" {
			body += "\n\n---\n"
			if req.AppVersion != "" {
				body += fmt.Sprintf("App version: %s\n", req.AppVersion)
			}
			if req.OSVersion != "" {
				body += fmt.Sprintf("macOS: %s\n", req.OSVersion)
			}
			body += fmt.Sprintf("Device ID: %s\n", req.DeviceID)
		}

		ticket, err := h.svc.Submit(c.Request.Context(), SubmitRequest{
			App:         App(req.AppID),
			Type:        TicketType(ticketType),
			Title:       req.Title,
			Description: body,
			Reporter:    req.Email,
		})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, ticket)
	}
}