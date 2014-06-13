go-junit-report
===============

Converts `go test` output to an xml report, suitable for applications that
expect junit xml reports (e.g. [Jenkins](http://jenkins-ci.org)).


Changes to original
===============
This version aggregates all tests into one test suite in order to be valid XML-format and being usable with Teamcity


Installation
------------

	go get github.com/jstemmer/go-junit-report

Usage
-----

	go test -v | go-junit-report > report.xml

