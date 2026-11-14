package tasks

import "github.com/gin-gonic/gin"

type DictionaryHandler struct {
	repo *DictionaryRepo
}

func NewDictionaryHandler(repo *DictionaryRepo) *DictionaryHandler {
	return &DictionaryHandler{repo: repo}
}

func (h *DictionaryHandler) listByType(dictionaryType string) gin.HandlerFunc {
	return func(c *gin.Context) {
		items, err := h.repo.ListByType(c.Request.Context(), dictionaryType)
		if err != nil {
			respondError(c, statusCodeForError(err), err)
			return
		}

		c.JSON(200, gin.H{"items": items})
	}
}

func (h *DictionaryHandler) ApplicantCategories(c *gin.Context) {
	h.listByType(DictionaryApplicantCategory)(c)
}
func (h *DictionaryHandler) Topics(c *gin.Context)      { h.listByType(DictionaryTopic)(c) }
func (h *DictionaryHandler) AppealKinds(c *gin.Context) { h.listByType(DictionaryAppealKind)(c) }
func (h *DictionaryHandler) AppealTypes(c *gin.Context) { h.listByType(DictionaryAppealType)(c) }
func (h *DictionaryHandler) PhysicalStatuses(c *gin.Context) {
	h.listByType(DictionaryPhysicalStatus)(c)
}
func (h *DictionaryHandler) SocialCategories(c *gin.Context) {
	h.listByType(DictionarySocialCategory)(c)
}
func (h *DictionaryHandler) Citizenships(c *gin.Context) { h.listByType(DictionaryCitizenship)(c) }

func (h *DictionaryHandler) AppealStatuses(c *gin.Context) {
	c.JSON(200, gin.H{"items": ListAppealStatuses()})
}
