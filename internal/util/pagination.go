package util

import (
	"fmt"
	"net/http"
	"strconv"
)

func ParsePagination(r *http.Request, defaultLimit, defaultOffset int) (limit, offset int, err error) {
	query := r.URL.Query()

	limit = defaultLimit
	offset = defaultOffset

	if lStr := query.Get("limit"); lStr != "" {
		limit, err = strconv.Atoi(lStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid limit parameter: %v", err)
		}
	}

	if oStr := query.Get("offset"); oStr != "" {
		offset, err = strconv.Atoi(oStr)
		if err != nil {
			return 0, 0, fmt.Errorf("invalid offset parameter: %v", err)
		}
	}

	return limit, offset, nil
}
