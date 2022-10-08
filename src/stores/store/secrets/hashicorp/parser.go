package hashicorp

import (
	"encoding/json"
	"time"

	"github.com/longfan78/quorum-key-manager/pkg/errors"

	"github.com/hashicorp/vault/api"

	"github.com/longfan78/quorum-key-manager/src/stores/entities"
)

func formatHashicorpSecret(id, value string, tags map[string]string, metadata *entities.Metadata) *entities.Secret {
	return &entities.Secret{
		ID:       id,
		Value:    value,
		Tags:     tags,
		Metadata: metadata,
	}
}

func formatHashicorpSecretMetadata(secret *api.Secret, version string) (*entities.Metadata, error) {
	jsonMetadata := secret.Data

	if version == "" {
		version = jsonMetadata["current_version"].(json.Number).String()
	}

	metadata := &entities.Metadata{
		Version: version,
	}

	secretVersion := jsonMetadata["versions"].(map[string]interface{})[version].(map[string]interface{})
	if secretVersion["deletion_time"].(string) != "" {
		deletionTime, err := time.Parse(time.RFC3339, secretVersion["deletion_time"].(string))
		if err != nil {
			return nil, errors.HashicorpVaultError("failed to parse deletion time from metadata")
		}

		metadata.DeletedAt = deletionTime
		metadata.Disabled = true
	}

	var err error
	metadata.CreatedAt, err = time.Parse(time.RFC3339, secretVersion["created_time"].(string))
	if err != nil {
		return nil, errors.HashicorpVaultError("failed to parse created time from metadata")
	}
	metadata.UpdatedAt = metadata.CreatedAt

	expirationDurationStr := jsonMetadata["delete_version_after"].(string)
	if expirationDurationStr != "0s" {
		expirationDuration, der := time.ParseDuration(expirationDurationStr)
		if der != nil {
			return nil, errors.HashicorpVaultError("failed to parse expiration time from metadata")
		}

		metadata.ExpireAt = metadata.CreatedAt.Add(expirationDuration)
	}

	return metadata, nil
}

func formatTags(tagsI map[string]interface{}) map[string]string {
	if tagsI == nil {
		return nil
	}

	tags := make(map[string]string)
	for i, tag := range tagsI {
		tags[i] = tag.(string)
	}

	return tags
}
