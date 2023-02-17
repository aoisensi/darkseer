package main

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aoisensi/darkseer/dmx"
	"github.com/qmuntal/gltf"
	"github.com/samber/lo"
)

var argScale = flag.Float64("scale", 0.02, "scale factor")

func main() {
	flag.Parse()
	for _, arg := range flag.Args() {
		for _, name := range lo.Must(filepath.Glob(arg)) {
			if !strings.HasSuffix(name, ".dmx") {
				log.Printf("⚠️ skipping \"%s\" (not a .dmx file)", name)
				continue
			}
			log.Println("ℹ️  start converting", name)
			f, err := os.Open(name)
			if err != nil {
				log.Println("❌", err)
				continue
			}
			defer f.Close()
			element, err := dmx.NewDecoder(f).Decode()
			if err != nil {
				log.Println("❌", err)
				continue
			}
			doc, err := convertModel(element)
			if err != nil {
				log.Println("❌", err)
				continue
			}
			nameGLTF := strings.TrimSuffix(name, ".dmx") + ".gltf"
			if err := gltf.Save(doc, nameGLTF); err != nil {
				log.Println("❌", err)
				continue
			}
			log.Println("✅ saved", nameGLTF)
		}
	}
}
