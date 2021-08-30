package legacy

import (
	"fmt"
	"io"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func LegacyETH1UnMarshal(m *eth.Eth1Data, dAtA []byte) error {
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
			return fmt.Errorf("proto: Eth1Data: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Eth1Data: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositRoot", wireType)
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
			m.DepositRoot = append(m.DepositRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.DepositRoot == nil {
				m.DepositRoot = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field DepositCount", wireType)
			}
			m.DepositCount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowBeaconBlock
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.DepositCount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHash", wireType)
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
			m.BlockHash = append(m.BlockHash[:0], dAtA[iNdEx:postIndex]...)
			if m.BlockHash == nil {
				m.BlockHash = []byte{}
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

func eth1data_marshalToSizedBuffer(m *eth.Eth1Data, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if len(m.BlockHash) > 0 {
		i -= len(m.BlockHash)
		copy(dAtA[i:], m.BlockHash)
		i = encodeVarintBeaconBlock(dAtA, i, uint64(len(m.BlockHash)))
		i--
		dAtA[i] = 0x1a
	}
	if m.DepositCount != 0 {
		i = encodeVarintBeaconBlock(dAtA, i, uint64(m.DepositCount))
		i--
		dAtA[i] = 0x10
	}
	if len(m.DepositRoot) > 0 {
		i -= len(m.DepositRoot)
		copy(dAtA[i:], m.DepositRoot)
		i = encodeVarintBeaconBlock(dAtA, i, uint64(len(m.DepositRoot)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func eth1data_size(m *eth.Eth1Data) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DepositRoot)
	if l > 0 {
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	if m.DepositCount != 0 {
		n += 1 + sovBeaconBlock(uint64(m.DepositCount))
	}
	l = len(m.BlockHash)
	if l > 0 {
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}
