/*
Copyright (c) 2020 Sergey Anisimov

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

package extracting

import (
	"bytes"
	"regexp"
	"strings"
)

var (
	notAlphaNum   = regexp.MustCompile("[^a-zA-Z0-9]+")
	matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")
)

func getMetricName(compName, capName, attrName, unitSuffix string) string {
	var buffer bytes.Buffer
	buffer.WriteString("smartthings_")
	if compName != "main" {
		buffer.WriteString(toMetricName(compName))
		buffer.WriteString("_")
	}
	buffer.WriteString(toMetricName(capName))
	buffer.WriteString("_")
	buffer.WriteString(toMetricName(attrName))
	if len(unitSuffix) != 0 {
		buffer.WriteString("_")
		buffer.WriteString(unitSuffix)
	}
	return buffer.String()
}

func toMetricName(deviceName string) string {
	deviceName = strings.ToLower(toSnakeCase(deviceName))
	return notAlphaNum.ReplaceAllString(deviceName, "_")
}

func toSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
