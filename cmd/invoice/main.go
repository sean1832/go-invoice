package main

import (
	"embed"
	"fmt"
	"html/template"
	"invoice/internal/core/model/storage"
	"invoice/internal/core/model/view"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// TODO: CLI
// invoice init
// invoice write lanzhang dingyuxu --date 2025/05/12
// invoice write lanzhang dingyuxu --date 2025/05/15 --hour 3 --rate 50 --description "house clean"
// invoice export --output invoice.pdf
// invoice send --title "invoice from lan zhang" --message "invoice attached. Kind regards, Lan Zhang" --email account@gmail.com

// invoice clear
// invoice provider ls
// invoice provider new "lanzhang"
// invoice provider rm "lanzhang"
// invoice client ls 							// list all clients
// invoice client new "dingyuxu" 				// add client profile (prompt info fields)
// invoice client rm "dingyuxu" 				// remove client profile
// invoice ls 									// list all invoices in the registry
// invoice restore -i 20251020_dingyuxu.json	// restore json file
// invoice restore 20251020_dingyuxu			// restore from registry (json file)

// `init` command initialize the program. create `db/` dir, prompt user to enter email address & app password.
// `provider` command manages provider informations
// `client` command manages client informations
// `write` command writes JSON in staging folder. everytime user write, add a row to entry
// `export` command render html to PDF and export
// `send` command send latest the PDF to client, save json to registry for archive
// `restore` command restore a cached invoice JSON in staging folder
// `clear` command clear staging folder

// row removal command needed?
// `send` command might be too fragil, can accedentaly send the wrong PDF to client. perhaps we group with export and send to prevent this.

//go:embed assets
var assets embed.FS

func main() {
	staticFS, err := fs.Sub(assets, "assets")
	if err != nil {
		fmt.Printf("error extracting embedded assets: %v", err)
		return
	}
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(staticFS))))

	data := setData()
	if data == nil {
		return
	}
	// parse embedded template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		funcMap := template.FuncMap{
			"formatDate": func(t time.Time) string {
				return t.Format("2006/01/02")
			},
		}
		tmpl := template.Must(template.New("invoice.html").Funcs(funcMap).ParseFS(assets, "assets/templates/invoice.html"))
		tmpl.ExecuteTemplate(w, "invoice.html", data)
	})

	fmt.Println("Listening on http://localhost:8080 ...")
	http.ListenAndServe(":8080", nil)
}

func setData() *view.ViewModel {
	const layout = "2006/01/02"
	date1, err := time.Parse(layout, "2025/10/30")
	if err != nil {
		fmt.Println("error parsing datetime")
		return nil
	}
	date2, err := time.Parse(layout, "2025/10/15")
	if err != nil {
		fmt.Println("error parsing datetime")
		return nil
	}
	exeDir, err := getExeDir()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	const clientPath = "db/clients/dingyuxu.json"
	const providerPath = "db/providers/lanzhang.json"
	if _, err := os.Stat(filepath.Join(*exeDir, clientPath)); err != nil {
		fmt.Printf("failed to open file or filepath not exist '%s': %v", clientPath, err)
		return nil
	}
	if _, err := os.Stat(filepath.Join(*exeDir, providerPath)); err != nil {
		fmt.Printf("failed to open file or filepath not exist '%s': %v", providerPath, err)
		return nil
	}

	clientJson, err := os.ReadFile(filepath.Join(*exeDir, clientPath))
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return nil
	}
	providerJson, err := os.ReadFile(filepath.Join(*exeDir, providerPath))
	if err != nil {
		fmt.Printf("failed to read file: %v\n", err)
		return nil
	}

	client, err := storage.NewClientFromJSON(clientJson)
	if err != nil {
		fmt.Printf("failed to create client from JSON: %v", err)
	}
	provider, err := storage.NewProviderFromJSON(providerJson)
	if err != nil {
		fmt.Printf("failed to create provider from JSON: %v", err)
	}

	model := storage.NewModel(*provider, *client).ToViewModel()
	model.AddItem(view.NewServiceEntry(date1, "Website Design", 1, 1500.00))
	model.AddItem(view.NewServiceEntryWithDetail(date2, "Algorithm Development", "Development of custom algorithms for data processing", 2, 2500.00))

	return model
}

func getExeDir() (*string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	exeDir := filepath.Dir(exePath)
	return &exeDir, nil
}
