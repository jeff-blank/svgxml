package svgxml

import (
	"encoding/xml"
)

// SVG is the base struct of an SVG image, corresponding to the <svg>...</svg>
// tag pair and its attributes and contents. The array fields of the struct are
// the supported tags and will be rendered in first-to-last order (e. g., any
// circle defined in [SVG.Circle] will visually cover any other elements within
// its bounds).
type SVG struct {
	Id         string      `xml:"id,attr"`
	Width      string      `xml:"width,attr"`
	Height     string      `xml:"height,attr"`
	ViewBox    string      `xml:"viewBox,attr,omitempty"`
	Title      string      `xml:"title,attr,omitempty"`
	Version    string      `xml:"version,attr"`
	XMLNS      string      `xml:"xmlns,attr"`
	XMLName    xml.Name    `xml:"svg"`
	Defs       DefsDef     `xml:"defs"`
	G          []GroupDef  `xml:"g"`
	A          []AnchorDef `xml:"a"`
	Path       []PathDef   `xml:"path"`
	Text       []TextDef   `xml:"text"`
	Rect       []RectDef   `xml:"rect"`
	Circle     []CircleDef `xml:"circle"`
	background RectDef
}

// <path> SVG/XML element
type PathDef struct {
	Id    string `xml:"id,attr"`
	D     string `xml:"d,attr"`
	Xform string `xml:"transform,attr,omitempty"`
	Style string `xml:"style,attr,omitempty"`
	Title string `xml:"title,omitempty"`
}

// <g> SVG/XML element
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

// <a> SVG/XML element
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

// <defs> SVG/XML element
type DefsDef struct {
	Id string `xml:"id,attr"`
}

// <tspan> SVG/XML element
type TSpanDef struct {
	Id    string `xml:"id,attr"`
	X     string `xml:"x,attr"`
	Y     string `xml:"y,attr"`
	Xform string `xml:"transform,attr,omitempty"`
	Label string `xml:",chardata"`
}

// <text> SVG/XML element
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

// <rect> SVG/XML element
type RectDef struct {
	Id     string `xml:"id,attr"`
	X      string `xml:"x,attr"`
	Y      string `xml:"y,attr"`
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
	Style  string `xml:"style,attr,omitempty"`
	Xform  string `xml:"transform,attr,omitempty"`
}

// <circle> SVG/XML element
type CircleDef struct {
	Id     string `xml:"id,attr"`
	CX     string `xml:"cx,attr"`
	CY     string `xml:"cy,attr"`
	Radius string `xml:"r,attr"`
	Style  string `xml:"style,attr,omitempty"`
	Xform  string `xml:"transform,attr,omitempty"`
}

// "viewBox" attribute of some elements
type ViewBoxDef struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Type constraint for functions
type NumberAttr interface {
	int | int32 | int64 | float32 | float64 | string
}

const (
	FindFirst = true
	FindAll   = false
)
