/* constants.go
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
	"strconv"
	"strings"
)

type LoggerType int

const (
	CONSOLE LoggerType = iota
	FILE
	XML
	SOCKET
)

type PropertyName int

const (
	FILENAME PropertyName = iota
	ROTATE
	FORMAT
	MAX_LINES
	MAX_SIZE
	MAX_RECORDS
	DAILY
	ENDPOINT
	PROTOCOL
)

var loggingLevels = newEnumMap()
var loggerTypes = newEnumMap()
var properties = newEnumMap()

func init() {
	loggingLevels.put(FINEST, "FINEST")
	loggingLevels.put(FINE, "FINE")
	loggingLevels.put(DEBUG, "DEBUG")
	loggingLevels.put(TRACE, "TRACE")
	loggingLevels.put(INFO, "INFO")
	loggingLevels.put(WARNING, "WARNING")
	loggingLevels.put(ERROR, "ERROR")
	loggingLevels.put(CRITICAL, "CRITICAL")

	loggerTypes.put(CONSOLE, "console")
	loggerTypes.put(FILE, "file")
	loggerTypes.put(XML, "xml")
	loggerTypes.put(SOCKET, "socket")

	properties.put(FILENAME, "filename")
	properties.put(ROTATE, "rotate")
	properties.put(FORMAT, "format")
	properties.put(MAX_LINES, "maxlines")
	properties.put(MAX_SIZE, "maxsize")
	properties.put(MAX_RECORDS, "maxrecords")
	properties.put(DAILY, "daily")
	properties.put(ENDPOINT, "endpoint")
	properties.put(PROTOCOL, "protocol")
}

func stringToLevel(levelString string) (lvl level, err error) {
	lvl = DEBUG
	err = nil

	v, ok := loggingLevels.name(levelString)
	if !ok {
		err = internalError{Message: levelString}
	} else {
		lvl = v.(level)
	}
	return
}

func stringToType(typeString string) (lType LoggerType, err error) {
	lType = CONSOLE
	err = nil

	v, ok := loggerTypes.name(typeString)
	if !ok {
		err = internalError{Message: typeString}
	} else {
		lType = v.(LoggerType)
	}
	return
}

func stringToPropertyName(p string) (property PropertyName, err error) {
	property = FILENAME
	err = nil

	v, ok := properties.name(p)
	if !ok {
		err = internalError{Message: p}
	} else {
		property = v.(PropertyName)
	}
	return
}

// Parse a number with K/M/G suffixes based on thousands (1000) or 2^10 (1024)
func strToNumSuffix(str string, mult int) int {
	num := 1
	if len(str) > 1 {
		switch str[len(str)-1] {
		case 'G', 'g':
			num *= mult
			fallthrough
		case 'M', 'm':
			num *= mult
			fallthrough
		case 'K', 'k':
			num *= mult
			str = str[0 : len(str)-1]
		}
	}
	parsed, _ := strconv.Atoi(str)
	return parsed * num
}

func stringToPropertyValue(p PropertyName, v string) (value interface{}, err error) {
	value = nil
	err = nil

	v = strings.Trim(v, " \r\n")
	switch p {
	case FILENAME:
		value = v
	case FORMAT:
		value = v
	case MAX_LINES:
		value = strToNumSuffix(v, 1000)
	case MAX_SIZE:
		value = strToNumSuffix(v, 1024)
	case MAX_RECORDS:
		value = strToNumSuffix(v, 1000)
	case DAILY:
		value = v != "false"
	case ROTATE:
		value = v != "false"
	case ENDPOINT:
		value = v
	case PROTOCOL:
		value = v
	default:
		err = internalError{Message: fmt.Sprintf("Unknown property \"%s=%s\"", p, v)}
	}

	return
}
