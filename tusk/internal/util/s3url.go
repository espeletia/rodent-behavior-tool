package util

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// DefaultRegion contains a default region for an S3 bucket, when a region
// cannot be determined, for example when the s3:// schema is used or when
// path style URL has been given without the region component in the
// fully-qualified domain name.
const DefaultRegion = "us-east-1"

var (
	ErrBucketNotFound    = errors.New("bucket name could not be found")
	ErrHostnameNotFound  = errors.New("hostname could not be found")
	ErrInvalidS3Endpoint = errors.New("an invalid S3 endpoint URL")

	// Pattern used to parse multiple path and host style S3 endpoint URLs.
	s3URLPattern          = regexp.MustCompile(`^(.+\.)?s3[.-](?:(accelerated|dualstack|website)[.-])?([a-z0-9-]+)\.`)
	minioURLPattern       = regexp.MustCompile(`^(.+\.)?minio[.-]?(?:(accelerated|dualstack|website)[.-])?([a-z0-9-]+)?\.?`)
	digitalOceanURLPatter = regexp.MustCompile(`^(.+\.)?digitaloceanspaces[.-]?(?:(accelerated|dualstack|website)[.-])?([a-z0-9-]+)?\.?`)
	localhostURLPattern   = regexp.MustCompile(`^(.+\.)?localhost[.-]?(?:(accelerated|dualstack|website)[.-])?([a-z0-9-]+)?\.?`)
	ipURLPattern          = regexp.MustCompile(`^(.+\.)?(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})[.-]?(?:(accelerated|dualstack|website)[.-])?([a-z0-9-]+)?\.?`)
)

type S3URIOpt func(*S3URI)

func WithScheme(s string) S3URIOpt {
	return func(s3u *S3URI) {
		s3u.Scheme = String(s)
	}
}

func WithBucket(s string) S3URIOpt {
	return func(s3u *S3URI) {
		s3u.Bucket = String(s)
	}
}

func WithKey(s string) S3URIOpt {
	return func(s3u *S3URI) {
		s3u.Key = String(s)
	}
}

func WithVersionID(s string) S3URIOpt {
	return func(s3u *S3URI) {
		s3u.VersionID = String(s)
	}
}

func WithRegion(s string) S3URIOpt {
	return func(s3u *S3URI) {
		s3u.Region = String(s)
	}
}

func WithNormalizedKey(b bool) S3URIOpt {
	return func(s3u *S3URI) {
		s3u.normalize = Bool(b)
	}
}

type S3URI struct {
	uri       *url.URL
	options   []S3URIOpt
	normalize *bool

	HostStyle   *bool
	PathStyle   *bool
	Accelerated *bool
	DualStack   *bool
	Website     *bool

	Scheme    *string
	Bucket    *string
	Key       *string
	VersionID *string
	Region    *string
}

func NewS3URI(opts ...S3URIOpt) *S3URI {
	return &S3URI{options: opts}
}

func (s3u *S3URI) Reset() *S3URI {
	return reset(s3u)
}

func (s3u *S3URI) Parse(v interface{}) (*S3URI, error) {
	return parse(s3u, v)
}

func (s3u *S3URI) ParseURL(u *url.URL) (*S3URI, error) {
	return parse(s3u, u)
}

func (s3u *S3URI) ParseString(s string) (*S3URI, error) {
	return parse(s3u, s)
}

func (s3u *S3URI) URI() *url.URL {
	return s3u.uri
}

func Parse(v interface{}) (*S3URI, error) {
	return NewS3URI().Parse(v)
}

func ParseURL(u *url.URL) (*S3URI, error) {
	return NewS3URI().ParseURL(u)
}

func ParseString(s string) (*S3URI, error) {
	return NewS3URI().ParseString(s)
}

func MustParse(s3u *S3URI, err error) *S3URI {
	if err != nil {
		panic(err)
	}
	return s3u
}

func Validate(v interface{}) bool {
	_, err := NewS3URI().Parse(v)
	return err == nil
}

func ValidateURL(u *url.URL) bool {
	_, err := NewS3URI().Parse(u)
	return err == nil
}

func ValidateString(s string) bool {
	_, err := NewS3URI().Parse(s)
	return err == nil
}

