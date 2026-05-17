package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gliedabrennung/messenger-core/internal/config"
	"github.com/gliedabrennung/messenger-core/internal/repository/postgres"
	"github.com/gliedabrennung/messenger-core/internal/usecase"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Seeder populates the database with demo data for development and presentation.
func main() {
	cfg := config.GetConfig()

	dbpool, err := pgxpool.New(context.Background(), cfg.DSN)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer dbpool.Close()

	// Run migrations first
	m, err := migrate.New("file://migrations", cfg.DSN)
	if err != nil {
		log.Fatalf("Migration init error: %v", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	repo := postgres.NewPostgresRepository(dbpool)
	authUC := usecase.NewAuthUseCase(repo, cfg.JWTSecret, cfg.JWTTTL)
	contactUC := usecase.NewContactUseCase(repo)
	msgUC := usecase.NewMessageUseCase(repo)
	ctx := context.Background()

	// ---- 1. Create users ----
	type seedUser struct {
		username string
		password string
		role     string
	}
	seedUsers := []seedUser{
		{"admin", "admin1234", "admin"},
		{"alice", "alice1234", "user"},
		{"bob", "bob12345", "user"},
		{"charlie", "charlie1234", "user"},
	}

	userIDs := make(map[string]int64)
	for _, su := range seedUsers {
		user, err := authUC.Register(ctx, su.username, su.password)
		if err != nil {
			fmt.Printf("  skip user %s: %v\n", su.username, err)
			// still need the ID — look it up
			existing, lookupErr := repo.GetByUsername(ctx, su.username)
			if lookupErr == nil {
				userIDs[su.username] = existing.ID
			}
			continue
		}
		userIDs[su.username] = user.ID
		fmt.Printf("  created user: %s (id=%d)\n", su.username, user.ID)
	}

	// ---- 2. Set up contacts (alice <-> bob, alice <-> charlie) ----
	type contactPair struct{ requester, accepter string }
	pairs := []contactPair{
		{"alice", "bob"},
		{"alice", "charlie"},
		{"bob", "charlie"},
	}
	for _, p := range pairs {
		uid, cid := userIDs[p.requester], userIDs[p.accepter]
		if uid == 0 || cid == 0 {
			continue
		}
		_ = contactUC.SendRequest(ctx, uid, cid)
		_ = contactUC.AcceptRequest(ctx, cid, uid)
		fmt.Printf("  contacts: %s <-> %s\n", p.requester, p.accepter)
	}

	// ---- 3. Seed messages ----
	type seedMsg struct{ from, to, content string }
	messages := []seedMsg{
		{"alice", "bob", "Hey Bob! How are you?"},
		{"bob", "alice", "Hi Alice! I'm doing great, thanks!"},
		{"alice", "bob", "Want to catch up later?"},
		{"bob", "alice", "Sure, I'm free at 6pm."},
		{"alice", "charlie", "Charlie, did you see the project requirements?"},
		{"charlie", "alice", "Yes! The Go backend looks solid."},
		{"charlie", "bob", "Bob, did you finish the tests?"},
		{"bob", "charlie", "Almost done, just adding more coverage."},
	}
	for _, sm := range messages {
		from, to := userIDs[sm.from], userIDs[sm.to]
		if from == 0 || to == 0 {
			continue
		}
		_, err := msgUC.SendMessage(ctx, from, to, "text", sm.content)
		if err != nil {
			fmt.Printf("  skip message: %v\n", err)
			continue
		}
		fmt.Printf("  message: [%s -> %s] %s\n", sm.from, sm.to, sm.content)
	}

	fmt.Println("\n✅ Seed complete!")
	fmt.Println("\nDemo credentials:")
	for _, su := range seedUsers {
		fmt.Printf("  %-10s / %s  (role: %s)\n", su.username, su.password, su.role)
	}
	os.Exit(0)
}
