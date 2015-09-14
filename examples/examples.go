/* examples.go
 *
 * Copyright (c) 2015, Michael Guzelevich <mguzelevich@gmail.com>
 * Copyright (c) 2010, Kyle Lemons <kyle@kylelemons.net> (creator).
 * All rights reserved.
 *
 * This software may be modified and distributed under the terms
 * of the New BSD license.  See the LICENSE file for details.
 */
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	l4g "github.com/mguzelevich/log4go"
)

const (
	filename = "test.runtime.log"
)

func TestConsoleLogWriter() {
	log := l4g.NewLogger()
	log.AddFilter("stdout", l4g.DEBUG, l4g.NewConsoleLogWriter())
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
}

func TestFileLogWriter() {
	// Get a new logger instance
	log := l4g.NewLogger()

	// Create a default logger that is logging messages of FINE or higher
	log.AddFilter("file", l4g.FINE, l4g.NewFileLogWriter(filename, false))
	log.Close()

	/* Can also specify manually via the following: (these are the defaults) */
	flw := l4g.NewFileLogWriter(filename, false)
	flw.SetFormat("[%D %T] [%L] (%S) %M")
	flw.SetRotate(false)
	flw.SetRotateSize(0)
	flw.SetRotateLines(0)
	flw.SetRotateDaily(false)
	log.AddFilter("file", l4g.FINE, flw)

	// Log some experimental messages
	log.Finest("runtime: Everything is created now (notice that I will not be printing to the file)")
	log.Info("runtime: The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))
	log.Critical("runtime: Time to close out!")

	// Close the log
	log.Close()

	// Print what was logged to the file (yes, I know I'm skipping error checking)
	fd, _ := os.Open(filename)
	in := bufio.NewReader(fd)
	fmt.Print("runtime: Messages logged to file were: (line numbers not included)\n")
	for lineno := 1; ; lineno++ {
		line, err := in.ReadString('\n')
		if err == io.EOF {
			break
		}
		fmt.Printf("%3d:\t%s", lineno, line)
	}
	fd.Close()

	// Remove the file so it's not lying around
	os.Remove(filename)
}

func TestSocketLogWriter() {
	log := l4g.NewLogger()
	log.AddFilter("network", l4g.FINEST, l4g.NewSocketLogWriter("udp", "192.168.1.255:12124"))

	// Run `nc -u -l -p 12124` or similar before you run this to see the following message
	log.Info("The time is now: %s", time.Now().Format("15:04:05 MST 2006/01/02"))

	// This makes sure the output stream buffer is written
	log.Close()
}

func TestLoadXmlConfiguration() {
	// Load the configuration (isn't this easy?)
	l4g.LoadConfiguration("example.xml")

	// And now we're ready!
	l4g.Finest("XML: This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("XML: Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("XML: About that time, eh chaps?")
}

func TestLoadYamlConfiguration() {
	// Load the configuration (isn't this easy?)
	l4g.LoadConfiguration("example.yaml")

	// And now we're ready!
	l4g.Finest("YAML: This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("YAML: Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("YAML: About that time, eh chaps?")
}

func TestLoadYamlCustomRootConfiguration() {
	// Load the configuration (isn't this easy?)
	l4g.YamlConfigRoot = "my_app.custom.debug.logging"
	l4g.LoadConfiguration("example.root.yaml")

	// And now we're ready!
	l4g.Finest("YAML custom root: This will only go to those of you really cool UDP kids!  If you change enabled=true.")
	l4g.Debug("YAML custom root: Oh no!  %d + %d = %d!", 2, 2, 2+2)
	l4g.Info("YAML custom root: About that time, eh chaps?")
}

func main() {
	TestLoadXmlConfiguration()
	TestLoadYamlConfiguration()
	TestLoadYamlCustomRootConfiguration()
	TestConsoleLogWriter()
	TestFileLogWriter()
	TestSocketLogWriter()
}
