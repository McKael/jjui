package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/idursun/jjui/internal/ui/actions"
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
	// Revset editing mode.
	if m.revsetModel.IsEditing() {
		return m.revsetModel.Update(intents.Cancel{}), true
	}
	// The expanded status help is an overlay over the active view; close it
	// before any view beneath it, unless a dialog is on top (mirrors
	// closeTopScope).
	if m.dialog() == nil && m.status.StatusExpanded() {
		m.status.SetStatusExpanded(false)
		return nil, true
	}
	// The flash (error popup) renders above views (ZOverlay), so it is
	// dismissed before them — except while the command history is open, where
	// View() hides the flash line and the history must close first.
	topScope, _ := m.stackedScope()
	if m.flash.Any() && topScope != actions.ScopeCommandHistory {
		m.flash.DeleteOldest()
		return nil, true
	}
	// Root views and dialogs are tracked in the ordered layer list.
	if len(m.layers) > 0 {
		return common.Close, true
	}
	// Split content (preview pane, bookmark pane, etc.).
	if m.splitContainer != nil && m.splitContainer.Close() {
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
