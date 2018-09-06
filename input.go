package s3

import (
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type inputModel struct {
	client *s3.S3

	method string

	credentials external.WithCredentialsValue
	region      external.WithRegion
	endpoint    string

	bucket *string
	key    *string

	expireTime time.Duration
	etag       *string

	content              io.ReadSeeker
	contentType          *string
	contentEncoding      *string
	contentMD5           *string
	contentLength        *int64
	serverSideEncryption s3.ServerSideEncryption
}

func newInput(method string, options []Option) *inputModel {
	in := &inputModel{method: method}
	in.applyOptions(options)
	return in
}

func (this *inputModel) applyOptions(options []Option) {
	for _, option := range options {
		if option != nil {
			option(this)
		}
	}
}

func (this *inputModel) validate() error {
	if this.method != GET && this.method != PUT {
		return ErrInvalidRequestMethod
	}
	if this.bucket == nil || len(*this.bucket) == 0 {
		return ErrBucketMissing
	}
	if this.key == nil || len(*this.key) == 0 {
		return ErrKeyMissing
	}
	if this.method == PUT && this.content == nil {
		return ErrContentMissing
	}
	return nil
}

func (this *inputModel) buildClient() error {
	config, err := this.buildConfig()
	if err != nil {
		return err
	}

	this.client = s3.New(config)
	if len(this.endpoint) > 0 {
		this.client.ForcePathStyle = true
	}

	return nil
}

func (this *inputModel) buildConfig() (aws.Config, error) {
	var configs []external.Config
	if this.credentials.AccessKeyID != "" {
		configs = append(configs, this.credentials)
	}
	if len(this.region) > 0 {
		configs = append(configs, external.WithRegion(this.region))
	}

	config, err := external.LoadDefaultAWSConfig(configs...)
	if err != nil {
		return aws.Config{}, err
	}

	if len(this.endpoint) > 0 {
		config.EndpointResolver = aws.ResolveWithEndpointURL(this.endpoint)
	}

	return config, err
}

func (this *inputModel) buildAWSRequest() (request *aws.Request) {
	if this.method == GET {
		request = this.buildGET()
	} else {
		request = this.buildPUT()
	}
	request.ExpireTime = this.expireTime
	return request
}

func (this *inputModel) buildGET() *aws.Request {
	parameters := s3.GetObjectInput{
		Bucket:      this.bucket,
		Key:         this.key,
		IfNoneMatch: this.etag,
	}
	request := this.client.GetObjectRequest(&parameters)
	request.ExpireTime = this.expireTime
	return request.Request
}

func (this *inputModel) buildPUT() *aws.Request {
	parameters := s3.PutObjectInput{
		Bucket: this.bucket,
		Key:    this.key,

		Body:            this.content,
		ContentMD5:      this.contentMD5,
		ContentType:     this.contentType,
		ContentLength:   this.contentLength,
		ContentEncoding: this.contentEncoding,

		ServerSideEncryption: this.serverSideEncryption,
	}
	request := this.client.PutObjectRequest(&parameters)
	request.ExpireTime = this.expireTime
	return request.Request
}
