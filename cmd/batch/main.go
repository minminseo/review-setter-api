package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/minminseo/recall-setter/infrastructure/db"
	"github.com/minminseo/recall-setter/infrastructure/repository"
	batchUsecase "github.com/minminseo/recall-setter/usecase/batch"
)

func main() {
	// ログ収集ツールとの連携想定でJSON形式で出力
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := db.NewDB(ctx)
	if err != nil {
		slog.Error("データベース接続に失敗しました。処理を続行できません。", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	batchRepository := repository.NewBatchRepository()
	batchUsecase := batchUsecase.NewBatchUsecase(batchRepository)

	runAlignedQuarterHourlyScheduler(batchUsecase)
}

// タイムアウト付きのContextを生成し、バッチ処理の単一の実行をカプセル化
func executeBatch(uc batchUsecase.IBatchUsecase, t time.Time) {
	slog.Info("15分間隔バッチ処理を開始します。", "実行時刻", t.Format(time.RFC3339))

	// バッチ処理一回ごとに独立したタイムアウト付きContextを生成
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	if err := uc.ExecuteUpdateOverdueScheduledDates(ctx); err != nil {
		// エラーハンドリングはユースケース層に委譲されており、ここではエラーの発生を検知するのみです。
		// 必要に応じて、リトライ処理や、監視システムへの通知などをここに追加することもできます。
	}
}

// IANAのタイムゾーンはUTCからのオフセットが全部15分単位なので、0, 15, 30, 45分のタイミングで実行
func runAlignedQuarterHourlyScheduler(uc batchUsecase.IBatchUsecase) {
	slog.Info("壁時計同期・15分間隔実行バッチスケジューラーを起動しました。")

	// 初回実行時刻の計算と待機
	now := time.Now()
	// 現在時刻の「分」を15で割った余りを計算し、次の15分マークまでの待機時間を算出
	remainder := now.Minute() % 15
	waitMinutes := 15 - remainder

	// 次の実行時刻を算出
	// 待機時間を足した後、Truncateで秒以下を切り捨てて正確な次の実行時刻を算出
	nextRun := now.Add(time.Duration(waitMinutes) * time.Minute).Truncate(time.Minute)

	slog.Info("最初のバッチ実行まで待機します。", "初回実行時刻", nextRun.Format(time.RFC3339))

	// 次の実行時間までゴルーチンを休止（CPUリソースを無駄に消費するビジーループ対策）
	time.Sleep(time.Until(nextRun))

	// 算出した初回実行時刻になったら、最初のバッチを実行（tickerの起動が0秒のタイミングからずれないようにゴルーチン使用）
	go executeBatch(uc, time.Now())

	// 初回実行後は、Tickerで15分ごとにバッチを実行するように設定
	ticker := time.NewTicker(15 * time.Minute)
	defer ticker.Stop()

	// ticker.Cからの通知を待ち、15分ごとにバッチを実行する無限ループに入る
	for execTime := range ticker.C {
		go executeBatch(uc, execTime)
	}
}
