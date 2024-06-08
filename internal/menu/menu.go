package menu

import (
	"github.com/templwind/sass-starter/internal/config"
	"github.com/templwind/sass-starter/internal/types"
)

type Menus map[string]*config.Menu

func (m Menus) Get(name string) *config.Menu {
	return m[name]
}

// GetContextualMenu adjusts menu items based on user context.
func (m Menus) GetContextualMenu(name string, userContext types.ACLContext) *config.Menu {
	originalMenu := m[name]
	contextualMenu := &config.Menu{
		MenuItems: make([]*config.MenuItem, 0, len(originalMenu.MenuItems)),
	}

	for _, item := range originalMenu.MenuItems {
		if shouldIncludeMenuItem(item, userContext) {
			contextualMenu.MenuItems = append(contextualMenu.MenuItems, item)
		}
	}
	return contextualMenu
}

func shouldIncludeMenuItem(item *config.MenuItem, userContext types.ACLContext) bool {
	// Check if the user's roles match the item's ACL roles
	for _, role := range item.Roles {
		if contains(userContext.Roles, role) {
			return true
		}
	}
	return false
}

func contains(slice []string, item string) bool {
	for _, elem := range slice {
		if elem == item {
			return true
		}
	}
	return false
}
