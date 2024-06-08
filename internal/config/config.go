package config

import (
	"embed"
	"sort"

	"github.com/biter777/countries"
	"github.com/goccy/go-yaml"
)

type Config struct {
	AppName            string   `yaml:"AppName"`
	DefaultDataDir     string   `yaml:"DefaultDataDir"`
	DatabaseFileName   string   `yaml:"DatabaseFileName"`
	RunMigrations      bool     `yaml:"RunMigrations"`
	MigrationsPath     string   `yaml:"MigrationsPath"`
	EmbeddedMigrations embed.FS `yaml:"-"`
	Auth               struct {
		TokenSecret   string `yaml:"TokenSecret"`
		TokenDuration string `yaml:"TokenDuration"`
	} `yaml:"Auth"`
	Setup struct {
		DefaultAdmin struct {
			Name     string `yaml:"Name"`
			Username string `yaml:"Username"`
			Email    string `yaml:"Email"`
			Password string `yaml:"Password"`
		} `yaml:"DefaultAdmin"`
		TestUser struct {
			Name     string `yaml:"Name"`
			Username string `yaml:"Username"`
			Email    string `yaml:"Email"`
			Password string `yaml:"Password"`
		} `yaml:"TestUser"`
	} `yaml:"Setup"`
	RoleBindings     map[string][]string `yaml:"RoleBindings"`
	DefaultRole      string              `yaml:"DefaultRole"`
	Menus            map[string]*Menu    `yaml:"Menus"`
	AllowedCountries map[string]bool     `yaml:"AllowedCountries,omitempty"`
	CountryCodeList  map[string]string   `yaml:"-"`
}

// Menu represents a collection of menu items.
type Menu struct {
	MenuItems []*MenuItem `yaml:"MenuItems"`
}

// MenuItem represents a single item in a menu.
type MenuItem struct {
	Icon        string      `yaml:"Icon"`
	Title       string      `yaml:"Title"`
	Subtitle    string      `yaml:"Subtitle"`
	MobileTitle string      `yaml:"MobileTitle"`
	Lead        string      `yaml:"Lead"`
	Link        string      `yaml:"Link"`
	IsAtEnd     bool        `yaml:"IsAtEnd"`
	IsDropdown  bool        `yaml:"IsDropdown"`
	IsHtmx      bool        `yaml:"IsHtmx"`
	SubItems    []*MenuItem `yaml:"SubItems"`
	Roles       []string    `yaml:"Roles"`
}

// Prepend adds an item to the beginning of the menu.
func (m *Menu) Prepend(item *MenuItem) {
	m.MenuItems = append([]*MenuItem{item}, m.MenuItems...)
}

// Append adds an item to the end of the menu.
func (m *Menu) Append(item *MenuItem) {
	m.MenuItems = append(m.MenuItems, item)
}

// InsertAt inserts an item at a specific index in the menu.
func (m *Menu) InsertAt(index int, item *MenuItem) {
	m.MenuItems = append(m.MenuItems[:index], append([]*MenuItem{item}, m.MenuItems[index:]...)...)
}

// GetMenuItems returns the items of the specified menu.
func (m *Menu) GetMenuItems() []*MenuItem {
	return m.MenuItems
}

// GetSubMenuItemsByURL returns the sub-items of a specific menu item based on the URL.
func (m *Menu) GetSubMenuItemsByURL(url string) []*MenuItem {
	for _, item := range m.MenuItems {
		if item.Link == url {
			return item.SubItems
		}
	}
	return []*MenuItem{}
}

// LoadConfigFromYamlBytes loads configuration from YAML bytes.
func LoadConfigFromYamlBytes(bytes []byte, cfg *Config) error {
	err := yaml.Unmarshal(bytes, cfg)
	if err != nil {
		return err
	}

	// Load country codes
	loadCountryCodes(cfg)

	return nil
}

func loadCountryCodes(cfg *Config) {
	// Initialize CountryCodeList
	cfg.CountryCodeList = make(map[string]string)

	allowed := cfg.AllowedCountries
	allCountries := countries.All()

	// Filter and populate CountryCodeList
	for _, country := range allCountries {
		alpha2 := country.Alpha2()
		if allowed[alpha2] {
			cfg.CountryCodeList[alpha2] = country.Info().Name
		}
	}

	// Convert map to a slice for sorting
	sortedCountries := make([]struct {
		Code string
		Name string
	}, 0, len(cfg.CountryCodeList))

	for code, name := range cfg.CountryCodeList {
		sortedCountries = append(sortedCountries, struct {
			Code string
			Name string
		}{Code: code, Name: name})
	}

	// Sort the slice by country name
	sort.Slice(sortedCountries, func(i, j int) bool {
		return sortedCountries[i].Name < sortedCountries[j].Name
	})

	// Clear and repopulate the map in sorted order
	cfg.CountryCodeList = make(map[string]string)
	for _, country := range sortedCountries {
		cfg.CountryCodeList[country.Code] = country.Name
	}
}
