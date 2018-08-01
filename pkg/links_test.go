package pkg

import (
	"reflect"
	"testing"
)

func TestLinkTree(t *testing.T) {
	links := []string{
		"/",
		"/a/",
		"/b/",
		"/a/c/",
		"/b/d",
		"/a/c/d",
		"/b/d/f",
	}

	tree := &LinkTree{}
	tree.AddLinks(links...)
	expected := &LinkTree{
		Value: "/",
		Childs: []*LinkTree{
			&LinkTree{
				Value: "/a/",
				Childs: []*LinkTree{
					&LinkTree{
						Value: "/a/c/",
						Childs: []*LinkTree{
							&LinkTree{
								Value:  "/a/c/d",
								Childs: []*LinkTree{},
							},
						},
					},
				},
			},
			&LinkTree{
				Value: "/b/",
				Childs: []*LinkTree{
					&LinkTree{
						Value: "/b/d",
						Childs: []*LinkTree{
							&LinkTree{
								Value:  "/b/d/f",
								Childs: []*LinkTree{},
							},
						},
					},
				},
			},
		},
	}
	if !reflect.DeepEqual(tree, expected) {
		t.Errorf("Trees are not equal: %v != %v", tree, expected)
	}

	expectedExport := []string{
		"/", "/a/", "/a/c/", "/a/c/d",
		"/b/", "/b/d", "/b/d/f",
	}
	resultExport := tree.GetLinks()
	if !reflect.DeepEqual(resultExport, expectedExport) {
		t.Errorf("Trees exports are not equal: %v != %v", resultExport, expectedExport)
	}
}
