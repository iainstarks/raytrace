package main

import (
	"github.com/cucumber/godog"
	"os"
	"flag"
	"testing"
	"github.com/cucumber/godog/colors"
)

// go test support for godog
var opt = godog.Options{Output: colors.Uncolored(os.Stdout), Format:"progress",}

func init() {
	godog.BindFlags("godog.", flag.CommandLine, &opt)
}

func TestMain(m *testing.M) {
	flag.Parse()
	opt.Paths = flag.Args()

//	status := godog.RunWithOptions("godogs", func(s *godog.Suite) {
//		FeatureContext(s)
//	}, opt)

	status := godog.TestSuite{
		Name: "raytrace",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer: InitializeScenario,
		Options: &opt,
	}.Run()
	if st := m.Run(); st > status {
		status = st
	}
	os.Exit(status)
}

func InitializeTestSuite(context *godog.TestSuiteContext) {
	context.BeforeSuite(func() {
		// initialize the global state here
		 })
}

func InitializeScenario(context *godog.ScenarioContext) {
	context.BeforeScenario(func(*godog.Scenario) {})
}

func aIsAPoint() error {
	return godog.ErrPending
}

func aIsAVector() error {
	return godog.ErrPending
}

func aIsNotAPoint() error {
	return godog.ErrPending
}

func aIsNotAVector() error {
	return godog.ErrPending
}

func aTuple(arg1, arg2, arg3, arg4, arg5, arg6, arg7, arg8 int) error {
	return godog.ErrPending
}

func aw(arg1, arg2 int) error {
	return godog.ErrPending
}

func ax(arg1, arg2 int) error {
	return godog.ErrPending
}

func ay(arg1, arg2 int) error {
	return godog.ErrPending
}

func az(arg1, arg2 int) error {
	return godog.ErrPending
}

func FeatureContext(s *godog.Suite) {
	s.Step(`^a is a point$`, aIsAPoint)
	s.Step(`^a is a vector$`, aIsAVector)
	s.Step(`^a is not a point$`, aIsNotAPoint)
	s.Step(`^a is not a vector$`, aIsNotAVector)
	s.Step(`^a ← tuple\((\d+)\.(\d+), -(\d+)\.(\d+), (\d+)\.(\d+), (\d+)\.(\d+)\)$`, aTuple)
	s.Step(`^a\.w = (\d+)\.(\d+)$`, aw)
	s.Step(`^a\.x = (\d+)\.(\d+)$`, ax)
	s.Step(`^a\.y = -(\d+)\.(\d+)$`, ay)
	s.Step(`^a\.z = (\d+)\.(\d+)$`, az)
}
/*
** Scenario: A tuple with w=1.0 is a point
** Given a ← tuple(4.3, -4.2, 3.1, 1.0)
** Then a.x = 4.3
** And a.y = -4.2
** And a.z = 3.1
** And a.w = 1.0
** And a is a point
** And a is not a vector
*/