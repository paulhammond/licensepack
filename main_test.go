package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var modules = ModuleSet{
	Modules: []Module{
		{
			Name:     "aa1",
			Licenses: []File{{Path: "a", Contents: "a"}},
		},
		{
			Name:     "aa2",
			Licenses: []File{{Path: "a", Contents: "a"}},
		},
		{
			Name:     "bb1",
			Licenses: []File{{Path: "b", Contents: "b"}},
		},
		{
			Name:     "bb2",
			Licenses: []File{{Path: "b", Contents: "b"}},
		},
		{
			Name:     "a.txt",
			Licenses: []File{{Path: "a.txt", Contents: "a"}},
		},
		{
			Name:     "A",
			Licenses: []File{{Path: "A", Contents: "a"}},
		},
		{
			Name:     "c",
			Licenses: []File{{Path: "c", Contents: "c"}},
		},
		{
			Name:     "d",
			Licenses: []File{{Path: "d", Contents: "d"}},
		},
		{
			Name: "cd",
			Licenses: []File{
				{Path: "c", Contents: "c"},
				{Path: "d", Contents: "d"},
			},
		},
	},
}

func TestFileGroups(t *testing.T) {
	g := modules.FileGroups()
	want := []*Group{
		{
			Names: []string{"aa1 a", "aa2 a", "a.txt a.txt", "A A"},
			Licenses: []File{
				{Path: "a", Contents: "a"},
			},
		},
		{
			Names: []string{"bb1 b", "bb2 b"},
			Licenses: []File{
				{Path: "b", Contents: "b"},
			},
		},
		{
			Names: []string{"c c", "cd c"},
			Licenses: []File{
				{Path: "c", Contents: "c"},
			},
		},
		{
			Names: []string{"d d", "cd d"},
			Licenses: []File{
				{Path: "d", Contents: "d"},
			},
		},
	}
	if diff := cmp.Diff(want, g); diff != "" {
		t.Errorf("FileGroups() (-want +got):\n%s", diff)
	}
}

func TestModuleGroups(t *testing.T) {
	g := modules.ModuleGroups()
	want := []*Group{
		{
			Names: []string{"aa1", "aa2"},
			Licenses: []File{
				{Path: "a", Contents: "a"},
			},
		},
		{
			Names: []string{"bb1", "bb2"},
			Licenses: []File{
				{Path: "b", Contents: "b"},
			},
		},
		{
			Names: []string{"a.txt"},
			Licenses: []File{
				{Path: "a.txt", Contents: "a"},
			},
		},
		{
			Names: []string{"A"},
			Licenses: []File{
				{Path: "A", Contents: "a"},
			},
		},
		{
			Names: []string{"c"},
			Licenses: []File{
				{Path: "c", Contents: "c"},
			},
		},
		{
			Names: []string{"d"},
			Licenses: []File{
				{Path: "d", Contents: "d"},
			},
		},
		{
			Names: []string{"cd"},
			Licenses: []File{
				{Path: "c", Contents: "c"},
				{Path: "d", Contents: "d"},
			},
		},
	}
	if diff := cmp.Diff(want, g); diff != "" {
		t.Errorf("FileGroups() (-want +got):\n%s", diff)
	}
}
