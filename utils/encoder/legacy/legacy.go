package legacy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	math_bits "math/bits"

	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

// Legacy is an old encoder type used by Prysm before they moved to SSZ completely.
// Older versions of KeyVault use this encoding.
type Legacy struct {
}

func New() *Legacy {
	return &Legacy{}
}

func (l *Legacy) Encode(obj interface{}) ([]byte, error) {
	switch t := obj.(type) {
	case *eth.AttestationData:
		return LegacyAttestationDataMarshal(t)
	case *eth.BeaconBlock:
		return LegacyBeaconBlockMarshal(t)
	case *validatorpb.SignRequest:
		return json.Marshal(t)
	}
	return nil, errors.New("type not supported")
}

func (l *Legacy) Decode(data []byte, v interface{}) error {
	switch t := v.(type) {
	case *eth.AttestationData:
		return LegacyAttestationDataUnMarshal(t, data)
	case *eth.BeaconBlock:
		return LegacyBeaconBlockUnMarshal(t, data)
	case *validatorpb.SignRequest:
		return json.Unmarshal(data, t)
	}
	return errors.New("type not supported")
}

func sovKeymanager(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozKeymanager(x uint64) (n int) {
	return sovKeymanager(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}

func encodeVarintKeymanager(dAtA []byte, offset int, v uint64) int {
	offset -= sovKeymanager(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}

func skipKeymanager(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowKeymanager
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowKeymanager
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowKeymanager
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthKeymanager
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupKeymanager
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthKeymanager
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}
