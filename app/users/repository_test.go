package users

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres" // Correct import
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *gorm.DB
var repoInstance *repo

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "admin",
				"POSTGRES_PASSWORD": "password",
				"POSTGRES_DB":       "mydatabase",
			},
			WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}()

	host, _ := container.Host(ctx)
	port, _ := container.MappedPort(ctx, "5432")

	dsn := fmt.Sprintf("postgres://admin:password@%s:%s/mydatabase?sslmode=disable", host, port.Port())

	// Connect to DB with GORM
	testDB, err = gorm.Open(pg.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}

	// Run migrations
	runMigrations()

	repoInstance = NewUserRepo(testDB)
	// Run tests
	code := m.Run()
	os.Exit(code)
}

func runMigrations() {
	// Open a raw database connection for running the extension command
	sqlDB, err := testDB.DB()
	if err != nil {
		log.Fatalf("failed to get raw db instance: %v", err)
	}

	// Ensure uuid-ossp extension is installed
	_, err = sqlDB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	if err != nil {
		log.Fatalf("failed to create uuid-ossp extension: %v", err)
	}

	// Now set up migration
	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Fatalf("migrate pg driver error: %v", err)
	}

	migrationsPath := "file://../../db/migrations"

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Fatalf("migrate instance error: %v", err)
	}

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migration failed: %v", err)
	}
}

func cleanDB() {
	var tableNames []string

	// Get all user-defined tables (excluding system tables)
	err := testDB.Raw(`
        SELECT tablename
        FROM pg_tables
        WHERE schemaname = 'public' and tablename != 'schema_migrations';
    `).Scan(&tableNames).Error
	if err != nil {
		log.Fatalf("failed to fetch table names: %v", err)
	}

	if len(tableNames) == 0 {
		return
	}

	// Build a TRUNCATE TABLE statement
	truncateSQL := "TRUNCATE TABLE "
	for i, name := range tableNames {
		truncateSQL += fmt.Sprintf(`"%s"`, name)
		if i < len(tableNames)-1 {
			truncateSQL += ", "
		}
	}
	truncateSQL += " RESTART IDENTITY CASCADE;"

	// Execute the truncate
	err = testDB.Exec(truncateSQL).Error
	if err != nil {
		log.Fatalf("failed to truncate tables: %v", err)
	}
}

func TestGetUsers(t *testing.T) {
	cleanDB()
	// Preparing test data
	user := &UserEntity{
		ID:        uuid.New(),
		Username:  "testuser",
		Email:     "testuser@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Insert the user into the DB
	if err := testDB.Create(user).Error; err != nil {
		t.Fatalf("failed to insert user: %v", err)
	}

	// Test the GetUsers method
	users, err := repoInstance.GetUsers()
	assert.NoError(t, err)
	assert.NotNil(t, users)
	assert.Len(t, users, 1)

	// Verify the user returned is the one inserted
	assert.Equal(t, user.Username, users[0].Username)
	assert.Equal(t, user.Email, users[0].Email)
}

func TestCreateUser(t *testing.T) {
	cleanDB() // Clean tables before running the test

	newUser := &UserEntity{
		Username: "johndoe",
		Email:    "john@example.com",
	}

	err := repoInstance.CreateUser(newUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Verify the user was inserted
	var found UserEntity
	result := testDB.First(&found, "username = ?", "johndoe")
	if result.Error != nil {
		t.Fatalf("expected user to be found, got error: %v", result.Error)
	}

	if found.Email != "john@example.com" {
		t.Errorf("expected email 'john@example.com', got %s", found.Email)
	}

	var users []*UserEntity
	result = testDB.Find(&users)
	if result.Error != nil {
		t.Fatalf("expected user to be found, got error: %v", result.Error)
	}

	assert.Equal(t, len(users), 1)
}
