// Code generated by internal/generate/tags/main.go; DO NOT EDIT.
package cloudfront

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	awstypes "github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/logging"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/types/option"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// listTags lists cloudfront service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func listTags(ctx context.Context, conn cloudfront.Client, identifier string) (tftags.KeyValueTags, error) {
	input := &cloudfront.ListTagsForResourceInput{
		Resource: aws.String(identifier),
	}

	output, err := conn.ListTagsForResource(ctx, input)

	if err != nil {
		return tftags.New(ctx, nil), err
	}

	return KeyValueTags(ctx, output.Tags.Items), nil
}

// ListTags lists cloudfront service tags and set them in Context.
// It is called from outside this package.
func (p *servicePackage) ListTags(ctx context.Context, meta any, identifier string) error {
	tags, err := listTags(ctx, *meta.(*conns.AWSClient).CloudFrontClient(ctx), identifier)

	if err != nil {
		return err
	}

	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(tags)
	}

	return nil
}

// []*SERVICE.Tag handling

// Tags returns cloudfront service tags.
func Tags(tags tftags.KeyValueTags) *awstypes.Tags {
	result := make([]awstypes.Tag, 0, len(tags))

	for k, v := range tags.Map() {
		tag := awstypes.Tag{
			Key:   aws.String(k),
			Value: aws.String(v),
		}

		result = append(result, tag)
	}

	return &awstypes.Tags{Items: result}
}

// KeyValueTags creates tftags.KeyValueTags from cloudfront service tags.
func KeyValueTags(ctx context.Context, tags []awstypes.Tag) tftags.KeyValueTags {
	m := make(map[string]*string, len(tags))

	for _, tag := range tags {
		m[aws.ToString(tag.Key)] = tag.Value
	}

	return tftags.New(ctx, m)
}

// getTagsIn returns cloudfront service tags from Context.
// nil is returned if there are no input tags.
func getTagsIn(ctx context.Context) []awstypes.Tag {
	tagsResult := make([]awstypes.Tag, 0)
	if inContext, ok := tftags.FromContext(ctx); ok {
		if tags := Tags(inContext.TagsIn.UnwrapOrDefault()); len(tags.Items) > 0 {
			for _, t := range tags.Items {
				tagsResult = append(tagsResult, t)
			}
		}
	}

	return nil
}

// setTagsOut sets cloudfront service tags in Context.
func setTagsOut(ctx context.Context, tags []awstypes.Tag) {
	if inContext, ok := tftags.FromContext(ctx); ok {
		inContext.TagsOut = option.Some(KeyValueTags(ctx, tags))
	}
}

// updateTags updates cloudfront service tags.
// The identifier is typically the Amazon Resource Name (ARN), although
// it may also be a different identifier depending on the service.
func updateTags(ctx context.Context, conn *cloudfront.Client, identifier string, oldTagsMap, newTagsMap any) error {

	oldTags := tftags.New(ctx, oldTagsMap)
	newTags := tftags.New(ctx, newTagsMap)

	ctx = tflog.SetField(ctx, logging.KeyResourceId, identifier)

	removedTags := oldTags.Removed(newTags)
	removedTags = removedTags.IgnoreSystem(names.CloudFront)
	removedTagKeys := removedTags.Keys()
	if len(removedTags) > 0 {
		input := &cloudfront.UntagResourceInput{
			Resource: aws.String(identifier),
			TagKeys:  &awstypes.TagKeys{Items: removedTagKeys},
		}

		_, err := conn.UntagResource(ctx, input)

		if err != nil {
			return fmt.Errorf("untagging resource (%s): %w", identifier, err)
		}
	}

	updatedTags := oldTags.Updated(newTags)
	updatedTags = updatedTags.IgnoreSystem(names.CloudFront)
	if len(updatedTags) > 0 {
		input := &cloudfront.TagResourceInput{
			Resource: aws.String(identifier),
			Tags:     Tags(updatedTags),
		}

		_, err := conn.TagResource(ctx, input)

		if err != nil {
			return fmt.Errorf("tagging resource (%s): %w", identifier, err)
		}
	}

	return nil
}

// UpdateTags updates cloudfront service tags.
// It is called from outside this package.
func (p *servicePackage) UpdateTags(ctx context.Context, meta any, identifier string, oldTags, newTags any) error {
	return updateTags(ctx, meta.(*conns.AWSClient).CloudFrontClient(ctx), identifier, oldTags, newTags)
}
