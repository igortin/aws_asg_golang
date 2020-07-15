package main

import (
	"errors"
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

	// Main logic
	for _, v := range res {
		fmt.Println("ASG name:", v.AutoScaleGroupName)
		log.Println("Trying to scale out ASG ...")
		isTrue, err := ScaleOutAsg(svc, v)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		if isTrue {
			log.Println("Scale out was succesfully completed")
		} else {
			log.Println("Error: Could not increase ASG, something wrong")
		}
	}
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
			AutoScaleGroupName:   *group.AutoScalingGroupName,
			LaunchTemplate:       *group.LaunchTemplate,
			TargetGroupARNs:      group.TargetGroupARNs,
			MinSize:              *group.MinSize,
			MaxSize:              *group.MaxSize,
			ServiceLinkedRoleARN: *group.ServiceLinkedRoleARN,
			DesiredSize:          *group.DesiredCapacity,
			Instances:            group.Instances,
		}

		m[*group.AutoScalingGroupARN] = &item
	}
	return m, nil
}

// ScaleOutAsg func increase number ASg instances
func ScaleOutAsg(svc *autoscaling.AutoScaling, group *AutoScaleGroup) (bool, error) {
	newDesireValue := group.DesiredSize + 1

	if newDesireValue > group.MaxSize {
		return false, errors.New("Error: ASG can not be scale out -> ASG max size exceeded")
	}

	if newDesireValue < group.MinSize {
		return false, errors.New("Error: ASG can not be scale out -> ASG min size bigger")
	}

	input := &autoscaling.SetDesiredCapacityInput{
		AutoScalingGroupName: aws.String(asgName),
		DesiredCapacity:      aws.Int64(newDesireValue),
		HonorCooldown:        aws.Bool(true),
	}

	_, err := svc.SetDesiredCapacity(input)
	if err != nil {
		return false, err
	}

	log.Printf("ASG name: %s, new desired number set: %d", group.AutoScaleGroupName, newDesireValue)
	return true, nil
}
