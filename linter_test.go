package sundrylint

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzerTimeParse(t *testing.T) {
	testCases := []struct {
		desc     string
		settings LinterSetting
	}{
		{
			desc:     "timeparse",
			settings: LinterSetting{},
		},
		{
			desc:     "iteroverzero",
			settings: LinterSetting{},
		},
		{
			desc:     "funcresultunused",
			settings: LinterSetting{},
		},
		{
			desc:     "rangeappendall",
			settings: LinterSetting{},
		},
	}

	for _, test := range testCases[3:] {
		t.Run(test.desc, func(t *testing.T) {
			a, err := NewAnalyzer(test.settings)
			if err != nil {
				t.Fatal(err)
			}

			analysistest.Run(t, analysistest.TestData(), a, test.desc)
		})
	}
}
