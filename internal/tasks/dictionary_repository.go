package tasks

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DictionaryRepo struct {
	pool *pgxpool.Pool
}

func NewDictionaryRepo(pool *pgxpool.Pool) *DictionaryRepo {
	return &DictionaryRepo{pool: pool}
}

func (r *DictionaryRepo) ListByType(ctx context.Context, dictionaryType string) ([]DictionaryItem, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT code, title
		FROM dictionaries
		WHERE dictionary_type = $1 AND is_active = true
		ORDER BY sort_order ASC, title ASC
	`, dictionaryType)
	if err != nil {
		return nil, dbError(err)
	}
	defer rows.Close()

	items := make([]DictionaryItem, 0)
	for rows.Next() {
		var item DictionaryItem
		if err := rows.Scan(&item.Code, &item.Title); err != nil {
			return nil, dbError(err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, dbError(err)
	}

	return items, nil
}

func (r *DictionaryRepo) Exists(ctx context.Context, dictionaryType string, code string) (bool, error) {
	var exists bool
	err := r.pool.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1 FROM dictionaries WHERE dictionary_type = $1 AND code = $2 AND is_active = true
		)
	`, dictionaryType, code).Scan(&exists)
	if err != nil {
		return false, dbError(err)
	}
	return exists, nil
}
