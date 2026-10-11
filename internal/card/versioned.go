package card

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteVersioned renames dir/{dark,light}_mode.svg to dir/{dark,light}_mode-<version>.svg
// and removes the numbered files from earlier runs, so only one pair is left. The README links to the versioned names:
// a new path is always a new URL, so GitHub's image proxy can never serve a stale card.
func WriteVersioned(dir, version string) error {
	old, err := filepath.Glob(filepath.Join(dir, "*_mode-*.svg"))
	if err != nil {
		return err
	}
	for _, f := range old {
		if err := os.Remove(f); err != nil {
			return fmt.Errorf("remove %s: %w", f, err)
		}
	}
	for _, name := range []string{"dark", "light"} {
		src := filepath.Join(dir, name+"_mode.svg")
		dst := filepath.Join(dir, fmt.Sprintf("%s_mode-%s.svg", name, version))
		if err := os.Rename(src, dst); err != nil {
			return fmt.Errorf("rename %s: %w", src, err)
		}
	}
	return nil
}
