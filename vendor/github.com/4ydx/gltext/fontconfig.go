// Copyright 2012 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gltext

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"time"
)

// Direction represents the direction in which strings should be rendered.
type Direction uint8

// FontConfig describes raster font metadata.
//
// It can be loaded from, or saved to a JSON encoded file,
// which should come with any bitmap font image.
type FontConfig struct {
	// The range of glyphs covered by this fontconfig
	// An array of Low, High values allowing the user to select disjoint subsets of the ttf
	RuneRanges RuneRanges

	// Glyphs holds a set of glyph descriptors, defining the location,
	// size and advance of each glyph in the sprite sheet.
	Glyphs Charset

	Image *image.NRGBA `json:"-"`

	Name string
}

// Load reads font configuration data from the given JSON encoded stream.
func (fc *FontConfig) Load(rootPath string) (err error) {
	file := fmt.Sprintf("%s/%s.config", rootPath, fc.Name)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, fc)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", time.Now())
	fc.Image, err = LoadFontImage(rootPath, fc.Name)
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", time.Now())
	fc.Glyphs.Scale(1)
	return nil
}

// Save writes font configuration data to the given stream as JSON data.
func (fc *FontConfig) Save(rootPath string) error {
	if _, err := os.Stat(rootPath); err != nil {
		if os.IsNotExist(err) {
			os.MkdirAll(rootPath, os.ModeDir|os.ModePerm)
		} else {
			return err
		}
	}
	data, err := json.Marshal(fc)
	if err != nil {
		return err
	}
	file := fmt.Sprintf("%s/%s.config", rootPath, fc.Name)
	err = ioutil.WriteFile(file, data, os.ModePerm)
	if err != nil {
		return err
	}
	if fc.Image == nil {
		return errors.New("Should not be nil.")
	}
	err = SaveImage(rootPath, fc.Name, fc.Image)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(file, data, os.ModePerm)
	return err
}

func LoadFontImage(rootPath, name string) (*image.NRGBA, error) {
	file := fmt.Sprintf("%s/%s.png", rootPath, name)
	return LoadImage(file)
}

func SaveImage(rootPath, name string, img *image.NRGBA) error {
	file := fmt.Sprintf("%s/%s.png", rootPath, name)
	image, err := os.Create(file)
	if err != nil {
		return err
	}
	defer image.Close()

	b := bufio.NewWriter(image)
	err = png.Encode(b, img)
	if err != nil {
		return err
	}
	return b.Flush()
}
