package detector

import (
	"github.com/hashicorp/go-version"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/schema"
)

type AwsInstanceNotSpecifiedTagDetector struct {
	*Detector
}

func (d *Detector) CreateAwsInstanceNotSpecifiedTagDetector() *AwsInstanceNotSpecifiedTagDetector {
	nd := &AwsInstanceNotSpecifiedTagDetector{Detector: d}
	nd.Name = "aws_instance_not_specified_tag"
	nd.IssueType = issue.NOTICE
	nd.TargetType = "resource"
	nd.Target = "aws_instance"
	nd.DeepCheck = false
	nd.Link = "https://github.com/wata727/tflint/blob/master/docs/aws_instance_not_specified_tag.md"
	return nd
}

func (d *AwsInstanceNotSpecifiedTagDetector) Detect(resource *schema.Resource, issues *[]*issue.Issue) {
	requiredTags := []string{
		"Project",
		"Environment",

	}

	optionalTags :=[]string{
		"Application",
		"Role",
	}

	v1, err := version.NewVersion(d.Config.TerraformVersion)
	// If `terraform_version` is not set, always detect.
	if err != nil {
		v1, _ = version.NewVersion("0.8.0")
	}
	v2, _ := version.NewVersion("0.8.8")

	tags, _ := resource.GetMapToken("tags")

	for _, tag := range requiredTags {
		if _, ok := tags[tag]; !ok && v1.LessThan(v2) {
			issue := &issue.Issue{
				Detector: d.Name,
				Type:     issue.ERROR,
				Message:  "Tag " + tag + " is not specified. This is required, please add it to instance",
				Line:     resource.Pos.Line,
				File:     resource.Pos.Filename,
				Link:     d.Link,
			}
			*issues = append(*issues, issue)
		}
	}

	for _, tag := range optionalTags {
		if _, ok := tags[tag]; !ok && v1.LessThan(v2) {
			issue := &issue.Issue{
				Detector: d.Name,
				Type:     d.IssueType,
				Message:  "Tag " + tag + " is not specified. This is optional, please verify your tags.",
				Line:     resource.Pos.Line,
				File:     resource.Pos.Filename,
				Link:     d.Link,
			}
			*issues = append(*issues, issue)
		}
	}

}
