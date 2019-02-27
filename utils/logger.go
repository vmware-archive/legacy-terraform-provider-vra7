package utils

import (
	"os"
	"path/filepath"

	logging "github.com/op/go-logging"
)

var (
	terraformVraProviderFileName = "vra-terraform.log"
	format                       = logging.MustStringFormatter(
		`%{color}%{time:2006-01-02T15:04:05.999Z07:00} %{shortfile} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
	)
)

// InitLog - initializes the log
func InitLog() {
	initTerraformVraProviderLog()
}

func initTerraformVraProviderLog() {
	//Setup terraform service log
	var backendList = []logging.Backend{}
	var logFile *os.File
	logFile, err := os.OpenFile("."+string(filepath.Separator)+terraformVraProviderFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		f, err2 := os.OpenFile(terraformVraProviderFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err2 != nil {
			panic("Error opening log file for write " + err2.Error())
		} else {
			logFile = f
		}
	}
	fileBackend := logging.NewLogBackend(logFile, "", 0)
	fileBackendFormatter := logging.NewBackendFormatter(fileBackend, format)

	// info and more severe messages should be sent to file backend
	fileBackendLeveled := logging.AddModuleLevel(fileBackendFormatter)
	fileBackendLeveled.SetLevel(logging.INFO, "")
	backendList = append(backendList, fileBackendLeveled)

	// debug and more severe messages should be send to console backend, log format is default
	consoleBackend := logging.NewLogBackend(os.Stderr, "", 0)
	consoleBackendFormatter := logging.NewBackendFormatter(consoleBackend, logging.DefaultFormatter)
	consoleBackendLeveled := logging.AddModuleLevel(consoleBackendFormatter)
	consoleBackendLeveled.SetLevel(logging.DEBUG, "")
	backendList = append(backendList, consoleBackendLeveled)

	logging.SetBackend(backendList...)
}
