package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"sort"

	"github.com/unidoc/unitype"
)

var s string = "ABC" // 咪妈咪哄567890 '《》？！，。~!@#$%^&*()_+-=\\|[]{};:\",.<>/?`、|￥×（）’“”‘：；【】\t²³μ㎡㎞㎝㎜㎏㎎℃"

func main() {
	tfnt, err := unitype.ParseFile("./Ubuntu-Medium.ttf")
	if err != nil {
		log.Fatalf("Error: %+v", err)
	}
	subfnt, err := tfnt.Subset([]rune(s))
	if err != nil {
		panic(err)
	}
	fmt.Printf("Subset font: %s\n", subfnt.String())
	check(subfnt)

	os.Remove("subset.ttf")
	err = subfnt.WriteFile("subset.ttf")
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
	}
	info()
}
func check(subfnt *unitype.Font) {
	buf := new(bytes.Buffer)
	err := subfnt.Write(buf)
	if err != nil {
		fmt.Printf("Failed writing: %+v\n", err)
		panic(err)
	}
	fmt.Printf("Subset font length: %d\n", buf.Len())
	err = unitype.ValidateBytes(buf.Bytes())
	if err != nil {
		fmt.Printf("Invalid subfnt: %+v\n", err)
		panic(err)
	} else {
		fmt.Printf("subset font is valid\n")
	}
}

func info() {
	fnt, err := unitype.ParseFile("./subset.ttf")
	if err != nil {
		panic(err)
	}
	fmt.Println("==========info============")
	fmt.Print(fnt.TableInfo("trec"))
	fmt.Print(fnt.TableInfo("head"))
	fmt.Print(fnt.TableInfo("os2"))
	fmt.Print(fnt.TableInfo("hhea"))
	fmt.Print(fnt.TableInfo("hmtx"))
	fmt.Print(fnt.TableInfo("cmap"))
	fmt.Print(fnt.TableInfo("loca"))
	fmt.Print(fnt.TableInfo("glyf"))
	fmt.Print(fnt.TableInfo("post"))
	fmt.Print(fnt.TableInfo("name"))
	var maps []map[rune]unitype.GlyphIndex
	var mapNames []string
	maps = append(maps, fnt.GetCmap(0, 3))
	mapNames = append(mapNames, "0,3")
	maps = append(maps, fnt.GetCmap(1, 0))
	mapNames = append(mapNames, "1,0")
	maps = append(maps, fnt.GetCmap(3, 1))
	mapNames = append(mapNames, "3,1")

	for i := range maps {
		var gids []unitype.GlyphIndex
		gidMap := map[unitype.GlyphIndex]rune{}
		for rune, gid := range maps[i] {
			gidMap[gid] = rune
			gids = append(gids, gid)
		}
		sort.Slice(gids, func(i, j int) bool {
			return gids[i] < gids[j]
		})
		for _, gid := range gids {
			fmt.Printf("%d/%s: %d - %c\n", i, mapNames[i], gid, gidMap[gid])
		}
	}

}
