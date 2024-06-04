package internal_scheduler

import "std-library/app/scheduler"

type DisallowConcurrentJob struct {
	job scheduler.Job
}

func DisallowConcurrent(job scheduler.Job) DisallowConcurrentJob {
	return DisallowConcurrentJob{job: job}
}
