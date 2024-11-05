#!/bin/bash

export ARTIFACT_DIR="/workspace/artifacts"
export OCI_STORAGE_USERNAME="$(jq -r '."quay-username"' /usr/local/konflux-test-infra/oci-storage)"
export OCI_STORAGE_TOKEN="$(jq -r '."quay-token"' /usr/local/konflux-test-infra/oci-storage)"

post_actions() {
    # Capture the exit code of the previous command
    local exit_code=$?
    local temp_annotation_file="$(mktemp)"

    # Change to the artifact directory
    cd "$ARTIFACT_DIR" || exit 1

    # Fetch and process manifest annotations
    if ! MANIFESTS=$(oras manifest fetch "${OCI_STORAGE_CONTAINER}" | jq .annotations); then
        echo "Error: Failed to fetch manifest from ${OCI_STORAGE_CONTAINER}"
        exit 1
    fi

    # Save the manifest annotations to a temporary file
    jq -n --argjson manifest "$MANIFESTS" '{ "manifest": $manifest }' > "${temp_annotation_file}"

    # Pull the container from OCI storage
    oras pull "${OCI_STORAGE_CONTAINER}"

    # Attempt to push back to OCI storage with retries
    local attempt=1
    while ! oras push "$OCI_STORAGE_CONTAINER" \
        --username="${OCI_STORAGE_USERNAME}" \
        --password="${OCI_STORAGE_TOKEN}" \
        --annotation-file "${temp_annotation_file}" \
        ./:application/vnd.acme.rocket.docs.layer.v1+tar; do
        
        if [ "$attempt" -ge 5 ]; then
            echo "Error: oras push failed after $attempt attempts."
            exit 1
        fi
        echo "Warning: oras push failed (attempt $attempt). Retrying in 5 seconds..."
        sleep 5
        attempt=$((attempt + 1))
    done

    # Exit with the original command's exit code
    exit "$exit_code"
}

# Ensure post_actions is called when the script exits
trap post_actions EXIT

# Install Ginkgo and run tests
go install github.com/onsi/ginkgo/v2/ginkgo@latest
ginkgo --junit-report "$ARTIFACT_DIR"/junit_sample.xml
