package legacy

import (
	"fmt"
	"io"

	eth "github.com/prysmaticlabs/prysm/proto/prysm/v1alpha1"
)

func LegacyBeaconBlockBodyUnMarshal(m *eth.BeaconBlockBody, dAtA []byte) error {
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
			return fmt.Errorf("proto: BeaconBlockBody: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: BeaconBlockBody: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RandaoReveal", wireType)
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
			m.RandaoReveal = append(m.RandaoReveal[:0], dAtA[iNdEx:postIndex]...)
			if m.RandaoReveal == nil {
				m.RandaoReveal = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Eth1Data", wireType)
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
			if m.Eth1Data == nil {
				m.Eth1Data = &eth.Eth1Data{}
			}
			if err := LegacyETH1UnMarshal(m.Eth1Data, dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Graffiti", wireType)
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
			m.Graffiti = append(m.Graffiti[:0], dAtA[iNdEx:postIndex]...)
			if m.Graffiti == nil {
				m.Graffiti = []byte{}
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ProposerSlashings", wireType)
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
			m.ProposerSlashings = append(m.ProposerSlashings, &eth.ProposerSlashing{})
			if err := LegacyProposerSlashingUnMarshal(m.ProposerSlashings[len(m.ProposerSlashings)-1], dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AttesterSlashings", wireType)
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
			m.AttesterSlashings = append(m.AttesterSlashings, &eth.AttesterSlashing{})
			if err := LegacyAttesterSlashingUnMarshal(m.AttesterSlashings[len(m.AttesterSlashings)-1], dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Attestations", wireType)
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
			m.Attestations = append(m.Attestations, &eth.Attestation{})
			if err := LegacyAttestationUnMarshal(m.Attestations[len(m.Attestations)-1], dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Deposits", wireType)
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
			m.Deposits = append(m.Deposits, &eth.Deposit{})
			if err := LegacyDepositUnMarshal(m.Deposits[len(m.Deposits)-1], dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoluntaryExits", wireType)
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
			m.VoluntaryExits = append(m.VoluntaryExits, &eth.SignedVoluntaryExit{})
			if err := LegacySignedExitUnMarshal(m.VoluntaryExits[len(m.VoluntaryExits)-1], dAtA[iNdEx:postIndex]); err != nil {
				return err
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

func beaconBlockBody_marshalToSizedBuffer(m *eth.BeaconBlockBody, dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	//if m.XXX_unrecognized != nil {
	//	i -= len(m.XXX_unrecognized)
	//	copy(dAtA[i:], m.XXX_unrecognized)
	//}
	if len(m.VoluntaryExits) > 0 {
		for iNdEx := len(m.VoluntaryExits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := signedVoluntary_marshalToSizedBuffer(m.VoluntaryExits[iNdEx], dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if len(m.Deposits) > 0 {
		for iNdEx := len(m.Deposits) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := deposit_marshalToSizedBuffer(m.Deposits[iNdEx], dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if len(m.Attestations) > 0 {
		for iNdEx := len(m.Attestations) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := attestation_marshalToSizedBuffer(m.Attestations[iNdEx], dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x32
		}
	}
	if len(m.AttesterSlashings) > 0 {
		for iNdEx := len(m.AttesterSlashings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := attesterSlashing_marshalToSizedBuffer(m.AttesterSlashings[iNdEx], dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x2a
		}
	}
	if len(m.ProposerSlashings) > 0 {
		for iNdEx := len(m.ProposerSlashings) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := proposerSlashing_marshalToSizedBuffer(m.ProposerSlashings[iNdEx], dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Graffiti) > 0 {
		i -= len(m.Graffiti)
		copy(dAtA[i:], m.Graffiti)
		i = encodeVarintBeaconBlock(dAtA, i, uint64(len(m.Graffiti)))
		i--
		dAtA[i] = 0x1a
	}
	if m.Eth1Data != nil {
		{
			size, err := eth1data_marshalToSizedBuffer(m.Eth1Data, dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintBeaconBlock(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if len(m.RandaoReveal) > 0 {
		i -= len(m.RandaoReveal)
		copy(dAtA[i:], m.RandaoReveal)
		i = encodeVarintBeaconBlock(dAtA, i, uint64(len(m.RandaoReveal)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func beaconBlockBody_size(m *eth.BeaconBlockBody) (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.RandaoReveal)
	if l > 0 {
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	if m.Eth1Data != nil {
		l = eth1data_size(m.Eth1Data)
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	l = len(m.Graffiti)
	if l > 0 {
		n += 1 + l + sovBeaconBlock(uint64(l))
	}
	if len(m.ProposerSlashings) > 0 {
		for _, e := range m.ProposerSlashings {
			l = proposerSlashing_size(e)
			n += 1 + l + sovBeaconBlock(uint64(l))
		}
	}
	if len(m.AttesterSlashings) > 0 {
		for _, e := range m.AttesterSlashings {
			l = attesterSlashing_size(e)
			n += 1 + l + sovBeaconBlock(uint64(l))
		}
	}
	if len(m.Attestations) > 0 {
		for _, e := range m.Attestations {
			l = attestation_size(e)
			n += 1 + l + sovBeaconBlock(uint64(l))
		}
	}
	if len(m.Deposits) > 0 {
		for _, e := range m.Deposits {
			l = deposit_size(e)
			n += 1 + l + sovBeaconBlock(uint64(l))
		}
	}
	if len(m.VoluntaryExits) > 0 {
		for _, e := range m.VoluntaryExits {
			l = signedVoluntary_size(e)
			n += 1 + l + sovBeaconBlock(uint64(l))
		}
	}
	//if m.XXX_unrecognized != nil {
	//	n += len(m.XXX_unrecognized)
	//}
	return n
}
