package propertytestingtemplate

import (
	"testing"

	"github.com/leanovate/gopter"
	"github.com/leanovate/gopter/gen"
	"github.com/leanovate/gopter/prop"
)

func PropertyTestingTemplate(t *testing.T) {
	parameters := gopter.DefaultTestParameters()
	properties := gopter.NewProperties(parameters)
	properties.Property("TestProperties", prop.ForAll(func(i int) bool {
		if i >= 0 {
			return true
		}
		t.Fail()
		return false
	}, gen.IntRange(0, 10)))
	properties.Run(gopter.ConsoleReporter(true))
}
