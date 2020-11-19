package converter

import (
	"fmt"

	// modules
	_ "github.com/yankghjh/selfhosted_store/cli/modules/docker-compose"

	"github.com/yankghjh/selfhosted_store/cli/project"
)

// Convert application form source format to target format
func Convert(srcFormat, dstFormat string, payload []byte) ([]byte, error) {
	decoder := project.LoadDecoder(srcFormat)
	if decoder == nil {
		return nil, fmt.Errorf("decoder %s not found", srcFormat)
	}

	applications, err := decoder(payload)
	if err != nil {
		return nil, err
	}

	encoder := project.LoadEncoder(dstFormat)
	if encoder == nil {
		return nil, fmt.Errorf("encoder %s not found", dstFormat)
	}

	result, err := encoder(applications)
	if err != nil {
		return nil, err
	}

	return result, nil
}
