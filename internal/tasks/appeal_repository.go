package tasks

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AppealRepo struct {
	pool *pgxpool.Pool
}

type AppealFileRecord struct {
	ID           int64
	AppealID     int64
	OriginalName string
	StoredName   string
	StoragePath  string
	ContentType  string
	SizeBytes    int64
}

type StatusHistoryRecord struct {
	ID             int64
	FromStatusCode *string
	ToStatusCode   string
	Comment        *string
	ChangedBy      int64
	ChangedAt      time.Time
}

func NewAppealRepo(pool *pgxpool.Pool) *AppealRepo {
	return &AppealRepo{pool: pool}
}

func (r *AppealRepo) Create(ctx context.Context, ownerUserID int64, input CreateAppealInput, files []SavedAppealFile) (int64, error) {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		cleanupSavedFiles(files)
		return 0, err
	}
	defer tx.Rollback(ctx)

	number := fmt.Sprintf("APL-%d", time.Now().UnixNano())
	var appealID int64
	err = tx.QueryRow(ctx, `
		INSERT INTO appeals (
			number, owner_user_id, applicant_category_code, topic_code, organization_name,
			appeal_kind_code, appeal_type_code, last_name, first_name, patronymic,
			identification_number, physical_status_code, social_category_code, email,
			is_abroad, citizenship_code, region, district, locality, street, house,
			apartment, note_address, postal_code, phone, content, planned_registration_date,
			current_status_code
		) VALUES (
			$1,$2,$3,$4,$5,
			$6,$7,$8,$9,$10,
			$11,$12,$13,$14,
			$15,$16,$17,$18,$19,$20,$21,
			$22,$23,$24,$25,$26,$27,
			$28
		) RETURNING id
	`,
		number, ownerUserID, input.ApplicantCategoryCode, input.TopicCode, input.OrganizationName,
		input.AppealKindCode, input.AppealTypeCode, input.LastName, input.FirstName, input.Patronymic,
		input.IdentificationNumber, input.PhysicalStatusCode, input.SocialCategoryCode, input.Email,
		input.IsAbroad, input.CitizenshipCode, input.Region, input.District, input.Locality, input.Street, input.House,
		input.Apartment, input.NoteAddress, input.PostalCode, input.Phone, input.Content, input.PlannedRegistrationDate,
		"new",
	).Scan(&appealID)
	if err != nil {
		cleanupSavedFiles(files)
		return 0, dbError(err)
	}

	for _, file := range files {
		_, err := tx.Exec(ctx, `
			INSERT INTO appeal_files (appeal_id, original_name, stored_name, storage_path, content_type, size_bytes)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, appealID, file.OriginalName, file.StoredName, file.StoragePath, file.ContentType, file.SizeBytes)
		if err != nil {
			cleanupSavedFiles(files)
			return 0, dbError(err)
		}
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO appeal_status_history (appeal_id, from_status_code, to_status_code, comment, changed_by_user_id)
		VALUES ($1, NULL, $2, NULL, $3)
	`, appealID, "new", ownerUserID)
	if err != nil {
		cleanupSavedFiles(files)
		return 0, dbError(err)
	}

	if err := tx.Commit(ctx); err != nil {
		cleanupSavedFiles(files)
		return 0, dbError(err)
	}

	return appealID, nil
}

