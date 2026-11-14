package tasks

import (
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AppealHandler struct {
	service *AppealService
	repo    *AppealRepo
}

type ChangeStatusRequest struct {
	StatusCode string  `json:"status_code" validate:"required"`
	Comment    *string `json:"comment"`
}

func NewAppealHandler(service *AppealService, repo *AppealRepo) *AppealHandler {
	return &AppealHandler{service: service, repo: repo}
}

func (h *AppealHandler) Create(c *gin.Context) {
	user, ok := CurrentUser(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, runtimeSettings.MaxBodySizeBytes)
	if err := c.Request.ParseMultipartForm(runtimeSettings.MaxBodySizeBytes); err != nil {
		respondError(c, http.StatusRequestEntityTooLarge, ErrRequestTooLarge)
		return
	}

	input, err := bindCreateAppealInput(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if err := VerifyCaptchaToken(c.Request.Context(), c.PostForm("captcha_token"), c.ClientIP()); err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	files := make([]*multipart.FileHeader, 0)
	if c.Request.MultipartForm != nil {
		files = c.Request.MultipartForm.File["files"]
	}

	id, err := h.service.Create(c.Request.Context(), user, input, files)
	if err != nil {
		var fieldErr fieldValidationError
		if errors.As(err, &fieldErr) {
			respondValidation(c, map[string]string(fieldErr))
			return
		}
		respondError(c, statusCodeForError(err), err)
		return
	}

	detail, err := h.service.Get(c.Request.Context(), id, user, user.Role == RoleAdmin)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	c.JSON(http.StatusCreated, detail)
}

func (h *AppealHandler) ListMine(c *gin.Context) {
	user, ok := CurrentUser(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	filter, err := parseAppealListFilter(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	items, total, err := h.service.List(c.Request.Context(), &user.UserID, filter)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	respondPaginated(c, items, total, filter.Page, filter.PageSize)
}

func (h *AppealHandler) ListAdmin(c *gin.Context) {
	filter, err := parseAppealListFilter(c)
	if err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	items, total, err := h.service.List(c.Request.Context(), nil, filter)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	respondPaginated(c, items, total, filter.Page, filter.PageSize)
}

func (h *AppealHandler) GetMine(c *gin.Context) {
	h.getOne(c, false)
}

func (h *AppealHandler) GetAdmin(c *gin.Context) {
	h.getOne(c, true)
}

func (h *AppealHandler) ChangeStatus(c *gin.Context) {
	user, ok := CurrentUser(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	appealID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidField)
		return
	}

	var req ChangeStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, http.StatusBadRequest, err)
		return
	}

	if errs := ValidateStruct(req); errs != nil {
		respondValidation(c, errs)
		return
	}

	if err := h.service.ChangeStatus(c.Request.Context(), appealID, req.StatusCode, req.Comment, user); err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	detail, err := h.service.Get(c.Request.Context(), appealID, user, true)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	c.JSON(http.StatusOK, detail)
}

func (h *AppealHandler) DownloadFile(c *gin.Context) {
	fileID, err := strconv.ParseInt(c.Param("fileID"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidField)
		return
	}

	file, err := h.repo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	if _, err := os.Stat(file.StoragePath); err != nil {
		respondError(c, http.StatusNotFound, ErrNotFound)
		return
	}

	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", `attachment; filename="`+filepath.Base(file.OriginalName)+`"`)
	c.File(file.StoragePath)
}

// DownloadFileMine отдаёт файл владельцу обращения (или админу) с проверкой прав.
func (h *AppealHandler) DownloadFileMine(c *gin.Context) {
	user, ok := CurrentUser(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	appealID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidField)
		return
	}
	fileID, err := strconv.ParseInt(c.Param("fileID"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidField)
		return
	}

	file, err := h.repo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}
	if file.AppealID != appealID {
		respondError(c, http.StatusNotFound, ErrNotFound)
		return
	}

	ownerID, err := h.repo.GetOwnerID(c.Request.Context(), appealID)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}
	if user.Role != RoleAdmin && ownerID != user.UserID {
		respondError(c, http.StatusForbidden, ErrForbidden)
		return
	}

	if _, err := os.Stat(file.StoragePath); err != nil {
		respondError(c, http.StatusNotFound, ErrNotFound)
		return
	}

	c.Header("Content-Type", file.ContentType)
	c.Header("Content-Disposition", `attachment; filename="`+filepath.Base(file.OriginalName)+`"`)
	c.File(file.StoragePath)
}

