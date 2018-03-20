// Package banner is a dynamic ASCII art banner generator.
//
// This package supports loading fonts at runtime and can parse
// arbitrary ASCII art fonts.
package banner

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

// GlyphTable describes the ordering and formatting expected by the
// banner package when parsing font data.
//
// This table can be used to generate new fonts using patorjk's timeless
// ASCCI generator at http://patorjk.com/software/taag/. Copy this table
// and paste it into the generator to get a nice chunk of ASCII goodness
// that can be used as a font with this package.
//
// You may need to add an extra new-line between each row to ensure
// the parser can separate them correctly.
const GlyphTable = `{ ! @ # $ % ^ & * 0
{ A a B b C c D d E e 1
{ F f G g H h I i J j 2
{ K k L l M m N n O o 3
{ P p Q q R r S s T t 4
{ U u V v +   X x Y   5
{ Z z - = _ W [ ]   } 6
{ 0 1 2 3 4 y      w  7
{ ( ) : ; " ' ? / \ 8 9`

var DefaultFont *Font

const numRows = 9

var base = [numRows][]byte{}

func init() {
	// Initialize the package level reference data structure.
	refLines := strings.Split(GlyphTable, "\n")
	for i, l := range refLines {
		base[i] = []byte(l)
	}
	var err error
	DefaultFont, err = NewFontString(fontRectangles)
	if err != nil {
		panic(err)
	}
}

// A Font is a data type that contains information necessary
// to render text using arbitrary ASCII art fonts.
type Font struct {
	size    int
	encoded map[byte]*letter
}

// NewFont attempts to read the font data stored in file.
// If this function cannot parse the font data, an error will be
// returned.
//
// Fonts can be created using http://patorjk.com/software/taag/.
// Select a font, then paste in the GlyphsTable text. Copy the output
// and place it in a file or use NewFontString to parse the font.
func NewFont(file string) (*Font, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return NewFontString(string(r))
}

// NewFontString attempts to parse the font data passed in
// as a string. If this function cannot parse the font data, an
// error will be returned.
func NewFontString(font string) (*Font, error) {
	fo := &Font{
		encoded: make(map[byte]*letter),
	}
	return fo, fo.generate(font)
}

func (f *Font) Size() int {
	return f.size
}

// A letter describes the width, height, and make up of a font glyph.
type letter struct {
	width  int
	height int
	lines  []string
}

// A Sequence is a sequence of letters representing an ASCII banner.
type Sequence []*letter

// New creates a new ASCII art banner using the given font and string input.
func New(f *Font, s string) Sequence {
	le := make(Sequence, len(s))
	q := f.encoded['?']
	for i, b := range []byte(s) {
		if l, ok := f.encoded[b]; ok {
			le[i] = l
		} else {
			le[i] = q
		}
	}
	return le
}

func (l Sequence) String() string {
	if len(l) == 0 {
		return ""
	}
	var res string
	for i := 0; i < len(l[0].lines); i++ {
		for _, le := range l {
			// TODO: Normalizing the widths should probably have already occurred.
			ln := le.lines[i]
			if len(ln) < le.width {
				ln = ln + strings.Repeat(" ", le.width-len(ln))
			}
			res = res + ln
		}
		res = res + "\n"
	}
	return res
}

func (l letter) String() string {
	return strings.Join(l.lines, "\n")
}

func parseRows(raw string) (lines [numRows][]byte, err error) {
	curr := 0
	seenNonEmpty := false
	for _, l := range strings.Split(raw, "\n") {
		b := []byte(l)
		empty := true
		for j := 0; j < len(l); j++ {
			c := l[j]
			switch c {
			case ' ', '_', '.', ',':
				// empty

			default:
				empty = false
				break
			}
		}
		if !empty {
			seenNonEmpty = true
			lines[curr] = append(lines[curr], append([]byte("\n"), b...)...)
			continue
		}
		if seenNonEmpty && curr < len(lines)-1 {
			curr++
			seenNonEmpty = false
		}
		if lines[curr] != nil && len(lines[curr]) > 0 {
			lines[curr] = append(lines[curr], []byte("\n")...)
		}
		lines[curr] = append(lines[curr], b...)
	}
	if curr+1 != len(lines) {
		return lines, errors.New("couldn't parse ascii raw properly: ensure each row is separated by at least one completely blank line")
	}
	return lines, nil
}

