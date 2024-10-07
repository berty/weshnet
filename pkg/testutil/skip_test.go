package testutil_test

import (
	"testing"

	"berty.tech/weshnet/v2/pkg/testutil"
)

func TestFilterRace(t *testing.T) {
	t.Run("always-run", func(t *testing.T) {})
	t.Run("skip-if-race", func(t *testing.T) {
		testutil.FilterRace(t, testutil.SkipIfRace)
	})
	t.Run("run-if-race", func(t *testing.T) {
		testutil.FilterRace(t, testutil.RunIfRace)
	})
	t.Run("always-skip", func(t *testing.T) {
		t.Skip()
	})
}
