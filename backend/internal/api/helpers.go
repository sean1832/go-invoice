package api

import (
	"encoding/json"
	"fmt"
	"invoice/internal/invoice"
	"os"
	"path/filepath"
	"strings"
)

func readJSON(path string, data any) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(fileData, data)
}

func writeJSON(path string, data any, indent int) error {
	fileData, err := json.MarshalIndent(data, "", strings.Repeat(" ", indent))
	if err != nil {
		return err
	}
	return os.WriteFile(path, fileData, 0644)
}

func isPathExist(path string) (bool, error) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

type identifiable interface {
	SetID(id string)
}

func getAllProfiles[T identifiable](root string) ([]T, error) {
	jsonFiles, err := filepath.Glob(filepath.Join(root, "*.json"))
	if err != nil {
		return nil, err
	}
	fileNum := len(jsonFiles)
	if fileNum <= 0 {
		return nil, os.ErrNotExist
	}
	var profiles = make([]T, fileNum)
	for i, file := range jsonFiles {
		var profile T
		if err := readJSON(file, &profile); err != nil {
			return nil, fmt.Errorf("failed to read profile data: %v", err)
		}
		filename := filepath.Base(file)
		ext := filepath.Ext(filename)
		id := strings.TrimSuffix(filename, ext)
		profile.SetID(id)
		profiles[i] = profile
	}

	return profiles, nil
}

func getAllInvoices(root string) ([]invoice.Invoice, error) {
	jsonFiles, err := filepath.Glob(filepath.Join(root, "*.json"))
	if err != nil {
		return nil, err
	}
	fileNum := len(jsonFiles)
	if fileNum <= 0 {
		return nil, os.ErrNotExist
	}
	var invoices = make([]invoice.Invoice, fileNum)
	for i, file := range jsonFiles {
		var invoice invoice.Invoice
		if err := readJSON(file, &invoice); err != nil {
			return nil, fmt.Errorf("failed to read profile data: %v", err)
		}
		invoices[i] = invoice
	}
	return invoices, nil
}
