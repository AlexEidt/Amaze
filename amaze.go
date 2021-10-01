// Alex Eidt
// Amaze Maze Generator.
// Check out the MazeLetter font at http://mazeletter.xyz/.

package main

import (
	"flag"
	"io/ioutil"
	"math"
	"math/rand"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
)

func main() {

	charset := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

	// Parse command line args
	fontType := flag.String("font", "metropolis", "Font to write letters in. One of: boxwood, labyrinth, metropolis, palazzo, path, sandcastle, underworld, utopia, ziggurat.")
	random := flag.Bool("r", false, "Generate Random Patterns.")
	chars := flag.String("chars", charset, "Character set to use for patterns.")
	background := flag.String("bg", "255,255,255", "Background color in format 'R,G,B' or in hex format.")
	bgOpacity := flag.Int("bgo", 255, "Background Color Opacity.")
	foreground := flag.String("fg", "0,0,0", "Foreground color in format 'R,G,B' or in hex format.")
	fgOpacity := flag.Int("fgo", 255, "Foreground Color Opacity.")
	fontSize := flag.Int("size", 100, "Font Size.")
	w := flag.Int("w", 1920, "Image width.")
	h := flag.Int("h", 1080, "Image height.")
	animate := flag.Bool("a", false, "Create one frame per character in character set.")
	limit := flag.Int("l", len(*chars), "Limit the number of animation frames created.")

	// Round fontSize to the nearest 10.
	*fontSize = int(math.Ceil(float64(*fontSize)/10.0)) * 10

	flag.Parse()
	args := flag.Args()
	filename := args[0]

	font := GetFont(filepath.Join("Fonts", *fontType+".ttf"))

	canvas := gg.NewContext(*w, *h)

	bR, bG, bB := ProcessRGB(*background)
	canvas.SetRGBA255(bR, bG, bB, *bgOpacity)
	canvas.Clear()

	face := truetype.NewFace(font, &truetype.Options{Size: float64(*fontSize)})
	canvas.SetFontFace(face)

	// groups := GroupSizes(canvas, charset)

	characters := *chars

	fR, fG, fB := ProcessRGB(*foreground)

	// Random Number Generator.
	var r *rand.Rand
	if *random {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	} else {
		r = nil
	}

	if *animate {
		text := ""
		count := 0
		for _, char := range characters {
			if *random {
				text += string((characters)[r.Intn(len(characters))])
			} else {
				text += string(char)
			}
			Amaze(canvas, r, strconv.Itoa(count)+filename, text, *random, fR, fG, fB, *fgOpacity)
			count++
			if count > *limit {
				break
			}
			canvas.SetRGBA255(bR, bG, bB, *bgOpacity)
			canvas.Clear()
		}
	} else {
		Amaze(canvas, r, filename, *chars, *random, fR, fG, fB, *fgOpacity)
	}
}

// Create the maze graphic given certain text.
func Amaze(canvas *gg.Context, r *rand.Rand, filename, text string, random bool, R, G, B, A int) {
	w, h := canvas.Width(), canvas.Height()

	size_w, size_h := canvas.MeasureString(string(text[0]))
	index := 0
	if random {
		index = r.Intn(len(text))
	}
	row := 0.0

	canvas.SetRGBA255(R, G, B, A)
	for row < float64(h)+size_h {
		column := 0.0
		for column < float64(w) {
			char := string(text[index])
			size_w, size_h = canvas.MeasureString(char)
			canvas.DrawString(char, column, row)
			if random {
				index = r.Intn(len(text))
			} else {
				index = (index + 1) % len(text)
			}
			column += size_w
		}
		row += size_h
	}
	canvas.SavePNG(filename)
}

// Parses the ttf file for the maze fonts into a Font object.
func GetFont(filename string) *truetype.Font {
	ttf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	font, err := truetype.Parse(ttf)
	if err != nil {
		panic(err)
	}
	return font
}

// Goes through all characters and groups characters with the same size.
// Use this function if some of the characters are not lining up.
func GroupSizes(canvas *gg.Context, charset string) []string {
	type Pair struct{ width, height float64 }
	sizes := make(map[Pair]string)
	for _, char := range charset {
		w, h := canvas.MeasureString(string(char))
		pair := Pair{w, h}
		if _, ok := sizes[pair]; !ok {
			sizes[pair] = ""
		}
		sizes[Pair{w, h}] += string(char)
	}
	groups := make([]string, len(sizes))
	index := 0
	for _, chars := range sizes {
		groups[index] = chars
		index++
	}
	return groups
}

// Parses RGB or hex strings into three separate R, G, B values between 0 and 255.
func ProcessRGB(rgb string) (int, int, int) {
	hex := map[byte]int{
		'0': 0, '1': 1, '2': 2, '3': 3,
		'4': 4, '5': 5, '6': 6, '7': 7,
		'8': 8, '9': 9, 'A': 10, 'B': 11,
		'C': 12, 'D': 13, 'E': 14, 'F': 15,
	}
	rgb = strings.ToUpper(rgb)
	if strings.HasPrefix(rgb, "#") {
		R := hex[rgb[1]]*16 + hex[rgb[2]]
		G := hex[rgb[3]]*16 + hex[rgb[4]]
		B := hex[rgb[5]]*16 + hex[rgb[6]]
		return R, G, B
	}
	colors := strings.Split(rgb, ",")
	if len(colors) != 3 {
		panic(rgb + "is not a valid RGB string.")
	}

	R, _ := strconv.Atoi(colors[0])
	G, _ := strconv.Atoi(colors[1])
	B, _ := strconv.Atoi(colors[2])

	return R, G, B
}
