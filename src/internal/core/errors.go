package core

import "errors"

var (
	ErrRepoNotFound = errors.New("repo: not found")
	ErrNotFound     = errors.New("not found")
)
