package main

import (
	"github.com/aws/aws-sdk-go/service/autoscaling"
)

// AutoScaleGroup structure
type AutoScaleGroup struct {
	AutoScalingGroupARN  string                                  `json:"AutoScalingGroupARN"`
	AutoScaleGroupName   string                                  `json:"AutoScaleGroupName"`
	LaunchTemplate       autoscaling.LaunchTemplateSpecification `json:"LaunchTemplate"`
	DesiredSize          int64                                   `json:"DesiredSize"`
	MaxSize              int64                                   `json:"MaxSize"`
	MinSize              int64                                   `json:"MinSize"`
	Status               string                                  `json:"Status"`
	TargetGroupARNs      []*string                               `json:"TargetGroupARNs"`
	ServiceLinkedRoleARN string                                  `json:"ServiceLinkedRoleARN"`
	Instances            []*autoscaling.Instance                 `json:"Instances"`
}
