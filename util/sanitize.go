package util

import (
	bm "github.com/microcosm-cc/bluemonday"
)

var policy *bm.Policy = nil

func setup_policy() {
	if policy == nil {
		policy = bm.UGCPolicy()
	}
}

func SanitizeBytes(content []byte) []byte {
	setup_policy()
	return policy.SanitizeBytes(content)
}

func Sanitize(content string) string {
	setup_policy()
	return policy.Sanitize(content)
}
