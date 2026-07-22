package segment

import "strings"

// Manual per-artifact metadata. Product classes come from findings F4
// (component distributions are product-class-conditional); capture-damage
// flags come from findings F9. Keyed by corpus basename prefix so version
// bumps in filenames keep matching. Unknown files map to "unknown" — the
// segment command warns on those instead of guessing silently.

var productClassRules = []struct{ prefix, class string }{
	{"anthropic--claude-ai--", "chat"},
	{"openai--chatgpt--", "chat"},
	{"google--gemini--", "chat"},
	{"xai--grok--", "chat"},
	{"anthropic--claude-code--", "coding-agent"},
	{"anysphere--cursor-ide--", "coding-agent"},
	{"windsurf--cascade--", "coding-agent"},
	{"cognition--devin--", "coding-agent"},
	{"openai--codex-cli--", "coding-agent"},
	{"cline--cline--", "coding-agent"},
	{"all-hands-ai--openhands--", "coding-agent"},
	{"aider-ai--aider--", "coding-agent"},
	{"vercel--v0--", "app-builder"},
	{"replit--replit-agent--", "app-builder"},
}

var fileFlagRules = []struct {
	prefix string
	flags  []string
}{
	{"cognition--devin--", []string{"wrapper", "aggregate-member", "session-fill"}},
	{"replit--replit-agent--", []string{"wrapper", "stripped-tags", "aggregate-member"}},
	{"vercel--v0--", []string{"truncated"}},
	{"anthropic--claude-code--", []string{"session-fill"}},
	{"openai--chatgpt--", []string{"session-fill"}},
	{"google--gemini--", []string{"session-fill"}},
	{"xai--grok--", []string{"session-fill"}},
	{"windsurf--cascade--", []string{"session-fill"}},
}

// Class returns the product class for a corpus file basename.
func Class(base string) string {
	for _, r := range productClassRules {
		if strings.HasPrefix(base, r.prefix) {
			return r.class
		}
	}
	return "unknown"
}

// FileFlags returns the capture-damage flags for a corpus file basename.
// The returned slice is a copy.
func FileFlags(base string) []string {
	for _, r := range fileFlagRules {
		if strings.HasPrefix(base, r.prefix) {
			return append([]string(nil), r.flags...)
		}
	}
	return nil
}
