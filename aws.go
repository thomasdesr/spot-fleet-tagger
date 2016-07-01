package main

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
)

func paginateSpotFleetRequest(c *ec2.EC2, r *ec2.DescribeSpotFleetInstancesInput) (*ec2.DescribeSpotFleetInstancesOutput, error) {
	result := &ec2.DescribeSpotFleetInstancesOutput{
		SpotFleetRequestId: r.SpotFleetRequestId,
		NextToken:          aws.String("do while"),
	}

	for result.NextToken != nil {
		t := time.Now()
		resp, err := c.DescribeSpotFleetInstances(r)
		metrics.awsApiLatency.WithLabelValues("DescribeSpotFleetInstances").Observe(time.Since(t).Seconds())
		if err != nil {
			return result, errors.Wrap(err, "failed to DescribeSpotFleetInstances")
		}

		result.ActiveInstances = append(result.ActiveInstances, resp.ActiveInstances...)
		result.NextToken = resp.NextToken
	}

	return result, nil
}

func getInstances(res *ec2.DescribeSpotFleetInstancesOutput) []*string {
	instanceIds := make([]*string, 0, len(res.ActiveInstances))

	for _, i := range res.ActiveInstances {
		instanceIds = append(instanceIds, i.InstanceId)
	}

	return instanceIds
}

func mapToAWSTags(m map[string]string) []*ec2.Tag {
	tags := make([]*ec2.Tag, 0)

	for k, v := range m {
		tags = append(tags, &ec2.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		})
	}

	return tags
}

func tagSpotFleetRequestIds(ec2Client *ec2.EC2, sfrs []string, tags map[string]string) error {
	for _, sfr := range sfrs {
		sfr, err := paginateSpotFleetRequest(ec2Client, &ec2.DescribeSpotFleetInstancesInput{
			SpotFleetRequestId: aws.String(sfr),
		})
		if err != nil {
			return errors.Wrap(err, "failed on paginateSpotFleetRequest")
		}

		t := time.Now()
		_, err = ec2Client.CreateTags(&ec2.CreateTagsInput{
			Resources: getInstances(sfr),
			Tags:      mapToAWSTags(tags),
		})
		metrics.awsApiLatency.WithLabelValues("CreateTags").Observe(time.Since(t).Seconds())
		if err != nil {
			return errors.Wrap(err, "failed to apply tags to instances")
		}
	}
	return nil
}
