module github.com/kiwibel/es-monitor

go 1.14

replace github.com/kiwibel/es-monitor/slack => ./pkg/slack

replace github.com/kiwibel/es-monitor/metric => ./pkg/metric

require (
	github.com/aws/aws-sdk-go v1.35.15 // indirect
	github.com/kiwibel/es-monitor/metric v0.0.0-00010101000000-000000000000
	github.com/kiwibel/es-monitor/slack v0.0.1
	github.com/slack-go/slack v0.7.2 // indirect
)