func (r *AppealRepo) List(ctx context.Context, ownerUserID *int64, filter AppealListFilter) ([]AppealListItemResponse, int64, error) {
	where := make([]string, 0)
	args := make([]any, 0)
	arg := 1

	if ownerUserID != nil {
		where = append(where, fmt.Sprintf("a.owner_user_id = $%d", arg))
		args = append(args, *ownerUserID)
		arg++
	}

	addEqual := func(column string, value *string) {
		if value == nil || *value == "" {
			return
		}
		where = append(where, fmt.Sprintf("%s = $%d", column, arg))
		args = append(args, *value)
		arg++
	}

	addEqual("a.citizenship_code", filter.CitizenshipCode)
	addEqual("a.region", filter.Region)
	addEqual("a.district", filter.District)
	addEqual("a.locality", filter.Locality)
	addEqual("a.street", filter.Street)
	addEqual("a.house", filter.House)
	addEqual("a.postal_code", filter.PostalCode)
	addEqual("a.apartment", filter.Apartment)
	addEqual("a.identification_number", filter.IdentificationNumber)
	addEqual("a.topic_code", filter.TopicCode)
	addEqual("a.applicant_category_code", filter.ApplicantCategoryCode)
	addEqual("a.appeal_kind_code", filter.AppealKindCode)
	addEqual("a.appeal_type_code", filter.AppealTypeCode)
	addEqual("a.current_status_code", filter.StatusCode)
	addEqual("s.group_code", filter.StatusGroupCode)
	addEqual("a.number", filter.Number)

	if filter.ApplicantName != nil && *filter.ApplicantName != "" {
		where = append(where, fmt.Sprintf("concat_ws(' ', a.last_name, a.first_name, coalesce(a.patronymic, '')) = $%d", arg))
		args = append(args, *filter.ApplicantName)
		arg++
	}

	if filter.PlannedDateFrom != nil {
		where = append(where, fmt.Sprintf("a.planned_registration_date >= $%d", arg))
		args = append(args, *filter.PlannedDateFrom)
		arg++
	}
	if filter.PlannedDateTo != nil {
		where = append(where, fmt.Sprintf("a.planned_registration_date <= $%d", arg))
		args = append(args, *filter.PlannedDateTo)
		arg++
	}
	if filter.CreatedAtFrom != nil {
		where = append(where, fmt.Sprintf("a.created_at >= $%d", arg))
		args = append(args, *filter.CreatedAtFrom)
		arg++
	}
	if filter.CreatedAtTo != nil {
		where = append(where, fmt.Sprintf("a.created_at <= $%d", arg))
		args = append(args, *filter.CreatedAtTo)
		arg++
	}

	whereSQL := ""
	if len(where) > 0 {
		whereSQL = "WHERE " + strings.Join(where, " AND ")
	}

	baseJoins := `
		FROM appeals a
		JOIN dictionaries ac ON ac.dictionary_type = 'applicant_category' AND ac.code = a.applicant_category_code
		JOIN dictionaries t ON t.dictionary_type = 'topic' AND t.code = a.topic_code
		JOIN dictionaries ak ON ak.dictionary_type = 'appeal_kind' AND ak.code = a.appeal_kind_code
		JOIN dictionaries atp ON atp.dictionary_type = 'appeal_type' AND atp.code = a.appeal_type_code
		LEFT JOIN dictionaries c ON c.dictionary_type = 'citizenship' AND c.code = a.citizenship_code
		JOIN appeal_statuses s ON s.code = a.current_status_code
	`

	countSQL := `SELECT count(*) ` + baseJoins + whereSQL
	var total int64
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&total); err != nil {
		return nil, 0, dbError(err)
	}

	page := filter.Page
	if page <= 0 {
		page = 1
	}
	pageSize := filter.PageSize
	if pageSize <= 0 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	orderBy := "a.created_at DESC"
	if filter.Sort == "created_at_asc" {
		orderBy = "a.created_at ASC"
	}

	argsWithPagination := append(args, pageSize, offset)
	limitArg := arg
	offsetArg := arg + 1

	query := `
		SELECT
			a.id, a.number, a.created_at, a.planned_registration_date,
			a.applicant_category_code, ac.title,
			a.topic_code, t.title,
			a.appeal_kind_code, ak.title,
			a.appeal_type_code, atp.title,
			a.last_name, a.first_name, a.patronymic,
			a.identification_number,
			a.citizenship_code, c.title,
			a.current_status_code
		` + baseJoins + whereSQL + `
		ORDER BY ` + orderBy + `
		LIMIT $` + fmt.Sprint(limitArg) + ` OFFSET $` + fmt.Sprint(offsetArg)

	rows, err := r.pool.Query(ctx, query, argsWithPagination...)
	if err != nil {
		return nil, 0, dbError(err)
	}
	defer rows.Close()

	items := make([]AppealListItemResponse, 0)
	for rows.Next() {
		var item AppealListItemResponse
		var lastName string
		var firstName string
		var patronymic *string
		var citizenshipCode *string
		var citizenshipTitle *string
		var statusCode string

		if err := rows.Scan(
			&item.ID, &item.Number, &item.CreatedAt, &item.PlannedRegistrationDate,
			&item.ApplicantCategory.Code, &item.ApplicantCategory.Title,
			&item.Topic.Code, &item.Topic.Title,
			&item.AppealKind.Code, &item.AppealKind.Title,
			&item.AppealType.Code, &item.AppealType.Title,
			&lastName, &firstName, &patronymic,
			&item.IdentificationNumber,
			&citizenshipCode, &citizenshipTitle,
			&statusCode,
		); err != nil {
			return nil, 0, dbError(err)
		}

		item.ApplicantName = composeApplicantName(lastName, firstName, patronymic)
		if citizenshipCode != nil && citizenshipTitle != nil {
			item.Citizenship = &DictionaryItem{Code: *citizenshipCode, Title: *citizenshipTitle}
		}
		if status, ok := GetAppealStatus(statusCode); ok {
			item.Status = status
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, dbError(err)
	}

	return items, total, nil
}

