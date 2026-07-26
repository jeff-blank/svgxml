package svgxml

import (
	"encoding/xml"
)

// SVG is the base struct of an SVG image. The array fields of
// the struct are the supported tags and will be rendered in
// first-to-last order (e. g., any circle defined in SVG.Circle
// will visually cover any other overlapping elements).
type SVG struct {
	Id      string      `xml:"id,attr"`
	Width   string      `xml:"width,attr"`
	Height  string      `xml:"height,attr"`
	ViewBox string      `xml:"viewBox,attr"`
	Title   string      `xml:"title,attr,omitempty"`
	Version string      `xml:"version,attr"`
	XMLNS   string      `xml:"xmlns,attr"`
	XMLName xml.Name    `xml:"svg"`
	Defs    DefsDef     `xml:"defs"`
	G       []GroupDef  `xml:"g"`
	A       []AnchorDef `xml:"a"`
	Path    []PathDef   `xml:"path"`
	Text    []TextDef   `xml:"text"`
	Rect    []RectDef   `xml:"rect"`
	Circle  []CircleDef `xml:"circle"`
}

type PathDef struct {
	Id    string `xml:"id,attr"`
	D     string `xml:"d,attr"`
	Xform string `xml:"transform,attr,omitempty"`
	Style string `xml:"style,attr,omitempty"`
	Title string `xml:"title,omitempty"`
}

type GroupDef struct {
	G      []GroupDef  `xml:"g,omitempty"`
	Id     string      `xml:"id,attr,omitempty"`
	Xform  string      `xml:"transform,attr,omitempty"`
	Style  string      `xml:"style,attr,omitempty"`
	A      []AnchorDef `xml:"a"`
	Path   []PathDef   `xml:"path"`
	Text   []TextDef   `xml:"text"`
	Rect   []RectDef   `xml:"rect"`
	Circle []CircleDef `xml:"circle"`
}

type AnchorDef struct {
	Id             string      `xml:"id,attr"`
	Download       string      `xml:"Download,attr"`
	Href           string      `xml:"href,attr"`
	HrefLang       string      `xml:"hreflang,attr"`
	ReferrerPolicy string      `xml:"referrerpolicy,attr"`
	Rel            string      `xml:"rel,attr"`
	Target         string      `xml:"target,attr"`
	Type           string      `xml:"type,attr"`
	Xform          string      `xml:"transform,attr,omitempty"`
	G              []GroupDef  `xml:"g"`
	Rect           []RectDef   `xml:"rect"`
	Circle         []CircleDef `xml:"circle"`
}

type DefsDef struct {
	Id string `xml:"id,attr"`
}

type TSpanDef struct {
	Id    string `xml:"id,attr"`
	X     string `xml:"x,attr"`
	Y     string `xml:"y,attr"`
	Xform string `xml:"transform,attr,omitempty"`
	Label string `xml:",chardata"`
}

type TextDef struct {
	Id         string     `xml:"id,attr"`
	X          string     `xml:"x,attr"`
	Y          string     `xml:"y,attr"`
	Style      string     `xml:"style,attr,omitempty"`
	Xform      string     `xml:"transform,attr,omitempty"`
	TextLength string     `xml:",attr,omitempty"`
	Label      string     `xml:",chardata"`
	TSpan      []TSpanDef `xml:"tspan"`
}

type RectDef struct {
	Id     string `xml:"id,attr"`
	X      string `xml:"x,attr"`
	Y      string `xml:"y,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
	Style  string `xml:"style,attr,omitempty"`
	Xform  string `xml:"transform,attr,omitempty"`
}

type CircleDef struct {
	Id     string `xml:"id,attr"`
	CX     string `xml:"cx,attr"`
	CY     string `xml:"cy,attr"`
	Radius string `xml:"r,attr"`
	Style  string `xml:"style,attr,omitempty"`
	Xform  string `xml:"transform,attr,omitempty"`
}

const (
	FindFirst = true
	FindAll   = false
)

type NumberAttr interface {
	int | int32 | int64 | float32 | float64
}

type ViewBoxDef struct {
	X, Y, Width, Height float64
}
