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
		yamlpath := path.Clean(c.Filenames[i])
		j := transformToJSON(yamlpath)

		if c.NewFile {
			noext := strings.TrimSuffix(yamlpath, path.Ext(yamlpath))
			jsonpath := noext + ".json"
			ioutil.WriteFile(jsonpath, []byte(j), 0644)
		} else {
			fmt.Print(j)
		}
	}

}

func transformToJSON(yamlpath string) string {
	bytes, err := ioutil.ReadFile(yamlpath)
	if err != nil {
		panic(err)
	}

	j, err := yaml.YAMLToJSON(bytes)
	if err != nil {
		panic(err)
	}

	return string(j)
}
