package documents

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
)

func Test_a(t *testing.T) {

}

//Test_PolicyReadSamples - read in sample JSON polcies and ensure they are correctly
//parsed into JSON and then Protobuffer formats
func Test_PolicyReadSamples(t *testing.T) {
	var policy *Policy
	var err error
	policy, err = ValidateJSONPolicyDoc("single.json")
	assert.Nil(t, err, "Error parsing JSON ")
	assert.NotNil(t, policy, "Policy should not be nil")
	assert.Equal(t, int64(7), policy.ParticipantCount, "Participant count incorrect")

	policy, err = ValidateJSONPolicyDoc("manager.json")
	assert.Nil(t, err, "Error parsing JSON ")
	assert.NotNil(t, policy, "Policy should not be nil")

	policy, err = ValidateJSONPolicyDoc("one_sg_two_of_three.json")
	assert.Nil(t, err, "Error parsing JSON ")
	assert.NotNil(t, policy, "Policy should not be nil")

	policy, err = ValidateJSONPolicyDoc("t_equals_p.json")
	assert.Nil(t, err, "Error parsing JSON ")
	assert.NotNil(t, policy, "Policy should not be nil")

	policy, err = ValidateJSONPolicyDoc("three_groups.json")
	assert.Nil(t, err, "Error parsing JSON ")
	assert.NotNil(t, policy, "Policy should not be nil")

	policy, err = ValidateJSONPolicyDoc("two_sg.json")
	assert.Nil(t, err, "Error parsing JSON ")
	assert.NotNil(t, policy, "Policy should not be nil")

	policy, err = ValidateJSONPolicyDoc("bad.json")
	assert.NotNil(t, err, "No Error parsing JSON ")
	print(policy)
}

func ValidateJSONPolicyDoc(filename string) (*Policy, error) {
	filepath := fmt.Sprintf("test-policy-documents/%s", filename)
	dat, err := ioutil.ReadFile(filepath)
	if err != nil {
		return &Policy{}, err
	}
	pol := &PolicyWrapper{}
	err = jsonpb.UnmarshalString(string(dat), pol)

	if err != nil {
		return nil, err
	}
	return pol.Policy, err
}