func (r *AppealRepo) GetByID(ctx context.Context, id int64) (*AppealDetailResponse, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT
			a.id, a.number, a.created_at, a.updated_at, a.owner_user_id,
			a.applicant_category_code, ac.title,
			a.topic_code, t.title,
			a.organization_name,
			a.appeal_kind_code, ak.title,
			a.appeal_type_code, atp.title,
			a.last_name, a.first_name, a.patronymic,
			a.identification_number,
			a.physical_status_code, ps.title,
			a.social_category_code, sc.title,
			a.email, a.is_abroad,
			a.citizenship_code, c.title,
			a.region, a.district, a.locality, a.street, a.house, a.apartment, a.note_address, a.postal_code,
			a.phone, a.content, a.planned_registration_date, a.current_status_code
		FROM appeals a
		JOIN dictionaries ac ON ac.dictionary_type = 'applicant_category' AND ac.code = a.applicant_category_code
		JOIN dictionaries t ON t.dictionary_type = 'topic' AND t.code = a.topic_code
		JOIN dictionaries ak ON ak.dictionary_type = 'appeal_kind' AND ak.code = a.appeal_kind_code
		JOIN dictionaries atp ON atp.dictionary_type = 'appeal_type' AND atp.code = a.appeal_type_code
		LEFT JOIN dictionaries ps ON ps.dictionary_type = 'physical_status' AND ps.code = a.physical_status_code
		LEFT JOIN dictionaries sc ON sc.dictionary_type = 'social_category' AND sc.code = a.social_category_code
		LEFT JOIN dictionaries c ON c.dictionary_type = 'citizenship' AND c.code = a.citizenship_code
		WHERE a.id = $1
	`, id)

	var detail AppealDetailResponse
	var physicalStatusCode *string
	var physicalStatusTitle *string
	var socialCategoryCode *string
	var socialCategoryTitle *string
	var citizenshipCode *string
	var citizenshipTitle *string
	var statusCode string

	err := row.Scan(
		&detail.ID, &detail.Number, &detail.CreatedAt, &detail.UpdatedAt, &detail.SubmittedByUserID,
		&detail.ApplicantCategory.Code, &detail.ApplicantCategory.Title,
		&detail.Topic.Code, &detail.Topic.Title,
		&detail.OrganizationName,
		&detail.AppealKind.Code, &detail.AppealKind.Title,
		&detail.AppealType.Code, &detail.AppealType.Title,
		&detail.LastName, &detail.FirstName, &detail.Patronymic,
		&detail.IdentificationNumber,
		&physicalStatusCode, &physicalStatusTitle,
		&socialCategoryCode, &socialCategoryTitle,
		&detail.Email, &detail.IsAbroad,
		&citizenshipCode, &citizenshipTitle,
		&detail.Region, &detail.District, &detail.Locality, &detail.Street, &detail.House, &detail.Apartment, &detail.NoteAddress, &detail.PostalCode,
		&detail.Phone, &detail.Content, &detail.PlannedRegistrationDate, &statusCode,
	)
	if err != nil {
		return nil, dbError(err)
	}

	if physicalStatusCode != nil && physicalStatusTitle != nil {
		detail.PhysicalStatus = &DictionaryItem{Code: *physicalStatusCode, Title: *physicalStatusTitle}
	}
	if socialCategoryCode != nil && socialCategoryTitle != nil {
		detail.SocialCategory = &DictionaryItem{Code: *socialCategoryCode, Title: *socialCategoryTitle}
	}
	if citizenshipCode != nil && citizenshipTitle != nil {
		detail.Citizenship = &DictionaryItem{Code: *citizenshipCode, Title: *citizenshipTitle}
	}
	if status, ok := GetAppealStatus(statusCode); ok {
		detail.Status = status
	}

	files, err := r.ListFiles(ctx, id)
	if err != nil {
		return nil, err
	}
	history, err := r.ListStatusHistory(ctx, id)
	if err != nil {
		return nil, err
	}

	detail.Files = files
	detail.StatusHistory = history
	return &detail, nil
}

func (r *AppealRepo) ListFiles(ctx context.Context, appealID int64) ([]AppealFileResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, original_name, content_type, size_bytes
		FROM appeal_files
		WHERE appeal_id = $1
		ORDER BY id ASC
	`, appealID)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()

	files := make([]AppealFileResponse, 0)
	for rows.Next() {
		var file AppealFileResponse
		if err := rows.Scan(&file.ID, &file.FileName, &file.ContentType, &file.SizeBytes); err != nil {
			return nil, dbError(err)
		}
		file.DownloadURL = fmt.Sprintf("/api/appeals/%d/files/%d/download", appealID, file.ID)
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		return nil, dbError(err)
	}

	return files, nil
}

