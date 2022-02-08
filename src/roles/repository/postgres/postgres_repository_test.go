package postgres

import (
	"context"
	"database/sql/driver"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/plagioriginal/user-microservice/domain"
	"github.com/stretchr/testify/assert"
)

type anyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}

func TestFetch_FailQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL;`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(errors.New("boom"))

	res, err := New(db).Fetch(context.TODO())

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	assert.Empty(t, res)
}

func TestFetch_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL;`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillDelayFor(time.Duration(time.Millisecond * 150)).
		WillReturnError(errors.New("doessn't matter"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*100))
	defer cancel()

	res, err := New(db).Fetch(ctx)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
	assert.Empty(t, res)
}

func TestFetch_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	createdAt := time.Now()
	roleIdOne := uuid.New()
	roleIdTwo := uuid.New()

	expectedResult := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})

	expectedResult.
		AddRow(roleIdOne, "slug", "Slug Role", createdAt, createdAt).
		AddRow(roleIdTwo, "slug2", "Slug Role2", createdAt, createdAt)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL;`
	mock.ExpectQuery(regexp.QuoteMeta(query)).
		WillReturnError(nil).
		WillReturnRows(expectedResult)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*150))
	defer cancel()

	res, err := New(db).Fetch(ctx)
	assert.Nil(t, err)
	assert.Equal(t, res, []domain.Role{
		{
			ID:        roleIdOne,
			RoleSlug:  "slug",
			RoleLabel: "Slug Role",
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			DeletedAt: time.Time{},
		},
		{
			ID:        roleIdTwo,
			RoleSlug:  "slug2",
			RoleLabel: "Slug Role2",
			CreatedAt: createdAt,
			UpdatedAt: createdAt,
			DeletedAt: time.Time{},
		},
	})
}

func TestGetBySlug_FailPrepareQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and role_slug=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("slug").
		WillReturnError(errors.New("boom"))

	res, err := New(db).GetBySlug(context.TODO(), "slug")

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	assert.Empty(t, res)
}

func TestGetBySlug_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and role_slug=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("slug").
		WillDelayFor(time.Duration(time.Millisecond * 150)).
		WillReturnError(errors.New("doessn't matter"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*100))
	defer cancel()

	res, err := New(db).GetBySlug(ctx, "slug")
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
	assert.Empty(t, res)
}

func TestGetBySlug_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	createdAt := time.Now()
	roleId := uuid.New()

	expectedResult := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})

	expectedResult.AddRow(roleId, "slug", "Slug Role", createdAt, createdAt)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and role_slug=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs("slug").
		WillReturnError(nil).
		WillReturnRows(expectedResult)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*150))
	defer cancel()

	res, err := New(db).GetBySlug(ctx, "slug")
	assert.Nil(t, err)
	assert.Equal(t, res, domain.Role{
		ID:        roleId,
		RoleSlug:  "slug",
		RoleLabel: "Slug Role",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		DeletedAt: time.Time{},
	})
}

func TestGetByUuid_FailPrepareQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and id=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id).
		WillReturnError(errors.New("boom"))

	res, err := New(db).GetByUUID(context.TODO(), id)

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	assert.Empty(t, res)
}

func TestGetByUuid_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and id=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id).
		WillDelayFor(time.Duration(time.Millisecond * 150)).
		WillReturnError(errors.New("doessn't matter"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*100))
	defer cancel()

	res, err := New(db).GetByUUID(ctx, id)
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
	assert.Empty(t, res)
}

func TestGetByUuid_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	createdAt := time.Now()
	roleId := uuid.New()

	expectedResult := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})

	expectedResult.AddRow(roleId, "slug", "Slug Role", createdAt, createdAt)

	query := `SELECT id, role_slug, role_label, created_at, updated_at FROM roles WHERE deleted_at IS NULL and id=$1;`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(roleId).
		WillReturnError(nil).
		WillReturnRows(expectedResult)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*150))
	defer cancel()

	res, err := New(db).GetByUUID(ctx, roleId)
	assert.Nil(t, err)
	assert.Equal(t, res, domain.Role{
		ID:        roleId,
		RoleSlug:  "slug",
		RoleLabel: "Slug Role",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		DeletedAt: time.Time{},
	})
}

func TestStore_FailPrepareQuery(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	query := `INSERT INTO roles (id, role_slug, role_label, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id, "slug", "label", anyTime{}, anyTime{}).
		WillReturnError(errors.New("boom"))

	res, err := New(db).Store(context.TODO(), domain.Role{
		ID:        id,
		RoleSlug:  "slug",
		RoleLabel: "label",
	})

	assert.Error(t, err)
	assert.Equal(t, err.Error(), "boom")
	assert.Empty(t, res)
}

func TestStore_TimeoutReached(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	id := uuid.New()
	query := `INSERT INTO roles (id, role_slug, role_label, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id, "slug", "label", anyTime{}, anyTime{}).
		WillDelayFor(time.Duration(time.Millisecond * 150)).
		WillReturnError(errors.New("doessn't matter"))

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*100))
	defer cancel()

	res, err := New(db).Store(ctx, domain.Role{
		ID:        id,
		RoleSlug:  "slug",
		RoleLabel: "label",
	})
	assert.Error(t, err)
	assert.Equal(t, err.Error(), "canceling query due to user request")
	assert.Empty(t, res)
}

func TestStore_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.Nil(t, err)

	createdAt := time.Now()
	roleId := uuid.New()

	expectedResult := sqlmock.NewRows([]string{"id", "role_slug", "role_label", "created_at", "updated_at"})
	expectedResult.AddRow(roleId, "slug", "label", createdAt, createdAt)

	id := uuid.New()
	query := `INSERT INTO roles (id, role_slug, role_label, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING *`
	mock.ExpectPrepare(regexp.QuoteMeta(query)).
		ExpectQuery().
		WithArgs(id, "slug", "label", anyTime{}, anyTime{}).
		WillReturnError(nil).
		WillReturnRows(expectedResult)

	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Millisecond*150))
	defer cancel()

	res, err := New(db).Store(ctx, domain.Role{
		ID:        id,
		RoleSlug:  "slug",
		RoleLabel: "label",
	})
	assert.Nil(t, err)
	assert.Equal(t, res, domain.Role{
		ID:        roleId,
		RoleSlug:  "slug",
		RoleLabel: "label",
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		DeletedAt: time.Time{},
	})
}
