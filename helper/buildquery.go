package helper

import (
	"fmt"
	"strings"
)

func BuildInQuery(baseQuery string, ids []int) (string, []any) {
	if len(ids) == 0 {
		panic("BuildInQuery: ids tidak boleh kosong")
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := strings.Replace(
		baseQuery,
		"(?)",
		fmt.Sprintf("(%s)", strings.Join(placeholders, ",")),
		1,
	)

	return query, args
}

func BuildInQueryString(baseQuery string, ids []string) (string, []any) {
	if len(ids) == 0 {
		panic("BuildInQuery: ids tidak boleh kosong")
	}

	placeholders := make([]string, len(ids))
	args := make([]any, len(ids))

	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}

	query := strings.Replace(
		baseQuery,
		"(?)",
		fmt.Sprintf("(%s)", strings.Join(placeholders, ",")),
		1,
	)

	return query, args
}
