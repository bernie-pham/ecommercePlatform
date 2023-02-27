package gapi

import (
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func FieldViolation(field string, err error) *errdetails.BadRequest_FieldViolation {
	return &errdetails.BadRequest_FieldViolation{
		Field:       field,
		Description: err.Error(),
	}
}

func InvalidArgumentError(violations []*errdetails.BadRequest_FieldViolation) error {
	badReq := &errdetails.BadRequest{FieldViolations: violations}
	statusInvalid := status.New(codes.InvalidArgument, "invalid paramaters")

	statusDetails, err := statusInvalid.WithDetails(badReq)
	if err != nil {
		return statusInvalid.Err()
	}
	return statusDetails.Err()
}

func MerchantEmailNotFound(field string) error {
	statusNotFound := status.New(codes.NotFound, "wrong email or not exist")
	notFoundReq := &errdetails.BadRequest{
		FieldViolations: []*errdetails.BadRequest_FieldViolation{
			&errdetails.BadRequest_FieldViolation{Field: field, Description: "not found"}}}
	statusDetails, err := statusNotFound.WithDetails(notFoundReq)
	if err != nil {
		return statusNotFound.Err()
	}
	return statusDetails.Err()
}
