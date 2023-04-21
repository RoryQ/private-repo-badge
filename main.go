package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/lucasb-eyer/go-colorful"
)

const (
	// This release tag will receive pngs for each tag discovered in the readme table.
	badgesReleaseTag = "readmebadges"
	parentModPath    = "github.com/RoryQ/private-repo-badge"
)

func main() {

	out, err := exec.Command("git", "tag").Output()
	if err != nil {
		panic(err)
	}

	tags := groupByTagPrefix(string(out))

	contents, err := os.ReadFile("README.md")
	if err != nil {
		panic(err)
	}

	table := readmeTable(contents)

	tablePrefixes := tagPrefixesFromTable(table)

	for _, prefix := range tablePrefixes {
		latest, ok := tags[prefix]
		if !ok {
			panic(prefix)
		}

		sortVersions(latest)

		_ = saveBadge(latest[0])

	}
}

func saveBadge(latestTag string) *os.File {
	badge := getBadge(latestTag)
	prefix, _ := tagToComponents(latestTag)
	filename := fmt.Sprintf("%s.png", escapeFilename(prefix))
	file, err := os.Create(filename)
	if err != nil {
		panic(err.Error())
	}
	_, err = file.Write(badge)
	if err != nil {
		panic(err.Error())
	}
	file.Seek(0, 0)
	return file
}

func getBadge(latestTag string) []byte {
	name, version := tagToComponents(latestTag)
	colour := getColor(parentModPath + name)

	hex := strings.TrimPrefix(colour.Hex(), "#")
	imgUrl := fmt.Sprintf("https://raster.shields.io/badge/%s-%s-%s", name, version, hex)

	println(imgUrl)

	resp, err := http.Get(imgUrl)
	if err != nil {
		panic(err.Error())
	}
	defer resp.Body.Close()

	out := bytes.Buffer{}
	io.Copy(&out, resp.Body)
	return out.Bytes()
}

func getColor(name string) colorful.Color {
	h := fnv.New64()
	h.Write([]byte(name))
	rand.Seed(int64(h.Sum64()))
	return colorful.HappyColor()
}

func sortVersions(versions []string) {
	destructure := func(s string) []int {
		_, version := tagToComponents(s)
		components := strings.Split(strings.TrimPrefix(version, "v"), ".")
		if len(components) != 3 {
			panic("not a semver? " + s)
		}

		return apply(components, func(item string, index int) int {
			return must(strconv.Atoi(item))
		})
	}

	sort.Slice(versions, func(i, j int) bool {
		return recursiveCompare(
			destructure(versions[i]),
			destructure(versions[j]),
		) > 0
	})
}

func recursiveCompare(versionA []int, versionB []int) int {
	if len(versionA) == 0 {
		return 0
	}

	a := versionA[0]
	b := versionB[0]

	if a > b {
		return 1
	} else if a < b {
		return -1
	}

	return recursiveCompare(versionA[1:], versionB[1:])
}

func tagToComponents(tag string) (prefix string, version string) {
	components := strings.Split(tag, "/")
	version = components[len(components)-1]
	prefix = strings.TrimSuffix(tag, "/"+version)
	return
}

func groupByTagPrefix(tags string) map[string][]string {
	lines := strings.Split(tags, "\n")

	groupBy := map[string][]string{}
	for _, line := range lines {
		prefix, _ := tagToComponents(line)
		groupBy[prefix] = append(groupBy[prefix], line)
	}

	return groupBy
}

func tagPrefixesFromTable(table string) []string {
	lines := strings.Split(table, "\n")

	tags := []string{}
	for i := 2; i < len(lines); i++ {
		cols := strings.Split(lines[i], "|")
		badgeCol := cols[len(cols)-2]
		tag := regexp.MustCompile(`download/`+badgesReleaseTag+`/(.*).png"`).FindAllStringSubmatch(badgeCol, 1)[0][1]
		tag = unescapeFilename(tag)
		tags = append(tags, tag)
	}

	return tags
}

func escapeFilename(s string) string {
	return strings.ReplaceAll(s, "/", "__")
}

func unescapeFilename(s string) string {
	return strings.ReplaceAll(s, "__", "/")
}

func readmeTable(contents []byte) string {
	re := regexp.MustCompile(`(?m)^[|]([^|]+[|])*\s+Latest Tag\s+\|\n^[|]([^|]+[|])*`)
	return re.FindAllString(string(contents), -1)[0]
}

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func apply[T any, R any](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, len(collection))

	for i, item := range collection {
		result[i] = iteratee(item, i)
	}

	return result
}
