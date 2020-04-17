package status

import (
	"fmt"
	"go-web-demo/library/ecode"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCodeConvert(t *testing.T) {
	var table = map[codes.Code]ecode.Code{
		codes.OK: ecode.OK,
		// codes.Canceled
		codes.Unknown:          ecode.ServerErr,
		codes.InvalidArgument:  ecode.RequestErr,
		codes.DeadlineExceeded: ecode.Deadline,
		codes.NotFound:         ecode.NotFound,
		// codes.AlreadyExists
		codes.PermissionDenied:  ecode.AccessDenied,
		codes.ResourceExhausted: ecode.LimitExceed,
		// codes.FailedPrecondition
		// codes.Aborted
		// codes.OutOfRange
		codes.Unimplemented: ecode.MethodNotAllowed,
		codes.Unavailable:   ecode.ServiceUnavailable,
		// codes.DataLoss
		codes.Unauthenticated: ecode.Unauthorized,
	}
	for k, v := range table {
		assert.Equal(t, ToECode(status.New(k, "-500")), v)
	}
	for k, v := range table {
		assert.Equal(t, TogRPCCode(v), k, fmt.Sprintf("togRPC code error: %d -> %d", v, k))
	}
}

func TestNoDetailsConvert(t *testing.T) {
	gst := status.New(codes.Unknown, "-2233")
	assert.Equal(t, ToECode(gst).Code(), -2233)

	gst = status.New(codes.Internal, "")
	assert.Equal(t, ToECode(gst).Code(), -500)
}
