package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/jpillora/opts"
)

type config struct {
	NewFile   bool     `opts:"help=create new file instead of writing to STDOUT"`
	Filenames []string `opts:"mode=arg,min=1"`
}

func main() {
	c := config{}
	opts.Parse(&c)

	for i := range c.Filenames {
		jsonpath := c.Filenames[i]
		y := transformToYAML(jsonpath)

		if c.NewFile {
			noext := strings.TrimSuffix(jsonpath, path.Ext(jsonpath))
			yamlpath := noext + ".yaml"
			ioutil.WriteFile(yamlpath, []byte(y), 0644)
		} else {
			fmt.Print(y)
		}
	}

}

func transformToYAML(jsonpath string) string {
	bytes, err := ioutil.ReadFile(jsonpath)
	if err != nil {
		panic(err)
	}

	y, err := yaml.JSONToYAML(bytes)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("---\n%s", y)
}
