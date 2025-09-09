package dtos

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_dtos_NewResponse(t *testing.T) {
	baseResponse := NewResponse(200, "success", 2)
	expected := BaseResponse{
		Status:  200,
		Message: "success",
		Data:    2,
	}
	require.Equal(t, expected, baseResponse, "NewResponse() return value not as expected")
}
