package internal_sys

import (
	"net/http"
	"std-library/app/internal/scheduler"
	"std-library/app/internal/web/http"
	actionlog "std-library/app/log"
	"std-library/app/web/errors"
	"std-library/json"
	"std-library/logs"
	"std-library/nets"
	"strings"
)

type SchedulerController struct {
	accessControl *internal_http.IPv4AccessControl
	scheduler     *internal_scheduler.SchedulerImpl
}

func NewSchedulerController(scheduler *internal_scheduler.SchedulerImpl) *SchedulerController {
	return &SchedulerController{
		accessControl: &internal_http.IPv4AccessControl{},
		scheduler:     scheduler,
	}
}

func (c *SchedulerController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := c.accessControl.Validate(nets.IP(r).String())
	if err != nil {
		errors.Forbidden("access denied", "IP_ACCESS_DENIED")
	}

	if r.Method == http.MethodGet && r.URL.Path == "/_sys/job" {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		w.Write(json.Stringify(c.scheduler.JobsInfo()))
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/_sys/job/")
	parts := strings.SplitN(path, "/", 1)
	if len(parts) != 1 || r.Method != http.MethodPost {
		errors.NotFound("not found")
	}

	job := parts[0]
	ctx := r.Context()
	logs.WarnWithCtx(ctx, "[MANUAL_OPERATION] trigger job manually, job=%s", job)
	actionlog.Context(&ctx, "manual_operation", true)
	id := actionlog.GetId(&ctx)
	c.scheduler.TriggerNow(job, id)
	w.WriteHeader(202)

	w.Write([]byte("job triggered, job=" + job + ", id=" + id))
	return
}
