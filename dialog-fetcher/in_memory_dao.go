package fetch_dialog

import (
	"fmt"
)

// InMemoryDAO is a fake database that might be used in tests
type InMemoryDAO struct {
	DialogEntries []DialogEntry
	ContextAliasEntries []ContextAliasEntry
}
// NewInMemoryDAO creates an implementation of DAO that is kept in memory
func NewInMemoryDAO() DAO {
	res := InMemoryDAO{DialogEntries: []DialogEntry{}, ContextAliasEntries: []ContextAliasEntry{}}
	return &res
}

// FetchByContextSubject fetches a piece of dialog by context and subject.
func (d InMemoryDAO)FetchByContextSubject(
	context string,
	subject string,
) (rv DialogEntry, err error) {
	for _, de := range d.DialogEntries {
		if de.Context == context && de.Subject == subject {
			return de,nil
		}
	}
	return DialogEntry{}, fmt.Errorf("Not found %s:%s", context, subject)
}

// FetchByDialogID fetches a piece of dialog using a unique UUID associated with the dialog
func (d InMemoryDAO)FetchByDialogID(dialogID string) (result DialogEntry, found bool, err error) {
	for _, de := range d.DialogEntries {
		if de.DialogID == dialogID {
			return de, true, nil
		}
	}
	return DialogEntry{}, false, fmt.Errorf("Not found %s", dialogID)
}

// FetchByAlias fetches a piece of dialog using a application/package ID, context alias, and subject
// https://github.com/adaptiveteam/dialog-library/tree/cultivate/aliases
func (d InMemoryDAO)FetchByAlias(
	packageName,
	contextAlias,
	subject string,
) (rv DialogEntry, err error) {
	applicationAlias := packageName+"#"+contextAlias
	
	for _, alias := range d.ContextAliasEntries {
		if alias.ApplicationAlias == applicationAlias {
			for _, de := range d.DialogEntries {
				if de.Context == alias.Context && de.Subject == subject {
					return de,nil
				}
			}
			return DialogEntry{}, fmt.Errorf("Not found %s:%s", applicationAlias, subject)
		}
	}

	return DialogEntry{}, fmt.Errorf("Not found %s", applicationAlias)
}

// Create s a new item in the dialog table
func (d *InMemoryDAO)Create(dialogEntry DialogEntry) error {
	d.DialogEntries = append(d.DialogEntries, dialogEntry)
	return nil
}
// CreateAlias creates a new alias in the aliases table
func (d *InMemoryDAO)CreateAlias(aliasEntry ContextAliasEntry) error {
	d.ContextAliasEntries = append(d.ContextAliasEntries, aliasEntry)
	return nil
}
