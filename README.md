
fork form github.com/unidoc/unitype

1、删除了测试用例
2、调整log库为slog
3、重写subset，原库未实现这个方法

我想要实现能够精准截取指定字，并生成一个较小的ttf字体文件，用于嵌入式等场景使用，源库虽然能够截取，但只能实现从前到后的截取，文件中依旧会存在大量无用字体数据
本库实现了subset方法，能够精准截取传入字符，并删除一切不相关的字库表数据，以求最精简字库文件生成


# UniType - truetype font library for golang.
This library is designed for parsing and editing truetype fonts.
Useful along with UniPDF for subsetting fonts for use in PDF files.

Contains a CLI for useful operations:
```bash
$ ./truecli
truecli - TrueCLI

Usage:
  truecli [command]

Available Commands:
  help        Help about any command
  info        Get font file info
  readwrite   Read and write font file
  subset      Subset font
  validate    Validate font file

Flags:
  -h, --help              help for truecli
      --loglevel string   Log level 'debug' and 'trace' give debug information

Use "truecli [command] --help" for more information about a command.
```

for example:
```bash
$ $ ./truecli info --trec myfnt.ttf
  trec: present with 22 table records
  DSIG: 6.78 kB
  GSUB: 368 B
  LTSH: 28.53 kB
  OS/2: 96 B
  VDMX: 1.47 kB
  cmap: 124.93 kB
  cvt: 672 B
  fpgm: 2.33 kB
  gasp: 16 B
  glyf: 13.12 MB
  hdmx: 484.97 kB
  head: 54 B
  hhea: 36 B
  hmtx: 114.09 kB
  kern: 3.43 kB
  loca: 114.09 kB
  maxp: 32 B
  name: 2.54 kB
  post: 283.82 kB
  prep: 676 B
  vhea: 36 B
  vmtx: 114.09 kB 
```

