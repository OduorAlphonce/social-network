package repositories

import (
	"database/sql"

	"github.com/gofrs/uuid/v5"
)

type sqliteGroupMembershipRepository struct {
	db *sql.DB
}

// NewGroupMembershipRepository creates a SQLite-backed group membership reader.
func NewGroupMembershipRepository(db *sql.DB) GroupMembershipRepository {
	return &sqliteGroupMembershipRepository{db: db}
}

func (r *sqliteGroupMembershipRepository) IsAcceptedGroupMember(groupID, userID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1
			FROM group_members
			WHERE group_id = ?
				AND user_id = ?
				AND status = 'accepted'
		)
	`, groupID.String(), userID.String()).Scan(&exists)
	return exists, err
}
