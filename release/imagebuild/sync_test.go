package imagebuild

import (
	"testing"
)

func TestStripBuildSuffix(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		// no build suffix — unchanged
		{"v3.6.7", "v3.6.7"},
		{"v3.6.7-k3s1", "v3.6.7-k3s1"},
		// build suffix stripped, k3s preserved
		{"v3.6.7-build20260415", "v3.6.7"},
		{"v3.6.7-k3s1-build20260415", "v3.6.7-k3s1"},
		{"v3.6.7-k3s2-build20260415", "v3.6.7-k3s2"},
		{"v1.32.3-k3s1-build20260101", "v1.32.3-k3s1"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := buildSuffixRE.ReplaceAllString(tt.input, "")
			if got != tt.want {
				t.Errorf("stripBuildSuffix(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateTagFormat(t *testing.T) {
	tests := []struct {
		tagName   string
		tagPrefix string
		want      bool
	}{
		// plain semver – valid
		{"v3.21.1", "", true},
		{"v1.32.3", "", true},
		// k3s prerelease – valid
		{"v3.6.7-k3s1", "", true},
		{"v1.32.3-k3s2", "", true},
		// build suffix – invalid (treated the same as any other unrecognised suffix)
		{"v3.6.7-build20260415", "", false},
		{"v3.6.7-k3s1-build20260415", "", false},
		// other prerelease suffixes – invalid
		{"v3.21.1-typha", "", false},
		{"v3.21.1-pod2daemon", "", false},
		{"v3.24.2-0.dev", "", false},
		{"v3.21.1-rc1", "", false},
		{"v3.21.1-alpha", "", false},
		{"v3.21.1-beta", "", false},
		// with tagPrefix
		{"kubernetes-v1.32.3", "kubernetes-", true},
		{"kubernetes-v1.32.3-k3s1", "kubernetes-", true},
		{"kubernetes-v1.32.3-rc1", "kubernetes-", false},
		// missing prefix
		{"v1.32.3", "kubernetes-", false},
	}

	for _, tt := range tests {
		name := tt.tagName
		if tt.tagPrefix != "" {
			name += " (prefix=" + tt.tagPrefix + ")"
		}
		t.Run(name, func(t *testing.T) {
			got := validateTagFormat(tt.tagName, tt.tagPrefix)
			if got != tt.want {
				t.Errorf("validateTagFormat(%q, %q) = %v, want %v", tt.tagName, tt.tagPrefix, got, tt.want)
			}
		})
	}
}
