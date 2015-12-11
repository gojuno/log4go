/* config.go
 *
 * Copyright (c) 2015, Michael Guzelevich <mguzelevich@gmail.com>
 * Copyright (c) 2010, Kyle Lemons <kyle@kylelemons.net> (creator).
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */
package log4go

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

type FilterItem struct {
	Enabled    bool
	Tag        string
	Level      level
	Type       LoggerType
	Properties map[PropertyName]interface{}
}

type LoggerCfg struct {
	Filters []*FilterItem
}

func (fi *FilterItem) getString(p PropertyName) string {
	return fi.getProperty(p).(string)
}

func (fi *FilterItem) getInt(p PropertyName) int {
	return fi.getProperty(p).(int)
}

func (fi *FilterItem) getBool(p PropertyName) bool {
	return fi.getProperty(p).(bool)
}

func (fi *FilterItem) getProperty(p PropertyName) interface{} {
	//err = nil

	v, ok := fi.Properties[p]

	switch p {
	case FILENAME:
		if !ok {
			v = "log.log"
			//err = Error{Message: "filename property is required"}
		}
	case FORMAT:
		if !ok {
			v = "[%D %T] [%L] (%S) %M"
		}
	case MAX_LINES:
		if !ok {
			v = 0
		}
	case MAX_SIZE:
		if !ok {
			v = 0
		}
	case MAX_RECORDS:
		if !ok {
			v = 0
		}
	case DAILY:
		if !ok {
			v = false
		}
	case ROTATE:
		if !ok {
			v = false
		}
		// default:
		// 	err = Error{Message: fmt.Sprintf("Unknown property \"%s=%s\"", p, v)}
	}
	return v
}

func loadFile(filename string) ([]byte, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return []byte(""), loadConfigurationError{filename, "Could not open file for reading", err}
	}

	contents, err := ioutil.ReadAll(fd)
	if err != nil {
		return []byte(""), loadConfigurationError{filename, "Could not read file", err}
	}
	return contents, nil
}

// Load XML configuration; see examples/example.xml for documentation
func (log Logger) ApplyConfiguration(lc *LoggerCfg) error {
	var filter LogWriter
	for _, fi := range lc.Filters {

		switch fi.Type {
		case CONSOLE:
			filter = getConsoleLogWriter(fi)
		case FILE:
			filter = getFileLogWriter(fi)
		case XML:
			filter = getXmlLogWriter(fi)
		case SOCKET:
			filter = NewSocketLogWriter(fi.getString(PROTOCOL), fi.getString(ENDPOINT))
		}

		if !fi.Enabled {
			continue
		}

		log[fi.Tag] = &Filter{fi.Level, filter}
	}
	return nil
}

func getConsoleLogWriter(fi *FilterItem) LogWriter {
	clw := NewConsoleLogWriter()
	clw.SetFormat(fi.getString(FORMAT))
	return clw
}

func getFileLogWriter(fi *FilterItem) LogWriter {
	flw := NewFileLogWriter(fi.getString(FILENAME), fi.getBool(ROTATE))
	flw.SetFormat(fi.getString(FORMAT))
	flw.SetRotateLines(fi.getInt(MAX_LINES))
	flw.SetRotateSize(fi.getInt(MAX_SIZE))
	flw.SetRotateDaily(fi.getBool(DAILY))
	return flw
}

func getXmlLogWriter(fi *FilterItem) LogWriter {
	xlw := NewXMLLogWriter(fi.getString(FILENAME), fi.getBool(ROTATE))
	xlw.SetFormat(fi.getString(FORMAT))
	xlw.SetRotateLines(fi.getInt(MAX_LINES))
	xlw.SetRotateSize(fi.getInt(MAX_SIZE))
	xlw.SetRotateDaily(fi.getBool(DAILY))

	return xlw
}

// Load XML configuration; see examples/example.xml for documentation
func (log *Logger) LoadConfiguration(filename string) {
	log.Close()
	// Open the configuration file
	contents, err := loadFile(filename)
	checkFatalError(err)

	ext := filepath.Ext(filename)
	switch ext {
	case ".xml":
		err = log.loadXmlConfiguration(contents)
	case ".yaml":
		err = log.loadYamlConfiguration(contents)
	default:
		checkFatalError(internalError{Message: "unknown filename extention [" + ext + "]"})
	}
	checkFatalError(err)
}
