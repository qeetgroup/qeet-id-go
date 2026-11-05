package qeetid

import (
	"encoding/json"
	"errors"
)

// UpdateGroupParams is the request body of GroupsService.Update. Fields left
// nil are not changed.
type UpdateGroupParams struct {
	// Name renames the group.
	Name *string
	// Description replaces the group's description.
	Description *string
	// ParentID moves the group under another group of the tenant.
	ParentID *string
	// MoveToTop moves the group to the top level, removing its parent. It
	// cannot be combined with ParentID.
	MoveToTop bool
}

// MarshalJSON encodes the parameters the way the API expects: parent_id is
// left out to keep the parent, null to remove it, or the new parent's ID.
func (p UpdateGroupParams) MarshalJSON() ([]byte, error) {
	body := map[string]any{}
	if p.Name != nil {
		body["name"] = *p.Name
	}
	if p.Description != nil {
		body["description"] = *p.Description
	}
	switch {
	case p.MoveToTop && p.ParentID != nil:
		return nil, errors.New("qeetid: UpdateGroupParams: set ParentID or MoveToTop, not both")
	case p.MoveToTop:
		body["parent_id"] = nil
	case p.ParentID != nil:
		body["parent_id"] = *p.ParentID
	}
	return json.Marshal(body)
}
