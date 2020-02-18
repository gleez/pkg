This is an implementation of a workerpool which can get expanded &amp; shrink dynamically. Workers can get added when needed and get dismissed when no longer are needed. Of-course this workerpool can be used just as a simple one with a fixed size.

Examples can be seen inside documents.

```
var (
	MaxWorker = 2   //os.Getenv("MAX_WORKERS")
	MaxQueue  = 200 //os.Getenv("MAX_QUEUE")
)

// Start the dispatcher.
pool := worker.New(MaxWorker, MaxQueue)


// let's create a job with the payload
var job = func() {
    foo()
}

// Push the work onto the queue.
pool.Queue(job)
```
