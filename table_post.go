/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package unitype

import (
	"errors"
	"fmt"
	"log/slog"
)

// postTable represents a PostScript (post) table.
// This table contains additional information needed for use on PostScript printers.
// Includes FontInfo dictionary entries and the PostScript names of all glyphs.
//
//   - version 1.0 is used the font file contains exactly the 258 glyphs in the standard Macintosh TrueType font file.
//     Glyph list on: https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6post.html
//   - version 2.0 is used for fonts that contain some glyphs not in the standard set or have different ordering.
//   - version 2.5 can handle nonstandard ordering of the standard mac glyphs via offsets.
//   - other versions do not contain post glyph name data.
type postTable struct {
	// header (all versions).
	version            fixed
	italicAngle        fixed // in degrees.
	underlinePosition  fword
	underlineThickness fword
	isFixedPitch       uint32
	minMemType42       uint32
	maxMemType42       uint32
	minMemType1        uint32
	maxMemType1        uint32

	// version 2.0 and 2.5 (partly).
	numGlyphs      uint16   // should equal maxp.numGlyphs
	glyphNameIndex []uint16 // len = numGlyphs

	// version 2.5.
	offsets []int8 // len = numGlyphs

	// Processed data:
	// TODO: Check `len = glyphNames` below, should be numGlyphs ?
	glyphNames []GlyphName // len = glyphNames, index is GlyphID (GID), glyphNames[GlyphID] -> GlyphName.
}

/*
 See https://developer.apple.com/fonts/TrueType-Reference-Manual/RM06/Chap6post.html
 and https://docs.microsoft.com/en-us/typography/opentype/spec/post
 for details regarding the format.
*/

