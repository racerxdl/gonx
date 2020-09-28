//+build ignore

package font

const _FONTNAME_fontName = "_FONTNAME_"

var _FONTNAME_Data = Data{
	CharWidth:  _FONTWIDTH_,
	CharHeight: _FONTHEIGHT_,
	XOffset:    _FONTXOFF_,
	YOffset:    _FONTYOFF_,
	Glyphs:     map[uint32]*Glyph{_GLYPHS_},
}

func init() {
	fonts["_FONTNAME_"] = &_FONTNAME_Data
}
