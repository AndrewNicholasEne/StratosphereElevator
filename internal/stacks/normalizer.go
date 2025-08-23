package stacks

import "strings"

func normalizeCreateInput(in CreateInput) (string, string, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" || len(name) > 200 {
		return "", "", ErrInvalidInput
	}
	// if slug provided, prefer it; else derive from name
	base := name
	if in.Slug != nil && strings.TrimSpace(*in.Slug) != "" {
		base = *in.Slug
	}
	slug := slugify(base)
	if slug == "" || len(slug) > 100 {
		return "", "", ErrInvalidInput
	}
	return name, slug, nil
}

// Minimal slugifier: lowercase, [a-z0-9-], collapse non-alnum to single '-'
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	dash := false
	for _, r := range s {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
			dash = false
		default:
			if !dash {
				b.WriteByte('-')
				dash = true
			}
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		out = "stack"
	}
	return out
}
