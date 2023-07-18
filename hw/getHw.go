////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Package hw contains files for identifying the hardware inside a computer

package hw

import (
	"fmt"
	"os/exec"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

const (
	// Initial header print, used for all other prints using fmt.Printf(header, s).
	header = "---------------------%s---------------------"

	// Beginning and ending headers for hardware information.
	begin = "Printing Hardware Information"
	end   = "End of Hardware Information"

	// All info headers to be printed.
	cpu       = "CPU INFO"
	gpu       = "GPU INFO"
	partition = "PARTITION INFO"
	diskUsage = "DISK USAGE INFO"
	diskHw    = "DISK HW INFO"
	ramUsage  = "RAM USAGE INFO"
)

// commandList is the list of headers that will be run and printed to the log.
// commandList will be iterated over in LogHardware rather than commandMap to
// ensure consistent printing order. The entry will be used to pull out the bash
// command to be run out of commandMap.
var commandList = []string{cpu, gpu, partition, diskUsage, diskHw, ramUsage}

// commandMap maps the header that will be printed to the bash command that will
// be run.
var commandMap = map[string][]string{
	cpu:       {"lscpu"},
	gpu:       {"bash", "-c", "lspci -vnnn | perl -lne 'print if /^\\d+\\:.+(\\[\\S+\\:\\S+\\])/' | grep VGA"},
	partition: {"lsblk"},
	diskUsage: {"df", "-h"},
	diskHw:    {"lshw", "-class", "disk", "-class", "storage"},
	ramUsage:  {"free", "-mt"},
}

// LogHardware iterates over commandList, running the command
// from the commandMap and printing out the results.
func LogHardware() {
	jww.INFO.Printf(header, begin)

	for _, cmdHeader := range commandList {
		// Pull command from map
		cmdList := commandMap[cmdHeader]

		// Run command
		var out []byte
		var err error
		cmd := strings.Join(cmdList, " ")
		if len(cmdList) == 1 {
			// Handle quirks of exec.Command and variable command length
			out, err = exec.Command(cmdList[0]).Output()
			if err != nil {
				// Print error, continue with other commands
				out = []byte((fmt.Sprintf("%s err: %s", cmdList, err)))
			}
		} else {
			out, err = exec.Command(cmdList[0], cmdList[1:]...).Output()
			if err != nil {
				// Print error, continue with other commands
				out = []byte((fmt.Sprintf("%s err: %s", cmdList, err)))
			}
		}

		jww.INFO.Printf("%s\r\n%s\r\n%s", fmt.Sprintf(header, cmdHeader), cmd, out)

	}

	jww.INFO.Printf(header, end)

	return
}
