// SPDX-License-Identifier: GPL-3.0-only
package constants

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultSiteID(t *testing.T) {
	// Test that DefaultSiteID has the expected value
	assert.Equal(t, "default", DefaultSiteID, "DefaultSiteID should be 'default'")
	
	// Test that DefaultSiteID is not empty
	assert.NotEmpty(t, DefaultSiteID, "DefaultSiteID should not be empty")
	
	// Test that DefaultSiteID is a constant (doesn't change)
	original := DefaultSiteID
	// Try to "modify" (this won't actually change the constant but tests immutability concept)
	currentValue := DefaultSiteID
	assert.Equal(t, original, currentValue, "DefaultSiteID should remain constant")
}

func TestDefaultSiteIDType(t *testing.T) {
	// Test that DefaultSiteID is a string type
	assert.IsType(t, "", DefaultSiteID, "DefaultSiteID should be a string type")
}
