package tmux

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jrarkaan/tmux-workspace/internal/config"
)

func TestStartWorkspaceCreatesSessionAndWindowsInOrder(t *testing.T) {
	runner := &fakeRunner{
		responses: []fakeResponse{
			{output: []byte("can't find session: backend-dev"), err: errors.New("exit status 1")},
		},
	}
	client := newInstalledTestClient(runner)
	workspace := testWorkspace(t)

	_, err := StartWorkspace(client, "backend-dev", workspace, WorkspaceStartOptions{NoAttach: true})
	if err != nil {
		t.Fatalf("StartWorkspace returned error: %v", err)
	}

	want := []fakeCall{
		{name: "tmux", args: []string{"has-session", "-t", "backend-dev"}},
		{name: "tmux", args: []string{"new-session", "-d", "-s", "backend-dev", "-n", "overview", "-c", workspace.Root}},
		{name: "tmux", args: []string{"send-keys", "-t", "backend-dev:overview", "export APP_ENV='development'; clear; pwd", "C-m"}},
		{name: "tmux", args: []string{"new-window", "-t", "backend-dev", "-n", "test-watch", "-c", workspace.Root}},
		{name: "tmux", args: []string{"send-keys", "-t", "backend-dev:test-watch", "export APP_ENV='development'; echo test", "C-m"}},
		{name: "tmux", args: []string{"new-window", "-t", "backend-dev", "-n", "logs", "-c", workspace.Root}},
		{name: "tmux", args: []string{"select-window", "-t", "backend-dev:overview"}},
	}

	if !reflect.DeepEqual(runner.calls, want) {
		t.Fatalf("calls = %#v, want %#v", runner.calls, want)
	}
}

func TestStartWorkspaceDoesNotAttachWhenNoAttach(t *testing.T) {
	runner := &fakeRunner{
		responses: []fakeResponse{
			{output: []byte("can't find session: backend-dev"), err: errors.New("exit status 1")},
		},
	}
	client := newInstalledTestClient(runner)

	_, err := StartWorkspace(client, "backend-dev", testWorkspace(t), WorkspaceStartOptions{NoAttach: true})
	if err != nil {
		t.Fatalf("StartWorkspace returned error: %v", err)
	}

	for _, call := range runner.calls {
		if len(call.args) > 0 && call.args[0] == "attach" {
			t.Fatalf("StartWorkspace unexpectedly attached: %#v", runner.calls)
		}
	}
}

func TestStartWorkspaceForceKillsExistingSessionBeforeRecreate(t *testing.T) {
	runner := &fakeRunner{}
	client := newInstalledTestClient(runner)
	workspace := testWorkspace(t)

	_, err := StartWorkspace(client, "backend-dev", workspace, WorkspaceStartOptions{NoAttach: true, Force: true})
	if err != nil {
		t.Fatalf("StartWorkspace returned error: %v", err)
	}

	wantPrefix := []fakeCall{
		{name: "tmux", args: []string{"has-session", "-t", "backend-dev"}},
		{name: "tmux", args: []string{"kill-session", "-t", "backend-dev"}},
		{name: "tmux", args: []string{"new-session", "-d", "-s", "backend-dev", "-n", "overview", "-c", workspace.Root}},
	}

	if len(runner.calls) < len(wantPrefix) || !reflect.DeepEqual(runner.calls[:len(wantPrefix)], wantPrefix) {
		t.Fatalf("call prefix = %#v, want %#v", runner.calls, wantPrefix)
	}
}

func TestStartWorkspaceExistingSessionWithoutForceDoesNotRecreate(t *testing.T) {
	runner := &fakeRunner{}
	client := newInstalledTestClient(runner)

	result, err := StartWorkspace(client, "backend-dev", testWorkspace(t), WorkspaceStartOptions{NoAttach: true})
	if err != nil {
		t.Fatalf("StartWorkspace returned error: %v", err)
	}

	wantCalls := []fakeCall{
		{name: "tmux", args: []string{"has-session", "-t", "backend-dev"}},
	}
	if !reflect.DeepEqual(runner.calls, wantCalls) {
		t.Fatalf("calls = %#v, want %#v", runner.calls, wantCalls)
	}

	wantMessages := []string{
		"Session already exists: backend-dev",
		"Skipping attach because --no-attach was provided.",
	}
	if !reflect.DeepEqual(result.Messages, wantMessages) {
		t.Fatalf("messages = %#v, want %#v", result.Messages, wantMessages)
	}
}

func testWorkspace(t *testing.T) config.Workspace {
	t.Helper()

	root := t.TempDir()
	return config.Workspace{
		Root: root,
		Env: map[string]string{
			"APP_ENV": "development",
		},
		Windows: []config.Window{
			{Name: "overview", Command: "clear; pwd"},
			{Name: "test-watch", Command: "echo test"},
			{Name: "logs"},
		},
	}
}
