package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ResourceData is an interface that all resource types (Client, Provider) must implement
type ResourceData interface {
	SetID(id string)
	HasRequiredFields() bool
}

// getResourceByID handles GET request for a single resource by ID
func getResourceByID(
	w http.ResponseWriter,
	r *http.Request,
	storageDir string,
	resourceType string,
	newResource func() ResourceData,
) {
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, fmt.Sprintf("%s ID is required", resourceType), http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storageDir, id+".json")
	resource := newResource()

	if err := readJSON(filePath, resource); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, fmt.Sprintf("%s not found for '%s'", resourceType, id), http.StatusNotFound)
		} else {
			writeRespErr(w, fmt.Sprintf("failed to read %s '%s'", resourceType, id), http.StatusInternalServerError)
		}
		return
	}

	writeRespOk(w, fmt.Sprintf("%s '%s'", resourceType, id), resource)
}

// updateResourceByID handles PUT request to update an existing resource
func updateResourceByID(
	w http.ResponseWriter,
	r *http.Request,
	storageDir string,
	resourceType string,
	newResource func() ResourceData,
) {
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, fmt.Sprintf("%s ID is required", resourceType), http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storageDir, id+".json")

	exists, err := isPathExist(filePath)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("failed to check %s '%s' existence: %v", resourceType, id, err), http.StatusInternalServerError)
		return
	}
	if !exists {
		writeRespErr(w, fmt.Sprintf("%s not found for '%s'", resourceType, id), http.StatusNotFound)
		return
	}

	resource := newResource()
	if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
		writeRespErr(w, fmt.Sprintf("invalid %s data for '%s': %v", resourceType, id, err), http.StatusBadRequest)
		return
	}

	if err := writeJSON(filePath, resource, 2); err != nil {
		writeRespErr(w, fmt.Sprintf("failed to update %s '%s'", resourceType, id), http.StatusInternalServerError)
		return
	}

	writeRespOk(w, fmt.Sprintf("updated %s '%s'", resourceType, id), resource)
}

// deleteResourceByID handles DELETE request to remove a resource
func deleteResourceByID(
	w http.ResponseWriter,
	r *http.Request,
	storageDir string,
	resourceType string,
) {
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, fmt.Sprintf("%s ID is required", resourceType), http.StatusBadRequest)
		return
	}

	filePath := filepath.Join(storageDir, id+".json")

	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, fmt.Sprintf("%s not found for '%s'", resourceType, id), http.StatusNotFound)
		} else {
			writeRespErr(w, fmt.Sprintf("failed to delete %s '%s'", resourceType, id), http.StatusInternalServerError)
		}
		return
	}

	writeRespOk(w, fmt.Sprintf("deleted %s '%s'", resourceType, id), nil)
}

// getAllResources handles GET request for listing all resources
func getAllResources(
	w http.ResponseWriter,
	storageDir string,
	resourceType string,
	getAll func(string) (any, error),
) {
	resources, err := getAll(storageDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, fmt.Sprintf("%s resource is currently empty", resourceType), http.StatusNotFound)
		} else {
			writeRespErr(w, fmt.Sprintf("failed to list %s informations", resourceType), http.StatusInternalServerError)
		}
		return
	}

	writeRespOk(w, fmt.Sprintf("list of %ss", resourceType), resources)
}

// createResource handles POST request to create a new resource
func createResource(
	w http.ResponseWriter,
	r *http.Request,
	storageDir string,
	resourceType string,
	newResource func() ResourceData,
) {
	if r.Body == nil || r.ContentLength == 0 {
		writeRespErr(w, "request body is empty", http.StatusBadRequest)
		return
	}

	resource := newResource()
	if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
		writeRespErr(w, fmt.Sprintf("invalid %s data", resourceType), http.StatusBadRequest)
		return
	}

	if !resource.HasRequiredFields() {
		writeRespErr(w, fmt.Sprintf("incomplete %s data", resourceType), http.StatusBadRequest)
		return
	}

	// Extract ID from the resource, or generate from name
	var tempData map[string]interface{}
	bodyBytes, _ := json.Marshal(resource)
	json.Unmarshal(bodyBytes, &tempData)
	id, ok := tempData["id"].(string)
	if !ok || id == "" {
		// if ID is not provided, generate from name
		name, ok := tempData["name"].(string)
		if !ok || name == "" {
			writeRespErr(w, "name or id field is required", http.StatusBadRequest)
			return
		}
		id = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	}

	filePath := filepath.Join(storageDir, id+".json")

	exists, err := isPathExist(filePath)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("failed to check %s '%s' existence: %v", resourceType, id, err), http.StatusInternalServerError)
		return
	}
	if exists {
		writeRespErr(w, fmt.Sprintf("%s already exists for '%s'", resourceType, id), http.StatusConflict)
		return
	}

	if err := writeJSON(filePath, resource, 2); err != nil {
		writeRespErr(w, fmt.Sprintf("failed to create %s '%s'", resourceType, id), http.StatusInternalServerError)
		return
	}

	writeRespWithStatus(w, fmt.Sprintf("created %s '%s'", resourceType, id), resource, http.StatusCreated)
}

// countResources counts the number of resource files in the given storage directory
func countResources(storageDir string) (int, error) {
	jsonFiles, err := filepath.Glob(filepath.Join(storageDir, "*.json"))
	if err != nil {
		return 0, err
	}
	return len(jsonFiles), nil
}
