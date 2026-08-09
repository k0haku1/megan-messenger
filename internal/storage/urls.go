package storage

import "context"

// URLResolver turns object keys into client-facing URLs.
type URLResolver struct {
	store ObjectStore
}

func NewURLResolver(store ObjectStore) *URLResolver {
	return &URLResolver{store: store}
}

func (r *URLResolver) Resolve(ctx context.Context, key string) string {
	if r == nil || r.store == nil || key == "" {
		return ""
	}
	// Legacy local filenames (no slash) from old disk store.
	if !containsSlash(key) {
		return "/uploads/avatars/" + key
	}
	url, err := r.store.PresignGet(ctx, key, PresignTTL)
	if err != nil {
		return ""
	}
	return url
}

func containsSlash(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '/' {
			return true
		}
	}
	return false
}

func (r *URLResolver) ResolveMany(ctx context.Context, keys []string) map[string]string {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		if key == "" {
			continue
		}
		out[key] = r.Resolve(ctx, key)
	}
	return out
}
