package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
	cdpruntime "github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

type PaperSize struct {
	Width  float64
	Height float64
}

var (
	PaperSizeA4     = PaperSize{Width: 8.27, Height: 11.69}  // in inches
	PaperSizeA3     = PaperSize{Width: 11.69, Height: 16.54} // in inches
	PaperSizeLetter = PaperSize{Width: 8.5, Height: 11.0}    // in inches
)

// ChromeService manages a headless Chrome browser instance
type ChromeService struct {
	browserCtx context.Context
	cancel     context.CancelFunc
}

// NewChromeService initializes a ChromeService instance. Uses remote Chrome if CHROME_REMOTE_URL is set.
func NewChromeService() (*ChromeService, error) {
	remoteURL := os.Getenv("CHROME_REMOTE_URL")
	if remoteURL != "" {
		log.Printf("ChromeService: Mode=Remote, URL=%s", remoteURL)
		return NewRemoteChromeService(remoteURL)
	}
	log.Printf("ChromeService: Mode=Local")
	return NewLocalChromeService()
}

// NewLocalChromeService initializes a local ChromeService instance
func NewLocalChromeService() (*ChromeService, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-background-networking", true),
		chromedp.Flag("mute-audio", true),
		chromedp.Flag("ignore-certificate-errors", true),
	)
	// os specific options
	if runtime.GOOS == "linux" {
		// linux
		log.Printf("ChromeService: Deteced Linux environment. Applying stability patches.")
		opts = append(opts,
			chromedp.NoSandbox,
			chromedp.Flag("disable-setuid-sandbox", true),
			chromedp.Flag("disable-dev-shm-usage", true),
			chromedp.Flag("disable-software-rasterizer", true),

			// These "Nuclear" flags are great for Alpine but CRASH Windows
			chromedp.Flag("no-zygote", true),
			chromedp.Flag("single-process", true),

			// Fix DBus & Temp paths on Linux containers
			chromedp.Env(
				"XDG_CONFIG_HOME=/tmp",
				"XDG_CACHE_HOME=/tmp",
				"HOME=/tmp",
				"DBUS_SESSION_BUS_ADDRESS=/dev/null",
			),
		)
	} else {
		// Windows needs very little configuration.
		log.Println("ChromeService: Detected Windows/Mac environment. Using standard configuration.")
	}

	// detect custom chrome path
	chromePath := getChromePath()
	if chromePath != "" {
		opts = append(opts, chromedp.ExecPath(chromePath))
		log.Printf("using chrome executable at: %s", chromePath)
	} else {
		log.Printf("using default chrome executable")
	}

	// initialize context
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("could not launch browser: %v", err)
	}

	return &ChromeService{
		browserCtx: ctx,
		cancel:     cancel,
	}, nil
}

// Close shuts down the browser context
func (s *ChromeService) Close() {
	s.cancel()
}

// NewRemoteChromeService creates a ChromeService that connects to a remote Chrome instance
func NewRemoteChromeService(remoteURL string) (*ChromeService, error) {
	wsURL, err := getWebSocketURL(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get remote chrome websocket: %v", err)
	}

	// create a remote allocator
	allocCtx, _ := chromedp.NewRemoteAllocator(context.Background(), wsURL)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	// initialize the browser
	if err := chromedp.Run(ctx); err != nil {
		cancel()
		return nil, fmt.Errorf("could not launch/connect to browser: %v", err)
	}

	return &ChromeService{
		browserCtx: ctx,
		cancel:     cancel,
	}, nil
}

func getWebSocketURL(remoteURL string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("%s/json/version", remoteURL))
	if err != nil {
		return "", fmt.Errorf("failed to get remote chrome version info: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	wsURL, ok := result["webSocketDebuggerUrl"].(string)
	if !ok {
		return "", fmt.Errorf("json/version response missing 'webSocketDebuggerUrl'")
	}
	return wsURL, nil
}

// GeneratePDF navigates to the specified URL and generates a PDF of the page
func (s *ChromeService) GeneratePDF(url string, timeout time.Duration, paperSize PaperSize, title string) ([]byte, error) {
	tabCtx, cancel := chromedp.NewContext(s.browserCtx)
	defer cancel()

	tabCtx, cancelTimemout := context.WithTimeout(tabCtx, timeout)
	defer cancelTimemout()

	chromedp.ListenTarget(tabCtx, func(ev interface{}) {
		if ev, ok := ev.(*cdpruntime.EventExceptionThrown); ok {
			log.Printf("[CHROME EXCEPTION] %s", ev.ExceptionDetails.Text)
		}
		if ev, ok := ev.(*cdpruntime.EventConsoleAPICalled); ok {
			for _, arg := range ev.Args {
				log.Printf("[CHROME CONSOLE] %s: %s", ev.Type, string(arg.Value))
			}
		}
	})

	var pdfBuffer []byte

	err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#pdf-render-complete, #pdf-render-error`, chromedp.ByQuery), // wait for id to shows up
		chromedp.Evaluate(fmt.Sprintf(`document.title = %q`, title), nil),
		emulation.SetEmulatedMedia().WithMedia("print"), // set emulated media to print for proper paper sizing
		chromedp.ActionFunc(func(ctx context.Context) error {
			// check if error element exists
			var errorNodes []*cdp.Node
			err := chromedp.Nodes(`#pdf-render-error`, &errorNodes, chromedp.ByQuery, chromedp.AtLeast(0)).Do(ctx)
			if err != nil {
				// error querying the DOM, not an application error
				return fmt.Errorf("could not check for error node: %v", err)
			}

			if len(errorNodes) > 0 {
				var errorText string
				_ = chromedp.TextContent(`#pdf-render-error`, &errorText, chromedp.ByQuery, chromedp.AtLeast(0)).Do(ctx)
				if errorText != "" {
					return fmt.Errorf("page rendered and error: %s", errorText)
				}
				return fmt.Errorf("page rendered and error") // e.g. invoice not found
			}

			// no error node, so #pdf-render-complete must be present.
			// pdf generation
			var printErr error
			pdfBuffer, _, printErr = page.PrintToPDF().
				WithPrintBackground(true).
				WithPaperHeight(paperSize.Height).
				WithPaperWidth(paperSize.Width).
				Do(ctx)
			return printErr
		}),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF for %s: %v", url, err)
	}

	return pdfBuffer, nil
}

func getChromePath() string {
	var chromePath string
	if v := os.Getenv("CHROME_BIN"); v != "" {
		chromePath = v
	} else if path, err := exec.LookPath("chromium-browser"); err == nil {
		// ^ FIXED: Alpine/Ubuntu uses "chromium-browser", not "chrome-browser"
		chromePath = path
	} else if path, err := exec.LookPath("google-chrome"); err == nil {
		chromePath = path
	} else if path, err := exec.LookPath("chromium"); err == nil {
		// Fallback for some other distros
		chromePath = path
	}
	return chromePath
}
