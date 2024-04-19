package basic

import (
	"fmt"
	"os"

	"github.com/tdewolff/canvas"
)

func mustLoadGameSVG(name string) *canvas.Canvas {
	c, err := loadGameSVG(name)
	if err != nil {
		panic(err)
	}
	return c
}

func loadGameSVG(name string) (*canvas.Canvas, error) {
	return loadSVG(fmt.Sprintf("assets/Game Symbols/NISEI_%s.svg", name))
}

func mustLoadGameAsset(name string) *canvas.Path {
	path, err := loadGameAsset(name)
	if err != nil {
		panic(err)
	}
	return path
}

func loadGameAsset(name string) (*canvas.Path, error) {
	return loadAsset(fmt.Sprintf("assets/Game Symbols/NISEI_%s.svg", name))
}

func loadAsset(filename string) (*canvas.Path, error) {
	svgData, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	path, err := canvas.ParseSVGPath(string(svgData))
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	return path, nil
}

func loadSVG(filename string) (*canvas.Canvas, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("opening file: %w", err)
	}

	path, err := canvas.ParseSVG(file)
	if err != nil {
		return nil, fmt.Errorf("parsing file: %w", err)
	}

	return path, nil
}
