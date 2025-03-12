
> fork form github.com/unidoc/unitype

1、删除了测试用例

2、调整log库为slog

3、重写subset，原库未实现这个方法

**我想要实现能够精准截取指定字，并生成一个较小的ttf字体文件，用于嵌入式等场景使用，源库虽然能够截取，但只能实现从前到后的截取，
文件中依旧会存在大量无用字体数据，本库实现了subset方法，能够精准截取传入字符，并删除一切不相关的字库表数据，以求最精简字库文件生成**

### UniType - truetype font library for golang.
This library is designed for parsing and editing truetype fonts.
Useful along with UniPDF for subsetting fonts for use in PDF files.


```text
Subset font: trec: present with 18 table records
GPOS: 76.37 kB
GSUB: 4.64 kB
OS/2: 96 B
VDMX: 5.71 kB
cmap: 984 B
cvt: 532 B
fpgm: 1.53 kB
gasp: 16 B
glyf: 131.51 kB
hdmx: 34.57 kB
head: 54 B
hhea: 36 B
hmtx: 4.93 kB
loca: 4.93 kB
maxp: 32 B
name: 724 B
post: 11.11 kB
prep: 441 B
--

Subset font length: 1478
subset font is valid
==========info============
trec: present with 10 table records
head: 54 B
maxp: 32 B
hhea: 36 B
hmtx: 16 B
loca: 20 B
glyf: 732 B
name: 188 B
OS/2: 96 B
post: 32 B
cmap: 100 B
--
head table: &unitype.headTable{majorVersion:0x1, minorVersion:0x0, fontRevision:54460, checksumAdjustment:0xb7f49fd1, magicNumber:0x5f0f3cf5, flags:0x19, unitsPerEm:0x3e8, created:3381223050, modified:3576827943, xMin:-1129, yMin:-205, xMax:3480, yMax:988, macStyle:0x0, lowestRecPPEM:0x9, fontDirectionHint:2, indexToLocFormat:1, glyphDataFormat:0}
os/2 table: &unitype.os2Table{version:0x3, xAvgCharWidth:614, usWeightClass:0x1f4, usWidthClass:0x5, fsType:0x0, ySubscriptXSize:700, ySubscriptYSize:650, ySubscriptXOffset:0, ySubscriptYOffset:140, ySuperscriptXSize:700, ySuperscriptYSize:650, ySuperscriptXOffset:0, ySuperscriptYOffset:477, yStrikeoutSize:89, yStrikeoutPosition:250, sFamilyClass:0, panose10:[]uint8{0x2, 0xb, 0x6, 0x4, 0x3, 0x6, 0x2, 0x3, 0x2, 0x4}, ulUnicodeRange1:0xe00002ff, ulUnicodeRange2:0x5000205b, ulUnicodeRange3:0x0, ulUnicodeRange4:0x0, achVendID:unitype.tag{0x44, 0x41, 0x4d, 0x41}, fsSelection:0x40, usFirstCharIndex:0x0, usLastCharIndex:0xfb04, sTypoAscender:776, sTypoDescender:-185, sTypoLineGap:56, usWinAscent:0x3a4, usWinDescent:0xbd, ulCodePageRange1:0x2000009f, ulCodePageRange2:0x56010000, sxHeight:523, sCapHeight:693, usDefaultChar:0x0, usBreakChar:0x20, usMaxContext:0x3, usLowerOpticalPointSize:0x0, usUpperOpticalPointSize:0x0}
hhea table: numHMetrics: 4
hmtx: hmetrics: 4, leftSideBearings: 0
cmap version: 0
cmap: encoding records: 2 subtables: 2
cmap: subtables: [4,0,3 4,3,1]
cmap subtable: 4,0,3: runes: 4
	0 - Charcode 0 (0x0) - rune  0
	1 - Charcode 65 (0x41) - rune  41
	2 - Charcode 66 (0x42) - rune  42
	3 - Charcode 67 (0x43) - rune  43
cmap subtable: 4,3,1: runes: 4
	0 - Charcode 0 (0x0) - rune  0
	1 - Charcode 65 (0x41) - rune  41
	2 - Charcode 66 (0x42) - rune  42
	3 - Charcode 67 (0x43) - rune  43
Loca table
- Short offsets: 0
- Long offsets: 5
glyf table present: 4 descriptions (0.71 kB)
post table present: 0 numGlyphs
- post glyphNameIndex: 0
- post glyphNames: 0
&unitype.postTable{version:196608, italicAngle:0, underlinePosition:-123, underlineThickness:20, isFixedPitch:0x0, minMemType42:0x0, maxMemType42:0x0, minMemType1:0x0, maxMemType1:0x0, numGlyphs:0x0, glyphNameIndex:[]uint16(nil), offsets:[]int8(nil), glyphNames:[]unitype.GlyphName(nil)}
name table
&unitype.nameTable{format:0x0, count:0xe, stringOffset:0xae, nameRecords:[]*unitype.nameRecord{(*unitype.nameRecord)(0xc0003a6390), (*unitype.nameRecord)(0xc0003a63c0), (*unitype.nameRecord)(0xc0003a63f0), (*unitype.nameRecord)(0xc0003a6420), (*unitype.nameRecord)(0xc0003a6450), (*unitype.nameRecord)(0xc0003a6480), (*unitype.nameRecord)(0xc0003a64b0), (*unitype.nameRecord)(0xc0003a64e0), (*unitype.nameRecord)(0xc0003a6510), (*unitype.nameRecord)(0xc0003a6540), (*unitype.nameRecord)(0xc0003a6570), (*unitype.nameRecord)(0xc0003a65a0), (*unitype.nameRecord)(0xc0003a65d0), (*unitype.nameRecord)(0xc0003a6600)}, langTagCount:0x0, langTagRecords:[]*unitype.langTagRecord(nil)}
0/0,3: 1 - A
0/0,3: 2 - B
0/0,3: 3 - C
2/3,1: 1 - A
2/3,1: 2 - B
2/3,1: 3 - C

Process finished with the exit code 0

```