func parse(s3u *S3URI, s interface{}) (*S3URI, error) {
	var (
		u   *url.URL
		err error
	)

	switch s := s.(type) {
	case string:
		u, err = url.Parse(s)
	case *url.URL:
		u = s
	default:
		return nil, fmt.Errorf("unable to parse unknown type: %T", s)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to parse given S3 endpoint URL: %w", err)
	}

	reset(s3u)
	s3u.uri = u

	switch u.Scheme {
	case "s3", "http", "https":
		s3u.Scheme = String(u.Scheme)
	default:
		return nil, fmt.Errorf("unable to parse schema type: %s", u.Scheme)
	}

	// Handle S3 endpoint URL with the schema s3:// that is neither
	// the host style nor the path style.
	if u.Scheme == "s3" {
		if u.Host == "" {
			return nil, ErrBucketNotFound
		}
		s3u.Bucket = String(u.Host)

		if u.Path != "" && u.Path != "/" {
			s3u.Key = String(u.Path[1:len(u.Path)])
		}
		s3u.Region = String(DefaultRegion)

		return s3u, nil
	}

	if u.Host == "" {
		return nil, ErrHostnameNotFound
	}

	matches := s3URLPattern.FindStringSubmatch(u.Host)
	if matches == nil || len(matches) < 1 {
		matches := digitalOceanURLPatter.FindStringSubmatch(u.Host)
		if matches == nil || len(matches) < 1 {
			matches = minioURLPattern.FindStringSubmatch(u.Host)
			if matches == nil || len(matches) < 1 {
				matches = localhostURLPattern.FindStringSubmatch(u.Host)
				if matches == nil || len(matches) < 1 {
					matches = ipURLPattern.FindStringSubmatch(u.Host)
					if matches == nil || len(matches) < 1 {
						return nil, ErrInvalidS3Endpoint
					}
				}
			}
		}
	}

	// prefix := matches[1]
	// usage := matches[2] // Type of the S3 bucket.
	// region := matches[3]

	prefix := ""
	usage := "" // Type of the S3 bucket.
	region := ""

	if prefix == "" {
		s3u.PathStyle = Bool(true)

		if u.Path != "" && u.Path != "/" {
			u.Path = u.Path[1:len(u.Path)]

			index := strings.Index(u.Path, "/")
			switch {
			case index == -1:
				s3u.Bucket = String(u.Path)
			case index == len(u.Path)-1:
				s3u.Bucket = String(u.Path[:index])
			default:
				s3u.Bucket = String(u.Path[:index])
				s3u.Key = String(u.Path[index+1:])
			}
		}
	} else {
		s3u.HostStyle = Bool(true)
		s3u.Bucket = String(prefix[:len(prefix)-1])

		if u.Path != "" && u.Path != "/" {
			s3u.Key = String(u.Path[1:len(u.Path)])
		}
	}

	const (
		// Used to denote type of the S3 bucket.
		accelerated = "accelerated"
		dualStack   = "dualstack"
		website     = "website"

		// Part of the amazonaws.com domain name.  Set when no region
		// could be ascertain correctly using the S3 endpoint URL.
		amazonAWS = "amazonaws"

		// Part of the query parameters.  Used when retrieving S3
		// object (key) of a particular version.
		versionID = "versionId"
	)

	// An S3 bucket can be either accelerated or website endpoint,
	// but not both.
	if usage == accelerated {
		s3u.Accelerated = Bool(true)
	} else if usage == website {
		s3u.Website = Bool(true)
	}

	// An accelerated S3 bucket can also be dualstack.
	if usage == dualStack || region == dualStack {
		s3u.DualStack = Bool(true)
	}

	// Handle the special case of an accelerated dualstack S3
	// endpoint URL:
	//   <BUCKET>.s3-accelerated.dualstack.amazonaws.com/<KEY>.
	// As there is no way to accertain the region solely based on
	// the S3 endpoint URL.
	if usage != accelerated {
		s3u.Region = String(DefaultRegion)
		if region != amazonAWS {
			s3u.Region = String(region)
		}
	}

	// Query string used when requesting a particular version of a given
	// S3 object (key).
	if s := u.Query().Get(versionID); s != "" {
		s3u.VersionID = String(s)
	}

	// Apply options that serve as overrides after the initial parsing
	// is completed.  This allows for bucket name, key, version ID, etc.,
	// to be overridden at the parsing stage.
	for _, o := range s3u.options {
		o(s3u)
	}

	// Remove trailing slash from the key name, so that the "key/" will
	// become "key" and similarly "a/complex/key/" will simply become
	// "a/complex/key" afer being normalized.
	if BoolValue(s3u.normalize) && s3u.Key != nil {
		k := StringValue(s3u.Key)
		if k[len(k)-1] == '/' {
			k = k[:len(k)-1]
		}
		s3u.Key = String(k)
	}

	return s3u, nil
}

