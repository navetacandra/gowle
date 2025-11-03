package config

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type RegexCache struct {
	Size   int
	Cursor int
	Keys   []string
	Values map[string]*regexp.Regexp
}

func (cache *RegexCache) Get(key string) (*regexp.Regexp, bool) {
	val, found := cache.Values[key]
	return val, found
}

func (cache *RegexCache) Set(key string, regex *regexp.Regexp) {
	if cache.Size == 0 {
		return
	}
	if oldKey := cache.Keys[cache.Cursor]; oldKey != "" {
		delete(cache.Values, cache.Keys[cache.Cursor])
	}

	cache.Keys[cache.Cursor] = key
	cache.Values[key] = regex
	cache.Cursor = (cache.Cursor + 1) % cache.Size
}

type GowleConfig struct {
	Watch      []string
	Ignore     []*regexp.Regexp
	Command    string
	Cwd        string
	regexCache RegexCache
}

func (config *GowleConfig) Load() (err error) {
	// reset state
	config.Watch = config.Watch[:0]
	config.Ignore = config.Ignore[:0]

	// init regexCache
	if config.regexCache.Size == 0 {
		config.regexCache.Size = 32
		config.regexCache.Keys = make([]string, 32)
		config.regexCache.Values = make(map[string]*regexp.Regexp)
	}

	cwd, _ := os.Getwd()
	path := filepath.Join(cwd, ".gowle")
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	config.Cwd = cwd

	sc := bufio.NewScanner(file)
	for sc.Scan() {
		line := sc.Text()
		line = strings.TrimSpace(line)

		if len(line) == 0 || strings.HasPrefix(line, "#") { // skip blank or commented
			continue
		}

		idx := strings.IndexByte(line, byte('='))
		if idx == -1 { // skip no-value var
			continue
		}

		key := line[:idx]
		val := line[idx+1:]

		switch key {
		case "WATCH":
			listParse(&val, &config.Watch)
		case "IGNORE":
			tmp := make([]string, 0, 10)
			listParse(&val, &tmp)
			createRegex(&tmp, &config.regexCache, &config.Ignore)
		case "COMMAND":
			config.Command = val
		}
	}

	return nil
}

func createRegex(src *[]string, cache *RegexCache, res *[]*regexp.Regexp) {
	for _, raw := range *src {
		if val, found := cache.Get(raw); found {
			*res = append(*res, val)
			continue
		}

		s := raw
		if s[0] != '^' { // match from start
			s = "^" + s
		}
		if s[len(s)-1] != '$' { // match untill end
			s = s + "$"
		}
		s = strings.ReplaceAll(s, ".", `\.`) // replace to regex friendly
		s = strings.ReplaceAll(s, "*", `.*`) // replace to regex friendly

		val := regexp.MustCompile(s)
		cache.Set(raw, val)
		*res = append(*res, val)
	}
}

func listParse(src *string, res *[]string) {
	start := 0
	inQuote := false

	for i, r := range *src {
		switch r {
		case '"':
			inQuote = !inQuote
		case ',':
			if !inQuote {
				item := (*src)[start:i]
				item = strings.Trim(item, "\"")
				item = strings.TrimSpace(item)
				if item != "" {
					*res = append(*res, item)
				}
				start = i + 1
			}
		}
	}

	if start < len(*src) {
		item := (*src)[start:]
		item = strings.Trim(item, "\"")
		item = strings.TrimSpace(item)
		if item != "" {
			*res = append(*res, item)
		}
	}
}
