package legacy

import (
	"fmt"
	"io"

	types "github.com/prysmaticlabs/eth2-types"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
	validatorpb "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1/validator-client"
)

func LegacySignRequestUnMarshal(m *validatorpb.SignRequest, dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowKeymanager
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
			return fmt.Errorf("proto: SignRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: SignRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PublicKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PublicKey = append(m.PublicKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PublicKey == nil {
				m.PublicKey = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SigningRoot", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SigningRoot = append(m.SigningRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.SigningRoot == nil {
				m.SigningRoot = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SignatureDomain", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SignatureDomain = append(m.SignatureDomain[:0], dAtA[iNdEx:postIndex]...)
			if m.SignatureDomain == nil {
				m.SignatureDomain = []byte{}
			}
			iNdEx = postIndex
		case 101:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Block", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &eth.BeaconBlock{}
			if err := LegacyBeaconBlockUnMarshal(v, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Object = &validatorpb.SignRequest_Block{v}
			iNdEx = postIndex
		case 102:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AttestationData", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &eth.AttestationData{}
			if err := LegacyAttestationDataUnMarshal(v, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Object = &validatorpb.SignRequest_AttestationData{v}
			iNdEx = postIndex
		case 103:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AggregateAttestationAndProof", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &eth.AggregateAttestationAndProof{}
			if err := LegacyAggregationAttAndProofUnMarshal(v, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Object = &validatorpb.SignRequest_AggregateAttestationAndProof{v}
			iNdEx = postIndex
		case 104:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Exit", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
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
				return ErrInvalidLengthKeymanager
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthKeymanager
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			v := &eth.VoluntaryExit{}
			if err := LegacyExitUnMarshal(v, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			m.Object = &validatorpb.SignRequest_Exit{v}
			iNdEx = postIndex
		case 105:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Slot", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Object = &validatorpb.SignRequest_Slot{types.Slot(v)}
		case 106:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Epoch", wireType)
			}
			var v uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowKeymanager
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Object = &validatorpb.SignRequest_Epoch{types.Epoch(v)}
		default:
			iNdEx = preIndex
			skippy, err := skipKeymanager(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if skippy < 0 {
				return ErrInvalidLengthKeymanager
			}
			if (iNdEx + skippy) < 0 {
				return ErrInvalidLengthKeymanager
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

func LegacySignRequestMarshal(m *validatorpb.SignRequest) ([]byte, error) {
	size := signReq_size(m)
	dAtA := make([]byte, size)
	n, err := signReq_marshalToSizedBuffer(m, dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func signReq_marshalToSizedBuffer(m *validatorpb.SignRequest, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if m.Object != nil {
		{
			size := obj_size(m)
			i -= size
			d := dAtA[i:]
			var err error
			switch t := m.GetObject().(type) {
			case *validatorpb.SignRequest_Block:
				s := signReqBlock_size(t.Block)
				_, err = signReqBlock_marshalToSizedBuffer(t.Block, d[:s])
				break
			case *validatorpb.SignRequest_BlockV2:
				break
			case *validatorpb.SignRequest_AttestationData:
				s := signReqAttData_size(t.AttestationData)
				_, err = signReqAttData_marshalToSizedBuffer(t.AttestationData, d[:s])
				break
			case *validatorpb.SignRequest_Slot:
				s := signReqSlot_size(uint64(t.Slot))
				_, err = slot_marshalToSizedBuffer(uint64(t.Slot), d[:s])
				break
			case *validatorpb.SignRequest_Epoch:
				s := signReqEpoch_size(uint64(t.Epoch))
				_, err = epoch_marshalToSizedBuffer(uint64(t.Epoch), d[:s])
				break
			case *validatorpb.SignRequest_AggregateAttestationAndProof:
				s := signReqAttAggProof_size(t.AggregateAttestationAndProof)
				_, err = signReqAttAggProof_marshalToSizedBuffer(t.AggregateAttestationAndProof, d[:s])
				break
			}
			if err != nil {
				return 0, err
			}
		}
	}
	if len(m.SignatureDomain) > 0 {
		i -= len(m.SignatureDomain)
		copy(dAtA[i:], m.SignatureDomain)
		i = encodeVarintKeymanager(dAtA, i, uint64(len(m.SignatureDomain)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.SigningRoot) > 0 {
		i -= len(m.SigningRoot)
		copy(dAtA[i:], m.SigningRoot)
		i = encodeVarintKeymanager(dAtA, i, uint64(len(m.SigningRoot)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PublicKey) > 0 {
		i -= len(m.PublicKey)
		copy(dAtA[i:], m.PublicKey)
		i = encodeVarintKeymanager(dAtA, i, uint64(len(m.PublicKey)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func signReq_size(m *validatorpb.SignRequest) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PublicKey)
	if l > 0 {
		n += 1 + l + sovKeymanager(uint64(l))
	}
	l = len(m.SigningRoot)
	if l > 0 {
		n += 1 + l + sovKeymanager(uint64(l))
	}
	l = len(m.SignatureDomain)
	if l > 0 {
		n += 1 + l + sovKeymanager(uint64(l))
	}
	if m.Object != nil {
		n += obj_size(m)
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}

func obj_size(m *validatorpb.SignRequest) (n int) {
	switch t := m.GetObject().(type) {
	case *validatorpb.SignRequest_Block:
		return signReqBlock_size(t.Block)
	case *validatorpb.SignRequest_BlockV2:
		return 0
	case *validatorpb.SignRequest_AttestationData:
		return signReqAttData_size(t.AttestationData)
		break
	case *validatorpb.SignRequest_Slot:
		return signReqSlot_size(uint64(t.Slot))
	case *validatorpb.SignRequest_Epoch:
		return signReqEpoch_size(uint64(t.Epoch))
	case *validatorpb.SignRequest_AggregateAttestationAndProof:
		return signReqAttAggProof_size(t.AggregateAttestationAndProof)
	}
	return 0
}

func signReqBlock_size(blk *eth.BeaconBlock) (n int) {
	if blk == nil {
		return 0
	}
	var l int
	_ = l
	if blk != nil {
		l = beaconBlock_size(blk)
		n += 2 + l + sovKeymanager(uint64(l))
	}
	return n
}

func signReqBlock_marshalToSizedBuffer(block *eth.BeaconBlock, dAtA []byte) (int, error) {
	i := len(dAtA)
	if block != nil {
		{
			size, err := beaconBlock_marshalToSizedBuffer(block, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintKeymanager(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6
		i--
		dAtA[i] = 0xaa
	}
	return len(dAtA) - i, nil
}

func signReqAttData_size(att *eth.AttestationData) (n int) {
	if att == nil {
		return 0
	}
	var l int
	_ = l
	if att != nil {
		l = attestationData_size(att)
		n += 2 + l + sovKeymanager(uint64(l))
	}
	return n
}

func signReqAttData_marshalToSizedBuffer(att *eth.AttestationData, dAtA []byte) (int, error) {
	i := len(dAtA)
	if att != nil {
		{
			size, err := attestationData_marshalToSizedBuffer(att, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintKeymanager(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6
		i--
		dAtA[i] = 0xb2
	}
	return len(dAtA) - i, nil
}

func signReqAttAggProof_size(att *eth.AggregateAttestationAndProof) (n int) {
	if att == nil {
		return 0
	}
	var l int
	_ = l
	if att != nil {
		l = aggAttProof_size(att)
		n += 2 + l + sovKeymanager(uint64(l))
	}
	return n
}

func signReqAttAggProof_marshalToSizedBuffer(att *eth.AggregateAttestationAndProof, dAtA []byte) (int, error) {
	i := len(dAtA)
	if att != nil {
		{
			size, err := aggAttProof_marshalToSizedBuffer(att, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintKeymanager(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x6
		i--
		dAtA[i] = 0xba
	}
	return len(dAtA) - i, nil
}

func signReqSlot_size(slot uint64) (n int) {
	//if m == nil {
	//	return 0
	//}
	var l int
	_ = l
	n += 2 + sovKeymanager(uint64(slot))
	return n
}

func slot_marshalToSizedBuffer(slot uint64, dAtA []byte) (int, error) {
	i := len(dAtA)
	i = encodeVarintKeymanager(dAtA, i, uint64(slot))
	i--
	dAtA[i] = 0x6
	i--
	dAtA[i] = 0xc8
	return len(dAtA) - i, nil
}

func signReqEpoch_size(epoch uint64) (n int) {
	//if m == nil {
	//	return 0
	//}
	var l int
	_ = l
	n += 2 + sovKeymanager(uint64(epoch))
	return n
}

func epoch_marshalToSizedBuffer(epoch uint64, dAtA []byte) (int, error) {
	i := len(dAtA)
	i = encodeVarintKeymanager(dAtA, i, uint64(epoch))
	i--
	dAtA[i] = 0x6
	i--
	dAtA[i] = 0xd0
	return len(dAtA) - i, nil
}
