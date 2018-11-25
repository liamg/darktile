package version

import (
	"reflect"
	"testing"
)

func TestSort(t *testing.T) {
	testcases := []struct {
		input  []string
		output []string
	}{
		{
			input: []string{
				"Package-0.4.tar.gz",
				"Package-0.1.tar.gz",
				"Package-0.10.1.tar.gz",
				"Package-0.10.tar.gz",
				"Package-0.2.tar.gz",
				"Package-0.3.1.tar.gz",
				"Package-0.3.2.tar.gz",
				"Package-0.3.tar.gz",
			},
			output: []string{
				"Package-0.1.tar.gz",
				"Package-0.2.tar.gz",
				"Package-0.3.tar.gz",
				"Package-0.3.1.tar.gz",
				"Package-0.3.2.tar.gz",
				"Package-0.4.tar.gz",
				"Package-0.10.tar.gz",
				"Package-0.10.1.tar.gz",
			},
		},
		{
			input: []string{
				"1.0-dev",
				"1.0a1",
				"1.0b1",
				"1.0RC1",
				"1.0rc1",
				"1.0",
				"1.0pl1",
				"1.1-dev",
				"1.2",
				"1.10",
			},
			output: []string{
				"1.0pl1",
				"1.0-dev",
				"1.0a1",
				"1.0b1",
				"1.0RC1",
				"1.0rc1",
				"1.0",
				"1.1-dev",
				"1.2",
				"1.10",
			},
		},
		{
			input: []string{
				"v1.0",
				"1.0.1",
				"dev-master",
			},
			output: []string{
				"v1.0",
				"1.0.1",
				"dev-master",
			},
		},
	}

	for _, testcase := range testcases {
		Sort(testcase.input)
		if !reflect.DeepEqual(testcase.input, testcase.output) {
			t.Errorf("Expected output %+v did not match actual %+v", testcase.output, testcase.input)
		}
	}
}