func (r *AppealRepo) ListStatusHistory(ctx context.Context, appealID int64) ([]AppealStatusHistoryResponse, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, from_status_code, to_status_code, comment, changed_by_user_id, created_at
		FROM appeal_status_history
		WHERE appeal_id = $1
		ORDER BY id ASC
	`, appealID)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()

	history := make([]AppealStatusHistoryResponse, 0)
	for rows.Next() {
		var item AppealStatusHistoryResponse
		var fromCode *string
		var toCode string
		if err := rows.Scan(&item.ID, &fromCode, &toCode, &item.Comment, &item.ChangedBy, &item.ChangedAt); err != nil {
			return nil, dbError(err)
		}
		if fromCode != nil {
			if fromStatus, ok := GetAppealStatus(*fromCode); ok {
				item.FromStatus = &fromStatus
			}
		}
		if toStatus, ok := GetAppealStatus(toCode); ok {
			item.ToStatus = toStatus
		}
		history = append(history, item)
	}

	if err := rows.Err(); err != nil {
		return nil, dbError(err)
	}

	return history, nil
}

func (r *AppealRepo) GetFileByID(ctx context.Context, fileID int64) (*AppealFileRecord, error) {
	var file AppealFileRecord
	err := r.pool.QueryRow(ctx, `
		SELECT id, appeal_id, original_name, stored_name, storage_path, content_type, size_bytes
		FROM appeal_files
		WHERE id = $1
	`, fileID).Scan(&file.ID, &file.AppealID, &file.OriginalName, &file.StoredName, &file.StoragePath, &file.ContentType, &file.SizeBytes)
	if err != nil {
		return nil, dbError(err)
	}
	return &file, nil
}

func (r *AppealRepo) GetOwnerID(ctx context.Context, appealID int64) (int64, error) {
	var ownerID int64
	err := r.pool.QueryRow(ctx, `SELECT owner_user_id FROM appeals WHERE id = $1`, appealID).Scan(&ownerID)
	if err != nil {
		return 0, dbError(err)
	}
	return ownerID, nil
}

func (r *AppealRepo) ChangeStatus(ctx context.Context, appealID int64, fromStatus string, toStatus string, comment *string, changedBy int64) error {
	tx, err := r.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE appeals SET current_status_code = $2, updated_at = now() WHERE id = $1`, appealID, toStatus)
	if err != nil {
		return dbError(err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO appeal_status_history (appeal_id, from_status_code, to_status_code, comment, changed_by_user_id)
		VALUES ($1, $2, $3, $4, $5)
	`, appealID, fromStatus, toStatus, comment, changedBy)
	if err != nil {
		return dbError(err)
	}

	return tx.Commit(ctx)
}

func composeApplicantName(lastName string, firstName string, patronymic *string) string {
	parts := []string{lastName, firstName}
	if patronymic != nil && *patronymic != "" {
		parts = append(parts, *patronymic)
	}
	return strings.Join(parts, " ")
}