// Reset fields in the S3URI type, and set boolean values to false.
func reset(s3u *S3URI) *S3URI {
	*s3u = S3URI{
		HostStyle:   Bool(false),
		PathStyle:   Bool(false),
		Accelerated: Bool(false),
		DualStack:   Bool(false),
		Website:     Bool(false),
	}
	return s3u
}

func String(s string) *string {
	return &s
}

func Bool(b bool) *bool {
	return &b
}

func StringValue(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func BoolValue(b *bool) bool {
	if b != nil {
		return *b
	}
	return false
}

func main() {
	s3u := NewS3URI()
	//	fmt.Println(s3u.ParseString("s3://test123"))
	//	fmt.Println(s3u.URI())
	//	spew.Dump(s3u)
	//	s3u.Bucket = String("test")
	//	fmt.Println(s3u.URI().String())
	//	spew.Dump(s3u)

	fmt.Println(s3u.ParseString("s3://test123/"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("s3://test123/key456"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("s3://test123/key456/"))
	spew.Dump(s3u)
	//	fmt.Println(s3u.ParseString("https://s3.amazonaws.com/test123"))
	//	fmt.Println(s3u.ParseString("https://s3.amazonaws.com/test123/"))
	//	fmt.Println(s3u.ParseString("https://s3.amazonaws.com/test123/key456"))
	//	fmt.Println(s3u.ParseString("https://s3.amazonaws.com/test123/key456/"))
	fmt.Println(s3u.ParseString("https://s3-eu-west-1.amazonaws.com/test123/key456/"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("https://s3.eu-west-1.amazonaws.com/test123/key456/"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("https://s3.dualstack.eu-west-1.amazonaws.com/test123/key456/"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("https://test123.s3-website-eu-west-1.amazonaws.com/key456/"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("https://test123.s3-accelerated.amazonaws.com/key456/"))
	spew.Dump(s3u)
	fmt.Println(s3u.ParseString("https://test123.s3-accelerated.dualstack.amazonaws.com/key456/"))
	spew.Dump(s3u)
	//	fmt.Println(s3u.ParseString("https://test123.s3.amazonaws.com/"))
	//	fmt.Println(s3u.ParseString("https://test123.s3.amazonaws.com/key456"))
	//	fmt.Println(s3u.ParseString("https://test123.s3.amazonaws.com/key456"))
	//	fmt.Println(s3u.ParseString("https://google.com")) // invalid S3 endpoint

	//	fmt.Println(s3u.ParseString("https://test123.s3.amazonaws.com/key456?versionId=123456&x=1&y=2&y=3;z"))
	//	fmt.Println(*s3u.Bucket, *s3u.Key, *s3u.Region, *s3u.PathStyle, *s3u.VersionID)
	//	fmt.Println(s3u.URI().Scheme)

	//	fmt.Println(s3u.ParseString("https://s3-eu-west-1.amazonaws.com/test123/key456?t=this+is+a+simple+%26+short+test."))

	//	u, _ := url.Parse("s3://test123/key456")
	//	fmt.Println(s3u.Parse(u))

	//	fmt.Println(MustParse(s3u.ParseString("s3://test123/key456")))
	//	// Will panic: no hostname
	//	// fmt.Println(MustParse(s3u.ParseString("")))

	//	s3u = NewS3URI(
	//		WithRegion("eu-west-1"),
	//		WithVersionID("12341234"),
	//		WithNormalizedKey(true),
	//	)
	//	spew.Dump(s3u.URI())
	//	fmt.Println(s3u.ParseString("https://test123.s3.amazonaws.com/key456/?versionId=123456&x=1&y=2&y=3;z"))
	//	fmt.Println(*s3u.Bucket, *s3u.Key, *s3u.Region, *s3u.PathStyle, *s3u.VersionID)
	//	fmt.Println(s3u.URI().Scheme)
	//	spew.Dump(s3u.URI())
	fmt.Println(Validate("https://test123.s3-accelerated.dualstack.amazonaws.com/key456/"))
	fmt.Println(Validate("ftp://google.com/"))
	fmt.Println(ParseString("ftp://google.com/"))
}
