package selfupdate

import "testing"

func TestReexecSelf_usesSeam(t *testing.T) {
	orig := reexecFn
	t.Cleanup(func() { reexecFn = orig })

	var gotExe string
	var gotArgs []string
	reexecFn = func(exe string, args []string) error {
		gotExe = exe
		gotArgs = args
		return nil
	}

	// reexecFn hides os.Executable details here; we only check args propagate.
	if err := ReexecSelf("a", "b", "c"); err != nil {
		t.Fatal(err)
	}
	if gotExe == "" {
		t.Fatal("exe should be non-empty")
	}
	if len(gotArgs) != 3 || gotArgs[0] != "a" || gotArgs[1] != "b" || gotArgs[2] != "c" {
		t.Fatalf("args = %v", gotArgs)
	}
}

func TestReexecSelf_propagatesSeamError(t *testing.T) {
	orig := reexecFn
	t.Cleanup(func() { reexecFn = orig })
	reexecFn = func(_ string, _ []string) error {
		return errSentinel
	}
	if err := ReexecSelf("x"); err != errSentinel {
		t.Fatalf("got %v, want errSentinel", err)
	}
}

var errSentinel = newSentinel()

type sentinelErr struct{}

func (s *sentinelErr) Error() string { return "sentinel" }

func newSentinel() error { return &sentinelErr{} }
