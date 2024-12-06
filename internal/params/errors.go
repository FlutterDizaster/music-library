package params

import (
	"github.com/FlutterDizaster/music-library/internal/apperrors"
)

var (
	ErrInvalidDateValue = apperrors.Error{
		Message: "invalid date value",
	}
	ErrInvalidDateRange = apperrors.Error{
		Message: "invalid date range",
	}

	ErrFailedToParseParameters = apperrors.Error{
		Message: "failed to parse parameters",
	}

	ErrInvalidOffsetValue = apperrors.Error{
		Message: "invalid offset value",
	}
	ErrInvalidLimitValue = apperrors.Error{
		Message: "invalid limit value",
	}
)
