package tmux

import (
	"errors"
	"reflect"
	"testing"
)

type fakeRunner struct {
	output []byte
	err    error
	calls  []fakeCall
}

type fakeCall struct {
	name string
	args []string
}

func (r *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	r.calls = append(r.calls, fakeCall{name: name, args: args})
	return r.output, r.err
}

func TestListSessionsParsesSessionNames(t *testing.T) {
	runner := &fakeRunner{output: []byte("backend-dev\nfrontend-astro-blog\n")}
	client := newInstalledTestClient(runner)

	sessions, err := client.ListSessions()
	if err != nil {
		t.Fatalf("ListSessions returned error: %v", err)
	}

	want := []string{"backend-dev", "frontend-astro-blog"}
	if !reflect.DeepEqual(sessions, want) {
		t.Fatalf("ListSessions() = %#v, want %#v", sessions, want)
	}
}

func TestListSessionsReturnsEmptyForNoServer(t *testing.T) {
	tests := []string{
		"no server running on /tmp/tmux-1000/default",
		"failed to connect to server",
		"error connecting to /tmp/tmux-1000/default",
		"no sessions",
	}

	for _, message := range tests {
		t.Run(message, func(t *testing.T) {
			runner := &fakeRunner{
				output: []byte(message),
				err:    errors.New("exit status 1"),
			}
			client := newInstalledTestClient(runner)

			sessions, err := client.ListSessions()
			if err != nil {
				t.Fatalf("ListSessions returned error: %v", err)
			}

			if len(sessions) != 0 {
				t.Fatalf("ListSessions() = %#v, want empty", sessions)
			}
		})
	}
}

func TestListWindowsParsesWindowOutput(t *testing.T) {
	runner := &fakeRunner{output: []byte("1:overview\n2:codex-impl\n3:test-watch\n")}
	client := newInstalledTestClient(runner)

	windows, err := client.ListWindows("backend-dev")
	if err != nil {
		t.Fatalf("ListWindows returned error: %v", err)
	}

	want := []WindowInfo{
		{Index: "1", Name: "overview"},
		{Index: "2", Name: "codex-impl"},
		{Index: "3", Name: "test-watch"},
	}
	if !reflect.DeepEqual(windows, want) {
		t.Fatalf("ListWindows() = %#v, want %#v", windows, want)
	}
}

func TestHasSessionReturnsTrueWhenRunnerSucceeds(t *testing.T) {
	runner := &fakeRunner{}
	client := newInstalledTestClient(runner)

	ok, err := client.HasSession("backend-dev")
	if err != nil {
		t.Fatalf("HasSession returned error: %v", err)
	}
	if !ok {
		t.Fatal("HasSession() = false, want true")
	}
}

func TestHasSessionReturnsFalseForMissingSession(t *testing.T) {
	runner := &fakeRunner{
		output: []byte("can't find session: backend-dev"),
		err:    errors.New("exit status 1"),
	}
	client := newInstalledTestClient(runner)

	ok, err := client.HasSession("backend-dev")
	if err != nil {
		t.Fatalf("HasSession returned error: %v", err)
	}
	if ok {
		t.Fatal("HasSession() = true, want false")
	}
}

func TestListWindowsReturnsSessionNotFound(t *testing.T) {
	runner := &fakeRunner{
		output: []byte("can't find session: backend-dev"),
		err:    errors.New("exit status 1"),
	}
	client := newInstalledTestClient(runner)

	_, err := client.ListWindows("backend-dev")
	if !errors.Is(err, ErrSessionNotFound) {
		t.Fatalf("ListWindows error = %v, want ErrSessionNotFound", err)
	}
}

func TestIsInstalledUsesLookPath(t *testing.T) {
	client := NewClientWithRunner(&fakeRunner{})
	client.lookPath = func(file string) (string, error) {
		return "", errors.New("not found")
	}

	if client.IsInstalled() {
		t.Fatal("IsInstalled() = true, want false")
	}
}

func newInstalledTestClient(runner Runner) *Client {
	client := NewClientWithRunner(runner)
	client.lookPath = func(file string) (string, error) {
		return "/usr/bin/tmux", nil
	}

	return client
}
