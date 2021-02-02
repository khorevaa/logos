package logos

import (
	"errors"
	"testing"
	"time"
)

func TestJob_Event(t *testing.T) {
	type fields struct {
		Name      string
		emitter   Emitter
		Start     time.Time
		KeyValues map[string]string
	}
	type args struct {
		eventType string
		event     string
		status    CompletionStatus
		nanos     int64
		err       error
		KeyValues map[string]string
	}

	emit := New("job_emmiter").EventEmitter()

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"emit event",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "Event",
				event:     "get_users",
			},
		},
		{
			"emit event with kv",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "EventKv",
				event:     "get_users",
				KeyValues: map[string]string{
					"connect_string": "localhost:1545",
					"user":           "admin",
					"table":          "users",
				},
			},
		},
		{
			"emit event error",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "EventErr",
				event:     "get_users",
				err:       errors.New("event error"),
			},
		},
		{
			"emit event error with kv",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "EventErrKv",
				event:     "get_users",
				err:       errors.New("event error with kv"),
				KeyValues: map[string]string{
					"connect_string": "localhost:1545",
					"user":           "admin",
					"table":          "users",
				},
			},
		},
		{
			"emit event complete",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "Complete",
				event:     "get_users",
				err:       errors.New("event error"),
			},
		},
		{
			"emit event complete with kv",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "CompleteKv",
				event:     "get_users",
				err:       errors.New("event error with kv"),
				KeyValues: map[string]string{
					"connect_string": "localhost:1545",
					"user":           "admin",
					"table":          "users",
				},
			},
		},
		{
			"emit event timing",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "Timing",
				event:     "fetch_users",
				nanos:     54000,
			},
		},
		{
			"emit event timing with kv",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
			},
			args{
				eventType: "TimingKv",
				event:     "fetch_users",
				nanos:     54000,
				KeyValues: map[string]string{
					"connect_string": "localhost:1545",
					"raw_sql":        "select * from users",
					"table":          "users",
				},
			},
		},
		{
			"emit event job with kv",
			fields{
				Name:    "users_job",
				emitter: emit,
				Start:   time.Now(),
				KeyValues: map[string]string{
					"connect_string": "localhost:1545",
				},
			},
			args{
				eventType: "TimingKv",
				event:     "fetch_users",
				nanos:     54000,
				KeyValues: map[string]string{
					"raw_sql": "select * from users",
					"table":   "users",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			j := &Job{
				Name:      tt.fields.Name,
				emitter:   tt.fields.emitter,
				Start:     tt.fields.Start,
				KeyValues: tt.fields.KeyValues,
			}

			switch tt.args.eventType {
			case "Event":
				j.Event(tt.args.event)
			case "EventKv":
				j.EventKv(tt.args.event, tt.args.KeyValues)
			case "EventErr":
				j.EventErr(tt.args.event, tt.args.err)
			case "EventErrKv":
				j.EventErrKv(tt.args.event, tt.args.err, tt.args.KeyValues)
			case "Complete":
				j.Complete(tt.args.status)
			case "CompleteKv":
				j.CompleteKv(tt.args.status, tt.args.KeyValues)
			case "Timing":
				j.Timing(tt.args.event, tt.args.nanos)
			case "TimingKv":
				j.TimingKv(tt.args.event, tt.args.nanos, tt.args.KeyValues)
			case "Gauge":
				j.Gauge(tt.args.event, float64(tt.args.nanos))
			case "GaugeKv":
				j.GaugeKv(tt.args.event, float64(tt.args.nanos), tt.args.KeyValues)

			}

		})
	}
}
