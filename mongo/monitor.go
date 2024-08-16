package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/event"
	"log/slog"
	actionlog "std-library/app/log"
	"strings"
	"time"
)

var defaultSlowOperation = 1 * time.Second.Nanoseconds()
var Monitor = &event.CommandMonitor{
	Succeeded: func(ctx context.Context, event *event.CommandSucceededEvent) {
		elapsed := event.CommandFinishedEvent.Duration.Nanoseconds()
		actionlog.Stat(&ctx, "mongo."+strings.ToLower(event.CommandName), float64(elapsed))

		if elapsed > defaultSlowOperation {
			actionlog.Context(&ctx, "slow_operation", true)
			slog.WarnContext(ctx, fmt.Sprintf("[SLOW_OPERATION] slow %s, duration %v, db: %s", event.CommandName, event.CommandFinishedEvent.Duration, event.DatabaseName))
		}
	},
	Failed: func(ctx context.Context, event *event.CommandFailedEvent) {
		elapsed := event.CommandFinishedEvent.Duration.Nanoseconds()
		actionlog.Stat(&ctx, "mongo."+strings.ToLower(event.CommandName), float64(elapsed))

		if elapsed > defaultSlowOperation {
			actionlog.Context(&ctx, "slow_operation", true)
			slog.WarnContext(ctx, fmt.Sprintf("[SLOW_OPERATION] slow %s, duration %v, db: %s", event.CommandName, event.CommandFinishedEvent.Duration, event.DatabaseName))
		}
	},
}

func SetDefaultSlowOperation(d time.Duration) {
	defaultSlowOperation = d.Nanoseconds()
}
