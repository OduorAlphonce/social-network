package repositories

import (
	"database/sql"

	"github.com/gofrs/uuid/v5"
	"learn.zone01kisumu.ke/git/qquinton/social-network/internal/models"
)

type sqliteGroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) GroupRepository {
	return &sqliteGroupRepository{db: db}
}

func (r *sqliteGroupRepository) CreateGroup(group *models.Group) error {
	query := `INSERT INTO groups (id, creator_id, title, description, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, group.ID.String(), group.CreatorID.String(), group.Title, group.Description, group.CreatedAt)
	return err
}

func (r *sqliteGroupRepository) GetGroupByID(id uuid.UUID) (*models.Group, error) {
	query := `SELECT id, creator_id, title, description, created_at FROM groups WHERE id = ?`
	row := r.db.QueryRow(query, id.String())

	var (
		rawID, rawCreatorID string
		g                   models.Group
		desc                sql.NullString
		createdAt           string
	)

	err := row.Scan(&rawID, &rawCreatorID, &g.Title, &desc, &createdAt)
	if err != nil {
		return nil, err
	}

	g.ID, _ = uuid.FromString(rawID)
	g.CreatorID, _ = uuid.FromString(rawCreatorID)
	g.Description = desc.String

	parsedCreatedAt, err := parseSQLiteTime(createdAt)
	if err != nil {
		return nil, err
	}
	g.CreatedAt = parsedCreatedAt

	return &g, nil
}

func (r *sqliteGroupRepository) ListGroups() ([]*models.Group, error) {
	query := `SELECT id, creator_id, title, description, created_at FROM groups ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []*models.Group
	for rows.Next() {
		var (
			rawID, rawCreatorID string
			g                   models.Group
			desc                sql.NullString
			createdAt           string
		)

		if err := rows.Scan(&rawID, &rawCreatorID, &g.Title, &desc, &createdAt); err != nil {
			return nil, err
		}

		g.ID, _ = uuid.FromString(rawID)
		g.CreatorID, _ = uuid.FromString(rawCreatorID)
		g.Description = desc.String

		parsedCreatedAt, err := parseSQLiteTime(createdAt)
		if err != nil {
			return nil, err
		}
		g.CreatedAt = parsedCreatedAt

		groups = append(groups, &g)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return groups, nil
}
