package service

import "errors"

var (
	ErrTaskNotFound  = errors.New("task not found")
	ErrInvalidStatus = errors.New("invalid task status")
	ErrTitleRequired = errors.New("title is required")
	ErrEmptyIDList   = errors.New("ids cannot be empty")
)
