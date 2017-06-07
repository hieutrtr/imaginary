package main

import (
	"net/http"
	"net/url"
)

type ImageSourceType string
type ImageSourceFactoryFunction func(*SourceConfig) ImageSource

type SourceConfig struct {
	AuthForwarding  bool
	Authorization   string
	MountPath       string
	Type            ImageSourceType
	AllowedOrigings []*url.URL
	MaxAllowedSize  int
	EnableCeph      bool
	CephConfig      string
	UseCephBlock    bool
	CephBlockURL    string
	EnableS3        bool
}

var imageSourceMap = make(map[ImageSourceType]ImageSource)
var imageSourceFactoryMap = make(map[ImageSourceType]ImageSourceFactoryFunction)

type ImageSource interface {
	Matches(*http.Request) bool
	GetImage(*http.Request) ([]byte, error)
	GetCache(*http.Request) ([]byte, error)
	Delete(*http.Request) error
}

func RegisterSource(sourceType ImageSourceType, factory ImageSourceFactoryFunction) {
	imageSourceFactoryMap[sourceType] = factory
}

func LoadSources(o ServerOptions) {
	for name, factory := range imageSourceFactoryMap {
		imageSourceMap[name] = factory(&SourceConfig{
			Type:            name,
			MountPath:       o.Mount,
			AuthForwarding:  o.AuthForwarding,
			Authorization:   o.Authorization,
			AllowedOrigings: o.AlloweOrigins,
			MaxAllowedSize:  o.MaxAllowedSize,
			EnableCeph:      o.EnableCeph,
			CephConfig:      o.CephConfig,
			UseCephBlock:    o.UseCephBlock,
			CephBlockURL:    o.CephBlockURL,
			EnableS3:        o.EnableS3,
		})
	}
}

func MatchSource(req *http.Request) ImageSource {
	for _, source := range imageSourceMap {
		if source != nil && source.Matches(req) {
			return source
		}
	}
	return nil
}
