package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/bianjieai/irita/utils/protoidl"
)

const (
	outputPrivacy = "output_privacy"
	outputCached  = "output_cached"
	description   = "description"
)

// TODO
func ParseMethods(content string) (methods []string, err error) {
	return
}

// TODO New MethodToMethodProperty process
//func MethodToMethodProperty(index int, method string) (methodProperty MethodProperty, err sdk.Error) {
//	return
//}

func MethodToMethodProperty(index int, method protoidl.Method) (methodProperty MethodProperty, err error) {
	// set default value
	opp := NoPrivacy
	opc := NoCached

	var err1 error
	if _, ok := method.Attributes[outputPrivacy]; ok {
		opp, err1 = OutputPrivacyEnumFromString(method.Attributes[outputPrivacy])
		if err1 != nil {
			return methodProperty, sdkerrors.Wrap(ErrInvalidOutputPrivacyEnum, method.Attributes[outputPrivacy])
		}
	}
	if _, ok := method.Attributes[outputCached]; ok {
		opc, err1 = OutputCachedEnumFromString(method.Attributes[outputCached])
		if err1 != nil {
			return methodProperty, sdkerrors.Wrap(ErrInvalidOutputPrivacyEnum, method.Attributes[outputCached])
		}
	}
	methodProperty = MethodProperty{
		ID:            int16(index),
		Name:          method.Name,
		Description:   method.Attributes[description],
		OutputPrivacy: opp,
		OutputCached:  opc,
	}
	return
}
