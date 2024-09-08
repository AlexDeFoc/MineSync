package globals

import (
	"fmt"
	"strings"
	"sync"
)

var (
	errors []string
	file   *string
	funcName *string
	mu     sync.Mutex
)

const (
	colorReset  = "\033[0m"
	colorError  = "\033[31m"
	colorInfo   = "\033[32m"
	colorDebug  = "\033[33m"
)

// SetFileName initializes or updates the global file name.
func SetFileName(name string) {
	if file == nil {
		file = new(string) // Initialize file pointer if it's nil
	}
	*file = name
}

// SetFuncName initializes or updates the global function name.
func SetFuncName(name string) {
	if funcName == nil {
		funcName = new(string) // Initialize funcName pointer if it's nil
	}
	*funcName = name
}

// Error logs an error message with context from file and function names if set.
func Error(text ...any) {
	mu.Lock()
	defer mu.Unlock()

	var builder strings.Builder
	for _, arg := range text {
		builder.WriteString(fmt.Sprint(arg, " | "))
	}
	newError := strings.TrimSuffix(builder.String(), " | ") // Remove trailing separator

	// Build the error message
	if file != nil {
		if funcName != nil {
			newError = fmt.Sprintf("%s %s: %s", *file, *funcName, newError)
		} else {
			newError = fmt.Sprintf("%s: %s", *file, newError)
		}
	}

	// Add the new error to the list, prepending it to keep the most recent first
	errors = append([]string{newError}, errors...)
}

// Debug logs a debug message.
func Debug(text ...any) {
    if LogLevel == "Debug" || LogLevel == "debug" {
        buildLog("DEBUG", text)
    }
}

// Info logs an info message.
func Info(text ...any) {
    if LogLevel == "Info" || LogLevel == "info" || LogLevel == "Debug" || LogLevel == "debug" {
        buildLog("INFO", text)
    }
}

// buildLog constructs and prints log messages for different levels.
func buildLog(level string, text ...any) {
	var builder strings.Builder

	for i, arg := range text {
		if i > 0 {
			builder.WriteString(" | ")
		}
		builder.WriteString(fmt.Sprint(arg))
	}

	switch level {
	case "DEBUG":
		fmt.Printf("%s%v:%s %v\n", colorDebug, level, colorReset, builder.String())
	case "INFO":
		fmt.Printf("%s%v:%s %v\n", colorInfo, level, colorReset, builder.String())
	}
}

// printErrors prints all errors with color.
func Print() {
	fmt.Printf("%sERRORS:%s\n", colorError, colorReset)
	for _, err := range errors {
		fmt.Printf("%s\n", err)
	}
    fmt.Scanln()
}
