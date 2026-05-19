package manager

import (
	"bufio"
	"os/exec"
	"strings"
	"time"

	"github.com/neur0map/glazepkg/internal/model"
)

type Mise struct{}

func (y *Mise) Name() model.Source { return model.SourceMise }

func (y *Mise) Available() bool { return commandExists("mise") }

func (y *Mise) Scan() ([]model.Package, error) {
	out, err := exec.Command("mise", "ls").Output()
	if err != nil {
		return nil, err
	}

	var pkgs []model.Package
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}
		pkgs = append(pkgs, model.Package{
			Name:        fields[0],
			Version:     fields[1],
			Source:      model.SourceMise,
			InstalledAt: time.Now(),
		})
	}
	return pkgs, nil
}

func (y *Mise) Search(query string) ([]model.Package, error) {
	out, err := exec.Command("mise", "search", query).Output()
	if err != nil || len(out) == 0 {
		return nil, nil
	}

	var pkgs []model.Package
	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "mise ERROR") {
			break
		}
		if strings.Contains(line, "Tool") && strings.Contains(line, "Description") {
			continue
		}
		tokens := strings.Split(line, "  ")
		name := strings.TrimSpace(tokens[0])
		description := strings.TrimSpace(tokens[len(tokens)-1])

		pkgs = append(pkgs, model.Package{Name: name, Description: description, Source: model.SourceMise})
	}

	return pkgs, nil
}

func (y *Mise) CheckUpdates(pkgs []model.Package) map[string]string {
	// return map of name → latest version
	// novotarq@burza:~/work/go/glazepkg|⇒  mise outdated yt-dlp node
	// name    requested  current     latest     source
	// node    latest     25.9.0      26.1.0     ~/work/go/glazepkg/mise.toml
	// yt-dlp  latest     2025.12.08  2026.03.17 ~/work/go/glazepkg/mise.toml
	return make(map[string]string)
}

func (y *Mise) InstallCmd(name string) *exec.Cmd {
	return exec.Command("mise", "use", name)
}

func (y *Mise) RemoveCmd(name string) *exec.Cmd {
	return exec.Command("mise", "uninstall", name)
}

func (y *Mise) UpgradeCmd(name string) *exec.Cmd {
	return exec.Command("mise", "upgrade", name)
}
