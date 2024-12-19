package xlog

import (
	"context"
	"regexp"
	"runtime"
	"sync"

	"golang.org/x/sync/errgroup"
)

// a List of directories that should be ignored by directory walking function.
// for example the versioning extension can register `.versions` directory to be
// ignored
var ignoredPaths = []*regexp.Regexp{
	regexp.MustCompile(`^\.`), // Ignore any hidden directory
}

// IgnorePath Register a pattern to be ignored when walking directories.
func IgnorePath(r *regexp.Regexp) {
	ignoredPaths = append(ignoredPaths, r)
}

func IsIgnoredPath(d string) bool {
	for _, v := range ignoredPaths {
		if v.MatchString(d) {
			return true
		}
	}

	return false
}

var pages []Page

func Pages(ctx context.Context) []Page {
	if pages == nil {
		populatePagesCache(ctx)
	}

	return pages[:]
}

// EachPage iterates on all available pages. many extensions
// uses it to get all pages and maybe parse them and extract needed information
func EachPage(ctx context.Context, f func(Page)) {
	if pages == nil {
		populatePagesCache(ctx)
	}

	currentPages := pages
	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			return
		default:
			f(p)
		}
	}
}

var concurrency = runtime.NumCPU() * 4

// EachPageCon Similar to EachPage but iterates concurrently
func EachPageCon(ctx context.Context, f func(Page)) {
	if pages == nil {
		populatePagesCache(ctx)
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(concurrency)

	currentPages := pages
	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			break
		default:
			grp.Go(func() (err error) { f(p); return })
		}
	}

	grp.Wait()
}

// MapPageCon Similar to EachPage but iterates concurrently
func MapPageCon[T any](ctx context.Context, f func(Page) T) []T {
	if pages == nil {
		populatePagesCache(ctx)
	}

	grp, ctx := errgroup.WithContext(ctx)
	grp.SetLimit(concurrency)

	currentPages := pages
	output := make([]T, 0, len(currentPages))
	var outputLck sync.Mutex

	for _, p := range currentPages {
		select {
		case <-ctx.Done():
			break
		default:
			grp.Go(func() (err error) {
				val := f(p)
				if any(val) == nil {
					return
				}

				outputLck.Lock()
				output = append(output, val)
				outputLck.Unlock()

				return
			})
		}
	}

	grp.Wait()

	return output
}

func clearPagesCache(p Page) (err error) {
	pages = nil
	return nil
}

func populatePagesCache(ctx context.Context) {
	pages = []Page{}
	for _, s := range sources {
		select {
		case <-ctx.Done():
			return
		default:
			s.Each(ctx, func(p Page) {
				pages = append(pages, p)
			})
		}
	}
}
