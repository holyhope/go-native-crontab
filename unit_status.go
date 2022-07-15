package god

import "context"

type unitStatus struct {
	exists bool
	loaded bool
}

func NewUnitStatus(exists bool, loaded bool) UnitStatus {
	return &unitStatus{
		exists: exists,
		loaded: loaded,
	}
}

func (status *unitStatus) Exists(ctx context.Context) bool {
	return status.exists
}

func (status *unitStatus) IsLoaded(ctx context.Context) bool {
	return status.loaded
}
