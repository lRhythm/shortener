package models

import "errors"

// ErrConflict - ошибка при конфликте данных, возникающая при операциях с хранилищем.
var ErrConflict = errors.New("conflict")
