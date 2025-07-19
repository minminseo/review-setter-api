package dbTest

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

var (
	username = "testuser"
	password = "testpassword"
	hostname = "localhost"
	dbName   = "testdb"
	port     int // コンテナ起動時に決定する
)

func CreateContainer() (*dockertest.Resource, *dockertest.Pool) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	pool.MaxWait = time.Minute * 2

	// Dockerコンテナ起動時の細かいオプションを指定する
	runOptions := &dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_USER=" + username,
			"POSTGRES_DB=" + dbName,
			"listen_addresses='*'",
		},
	}

	// ホスト設定の関数
	hostConfigFunc := func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	}

	// コンテナを起動
	resource, err := pool.RunWithOptions(runOptions, hostConfigFunc)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// コンテナの有効期限を設定 (2分)
	if err := resource.Expire(120); err != nil {
		log.Fatalf("Could not set resource expiration: %s", err)
	}

	return resource, pool
}

func CloseContainer(resource *dockertest.Resource, pool *dockertest.Pool) {
	// コンテナの終了
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func ConnectDB(resource *dockertest.Resource, pool *dockertest.Pool) *sql.DB {
	// DB(コンテナ)との接続
	var db *sql.DB
	if err := pool.Retry(func() error {
		var err error
		portStr := resource.GetPort("5432/tcp")
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return err
		}

		dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
			username, password, hostname, portStr, dbName)

		db, err = sql.Open("pgx", dbURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}
	return db
}

func SetupTestDB(db *sql.DB, migrationsPath string) error {
	// マイグレーションファイルのパスを取得
	migrationsAbsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to get migrations path: %w", err)
	}

	// データベースドライバーを作成
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to create database driver: %w", err)
	}

	// マイグレーションインスタンスを作成
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsAbsPath),
		"postgres",
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	// マイグレーションを実行
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func GetDSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		username, password, hostname, port, dbName)
}
