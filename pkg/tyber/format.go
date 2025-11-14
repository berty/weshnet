package tyber

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Detail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func safeJSONMarshal(v interface{}) string {
	bs, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return fmt.Sprintf(`"%s"`, err.Error())
	}
	return string(bs)
}

func JSONDetail(name string, val interface{}) Detail {
	return Detail{Name: name, Description: safeJSONMarshal(val)}
}

func WithJSONDetail(name string, val interface{}) StepMutator {
	return func(s Step) Step {
		s.Details = append(s.Details, JSONDetail(name, val))
		return s
	}
}

func WithError(err error) StepMutator {
	return func(s Step) Step {
		s.Details = append(s.Details, Detail{Name: "Error", Description: err.Error()})
		return s
	}
}

type LogType string

const (
	TraceType     LogType = "trace"
	StepType      LogType = "step"
	EventType     LogType = "event"
	SubscribeType LogType = "subcribe"
)

var KnownLogTypes = []LogType{TraceType, StepType, EventType, SubscribeType}

func (lt LogType) IsKnown() bool {
	return slices.Contains(KnownLogTypes, lt)
}

type StatusType string

const (
	Running   StatusType = "running"
	Succeeded StatusType = "succeeded"
	Failed    StatusType = "failed"
)

type Trace struct {
	TraceID string `json:"traceID"`
}

type Event struct {
	Details []Detail `json:"details"`
}

func FormatTraceLogFields(ctx context.Context) []zapcore.Field {
	return []zapcore.Field{
		zap.String("tyberLogType", string(TraceType)),
		zap.Any("trace", Trace{
			TraceID: GetTraceIDFromContext(ctx),
		}),
	}
}

func FormatEventLogFields(_ context.Context, details []Detail) []zapcore.Field {
	return []zapcore.Field{
		zap.String("tyberLogType", string(EventType)),
		zap.Any("event", Event{
			Details: details,
		}),
	}
}

func ZapFieldsToDetails(fields ...zap.Field) []Detail {
	dets := make([]Detail, len(fields))
	for i, field := range fields {
		dets[i] = Detail{Name: field.Key, Description: field.String}
	}
	return dets
}
