package legacy

import (
	"fmt"
	"io"

	types "github.com/prysmaticlabs/eth2-types"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func LegacyAggregationAttAndProofUnMarshal(m *eth.AggregateAttestationAndProof, dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAttestation
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: AggregateAttestationAndProof: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AggregateAttestationAndProof: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field AggregatorIndex", wireType)
			}
			m.AggregatorIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.AggregatorIndex |= types.ValidatorIndex(uint64(b&0x7F) << shift)
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SelectionProof", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SelectionProof = append(m.SelectionProof[:0], dAtA[iNdEx:postIndex]...)
			if m.SelectionProof == nil {
				m.SelectionProof = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Aggregate", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthAttestation
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthAttestation
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Aggregate == nil {
				m.Aggregate = &eth.Attestation{}
			}
			if err := LegacyAttestationUnMarshal(m.Aggregate, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAttestation(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthAttestation
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthAttestation
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			//m.XXX_unrecognized = append(m.XXX_unrecognized, dAtA[iNdEx:iNdEx+skippy]...)
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}

func LegacyAggregationAttAndProofMarshal(m *eth.AggregateAttestationAndProof) (dAtA []byte, err error) {
	size := aggAttProof_size(m)
	dAtA = make([]byte, size)
	n, err := aggAttProof_marshalToSizedBuffer(m, dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func aggAttProof_size(m *eth.AggregateAttestationAndProof) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.AggregatorIndex != 0 {
		n += 1 + sovAttestation(uint64(m.AggregatorIndex))
	}
	l = len(m.SelectionProof)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	if m.Aggregate != nil {
		l = attestation_size(m.Aggregate)
		n += 1 + l + sovAttestation(uint64(l))
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}

func aggAttProof_marshalToSizedBuffer(m *eth.AggregateAttestationAndProof, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if m.Aggregate != nil {
		{
			size, err := attestation_marshalToSizedBuffer(m.Aggregate, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAttestation(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SelectionProof) > 0 {
		i -= len(m.SelectionProof)
		copy(dAtA[i:], m.SelectionProof)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.SelectionProof)))
		i--
		dAtA[i] = 0x12
	}
	if m.AggregatorIndex != 0 {
		i = encodeVarintAttestation(dAtA, i, uint64(m.AggregatorIndex))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}
