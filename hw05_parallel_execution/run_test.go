package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

/*

// переписал на табличные тесты

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})
}
*/

func TestRunTable(t *testing.T) {
	tests := []struct {
		testName       string
		workersCount   int
		maxErrorsCount int
		goodTasksCount int
		errTasksCount  int
	}{
		{
			testName:     "if were errors in first M tasks, than finished not more N+M tasks",
			workersCount: 10, maxErrorsCount: 23, goodTasksCount: 0, errTasksCount: 50,
		},
		{testName: "tasks without errors", workersCount: 5, maxErrorsCount: 1, goodTasksCount: 50, errTasksCount: 0},
		{testName: "errors more than max errors", workersCount: 5, maxErrorsCount: 20, goodTasksCount: 20, errTasksCount: 30},
		{testName: "errors less than max errors", workersCount: 5, maxErrorsCount: 40, goodTasksCount: 20, errTasksCount: 30},
		{testName: "max errors = 0", workersCount: 5, maxErrorsCount: 0, goodTasksCount: 20, errTasksCount: 30},
		{testName: "negativ max errors", workersCount: 5, maxErrorsCount: -10, goodTasksCount: 20, errTasksCount: 30},
		{testName: "workersCount = 0", workersCount: 0, maxErrorsCount: 10, goodTasksCount: 20, errTasksCount: 30},
		{testName: "tasks Count = 0", workersCount: 5, maxErrorsCount: 10, goodTasksCount: 0, errTasksCount: 0},
	}

	defer goleak.VerifyNone(t)

	for _, ts := range tests {
		t.Run(ts.testName, func(t *testing.T) {
			tasksCount := ts.goodTasksCount + ts.errTasksCount
			tasks := make([]Task, 0, tasksCount)

			var runTasksCount int32
			var sumTime time.Duration

			// генерируем задачи. вначале - с возвратом шибки, затем - без ошибок
			for i := 0; i < tasksCount; i++ {
				var err error
				taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
				sumTime += taskSleep

				if i < ts.errTasksCount {
					err = fmt.Errorf("error from task %d", i)
				}
				tasks = append(tasks, func() error {
					time.Sleep(taskSleep)
					atomic.AddInt32(&runTasksCount, 1)
					return err
				})
			}

			start := time.Now()
			err := Run(tasks, ts.workersCount, ts.maxErrorsCount)
			elapsedTime := time.Since(start)

			// Правило для параметра maxErrCount: если maxErrCount = 0 - проверка ошибок выключена (на ошибки не проверяем).
			// Контроль результатов теста:
			// если указаны некорректные параметры (нет обработчиков или нет задач),
			// то проверяем результат на ошибку ErrInvalidParameters
			// иначе, если проверка ошибок включена и ожидаем кол-во ошибок больше чем maxErrCount,
			// то идем по ветке контроля срабатывания прерывания по кол-ву ошибок
			// иначе - идем по ветке контроля обработки всего пакета задач.
			switch {
			case ts.workersCount <= 0 || tasksCount == 0:
				require.Truef(t, errors.Is(err, ErrInvalidParameters), "actual err - %v", err)
			case ts.maxErrorsCount > 0 && ts.errTasksCount >= ts.maxErrorsCount:
				require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
				require.LessOrEqual(t, runTasksCount, int32(ts.workersCount+ts.maxErrorsCount), "extra tasks were started")
			default:
				require.NoError(t, err)
				require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
				require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
			}
		})
	}
}
