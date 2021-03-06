package goblin

import (
	"fmt"
	"strings"
)

func splitPath(path string) ([]string, error) {
	tokens := strings.Split(path, pathSeparator)

	// Special case the root signifier because it's illegal in any other
	// usage and this saves us from doing this check in every iteration
	// of the loop below.
	if len(tokens) == 1 && tokens[0] == filesystemRootPath {
		return tokens, nil
	}

	err := validatePath(tokens)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

func validatePath(path []string) error {
	switch {
	case len(path) == 1 && path[0] == filesystemRootPath:
		// "." is a special case
		return nil
	case len(path) == 0:
		return fmt.Errorf("path cannot be empty")
	case path[0] == pathSeparator:
		return fmt.Errorf("path cannot be an absolute path")
	}

	for _, pathToken := range path {
		trimPathToken := strings.TrimSpace(pathToken)
		switch {
		case trimPathToken == "." || trimPathToken == "..":
			return fmt.Errorf("%s is not allowed in paths", trimPathToken)
		case trimPathToken == "":
			return fmt.Errorf(
				"path cannot contain empty segments: %s",
				strings.Join(path, pathSeparator),
			)
		}
	}

	return nil
}
