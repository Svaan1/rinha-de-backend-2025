package queue

type Worker struct {
	ID         int
	JobChannel chan Job
	WorkerPool chan chan Job
	QuitChan   chan bool
}

func NewWorker(id int, workerPool chan chan Job) Worker {
	return Worker{
		ID:         id,
		JobChannel: make(chan Job),
		WorkerPool: workerPool,
		QuitChan:   make(chan bool),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			w.WorkerPool <- w.JobChannel

			select {
			case job := <-w.JobChannel:
				job.Execute()
			case <-w.QuitChan:
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	go func() {
		w.QuitChan <- true
	}()
}
