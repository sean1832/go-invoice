package invoice

import (
	"encoding/json"
	"fmt"
	"os"
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

func LoadInvoice(invoiceRoot string, id string) (*Invoice, error) {
	inv := &Invoice{}
	filepath := filepath.Join(invoiceRoot, id+".json")
	err := loadResourceFromFile(filepath, inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
}

func SaveInvoice(invoiceRoot string, inv *Invoice) error {
	filepath := filepath.Join(invoiceRoot, inv.ID+".json")
	data, err := json.MarshalIndent(inv, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal invoice '%s': %w", inv.ID, err)
	}
	err = os.WriteFile(filepath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write invoice file '%s': %w", filepath, err)
	}
	return nil
}

func loadResourceFromFile(filepath string, resource any) error {
	// read json
	file, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file '%s': %w", filepath, err)
	}
	return json.Unmarshal(file, resource)
}
