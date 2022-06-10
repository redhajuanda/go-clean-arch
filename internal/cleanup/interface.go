package cleanup

import "context"

type IService interface {
	CleanUpHTTPLog(ctx context.Context) error
}