func (f *font) parsePost(r *byteReader) (*postTable, error) {
	slog.Debug("Parsing post table")
	if f.maxp == nil {
		// maxp table required for numGlyphs check. Could probably be omitted, can consider
		// if run into those cases where post is present and maxp is not (and all other information present).
		slog.Debug("Required maxp table missing")
		return nil, errRequiredField
	}

	tr, has, err := f.seekToTable(r, "post")
	if err != nil {
		return nil, err
	}
	if !has {
		slog.Debug("Post table not present")
		return nil, nil
	}

	start := r.Offset()

	t := &postTable{}
	err = r.read(&t.version, &t.italicAngle, &t.underlinePosition, &t.underlineThickness, &t.isFixedPitch)
	if err != nil {
		return nil, err
	}
	err = r.read(&t.minMemType42, &t.maxMemType42, &t.minMemType1, &t.maxMemType1)
	if err != nil {
		return nil, err
	}

	slog.Debug(fmt.Sprintf("Version: %v %v 0x%X", t.version, t.version.Float64(), t.version))
	switch uint32(t.version) {
	case 0x00010000: // 1.0 - font files contains exactly the 258 standard Macintosh glyphs.
		if t.numGlyphs != 258 {
			slog.Debug("Should have the mac number of glyph names")
			// TODO(gunnsth): If this is too strict, can just set the first 258 glyphnames.
			return nil, errRangeCheck
		}
		t.glyphNames = make([]GlyphName, int(t.numGlyphs))
		for i := range macGlyphNames {
			t.glyphNames[i] = macGlyphNames[i]
		}

	case 0x00020000: // 2.0
		slog.Debug("Version: 2.0")
		err = r.read(&t.numGlyphs)
		if err != nil {
			return nil, err
		}
		slog.Debug(fmt.Sprintf("numGlyphs: %d", t.numGlyphs))
		if t.numGlyphs != f.maxp.numGlyphs {
			slog.Debug(fmt.Sprintf("post numGlyphs != maxp.numGlyphs (%d != %d)", t.numGlyphs, f.maxp.numGlyphs))
			return nil, errRangeCheck
		}
		err = r.readSlice(&t.glyphNameIndex, int(t.numGlyphs))
		if err != nil {
			return nil, err
		}
		newGlyphs := 0
		for _, ni := range t.glyphNameIndex {
			if ni >= 258 && ni <= 32767 {
				newGlyphs++
			}
		}
		slog.Debug(fmt.Sprintf("newGlyphs: %d", newGlyphs))
		var names []string
		for i := 0; i < newGlyphs; i++ {
			if r.Offset()-start >= int64(tr.length) {
				slog.Debug("ERROR: Reading outside post table")
				slog.Debug(fmt.Sprintf("%d > %d", r.Offset()-start, tr.length))
				return nil, errors.New("reading outside table")
			}
			var numChars int8
			err = r.read(&numChars)
			if err != nil {
				return nil, err
			}
			if numChars == 0 {
				break
			}

			name := make([]byte, numChars)
			err = r.readBytes(&name, int(numChars))
			if err != nil {
				slog.Debug(fmt.Sprintf("ERROR: %v", err))
				return nil, err
			}

			names = append(names, string(name))
		}
		if len(names) != newGlyphs {
			slog.Debug(fmt.Sprintf("newGlyphs != len(names) (%d != %d)", len(names), newGlyphs))
			return nil, errors.New("mismatching number of names loaded")
		}

		t.glyphNames = make([]GlyphName, int(t.numGlyphs))
		for i := 0; i < int(t.numGlyphs); i++ {
			var name GlyphName

			ni := t.glyphNameIndex[i]
			if ni < 258 {
				name = macGlyphNames[ni]
			} else if ni <= 32767 {
				ni -= 258
				if int(ni) >= len(names) {
					slog.Debug(fmt.Sprintf("ERROR: Glyph %d referring to outside name list (%d)", i, ni))
					// Let's be strict initially and slack if we find that it is needed.
					return nil, errRangeCheck
				}
				name = GlyphName(names[ni])
			}
			slog.Debug(fmt.Sprintf("GID %d -> '%s'", i, name))
			t.glyphNames[i] = name
		}
		slog.Debug(fmt.Sprintf("len(names) = %d", len(names)))

	case 0x00025000: // 2.5
		slog.Debug("Version: 2.5")
		err = r.read(&t.numGlyphs)
		if err != nil {
			return nil, err
		}
		if t.numGlyphs != f.maxp.numGlyphs {
			slog.Debug(fmt.Sprintf("post numGlyphs != maxp.numGlyphs (%d != %d)", t.numGlyphs, f.maxp.numGlyphs))
			return nil, errRangeCheck
		}
		err = r.readSlice(&t.offsets, int(t.numGlyphs))
		if err != nil {
			return nil, err
		}
		t.glyphNames = make([]GlyphName, int(t.numGlyphs))
		for i := 0; i < int(t.numGlyphs); i++ {
			nameIndex := i + 1 + int(t.offsets[i])
			if nameIndex < 0 || nameIndex > 257 {
				slog.Debug(fmt.Sprintf("ERROR: name index outside range (%d)", nameIndex))
				continue
			}
			t.glyphNames[i] = macGlyphNames[nameIndex]
			slog.Debug(fmt.Sprintf("2.5 I: %d -> %s", i, t.glyphNames[i]))
		}

	case 0x00030000: // 3.0
		slog.Debug("Version 3.0 - no postscript data")
	default:
		slog.Debug(fmt.Sprintf("Unsupported version of post (%d) - no post data loaded", t.version))
	}

	return t, nil
}

func (f *font) writePost(w *byteWriter) error {
	if f.post == nil {
		return nil
	}
	t := f.post

	// TODO(gunnsth): Write out with v1.0 or v2.0.
	version := t.version
	if version != 0x00010000 {
		// Include no postscript data.
		// TODO(gunnsth): support writing v2.0.
		version = 0x00030000
		t.version = version
	}

	err := w.write(t.version, t.italicAngle, t.underlinePosition, t.underlineThickness, t.isFixedPitch)
	if err != nil {
		return err
	}
	err = w.write(t.minMemType42, t.maxMemType42, t.minMemType1, t.maxMemType1)
	if err != nil {
		return err
	}

	return nil
}
