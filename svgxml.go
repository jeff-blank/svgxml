package svgxml

import (
	"encoding/xml"
	"fmt"
	"os"
	"strconv"
	s "strings"
)

// NewFromBytes accepts a byte slice of SVG XML data and returns an SVG
// object pointer or an error.
func NewFromBytes(xmlBytes []byte) (*SVG, error) {
	var svgData SVG
	err := xml.Unmarshal(xmlBytes, &svgData)
	if err != nil {
		return nil, fmt.Errorf("NewFromBytes(): parse XML: %w", err)
	}
	return &svgData, nil
}

// NewFromFile accepts a filename containing SVG XML data and returns an
// object pointer or an error.
//
// This method calls [NewFromBytes] after successfully reading from the file
// named in the filename parameter.
func NewFromFile(filename string) (*SVG, error) {
	xmlBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("NewFromFile(): read SVG XML file: %w", err)
	}
	svgData, err := NewFromBytes(xmlBytes)
	if err != nil {
		err = fmt.Errorf("NewFromFile(): %w", err)
	}
	return svgData, err
}

// NewRect returns a [RectDef] object with the supplied parameters. x, y, width,
// and height can be all strings or a mix of int or float types (but not a mix
// of strings and ints/floats).
func NewRect[N NumberAttr](id, style string, x, y, width, height N) RectDef {
	newRect := RectDef{
		Id:     id,
		Style:  style,
		X:      anyNumberToString(x, 6),
		Y:      anyNumberToString(y, 6),
		Width:  anyNumberToString(width, 6),
		Height: anyNumberToString(height, 6),
	}
	return newRect
}

// GetXml returns the current object as an XML byte slice. If a background
// has been defined with [SVG.AddBackground], it is inserted into a copy of S
// as the bottom-most visual element, and the copy is marshaled to XML;
// otherwise S is marshaled as-is.
func (S *SVG) GetXml() ([]byte, error) {
	xmlBytes, err := xml.Marshal(placeBackgroundRect(S))
	if err != nil {
		return nil, fmt.Errorf("GetXml(): marshal: %w", err)
	}
	return append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), xmlBytes...), nil
}

// GetXmlIndented returns the current object as an XML byte slice. If a background
// has been defined with [SVG.AddBackground], it is inserted into a copy of S
// as the bottom-most visual element, and the copy is marshaled to XML;
// otherwise S is marshaled as-is.
//
// The indentPrefix and indentString parameters are passed directly to
// [xml.MarshalIndent].
func (S *SVG) GetXmlIndented(indentPrefix, indentString string) ([]byte, error) {
	xmlBytes, err := xml.MarshalIndent(placeBackgroundRect(S), indentPrefix, indentString)
	if err != nil {
		return nil, fmt.Errorf("GetXmlIndented(): marshal: %w", err)
	}
	return append([]byte(`<?xml version="1.0" encoding="UTF-8"?>`+"\n"), xmlBytes...), nil
}

// WriteFile writes an XML representation of the current object to the
// supplied filename.
//
// The method calls the [SVG.GetXml] method and writes the XML data to disk
// rather than returning a byte slice to the caller. The data is first written
// to a temporary file; if that succeeds, the file with the specified filename
// is unlinked and the temporary file renamed.
func (S *SVG) WriteFile(filename string) error {

	xmlBytes, err := S.GetXml()
	if err != nil {
		return fmt.Errorf("WriteFile(): generate XML: %w", err)
	}

	if err = doWriteFile(xmlBytes, filename); err != nil {
		return fmt.Errorf("WriteFile(): %w", err)
	}

	return nil
}

// WriteFileIndented writes an XML representation of the current object to the
// supplied filename.
//
// The method calls the [SVG.GetXmlIndented] method and writes the XML data to disk
// rather than returning a byte slice to the caller. The data is first written
// to a temporary file; if that succeeds, the file with the specified filename
// is unlinked and the temporary file renamed.
func (S *SVG) WriteFileIndented(filename, indentPrefix, indentString string) error {

	xmlBytes, err := S.GetXmlIndented(indentPrefix, indentString)
	if err != nil {
		return fmt.Errorf("WriteFileIndented(): generate XML: %w", err)
	}

	if err = doWriteFile(xmlBytes, filename); err != nil {
		return fmt.Errorf("WriteFileIndented(): %w", err)
	}

	return nil
}

