package application

import (
	"aviasales/internal/services/storage"
	"aviasales/pkg/entities"
	"aviasales/pkg/logger"
	"context"
	"encoding/xml"
	"os"
	"sync"
)

type WorkerParserQueue struct {
	Ctx      context.Context
	FileName string
	Storage  storage.IStorage
}

// SpawnWorkers creates <count> workers that process messages.
func SpawnWorkers(
	workersLimit int,
) *WorkerPool {
	ret := &WorkerPool{
		ch: make(chan WorkerParserQueue, workersLimit),
		wg: &sync.WaitGroup{},
	}

	ret.wg.Add(workersLimit)

	for i := 0; i < workersLimit; i++ {
		w := worker{}
		go w.runLoop(ret.ch, ret.wg)
	}
	return ret
}

// WorkerPool processes Daemon SendMessage's
// until Flush() is called.
type WorkerPool struct {
	wg        *sync.WaitGroup
	ch        chan WorkerParserQueue
	isClosing sync.Mutex
	isClosed  bool
}

func (w *WorkerPool) Flush() {
	w.isClosing.Lock()
	defer w.isClosing.Unlock()

	if !w.isClosed {
		close(w.ch)
		w.wg.Wait()
		w.isClosed = true
	}
}

func (w *WorkerPool) Put(queue *WorkerParserQueue) {
	w.ch <- *queue
}

type worker struct{}

func (w *worker) runLoop(
	in <-chan WorkerParserQueue,
	wg *sync.WaitGroup,
) {
	defer wg.Done()

	for queue := range in {
		w.processQueue(
			queue.Ctx,
			queue.FileName,
			queue.Storage,
		)
	}
}

func (w *worker) processQueue(
	rootCtx context.Context,
	fileName string,
	storage storage.IStorage,
) {
	ctx := logger.With(rootCtx, "fileName", fileName)

	xmlFile, err := os.Open(fileName)
	if err != nil {
		logger.Error(ctx, "unable to read file", err)
		return
	}
	defer func() {
		_ = xmlFile.Close()
	}()

	decoder := xml.NewDecoder(xmlFile)

	counter := 0
	var inElement string

loop:
	for {
		select {
		case <-ctx.Done():
			logger.Info(ctx, "worker canceled", "count", counter)
			return
		default:
			t, _ := decoder.Token()
			if t == nil {
				break loop
			}

			switch se := t.(type) {
			case xml.StartElement:
				inElement = se.Name.Local
				if inElement == "Flights" {
					var itinerary entities.Itinerary
					_ = decoder.DecodeElement(&itinerary, &se)

					storage.AddItinerary(itinerary)
					counter++
				}
			default:

			}
		}
	}

	logger.Info(ctx, "added itineraries", "count", counter)
}
