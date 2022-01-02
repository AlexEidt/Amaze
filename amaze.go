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

	// List of font names.
	font_names := []string{
		"boxwood",
		"labyrinth",
		"metropolis",
		"palazzo",
		"path",
		"sandcastle",
		"underworld",
		"utopia",
		"ziggurat",
	}

	// Parse command line args
	fontType := flag.String("font", "metropolis", "Font to write letters in. One of: "+strings.Join(font_names, ", "))
	random := flag.Bool("r", false, "Generate Random Patterns.")
	chars := flag.String("chars", charset, "Character set to use for patterns.")
	background := flag.String("bg", "255,255,255,255", "Background color in RGBA format.")
	foreground := flag.String("fg", "0,0,0,255", "Foreground color in RGBA format.")
	fontSize := flag.Int("size", 100, "Font Size.")
	w := flag.Int("w", 1920, "Image width.")
	h := flag.Int("h", 1080, "Image height.")
	animate := flag.Bool("a", false, "Create one frame per character in character set.")
	limit := flag.Int("l", len(*chars), "Limit the number of animation frames created.")

	flag.Parse()
	args := flag.Args()
	filename := args[0]

	bR, bG, bB, bA := ProcessRGB(*background)
	fR, fG, fB, fA := ProcessRGB(*foreground)
	canvas := gg.NewContext(*w, *h)
	canvas.SetRGBA255(bR, bG, bB, bA)
	canvas.Clear()
	canvas.SetRGBA255(fR, fG, fB, fA)

	font := GetFont(filepath.Join("Fonts", *fontType+".ttf"))
	face := truetype.NewFace(
		font,
		&truetype.Options{
			// Round *fontSize to the nearest 10.
			Size: float64(int(math.Ceil(float64(*fontSize)/10.0)) * 10),
		},
	)
	canvas.SetFontFace(face)

	if *random {
		runes := []rune(*chars)
		// Randomly shuffle input string.
		rand.Shuffle(len(runes), func(i, j int) {
			runes[i], runes[j] = runes[j], runes[i]
		})
		*chars = string(runes)
	}

	if *animate {
		text := ""
		for i, char := range *chars {
			text += string(char)
			canvas.SetRGBA255(fR, fG, fB, fA)
			Amaze(canvas, strconv.Itoa(i)+filename, text, *random)
			if i > *limit {
				break
			}
			canvas.SetRGBA255(bR, bG, bB, bA)
			canvas.Clear()
		}
	} else {
		Amaze(canvas, filename, *chars, *random)
	}
}

// Create the maze graphic given certain text.
func Amaze(canvas *gg.Context, filename, text string, random bool) {
	w, h := canvas.Width(), canvas.Height()
	size_w, size_h := canvas.MeasureString(string(text[0]))
	index := 0
	var r *rand.Rand
	if random {
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		index = r.Intn(len(text))
	}
	row := 0.0

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

// Parses RGBA string into separate R, G, B, A values between 0 and 255.
func ProcessRGB(rgba string) (int, int, int, int) {
	colors := strings.Split(rgba, ",")
	if len(colors) != 4 {
		panic(rgba + "is not a valid RGBA string.")
	}

	R, _ := strconv.Atoi(colors[0])
	G, _ := strconv.Atoi(colors[1])
	B, _ := strconv.Atoi(colors[2])
	A, _ := strconv.Atoi(colors[3])

	return R, G, B, A
}
