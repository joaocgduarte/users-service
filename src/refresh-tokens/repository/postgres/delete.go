package postgres

import (
	"context"

	"github.com/google/uuid"
)

func (r PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE id = $1
	`

	stmt, err := r.Db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, id)
	return err
}
