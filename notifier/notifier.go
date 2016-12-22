// Package notifier manages notifications for consul-alerts
package notifier

import "time"

const (
	SYSTEM_HEALTHY  string = "HEALTHY"
	SYSTEM_UNSTABLE string = "UNSTABLE"
	SYSTEM_CRITICAL string = "CRITICAL"
)

const header = `%s is %s.

Fail: %d, Warn: %d, Pass: %d
`

type Message struct {
	Node      string
	ServiceId string
	Service   string
	CheckId   string
	Check     string
	Status    string
	Output    string
	Notes     string
	Interval  int
	RmdCheck  time.Time
	NotifList map[string]bool
	Timestamp time.Time
}

type Messages []Message

type Notifier interface {
	Notify(alerts Messages) bool
	NotifierName() string
}

type Notifiers struct {
	Email     *EmailNotifier
	Log       *LogNotifier
	Influxdb  *InfluxdbNotifier
	Slack     *SlackNotifier
	PagerDuty *PagerDutyNotifier
	HipChat   *HipChatNotifier
	OpsGenie  *OpsGenieNotifier
	AwsSns    *AwsSnsNotifier
	VictorOps *VictorOpsNotifier
	Custom    []string
}

func (m Message) IsCritical() bool {
	return m.Status == "critical"
}

func (m Message) IsWarning() bool {
	return m.Status == "warning"
}

func (m Message) IsPassing() bool {
	return m.Status == "passing"
}

func (m Messages) Summary() (overallStatus string, pass, warn, fail int) {
	hasCritical := false
	hasWarnings := false
	for _, message := range m {
		switch {
		case message.IsCritical():
			hasCritical = true
			fail++
		case message.IsWarning():
			hasWarnings = true
			warn++
		case message.IsPassing():
			pass++
		}
	}
	if hasCritical {
		overallStatus = SYSTEM_CRITICAL
	} else if hasWarnings {
		overallStatus = SYSTEM_UNSTABLE
	} else {
		overallStatus = SYSTEM_HEALTHY
	}
	return
}
