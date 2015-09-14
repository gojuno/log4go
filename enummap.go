/* enummap.go
 *
 * Copyright (c) 2015, Michael Guzelevich <mguzelevich@gmail.com>
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */
package log4go

type enumMap struct {
	ab map[interface{}]interface{}
	ba map[interface{}]interface{}
}

func newEnumMap() *enumMap {
	return &enumMap{make(map[interface{}]interface{}), make(map[interface{}]interface{})}
}

func (m *enumMap) put(name, value interface{}) *enumMap {
	m.ab[name] = value
	m.ba[value] = name
	return m
}

func (m *enumMap) value(name interface{}) (value interface{}, exists bool) {
	value, exists = m.ab[name]
	return
}

func (m *enumMap) name(value interface{}) (name interface{}, exists bool) {
	name, exists = m.ba[value]
	return
}
