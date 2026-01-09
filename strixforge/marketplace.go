package main

import (
	"context"
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/daveweinstein1/strixforge/pkg/containerhub"
)

// MarketplaceModel handles the UI for browsing and installing containers
type MarketplaceModel struct {
	manager   *containerhub.Manager
	list      list.Model
	spinner   spinner.Model
	isLoading bool
	err       error
	width     int
	height    int
	onBack    func()
	onInstall func(containerhub.Image, string) // image, tag

	// Selection state
	selectedImage *containerhub.Image
	selectedTag   string
	view          string // "list", "tags", "installing"
}

// Item wraps marketplace.Image for the list
type item struct {
	image containerhub.Image
}

func (i item) Title() string       { return i.image.Name }
func (i item) Description() string { return i.image.Description }
func (i item) FilterValue() string { return i.image.Name }

func NewMarketplaceModel(mgr *containerhub.Manager, width, height int, backFunc func()) MarketplaceModel {
	l := list.New(nil, list.NewDefaultDelegate(), width, height-4)
	l.Title = "Container Hub"
	l.SetShowHelp(false)
	l.SetFilteringEnabled(true)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return MarketplaceModel{
		manager:   mgr,
		list:      l,
		spinner:   s,
		isLoading: true,
		view:      "list",
		width:     width,
		height:    height,
		onBack:    backFunc,
	}
}

func (m MarketplaceModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.fetchImagesCmd,
	)
}

func (m MarketplaceModel) fetchImagesCmd() tea.Msg {
	ctx := context.Background()
	images, err := m.manager.FetchAllImages(ctx)
	if err != nil {
		return errMsg{err}
	}

	items := make([]list.Item, len(images))
	for i, img := range images {
		items[i] = item{image: img}
	}
	return itemsMsg(items)
}

type itemsMsg []list.Item
type errMsg struct{ err error }

func (m MarketplaceModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.list.FilterState() == list.Filtering {
			break
		}

		switch msg.String() {
		case "q", "esc":
			if m.view == "tags" {
				m.view = "list"
				m.selectedImage = nil
				return m, nil
			}
			if m.onBack != nil {
				m.onBack()
			}
			return m, nil

		case "enter":
			if m.view == "list" {
				if i, ok := m.list.SelectedItem().(item); ok {
					m.selectedImage = &i.image

					// If image has tags pre-loaded, use them
					// Otherwise we'd need to fetch them. For MVP, we assumed pre-fetched or listable.
					// Let's assume fetching on select if empty?
					// For now, simple transition
					m.view = "tags"

					// Trigger tag fetch if needed (TODO)
				}
			} else if m.view == "tags" {
				// Select tag and install
				// TODO: Implement tag selection UI
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.list.SetSize(msg.Width, msg.Height-4)

	case itemsMsg:
		m.isLoading = false
		m.list.SetItems(msg)

	case errMsg:
		m.isLoading = false
		m.err = msg.err
	}

	if m.isLoading {
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	m.list, cmd = m.list.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m MarketplaceModel) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v", m.err)
	}

	if m.isLoading {
		return fmt.Sprintf("\n %s Loading container hub data...\n", m.spinner.View())
	}

	if m.view == "list" {
		return m.list.View()
	}

	if m.view == "tags" && m.selectedImage != nil {
		s := fmt.Sprintf("\n Selected: %s\n\n", m.selectedImage.Name)
		s += fmt.Sprintf(" source: %s\n", m.selectedImage.Source)
		s += fmt.Sprintf(" url: %s\n\n", m.selectedImage.URL)
		s += " Available Tags (detected/static):\n"

		if len(m.selectedImage.Tags) == 0 {
			s += " (Fetching tags not implemented in UI MVP yet)\n"
			s += " Press ENTER to install :latest (or equivalent)\n"
		} else {
			for _, t := range m.selectedImage.Tags {
				s += fmt.Sprintf(" - %s\n", t.Name)
			}
		}

		s += "\n [ESC] Back\n"
		return s
	}

	return "Unknown view"
}
