package retries

import (
	"fmt"
	"time"
)

type GenericRetry struct {
	Action            string      `json:"action"`
	MaxRetries        int         `json:"maxRetries"`
	Attempts          int         `json:"attempts"`
	LastAttemptTime   string      `json:"lastAttemptTime"`
	LastAttemptResult RetryResult `json:"lastAttemptResult"`
}

type RetryResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type Getter interface {
	GetRetries() []GenericRetry
}

type Setter interface {
	SetRetries([]GenericRetry)
}

type Retrier interface {
	Getter
	Setter
}

func NewRetry(name string, maxRetries int) GenericRetry {
	return GenericRetry{
		Action:     name,
		MaxRetries: maxRetries,
	}
}

func (r GenericRetry) CanRetry(obj Retrier) bool {
	var instance = findOrCreateRetry(obj, r)
	if instance != nil {
		if instance.MaxRetries == 0 {
			return true
		}

		if instance.Attempts < instance.MaxRetries {
			return true
		}
	}

	return false
}

func (r GenericRetry) Success(obj Retrier, message string) {
	r.set(obj, true, message)
}

func (r GenericRetry) Failure(obj Retrier, message string) {
	r.set(obj, false, message)
}

func (r GenericRetry) Successf(obj Retrier, message string, args ...any) {
	r.set(obj, true, fmt.Sprintf(message, args))
}

func (r GenericRetry) Failuref(obj Retrier, message string, args ...any) {
	r.set(obj, false, fmt.Sprintf(message, args))
}

func (r GenericRetry) set(obj Retrier, success bool, message string) {
	var instance = findOrCreateRetry(obj, r)
	if instance != nil {
		instance.Attempts += 1
		instance.LastAttemptTime = time.Now().String()
		instance.LastAttemptResult.Success = success
		instance.LastAttemptResult.Message = message
	}
}

func findOrCreateRetry(obj Retrier, r GenericRetry) *GenericRetry {
	for _, gr := range obj.GetRetries() {
		if gr.Action == r.Action {
			return &gr
		}
	}

	obj.SetRetries(append(obj.GetRetries(), r))

	return &obj.GetRetries()[len(obj.GetRetries())-1]
}
