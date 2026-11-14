package tasks

import (
	"context"
	"mime/multipart"
	"strings"
)

type AppealService struct {
	repo         *AppealRepo
	dictionaries *DictionaryRepo
}

func NewAppealService(repo *AppealRepo, dictionaries *DictionaryRepo) *AppealService {
	return &AppealService{repo: repo, dictionaries: dictionaries}
}

func (s *AppealService) Create(ctx context.Context, user *UserContext, input CreateAppealInput, files []*multipart.FileHeader) (int64, error) {
	if errs := ValidateStruct(input); errs != nil {
		return 0, validationError(errs)
	}

	if err := s.validateCreateInput(ctx, input); err != nil {
		return 0, err
	}

	savedFiles, err := saveAppealFiles(files)
	if err != nil {
		return 0, err
	}

	return s.repo.Create(ctx, user.UserID, input, savedFiles)
}

func (s *AppealService) validateCreateInput(ctx context.Context, input CreateAppealInput) error {
	checks := []struct {
		dictionaryType string
		code           string
	}{
		{DictionaryApplicantCategory, input.ApplicantCategoryCode},
		{DictionaryTopic, input.TopicCode},
		{DictionaryAppealKind, input.AppealKindCode},
		{DictionaryAppealType, input.AppealTypeCode},
	}

	for _, check := range checks {
		exists, err := s.dictionaries.Exists(ctx, check.dictionaryType, check.code)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidField
		}
	}

	if input.PhysicalStatusCode != nil {
		exists, err := s.dictionaries.Exists(ctx, DictionaryPhysicalStatus, *input.PhysicalStatusCode)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidField
		}
	}
	if input.SocialCategoryCode != nil {
		exists, err := s.dictionaries.Exists(ctx, DictionarySocialCategory, *input.SocialCategoryCode)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidField
		}
	}
	if input.CitizenshipCode != nil {
		exists, err := s.dictionaries.Exists(ctx, DictionaryCitizenship, *input.CitizenshipCode)
		if err != nil {
			return err
		}
		if !exists {
			return ErrInvalidField
		}
	}

	if input.ApplicantCategoryCode == "legal" && input.OrganizationName == nil {
		return ErrRequiredIfField
	}
	if input.ApplicantCategoryCode != "legal" {
		input.OrganizationName = nil
	}

	if input.IsAbroad {
		if strings.TrimSpace(input.Locality) == "" {
			return ErrRequiredField
		}
	} else if strings.TrimSpace(input.Locality) == "" || strings.TrimSpace(input.Street) == "" || strings.TrimSpace(input.House) == "" {
		return ErrRequiredField
	}

	return nil
}

func (s *AppealService) List(ctx context.Context, ownerUserID *int64, filter AppealListFilter) ([]AppealListItemResponse, int64, error) {
	return s.repo.List(ctx, ownerUserID, filter)
}

func (s *AppealService) Get(ctx context.Context, appealID int64, currentUser *UserContext, adminView bool) (*AppealDetailResponse, error) {
	detail, err := s.repo.GetByID(ctx, appealID)
	if err != nil {
		return nil, err
	}
	if !adminView && detail.SubmittedByUserID != currentUser.UserID {
		return nil, ErrForbidden
	}
	detail.AllowedStatusTransitions = AllowedTransitions(detail.Status.Code, adminView && currentUser.Role == RoleAdmin)
	return detail, nil
}

func (s *AppealService) ChangeStatus(ctx context.Context, appealID int64, toStatus string, comment *string, currentUser *UserContext) error {
	detail, err := s.repo.GetByID(ctx, appealID)
	if err != nil {
		return err
	}

	if !IsAllowedTransition(detail.Status.Code, toStatus, currentUser.Role == RoleAdmin) {
		return ErrInvalidStatusTransition
	}

	_, ok := GetAppealStatus(toStatus)
	if !ok {
		return ErrInvalidField
	}

	return s.repo.ChangeStatus(ctx, appealID, detail.Status.Code, toStatus, comment, currentUser.UserID)
}

type fieldValidationError map[string]string

func (e fieldValidationError) Error() string { return "validation error" }

func validationError(errs map[string]string) error {
	return fieldValidationError(errs)
}
