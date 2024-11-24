package repositories

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"strings"
	"test-task/internal/dto"
	"test-task/internal/entities"
	errs "test-task/internal/errors"
)

type Songs struct {
	db *sqlx.DB
}

func NewSongsRepository(db *sqlx.DB) *Songs {
	return &Songs{db: db}
}

func (repo *Songs) Add(ctx context.Context, song entities.Song, verses []string) (int64, error) {

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("couldn't begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := "INSERT INTO songs (\"group\", name, link, release_date) " +
		"VALUES (:group, :name, :link, :release_date) RETURNING id"
	rows, err := tx.NamedQuery(query, song)
	defer rows.Close()

	if err != nil {
		return 0, fmt.Errorf("couldn't insert song: %w", err)
	}

	var id int64
	if rows.Next() {
		if err = rows.Scan(&id); err != nil {
			return 0, fmt.Errorf("couldn't scan rows: %w", err)
		}
	} else {
		return 0, fmt.Errorf("failed to get inserted song id")
	}
	rows.Close() //VERY IMPORTANT TO CALL BEFORE NEXT QUERY, OTHERWISE "pq: unexpected Parse response 'C'" !!!!!

	queryBuilder := squirrel.Insert("verses").
		Columns("song_id", "verse_number", "verse_text").
		PlaceholderFormat(squirrel.Dollar)

	for i := 0; i < len(verses); i++ {
		queryBuilder = queryBuilder.Values(id, i, verses[i])
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("couldn't create query to insert verses: %w", err)
	}

	log.Debugf("verses query: %v, args: %v", query, args)
	_, err = tx.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("couldn't insert verses: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("couldn't commit transaction: %w", err)
	}

	return id, nil
}

func (repo *Songs) Get(ctx context.Context, params dto.GetSongsRequest) ([]entities.Song, error) {

	queryBuilder := squirrel.Select("*").
		From("songs").PlaceholderFormat(squirrel.Dollar)

	if params.Group != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"LOWER(\"group\")": strings.ToLower(params.Group)})
	}
	if params.Name != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"LOWER(name)": strings.ToLower(params.Name)})
	}
	if params.Link != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"LOWER(link)": strings.ToLower(params.Link)})
	}

	releaseAfterFilter, releaseBeforeFilter := false, false
	if !params.ReleaseBefore.IsZero() {
		releaseBeforeFilter = true
		queryBuilder = queryBuilder.Where(squirrel.Lt{"release_date": params.ReleaseBefore})
	}
	if !params.ReleaseAfter.IsZero() {
		releaseAfterFilter = true
		queryBuilder = queryBuilder.Where(squirrel.Gt{"release_date": params.ReleaseAfter})
	}

	//can't add OrderBy inside both ifs because squirrel won't override:
	//it turns into "...ORDER BY release_date DESC, release_date ASC" LOL
	if releaseAfterFilter { //order asc for a case with both release_after and release_before defined
		queryBuilder = queryBuilder.OrderBy("release_date ASC")
	} else if releaseBeforeFilter {
		queryBuilder = queryBuilder.OrderBy("release_date DESC")
	}

	queryBuilder = queryBuilder.
		Offset(uint64((params.Page - 1) * params.PageSize)).
		Limit(uint64(params.PageSize))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("couldn't create query to get songs: %w", err)
	}

	log.Debugf("query: %v, args: %v", query, args)

	songs := make([]entities.Song, 0)
	err = repo.db.SelectContext(ctx, &songs, query, args...)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("couldn't get songs: %w", err)
	}

	return songs, nil
}

func (repo *Songs) GetVerses(ctx context.Context, songID int64, page int, pageSize int) ([]entities.Verse, error) {

	var songExists bool
	err := repo.db.GetContext(ctx, &songExists,
		"SELECT EXISTS(SELECT 1 FROM songs WHERE id = $1)", songID)
	if err != nil {
		return nil, fmt.Errorf("couldn't check song existance: %w", err)
	}
	if !songExists {
		return nil, fmt.Errorf("%w: song with id %d doesn't exist", errs.NotFound, songID)
	}

	res := make([]entities.Verse, 0)
	err = repo.db.SelectContext(ctx, &res,
		"SELECT * FROM verses WHERE song_id = $1 LIMIT $2 OFFSET $3", songID, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, fmt.Errorf("couldn't get verses: %w", err)
	}

	return res, nil
}

func (repo *Songs) Update(ctx context.Context, song entities.Song) error {

	res, err := repo.db.NamedExecContext(ctx,
		"UPDATE songs SET name = :name, \"group\" = :group, release_date = :release_date, link = :link"+
			" WHERE id = :id", song)
	if err != nil {
		return fmt.Errorf("couldn't update song: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't get RowsAffected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: song with id %d doesn't exist", errs.NotFound, song.ID)
	}
	return nil
}

func (repo *Songs) Delete(ctx context.Context, ID int64) error {

	res, err := repo.db.ExecContext(ctx, "DELETE FROM songs WHERE id = $1", ID)
	if err != nil {
		return fmt.Errorf("couldn't delete song: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("couldn't get RowsAffected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("%w: song with id %d doesn't exist", errs.NotFound, ID)
	}

	return nil
}
