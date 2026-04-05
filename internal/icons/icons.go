package icons

import "strings"

// IconStyle controls which set of icons to use.
type IconStyle int

const (
	Auto     IconStyle = iota // detect from environment
	NerdFont                  // Nerd Font glyphs
	Emoji                     // Unicode emoji
	None                      // no icons
)

type entry struct {
	keys     []string
	nerdFont string
	emoji    string
}

var table = []entry{
	{keys: []string{"node", "nodejs", "npm", "npx"}, nerdFont: "\ue718", emoji: "🟢"},
	{keys: []string{"postgres", "postgresql", "pg"}, nerdFont: "\ue76e", emoji: "🐘"},
	{keys: []string{"docker", "docker-proxy", "dockerd", "containerd"}, nerdFont: "\uf308", emoji: "🐳"},
	{keys: []string{"redis", "redis-server"}, nerdFont: "\ue76d", emoji: "⚡"},
	{keys: []string{"mongod", "mongo", "mongodb"}, nerdFont: "\ue7a4", emoji: "🍃"},
	{keys: []string{"nginx"}, nerdFont: "\ue776", emoji: "🌐"},
	{keys: []string{"python", "python3", "python2", "uvicorn", "gunicorn", "flask", "django"}, nerdFont: "\ue606", emoji: "🐍"},
	{keys: []string{"ruby", "rails", "puma", "unicorn"}, nerdFont: "\ue739", emoji: "💎"},
	{keys: []string{"java", "mvn", "gradle", "spring"}, nerdFont: "\ue738", emoji: "☕"},
	{keys: []string{"go", "air", "gin"}, nerdFont: "\ue724", emoji: "🔵"},
	{keys: []string{"php", "php-fpm"}, nerdFont: "\ue73d", emoji: "🐘"},
	{keys: []string{"mysql", "mysqld"}, nerdFont: "\ue704", emoji: "🐬"},
	{keys: []string{"caddy"}, nerdFont: "\uf0c2", emoji: "☁️"},
	{keys: []string{"deno"}, nerdFont: "\ue71a", emoji: "🦕"},
	{keys: []string{"bun"}, nerdFont: "\ue71a", emoji: "🐾"},
}

const (
	defaultNerdFont = "\uf013" // gear icon
	defaultEmoji    = "⚙️"
)

// Resolve returns the icon for a process name in the given style.
func Resolve(processName string, style IconStyle) string {
	if style == None {
		return ""
	}
	lower := strings.ToLower(processName)
	for _, e := range table {
		for _, key := range e.keys {
			if strings.Contains(lower, key) {
				if style == NerdFont {
					return e.nerdFont
				}
				return e.emoji
			}
		}
	}
	if style == NerdFont {
		return defaultNerdFont
	}
	return defaultEmoji
}
