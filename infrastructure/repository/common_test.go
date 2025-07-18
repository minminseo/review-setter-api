package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/go-testfixtures/testfixtures/v3"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minminseo/recall-setter/infrastructure/db"
	dbTest "github.com/minminseo/recall-setter/infrastructure/db/db_test"
	"github.com/minminseo/recall-setter/infrastructure/db/dbgen"
	"github.com/ory/dockertest/v3"
)

var (
	testDB       *sql.DB
	testPool     *pgxpool.Pool
	testDBPool   *pgxpool.Pool // For TransactionManager
	testQueries  *dbgen.Queries
	testFixtures *testfixtures.Loader
	dockerPool   *dockertest.Pool
	resource     *dockertest.Resource
)

func TestMain(m *testing.M) {
	var err error

	// コンテナの作成と起動
	resource, dockerPool = dbTest.CreateContainer()

	// データベースに接続
	testDB = dbTest.ConnectDB(resource, dockerPool)

	// マイグレーションの実行
	if err = dbTest.SetupTestDB(testDB, "../../migrations"); err != nil {
		dbTest.CloseContainer(resource, dockerPool)
		panic(err)
	}

	// pgx接続プールを作成
	testPool, err = pgxpool.New(context.Background(), dbTest.GetDSN())
	if err != nil {
		dbTest.CloseContainer(resource, dockerPool)
		panic(err)
	}

	// テスト用のクエリインスタンスを作成
	testQueries = dbgen.New(testPool)

	// TransactionManager用にプールのエイリアスを作成
	testDBPool = testPool

	// フィクスチャの読み込み
	testFixtures, err = testfixtures.New(
		testfixtures.Database(testDB),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory("../fixtures"),
	)
	if err != nil {
		dbTest.CloseContainer(resource, dockerPool)
		panic(err)
	}

	// テストの実行
	m.Run()

	// クリーンアップ
	if testPool != nil {
		testPool.Close()
	}
	if testDB != nil {
		testDB.Close()
	}

	dbTest.CloseContainer(resource, dockerPool)
}

// PrepareTestDatabase loads fixtures into the test database
func PrepareTestDatabase(t *testing.T) {
	t.Helper()
	if err := testFixtures.Load(); err != nil {
		t.Fatalf("Could not load fixtures: %v", err)
	}
}

// GetTestDB returns the test database connection
func GetTestDB() *sql.DB {
	return testDB
}

// GetTestPool returns the test pgx pool
func GetTestPool() *pgxpool.Pool {
	return testPool
}

// GetTestContext returns a context with the test database queries
func GetTestContext() context.Context {
	return db.WithQueries(context.Background(), testQueries)
}

// CleanupTestDatabase truncates all tables (optional cleanup function)
func CleanupTestDatabase(t *testing.T) {
	t.Helper()

	tables := []string{
		"email_verifications",
		"review_dates",
		"review_items",
		"review_boxes",
		"pattern_steps",
		"review_patterns",
		"categories",
		"users",
	}

	for _, table := range tables {
		if _, err := testDB.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)); err != nil {
			t.Errorf("Failed to truncate table %s: %v", table, err)
		}
	}
}
