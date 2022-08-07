package url_shortener

import (
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if _, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, pathsToUrls[r.URL.Path], 301)
		}
		fallback.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
type ymlConf struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func parseYAML(yml []byte) ([]ymlConf, error) {
	var ymlMap []ymlConf
	err := yaml.Unmarshal(yml, &ymlMap)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return ymlMap, err
}

func buildMap(yml []ymlConf) map[string]string {
	paths := make(map[string]string)
	for _, y := range yml {
		paths[y.Path] = y.Url
	}
	return paths
}
