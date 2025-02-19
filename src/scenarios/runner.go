package scenarios

import (
	"fmt"
	"log"
	"time"

	"github.com/hidracloud/hidra/src/models"
	"github.com/hidracloud/hidra/src/utils"
)

// InitializeScenario initialize a new scenario
func InitializeScenario(s models.Scenario) models.IScenario {
	srunner := Sample[s.Kind]()
	srunner.Init()

	return srunner
}

// RunIScenario run already initialize scenario
func RunIScenario(s models.Scenario, srunner models.IScenario) *models.ScenarioResult {
	metric := models.ScenarioResult{}
	metric.Scenario = s
	metric.StepResults = make([]*models.StepResult, 0)
	metric.StartDate = time.Now()

	defer srunner.Close()

	for _, step := range s.Steps {
		smetric := models.StepResult{}
		smetric.Step = step
		smetric.StartDate = time.Now()
		customMetrics, err := srunner.RunStep(step.Type, step.Params, step.Timeout)

		smetric.Metrics = customMetrics
		smetric.EndDate = time.Now()
		metric.StepResults = append(metric.StepResults, &smetric)

		if step.Negate && err == nil {
			metric.Error = fmt.Errorf("expected fail")
			metric.EndDate = time.Now()
			return &metric
		}

		if err != nil && !step.Negate {
			metric.Error = err
			metric.EndDate = time.Now()
			return &metric
		}
	}

	metric.EndDate = time.Now()
	return &metric
}

// RunScenario Run one scenario
func RunScenario(s models.Scenario, name, desc string) *models.ScenarioResult {
	utils.LogDebug("[%s] Running new scenario, \"%s\"\n", name, desc)
	srunner := InitializeScenario(s)
	return RunIScenario(s, srunner)
}

// PrettyPrintScenarioResults Print scenario metrics
func PrettyPrintScenarioResults(m *models.ScenarioResult, name, desc string) {
	log.Printf("[%s] Metrics results for: %s\n", name, desc)
	if m.Error == nil {
		log.Printf("[%s] Scenario ran without issues\n", name)
	} else {
		log.Printf("[%s] Scenario ran with issues: %s\n", name, m.Error)

	}
	log.Printf("[%s] Total scenario duration: %d (ms)\n", name, m.EndDate.Sub(m.StartDate).Milliseconds())

	for _, s := range m.StepResults {
		log.Printf("[%s]   |_ %s duration: %d (ms)\n", name, s.Step.Type, s.EndDate.Sub(s.StartDate).Milliseconds())
	}
}

// PrettyPrintScenarioResults2String Print scenario metrics
func PrettyPrintScenarioResults2String(m *models.ScenarioResult, name, desc string) string {
	out := ""

	out += fmt.Sprintf("[%s] Metrics results for: %s\n", name, desc)
	if m.Error == nil {
		out += fmt.Sprintf("[%s] Scenario ran without issues\n", name)
	} else {
		out += fmt.Sprintf("[%s] Scenario ran with issues: %s\n", name, m.Error)

	}

	out += fmt.Sprintf("[%s] Total scenario duration: %d (ms)\n", name, m.EndDate.Sub(m.StartDate).Milliseconds())

	for _, s := range m.StepResults {
		out += fmt.Sprintf("[%s]   |_ %s duration: %d (ms)\n", name, s.Step.Type, s.EndDate.Sub(s.StartDate).Milliseconds())
	}

	out += fmt.Sprintf("[%s] Total scenario duration: %d (ms)\n", name, m.EndDate.Sub(m.StartDate).Milliseconds())
	return out
}
