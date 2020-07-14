package main

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

var (
	profileaName              = "private"
	awsRegion                 = "us-east-1"
	tagKey                    = "tag:role"
	tagValue                  = "etcd"
	instanceStateFilter       = "stopped"
	asgName                   = "GOOD-ASG"
	limit               int64 = 10
)

func main() {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(awsRegion)},
		Profile: profileaName,
	})
	if err != nil {
		log.Println("Error: Can not create new session.")
	}

	svc := autoscaling.New(sess)

	data, err := getAsgData(svc)
	if err != nil {
		log.Printf("Error: %v", err)
	}

	res, err := parseAsgData(data)
	if err != nil {
		log.Printf("Error: %v", err)
	}
	fmt.Printf("ASG: %v", res)
}

func getAsgData(svc *autoscaling.AutoScaling) (*autoscaling.DescribeAutoScalingGroupsOutput, error) {

	input := &autoscaling.DescribeAutoScalingGroupsInput{
		MaxRecords: aws.Int64(limit),
		AutoScalingGroupNames: []*string{
			aws.String(asgName),
		},
	}

	output, err := svc.DescribeAutoScalingGroups(input)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func parseAsgData(data *autoscaling.DescribeAutoScalingGroupsOutput) (map[string]*AutoScaleGroup, error) {

	m := map[string]*AutoScaleGroup{}

	// iterate []*Group
	for _, group := range data.AutoScalingGroups {

		// init instance of structure AutoScaleGroup{}
		item := AutoScaleGroup{
			AutoScaleGroupName: *group.AutoScalingGroupName,
			LaunchTemplate:     *group.LaunchTemplate,
			TargetGroupARNs:    group.TargetGroupARNs,
		}

		m[*group.AutoScalingGroupARN] = &item
	}
	return m, nil
}
