package pool

import (
	"context"
	"runtime"
	"sync"
)

// Task представляет задачу, которая будет выполняться воркером.
type Task func(workerID int)

// WorkerPool представляет пул воркеров для выполнения задач.
type WorkerPool struct {
	workerID      int                // Уникальный идентификатор воркера
	maxWorkers    int                // Максимальное количество воркеров
	tasks         chan Task          // Канал задач для выполнения
	cancel        context.CancelFunc // Функция отмены контекста
	wg            sync.WaitGroup     // Группа ожидания для синхронизации завершения работы воркеров
	mu            sync.Mutex         // Мьютекс для безопасного изменения состояния пула
	ctx           context.Context    // Контекст для завершения работы воркеров
	stopWorkerCh  chan struct{}      // Канал для остановки воркеров
	activeWorkers int                // Текущее количество активных воркеров
}

// NewWorkerPool создает новый пул воркеров с заданным контекстом.
// По умолчанию количество воркеров равно количеству доступных ядер процессора.
func NewWorkerPool(ctx context.Context) *WorkerPool {
	ctx, cancel := context.WithCancel(ctx)
	numCPU := runtime.NumCPU()
	wp := &WorkerPool{
		maxWorkers:   numCPU,
		tasks:        make(chan Task),
		cancel:       cancel,
		ctx:          ctx,
		workerID:     0,
		stopWorkerCh: make(chan struct{}),
	}
	wp.startWorkers(numCPU)
	return wp
}

// startWorkers запускает указанное количество воркеров.
func (wp *WorkerPool) startWorkers(num int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < num; i++ {
		wp.wg.Add(1)
		wp.workerID++
		wp.activeWorkers++
		go wp.worker(wp.workerID)
	}
}

// stopWorkers останавливает указанное количество воркеров.
func (wp *WorkerPool) stopWorkers(num int) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	for i := 0; i < num; i++ {
		wp.stopWorkerCh <- struct{}{}
		wp.activeWorkers--
	}
}

// worker представляет собой основной цикл выполнения задач для одного воркера.
func (wp *WorkerPool) worker(workerID int) {
	defer wp.wg.Done()
	for {
		select {
		case <-wp.ctx.Done(): // Завершение работы воркера при отмене контекста
			return
		case <-wp.stopWorkerCh: // Завершение работы воркера по сигналу остановки
			return
		case task, ok := <-wp.tasks: // Получение задачи из канала
			if !ok {
				return
			}
			task(workerID)
		}
	}
}

// Submit добавляет задачу в пул для выполнения.
// Если пул завершает работу, задача не будет добавлена.
func (wp *WorkerPool) Submit(task Task) {
	select {
	case wp.tasks <- task:
	default:
	}
}

// Stop завершает работу всех воркеров и закрывает пул.
func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.wg.Wait() // Ожидание завершения всех воркеров
	wp.mu.Lock()
	close(wp.tasks) // Закрытие канала задач
	wp.mu.Unlock()
}

// Resize изменяет количество воркеров в пуле.
func (wp *WorkerPool) Resize(newSize int) {
	wp.mu.Lock()
	diff := newSize - wp.maxWorkers
	wp.maxWorkers = newSize
	wp.mu.Unlock()

	if diff > 0 {
		wp.startWorkers(diff) // Добавление новых воркеров
	} else if diff < 0 {
		wp.stopWorkers(-diff) // Остановка лишних воркеров
	}
}
