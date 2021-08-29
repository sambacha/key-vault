package encoder

import (
	"fmt"
	"io"

	types "github.com/prysmaticlabs/eth2-types"
	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func LegacyAttestationDataUnMarshal(m *eth.AttestationData, dAtA []byte) error {
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
			return fmt.Errorf("proto: AttestationData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AttestationData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slot", wireType)
			}
			m.Slot = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Slot |= types.Slot(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CommitteeIndex", wireType)
			}
			m.CommitteeIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAttestation
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CommitteeIndex |= types.CommitteeIndex(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BeaconBlockRoot", wireType)
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
			m.BeaconBlockRoot = append(m.BeaconBlockRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.BeaconBlockRoot == nil {
				m.BeaconBlockRoot = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Source", wireType)
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
			if m.Source == nil {
				m.Source = &eth.Checkpoint{}
			}
			if err := checkpoint_unmarshal(m.Source, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Target", wireType)
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
			if m.Target == nil {
				m.Target = &eth.Checkpoint{}
			}
			if err := checkpoint_unmarshal(m.Target, dAtA[iNdEx:postIndex]); err != nil {
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

func LegacyAttestationDataMarshal(m *eth.AttestationData) ([]byte, error) {
	size := attestationData_size(m)
	dAtA := make([]byte, size)
	n, err := attestationData_marshalToSizedBuffer(m, dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func attestationData_size(m *eth.AttestationData) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Slot != 0 {
		n += 1 + sovAttestation(uint64(m.Slot))
	}
	if m.CommitteeIndex != 0 {
		n += 1 + sovAttestation(uint64(m.CommitteeIndex))
	}
	l = len(m.BeaconBlockRoot)
	if l > 0 {
		n += 1 + l + sovAttestation(uint64(l))
	}
	if m.Source != nil {
		l = checkpoint_size(m.Source)
		n += 1 + l + sovAttestation(uint64(l))
	}
	if m.Target != nil {
		l = checkpoint_size(m.Target)
		n += 1 + l + sovAttestation(uint64(l))
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}

func attestationData_marshalToSizedBuffer(m *eth.AttestationData, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if m.Target != nil {
		{
			size, err := checkpoint_marshalToSizedBuffer(m.Target, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAttestation(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x2a
	}
	if m.Source != nil {
		{
			size, err := checkpoint_marshalToSizedBuffer(m.Source, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintAttestation(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x22
	}
	if len(m.BeaconBlockRoot) > 0 {
		i -= len(m.BeaconBlockRoot)
		copy(dAtA[i:], m.BeaconBlockRoot)
		i = encodeVarintAttestation(dAtA, i, uint64(len(m.BeaconBlockRoot)))
		i--
		dAtA[i] = 0x1a
	}
	if m.CommitteeIndex != 0 {
		i = encodeVarintAttestation(dAtA, i, uint64(m.CommitteeIndex))
		i--
		dAtA[i] = 0x10
	}
	if m.Slot != 0 {
		i = encodeVarintAttestation(dAtA, i, uint64(m.Slot))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}
