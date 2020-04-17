package status

import (
	"go-web-demo/library/ecode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// togRPCCode convert ecode.Codo to gRPC code
func TogRPCCode(code ecode.Codes) codes.Code {
	switch code.Code() {
	case ecode.OK.Code():
		return codes.OK
	case ecode.RequestErr.Code():
		return codes.InvalidArgument
	case ecode.NotFound.Code():
		return codes.NotFound
	case ecode.Unauthorized.Code():
		return codes.Unauthenticated
	case ecode.AccessDenied.Code():
		return codes.PermissionDenied
	case ecode.LimitExceed.Code():
		return codes.ResourceExhausted
	case ecode.MethodNotAllowed.Code():
		return codes.Unimplemented
	case ecode.Deadline.Code():
		return codes.DeadlineExceeded
	case ecode.ServiceUnavailable.Code():
		return codes.Unavailable
	}
	return codes.Unknown
}

func ToECode(gst *status.Status) ecode.Code {
	gcode := gst.Code()
	switch gcode {
	case codes.OK:
		return ecode.OK
	case codes.InvalidArgument:
		return ecode.RequestErr
	case codes.NotFound:
		return ecode.NotFound
	case codes.PermissionDenied:
		return ecode.AccessDenied
	case codes.Unauthenticated:
		return ecode.Unauthorized
	case codes.ResourceExhausted:
		return ecode.LimitExceed
	case codes.Unimplemented:
		return ecode.MethodNotAllowed
	case codes.DeadlineExceeded:
		return ecode.Deadline
	case codes.Unavailable:
		return ecode.ServiceUnavailable
	case codes.Unknown:
		return ecode.String(gst.Message())
	}
	return ecode.ServerErr
}
