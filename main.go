package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/kevinburke/ssh_config"
	"github.com/manifoldco/promptui"
)

type NameSorter []Profile

type bellSkipper struct{}

// Write implements an io.WriterCloser over os.Stderr, but it skips the terminal
// bell character.
func (bs *bellSkipper) Write(b []byte) (int, error) {
	const charBell = 7 // c.f. readline.CharBell
	if len(b) == 1 && b[0] == charBell {
		return 0, nil
	}
	return os.Stderr.Write(b)
}

// Close implements an io.WriterCloser over os.Stderr.
func (bs *bellSkipper) Close() error {
	return os.Stderr.Close()
}

func (a NameSorter) Len() int           { return len(a) }
func (a NameSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a NameSorter) Less(i, j int) bool { return a[i].Name < a[j].Name }

type Profile struct {
	Name        string
	Host        string
	User        string
	Key         string
	Port        string
	Title       string
	Description string
}

func main() {
	profiles := getProfilesFromSSHConfig()

	params := os.Args[1:]

	if len(params) == 1 {
		if params[0] == "-h" {
			printHelp()
			os.Exit(0)
		}
		p, e := findProfileByName(params[0], &profiles)
		if e == nil {
			doProfile(p)
			return
		} else {
			doSsh(&params)
			return
		}
	}
	if len(params) > 1 {
		doSsh(&params)
		return
	}

	sort.Sort(NameSorter(profiles))

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   `{{"›" | green}} {{ .Title | green }} ({{ .User | green }}{{"@" | green}}{{ .Host | green }})`,
		Inactive: `{{"›" | white}} {{ .Title | white }} ({{ .User | faint }}{{"@" | faint}}{{ .Host | faint }})`,
		Selected: `{{"›" | green}} {{ .Title | green | faint }}`,
		Details: `
{{ .Title | yellow }}
{{ "Name:" | faint }}	{{ .Name }}
{{ "User:" | faint }}	{{ .User }}
{{ "Host:" | faint }}	{{ .Host }}
{{ "Key:" | faint }}	{{ .Key }}
{{ "Port:" | faint }}	{{ .Port }}
{{ "Description:" | faint }}	{{ .Description }}
`,
	}

	prompt := promptui.Select{
		Label:     "Profiles:",
		Items:     profiles,
		Templates: templates,
		Size:      60,
		Stdout:    &bellSkipper{},
	}

	selected, _, err := prompt.Run()

	if err != nil {
		os.Exit(0)
	}

	doProfile(&profiles[selected])

}

func doProfile(con *Profile) {
	cmd := exec.Command("ssh", con.Name)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func doSsh(params *[]string) {
	cmd := exec.Command("ssh", *params...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Run()
}

func printHelp() {
	cmd := exec.Command("ssh")
	out, _ := cmd.CombinedOutput()
	fmt.Println(string(out))
}

func findProfileByName(name string, profiles *[]Profile) (*Profile, error) {
	for _, p := range *profiles {
		if p.Name == name {
			return &p, nil
		}
	}
	return nil, errors.New("no-profile")
}

func getProfilesFromSSHConfig() []Profile {
	profiles := []Profile{}
	f, _ := os.Open(filepath.Join(os.Getenv("HOME"), ".ssh", "config"))
	cfg, _ := ssh_config.Decode(f)
	for _, host := range cfg.Hosts {
		name := host.Patterns[0].String()
		if name == "" || name == "*" {
			continue
		}
		p := Profile{Name: name, Title: name}
		noQssh := false
		for _, node := range host.Nodes {
			line := strings.TrimSpace(node.String())
			if strings.HasPrefix(line, "#Title ") {
				p.Title = strings.ReplaceAll(line, "#Title ", "")
			}
			if strings.HasPrefix(line, "#Description ") {
				p.Description = strings.ReplaceAll(line, "#Description ", "")
			}
			if strings.HasPrefix(line, "HostName ") {
				p.Host = strings.ReplaceAll(line, "HostName ", "")
			}

			if strings.HasPrefix(line, "User ") {
				p.User = strings.ReplaceAll(line, "User ", "")
			}

			key, _ := cfg.Get(p.Name, "IdentityFile")
			if key == "" {
				key = "~/.ssh/id_rsa"
			}
			p.Key = key

			if strings.HasPrefix(line, "IdentityFile ") {
				p.Key = strings.ReplaceAll(line, "IdentityFile ", "")
				if p.Key == "" {
					p.Key = "adlkajdf"
				}
			}

			port, _ := cfg.Get(p.Name, "Port")
			if port == "" {
				port = "22"
			}
			p.Port = port

			if strings.HasPrefix(line, "#noqssh") {
				noQssh = true
			}

		}
		if !noQssh {
			profiles = append(profiles, p)
		}
	}
	return profiles
}
