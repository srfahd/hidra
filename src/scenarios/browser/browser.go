// Package browser run tests in real browser
package browser

import (
	"context"
	"fmt"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/hidracloud/hidra/src/models"
	"github.com/hidracloud/hidra/src/scenarios"
)

// Scenario Represent an browser scenario
type Scenario struct {
	models.Scenario
	ctx context.Context
}

// RCA generate RCAs for scenario
func (h *Scenario) RCA(result *models.ScenarioResult) error {
	log.Println("Chrome RCA")
	return nil
}

func (b *Scenario) urlShouldBe(c map[string]string) ([]models.Metric, error) {
	if _, ok := c["url"]; !ok {
		return nil, fmt.Errorf("url parameter missing")
	}

	var url string
	err := chromedp.Run(b.ctx,
		chromedp.Location(&url),
	)

	if err != nil {
		return nil, err
	}

	if url != c["url"] {
		return nil, fmt.Errorf("url is not %s is %s", c["url"], url)
	}
	return nil, nil
}

func (b *Scenario) textShouldBe(c map[string]string) ([]models.Metric, error) {
	if _, ok := c["text"]; !ok {
		return nil, fmt.Errorf("text parameter missing")
	}

	if _, ok := c["selector"]; !ok {
		return nil, fmt.Errorf("selector parameter missing")
	}

	var text string
	err := chromedp.Run(b.ctx,
		chromedp.Text(c["selector"], &text),
	)

	if err != nil {
		return nil, err
	}

	if text != c["text"] {
		return nil, fmt.Errorf("text is not %s is %s", c["text"], text)
	}
	return nil, nil
}

func (b *Scenario) sendKeys(c map[string]string) ([]models.Metric, error) {
	if _, ok := c["keys"]; !ok {
		return nil, fmt.Errorf("keys parameter missing")
	}

	if _, ok := c["selector"]; !ok {
		return nil, fmt.Errorf("selector parameter missing")
	}

	err := chromedp.Run(b.ctx,
		chromedp.SendKeys(c["selector"], c["keys"]),
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *Scenario) waitVisible(c map[string]string) ([]models.Metric, error) {
	if _, ok := c["selector"]; !ok {
		return nil, fmt.Errorf("selector parameter missing")
	}

	err := chromedp.Run(b.ctx,
		chromedp.WaitVisible(c["selector"]),
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *Scenario) click(c map[string]string) ([]models.Metric, error) {
	if _, ok := c["selector"]; !ok {
		return nil, fmt.Errorf("selector parameter missing")
	}

	err := chromedp.Run(b.ctx,
		chromedp.Click(c["selector"], chromedp.NodeVisible),
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (b *Scenario) navigateTo(c map[string]string) ([]models.Metric, error) {
	if _, ok := c["url"]; !ok {
		return nil, fmt.Errorf("url parameter missing")
	}

	err := chromedp.Run(b.ctx,
		chromedp.Navigate(c["url"]),
	)

	if err != nil {
		return nil, err
	}

	return nil, nil
}

// Description return scenario description
func (b *Scenario) Description() string {
	return "Run scenario in real browser"
}

// Close closes the scenario
func (s *Scenario) Close() {

}

// Init initialize scenario
func (b *Scenario) Init() {
	b.StartPrimitives()

	// Initialize chrome context
	b.ctx, _ = chromedp.NewContext(
		context.Background(),
	)

	b.RegisterStep("navigateTo", models.StepDefinition{
		Description: "It navigates to a url",
		Params: []models.StepParam{
			{Name: "url", Optional: false, Description: "The url to navigate to"},
		},
		Fn: b.navigateTo,
	})

	b.RegisterStep("urlShouldBe", models.StepDefinition{
		Description: "It checks if the url is the expected one",
		Params: []models.StepParam{
			{Name: "url", Optional: false, Description: "The url to check"},
		},
		Fn: b.urlShouldBe,
	})

	b.RegisterStep("textShouldBe", models.StepDefinition{
		Description: "It checks if the text is the expected one",
		Params: []models.StepParam{
			{Name: "text", Optional: false, Description: "The text to check"},
			{Name: "selector", Optional: false, Description: "The selector to check"},
		},
		Fn: b.textShouldBe,
	})

	b.RegisterStep("sendKeys", models.StepDefinition{
		Description: "It sends keys to a selector",
		Params: []models.StepParam{
			{Name: "keys", Optional: false, Description: "The keys to send to the selector. "},
			{Name: "selector", Optional: false, Description: "The selector to send the keys to. "},
		},
		Fn: b.sendKeys,
	})

	b.RegisterStep("waitVisible", models.StepDefinition{
		Description: "It waits until the selector is visible",
		Params: []models.StepParam{
			{Name: "selector", Optional: false, Description: "The selector to wait for. "},
		},
		Fn: b.waitVisible,
	})

	b.RegisterStep("click", models.StepDefinition{
		Description: "It clicks on a selector",
		Params: []models.StepParam{
			{Name: "selector", Optional: false, Description: "The selector to click. "},
		},
		Fn: b.click,
	})
}

func init() {
	scenarios.Add("browser", func() models.IScenario {
		return &Scenario{}
	})
}
