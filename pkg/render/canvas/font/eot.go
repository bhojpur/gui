package font

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

import (
	"fmt"
)

// ParseEOT parses the EOT font format and returns its contained SFNT font format (TTF or OTF). See https://www.w3.org/Submission/EOT/
func ParseEOT(b []byte) ([]byte, error) {
	r := NewBinaryReader(b)
	_ = r.ReadUint32LE()             // EOTSize
	fontDataSize := r.ReadUint32LE() // FontDataSize
	version := r.ReadUint32LE()      // Version
	if version != 0x00010000 && version != 0x00020001 && version != 0x00020002 {
		return nil, fmt.Errorf("unsupported version")
	}
	flags := r.ReadUint32LE()       // Flags
	_ = r.ReadBytes(10)             // FontPANOSE
	_ = r.ReadByte()                // Charset
	_ = r.ReadByte()                // Italic
	_ = r.ReadUint32LE()            // Weight
	_ = r.ReadUint16LE()            // fsType
	magicNumber := r.ReadUint16LE() // MagicNumber
	if magicNumber != 0x504C {
		return nil, fmt.Errorf("invalid magic number")
	}
	_ = r.ReadBytes(24) // Unicode and CodePage ranges
	checkSumAdjustment := r.ReadUint32LE()
	_ = r.ReadBytes(16)  // Reserved
	_ = r.ReadUint16LE() // Padding1

	familyNameSize := r.ReadUint16LE()      // FamilyNameSize
	_ = r.ReadBytes(uint32(familyNameSize)) // FamilyName
	_ = r.ReadUint16LE()                    // Padding2

	styleNameSize := r.ReadUint16LE()      // StyleNameSize
	_ = r.ReadBytes(uint32(styleNameSize)) // Stylename
	_ = r.ReadUint16LE()                   // Padding3

	versionNameSize := r.ReadUint16LE()      // VersionNameSize
	_ = r.ReadBytes(uint32(versionNameSize)) // VersionName
	_ = r.ReadUint16LE()                     // Padding4

	fullNameSize := r.ReadUint16LE()      // FullNameSize
	_ = r.ReadBytes(uint32(fullNameSize)) // FullName

	if version == 0x00020001 || version == 0x00020002 {
		_ = r.ReadUint16LE()                    // Padding5
		rootStringSize := r.ReadUint16LE()      // RootStringSize
		_ = r.ReadBytes(uint32(rootStringSize)) // RootString
	}
	if version == 0x00020002 {
		_ = r.ReadUint32LE()                   // RootStringCheckSum
		_ = r.ReadUint32LE()                   // EUDCCodePage
		_ = r.ReadUint16LE()                   // Padding6
		signatureSize := r.ReadUint16LE()      // SignatureSize
		_ = r.ReadBytes(uint32(signatureSize)) // Signature
		_ = r.ReadUint32LE()                   // EUDCFlags
		eudcFontSize := r.ReadUint32LE()       // EUDCFontSize
		_ = r.ReadBytes(uint32(eudcFontSize))  // EUDCFontData
	}

	fontData := r.ReadBytes(fontDataSize)
	if r.EOF() {
		return nil, ErrInvalidFontData
	}

	isCompressed := (flags & 0x00000004) != 0
	isXORed := (flags & 0x10000000) != 0

	if isXORed {
		for i := 0; i < len(fontData); i++ {
			fontData[i] ^= 0x50
		}
	}

	if isCompressed {
		// TODO: (EOT) see https://www.w3.org/Submission/MTX/
		return nil, fmt.Errorf("EOT compression not supported")
	}

	_ = checkSumAdjustment
	// TODO: (EOT) verify or recalculate master checksum
	//checksum := 0xB1B0AFBA - calcChecksum(w.Bytes())
	//if checkSumAdjustment != checksum {
	//return nil, 0, fmt.Errorf("bad checksum")
	//}

	// replace overal checksum in head table
	//buf := w.Bytes()
	//binary.BigEndian.PutUint32(buf[iCheckSumAdjustment:], checkSumAdjustment)
	return fontData, nil
}
