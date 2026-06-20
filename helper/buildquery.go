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

func BuildInQueryWithArgs(baseQuery string, ids []int, extraArgs ...any) (string, []any) {
	if len(ids) == 0 {
		panic("BuildInQuery: ids tidak boleh kosong")
	}

	placeholders := make([]string, len(ids))
	args := make([]any, 0, len(ids)+len(extraArgs))

	for i, id := range ids {
		placeholders[i] = "?"
		args = append(args, id)
	}

	query := strings.Replace(
		baseQuery,
		"(?)",
		fmt.Sprintf("(%s)", strings.Join(placeholders, ",")),
		1,
	)

	// append arg tambahan (misal tahun)
	args = append(args, extraArgs...)

	return query, args
}
