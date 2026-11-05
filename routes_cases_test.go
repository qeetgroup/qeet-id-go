package qeetid

import (
	"context"
)

// routeCase is one API call of the client, with the request it must send.
type routeCase struct {
	name       string
	method     string
	path       string
	public     bool
	idempotent bool
	call       func(ctx context.Context, c *Client) error
}

// routeCases covers every generated method: each must send the documented
// method and path, authenticate unless the endpoint is public, and send an
// Idempotency-Key exactly when the endpoint accepts one.
var routeCases = []routeCase{
	{"Tenants.List", "GET", "/v1/tenants", false, false, func(ctx context.Context, c *Client) error { _, err := c.Tenants.List(ctx, nil); return err }},
	{"Tenants.All", "GET", "/v1/tenants", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.Tenants.All(ctx, nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"Tenants.Create", "POST", "/v1/tenants", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Tenants.Create(ctx, CreateTenantParams{})
		return err
	}},
	{"Tenants.Get", "GET", "/v1/tenants/tenantID-1", false, false, func(ctx context.Context, c *Client) error { _, err := c.Tenants.Get(ctx, "tenantID-1"); return err }},
	{"Tenants.Update", "PATCH", "/v1/tenants/tenantID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Tenants.Update(ctx, "tenantID-1", UpdateTenantParams{})
		return err
	}},
	{"Tenants.Delete", "DELETE", "/v1/tenants/tenantID-1", false, false, func(ctx context.Context, c *Client) error { _, err := c.Tenants.Delete(ctx, "tenantID-1"); return err }},
	{"Tenants.CancelDeletion", "POST", "/v1/tenants/tenantID-1/deletion/cancel", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Tenants.CancelDeletion(ctx, "tenantID-1")
		return err
	}},
	{"Users.Create", "POST", "/v1/users", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Users.Create(ctx, CreateUserParams{})
		return err
	}},
	{"Users.Get", "GET", "/v1/users/userID-1", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Get(ctx, "userID-1"); return err }},
	{"Users.Update", "PATCH", "/v1/users/userID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.Update(ctx, "userID-1", UpdateUserParams{})
		return err
	}},
	{"Users.Delete", "DELETE", "/v1/users/userID-1", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Delete(ctx, "userID-1"); return err }},
	{"Users.List", "GET", "/v1/users", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.List(ctx, nil); return err }},
	{"Users.All", "GET", "/v1/users", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.Users.All(ctx, nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"Users.SetPassword", "POST", "/v1/users/userID-1/password", false, false, func(ctx context.Context, c *Client) error {
		return c.Users.SetPassword(ctx, "userID-1", SetPasswordParams{})
	}},
	{"Users.BulkCreate", "POST", "/v1/users/bulk", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.BulkCreate(ctx, BulkCreateUsersParams{})
		return err
	}},
	{"Users.ListDeleted", "GET", "/v1/users/deleted", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.ListDeleted(ctx, nil); return err }},
	{"Users.Restore", "POST", "/v1/users/userID-1/restore", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Restore(ctx, "userID-1"); return err }},
	{"Users.Purge", "DELETE", "/v1/users/userID-1/purge", false, false, func(ctx context.Context, c *Client) error { return c.Users.Purge(ctx, "userID-1") }},
	{"Users.RemoveMembership", "DELETE", "/v1/users/userID-1/membership", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.RemoveMembership(ctx, "userID-1")
		return err
	}},
	{"Users.ResetMFA", "DELETE", "/v1/users/userID-1/mfa", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.ResetMFA(ctx, "userID-1"); return err }},
	{"Users.SetMFARequired", "PUT", "/v1/users/userID-1/mfa-required", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.SetMFARequired(ctx, "userID-1", SetMFARequiredParams{})
		return err
	}},
	{"Users.Security", "GET", "/v1/users/userID-1/security", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Security(ctx, "userID-1"); return err }},
	{"Users.Access", "GET", "/v1/users/userID-1/access", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Access(ctx, "userID-1"); return err }},
	{"Users.Stats", "GET", "/v1/users/stats", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Stats(ctx); return err }},
	{"Users.Trends", "GET", "/v1/users/trends", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Trends(ctx); return err }},
	{"Users.ListSuspended", "GET", "/v1/users/suspended", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.ListSuspended(ctx, nil); return err }},
	{"Users.Reactivate", "POST", "/v1/users/userID-1/reactivate", false, false, func(ctx context.Context, c *Client) error { _, err := c.Users.Reactivate(ctx, "userID-1"); return err }},
	{"Users.ReviewSuspension", "PUT", "/v1/users/userID-1/suspension", false, false, func(ctx context.Context, c *Client) error {
		return c.Users.ReviewSuspension(ctx, "userID-1", ReviewSuspensionParams{})
	}},
	{"Users.ListSessions", "GET", "/v1/users/userID-1/sessions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListSessions(ctx, "userID-1")
		return err
	}},
	{"Users.RevokeSessions", "POST", "/v1/users/userID-1/sessions/revoke-all", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.RevokeSessions(ctx, "userID-1")
		return err
	}},
	{"Users.EndSession", "DELETE", "/v1/tenants/tenantID-1/users/userID-1/sessions/sessionID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Users.EndSession(ctx, "tenantID-1", "userID-1", "sessionID-1")
	}},
	{"Users.ListTrustedDevices", "GET", "/v1/users/userID-1/trusted-devices", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListTrustedDevices(ctx, "userID-1")
		return err
	}},
	{"Users.ForgetTrustedDevices", "DELETE", "/v1/users/userID-1/trusted-devices", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ForgetTrustedDevices(ctx, "userID-1")
		return err
	}},
	{"Users.ForgetTrustedDevice", "DELETE", "/v1/users/userID-1/trusted-devices/deviceID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ForgetTrustedDevice(ctx, "userID-1", "deviceID-1")
		return err
	}},
	{"Users.ListPushDevices", "GET", "/v1/users/userID-1/push-devices", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListPushDevices(ctx, "userID-1")
		return err
	}},
	{"Users.DisconnectPushDevice", "DELETE", "/v1/users/userID-1/push-devices/deviceID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Users.DisconnectPushDevice(ctx, "userID-1", "deviceID-1")
	}},
	{"Users.ListPasskeys", "GET", "/v1/tenants/tenantID-1/users/userID-1/passkeys", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListPasskeys(ctx, "tenantID-1", "userID-1")
		return err
	}},
	{"Users.DeletePasskey", "DELETE", "/v1/tenants/tenantID-1/users/userID-1/passkeys/passkeyID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Users.DeletePasskey(ctx, "tenantID-1", "userID-1", "passkeyID-1")
	}},
	{"Users.SignIns", "GET", "/v1/tenants/tenantID-1/users/userID-1/sign-ins", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.SignIns(ctx, "tenantID-1", "userID-1", nil)
		return err
	}},
	{"Users.StartEmailVerification", "POST", "/v1/users/userID-1/verify/email/start", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.StartEmailVerification(ctx, "userID-1", StartEmailVerificationParams{})
		return err
	}},
	{"Users.ConfirmEmailVerification", "POST", "/v1/users/userID-1/verify/email/confirm", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ConfirmEmailVerification(ctx, "userID-1", ConfirmVerificationParams{})
		return err
	}},
	{"Users.StartPhoneVerification", "POST", "/v1/users/userID-1/verify/phone/start", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.StartPhoneVerification(ctx, "userID-1", StartPhoneVerificationParams{})
		return err
	}},
	{"Users.ConfirmPhoneVerification", "POST", "/v1/users/userID-1/verify/phone/confirm", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ConfirmPhoneVerification(ctx, "userID-1", ConfirmVerificationParams{})
		return err
	}},
	{"Users.ListSocialIdentities", "GET", "/v1/users/userID-1/social/identities", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListSocialIdentities(ctx, "userID-1")
		return err
	}},
	{"Groups.Create", "POST", "/v1/groups", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.Create(ctx, CreateGroupParams{})
		return err
	}},
	{"Groups.Update", "PATCH", "/v1/groups/groupID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.Update(ctx, "groupID-1", UpdateGroupParams{})
		return err
	}},
	{"Groups.Delete", "DELETE", "/v1/groups/groupID-1", false, false, func(ctx context.Context, c *Client) error { return c.Groups.Delete(ctx, "groupID-1") }},
	{"Groups.List", "GET", "/v1/tenants/tenantID-1/groups", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.List(ctx, "tenantID-1", nil)
		return err
	}},
	{"Groups.ListMembers", "GET", "/v1/groups/groupID-1/members", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.ListMembers(ctx, "groupID-1")
		return err
	}},
	{"Groups.AddMember", "POST", "/v1/groups/groupID-1/members/userID-1", false, false, func(ctx context.Context, c *Client) error { return c.Groups.AddMember(ctx, "groupID-1", "userID-1") }},
	{"Groups.RemoveMember", "DELETE", "/v1/groups/groupID-1/members/userID-1", false, false, func(ctx context.Context, c *Client) error { return c.Groups.RemoveMember(ctx, "groupID-1", "userID-1") }},
	{"Groups.ListRoles", "GET", "/v1/tenants/tenantID-1/groups/groupID-1/roles", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.ListRoles(ctx, "tenantID-1", "groupID-1")
		return err
	}},
	{"Groups.GrantRole", "POST", "/v1/tenants/tenantID-1/groups/groupID-1/roles/roleID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Groups.GrantRole(ctx, "tenantID-1", "groupID-1", "roleID-1")
	}},
	{"Groups.RevokeRole", "DELETE", "/v1/tenants/tenantID-1/groups/groupID-1/roles/roleID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Groups.RevokeRole(ctx, "tenantID-1", "groupID-1", "roleID-1")
	}},
	{"Groups.Import", "POST", "/v1/groups/import", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.Import(ctx, ImportGroupsParams{})
		return err
	}},
	{"Groups.ListTemplates", "GET", "/v1/groups/templates", false, false, func(ctx context.Context, c *Client) error { _, err := c.Groups.ListTemplates(ctx); return err }},
	{"Groups.CreateTemplate", "POST", "/v1/groups/templates", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.CreateTemplate(ctx, CreateGroupTemplateParams{})
		return err
	}},
	{"Groups.DeleteTemplate", "DELETE", "/v1/groups/templates/templateID-1", false, false, func(ctx context.Context, c *Client) error { return c.Groups.DeleteTemplate(ctx, "templateID-1") }},
	{"Groups.ApplyTemplate", "POST", "/v1/groups/templates/templateID-1/apply", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Groups.ApplyTemplate(ctx, "templateID-1", ApplyGroupTemplateParams{})
		return err
	}},
	{"Invitations.Create", "POST", "/v1/invites", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Invitations.Create(ctx, CreateInvitationParams{})
		return err
	}},
	{"Invitations.List", "GET", "/v1/tenants/tenantID-1/invites", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Invitations.List(ctx, "tenantID-1")
		return err
	}},
	{"Invitations.Revoke", "DELETE", "/v1/invites/invitationID-1", false, false, func(ctx context.Context, c *Client) error { return c.Invitations.Revoke(ctx, "invitationID-1") }},
	{"Invitations.Resend", "POST", "/v1/invites/invitationID-1/resend", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Invitations.Resend(ctx, "invitationID-1")
		return err
	}},
	{"Domains.List", "GET", "/v1/tenants/tenantID-1/domains", false, false, func(ctx context.Context, c *Client) error { _, err := c.Domains.List(ctx, "tenantID-1"); return err }},
	{"Domains.Create", "POST", "/v1/tenants/tenantID-1/domains", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Create(ctx, "tenantID-1", CreateDomainParams{})
		return err
	}},
	{"Domains.Verify", "POST", "/v1/tenants/tenantID-1/domains/domainID-1/verify", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Verify(ctx, "tenantID-1", "domainID-1")
		return err
	}},
	{"Domains.Update", "PATCH", "/v1/tenants/tenantID-1/domains/domainID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Update(ctx, "tenantID-1", "domainID-1", UpdateDomainParams{})
		return err
	}},
	{"Domains.Delete", "DELETE", "/v1/tenants/tenantID-1/domains/domainID-1", false, false, func(ctx context.Context, c *Client) error { return c.Domains.Delete(ctx, "tenantID-1", "domainID-1") }},
	{"Domains.LoginDomain", "GET", "/v1/tenants/tenantID-1/login-domain", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Domains.LoginDomain(ctx, "tenantID-1")
		return err
	}},
	{"Domains.CheckLoginDomain", "POST", "/v1/tenants/tenantID-1/login-domain/check", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Domains.CheckLoginDomain(ctx, "tenantID-1")
		return err
	}},
	{"Domains.ListBreakGlassMembers", "GET", "/v1/tenants/tenantID-1/break-glass-members", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Domains.ListBreakGlassMembers(ctx, "tenantID-1")
		return err
	}},
	{"Domains.AddBreakGlassMember", "PUT", "/v1/tenants/tenantID-1/break-glass-members/userID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Domains.AddBreakGlassMember(ctx, "tenantID-1", "userID-1")
	}},
	{"Domains.RemoveBreakGlassMember", "DELETE", "/v1/tenants/tenantID-1/break-glass-members/userID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Domains.RemoveBreakGlassMember(ctx, "tenantID-1", "userID-1")
	}},
	{"Branding.Get", "GET", "/v1/tenants/tenantID-1/branding", false, false, func(ctx context.Context, c *Client) error { _, err := c.Branding.Get(ctx, "tenantID-1"); return err }},
	{"Branding.Update", "PUT", "/v1/tenants/tenantID-1/branding", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Branding.Update(ctx, "tenantID-1", Branding{})
		return err
	}},
	{"Roles.List", "GET", "/v1/tenants/tenantID-1/roles", false, false, func(ctx context.Context, c *Client) error { _, err := c.Roles.List(ctx, "tenantID-1"); return err }},
	{"Roles.Create", "POST", "/v1/tenants/tenantID-1/roles", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Roles.Create(ctx, "tenantID-1", CreateRoleParams{})
		return err
	}},
	{"Roles.Update", "PATCH", "/v1/tenants/tenantID-1/roles/roleID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Roles.Update(ctx, "tenantID-1", "roleID-1", UpdateRoleParams{})
		return err
	}},
	{"Roles.Delete", "DELETE", "/v1/tenants/tenantID-1/roles/roleID-1", false, false, func(ctx context.Context, c *Client) error { return c.Roles.Delete(ctx, "tenantID-1", "roleID-1") }},
	{"Roles.ListPermissions", "GET", "/v1/roles/roleID-1/permissions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Roles.ListPermissions(ctx, "roleID-1")
		return err
	}},
	{"Roles.SetPermissions", "PUT", "/v1/tenants/tenantID-1/roles/roleID-1/permissions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Roles.SetPermissions(ctx, "tenantID-1", "roleID-1", SetRolePermissionsParams{})
		return err
	}},
	{"Roles.GrantPermission", "POST", "/v1/roles/roleID-1/permissions/permissionID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Roles.GrantPermission(ctx, "roleID-1", "permissionID-1")
	}},
	{"Roles.RevokePermission", "DELETE", "/v1/roles/roleID-1/permissions/permissionID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Roles.RevokePermission(ctx, "roleID-1", "permissionID-1")
	}},
	{"Roles.Holders", "GET", "/v1/tenants/tenantID-1/roles/roleID-1/holders", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Roles.Holders(ctx, "tenantID-1", "roleID-1", nil)
		return err
	}},
	{"Roles.AssignToUser", "POST", "/v1/users/userID-1/tenants/tenantID-1/roles/roleID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Roles.AssignToUser(ctx, "userID-1", "tenantID-1", "roleID-1")
	}},
	{"Roles.UnassignFromUser", "DELETE", "/v1/users/userID-1/tenants/tenantID-1/roles/roleID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Roles.UnassignFromUser(ctx, "userID-1", "tenantID-1", "roleID-1")
	}},
	{"Permissions.List", "GET", "/v1/permissions", false, false, func(ctx context.Context, c *Client) error { _, err := c.Permissions.List(ctx); return err }},
	{"Permissions.Check", "GET", "/v1/check", false, false, func(ctx context.Context, c *Client) error { _, err := c.Permissions.Check(ctx, nil); return err }},
	{"Permissions.Effective", "GET", "/v1/users/userID-1/tenants/tenantID-1/permissions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Permissions.Effective(ctx, "userID-1", "tenantID-1")
		return err
	}},
	{"ABACPolicies.List", "GET", "/v1/tenants/tenantID-1/abac/policies", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.List(ctx, "tenantID-1")
		return err
	}},
	{"ABACPolicies.Create", "POST", "/v1/tenants/tenantID-1/abac/policies", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Create(ctx, "tenantID-1", ABACPolicyParams{})
		return err
	}},
	{"ABACPolicies.Get", "GET", "/v1/tenants/tenantID-1/abac/policies/policyID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Get(ctx, "tenantID-1", "policyID-1")
		return err
	}},
	{"ABACPolicies.Update", "PATCH", "/v1/tenants/tenantID-1/abac/policies/policyID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Update(ctx, "tenantID-1", "policyID-1", ABACPolicyParams{})
		return err
	}},
	{"ABACPolicies.Delete", "DELETE", "/v1/tenants/tenantID-1/abac/policies/policyID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.ABACPolicies.Delete(ctx, "tenantID-1", "policyID-1")
	}},
	{"ABACPolicies.ListVersions", "GET", "/v1/tenants/tenantID-1/abac/policies/policyID-1/versions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.ListVersions(ctx, "tenantID-1", "policyID-1")
		return err
	}},
	{"ABACPolicies.RestoreVersion", "POST", "/v1/tenants/tenantID-1/abac/policies/policyID-1/versions/version-1/restore", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.RestoreVersion(ctx, "tenantID-1", "policyID-1", "version-1")
		return err
	}},
	{"ABACPolicies.Evaluate", "POST", "/v1/tenants/tenantID-1/abac/evaluate", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Evaluate(ctx, "tenantID-1", ABACEvaluateParams{}, nil)
		return err
	}},
	{"Relationships.List", "GET", "/v1/tenants/tenantID-1/relation-tuples", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.List(ctx, "tenantID-1", nil)
		return err
	}},
	{"Relationships.Create", "POST", "/v1/tenants/tenantID-1/relation-tuples", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.Create(ctx, "tenantID-1", CreateRelationTupleParams{})
		return err
	}},
	{"Relationships.Delete", "DELETE", "/v1/tenants/tenantID-1/relation-tuples/tupleID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Relationships.Delete(ctx, "tenantID-1", "tupleID-1")
	}},
	{"Relationships.Check", "POST", "/v1/tenants/tenantID-1/relation-tuples/check", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.Check(ctx, "tenantID-1", CheckRelationParams{}, nil)
		return err
	}},
	{"Relationships.ListTypes", "GET", "/v1/tenants/tenantID-1/relation-tuples/types", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.ListTypes(ctx, "tenantID-1")
		return err
	}},
	{"Relationships.ListObjects", "GET", "/v1/tenants/tenantID-1/relation-tuples/objects", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.ListObjects(ctx, "tenantID-1", nil)
		return err
	}},
	{"Relationships.AllObjects", "GET", "/v1/tenants/tenantID-1/relation-tuples/objects", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.Relationships.AllObjects(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"Relationships.Graph", "GET", "/v1/tenants/tenantID-1/relation-tuples/graph", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.Graph(ctx, "tenantID-1", nil)
		return err
	}},
	{"AuthZEN.Evaluate", "POST", "/v1/tenants/tenantID-1/access/v1/evaluation", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthZEN.Evaluate(ctx, "tenantID-1", EvaluationParams{})
		return err
	}},
	{"AuthZEN.EvaluateBatch", "POST", "/v1/tenants/tenantID-1/access/v1/evaluations", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthZEN.EvaluateBatch(ctx, "tenantID-1", BatchEvaluationParams{})
		return err
	}},
	{"AuthZEN.ListDecisions", "GET", "/v1/tenants/tenantID-1/authorization/decisions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthZEN.ListDecisions(ctx, "tenantID-1", nil)
		return err
	}},
	{"AuthZEN.AllDecisions", "GET", "/v1/tenants/tenantID-1/authorization/decisions", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.AuthZEN.AllDecisions(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"AuthPolicy.Get", "GET", "/v1/tenants/tenantID-1/auth-policy", false, false, func(ctx context.Context, c *Client) error { _, err := c.AuthPolicy.Get(ctx, "tenantID-1"); return err }},
	{"AuthPolicy.Update", "PUT", "/v1/tenants/tenantID-1/auth-policy", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthPolicy.Update(ctx, "tenantID-1", AuthPolicy{})
		return err
	}},
	{"AuthPolicy.Usage", "GET", "/v1/tenants/tenantID-1/auth-policy/usage", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthPolicy.Usage(ctx, "tenantID-1")
		return err
	}},
	{"SecurityPolicy.Get", "GET", "/v1/tenants/tenantID-1/policy", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SecurityPolicy.Get(ctx, "tenantID-1")
		return err
	}},
	{"SecurityPolicy.Replace", "PUT", "/v1/tenants/tenantID-1/policy", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SecurityPolicy.Replace(ctx, "tenantID-1", SecurityPolicy{})
		return err
	}},
	{"SecurityPolicy.Update", "PATCH", "/v1/tenants/tenantID-1/policy", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SecurityPolicy.Update(ctx, "tenantID-1", UpdateSecurityPolicyParams{})
		return err
	}},
	{"IPRules.List", "GET", "/v1/tenants/tenantID-1/ip-rules", false, false, func(ctx context.Context, c *Client) error { _, err := c.IPRules.List(ctx, "tenantID-1"); return err }},
	{"IPRules.Create", "POST", "/v1/tenants/tenantID-1/ip-rules", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.Create(ctx, "tenantID-1", CreateIPRuleParams{})
		return err
	}},
	{"IPRules.Delete", "DELETE", "/v1/tenants/tenantID-1/ip-rules/ruleID-1", false, false, func(ctx context.Context, c *Client) error { return c.IPRules.Delete(ctx, "tenantID-1", "ruleID-1") }},
	{"IPRules.Check", "POST", "/v1/tenants/tenantID-1/ip-rules/check", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.Check(ctx, "tenantID-1", CheckIPParams{})
		return err
	}},
	{"IPRules.SetEnforcement", "PUT", "/v1/tenants/tenantID-1/ip-rules/config", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.SetEnforcement(ctx, "tenantID-1", SetIPEnforcementParams{})
		return err
	}},
	{"Threats.Bots", "GET", "/v1/tenants/tenantID-1/security/bots", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.Bots(ctx, "tenantID-1", nil)
		return err
	}},
	{"Threats.BotSettings", "GET", "/v1/tenants/tenantID-1/security/bots/settings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.BotSettings(ctx, "tenantID-1")
		return err
	}},
	{"Threats.UpdateBotSettings", "PUT", "/v1/tenants/tenantID-1/security/bots/settings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.UpdateBotSettings(ctx, "tenantID-1", BotSettings{})
		return err
	}},
	{"Threats.RiskSettings", "GET", "/v1/tenants/tenantID-1/security/risk-settings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.RiskSettings(ctx, "tenantID-1")
		return err
	}},
	{"Threats.UpdateRiskSettings", "PUT", "/v1/tenants/tenantID-1/security/risk-settings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.UpdateRiskSettings(ctx, "tenantID-1", RiskSettings{})
		return err
	}},
	{"Threats.ListAnomalies", "GET", "/v1/tenants/tenantID-1/security/anomalies", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.ListAnomalies(ctx, "tenantID-1", nil)
		return err
	}},
	{"Threats.AllAnomalies", "GET", "/v1/tenants/tenantID-1/security/anomalies", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.Threats.AllAnomalies(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"Threats.AnomalySummary", "GET", "/v1/tenants/tenantID-1/security/anomalies/summary", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Threats.AnomalySummary(ctx, "tenantID-1")
		return err
	}},
	{"Threats.ResolveAnomaly", "POST", "/v1/tenants/tenantID-1/security/anomalies/anomalyID-1/resolve", false, false, func(ctx context.Context, c *Client) error {
		return c.Threats.ResolveAnomaly(ctx, "tenantID-1", "anomalyID-1")
	}},
	{"OIDCClients.List", "GET", "/v1/tenants/tenantID-1/oidc/clients", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.List(ctx, "tenantID-1")
		return err
	}},
	{"OIDCClients.Create", "POST", "/v1/tenants/tenantID-1/oidc/clients", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.Create(ctx, "tenantID-1", CreateOIDCClientParams{})
		return err
	}},
	{"OIDCClients.Get", "GET", "/v1/tenants/tenantID-1/oidc/clients/clientID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.Get(ctx, "tenantID-1", "clientID-1")
		return err
	}},
	{"OIDCClients.Update", "PATCH", "/v1/tenants/tenantID-1/oidc/clients/clientID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.Update(ctx, "tenantID-1", "clientID-1", UpdateOIDCClientParams{})
		return err
	}},
	{"OIDCClients.Delete", "DELETE", "/v1/tenants/tenantID-1/oidc/clients/clientID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.OIDCClients.Delete(ctx, "tenantID-1", "clientID-1")
	}},
	{"OIDCClients.RotateSecret", "POST", "/v1/tenants/tenantID-1/oidc/clients/clientID-1/rotate-secret", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.RotateSecret(ctx, "tenantID-1", "clientID-1")
		return err
	}},
	{"OIDCClients.AuthFailures", "GET", "/v1/tenants/tenantID-1/oidc/clients/clientID-1/auth-failures", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.AuthFailures(ctx, "tenantID-1", "clientID-1")
		return err
	}},
	{"OIDCClients.ListShadowAI", "GET", "/v1/tenants/tenantID-1/oidc/clients/shadow-ai", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.ListShadowAI(ctx, "tenantID-1")
		return err
	}},
	{"OIDCClients.AcknowledgeShadowAI", "POST", "/v1/tenants/tenantID-1/oidc/clients/clientID-1/review", false, false, func(ctx context.Context, c *Client) error {
		return c.OIDCClients.AcknowledgeShadowAI(ctx, "tenantID-1", "clientID-1")
	}},
	{"OAuthResources.List", "GET", "/v1/tenants/tenantID-1/oidc/resources", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OAuthResources.List(ctx, "tenantID-1")
		return err
	}},
	{"OAuthResources.Create", "POST", "/v1/tenants/tenantID-1/oidc/resources", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OAuthResources.Create(ctx, "tenantID-1", CreateOAuthResourceParams{})
		return err
	}},
	{"OAuthResources.Delete", "DELETE", "/v1/tenants/tenantID-1/oidc/resources/resourceID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.OAuthResources.Delete(ctx, "tenantID-1", "resourceID-1")
	}},
	{"OAuthGrants.List", "GET", "/v1/tenants/tenantID-1/oauth/grants", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OAuthGrants.List(ctx, "tenantID-1", nil)
		return err
	}},
	{"OAuthGrants.All", "GET", "/v1/tenants/tenantID-1/oauth/grants", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.OAuthGrants.All(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"OAuthGrants.Revoke", "DELETE", "/v1/tenants/tenantID-1/oauth/grants/grantID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.OAuthGrants.Revoke(ctx, "tenantID-1", "grantID-1")
	}},
	{"OAuthGrants.RevokeMany", "POST", "/v1/tenants/tenantID-1/oauth/grants/revoke", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OAuthGrants.RevokeMany(ctx, "tenantID-1", RevokeOAuthGrantsParams{})
		return err
	}},
	{"DeviceAuthorizations.List", "GET", "/v1/tenants/tenantID-1/oauth/devices", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.DeviceAuthorizations.List(ctx, "tenantID-1")
		return err
	}},
	{"DeviceAuthorizations.Revoke", "DELETE", "/v1/tenants/tenantID-1/oauth/devices/authorizationID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.DeviceAuthorizations.Revoke(ctx, "tenantID-1", "authorizationID-1")
	}},
	{"ServicePrincipals.Create", "POST", "/v1/service-principals", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.ServicePrincipals.Create(ctx, CreateServicePrincipalParams{})
		return err
	}},
	{"ServicePrincipals.List", "GET", "/v1/tenants/tenantID-1/service-principals", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ServicePrincipals.List(ctx, "tenantID-1")
		return err
	}},
	{"ServicePrincipals.Disable", "DELETE", "/v1/service-principals/principalID-1", false, false, func(ctx context.Context, c *Client) error { return c.ServicePrincipals.Disable(ctx, "principalID-1") }},
	{"ServicePrincipals.Enable", "POST", "/v1/tenants/tenantID-1/service-principals/principalID-1/enable", false, false, func(ctx context.Context, c *Client) error {
		return c.ServicePrincipals.Enable(ctx, "tenantID-1", "principalID-1")
	}},
	{"ServicePrincipals.RotateSecret", "POST", "/v1/tenants/tenantID-1/service-principals/principalID-1/rotate-secret", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.ServicePrincipals.RotateSecret(ctx, "tenantID-1", "principalID-1", RotateSecretParams{})
		return err
	}},
	{"SocialProviders.List", "GET", "/v1/tenants/tenantID-1/social/providers", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SocialProviders.List(ctx, "tenantID-1")
		return err
	}},
	{"SocialProviders.Upsert", "POST", "/v1/social/providers", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SocialProviders.Upsert(ctx, UpsertSocialProviderParams{})
		return err
	}},
	{"SocialProviders.Update", "PATCH", "/v1/social/providers/provider-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SocialProviders.Update(ctx, "provider-1", UpdateSocialProviderParams{})
		return err
	}},
	{"SocialProviders.Delete", "DELETE", "/v1/social/providers/provider-1", false, false, func(ctx context.Context, c *Client) error { return c.SocialProviders.Delete(ctx, "provider-1") }},
	{"SocialProviders.UnlinkIdentity", "DELETE", "/v1/social/identities/identityID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SocialProviders.UnlinkIdentity(ctx, "identityID-1")
	}},
	{"SAMLConnections.List", "GET", "/v1/tenants/tenantID-1/saml", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.List(ctx, "tenantID-1")
		return err
	}},
	{"SAMLConnections.Create", "POST", "/v1/tenants/tenantID-1/saml", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Create(ctx, "tenantID-1", SAMLConnectionParams{})
		return err
	}},
	{"SAMLConnections.Import", "POST", "/v1/tenants/tenantID-1/saml/import", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Import(ctx, "tenantID-1", ImportSAMLConnectionParams{})
		return err
	}},
	{"SAMLConnections.Get", "GET", "/v1/tenants/tenantID-1/saml/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Get(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"SAMLConnections.Update", "PATCH", "/v1/tenants/tenantID-1/saml/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Update(ctx, "tenantID-1", "connectionID-1", SAMLConnectionParams{})
		return err
	}},
	{"SAMLConnections.Delete", "DELETE", "/v1/tenants/tenantID-1/saml/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLConnections.Delete(ctx, "tenantID-1", "connectionID-1")
	}},
	{"SAMLConnections.ApplyMetadata", "POST", "/v1/tenants/tenantID-1/saml/connectionID-1/metadata", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.ApplyMetadata(ctx, "tenantID-1", "connectionID-1", ApplySAMLMetadataParams{})
		return err
	}},
	{"SAMLConnections.AddCertificate", "POST", "/v1/tenants/tenantID-1/saml/connectionID-1/certificates", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.AddCertificate(ctx, "tenantID-1", "connectionID-1", AddSAMLCertificateParams{})
		return err
	}},
	{"SAMLConnections.RemoveCertificate", "DELETE", "/v1/tenants/tenantID-1/saml/connectionID-1/certificates/certificateID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLConnections.RemoveCertificate(ctx, "tenantID-1", "connectionID-1", "certificateID-1")
	}},
	{"SAMLConnections.Test", "POST", "/v1/tenants/tenantID-1/saml/connectionID-1/test", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Test(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"SAMLConnections.ListGroupMappings", "GET", "/v1/tenants/tenantID-1/saml/connectionID-1/group-mappings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.ListGroupMappings(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"SAMLConnections.CreateGroupMapping", "POST", "/v1/tenants/tenantID-1/saml/connectionID-1/group-mappings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.CreateGroupMapping(ctx, "tenantID-1", "connectionID-1", CreateGroupMappingParams{})
		return err
	}},
	{"SAMLConnections.DeleteGroupMapping", "DELETE", "/v1/tenants/tenantID-1/saml/connectionID-1/group-mappings/mappingID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLConnections.DeleteGroupMapping(ctx, "tenantID-1", "connectionID-1", "mappingID-1")
	}},
	{"SAMLServiceProviders.List", "GET", "/v1/tenants/tenantID-1/saml-providers", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.List(ctx, "tenantID-1")
		return err
	}},
	{"SAMLServiceProviders.Create", "POST", "/v1/tenants/tenantID-1/saml-providers", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.Create(ctx, "tenantID-1", SAMLServiceProviderParams{})
		return err
	}},
	{"SAMLServiceProviders.Get", "GET", "/v1/tenants/tenantID-1/saml-providers/providerID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.Get(ctx, "tenantID-1", "providerID-1")
		return err
	}},
	{"SAMLServiceProviders.Update", "PATCH", "/v1/tenants/tenantID-1/saml-providers/providerID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.Update(ctx, "tenantID-1", "providerID-1", SAMLServiceProviderParams{})
		return err
	}},
	{"SAMLServiceProviders.Delete", "DELETE", "/v1/tenants/tenantID-1/saml-providers/providerID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.Delete(ctx, "tenantID-1", "providerID-1")
	}},
	{"SAMLServiceProviders.ListAssignments", "GET", "/v1/tenants/tenantID-1/saml-providers/providerID-1/assignments", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.ListAssignments(ctx, "tenantID-1", "providerID-1")
		return err
	}},
	{"SAMLServiceProviders.AssignUser", "POST", "/v1/tenants/tenantID-1/saml-providers/providerID-1/assignments/users/userID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.AssignUser(ctx, "tenantID-1", "providerID-1", "userID-1")
	}},
	{"SAMLServiceProviders.UnassignUser", "DELETE", "/v1/tenants/tenantID-1/saml-providers/providerID-1/assignments/users/userID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.UnassignUser(ctx, "tenantID-1", "providerID-1", "userID-1")
	}},
	{"SAMLServiceProviders.AssignGroup", "POST", "/v1/tenants/tenantID-1/saml-providers/providerID-1/assignments/groups/groupID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.AssignGroup(ctx, "tenantID-1", "providerID-1", "groupID-1")
	}},
	{"SAMLServiceProviders.UnassignGroup", "DELETE", "/v1/tenants/tenantID-1/saml-providers/providerID-1/assignments/groups/groupID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.UnassignGroup(ctx, "tenantID-1", "providerID-1", "groupID-1")
	}},
	{"OIDCConnections.List", "GET", "/v1/tenants/tenantID-1/oidc-connections", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.List(ctx, "tenantID-1")
		return err
	}},
	{"OIDCConnections.Create", "POST", "/v1/tenants/tenantID-1/oidc-connections", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.Create(ctx, "tenantID-1", CreateOIDCConnectionParams{})
		return err
	}},
	{"OIDCConnections.Get", "GET", "/v1/tenants/tenantID-1/oidc-connections/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.Get(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"OIDCConnections.Update", "PATCH", "/v1/tenants/tenantID-1/oidc-connections/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.Update(ctx, "tenantID-1", "connectionID-1", UpdateOIDCConnectionParams{})
		return err
	}},
	{"OIDCConnections.Delete", "DELETE", "/v1/tenants/tenantID-1/oidc-connections/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.OIDCConnections.Delete(ctx, "tenantID-1", "connectionID-1")
	}},
	{"LDAP.List", "GET", "/v1/tenants/tenantID-1/ldap", false, false, func(ctx context.Context, c *Client) error { _, err := c.LDAP.List(ctx, "tenantID-1"); return err }},
	{"LDAP.Create", "POST", "/v1/tenants/tenantID-1/ldap", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Create(ctx, "tenantID-1", CreateLDAPConnectionParams{})
		return err
	}},
	{"LDAP.Get", "GET", "/v1/tenants/tenantID-1/ldap/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Get(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"LDAP.Update", "PATCH", "/v1/tenants/tenantID-1/ldap/connectionID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Update(ctx, "tenantID-1", "connectionID-1", UpdateLDAPConnectionParams{})
		return err
	}},
	{"LDAP.Delete", "DELETE", "/v1/tenants/tenantID-1/ldap/connectionID-1", false, false, func(ctx context.Context, c *Client) error { return c.LDAP.Delete(ctx, "tenantID-1", "connectionID-1") }},
	{"LDAP.Test", "POST", "/v1/tenants/tenantID-1/ldap/connectionID-1/test", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Test(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"LDAP.ListGroupMappings", "GET", "/v1/tenants/tenantID-1/ldap/connectionID-1/group-mappings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.ListGroupMappings(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"LDAP.CreateGroupMapping", "POST", "/v1/tenants/tenantID-1/ldap/connectionID-1/group-mappings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.CreateGroupMapping(ctx, "tenantID-1", "connectionID-1", CreateGroupMappingParams{})
		return err
	}},
	{"LDAP.DeleteGroupMapping", "DELETE", "/v1/tenants/tenantID-1/ldap/connectionID-1/group-mappings/mappingID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.LDAP.DeleteGroupMapping(ctx, "tenantID-1", "connectionID-1", "mappingID-1")
	}},
	{"Directory.ListConnections", "GET", "/v1/tenants/tenantID-1/directory/connections", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListConnections(ctx, "tenantID-1")
		return err
	}},
	{"Directory.ListMappings", "GET", "/v1/tenants/tenantID-1/directory/connections/connectionID-1/mappings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListMappings(ctx, "tenantID-1", "connectionID-1")
		return err
	}},
	{"Directory.CreateMapping", "POST", "/v1/tenants/tenantID-1/directory/connections/connectionID-1/mappings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.CreateMapping(ctx, "tenantID-1", "connectionID-1", DirectoryMapping{})
		return err
	}},
	{"Directory.UpdateMapping", "PUT", "/v1/tenants/tenantID-1/directory/connections/connectionID-1/mappings/mappingID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.UpdateMapping(ctx, "tenantID-1", "connectionID-1", "mappingID-1", DirectoryMapping{})
		return err
	}},
	{"Directory.DeleteMapping", "DELETE", "/v1/tenants/tenantID-1/directory/connections/connectionID-1/mappings/mappingID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Directory.DeleteMapping(ctx, "tenantID-1", "connectionID-1", "mappingID-1")
	}},
	{"Directory.PreviewMappings", "POST", "/v1/tenants/tenantID-1/directory/connections/connectionID-1/mappings/preview", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.PreviewMappings(ctx, "tenantID-1", "connectionID-1", PreviewMappingsParams{})
		return err
	}},
	{"Directory.ListErrors", "GET", "/v1/tenants/tenantID-1/directory/errors", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListErrors(ctx, "tenantID-1", nil)
		return err
	}},
	{"Directory.ResolveError", "POST", "/v1/tenants/tenantID-1/directory/errors/eventID-1/resolve", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ResolveError(ctx, "tenantID-1", "eventID-1")
		return err
	}},
	{"Directory.RetryError", "POST", "/v1/tenants/tenantID-1/directory/errors/eventID-1/retry", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.RetryError(ctx, "tenantID-1", "eventID-1")
		return err
	}},
	{"Directory.ListActivity", "GET", "/v1/tenants/tenantID-1/directory/activity", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListActivity(ctx, "tenantID-1", nil)
		return err
	}},
	{"Directory.ListAlerts", "GET", "/v1/tenants/tenantID-1/directory/alerts", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListAlerts(ctx, "tenantID-1")
		return err
	}},
	{"Directory.CreateAlert", "POST", "/v1/tenants/tenantID-1/directory/alerts", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.CreateAlert(ctx, "tenantID-1", CreateDirectoryAlertParams{})
		return err
	}},
	{"Directory.UpdateAlert", "PATCH", "/v1/tenants/tenantID-1/directory/alerts/alertID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.UpdateAlert(ctx, "tenantID-1", "alertID-1", UpdateDirectoryAlertParams{})
		return err
	}},
	{"Directory.DeleteAlert", "DELETE", "/v1/tenants/tenantID-1/directory/alerts/alertID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Directory.DeleteAlert(ctx, "tenantID-1", "alertID-1")
	}},
	{"Directory.Profile", "GET", "/v1/tenants/tenantID-1/directory/profiles/userID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Directory.Profile(ctx, "tenantID-1", "userID-1")
		return err
	}},
	{"SCIM.Config", "GET", "/v1/tenants/tenantID-1/scim", false, false, func(ctx context.Context, c *Client) error { _, err := c.SCIM.Config(ctx, "tenantID-1"); return err }},
	{"SCIM.ListTokens", "GET", "/v1/tenants/tenantID-1/scim/tokens", false, false, func(ctx context.Context, c *Client) error { _, err := c.SCIM.ListTokens(ctx, "tenantID-1"); return err }},
	{"SCIM.CreateToken", "POST", "/v1/tenants/tenantID-1/scim/tokens", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.CreateToken(ctx, "tenantID-1", CreateSCIMTokenParams{})
		return err
	}},
	{"SCIM.RotateToken", "POST", "/v1/tenants/tenantID-1/scim/tokens/tokenID-1/rotate", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.RotateToken(ctx, "tenantID-1", "tokenID-1", RotateSCIMTokenParams{})
		return err
	}},
	{"SCIM.RevokeToken", "DELETE", "/v1/tenants/tenantID-1/scim/tokens/tokenID-1", false, false, func(ctx context.Context, c *Client) error { return c.SCIM.RevokeToken(ctx, "tenantID-1", "tokenID-1") }},
	{"SCIM.ListUsers", "GET", "/v1/tenants/tenantID-1/scim/users", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.ListUsers(ctx, "tenantID-1", nil)
		return err
	}},
	{"SCIM.AllUsers", "GET", "/v1/tenants/tenantID-1/scim/users", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.SCIM.AllUsers(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"SCIM.ListGroups", "GET", "/v1/tenants/tenantID-1/scim/groups", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.ListGroups(ctx, "tenantID-1", nil)
		return err
	}},
	{"SCIM.AllGroups", "GET", "/v1/tenants/tenantID-1/scim/groups", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.SCIM.AllGroups(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"AdminPortalLinks.Create", "POST", "/v1/tenants/tenantID-1/admin-portal/links", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AdminPortalLinks.Create(ctx, "tenantID-1", CreateAdminPortalLinkParams{})
		return err
	}},
	{"AdminPortalLinks.List", "GET", "/v1/tenants/tenantID-1/admin-portal/links", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AdminPortalLinks.List(ctx, "tenantID-1")
		return err
	}},
	{"AdminPortalLinks.Revoke", "DELETE", "/v1/tenants/tenantID-1/admin-portal/links/linkID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.AdminPortalLinks.Revoke(ctx, "tenantID-1", "linkID-1")
	}},
	{"APIKeys.Create", "POST", "/v1/api-keys", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.APIKeys.Create(ctx, CreateAPIKeyParams{})
		return err
	}},
	{"APIKeys.List", "GET", "/v1/tenants/tenantID-1/api-keys", false, false, func(ctx context.Context, c *Client) error { _, err := c.APIKeys.List(ctx, "tenantID-1"); return err }},
	{"APIKeys.Rotate", "POST", "/v1/tenants/tenantID-1/api-keys/keyID-1/rotate", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.APIKeys.Rotate(ctx, "tenantID-1", "keyID-1", RotateSecretParams{})
		return err
	}},
	{"APIKeys.Revoke", "DELETE", "/v1/api-keys/keyID-1", false, false, func(ctx context.Context, c *Client) error { return c.APIKeys.Revoke(ctx, "keyID-1") }},
	{"Agents.Create", "POST", "/v1/tenants/tenantID-1/agents", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Create(ctx, "tenantID-1", CreateAgentParams{})
		return err
	}},
	{"Agents.List", "GET", "/v1/tenants/tenantID-1/agents", false, false, func(ctx context.Context, c *Client) error { _, err := c.Agents.List(ctx, "tenantID-1"); return err }},
	{"Agents.Update", "PATCH", "/v1/tenants/tenantID-1/agents/agentID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Agents.Update(ctx, "tenantID-1", "agentID-1", UpdateAgentParams{})
	}},
	{"Agents.Delete", "DELETE", "/v1/tenants/tenantID-1/agents/agentID-1", false, false, func(ctx context.Context, c *Client) error { return c.Agents.Delete(ctx, "tenantID-1", "agentID-1") }},
	{"Agents.Suspend", "POST", "/v1/tenants/tenantID-1/agents/agentID-1/suspend", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Suspend(ctx, "tenantID-1", "agentID-1")
		return err
	}},
	{"Agents.Resume", "POST", "/v1/tenants/tenantID-1/agents/agentID-1/resume", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Resume(ctx, "tenantID-1", "agentID-1")
		return err
	}},
	{"Agents.Decommission", "POST", "/v1/tenants/tenantID-1/agents/agentID-1/decommission", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Decommission(ctx, "tenantID-1", "agentID-1")
		return err
	}},
	{"Agents.SuspendAll", "POST", "/v1/tenants/tenantID-1/agents/kill-all", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.SuspendAll(ctx, "tenantID-1")
		return err
	}},
	{"Agents.ResumeAll", "POST", "/v1/tenants/tenantID-1/agents/resume-all", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.ResumeAll(ctx, "tenantID-1")
		return err
	}},
	{"Agents.RotateSecret", "POST", "/v1/tenants/tenantID-1/agents/agentID-1/rotate-secret", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.RotateSecret(ctx, "tenantID-1", "agentID-1", RotateSecretParams{})
		return err
	}},
	{"Agents.ListSponsoredBy", "GET", "/v1/tenants/tenantID-1/agents/sponsored-by/userID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.ListSponsoredBy(ctx, "tenantID-1", "userID-1")
		return err
	}},
	{"Agents.TransferSponsorship", "POST", "/v1/tenants/tenantID-1/agents/sponsored-by/userID-1/transfer", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.TransferSponsorship(ctx, "tenantID-1", "userID-1", TransferSponsorshipParams{})
		return err
	}},
	{"Agents.Token", "POST", "/v1/agents/token", true, false, func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Token(ctx, AgentTokenParams{})
		return err
	}},
	{"Credentials.Issue", "POST", "/v1/tenants/tenantID-1/credentials", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Credentials.Issue(ctx, "tenantID-1", IssueCredentialParams{})
		return err
	}},
	{"Credentials.List", "GET", "/v1/tenants/tenantID-1/credentials", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Credentials.List(ctx, "tenantID-1")
		return err
	}},
	{"Credentials.Revoke", "POST", "/v1/tenants/tenantID-1/credentials/credentialID-1/revoke", false, false, func(ctx context.Context, c *Client) error {
		return c.Credentials.Revoke(ctx, "tenantID-1", "credentialID-1")
	}},
	{"Credentials.Verify", "POST", "/v1/credentials/verify", true, false, func(ctx context.Context, c *Client) error {
		_, err := c.Credentials.Verify(ctx, VerifyCredentialParams{})
		return err
	}},
	{"AuthHooks.List", "GET", "/v1/tenants/tenantID-1/auth-hooks", false, false, func(ctx context.Context, c *Client) error { _, err := c.AuthHooks.List(ctx, "tenantID-1"); return err }},
	{"AuthHooks.Create", "POST", "/v1/tenants/tenantID-1/auth-hooks", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.Create(ctx, "tenantID-1", CreateAuthHookParams{})
		return err
	}},
	{"AuthHooks.Update", "PATCH", "/v1/tenants/tenantID-1/auth-hooks/hookID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.Update(ctx, "tenantID-1", "hookID-1", UpdateAuthHookParams{})
		return err
	}},
	{"AuthHooks.Delete", "DELETE", "/v1/tenants/tenantID-1/auth-hooks/hookID-1", false, false, func(ctx context.Context, c *Client) error { return c.AuthHooks.Delete(ctx, "tenantID-1", "hookID-1") }},
	{"AuthHooks.Test", "POST", "/v1/tenants/tenantID-1/auth-hooks/hookID-1/test", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.Test(ctx, "tenantID-1", "hookID-1")
		return err
	}},
	{"Webhooks.Create", "POST", "/v1/webhooks", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.Create(ctx, CreateWebhookParams{})
		return err
	}},
	{"Webhooks.Get", "GET", "/v1/webhooks/webhookID-1", false, false, func(ctx context.Context, c *Client) error { _, err := c.Webhooks.Get(ctx, "webhookID-1"); return err }},
	{"Webhooks.List", "GET", "/v1/tenants/tenantID-1/webhooks", false, false, func(ctx context.Context, c *Client) error { _, err := c.Webhooks.List(ctx, "tenantID-1"); return err }},
	{"Webhooks.Update", "PATCH", "/v1/tenants/tenantID-1/webhooks/webhookID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.Update(ctx, "tenantID-1", "webhookID-1", UpdateWebhookParams{})
		return err
	}},
	{"Webhooks.Disable", "DELETE", "/v1/webhooks/webhookID-1", false, false, func(ctx context.Context, c *Client) error { return c.Webhooks.Disable(ctx, "webhookID-1") }},
	{"Webhooks.Enable", "POST", "/v1/tenants/tenantID-1/webhooks/webhookID-1/enable", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.Enable(ctx, "tenantID-1", "webhookID-1")
		return err
	}},
	{"Webhooks.RotateSecret", "POST", "/v1/tenants/tenantID-1/webhooks/webhookID-1/rotate-secret", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.RotateSecret(ctx, "tenantID-1", "webhookID-1")
		return err
	}},
	{"Webhooks.Test", "POST", "/v1/webhooks/webhookID-1/test", false, false, func(ctx context.Context, c *Client) error { _, err := c.Webhooks.Test(ctx, "webhookID-1"); return err }},
	{"Webhooks.ListEventTypes", "GET", "/v1/tenants/tenantID-1/webhooks/events", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.ListEventTypes(ctx, "tenantID-1")
		return err
	}},
	{"Webhooks.ListDeliveries", "GET", "/v1/webhooks/webhookID-1/deliveries", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.ListDeliveries(ctx, "webhookID-1", nil)
		return err
	}},
	{"Webhooks.AllDeliveries", "GET", "/v1/webhooks/webhookID-1/deliveries", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.Webhooks.AllDeliveries(ctx, "webhookID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"Webhooks.ListAttempts", "GET", "/v1/tenants/tenantID-1/webhooks/webhookID-1/deliveries/deliveryID-1/attempts", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.ListAttempts(ctx, "tenantID-1", "webhookID-1", "deliveryID-1")
		return err
	}},
	{"Webhooks.RetryDelivery", "POST", "/v1/webhooks/webhookID-1/deliveries/deliveryID-1/retry", false, false, func(ctx context.Context, c *Client) error {
		return c.Webhooks.RetryDelivery(ctx, "webhookID-1", "deliveryID-1")
	}},
	{"Webhooks.RetryGivenUp", "POST", "/v1/tenants/tenantID-1/webhooks/webhookID-1/deliveries/retry", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.RetryGivenUp(ctx, "tenantID-1", "webhookID-1")
		return err
	}},
	{"Secrets.List", "GET", "/v1/tenants/tenantID-1/secrets", false, false, func(ctx context.Context, c *Client) error { _, err := c.Secrets.List(ctx, "tenantID-1"); return err }},
	{"Secrets.Create", "POST", "/v1/tenants/tenantID-1/secrets", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.Create(ctx, "tenantID-1", CreateSecretParams{})
		return err
	}},
	{"Secrets.Update", "PATCH", "/v1/tenants/tenantID-1/secrets/secretID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.Update(ctx, "tenantID-1", "secretID-1", UpdateSecretParams{})
		return err
	}},
	{"Secrets.Delete", "DELETE", "/v1/tenants/tenantID-1/secrets/secretID-1", false, false, func(ctx context.Context, c *Client) error { return c.Secrets.Delete(ctx, "tenantID-1", "secretID-1") }},
	{"Secrets.Reveal", "POST", "/v1/tenants/tenantID-1/secrets/secretID-1/reveal", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.Reveal(ctx, "tenantID-1", "secretID-1")
		return err
	}},
	{"Secrets.GetByName", "GET", "/v1/vault/name-1", false, false, func(ctx context.Context, c *Client) error { _, err := c.Secrets.GetByName(ctx, "name-1"); return err }},
	{"TokenVault.RegisterProvider", "POST", "/v1/tenants/tenantID-1/vault/tokens/providers", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.RegisterProvider(ctx, "tenantID-1", RegisterTokenVaultProviderParams{})
		return err
	}},
	{"TokenVault.ListProviders", "GET", "/v1/tenants/tenantID-1/vault/tokens/providers", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.ListProviders(ctx, "tenantID-1")
		return err
	}},
	{"TokenVault.DeleteProvider", "DELETE", "/v1/tenants/tenantID-1/vault/tokens/providers/provider-1", false, false, func(ctx context.Context, c *Client) error {
		return c.TokenVault.DeleteProvider(ctx, "tenantID-1", "provider-1")
	}},
	{"TokenVault.ListConnections", "GET", "/v1/vault/tokens", false, false, func(ctx context.Context, c *Client) error { _, err := c.TokenVault.ListConnections(ctx); return err }},
	{"TokenVault.AccessToken", "GET", "/v1/vault/tokens/provider-1/access-token", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.AccessToken(ctx, "provider-1")
		return err
	}},
	{"TokenVault.Disconnect", "DELETE", "/v1/vault/tokens/provider-1", false, false, func(ctx context.Context, c *Client) error { return c.TokenVault.Disconnect(ctx, "provider-1") }},
	{"Audit.List", "GET", "/v1/tenants/tenantID-1/audit", false, false, func(ctx context.Context, c *Client) error { _, err := c.Audit.List(ctx, "tenantID-1", nil); return err }},
	{"Audit.All", "GET", "/v1/tenants/tenantID-1/audit", false, false, func(ctx context.Context, c *Client) error {
		for _, err := range c.Audit.All(ctx, "tenantID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{"Audit.Verify", "GET", "/v1/tenants/tenantID-1/audit/verify", false, false, func(ctx context.Context, c *Client) error { _, err := c.Audit.Verify(ctx, "tenantID-1"); return err }},
	{"Audit.LastVerification", "GET", "/v1/tenants/tenantID-1/audit/verification", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Audit.LastVerification(ctx, "tenantID-1")
		return err
	}},
	{"Audit.ListAnomalies", "GET", "/v1/tenants/tenantID-1/audit/anomalies", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Audit.ListAnomalies(ctx, "tenantID-1", nil)
		return err
	}},
	{"Audit.AnomalySummary", "GET", "/v1/tenants/tenantID-1/audit/anomalies/summary", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Audit.AnomalySummary(ctx, "tenantID-1")
		return err
	}},
	{"Audit.ResolveAnomaly", "POST", "/v1/tenants/tenantID-1/audit/anomalies/anomalyID-1/resolve", false, false, func(ctx context.Context, c *Client) error {
		return c.Audit.ResolveAnomaly(ctx, "tenantID-1", "anomalyID-1")
	}},
	{"Audit.AnomalySettings", "GET", "/v1/tenants/tenantID-1/audit/anomaly-settings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Audit.AnomalySettings(ctx, "tenantID-1")
		return err
	}},
	{"Audit.UpdateAnomalySettings", "PUT", "/v1/tenants/tenantID-1/audit/anomaly-settings", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Audit.UpdateAnomalySettings(ctx, "tenantID-1", AuditAnomalySettings{})
		return err
	}},
	{"Analytics.Overview", "GET", "/v1/tenants/tenantID-1/analytics/overview", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Analytics.Overview(ctx, "tenantID-1", nil)
		return err
	}},
	{"Compliance.GenerateEvidence", "POST", "/v1/tenants/tenantID-1/compliance/framework-1/evidence", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.GenerateEvidence(ctx, "tenantID-1", "framework-1")
		return err
	}},
	{"Compliance.ListEvidence", "GET", "/v1/tenants/tenantID-1/compliance/framework-1/evidence", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.ListEvidence(ctx, "tenantID-1", "framework-1")
		return err
	}},
	{"Compliance.GetEvidence", "GET", "/v1/tenants/tenantID-1/compliance/evidence/runID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.GetEvidence(ctx, "tenantID-1", "runID-1")
		return err
	}},
	{"Compliance.ListControlOwners", "GET", "/v1/tenants/tenantID-1/compliance/controls/owners", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.ListControlOwners(ctx, "tenantID-1")
		return err
	}},
	{"Compliance.SetControlOwner", "PUT", "/v1/tenants/tenantID-1/compliance/controls/controlID-1/owner", false, false, func(ctx context.Context, c *Client) error {
		return c.Compliance.SetControlOwner(ctx, "tenantID-1", "controlID-1", SetControlOwnerParams{})
	}},
	{"Compliance.RemoveControlOwner", "DELETE", "/v1/tenants/tenantID-1/compliance/controls/controlID-1/owner", false, false, func(ctx context.Context, c *Client) error {
		return c.Compliance.RemoveControlOwner(ctx, "tenantID-1", "controlID-1")
	}},
	{"Compliance.ListAudits", "GET", "/v1/tenants/tenantID-1/compliance/audits", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.ListAudits(ctx, "tenantID-1", nil)
		return err
	}},
	{"Compliance.CreateAudit", "POST", "/v1/tenants/tenantID-1/compliance/audits", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.CreateAudit(ctx, "tenantID-1", CreateExternalAuditParams{})
		return err
	}},
	{"Compliance.DeleteAudit", "DELETE", "/v1/tenants/tenantID-1/compliance/audits/auditID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.Compliance.DeleteAudit(ctx, "tenantID-1", "auditID-1")
	}},
	{"Retention.Get", "GET", "/v1/tenants/tenantID-1/retention", false, false, func(ctx context.Context, c *Client) error { _, err := c.Retention.Get(ctx, "tenantID-1"); return err }},
	{"Retention.Update", "PUT", "/v1/tenants/tenantID-1/retention", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Retention.Update(ctx, "tenantID-1", RetentionPolicy{})
		return err
	}},
	{"Retention.Preview", "POST", "/v1/tenants/tenantID-1/retention/preview", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Retention.Preview(ctx, "tenantID-1")
		return err
	}},
	{"Retention.Run", "POST", "/v1/tenants/tenantID-1/retention/run", false, false, func(ctx context.Context, c *Client) error { _, err := c.Retention.Run(ctx, "tenantID-1"); return err }},
	{"GDPR.CreatePurge", "POST", "/v1/gdpr/purge", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.CreatePurge(ctx, CreatePurgeRequestParams{})
		return err
	}},
	{"GDPR.CancelPurge", "DELETE", "/v1/gdpr/purge/requestID-1", false, false, func(ctx context.Context, c *Client) error { return c.GDPR.CancelPurge(ctx, "requestID-1") }},
	{"GDPR.ListPurges", "GET", "/v1/tenants/tenantID-1/gdpr/purge", false, false, func(ctx context.Context, c *Client) error { _, err := c.GDPR.ListPurges(ctx, "tenantID-1"); return err }},
	{"GDPR.CreateExport", "POST", "/v1/gdpr/export", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.CreateExport(ctx, CreateExportRequestParams{})
		return err
	}},
	{"GDPR.ListExports", "GET", "/v1/tenants/tenantID-1/gdpr/export", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.ListExports(ctx, "tenantID-1")
		return err
	}},
	{"GDPR.GetExport", "GET", "/v1/tenants/tenantID-1/gdpr/export/requestID-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.GetExport(ctx, "tenantID-1", "requestID-1")
		return err
	}},
	{"LogSinks.List", "GET", "/v1/tenants/tenantID-1/log-sinks", false, false, func(ctx context.Context, c *Client) error { _, err := c.LogSinks.List(ctx, "tenantID-1"); return err }},
	{"LogSinks.Create", "POST", "/v1/tenants/tenantID-1/log-sinks", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.Create(ctx, "tenantID-1", LogSinkParams{})
		return err
	}},
	{"LogSinks.SetEnabled", "PATCH", "/v1/tenants/tenantID-1/log-sinks/sinkID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.LogSinks.SetEnabled(ctx, "tenantID-1", "sinkID-1", SetLogSinkEnabledParams{})
	}},
	{"LogSinks.Delete", "DELETE", "/v1/tenants/tenantID-1/log-sinks/sinkID-1", false, false, func(ctx context.Context, c *Client) error { return c.LogSinks.Delete(ctx, "tenantID-1", "sinkID-1") }},
	{"LogSinks.Test", "POST", "/v1/tenants/tenantID-1/log-sinks/sinkID-1/test", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.Test(ctx, "tenantID-1", "sinkID-1")
		return err
	}},
	{"LogSinks.TestSettings", "POST", "/v1/tenants/tenantID-1/log-sinks/test", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.TestSettings(ctx, "tenantID-1", LogSinkParams{})
		return err
	}},
	{"RateLimits.Get", "GET", "/v1/tenants/tenantID-1/rate-limits", false, false, func(ctx context.Context, c *Client) error { _, err := c.RateLimits.Get(ctx, "tenantID-1"); return err }},
	{"RateLimits.Update", "PUT", "/v1/tenants/tenantID-1/rate-limits", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.RateLimits.Update(ctx, "tenantID-1", UpdateRateLimitsParams{})
		return err
	}},
	{"RateLimits.Reset", "DELETE", "/v1/tenants/tenantID-1/rate-limits", false, false, func(ctx context.Context, c *Client) error { return c.RateLimits.Reset(ctx, "tenantID-1") }},
	{"EmailTemplates.List", "GET", "/v1/tenants/tenantID-1/email-templates", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.List(ctx, "tenantID-1")
		return err
	}},
	{"EmailTemplates.Get", "GET", "/v1/tenants/tenantID-1/email-templates/key-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Get(ctx, "tenantID-1", "key-1")
		return err
	}},
	{"EmailTemplates.Update", "PUT", "/v1/tenants/tenantID-1/email-templates/key-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Update(ctx, "tenantID-1", "key-1", UpdateEmailTemplateParams{})
		return err
	}},
	{"EmailTemplates.Reset", "DELETE", "/v1/tenants/tenantID-1/email-templates/key-1", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Reset(ctx, "tenantID-1", "key-1")
		return err
	}},
	{"EmailTemplates.Preview", "POST", "/v1/tenants/tenantID-1/email-templates/key-1/preview", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Preview(ctx, "tenantID-1", "key-1", PreviewEmailTemplateParams{})
		return err
	}},
	{"EmailTemplates.SendTest", "POST", "/v1/tenants/tenantID-1/email-templates/key-1/test", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.SendTest(ctx, "tenantID-1", "key-1", SendTestEmailParams{})
		return err
	}},
	{"EmailDelivery.ListMessages", "GET", "/v1/tenants/tenantID-1/email/messages", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailDelivery.ListMessages(ctx, "tenantID-1", nil)
		return err
	}},
	{"EmailDelivery.ListSuppressions", "GET", "/v1/tenants/tenantID-1/email/suppressions", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.EmailDelivery.ListSuppressions(ctx, "tenantID-1", nil)
		return err
	}},
	{"EmailDelivery.DeleteSuppression", "DELETE", "/v1/tenants/tenantID-1/email/suppressions/suppressionID-1", false, false, func(ctx context.Context, c *Client) error {
		return c.EmailDelivery.DeleteSuppression(ctx, "tenantID-1", "suppressionID-1")
	}},
	{"Billing.ListPlans", "GET", "/v1/billing/plans", false, false, func(ctx context.Context, c *Client) error { _, err := c.Billing.ListPlans(ctx); return err }},
	{"Billing.Subscription", "GET", "/v1/tenants/tenantID-1/billing/subscription", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Subscription(ctx, "tenantID-1")
		return err
	}},
	{"Billing.Subscribe", "PUT", "/v1/tenants/tenantID-1/billing/subscription", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Subscribe(ctx, "tenantID-1", SubscribeParams{})
		return err
	}},
	{"Billing.CancelSubscription", "POST", "/v1/tenants/tenantID-1/billing/subscription/cancel", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.CancelSubscription(ctx, "tenantID-1")
		return err
	}},
	{"Billing.StartTrial", "POST", "/v1/tenants/tenantID-1/billing/trial", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.StartTrial(ctx, "tenantID-1", StartTrialParams{})
		return err
	}},
	{"Billing.Checkout", "POST", "/v1/tenants/tenantID-1/billing/checkout", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Checkout(ctx, "tenantID-1", CheckoutParams{})
		return err
	}},
	{"Billing.ListInvoices", "GET", "/v1/tenants/tenantID-1/billing/invoices", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.ListInvoices(ctx, "tenantID-1")
		return err
	}},
	{"Billing.PayInvoice", "POST", "/v1/tenants/tenantID-1/billing/invoices/invoiceID-1/pay", false, true, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.PayInvoice(ctx, "tenantID-1", "invoiceID-1", PayInvoiceParams{})
		return err
	}},
	{"Billing.PaymentMethod", "GET", "/v1/tenants/tenantID-1/billing/payment-method", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.PaymentMethod(ctx, "tenantID-1")
		return err
	}},
	{"Billing.Profile", "GET", "/v1/tenants/tenantID-1/billing/profile", false, false, func(ctx context.Context, c *Client) error { _, err := c.Billing.Profile(ctx, "tenantID-1"); return err }},
	{"Billing.UpdateProfile", "PUT", "/v1/tenants/tenantID-1/billing/profile", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Billing.UpdateProfile(ctx, "tenantID-1", BillingProfile{})
		return err
	}},
	{"Entitlements.Get", "GET", "/v1/tenants/tenantID-1/entitlements", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Entitlements.Get(ctx, "tenantID-1")
		return err
	}},
	{"Entitlements.Usage", "GET", "/v1/tenants/tenantID-1/entitlements/usage", false, false, func(ctx context.Context, c *Client) error {
		_, err := c.Entitlements.Usage(ctx, "tenantID-1")
		return err
	}},
	{"SigningKeys.List", "GET", "/v1/oidc/signing-keys", false, false, func(ctx context.Context, c *Client) error { _, err := c.SigningKeys.List(ctx); return err }},
}
