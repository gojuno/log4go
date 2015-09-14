/* config_xml.go
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
	"encoding/xml"
)

type xmlProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type xmlFilter struct {
	Enabled  string        `xml:"enabled,attr"`
	Tag      string        `xml:"tag"`
	Level    string        `xml:"level"`
	Type     string        `xml:"type"`
	Property []xmlProperty `xml:"property"`
}

type xmlLoggerConfig struct {
	Filter []xmlFilter `xml:"filter"`
}

func xmlNewFilterCfg(enabled string, tag string, fType string, lvl string) (*FilterItem, error) {
	f := FilterItem{
		Enabled:    enabled != "false",
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

	return &f, nil
}

func xmlToConfiguration(xc *xmlLoggerConfig) (*LoggerCfg, error) {
	lc := new(LoggerCfg)
	for _, xmlfilt := range xc.Filter {
		f, err := xmlNewFilterCfg(xmlfilt.Enabled, xmlfilt.Tag, xmlfilt.Type, xmlfilt.Level)
		if err != nil {
			return nil, err
		}

		lc.Filters = append(lc.Filters, f)
		for _, p := range xmlfilt.Property {
			pName, err := stringToPropertyName(p.Name)
			if err != nil {
				return nil, err
			}
			value, err := stringToPropertyValue(pName, p.Value)
			if err != nil {
				return nil, err
			}
			f.Properties[pName] = value

		}
	}
	return lc, nil
}

// Load configuration from XML content
func (log *Logger) loadXmlConfiguration(contents []byte) error {
	xc := new(xmlLoggerConfig)

	// TODO: replace errors to typed

	err := xml.Unmarshal(contents, xc)
	if err != nil {
		// fmt.Fprintf(os.Stderr, "LoadConfiguration: Error: Could not parse XML configuration in %q: %s\n", filename, err)
		return err
	}

	lc, err := xmlToConfiguration(xc)
	if err != nil {
		return err
	}

	err = log.ApplyConfiguration(lc)
	if err != nil {
		return err
	}
	return nil
}
