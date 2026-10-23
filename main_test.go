package main

import (
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
	"testing"

	"github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/tools/router"
)

func TestExpandDefaultArgs(t *testing.T) {
	item := Item{
		Alias:       "gl",
		URL:         "https://gitlab.something/%s/%s/-/merge_requests",
		DefaultArgs: []string{"extranet", "extrota"},
	}
	for _, tc := range []struct {
		query string
		args  []string
		url   string
	}{
		{"gl", []string{"extranet", "extrota"}, "https://gitlab.something/extranet/extrota/-/merge_requests"},
		{"gl   ", []string{"extranet", "extrota"}, "https://gitlab.something/extranet/extrota/-/merge_requests"},
		{"gl other project", []string{"other", "project"}, "https://gitlab.something/other/project/-/merge_requests"},
		{"gl other", []string{"other"}, "https://gitlab.something/other/%s/-/merge_requests"},
	} {
		t.Run(tc.query, func(t *testing.T) {
			got := expand(item, tc.query)
			if got.URL != tc.url || !slices.Equal(got.Args, tc.args) {
				t.Fatalf("got URL %q, args %v; want URL %q, args %v", got.URL, got.Args, tc.url, tc.args)
			}
		})
	}
	search := Item{Alias: "s", URL: "https://example.com/?q=%s"}
	if got := expand(search, "s two words"); got.URL != "https://example.com/?q=two words" {
		t.Fatalf("last explicit argument must retain spaces: %q", got.URL)
	}
	if got := expand(search, "s"); got.URL != search.URL || len(got.Args) != 0 {
		t.Fatalf("an alias without defaults changed: %+v", got)
	}
}

func TestDefaultArgsStorage(t *testing.T) {
	app := pocketbase.NewWithConfig(pocketbase.Config{DefaultDataDir: t.TempDir()})
	t.Cleanup(func() { app.ResetBootstrapState() })
	if err := app.Bootstrap(); err != nil {
		t.Fatal(err)
	}
	// Start with the old schema and a record to verify the migration's backfill.
	var initial core.MigrationsList
	initial.Add(core.AppMigrations.Item(0))
	if _, err := core.NewMigrationsRunner(app, initial).Up(); err != nil {
		t.Fatal(err)
	}
	collection, err := app.FindCollectionByNameOrId("items")
	if err != nil {
		t.Fatal(err)
	}
	legacy := core.NewRecord(collection)
	legacy.Set("alias", "legacy")
	legacy.Set("url", "https://example.com/")
	if err := app.Save(legacy); err != nil {
		t.Fatal(err)
	}
	if err := app.RunAppMigrations(); err != nil {
		t.Fatal(err)
	}
	legacy, err = app.FindRecordById("items", legacy.Id)
	if err != nil {
		t.Fatal(err)
	}
	if got := legacy.GetString("default_args"); got != "[]" {
		t.Fatalf("migrated default_args = %q; want []", got)
	}
	form := url.Values{
		"alias":        {"gl"},
		"url":          {"https://gitlab.something/%s/%s/-/merge_requests"},
		"tags":         {"  git  work \t"},
		"default_args": {" extranet  extrota "},
	}
	req := httptest.NewRequest("POST", "/items", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	item, err := bindItem(&core.RequestEvent{Event: router.Event{Request: req}})
	if err != nil {
		t.Fatal(err)
	}
	if err := createItem(app, item); err != nil {
		t.Fatal(err)
	}
	stored, err := app.FindFirstRecordByData("items", "alias", "gl")
	if err != nil {
		t.Fatal(err)
	}
	if got := stored.GetString("tags"); got != `["git","work"]` {
		t.Fatalf("stored tags = %s; want [\"git\",\"work\"]", got)
	}
	if got := stored.GetString("default_args"); got != `["extranet","extrota"]` {
		t.Fatalf("stored default_args = %s; want [\"extranet\",\"extrota\"]", got)
	}
	for _, query := range []string{"gl", "gl "} {
		got := getItems(app, query).Expansion
		if got.URL != "https://gitlab.something/extranet/extrota/-/merge_requests" || !slices.Equal(got.Args, item.DefaultArgs) {
			t.Fatalf("query %q: unexpected expansion %+v", query, got)
		}
	}
	if err := createItem(app, Item{Alias: "plain", URL: "https://example.com/"}); err != nil {
		t.Fatal(err)
	}
	plain, err := app.FindFirstRecordByData("items", "alias", "plain")
	if err != nil {
		t.Fatal(err)
	}
	if got := plain.GetString("default_args"); got != "[]" {
		t.Fatalf("omitted default_args = %q; want []", got)
	}
}
