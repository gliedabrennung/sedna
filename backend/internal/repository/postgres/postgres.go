package postgres

import (
	"context"
	"errors"

	"github.com/gliedabrennung/messenger-core/internal/entity"
	"github.com/gliedabrennung/messenger-core/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const uniqueViolationCode = "23505"

type Repository struct {
	db *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{db: pool}
}

func (repo *Repository) Create(ctx context.Context, user *entity.User) error {
	if user.Role == "" {
		user.Role = "user"
	}
	query := `
		INSERT INTO users (username, password, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`
	err := repo.db.QueryRow(ctx, query, user.Username, user.Password, user.Role).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
			return repository.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}

func (repo *Repository) GetByUsername(ctx context.Context, username string) (*entity.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users
		WHERE username = $1`
	user := &entity.User{}
	err := repo.db.QueryRow(ctx, query, username).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (repo *Repository) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	query := `
		SELECT id, username, password, role, created_at, updated_at
		FROM users
		WHERE id = $1`
	user := &entity.User{}
	err := repo.db.QueryRow(ctx, query, id).
		Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

func (repo *Repository) GetAll(ctx context.Context, limit, offset int) ([]entity.User, error) {
	query := `SELECT id, username, role, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2`
	rows, err := repo.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Role, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (repo *Repository) DeleteUser(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := repo.db.Exec(ctx, query, id)
	return err
}

func (repo *Repository) SaveMessage(ctx context.Context, msg *entity.Message) error {
	query := `
		INSERT INTO messages (from_id, to_id, type, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, status, created_at`
	err := repo.db.QueryRow(ctx, query, msg.FromID, msg.ToID, msg.Type, msg.Content).
		Scan(&msg.ID, &msg.Status, &msg.CreatedAt)
	return err
}

func (repo *Repository) GetMessagesBetween(ctx context.Context, user1, user2 int64, limit, offset int) ([]entity.Message, error) {
	query := `
		SELECT id, from_id, to_id, type, content, status, created_at
		FROM messages
		WHERE (from_id = $1 AND to_id = $2) OR (from_id = $2 AND to_id = $1)
		ORDER BY created_at ASC
		LIMIT $3 OFFSET $4`
	rows, err := repo.db.Query(ctx, query, user1, user2, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []entity.Message
	for rows.Next() {
		var m entity.Message
		if err := rows.Scan(&m.ID, &m.FromID, &m.ToID, &m.Type, &m.Content, &m.Status, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}

func (repo *Repository) UpdateMessageStatuses(ctx context.Context, fromID, toID int64, status string) error {
	query := `UPDATE messages SET status = $3 WHERE from_id = $1 AND to_id = $2 AND status != $3 AND status != 'read'`
	_, err := repo.db.Exec(ctx, query, fromID, toID, status)
	return err
}

func (repo *Repository) GetUnreadCounts(ctx context.Context, userID int64) (map[int64]int, error) {
	query := `
		SELECT from_id, COUNT(*) 
		FROM messages 
		WHERE to_id = $1 AND status != 'read' 
		GROUP BY from_id`
	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[int64]int)
	for rows.Next() {
		var fromID int64
		var count int
		if err := rows.Scan(&fromID, &count); err != nil {
			return nil, err
		}
		counts[fromID] = count
	}
	return counts, rows.Err()
}

func (repo *Repository) AddContact(ctx context.Context, userID, contactID int64) error {
	query := `INSERT INTO contacts (user_id, contact_id, status) VALUES ($1, $2, 'pending') ON CONFLICT DO NOTHING`
	_, err := repo.db.Exec(ctx, query, userID, contactID)
	return err
}

func (repo *Repository) UpdateStatus(ctx context.Context, userID, contactID int64, status string) error {
	query := `UPDATE contacts SET status = $3 WHERE user_id = $1 AND contact_id = $2`
	_, err := repo.db.Exec(ctx, query, userID, contactID, status)
	return err
}

func (repo *Repository) GetContacts(ctx context.Context, userID int64) ([]entity.User, error) {
	query := `
		SELECT u.id, u.username, u.role, u.created_at, u.updated_at
		FROM users u
		JOIN contacts c ON u.id = c.contact_id
		WHERE c.user_id = $1 AND c.status = 'accepted'
		UNION
		SELECT u.id, u.username, u.role, u.created_at, u.updated_at
		FROM users u
		JOIN contacts c ON u.id = c.user_id
		WHERE c.contact_id = $1 AND c.status = 'accepted'
	`
	rows, err := repo.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var u entity.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Role, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (repo *Repository) DeleteContact(ctx context.Context, userID, contactID int64) error {
	query := `DELETE FROM contacts WHERE (user_id = $1 AND contact_id = $2) OR (user_id = $2 AND contact_id = $1)`
	_, err := repo.db.Exec(ctx, query, userID, contactID)
	return err
}
