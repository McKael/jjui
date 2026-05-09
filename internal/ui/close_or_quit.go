package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/idursun/jjui/internal/ui/common"
	"github.com/idursun/jjui/internal/ui/intents"
)

// closeTopView closes the topmost closable view or input state (preview pane,
// stacked views, expanded status, flash messages, active operations, etc.).
// It returns the resulting command and true when something was closed, or
// (nil, false) when there was nothing to close. Used by the CloseOrQuit intent
// to decide between closing a view and quitting the application.
func (m *Model) closeTopView() (tea.Cmd, bool) {
	// Password prompt takes highest priority.
	if m.password != nil {
		return m.password.Update(intents.Cancel{}), true
	}
	// Status input mode.
	if m.status.IsFocused() {
		return m.status.Update(intents.Cancel{}), true
	}
	// Revset editing mode.
	if m.revsetModel.IsEditing() {
		return m.revsetModel.Update(intents.Cancel{}), true
	}
	// Stacked views (help, bookmarks, git, undo, redo, etc.).
	if m.stacked != nil || m.diff != nil || m.oplog != nil {
		return common.Close, true
	}
	// Preview panel.
	if m.splitContainer != nil && m.splitContainer.IsContentActive(previewContentID) {
		m.splitContainer.Close()
		return nil, true
	}
	// Flash messages.
	if m.flash.Any() {
		m.flash.DeleteOldest()
		return nil, true
	}
	// Expanded status.
	if m.status.StatusExpanded() {
		m.status.ToggleStatusExpand()
		return nil, true
	}
	// Operations in revisions (squash, rebase, etc.) - route cancel to revisions.
	if !m.revisions.InNormalMode() {
		return m.revisions.Update(intents.Cancel{}), true
	}
	// Quick search.
	if m.revisions.HasQuickSearch() {
		return m.revisions.Update(intents.Cancel{}), true
	}
	return nil, false
}