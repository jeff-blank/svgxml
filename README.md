# svgxml

This has grown out of my need for a way in Go to manage the structures of SVG
images in a way that's more straightforward than working with `encoding/xml`
directly and repeating code across projects. It's incredibly incomplete relative
to the SVG spec but works OK for what's implemented.

As this is intended for managing the _structure_ of an image, the focus is on
finding, examining, and modifying existing elements. Adding elements can also
be done easily if not simply; array management (for element ordering) and
generating and setting the values of the XML attributes is left to the user.
For the most part, no major abstractions are provided for creating or modifying
elements.

What little documentation I have is mainly what's been generated from my code
comments; see
[pkg.go.dev](https://pkg.go.dev/github.com/jeff-blank/svgxml).
