package eeslism

import (
	"testing"
)

func Test_StandardPlan(t *testing.T) {
	t.Skip("Skipping complex test - PCM sample requires additional setup")
	Entry("../samples/standard-plan-no-hcap-PCM-CM-fsolm.txt", "../Base")
}
