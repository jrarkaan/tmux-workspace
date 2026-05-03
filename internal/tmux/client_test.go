package tmux

import (
	"errors"
	"reflect"
	"testing"
)

type fakeRunner struct {
	output    []byte
	err       error
	responses []fakeResponse
	calls     []fakeCall
}

type fakeCall struct {
	name string
	args []string
}

type fakeResponse struct {
	output []byte
	err    error
}

func (r *fakeRunner) Run(name string, args ...string) ([]byte, error) {
	r.calls = append(r.calls, fakeCall{name: name, args: args})
	if len(r.responses) > 0 {
		response := r.responses[0]
		r.responses = r.responses[1:]
		return response.output, response.err
	}

	return r.output, r.err
}

func (r *fakeRunner) RunInteractive(name string, args ...string) error {
	r.calls = append(r.calls, fakeCall{name: name, args: args})
	if len(r.responses) > 0 {
		response := r.responses[0]
		r.responses = r.responses[1:]
		return response.err
	}

	return r.err
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

func TestClientMutationCommandsBuildExpectedArgs(t *testing.T) {
	tests := []struct {
		name string
		run  func(*Client) error
		want fakeCall
	}{
		{
			name: "new session",
			run: func(client *Client) error {
				return client.NewSession("backend-dev", "overview", "/tmp/backend")
			},
			want: fakeCall{name: "tmux", args: []string{"new-session", "-d", "-s", "backend-dev", "-n", "overview", "-c", "/tmp/backend"}},
		},
		{
			name: "new window",
			run: func(client *Client) error {
				return client.NewWindow("backend-dev", "test-watch", "/tmp/backend")
			},
			want: fakeCall{name: "tmux", args: []string{"new-window", "-t", "backend-dev", "-n", "test-watch", "-c", "/tmp/backend"}},
		},
		{
			name: "send keys",
			run: func(client *Client) error {
				return client.SendKeys("backend-dev:overview", "git status")
			},
			want: fakeCall{name: "tmux", args: []string{"send-keys", "-t", "backend-dev:overview", "git status", "C-m"}},
		},
		{
			name: "select window",
			run: func(client *Client) error {
				return client.SelectWindow("backend-dev", "overview")
			},
			want: fakeCall{name: "tmux", args: []string{"select-window", "-t", "backend-dev:overview"}},
		},
		{
			name: "attach",
			run: func(client *Client) error {
				return client.Attach("backend-dev")
			},
			want: fakeCall{name: "tmux", args: []string{"attach", "-t", "backend-dev"}},
		},
		{
			name: "kill session",
			run: func(client *Client) error {
				return client.KillSession("backend-dev")
			},
			want: fakeCall{name: "tmux", args: []string{"kill-session", "-t", "backend-dev"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			runner := &fakeRunner{}
			client := newInstalledTestClient(runner)

			if err := tt.run(client); err != nil {
				t.Fatalf("command returned error: %v", err)
			}

			if !reflect.DeepEqual(runner.calls, []fakeCall{tt.want}) {
				t.Fatalf("calls = %#v, want %#v", runner.calls, []fakeCall{tt.want})
			}
		})
	}
}

func newInstalledTestClient(runner Runner) *Client {
	client := NewClientWithRunner(runner)
	client.lookPath = func(file string) (string, error) {
		return "/usr/bin/tmux", nil
	}

	return client
}
