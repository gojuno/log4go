/* config_yaml.go
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
	"strings"

	"gopkg.in/yaml.v2"
)

type yamlFilterProperties map[string]string

type yamlFilter struct {
	Enabled    bool                 "enabled"
	Type       string               "type"
	Level      string               "level"
	Properties yamlFilterProperties ",flow"
}

type yamlLoggerConfig struct {
	Logging map[string]yamlFilter ",flow"
}

func unmarshalYamlSible(contents []byte, startWith string) (*yamlLoggerConfig, error) {
	// err := yaml.Unmarshal(contents, yc)
	// if err != nil {
	// 	// fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
	// 	return err
	// }

	m := map[interface{}]interface{}{}
	if err := yaml.Unmarshal(contents, &m); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse YAML configuration: %s\n", err)
		return nil, err
	}

	var key interface{}
	var value interface{}
	value = m
	for _, k := range strings.Split(startWith, ".") {
		v, _ := value.(map[interface{}]interface{})[k]
		value = v.(interface{})
		key = k
		// if ok {
		//   panic("panik")
		// }
	}

	r, err := yaml.Marshal(map[interface{}]interface{}{key: value})
	if err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse YAML configuration: %s\n", err)
		return nil, err
	}

	yc := new(yamlLoggerConfig)
	if err := yaml.Unmarshal(r, yc); err != nil {
		fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse YAML configuration: %s\n", err)
		return nil, err
	}
	return yc, nil
}

func yamlNewFilterCfg(enabled bool, tag string, fType string, lvl string, properties yamlFilterProperties) (*FilterItem, error) {
	f := FilterItem{
		Enabled:    enabled,
		Tag:        tag,
		Properties: map[PropertyName]interface{}{},
	}

	t, err := stringToType(fType)
	if err != nil {
		return nil, configurationFieldError{
			"could not parse logger type",
			"type",
			fType,
			err,
		}
	}
	f.Type = t

	l, err := stringToLevel(lvl)
	if err != nil {
		return nil, configurationFieldError{
			"could not parse logging level",
			"level",
			lvl,
			err,
		}
	}
	f.Level = l

	for pKey, pValue := range properties {
		pName, err := stringToPropertyName(pKey)
		if err != nil {
			return nil, err
		}
		value, err := stringToPropertyValue(pName, pValue)
		if err != nil {
			return nil, err
		}
		f.Properties[pName] = value
	}

	return &f, nil
}

func yamlToConfiguration(yc *yamlLoggerConfig) (*LoggerCfg, error) {
	lc := new(LoggerCfg)
	for tag, desc := range yc.Logging {
		f, err := yamlNewFilterCfg(desc.Enabled, tag, desc.Type, desc.Level, desc.Properties)
		if err != nil {
			return nil, err
		}
		lc.Filters = append(lc.Filters, f)
	}
	return lc, nil
}

// Load configuration from YAML content
func (log *Logger) loadYamlConfiguration(contents []byte) error {
	// TODO: replace errors to typed
	yc, err := unmarshalYamlSible(contents, YamlConfigRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[log4go] unmarshalYamlSible error\n")
		return err
	}

	lc, err := yamlToConfiguration(yc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[log4go] yamlToConfiguration error\n")
		return err
	}

	err = log.ApplyConfiguration(lc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[log4go] ApplyConfiguration error\n")
		return err
	}
	return nil
}
