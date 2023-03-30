package main

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/narqo/go-badge"
)

func main() {
	badge := getBadge(os.Args[1], os.Args[2])
	fmt.Printf("%s", badge)
}

func getBadge(module, tag string) []byte {
	color := getColor(module)

	badge, err := badge.RenderBytes(tag, module, badge.Color(color.Hex()))
	if err != nil {
		panic(err)
	}

	return badge
}

func getColor(name string) colorful.Color {
	h := fnv.New64()
	h.Write([]byte(name))
	rand.Seed(int64(h.Sum64()))
	return colorful.HappyColor()
}
