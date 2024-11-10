package pagination

import "testing"

func TestComputePages(t *testing.T) {
	testCases := []struct {
		maxPages  int
		offset    int
		pageSize  int
		totalSize int
		expected  []Page
	}{
		{5, 0, 5, 5, []Page{{1, 0, 5, true, RegularPage}}},
		{5, 0, 5, 6, []Page{{1, 0, 5, true, RegularPage}, {2, 5, 1, false, RegularPage}}},
		{5, 0, 5, 10, []Page{{1, 0, 5, true, RegularPage}, {2, 5, 5, false, RegularPage}}},
	}

	for _, tc := range testCases {
		actual := ComputePages(tc.maxPages, tc.offset, tc.pageSize, tc.totalSize)
		for i := range tc.expected {
			if len(tc.expected) != len(actual) {
				t.Errorf("unequal page sizes, want %d, got %d", len(tc.expected), len(actual))
			}
			if tc.expected[i].Active != actual[i].Active ||
				tc.expected[i].Offset != actual[i].Offset ||
				tc.expected[i].Size != actual[i].Size ||
				tc.expected[i].Number != actual[i].Number ||
				tc.expected[i].ItemType != actual[i].ItemType {
				t.Errorf("invalid pagination result, want %+v, got %+v", tc.expected[i], actual[i])
			}
		}
	}
}
