package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-invoice/internal/invoice"
	"go-invoice/internal/types"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type resourceType string

const (
	InvoiceType  resourceType = "invoice"
	ClientType   resourceType = "client"
	ProviderType resourceType = "provider"
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
	resourceType resourceType,
	newResource func() ResourceData,
) {
	logger := slog.With("url", r.RequestURI, "method", r.Method)
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
			logger.Error("resource not found", "error", err)
		} else {
			writeRespErr(w, fmt.Sprintf("failed to read %s '%s'", resourceType, id), http.StatusInternalServerError)
			logger.Error("failed to read resource", "error", err)
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
	resourceType resourceType,
	newResource func() ResourceData,
) {
	logger := slog.With("url", r.RequestURI, "method", r.Method)
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, fmt.Sprintf("%s ID is required", resourceType), http.StatusBadRequest)
		logger.Error("failed to parse ID from URL")
		return
	}

	filePath := filepath.Join(storageDir, id+".json")

	exists, err := isPathExist(filePath)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("failed to check %s '%s' existence: %v", resourceType, id, err), http.StatusInternalServerError)
		logger.Error("failed to check resource existence", "error", err)
		return
	}
	if !exists {
		writeRespErr(w, fmt.Sprintf("%s not found for '%s'", resourceType, id), http.StatusNotFound)
		logger.Error("resource not found")
		return
	}

	resource := newResource()
	if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
		writeRespErr(w, fmt.Sprintf("invalid %s data for '%s': %v", resourceType, id, err), http.StatusBadRequest)
		logger.Error("invalid resource data", "error", err)
		return
	}

	if err := writeJSON(filePath, resource, 2); err != nil {
		writeRespErr(w, fmt.Sprintf("failed to update %s '%s'", resourceType, id), http.StatusInternalServerError)
		logger.Error("failed to update resource", "error", err)
		return
	}

	writeRespOk(w, fmt.Sprintf("updated %s '%s'", resourceType, id), resource)
}

// deleteResourceByID handles DELETE request to remove a resource
func deleteResourceByID(
	w http.ResponseWriter,
	r *http.Request,
	storageDir string,
	resourceType resourceType,
) {
	logger := slog.With("url", r.RequestURI, "method", r.Method)
	id := r.PathValue("id")
	if id == "" {
		writeRespErr(w, fmt.Sprintf("%s ID is required", resourceType), http.StatusBadRequest)
		logger.Error("failed to parse ID from URL")
		return
	}

	filePath := filepath.Join(storageDir, id+".json")

	if err := os.Remove(filePath); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, fmt.Sprintf("%s not found for '%s'", resourceType, id), http.StatusNotFound)
			logger.Error("resource item not found")
		} else {
			writeRespErr(w, fmt.Sprintf("failed to delete %s '%s'", resourceType, id), http.StatusInternalServerError)
			logger.Error("failed to delete resource item")
		}
		return
	}

	writeRespOk(w, fmt.Sprintf("deleted %s '%s'", resourceType, id), nil)
}

// getAllResources handles GET request for listing all resources
func getAllResources(
	w http.ResponseWriter,
	r *http.Request,
	storageDir string,
	resourceType resourceType,
	getAll func(string) (any, error),
) {
	logger := slog.With("url", r.RequestURI, "method", r.Method)
	resources, err := getAll(storageDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, fmt.Sprintf("%s resource is currently empty", resourceType), http.StatusNotFound)
			logger.Error("resource is empty")
		} else {
			writeRespErr(w, fmt.Sprintf("failed to list %s informations", resourceType), http.StatusInternalServerError)
			logger.Error("failed to list resource informations")
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
	resourceType resourceType,
	newResource func() ResourceData,
) {
	logger := slog.With("url", r.RequestURI, "method", r.Method)
	if r.Body == nil || r.ContentLength == 0 {
		writeRespErr(w, "request body is empty", http.StatusBadRequest)
		logger.Error("request body is empty")
		return
	}

	resource := newResource()
	if err := json.NewDecoder(r.Body).Decode(resource); err != nil {
		writeRespErr(w, fmt.Sprintf("invalid %s data", resourceType), http.StatusBadRequest)
		logger.Error("invalide resource data")
		return
	}

	if !resource.HasRequiredFields() {
		writeRespErr(w, fmt.Sprintf("incomplete %s data", resourceType), http.StatusBadRequest)
		logger.Error("incomplete resource data")
		return
	}

	// Extract ID from the resource, or generate from name
	var tempData map[string]interface{}
	bodyBytes, _ := json.Marshal(resource)
	json.Unmarshal(bodyBytes, &tempData)

	var id string
	switch resourceType {
	case ProviderType, ClientType:
		existingID, ok := tempData["id"].(string)
		if !ok || existingID == "" {
			// if ID is not provided, generate from name
			name, ok := tempData["name"].(string)
			if !ok || name == "" {
				writeRespErr(w, "name or id field is required", http.StatusBadRequest)
				logger.Error("failed to generate id from name when ID is not avaliable")
				return
			}
			id = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
		} else {
			id = existingID
		}
	case InvoiceType:
		// id generation: INV-YYMMDDXX
		dateStr := types.Today().Format("060102")
		pattern := fmt.Sprintf("INV-%s*.json", dateStr)
		jsonFiles, err := filepath.Glob(filepath.Join(storageDir, pattern))
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			writeRespErr(w, "failed to generate invoice ID", http.StatusInternalServerError)
			logger.Error("failed to generate invoice ID", "error", err)
			return
		}
		suffix := invoice.FindMaxSuffixFromFilename(jsonFiles) + 1
		id = fmt.Sprintf("INV-%s%02d", dateStr, suffix)

		// apply default email template if not set
		inv, ok := resource.(*invoice.Invoice)
		if !ok {
			writeRespErr(w, "invalid invoice data", http.StatusInternalServerError)
			logger.Error("invalid invoice data")
			return
		}
		if inv.EmailTemplateID == "" {
			inv.EmailTemplateID = "default"
		}
	default:
		writeRespErr(w, "invalid resource type, this is likely an internal error", http.StatusInternalServerError)
		logger.Error("unsupported resource type", "resourceType", resourceType)
		return
	}

	resource.SetID(id)

	filePath := filepath.Join(storageDir, id+".json")

	exists, err := isPathExist(filePath)
	if err != nil {
		writeRespErr(w, fmt.Sprintf("failed to check %s '%s' existence: %v", resourceType, id, err), http.StatusInternalServerError)
		logger.Error("failed to check resource existence", "error", err)
		return
	}
	if exists {
		writeRespErr(w, fmt.Sprintf("%s already exists for '%s'", resourceType, id), http.StatusConflict)
		logger.Error("resource already exists")
		return
	}

	if err := writeJSON(filePath, resource, 2); err != nil {
		writeRespErr(w, fmt.Sprintf("failed to create %s '%s'", resourceType, id), http.StatusInternalServerError)
		logger.Error("failed to create resource", "error", err)
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
