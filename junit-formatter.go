package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"runtime"
	"strings"
	"errors"
)

type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

type JUnitTestCase struct {
	XMLName   xml.Name      `xml:"testcase"`
	Classname string        `xml:"classname,attr"`
	Name      string        `xml:"name,attr"`
	Time      string        `xml:"time,attr"`
	Failure   *JUnitFailure `xml:"failure,omitempty"`
}

type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

func NewJUnitProperty(name, value string) JUnitProperty {
	return JUnitProperty{
		Name:  name,
		Value: value,
	}
}

// JUnitReportXML writes a junit xml representation of the given report to w
// in the format described at http://windyroad.org/dl/Open%20Source/JUnit.xsd
func JUnitReportXML(report *Report, w io.Writer) error {
	suites := []JUnitTestSuite{}
	if len(report.Packages) <= 0 {
		return errors.New("No report found")
	}
	pkg := report.Packages[0]
	packageName := pkg.Name
	packageName = packageName[:strings.LastIndex(packageName, "/")]
	ts := JUnitTestSuite{
		Tests:      0,
		Failures:   0,
		Time:       "",
		Name:       packageName,
		Properties: []JUnitProperty{},
		TestCases:  []JUnitTestCase{},
	}

	// properties
	ts.Properties = append(ts.Properties, NewJUnitProperty("go.version", runtime.Version()))
	var time = 0
	// convert Report to JUnit test suites
	for _, pkgCurrent := range report.Packages {
		classname := pkgCurrent.Name
		if idx := strings.LastIndex(classname, "/"); idx > -1 && idx < len(pkgCurrent.Name) {
			classname = pkgCurrent.Name[idx+1:]
		}
		time+=pkgCurrent.Time

		// individual test cases
		for _, test := range pkgCurrent.Tests {
			ts.Tests = ts.Tests+1
			testCase := JUnitTestCase{
				Classname: classname,
				Name:      test.Name,
				Time:      formatTime(test.Time),
				Failure:   nil,
			}

			if test.Result == FAIL {
				ts.Failures += 1

				testCase.Failure = &JUnitFailure{
					Message:  "Failed",
					Type:     "",
					Contents: strings.Join(test.Output, "\n"),
				}
			}

			ts.TestCases = append(ts.TestCases, testCase)
		}
	}
	ts.Time = formatTime(time)
	suites = append(suites, ts)

	// to xml
	bytes, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(w)

	// remove newline from xml.Header, because xml.MarshalIndent starts with a newline
	writer.WriteString(xml.Header[:len(xml.Header) - 1])
	writer.WriteByte('\n')
	writer.Write(bytes)
	writer.WriteByte('\n')
	writer.Flush()

	return nil
}

func countFailures(tests []Test) (result int) {
	for _, test := range tests {
		if test.Result == FAIL {
			result += 1
		}
	}
	return
}

func formatTime(time int) string {
	return fmt.Sprintf("%.3f", float64(time)/1000.0)
}
