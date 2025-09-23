package persistence

import (
	"context"
	"fmt"
	"github.com/eldius/rpi-system-monitor/internal/model"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/tsdb"
	"github.com/prometheus/prometheus/tsdb/chunkenc"
	"log/slog"
	"math"
	"time"
)

func openDB() (*tsdb.DB, error) {
	db, err := tsdb.Open(".db/tsdb.db", slog.With("pkg", "persistence"), nil, tsdb.DefaultOptions(), nil)
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	return db, nil
}

func Persist(ctx context.Context, result *model.ProbesResult) error {
	db, err := openDB()
	if err != nil {
		return fmt.Errorf("opening database: %w", err)
	}
	defer func() {
		_ = db.Close()
	}()

	if err := saveTemperature(ctx, db, result.Temp, result.Timestamp.Unix()); err != nil {
		return fmt.Errorf("appending temperature: %w", err)
	}

	if err := saveMemoryUsage(ctx, db, result.Memory, result.Timestamp.Unix()); err != nil {
		return fmt.Errorf("appending memory: %w", err)
	}

	if err := saveCPUUsage(ctx, db, result.CPU, result.Timestamp.Unix()); err != nil {
		return fmt.Errorf("appending cpu usage: %w", err)
	}
	return nil
}

func Get(ctx context.Context) ([]model.ProbesResult, error) {
	db, err := openDB()
	if err != nil {
		return nil, fmt.Errorf("opening database: %w", err)
	}
	tempTS, err := getTemperature(ctx, db)
	if err != nil {
		return nil, err
	}
	memTS, err := getMemoryUsage(ctx, db)
	if err != nil {
		return nil, err
	}
	cpuTS, err := getCPUUsage(ctx, db)

	var result []model.ProbesResult
	for k, v := range tempTS {
		result = append(result, model.ProbesResult{
			Timestamp: time.Unix(k, 0),
			Temp: model.TemperatureResult{
				Temperature:    v[temperatureIdx],
				RawTemperature: int64(v[rawTemperatureIdx]),
			},
			Memory: model.MemoryResult{
				MemoryUsagePercentage: memTS[k][memoryUsagePercentageIdx],
				UsedMemory:            int64(memTS[k][usedMemoryIdx]),
				TotalMemory:           int64(memTS[k][totalMemoryIdx]),
			},
			CPU: model.CPUResult{
				CPUCount: int64(cpuTS[k][cpuCountIdx]),
				CPUUsage: cpuTS[k][cpuUsageIdx],
			},
		})
	}

	return result, nil
}

func saveTemperature(ctx context.Context, db *tsdb.DB, result model.TemperatureResult, timestamp int64) error {
	if err := persist(ctx, db, timestamp, result.Temperature, []string{dimensionLabelName, temperature, "unit", "celsius"}); err != nil {
		return fmt.Errorf("appending temperature: %w", err)
	}

	if err := persist(ctx, db, timestamp, float64(result.RawTemperature), []string{dimensionLabelName, rawTemperature}); err != nil {
		return fmt.Errorf("appending raw temperature: %w", err)
	}

	return nil
}

func getTemperature(ctx context.Context, db *tsdb.DB) (map[int64][2]float64, error) {
	t, err := fetch(ctx, db, [2]string{dimensionLabelName, temperature})
	if err != nil {
		return nil, fmt.Errorf("fetching temperature: %w", err)
	}
	rt, err := fetch(ctx, db, [2]string{dimensionLabelName, rawTemperature})
	if err != nil {
		return nil, fmt.Errorf("fetching raw temperature: %w", err)
	}
	result := make(map[int64][2]float64)
	for k, v := range t {
		result[k] = [2]float64{
			v,
			rt[k],
		}
	}

	return result, nil
}

