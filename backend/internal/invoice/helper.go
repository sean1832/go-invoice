package invoice

import (
	"path/filepath"
	"strconv"
)

// findMaxSuffix finds the maximum numeric suffix from a list of filenames.
// Expected filename format: INV-YYMMDDSS.json, where SS is a two-digit suffix.
func FindMaxSuffixFromFilename(filenames []string) int {
	maxSuffix := 0
	for _, filename := range filenames {
		// INV-25110203.json
		base := filepath.Base(filename)

		// trim the extensions. Also handle case like `.json.bak`
		idStr := base[0 : len(base)-len(filepath.Ext(base))]
		if len(idStr) < 2 {
			continue // malformed, skip
		}

		suffixStr := idStr[len(idStr)-2:]

		suffix, err := strconv.Atoi(suffixStr)
		if err != nil {
			continue // not numeric, skip
		}

		if suffix > maxSuffix {
			maxSuffix = suffix
		}
	}

	return maxSuffix
}
