/* errors.go
 *
 * Copyright (c) 2015, Michael Guzelevich <mguzelevich@gmail.com>
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */
package log4go

import (
	"fmt"
	"os"
)

type internalError struct {
	Id      string
	Message string
}

func (e internalError) Error() string {
	return fmt.Sprintf("Log4go error %s: %s", e.Id, e.Message)
}

type loadConfigurationError struct {
	Filename string
	Message  string
	Err      error
}

func (e loadConfigurationError) Error() string {
	return fmt.Sprintf("loadConfigurationError error: file=[%s]: %s [%s]", e.Filename, e.Message, e.Err)
}

type configurationFieldError struct {
	Message   string
	FieldName string
	Value     string
	Err       error
}

func (e configurationFieldError) Error() string {
	return fmt.Sprintf("error: [%s]: %s=%s]", e.Message, e.FieldName, e.Value, e.Err)
}

func checkFatalError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}
