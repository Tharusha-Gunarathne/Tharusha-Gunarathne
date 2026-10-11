package card

import (
	"fmt"
	"os"
	"path/filepath"
)

// WriteVersioned copies dir/{dark,light}_mode.svg to dir/{dark,light}_mode-<version>.svg
// and removes copies from earlier runs. The README links to the versioned names:
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
		data, err := os.ReadFile(filepath.Join(dir, name+"_mode.svg"))
		if err != nil {
			return err
		}
		dst := filepath.Join(dir, fmt.Sprintf("%s_mode-%s.svg", name, version))
		if err := os.WriteFile(dst, data, 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dst, err)
		}
	}
	return nil
}
