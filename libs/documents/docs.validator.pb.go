// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: docs.proto

package documents

import (
	fmt "fmt"
	math "math"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/mwitkow/go-proto-validators"
	_ "github.com/gogo/protobuf/gogoproto"
	regexp "regexp"
	github_com_mwitkow_go_proto_validators "github.com/mwitkow/go-proto-validators"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

var _regex_SignedEnvelope_SignerCID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)

func (this *SignedEnvelope) Validate() error {
	if !(len(this.Signature) > 20) {
		return github_com_mwitkow_go_proto_validators.FieldError("Signature", fmt.Errorf(`value '%v' must have a length greater than '20'`, this.Signature))
	}
	if !_regex_SignedEnvelope_SignerCID.MatchString(this.SignerCID) {
		return github_com_mwitkow_go_proto_validators.FieldError("SignerCID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.SignerCID))
	}
	return nil
}
func (this *Envelope) Validate() error {
	if this.Header != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.Header); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("Header", err)
		}
	}
	return nil
}
func (this *Header) Validate() error {
	if !(this.DateTime > 1564050341) {
		return github_com_mwitkow_go_proto_validators.FieldError("DateTime", fmt.Errorf(`value '%v' must be greater than '1564050341'`, this.DateTime))
	}
	if !(this.DateTime < 32521429541) {
		return github_com_mwitkow_go_proto_validators.FieldError("DateTime", fmt.Errorf(`value '%v' must be less than '32521429541'`, this.DateTime))
	}
	if len(this.Recipients) > 20 {
		return github_com_mwitkow_go_proto_validators.FieldError("Recipients", fmt.Errorf(`value '%v' must contain at most 20 elements`, this.Recipients))
	}
	for _, item := range this.Recipients {
		if item != nil {
			if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(item); err != nil {
				return github_com_mwitkow_go_proto_validators.FieldError("Recipients", err)
			}
		}
	}
	return nil
}

var _regex_Recipient_CID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)

func (this *Recipient) Validate() error {
	if !_regex_Recipient_CID.MatchString(this.CID) {
		return github_com_mwitkow_go_proto_validators.FieldError("CID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.CID))
	}
	return nil
}
func (this *IDDocument) Validate() error {
	if !(this.Timestamp > 1564050341) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be greater than '1564050341'`, this.Timestamp))
	}
	if !(this.Timestamp < 32521429541) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be less than '32521429541'`, this.Timestamp))
	}
	return nil
}

var _regex_OrderDocument_PrincipalCID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)
var _regex_OrderDocument_BeneficiaryCID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)

func (this *OrderDocument) Validate() error {
	if !(this.Coin > -1) {
		return github_com_mwitkow_go_proto_validators.FieldError("Coin", fmt.Errorf(`value '%v' must be greater than '-1'`, this.Coin))
	}
	if !(this.Coin < 999) {
		return github_com_mwitkow_go_proto_validators.FieldError("Coin", fmt.Errorf(`value '%v' must be less than '999'`, this.Coin))
	}
	if !_regex_OrderDocument_PrincipalCID.MatchString(this.PrincipalCID) {
		return github_com_mwitkow_go_proto_validators.FieldError("PrincipalCID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.PrincipalCID))
	}
	if !_regex_OrderDocument_BeneficiaryCID.MatchString(this.BeneficiaryCID) {
		return github_com_mwitkow_go_proto_validators.FieldError("BeneficiaryCID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.BeneficiaryCID))
	}
	if this.Reference == "" {
		return github_com_mwitkow_go_proto_validators.FieldError("Reference", fmt.Errorf(`value '%v' must not be an empty string`, this.Reference))
	}
	if !(this.Timestamp > 1564050341) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be greater than '1564050341'`, this.Timestamp))
	}
	if !(this.Timestamp < 32521429541) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be less than '32521429541'`, this.Timestamp))
	}
	if this.OrderPart2 != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.OrderPart2); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("OrderPart2", err)
		}
	}
	if this.OrderPart3 != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.OrderPart3); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("OrderPart3", err)
		}
	}
	if this.OrderPart4 != nil {
		if err := github_com_mwitkow_go_proto_validators.CallValidatorIfExists(this.OrderPart4); err != nil {
			return github_com_mwitkow_go_proto_validators.FieldError("OrderPart4", err)
		}
	}
	return nil
}

var _regex_OrderPart2_PreviousOrderCID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)

func (this *OrderPart2) Validate() error {
	if !_regex_OrderPart2_PreviousOrderCID.MatchString(this.PreviousOrderCID) {
		return github_com_mwitkow_go_proto_validators.FieldError("PreviousOrderCID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.PreviousOrderCID))
	}
	if !(this.Timestamp > 1564050341) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be greater than '1564050341'`, this.Timestamp))
	}
	if !(this.Timestamp < 32521429541) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be less than '32521429541'`, this.Timestamp))
	}
	return nil
}

var _regex_OrderPart3_PreviousOrderCID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)

func (this *OrderPart3) Validate() error {
	if !_regex_OrderPart3_PreviousOrderCID.MatchString(this.PreviousOrderCID) {
		return github_com_mwitkow_go_proto_validators.FieldError("PreviousOrderCID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.PreviousOrderCID))
	}
	if !(this.Timestamp > 1564050341) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be greater than '1564050341'`, this.Timestamp))
	}
	if !(this.Timestamp < 32521429541) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be less than '32521429541'`, this.Timestamp))
	}
	return nil
}

var _regex_OrderPart4_PreviousOrderCID = regexp.MustCompile(`^Q[[:alnum:]]{45}$|^$`)

func (this *OrderPart4) Validate() error {
	if !_regex_OrderPart4_PreviousOrderCID.MatchString(this.PreviousOrderCID) {
		return github_com_mwitkow_go_proto_validators.FieldError("PreviousOrderCID", fmt.Errorf(`value '%v' must be a string conforming to regex "^Q[[:alnum:]]{45}$|^$"`, this.PreviousOrderCID))
	}
	if !(this.Timestamp > 1564050341) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be greater than '1564050341'`, this.Timestamp))
	}
	if !(this.Timestamp < 32521429541) {
		return github_com_mwitkow_go_proto_validators.FieldError("Timestamp", fmt.Errorf(`value '%v' must be less than '32521429541'`, this.Timestamp))
	}
	return nil
}
func (this *Policy) Validate() error {
	return nil
}
func (this *PlainTestMessage1) Validate() error {
	return nil
}
func (this *EncryptTestMessage1) Validate() error {
	return nil
}
func (this *SimpleString) Validate() error {
	return nil
}
