package helper

import (
	"database/sql"
	"strings"
)

func EmptyStringIfNull(s string) string {
	if len(strings.TrimSpace(s)) == 0 {
		return ""
	}
	return s
}

func EmptyIntIfNull(s int) int {
	if s == 0 {
		return 0
	}
	return s
}

func GetNullStringValue(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
