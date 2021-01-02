// +build ignore

// This program generates fontname.go. It can be invoked by running
// go generate
package main

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	charRegex     = regexp.MustCompile(`STARTCHAR([\S\s]*?)ENDCHAR`)
	boundBoxRegex = regexp.MustCompile(`FONTBOUNDINGBOX (.*)`)
)

func hexByteToBytes(hexByte string) ([]byte, error) {
	numBits := len(hexByte) * 4 // Hex string chars = nibbles
	v, err := strconv.ParseUint(hexByte, 16, numBits)
	if err != nil {
		fmt.Printf("Error parsing %q as number: %s\n", hexByte, err)
		return nil, err
	}
	bits := make([]byte, 0)

	for i := numBits - 1; i >= 0; i-- {
		mask := 1 << uint(i)
		if v&uint64(mask) > 0 {
			bits = append(bits, 1)
		} else {
			bits = append(bits, 0)
		}
	}

	return bits, nil
}

func buildGlyph(glyphInfo string) (uint32, string) {
	glyphFragment := ""
	lines := strings.Split(glyphInfo, "\n")

	var data [][]byte
	name := ""
	encoding := uint32(0)

	for _, v := range lines {
		switch {
		case strings.Contains(v, "STARTCHAR"):
			name = strings.Trim(v[len("STARTCHAR"):], " \r\n")
			//fmt.Println("GLYPH NAME", name)
		case strings.Contains(v, "ENCODING"):
			enc := strings.Trim(v[len("ENCODING"):], " \r\n")
			enc64, err := strconv.ParseUint(enc, 10, 32)
			if err != nil {
				panic(err)
			}
			encoding = uint32(enc64)
			//fmt.Println("ENCODING", encoding)
		case strings.Contains(v, "BITMAP"):
		case strings.Contains(v, "SWIDTH"):
		case strings.Contains(v, "DWIDTH"):
		case strings.Contains(v, "BBX"):
		case strings.Contains(v, "BITMAP"):
		case strings.Contains(v, "ENDCHAR"):
		default:
			rowData, err := hexByteToBytes(v)
			if err != nil {
				continue
			}
			data = append(data, rowData)
		}
	}

	//fmt.Println(name, encoding, data)

	contentFragment := ""

	for _, v := range data {
		frag := "{"
		for z, b := range v {
			frag += fmt.Sprintf(" %d", b)
			if z != len(v)-1 {
				frag += ","
			}
		}
		frag += " },\n"
		contentFragment += frag
	}

	glyphFragment = fmt.Sprintf(`&Glyph{
			Name: "%s",
			Data: [][]byte{
				%s
			},
		}`, name, contentFragment)

	z, err := format.Source([]byte(glyphFragment))

	if err != nil {
		panic(err)
	}

	glyphFragment = string(z)

	return encoding, glyphFragment
}

func buildFont(tpl, name, bdf string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.Replace(name, "-", "_", -1)
	name = strings.Replace(name, ".", "_", -1)
	tpl = strings.Replace(tpl, "_FONTNAME_", name, -1)

	boundBox := strings.Split(boundBoxRegex.FindString(bdf), " ")[1:]
	tpl = strings.Replace(tpl, "_FONTWIDTH_", boundBox[0], -1)
	tpl = strings.Replace(tpl, "_FONTHEIGHT_", boundBox[1], -1)
	tpl = strings.Replace(tpl, "_FONTXOFF_", boundBox[2], -1)
	tpl = strings.Replace(tpl, "_FONTYOFF_", boundBox[3], -1)

	chars := charRegex.FindAllStringSubmatch(bdf, -1)

	glyphList := ""
	for _, v := range chars {
		encoding, frag := buildGlyph(v[0])
		glyphList += fmt.Sprintf("%d: %s,\n", encoding, frag)
	}

	tpl = strings.Replace(tpl, "_GLYPHS_", "\n"+glyphList, -1)

	z, err := format.Source([]byte(tpl))

	if err != nil {
		fmt.Printf("Error formatting generated code: %s\n%s\n", err, string(tpl))
		panic(err)
	}

	return string(z)
}

func main() {
	var files []string

	root := "."
	err := filepath.Walk(root, func(f string, info os.FileInfo, err error) error {
		if strings.Contains(f, ".bdf") {
			files = append(files, f)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}

	tplB, err := ioutil.ReadFile("font.templ.go")
	if err != nil {
		panic(err)
	}

	tpl := string(tplB)

	tpl = strings.Replace(tpl, "//+build ignore\n\n", "", -1) // Remove build ignore

	for _, file := range files {
		fmt.Printf("Building font %s\n", file)
		fontData, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		fontName := strings.Replace(strings.ToLower(path.Base(file)), ".bdf", "", -1)
		bdf := string(fontData)
		resultFile := buildFont(tpl, fontName, bdf)
		filename := "font_gen_" + fontName + ".go"
		err = ioutil.WriteFile(filename, []byte(resultFile), 0660)
		if err != nil {
			panic(err)
		}
	}
}