// generate parses and prepares the font for use in rendering text.
func (f *Font) generate(raw string) error {
	var err error

	// Parse the font data into rows of glyphs.
	rawRows, err := parseRows(raw)
	if err != nil {
		return err
	}

	// Gather the dimensions of the font.
	heights := make([]int, len(rawRows))
	lengths := make([]int, len(rawRows))
	lowH := 10
	maxH := 0
	for i, l := range rawRows {
		pli := bytes.Split(l, []byte("\n"))
		for _, line := range pli {
			if len(line) > lengths[i] {
				lengths[i] = len(line)
			}
		}
		h := len(pli)
		heights[i] = h
		if h > maxH {
			maxH = h
		} else if h < lowH {
			lowH = h
		}
	}

	// Normalize the size of the font.
	if lowH != maxH {
		// TODO: There still seems to be issues here.
		// Best bet is to try to ensure all the rows parse to the same heights.
		// The "faughnt" format's ordering of characters is deliberate to help
		// ensure the letters end up aligned: the { at the beginning of each line
		// usually sets the pace for the rest of the line nicely.
		targetH := maxH
		for i, lb := range rawRows {
			l := append([]byte{}, bytes.Trim(lb, "\n")...)
			h := bytes.Count(l, []byte("\n"))
			if h < targetH {
				// This row is too short, add lines of blank space to the
				// end until it's the right size.
				for h < targetH {
					l = append(append(l, []byte("\n")...), bytes.Repeat([]byte(" "), lengths[i])...)
					h++
				}
			} else if targetH > h {
				// This row is too tall, remove the first line and hope that's
				// enough.
				idx := bytes.Index(rawRows[i], []byte("\n"))
				l = l[idx+1:]
			}
			// Update font data.
			rawRows[i] = l
			// And font metadata.
			heights[i] = bytes.Count(rawRows[i], []byte("\n"))
		}
	}

	// Take for granted that the heights are identical from now on.
	f.size = heights[0]

	// Separate individual font glyphs.
	for i, l := range rawRows {
		// Remove spaces from the reference text to simplify iterating over it.
		ref := bytes.Replace(base[i], []byte(" "), []byte{}, -1)

		// p is the current column being worked on.
		p := 0
		currChar := ref[p]

		// Split up the row's data
		lns := bytes.Split(l, []byte("\n"))

		// The tailing column position that has been processed.
		last := 0
		// Currently skipping columns (last column processed was all spaces)
		skip := true

		// Check each column for white-space. Use this as a signal
		// to determine when a character begins and ends.
		for j := 0; j < lengths[i]; j++ {
			// Because characters span multiple lines, check the entire height
			// of the column for whitespace.
			spaceColumn := true
			for k := 0; k < f.size; k++ {
				ln := lns[k]
				if len(ln) > j && ln[j] != ' ' {
					// Found something that wasn't a space!
					spaceColumn = false
					// No need to continue checking for non-space.
					break
				}
			}

			if !spaceColumn {
				// Not a space, are we already skipping?
				if skip {
					// Yes, this is the start of a new letter.
					last = j
					skip = false
					continue
				}
				// Otherwise, we already knew we were inside a letter, and we are still in it.
				// No need to do anything but loop back around and process the next column.
				continue
			}

			// Current column is all spaces. Are we already skipping?
			if !skip {
				// No, this is the end of a letter.
				// Start skipping whitespace.
				skip = true
				// And add everything from last to j in all rows to the current glyph.
				ltr := &letter{}
				w := 0
				// Again, operate on the entire height and width of the row.
				for k := 0; k < f.size; k++ {
					ln := lns[k]
					var lne string
					// Be sure to handle short lines in the mix of things.
					if len(ln) > j {
						lne = string(ln[last:j])
					} else if last > len(ln) {
						// The last column was already past the end of this line.
						lne = ""
					} else {
						// Add everything since last time, even though it's a bit short.
						lne = string(ln[last:])
					}
					ltr.lines = append(ltr.lines, lne)
					if len(lne) > w {
						w = len(lne)
					}
				}
				// Done encoding, set the width and size.
				ltr.width = w
				ltr.height = f.size
				// Add the letter to our encoding.
				f.encoded[currChar] = ltr
				// And move on.
				p++
				if p >= len(ref) {
					// We've reached the end of the current row in rawLines.
					break
				}
				// Otherwise, set the next iteration up
				currChar = ref[p]
				last = j
				// And carry on with the next character in the current row.
				continue
			}
			// We were already skipping and we found another blank column.
			// Keep on skipping.
			last = j
		}
	}

	// Don't forget the space!
	f.encoded[' '] = &letter{width: 5, height: f.size, lines: make([]string, f.size)}
	return nil
}
