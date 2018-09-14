package aws

import (
	"testing"
)

func TestParseEc2ResourceId(t *testing.T) {
	type TestCase struct {
		ResourceId     string
		ExpectedPrefix string
		ExpectedUnique string
	}
	testCases := []TestCase{
		{
			ResourceId:     "",
			ExpectedPrefix: "",
			ExpectedUnique: "",
		},
		{
			ResourceId:     "vpc12345678",
			ExpectedPrefix: "",
			ExpectedUnique: "",
		},
		{
			ResourceId:     "vpc-12345678",
			ExpectedPrefix: "vpc",
			ExpectedUnique: "12345678",
		},
		{
			ResourceId:     "i-12345678abcdef00",
			ExpectedPrefix: "i",
			ExpectedUnique: "12345678abcdef00",
		},
	}

	for i, testCase := range testCases {
		resultPrefix, resultUnique := parseEc2ResourceId(testCase.ResourceId)
		if resultPrefix != testCase.ExpectedPrefix || resultUnique != testCase.ExpectedUnique {
			t.Errorf(
				"test case %d: got (%s, %s), but want (%s, %s)",
				i, resultPrefix, resultUnique, testCase.ExpectedPrefix, testCase.ExpectedUnique,
			)
		}
	}
}
