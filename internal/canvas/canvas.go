package canvas

import (
	"board/internal/color"
	"fmt"
	"sync"
)

type Canvas struct {
	width, height int
	fields        [][]color.Color
	mtx           sync.RWMutex
}

func NewCanvas(
	width, height int,
) *Canvas {
	fields := make([][]color.Color, width)
	for i := range fields {
		fields[i] = make([]color.Color, height)
	}
	return &Canvas{
		width:  width,
		height: height,
		fields: fields,
	}
}

func (c *Canvas) Clear() {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	c.fields = make([][]color.Color, c.width)
	for i := range c.fields {
		c.fields[i] = make([]color.Color, c.height)
	}
}

func (c *Canvas) Paint(x, y int, paintColor color.Color) error {
	c.mtx.Lock()
	defer c.mtx.Unlock()
	if x >= len(c.fields) {
		return fmt.Errorf("X out of range")
	}
	if y >= len(c.fields[x]) {
		return fmt.Errorf("Y out of range")
	}
	c.fields[x][y] = paintColor

	return nil
}

func (c *Canvas) Snapshot() [][]color.Color {
	c.mtx.RLock()
	defer c.mtx.RUnlock()
	newFields := make([][]color.Color, c.width)
	for i := range newFields {
		newFields[i] = make([]color.Color, c.height)
		copy(newFields[i], c.fields[i])
	}
	return newFields
}
