package uploader

import (
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Stuhub-io/config"
	"github.com/Stuhub-io/core/domain"
	"github.com/Stuhub-io/core/ports"
	"github.com/cloudinary/cloudinary-go/api"
)

type CloudinaryUploader struct {
	cfg config.Config
}

func NewCloudinaryUploader(cfg config.Config) ports.Uploader {
	return &CloudinaryUploader{cfg}
}

func ignoredSignatureKey(key string) bool {
	switch key {
	case "file", "cloud_name", "resource_type", "api_key":
		return true
	}
	return false
}

func (u *CloudinaryUploader) SignRequest(input domain.SignUrlInput) (s *domain.SignedUrl, e *domain.Error) {
	// https://cloudinary.com/documentation/upload_images#generating_authentication_signatures
	// All parameters added to the method call should be included except: file, cloud_name, resource_type and your api_key.

	var arrayKey = regexp.MustCompile(`(.*)\[\d+]`)

	requestParams, err := url.ParseQuery(input.AdditionalQuery)
	if err != nil {
		return nil, domain.NewErr("Invalid Query", domain.BadRequestCode)
	}

	requestParams.Set("public_id", input.PublicID)

	requestKeys := make([]string, 0, len(requestParams))
	for k := range requestParams {
		requestKeys = append(requestKeys, k)
	}
	sort.Strings(requestKeys)

	signatureParams := make(url.Values)

	for _, k := range requestKeys {
		switch {
		case arrayKey.MatchString(k):
			kName := arrayKey.FindStringSubmatch(k)[1]
			signatureParams[kName] = append(signatureParams[kName], requestParams[k][0])
		case ignoredSignatureKey(k):
			// omit
		default:
			signatureParams[k] = requestParams[k]
		}
	}

	for k, v := range signatureParams {
		signatureParams[k] = []string{strings.Join(v, ",")}
	}

	signatureParams.Set("timestamp", strconv.FormatInt(time.Now().Unix(), 10))

	signature, err := api.SignParameters(signatureParams, u.cfg.CloudinarySecretKey)

	if err != nil {
		return nil, domain.NewErr("Invalid signature string input", domain.BadRequestCode)
	}

	return &domain.SignedUrl{
		Url:       u.cfg.CloudinaryBaseURL + "/" + input.ResourceType.String() + "/upload",
		Signature: signature,
		ApiKey:    u.cfg.CloudinaryApiKey,
		Timestamp: signatureParams.Get("timestamp"),
	}, nil
}
