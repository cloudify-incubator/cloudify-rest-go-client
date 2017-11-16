/*
Copyright (c) 2017 GigaSpaces Technologies Ltd. All rights reserved

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

//printBottomLine - print "+-...-+" line as bottom
func printBottomLine(columnSizes []int) {
	fmt.Printf("+")
	for _, size := range columnSizes {
		fmt.Print(strings.Repeat("-", size+2))
		fmt.Printf("+")
	}
	fmt.Printf("\n")
}

//printLine - print "| <text> |" from text columns/lines
func printLine(columnSizes []int, lines []string) {
	fmt.Printf("|")
	for col, size := range columnSizes {
		fmt.Print(" " + lines[col] + " ")
		fmt.Print(strings.Repeat(" ", size-utf8.RuneCountInString(lines[col])))
		fmt.Printf("|")
	}
	fmt.Printf("\n")
}

//PrintTable - print table with column titles and several lines
func PrintTable(titles []string, lines [][]string) {
	columnSizes := make([]int, len(titles))

	// column title sizes
	for col, name := range titles {
		if columnSizes[col] < utf8.RuneCountInString(name) {
			columnSizes[col] = utf8.RuneCountInString(name)
		}
	}

	// column value sizes
	for _, values := range lines {
		for col, name := range values {
			if columnSizes[col] < utf8.RuneCountInString(name) {
				columnSizes[col] = utf8.RuneCountInString(name)
			}
		}
	}

	printBottomLine(columnSizes)
	// titles
	printLine(columnSizes, titles)
	printBottomLine(columnSizes)
	// lines
	for _, values := range lines {
		printLine(columnSizes, values)
	}
	printBottomLine(columnSizes)
}

//CliArgumentsList - return clean list of arguments and options
func CliArgumentsList(osArgs []string) (arguments []string, options []string) {
	for pos, str := range osArgs {
		if str[:1] == "-" {
			return osArgs[:pos], osArgs[pos:]
		}
	}
	return osArgs, []string{}
}

//DirZipArchive - create archive from directory and return as bytes array
func DirZipArchive(parentDir string) ([]byte, error) {
	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)

	// Create a new zip archive.
	w := zip.NewWriter(buf)

	log.Printf("Looking into %s", parentDir)
	errWalk := filepath.Walk(parentDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			f, errCreate := w.Create("parent/" + path[len(parentDir):])
			if errCreate != nil {
				return errCreate
			}

			content, errRead := ioutil.ReadFile(path)
			if errRead != nil {
				return errRead
			}

			_, errWrite := f.Write(content)
			if errWrite != nil {
				return errWrite
			}
			log.Printf("Attached: %s", path[len(parentDir):])
		}
		return nil
	})

	if errWalk != nil {
		return nil, errWalk
	}

	// Make sure to check the error on Close.
	errZip := w.Close()
	if errZip != nil {
		return nil, errZip
	}
	return buf.Bytes(), nil
}

//InList - return true if string is already in list
func InList(source []string, value string) bool {
	for _, inList := range source {
		if inList == value {
			return true
		}
	}
	return false
}
