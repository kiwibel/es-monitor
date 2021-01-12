package metric

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func PutMetric(cluster string, statusCode float64) error {

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create new cloudwatch client.
	svc := cloudwatch.New(sess)
	fmt.Println("Publishing metric data points to Amazon CloudWatch...")
	_, err := svc.PutMetricData(&cloudwatch.PutMetricDataInput{
		Namespace: aws.String("Elasticsearch/Cluster"),
		MetricData: []*cloudwatch.MetricDatum{
			&cloudwatch.MetricDatum{
				MetricName: aws.String("ClusterStatus"),
				Unit:       aws.String("Count"),
				Value:      aws.Float64(statusCode),
				Dimensions: []*cloudwatch.Dimension{
					&cloudwatch.Dimension{
						Name:  aws.String("ESCluster"),
						Value: aws.String(cluster),
					},
				},
			},
		},
	})

	if err != nil {
		fmt.Println("Error adding metrics:", err.Error())
	}
	return err
	// Get information about metrics
	//result, err := svc.ListMetrics(&cloudwatch.ListMetricsInput{
	//	Namespace: aws.String("Elasticsearch/Cluster"),
	//})
	//if err != nil {
	//	fmt.Println("Error getting metrics:", err.Error())
	//	return
	//}
	//
	//for _, metric := range result.Metrics {
	//	fmt.Println(*metric.MetricName)
	//
	//	for _, dim := range metric.Dimensions {
	//		fmt.Println(*dim.Name+":", *dim.Value)
	//		fmt.Println()
	//	}
	//}
}
