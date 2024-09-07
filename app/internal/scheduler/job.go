package internal_scheduler

import "github.com/odycenter/std-library/app/scheduler"

type DisallowConcurrentJob struct {
	job scheduler.Job
}

func DisallowConcurrent(job scheduler.Job) DisallowConcurrentJob {
	return DisallowConcurrentJob{job: job}
}
