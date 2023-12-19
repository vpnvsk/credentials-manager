package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/vpnvsk/p_s/internal/models"
	"net/http"
)

// @Summary Add new password
// @Security ApiKeyAuth
// @Tags Manager
// @Description Add new password
// @ID create-password
// @Accept  json
// @Produce  json
// @Param input body models.Credentials true "credentials info"
// @Success 200 {uuid} uuid.UUID "uuid"
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/ps [post]
func (h *Handler) createCredentials(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	var input models.Credentials
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Credentials.CreateCredentials(userId, input)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

type getAllCredentialsResponse struct {
	Data []models.CredentialsList `json:"data"`
}

// @Summary Get all credentials
// @Security ApiKeyAuth
// @Tags Manager
// @Description Get all credentials
// @ID get-all-password
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllCredentialsResponse
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/ps [get]
func (h *Handler) getAllCredentials(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	list, err := h.services.Credentials.GetAllCredentials(userId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllCredentialsResponse{
		Data: list,
	})

}

// @Summary Get Credentials By Id
// @Security ApiKeyAuth
// @Tags Manager
// @Description get password by id
// @ID get-password-by-id
// @Accept  json
// @Produce  json
// @Success 200 {object} models.CredentialsItemGet
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/ps/:id [get]
func (h *Handler) getCredentialsById(c *gin.Context) {
	id := c.Param("id")

	psId, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	item, err := h.services.Credentials.GetCredentialsByID(userId, psId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

// @Summary Update credentials
// @Security ApiKeyAuth
// @Tags Manager
// @Description Update password
// @ID update-credentials
// @Accept  json
// @Produce  json
// @Param input body models.CredentialsItemUpdate true "credentials info"
// @Success 200 {string} string ok
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/ps/:id [put]
func (h *Handler) updateCredentials(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id := c.Param("id")

	psId, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}
	var input models.CredentialsItemUpdate
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.services.UpdateCredentials(userId, psId, input); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, statusResponse{"ok"})
}

// @Summary Delete credentials
// @Security ApiKeyAuth
// @Tags Manager
// @Description Delete password
// @ID delete-credentials
// @Accept  json
// @Produce  json
// @Success 200 {string} string ok
// @Failure 400,403,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /api/ps/:id [delete]
func (h *Handler) deleteCredentials(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	id := c.Param("id")

	psId, err := uuid.Parse(id)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	err = h.services.DeleteCredentials(userId, psId)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, statusResponse{"ok"})
}
