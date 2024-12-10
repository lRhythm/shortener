/*
Package logs - пользовательское логирование сервиса.
*/
package logs

import "github.com/sirupsen/logrus"

// New - создание объекта логирования.
func New() *logrus.Logger {
	var logger = logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	return logger
}
