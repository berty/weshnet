package tyber

import (
	"context"
	"runtime/debug"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// types

type Step struct {
	ParentTraceID   string     `json:"parentTraceID"`
	Details         []Detail   `json:"details"`
	Status          StatusType `json:"status"`
	EndTrace        bool       `json:"endTrace"`
	UpdateTraceName string     `json:"updateTraceName"`
	ForceReopen     bool       `json:"forceReopen"`
}

type StepMutator = func(Step) Step

// zap format

func FormatStepLogFields(ctx context.Context, details []Detail, mutators ...StepMutator) []zapcore.Field {
	s := Step{
		ParentTraceID: GetTraceIDFromContext(ctx),
		Status:        Succeeded,
		Details:       details,
		EndTrace:      false,
		ForceReopen:   false,
	}
	for _, m := range mutators {
		s = m(s)
	}

	// Add debug if a there is no parent trace ID
	if s.ParentTraceID == noTraceID {
		s.Details = append(s.Details, Detail{Name: "StackTrace", Description: string(debug.Stack())})
	}

	return []zapcore.Field{
		zap.String("tyberLogType", string(StepType)),
		zap.Any("step", s),
	}
}

// constant mutators

func ForceReopen(s Step) Step {
	s.ForceReopen = true
	return s
}

func EndTrace(s Step) Step {
	s.EndTrace = true
	return s
}

func Fatal(s Step) Step {
	s.EndTrace = true
	s.Status = Failed
	return s
}

// variable mutators

func Status(st StatusType) StepMutator {
	return func(s Step) Step {
		s.Status = st
		return s
	}
}

func UpdateTraceName(newTitle string) StepMutator {
	return func(s Step) Step {
		s.UpdateTraceName = newTitle
		return s
	}
}

func WithDetail(name string, description string) StepMutator {
	return func(s Step) Step {
		s.Details = append(s.Details, Detail{Name: name, Description: description})
		return s
	}
}
