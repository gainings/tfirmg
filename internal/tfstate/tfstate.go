package tfstate

import (
	"context"
	"encoding/json"
	"github.com/gainings/tfirg/internal/model/resource"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/pkg/errors"
)

type TFState struct {
	Version          int    `json:"version"`
	TerraformVersion string `json:"terraform_version"`
	Resources        []struct {
		Module    *string `json:"module"`
		Mode      string  `json:"mode"`
		Type      string  `json:"type"`
		Name      string  `json:"name"`
		Provider  string  `json:"provider"`
		Instances []struct {
			IndexKey   json.RawMessage        `json:"index_key"`
			Attributes map[string]interface{} `json:"attributes"`
		} `json:"instances"`
	} `json:"resources"`
}

func LoadTFState(ctx context.Context, loc string) (*TFState, error) {
	u, err := url.Parse(loc)
	if err != nil {
		return nil, err
	}

	var rc io.ReadCloser
	//TODO: support gcs, azure and more...
	switch u.Scheme {
	case "s3":
		key := strings.TrimPrefix(u.Path, "/")
		rc, err = readFromS3(ctx, u.Host, key)
	case "file":
		rc, err = os.Open(u.Path)
	default:
		err = errors.Errorf("URL scheme %s is not supported", u.Scheme)
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read tfstate from %s", u.String())
	}
	defer rc.Close()

	var tfState *TFState
	decoder := json.NewDecoder(rc)
	if err := decoder.Decode(&tfState); err != nil {
		return nil, err
	}
	return tfState, nil
}
func readFromS3(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	return result.Body, nil
}

func getTFStateLocation(filename string) (string, error) {
	//TODO: hcl parse and get tfstate path
	return "", nil
}

type Transformer struct {
	rc resource.ResourceCreator
}

func NewTransformer(rc resource.ResourceCreator) Transformer {
	return Transformer{
		rc: rc,
	}
}

func (t Transformer) TransformToResources(tfstate *TFState) resource.Resources {
	var rs resource.Resources
	for _, r := range tfstate.Resources {
		if r.Mode == "data" {
			continue
		}
		for _, i := range r.Instances {
			r := t.rc.Create(r.Type, r.Name, string(i.IndexKey), r.Module, i.Attributes)
			if r.ID == nil {
				continue
			}
			rs = append(rs, r)
		}
	}
	return rs
}