// FindPathsById traverses an [SVG] object and returns a slice of [PathDef]
// with matching ids.
//
// The supplied id may be a string or a compiled regular expression from
// the [regexp] package.
//
// Use the [FindFirst] or [FindAll] constant as the second parameter to direct
// cause the method to return only the first matching element or all matching
// elements, respectively.
func (S *SVG) FindPathsById(id any, findFirst bool) ([]*PathDef, error) {
	paths, err := traverseObjectById(S, "PathDef", id, findFirst)
	if err != nil {
		err = fmt.Errorf("FindElementsById(): %w", err)
	}
	results := make([]*PathDef, len(paths))
	for i, p := range paths {
		results[i] = p.(*PathDef)
	}
	return results, err
}

// FindGroupsById traverses a [GroupDef] object and returns a slice of
// [PathDef] with matching ids.
//
// The supplied id may be a string or a compiled regular expression from
// the [regexp] package.
//
// Use the [FindFirst] or [FindAll] constant as the second parameter to direct
// cause the method to return only the first matching element or all matching
// elements, respectively.
func (G *GroupDef) FindPathsById(id any, findFirst bool) ([]*PathDef, error) {
	paths, err := traverseObjectById(G, "PathDef", id, findFirst)
	if err != nil {
		err = fmt.Errorf("FindPathsById(): %w", err)
	}
	results := make([]*PathDef, len(paths))
	for i, p := range paths {
		results[i] = p.(*PathDef)
	}
	return results, err
}

// FindGroupsById traverses a [GroupDef] object and returns a slice of
// [GroupDef] with matching ids.
//
// The supplied id may be a string or a compiled regular expression from
// the [regexp] package.
//
// Use the [FindFirst] or [FindAll] constant as the second parameter to direct
// cause the method to return only the first matching element or all matching
// elements, respectively.
func (S *SVG) FindGroupsById(id any, findFirst bool) ([]*GroupDef, error) {
	groups, err := traverseObjectById(S, "GroupDef", id, findFirst)
	if err != nil {
		err = fmt.Errorf("FindGroupsById(): %w", err)
	}
	results := make([]*GroupDef, len(groups))
	for i, g := range groups {
		results[i] = g.(*GroupDef)
	}
	return results, err
}

// AddBackground takes a color definition, such as "#<hex>" or a name like
// "red", and creates a rect matching the image's viewBox or x/y/width/height
// and with the specified fill color.
func (S *SVG) AddBackground(colorDef string) error {
	var bgRect RectDef
	viewBox, err := S.GetViewBox()
	if err != nil {
		return fmt.Errorf("AddBackground(): get viewBox: %w", err)
	} else if viewBox == (ViewBoxDef{}) {
		bgRect = NewRect("svgxml_backgroundColor", "fill:"+colorDef, "0", "0", S.Width, S.Height)
	} else {
		bgRect = NewRect("svgxml_backgroundColor", "fill:"+colorDef, viewBox.X, viewBox.Y, viewBox.Width, viewBox.Height)
	}
	S.background = bgRect
	return nil
}

// RemoveBackground clears any rect added by [SVG.AddBackground]
func (S *SVG) RemoveBackground() {
	S.background = RectDef{}
}

// GetBackground returns the [RectDef] created with [SVG.AddBackground]. If no
// background has been defined, the return value will be a [RectDef] with all
// zero values.
func (S *SVG) GetBackground() RectDef {
	return S.background
}

// GetViewBox returns the value of the "viewBox" XML attribute as a
// [ViewBoxDef] struct of float64 values.
func (S *SVG) GetViewBox() (ViewBoxDef, error) {
	var result [4]float64

	if S.ViewBox == "" {
		return ViewBoxDef{}, nil
	}
	corners := s.Split(s.ReplaceAll(S.ViewBox, ",", " "), " ")
	if len(corners) != 4 {
		return ViewBoxDef{}, fmt.Errorf("invalid viewBox value '%s'", S.ViewBox)
	}
	rInd := 0
	for _, c := range corners {
		if len(c) == 0 {
			continue
		}
		val, err := strconv.ParseFloat(c, 64)
		if err != nil {
			return ViewBoxDef{}, fmt.Errorf("parse viewBox component '%s' as float64: %w", c, err)
		}
		result[rInd] = val
		rInd++
	}

	return ViewBoxDef{X: result[0], Y: result[1], Width: result[2], Height: result[3]}, nil
}

// GetViewBox sets the image's "viewBox" XML attribute from the values in the
// supplied [ViewBoxDef] struct.
func (S *SVG) SetViewBox(box ViewBoxDef) {
	S.ViewBox = fmt.Sprintf("%.6f %.6f %.6f %.6f", box.X, box.Y, box.Width, box.Height)
}
