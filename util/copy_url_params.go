package util

import "net/url"

func CopyUrlParams(src url.Values, dest url.Values, keys []string) {
	if keys == nil || len(keys) == 0 {
		for key := range src {
			dest.Set(key, src.Get(key))
		}
	} else {
		for _, key := range keys {
			val := src.Get(key)
			if val != "" {
				dest.Set(key, val)
			}
		}
	}
}
