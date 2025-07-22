package queue

type Queue struct {
	WorkerPool chan chan Job
	MaxWorkers int
	JobQueue   chan Job
}

func New(maxJobs, maxWorkers int) Queue {
	workerPool := make(chan chan Job, maxWorkers)
	jobQueue := make(chan Job, maxJobs)
	return Queue{
		WorkerPool: workerPool,
		MaxWorkers: maxWorkers,
		JobQueue:   jobQueue,
	}
}

func (q *Queue) Enqueue(job Job) {
	q.JobQueue <- job
}

func (q *Queue) dispatch() {
	for job := range q.JobQueue {
		jobChannel := <-q.WorkerPool
		jobChannel <- job
	}
}

func (q *Queue) Run() {
	for i := 1; i <= q.MaxWorkers; i++ {
		worker := NewWorker(i, q.WorkerPool)
		worker.Start()
	}

	go q.dispatch()
}
