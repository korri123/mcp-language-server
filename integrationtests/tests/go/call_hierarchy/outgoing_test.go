package callhierarchy_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/isaacphi/mcp-language-server/integrationtests/tests/common"
	"github.com/isaacphi/mcp-language-server/integrationtests/tests/go/internal"
	"github.com/isaacphi/mcp-language-server/internal/tools"
)

func TestOutgoingCalls(t *testing.T) {
	suite := internal.GetTestSuite(t)

	ctx, cancel := context.WithTimeout(suite.Context, 10*time.Second)
	defer cancel()

	tests := []struct {
		name         string
		symbolName   string
		expectedText string
		snapshotName string
	}{
		{
			name:         "Function with calls in other file",
			symbolName:   "ConsumerFunction",
			expectedText: "HelperFunction",
			snapshotName: "outgoing-other-file",
		},
		{
			name:         "Function with calls in same file",
			symbolName:   "main",
			expectedText: "FooBar",
			snapshotName: "outgoing-same-file",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Call the GetOutgoingCalls tool
			result, err := tools.GetCallees(ctx, suite.Client, tc.symbolName, 1)
			if err != nil {
				t.Fatalf("Failed to find outgoing calls: %v", err)
			}

			// Check that the result contains relevant information
			if !strings.Contains(result, tc.expectedText) {
				t.Errorf("Outgoing calls do not contain expected text: %s", tc.expectedText)
			}

			// Use snapshot testing to verify exact output
			common.SnapshotTest(t, "go", "call_hierarchy", tc.snapshotName, result)
		})
	}

}