func (h *AppealHandler) getOne(c *gin.Context, adminView bool) {
	user, ok := CurrentUser(c)
	if !ok {
		respondError(c, http.StatusUnauthorized, ErrUnauthorized)
		return
	}

	appealID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		respondError(c, http.StatusBadRequest, ErrInvalidField)
		return
	}

	detail, err := h.service.Get(c.Request.Context(), appealID, user, adminView)
	if err != nil {
		respondError(c, statusCodeForError(err), err)
		return
	}

	c.JSON(http.StatusOK, detail)
}

func parseAppealListFilter(c *gin.Context) (AppealListFilter, error) {
	page, _ := strconv.Atoi(defaultString(c.Query("page"), "1"))
	pageSize, _ := strconv.Atoi(defaultString(c.Query("page_size"), "20"))
	filter := AppealListFilter{
		Page:                  page,
		PageSize:              pageSize,
		Sort:                  defaultString(c.Query("sort"), "created_at_desc"),
		CitizenshipCode:       optionalString(c.Query("citizenship")),
		Region:                optionalString(c.Query("region")),
		District:              optionalString(c.Query("district")),
		Locality:              optionalString(c.Query("locality")),
		Street:                optionalString(c.Query("street")),
		House:                 optionalString(c.Query("house")),
		PostalCode:            optionalString(c.Query("postal_code")),
		Apartment:             optionalString(c.Query("apartment")),
		IdentificationNumber:  optionalString(c.Query("identification_number")),
		ApplicantName:         optionalString(c.Query("applicant_name")),
		TopicCode:             optionalString(c.Query("topic_code")),
		ApplicantCategoryCode: optionalString(c.Query("applicant_category_code")),
		AppealKindCode:        optionalString(c.Query("appeal_kind_code")),
		AppealTypeCode:        optionalString(c.Query("appeal_type_code")),
		StatusCode:            optionalString(c.Query("status_code")),
		StatusGroupCode:       optionalString(c.Query("status_group_code")),
		Number:                optionalString(c.Query("number")),
	}

	var err error
	filter.PlannedDateFrom, err = optionalDate(c.Query("planned_date_from"))
	if err != nil {
		return filter, err
	}
	filter.PlannedDateTo, err = optionalDate(c.Query("planned_date_to"))
	if err != nil {
		return filter, err
	}
	filter.CreatedAtFrom, err = optionalDateTime(c.Query("created_at_from"))
	if err != nil {
		return filter, err
	}
	filter.CreatedAtTo, err = optionalDateTime(c.Query("created_at_to"))
	if err != nil {
		return filter, err
	}

	return filter, nil
}

func bindCreateAppealInput(c *gin.Context) (CreateAppealInput, error) {
	input := CreateAppealInput{
		ApplicantCategoryCode: c.PostForm("applicant_category_code"),
		TopicCode:             c.PostForm("topic_code"),
		OrganizationName:      optionalString(c.PostForm("organization_name")),
		AppealKindCode:        c.PostForm("appeal_kind_code"),
		AppealTypeCode:        c.PostForm("appeal_type_code"),
		LastName:              c.PostForm("last_name"),
		FirstName:             c.PostForm("first_name"),
		Patronymic:            optionalString(c.PostForm("patronymic")),
		IdentificationNumber:  c.PostForm("identification_number"),
		PhysicalStatusCode:    optionalString(c.PostForm("physical_status_code")),
		SocialCategoryCode:    optionalString(c.PostForm("social_category_code")),
		Email:                 c.PostForm("email"),
		IsAbroad:              c.PostForm("is_abroad") == "true",
		CitizenshipCode:       optionalString(c.PostForm("citizenship_code")),
		Region:                optionalString(c.PostForm("region")),
		District:              optionalString(c.PostForm("district")),
		Locality:              c.PostForm("locality"),
		Street:                c.PostForm("street"),
		House:                 c.PostForm("house"),
		Apartment:             optionalString(c.PostForm("apartment")),
		NoteAddress:           optionalString(c.PostForm("note_address")),
		PostalCode:            optionalString(c.PostForm("postal_code")),
		Phone:                 c.PostForm("phone"),
		Content:               c.PostForm("content"),
	}

	if value := c.PostForm("planned_registration_date"); value != "" {
		parsed, err := time.Parse(time.RFC3339, value)
		if err != nil {
			parsed, err = time.Parse("2006-01-02", value)
			if err != nil {
				return input, ErrInvalidField
			}
		}
		input.PlannedRegistrationDate = &parsed
	}

	return input, nil
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	copy := value
	return &copy
}

func defaultString(value string, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func optionalDate(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return nil, ErrInvalidField
	}
	return &parsed, nil
}

func optionalDateTime(value string) (*time.Time, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, ErrInvalidField
	}
	return &parsed, nil
}
