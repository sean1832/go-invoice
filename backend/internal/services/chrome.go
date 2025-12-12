package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/page"
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

// NewChromeService initializes a new ChromeService with a headless browser context
func NewChromeService() (*ChromeService, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Headless,
		chromedp.DisableGPU,
		// STRICTLY REQUIRED for running in Docker/Linux without a display
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true), // prevents OOM in container /dev/shm

		// reduce memory usage for long-running services
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-background-networking", true),
	)

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

// GeneratePDF navigates to the specified URL and generates a PDF of the page
func (s *ChromeService) GeneratePDF(url string, timeout time.Duration, paperSize PaperSize, title string) ([]byte, error) {
	tabCtx, cancel := chromedp.NewContext(s.browserCtx)
	defer cancel()

	tabCtx, cancelTimemout := context.WithTimeout(tabCtx, timeout)
	defer cancelTimemout()

	var pdfBuffer []byte

	err := chromedp.Run(tabCtx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`#pdf-render-complete, #pdf-render-error`, chromedp.ByQuery), // wait for id to shows up
		chromedp.Evaluate(fmt.Sprintf(`document.title = %q`, title), nil),
		emulation.SetEmulatedMedia().WithMedia("print"), // set emulated media to print for proper paper sizing
		chromedp.ActionFunc(func(ctx context.Context) error {
			// check if error element exists
			var errorNodes []*cdp.Node
			err := chromedp.Nodes(`#pdf-reader-error`, &errorNodes, chromedp.ByQuery, chromedp.AtLeast(0)).Do(ctx)
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
	} else if path, err := exec.LookPath("chrome-browser"); err == nil {
		chromePath = path
	} else if path, err := exec.LookPath("google-chrome"); err == nil {
		chromePath = path
	}
	return chromePath
}
