package helpers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func ExtractIDFromPath(r *http.Request) (uint, error) {
	path := strings.Trim(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return 0, errors.New("invalid path format")
	}

	id, err := strconv.Atoi(parts[2])

	if err != nil || id <= 0 {
		return 0, errors.New("invalid ID format")
	}

	return uint(id), nil
}

func ParseLimitParam(r *http.Request, defaultValue, maxValue int) int {
	limitStr := r.URL.Query().Get("limit")
	if limitStr == "" {
		return defaultValue
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		return defaultValue
	}

	if limit <= 0 {
		return defaultValue
	}

	if limit > maxValue {
		return maxValue
	}

	return limit
}
