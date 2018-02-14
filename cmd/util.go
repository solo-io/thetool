package cmd

import (
	"crypto/sha256"
	"fmt"

	"github.com/solo-io/thetool/pkg/feature"
)

func loadEnabledFeatures() ([]feature.Feature, error) {
	features, err := feature.LoadFromFile(dataFile)
	if err != nil {
		return nil, err
	}
	var enabled []feature.Feature
	for _, f := range features {
		if f.Enabled {
			enabled = append(enabled, f)
		}
	}
	return enabled, nil
}

// featuresHash generates a hash for particular envoy and glue build
// based on the features included
func featuresHash(features []feature.Feature) string {
	hash := sha256.New()
	for _, f := range features {
		hash.Write([]byte(f.Name))
		hash.Write([]byte(f.Version))
	}

	return fmt.Sprintf("%x", hash.Sum(nil))[:8]
}
