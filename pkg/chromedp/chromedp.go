package chromedp

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/chromedp"
)

// ScrapeArticleText scrapes all text content from article tags on a given page
func ScrapeArticleText(pageURL string) (string, error) {
	// Create a new context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// Create a new Chrome instance
	ctx, cancel = chromedp.NewContext(ctx)
	defer cancel()

	var article string

	// Run chromedp tasks
	err := chromedp.Run(ctx,
		// Navigate to the URL
		chromedp.Navigate(pageURL),
		// Wait for the article tags to load
		chromedp.WaitVisible("article", chromedp.ByQuery),
		// Extract text from all article tags
		chromedp.Evaluate(`document.querySelector('article').textContent.replaceAll('\t',"").replaceAll('\n','')`, &article),
	)

	if err != nil {
		return "", fmt.Errorf("failed to scrape articles from %s: %w", pageURL, err)
	}

	return article, nil
}
