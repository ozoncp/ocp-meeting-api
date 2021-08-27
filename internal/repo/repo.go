package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/ozoncp/ocp-meeting-api/internal/models"
)

type Repo interface {
	Add(ctx context.Context, meeting *models.Meeting) error
	AddMany(ctx context.Context, meetings []models.Meeting) error
	Describe(ctx context.Context, meetingId uint64) (*models.Meeting, error)
	Update(ctx context.Context, meeting models.Meeting) (bool, error)
	List(ctx context.Context, limit, offset uint64) ([]models.Meeting, error)
	Remove(ctx context.Context, id uint64) (bool, error)
}

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) Repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Add(ctx context.Context, meeting *models.Meeting) error {
	query := sq.
		Insert("meeting").
		Columns("user_id", "link", "start", "\"end\"").
		Values(meeting.UserId, meeting.Link, meeting.Start, meeting.End).
		Suffix("RETURNING \"id\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	err := query.QueryRowContext(ctx).Scan(&meeting.Id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) AddMany(ctx context.Context, meetings []models.Meeting) error {
	query := sq.Insert("meeting").
		Columns("user_id", "link", "start", "\"end\"").
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	for _, meeting := range meetings {
		query = query.Values(meeting.UserId, meeting.Link, meeting.Start, meeting.End)
	}

	if _, err := query.QueryContext(ctx); err != nil {
		return err
	}
	return nil
}

func (r *repo) Describe(ctx context.Context, id uint64) (*models.Meeting, error) {
	query := sq.Select("id", "user_id", "link", "start", "\"end\"").
		From("meeting").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	var meeting models.Meeting

	if err := query.QueryRowContext(ctx).
		Scan(&meeting.Id, &meeting.UserId, &meeting.Link, &meeting.Start, &meeting.End); err != nil {
		return nil, err
	}

	return &meeting, nil
}

func (r *repo) Update(ctx context.Context, meeting models.Meeting) (bool, error) {
	query := sq.Update("meeting").
		Set("user_id", meeting.UserId).
		Set("link", meeting.Link).
		Set("start", meeting.Start).
		Set("\"end\"", meeting.End).
		Where(sq.Eq{"id": meeting.Id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	exec, err := query.ExecContext(ctx)
	if err != nil {
		return false, err
	}

	rowsAffected, err := exec.RowsAffected()
	if err != nil {
		return false, err
	}

	return rowsAffected > 0, nil
}

func (r *repo) List(ctx context.Context, limit, offset uint64) ([]models.Meeting, error) {
	query := sq.Select("id", "user_id", "link", "start", "\"end\"").
		From("meeting").
		Offset(offset).
		Limit(limit).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	rows, err := query.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	meetings := make([]models.Meeting, 0, limit)
	for rows.Next() {
		meeting := models.Meeting{}
		if err := rows.Scan(&meeting.Id, &meeting.UserId, &meeting.Link, &meeting.Start, &meeting.End); err != nil {
			return nil, err
		}
		meetings = append(meetings, meeting)
	}
	return meetings, nil
}

func (r *repo) Remove(ctx context.Context, id uint64) (bool, error) {
	query := sq.Delete("meeting").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).
		RunWith(r.db)

	ret, err := query.ExecContext(ctx)
	if err != nil {
		return false, err
	}

	rowsDeleted, err := ret.RowsAffected()
	if err != nil {
		return false, err
	}
	return rowsDeleted > 0, err

}
