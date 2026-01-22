package browser

import (
	"context"

	"github.com/chromedp/chromedp"
)

// SetupContext creates a chromedp context with the specified options.
// Returns the context and a cancel function that should be called when done.
func SetupContext(ctx context.Context, headless bool, userAgent string) (context.Context, context.CancelFunc) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.UserAgent(userAgent),
		chromedp.Flag("headless", headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	allocCtx, cancel1 := chromedp.NewExecAllocator(ctx, opts...)
	browserCtx, cancel2 := chromedp.NewContext(allocCtx)

	// Return combined cancel function
	cancelFunc := func() {
		cancel2()
		cancel1()
	}

	return browserCtx, cancelFunc
}
