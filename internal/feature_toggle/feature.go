package feature_toggle

import (
	"context"
	"github.com/spf13/viper"
)

type Feature func(ctx context.Context) error

// FeatureToggle toggles a feature based on a config property
func FeatureToggle(ctx context.Context, propKey string, f Feature) error {
	enabled := viper.GetBool(propKey)

	if enabled {
		return f(ctx)
	}
	return nil
}
