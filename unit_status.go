package god

import "context"

type unitStatus struct {
	exists  bool
	enabled bool
}

func NewUnitStatus(exists bool, enabled bool) UnitStatus {
	return &unitStatus{
		exists:  exists,
		enabled: enabled,
	}
}

func (status *unitStatus) Exists(ctx context.Context) bool {
	return status.exists
}

func (status *unitStatus) IsEnabled(ctx context.Context) bool {
	return status.enabled
}