func saveMemoryUsage(ctx context.Context, db *tsdb.DB, result model.MemoryResult, timestamp int64) error {
	if err := persist(ctx, db, timestamp, result.MemoryUsagePercentage, []string{dimensionLabelName, memoryUsagePercentage, "unit", "percent"}); err != nil {
		return fmt.Errorf("appending memory usage percentage: %w", err)
	}
	if err := persist(ctx, db, timestamp, float64(result.UsedMemory), []string{dimensionLabelName, usedMemory, "unit", "bytes"}); err != nil {
		return fmt.Errorf("appending used memory: %w", err)
	}
	if err := persist(ctx, db, timestamp, float64(result.TotalMemory), []string{dimensionLabelName, totalMemory, "unit", "bytes"}); err != nil {
		return fmt.Errorf("appending total memory: %w", err)
	}
	return nil
}

func getMemoryUsage(ctx context.Context, db *tsdb.DB) (map[int64][3]float64, error) {
	var result = make(map[int64][3]float64)
	mup, err := fetch(ctx, db, [2]string{dimensionLabelName, memoryUsagePercentage})
	if err != nil {
		return nil, err
	}
	um, err := fetch(ctx, db, [2]string{dimensionLabelName, usedMemory})
	if err != nil {
		return nil, err
	}
	tm, err := fetch(ctx, db, [2]string{dimensionLabelName, totalMemory})
	if err != nil {
		return nil, err
	}

	for k, v := range mup {
		result[k] = [3]float64{
			v,
			um[k],
			tm[k],
		}
	}
	return result, nil
}

func saveCPUUsage(ctx context.Context, db *tsdb.DB, result model.CPUResult, timestamp int64) error {
	if err := persist(ctx, db, timestamp, result.CPUUsage, []string{dimensionLabelName, cpuUsage, "unit", "percent"}); err != nil {
		return fmt.Errorf("appending cpu usage: %w", err)
	}
	if err := persist(ctx, db, timestamp, float64(result.CPUCount), []string{dimensionLabelName, cpuCount, "unit", "count"}); err != nil {
		return fmt.Errorf("appending cpu count: %w", err)
	}

	return nil
}

func getCPUUsage(ctx context.Context, db *tsdb.DB) (map[int64][2]float64, error) {
	cc, err := fetch(ctx, db, [2]string{dimensionLabelName, cpuCount})
	if err != nil {
		return nil, fmt.Errorf("fetching cpu usage: %w", err)
	}
	cu, err := fetch(ctx, db, [2]string{dimensionLabelName, cpuUsage})
	if err != nil {
		return nil, fmt.Errorf("fetching cpu usage: %w", err)
	}

	result := make(map[int64][2]float64)
	for k, v := range cc {
		result[k] = [2]float64{
			v,
			cu[k],
		}
	}
	return result, nil
}

func persist(ctx context.Context, db *tsdb.DB, timestamp int64, value float64, lbl []string) error {
	appender := db.Appender(ctx)
	defer func() {
		_ = appender.Rollback()
	}()
	ref, err := appender.Append(0, labels.FromStrings(lbl...), timestamp, value)
	if err != nil {
		return fmt.Errorf("appending temperature: %w", err)
	}
	fmt.Println("ref:", ref)
	return appender.Commit()
}

func fetch(ctx context.Context, db *tsdb.DB, lbl [2]string) (map[int64]float64, error) {
	querier, err := db.Querier(math.MinInt64, math.MaxInt64)
	if err != nil {
		return nil, fmt.Errorf("opening querier: %w", err)
	}
	defer func() {
		_ = querier.Close()
	}()
	queryResult := querier.Select(ctx, true, nil, labels.MustNewMatcher(labels.MatchEqual, lbl[0], lbl[1]))

	var result = make(map[int64]float64)
	for queryResult.Next() {
		series := queryResult.At()
		fmt.Println("series:", series.Labels().String())

		it := series.Iterator(nil)
		for it.Next() == chunkenc.ValFloat {
			ts, v := it.At() // We ignore the timestamp here, only to have a predictable output we can test against (below)
			fmt.Println("sample", v)
			result[ts] = v
		}

		fmt.Println("it.Err():", it.Err())
	}

	return result, nil
}
