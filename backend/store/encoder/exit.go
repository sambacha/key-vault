package encoder

import (
	"fmt"
	"io"

	types "github.com/prysmaticlabs/eth2-types"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func LegacySignedExitUnMarshal(m *eth.SignedVoluntaryExit, dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBeaconBlock
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
			return fmt.Errorf("proto: SignedVoluntaryExit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignedVoluntaryExit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBeaconBlock
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
				return ErrInvalidLengthBeaconBlock
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthBeaconBlock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if m.Exit == nil {
				m.Exit = &eth.VoluntaryExit{}
			}
			if err := LegacyExitUnMarshal(m.Exit, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBeaconBlock
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
				return ErrInvalidLengthBeaconBlock
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthBeaconBlock
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipBeaconBlock(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBeaconBlock
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBeaconBlock
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

func LegacyExitUnMarshal(m *eth.VoluntaryExit, dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowBeaconBlock
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
			return fmt.Errorf("proto: VoluntaryExit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: VoluntaryExit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epoch", wireType)
			}
			m.Epoch = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBeaconBlock
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Epoch |= types.Epoch(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ValidatorIndex", wireType)
			}
			m.ValidatorIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBeaconBlock
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ValidatorIndex |= types.ValidatorIndex(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipBeaconBlock(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthBeaconBlock
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthBeaconBlock
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

func voluntary_marshalToSizedBuffer(m *eth.VoluntaryExit, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if m.ValidatorIndex != 0 {
		i = encodeVarintBeaconBlock(dAtA, i, uint64(m.ValidatorIndex))
		i--
		dAtA[i] = 0x10
	}
	if m.Epoch != 0 {
		i = encodeVarintBeaconBlock(dAtA, i, uint64(m.Epoch))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func signedVoluntary_marshalToSizedBuffer(m *eth.SignedVoluntaryExit, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintBeaconBlock(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x12
	}
	if m.Exit != nil {
		{
			size, err := voluntary_marshalToSizedBuffer(m.Exit, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func signedVoluntary_size(m *eth.SignedVoluntaryExit) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Exit != nil {
		l = voluntary_size(m.Exit)
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}

func voluntary_size(m *eth.VoluntaryExit) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Epoch != 0 {
		n += 1 + sovBeaconBlock(uint64(m.Epoch))
	}
	if m.ValidatorIndex != 0 {
		n += 1 + sovBeaconBlock(uint64(m.ValidatorIndex))
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}
