package log

import (
	"crypto/sha256"
	"encoding/hex"

	"go.uber.org/zap"
)

// PIIMode indicates how to resolve PII fields in log statements.
type PIIMode uint8

const (
	// PIIModeNone indicates that PII fields shall be left as is.
	PIIModeNone PIIMode = 0

	// PIIModeHash indicates that the value part of a PII field shall
	// be hashed (SHA256). The key of the field stays untouched.
	PIIModeHash PIIMode = 1

	// PIIModeMask indicates that the value part of a PII field shall
	// be masked. If this mode is selected a mask function needs to be
	// provided under the MaskFunc property of this package. If no
	// MaskFunc is provided, PII fields will be omitted in the logs
	// using this mode.
	PIIModeMask PIIMode = 2

	// PIIModeRemove indicates that PII fields shall be omitted
	// completely from the final logs.
	PIIModeRemove PIIMode = 3
)

var (
	piiModes = map[PIIMode]struct{}{
		PIIModeNone:   {},
		PIIModeHash:   {},
		PIIModeMask:   {},
		PIIModeRemove: {},
	}

	// MaskFunc gets called on PII resolvers, when PII mode "mask" is chosen.
	// The function shall be thread-safe. When no function is provided, but
	// the mask PII mode is chosen, any PII fields will be omitted.
	MaskFunc func(key, value string) ResolvedPIIField
)

type field struct {
	key   string
	value string
}

func (f *field) resolve(piiMode PIIMode) zap.Field {
	switch piiMode {
	case PIIModeNone:
		return zap.String(f.key, f.value)
	case PIIModeHash:
		return zap.String(f.key, hash(f.value))
	case PIIModeMask:
		if MaskFunc == nil {
			return zap.Skip()
		}

		return MaskFunc(f.key, f.value).zapField()
	case PIIModeRemove:
		return zap.Skip()
	default:
		return zap.Skip()
	}
}

// PII is used to create standard PII field. When the field gets logged
// the actual PII is handled based on the current PII mode of the logger.
func PII(key, value string) *field {
	return &field{
		key:   key,
		value: value,
	}
}

// The CustomResolveFunc is passed to the CustomPII function of this
// package to handle the PII resolution in a customised way before a
// specific field gets logged.
type CustomResolveFunc func(mode PIIMode, key, value string) ResolvedPIIField

// CustomPII is used to create a PII field with a custom resolve function
// of type CustomResolveFunc. When the field gets logged the actual PII
// shall be handled appropriately by the custom function based on the
// given PII mode from the logger. The custom resolve function shall be
// thread-safe.
func CustomPII(key, value string, resolveFunc CustomResolveFunc) *customPIIField {
	if key == "" || value == "" || resolveFunc == nil {
		return nil
	}

	return &customPIIField{
		key:               key,
		value:             value,
		customResolveFunc: resolveFunc,
	}
}

type customPIIField struct {
	key               string
	value             string
	customResolveFunc CustomResolveFunc
}

func (f *customPIIField) resolve(piiMode PIIMode) zap.Field {
	return f.customResolveFunc(piiMode, f.key, f.value).zapField()
}

type ResolvedPIIField struct {
	Key   string
	Value string
}

func (f ResolvedPIIField) zapField() zap.Field {
	return zap.String(f.Key, f.Value)
}

func hash(in string) string {
	hashVal := sha256.Sum256([]byte(in))

	return hex.EncodeToString(hashVal[:])
}
