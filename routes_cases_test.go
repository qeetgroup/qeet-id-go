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
	basic      bool // the request authenticates an OAuth client with HTTP Basic
	call       func(ctx context.Context, c *Client) error
}

// routeCases covers every generated method: each must send the documented
// method and path, authenticate unless the endpoint is public, and send an
// Idempotency-Key exactly when the endpoint accepts one.
var routeCases = []routeCase{
	{name: "Organizations.List", method: "GET", path: "/v1/organizations", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Organizations.List(ctx, nil); return err }},
	{name: "Organizations.All", method: "GET", path: "/v1/organizations", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Organizations.All(ctx, nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Organizations.Create", method: "POST", path: "/v1/organizations", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Organizations.Create(ctx, CreateOrganizationParams{})
		return err
	}},
	{name: "Organizations.Get", method: "GET", path: "/v1/organizations/organizationID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Organizations.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "Organizations.Update", method: "PATCH", path: "/v1/organizations/organizationID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Organizations.Update(ctx, "organizationID-1", UpdateOrganizationParams{})
		return err
	}},
	{name: "Organizations.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Organizations.Delete(ctx, "organizationID-1")
		return err
	}},
	{name: "Organizations.CancelDeletion", method: "POST", path: "/v1/organizations/organizationID-1/deletion/cancel", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Organizations.CancelDeletion(ctx, "organizationID-1")
		return err
	}},
	{name: "Users.Create", method: "POST", path: "/v1/users", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.Create(ctx, CreateUserParams{})
		return err
	}},
	{name: "Users.Get", method: "GET", path: "/v1/users/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Get(ctx, "userID-1"); return err }},
	{name: "Users.Update", method: "PATCH", path: "/v1/users/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.Update(ctx, "userID-1", UpdateUserParams{})
		return err
	}},
	{name: "Users.Delete", method: "DELETE", path: "/v1/users/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Delete(ctx, "userID-1"); return err }},
	{name: "Users.List", method: "GET", path: "/v1/users", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.List(ctx, nil); return err }},
	{name: "Users.All", method: "GET", path: "/v1/users", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Users.All(ctx, nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Users.SetPassword", method: "POST", path: "/v1/users/userID-1/password", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Users.SetPassword(ctx, "userID-1", SetPasswordParams{})
	}},
	{name: "Users.BulkCreate", method: "POST", path: "/v1/users/bulk", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.BulkCreate(ctx, BulkCreateUsersParams{})
		return err
	}},
	{name: "Users.ListDeleted", method: "GET", path: "/v1/users/deleted", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.ListDeleted(ctx, nil); return err }},
	{name: "Users.AllDeleted", method: "GET", path: "/v1/users/deleted", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Users.AllDeleted(ctx, nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Users.Restore", method: "POST", path: "/v1/users/userID-1/restore", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Restore(ctx, "userID-1"); return err }},
	{name: "Users.Purge", method: "DELETE", path: "/v1/users/userID-1/purge", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Users.Purge(ctx, "userID-1") }},
	{name: "Users.RemoveMembership", method: "DELETE", path: "/v1/users/userID-1/membership", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.RemoveMembership(ctx, "userID-1")
		return err
	}},
	{name: "Users.ResetMFA", method: "DELETE", path: "/v1/users/userID-1/mfa", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.ResetMFA(ctx, "userID-1"); return err }},
	{name: "Users.SetMFARequired", method: "PUT", path: "/v1/users/userID-1/mfa-required", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.SetMFARequired(ctx, "userID-1", SetMFARequiredParams{})
		return err
	}},
	{name: "Users.Security", method: "GET", path: "/v1/users/userID-1/security", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Security(ctx, "userID-1"); return err }},
	{name: "Users.Access", method: "GET", path: "/v1/users/userID-1/access", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Access(ctx, "userID-1"); return err }},
	{name: "Users.Stats", method: "GET", path: "/v1/users/stats", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Stats(ctx); return err }},
	{name: "Users.Trends", method: "GET", path: "/v1/users/trends", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Trends(ctx); return err }},
	{name: "Users.ListSuspended", method: "GET", path: "/v1/users/suspended", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.ListSuspended(ctx, nil); return err }},
	{name: "Users.Reactivate", method: "POST", path: "/v1/users/userID-1/reactivate", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Users.Reactivate(ctx, "userID-1"); return err }},
	{name: "Users.ReviewSuspension", method: "PUT", path: "/v1/users/userID-1/suspension", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Users.ReviewSuspension(ctx, "userID-1", ReviewSuspensionParams{})
	}},
	{name: "Users.ListSessions", method: "GET", path: "/v1/users/userID-1/sessions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListSessions(ctx, "userID-1")
		return err
	}},
	{name: "Users.RevokeSessions", method: "POST", path: "/v1/users/userID-1/sessions/revoke-all", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.RevokeSessions(ctx, "userID-1")
		return err
	}},
	{name: "Users.EndSession", method: "DELETE", path: "/v1/organizations/organizationID-1/users/userID-1/sessions/sessionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Users.EndSession(ctx, "organizationID-1", "userID-1", "sessionID-1")
	}},
	{name: "Users.ListTrustedDevices", method: "GET", path: "/v1/users/userID-1/trusted-devices", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListTrustedDevices(ctx, "userID-1")
		return err
	}},
	{name: "Users.ForgetTrustedDevices", method: "DELETE", path: "/v1/users/userID-1/trusted-devices", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ForgetTrustedDevices(ctx, "userID-1")
		return err
	}},
	{name: "Users.ForgetTrustedDevice", method: "DELETE", path: "/v1/users/userID-1/trusted-devices/deviceID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ForgetTrustedDevice(ctx, "userID-1", "deviceID-1")
		return err
	}},
	{name: "Users.ListPushDevices", method: "GET", path: "/v1/users/userID-1/push-devices", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListPushDevices(ctx, "userID-1")
		return err
	}},
	{name: "Users.DisconnectPushDevice", method: "DELETE", path: "/v1/users/userID-1/push-devices/deviceID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Users.DisconnectPushDevice(ctx, "userID-1", "deviceID-1")
	}},
	{name: "Users.ListPasskeys", method: "GET", path: "/v1/organizations/organizationID-1/users/userID-1/passkeys", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListPasskeys(ctx, "organizationID-1", "userID-1")
		return err
	}},
	{name: "Users.DeletePasskey", method: "DELETE", path: "/v1/organizations/organizationID-1/users/userID-1/passkeys/passkeyID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Users.DeletePasskey(ctx, "organizationID-1", "userID-1", "passkeyID-1")
	}},
	{name: "Users.SignIns", method: "GET", path: "/v1/organizations/organizationID-1/users/userID-1/sign-ins", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.SignIns(ctx, "organizationID-1", "userID-1", nil)
		return err
	}},
	{name: "Users.StartEmailVerification", method: "POST", path: "/v1/users/userID-1/verify/email/start", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.StartEmailVerification(ctx, "userID-1")
		return err
	}},
	{name: "Users.ConfirmEmailVerification", method: "POST", path: "/v1/users/userID-1/verify/email/confirm", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ConfirmEmailVerification(ctx, "userID-1", ConfirmVerificationParams{})
		return err
	}},
	{name: "Users.StartPhoneVerification", method: "POST", path: "/v1/users/userID-1/verify/phone/start", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.StartPhoneVerification(ctx, "userID-1", StartPhoneVerificationParams{})
		return err
	}},
	{name: "Users.ConfirmPhoneVerification", method: "POST", path: "/v1/users/userID-1/verify/phone/confirm", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ConfirmPhoneVerification(ctx, "userID-1", ConfirmVerificationParams{})
		return err
	}},
	{name: "Users.ListSocialIdentities", method: "GET", path: "/v1/users/userID-1/social/identities", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Users.ListSocialIdentities(ctx, "userID-1")
		return err
	}},
	{name: "Groups.Create", method: "POST", path: "/v1/groups", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.Create(ctx, CreateGroupParams{})
		return err
	}},
	{name: "Groups.Update", method: "PATCH", path: "/v1/groups/groupID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.Update(ctx, "groupID-1", UpdateGroupParams{})
		return err
	}},
	{name: "Groups.Delete", method: "DELETE", path: "/v1/groups/groupID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Groups.Delete(ctx, "groupID-1") }},
	{name: "Groups.List", method: "GET", path: "/v1/organizations/organizationID-1/groups", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.List(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Groups.ListMembers", method: "GET", path: "/v1/groups/groupID-1/members", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.ListMembers(ctx, "groupID-1")
		return err
	}},
	{name: "Groups.AddMember", method: "POST", path: "/v1/groups/groupID-1/members/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Groups.AddMember(ctx, "groupID-1", "userID-1") }},
	{name: "Groups.RemoveMember", method: "DELETE", path: "/v1/groups/groupID-1/members/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Groups.RemoveMember(ctx, "groupID-1", "userID-1") }},
	{name: "Groups.ListRoles", method: "GET", path: "/v1/organizations/organizationID-1/groups/groupID-1/roles", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.ListRoles(ctx, "organizationID-1", "groupID-1")
		return err
	}},
	{name: "Groups.GrantRole", method: "POST", path: "/v1/organizations/organizationID-1/groups/groupID-1/roles/roleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Groups.GrantRole(ctx, "organizationID-1", "groupID-1", "roleID-1")
	}},
	{name: "Groups.RevokeRole", method: "DELETE", path: "/v1/organizations/organizationID-1/groups/groupID-1/roles/roleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Groups.RevokeRole(ctx, "organizationID-1", "groupID-1", "roleID-1")
	}},
	{name: "Groups.Import", method: "POST", path: "/v1/groups/import", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.Import(ctx, ImportGroupsParams{})
		return err
	}},
	{name: "Groups.ListTemplates", method: "GET", path: "/v1/groups/templates", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Groups.ListTemplates(ctx); return err }},
	{name: "Groups.CreateTemplate", method: "POST", path: "/v1/groups/templates", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.CreateTemplate(ctx, CreateGroupTemplateParams{})
		return err
	}},
	{name: "Groups.DeleteTemplate", method: "DELETE", path: "/v1/groups/templates/templateID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Groups.DeleteTemplate(ctx, "templateID-1") }},
	{name: "Groups.ApplyTemplate", method: "POST", path: "/v1/groups/templates/templateID-1/apply", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Groups.ApplyTemplate(ctx, "templateID-1", ApplyGroupTemplateParams{})
		return err
	}},
	{name: "Invitations.Create", method: "POST", path: "/v1/invites", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Invitations.Create(ctx, CreateInvitationParams{})
		return err
	}},
	{name: "Invitations.List", method: "GET", path: "/v1/organizations/organizationID-1/invites", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Invitations.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Invitations.Revoke", method: "DELETE", path: "/v1/invites/invitationID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Invitations.Revoke(ctx, "invitationID-1") }},
	{name: "Invitations.Resend", method: "POST", path: "/v1/invites/invitationID-1/resend", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Invitations.Resend(ctx, "invitationID-1")
		return err
	}},
	{name: "Domains.List", method: "GET", path: "/v1/organizations/organizationID-1/domains", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Domains.Create", method: "POST", path: "/v1/organizations/organizationID-1/domains", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Create(ctx, "organizationID-1", CreateDomainParams{})
		return err
	}},
	{name: "Domains.Verify", method: "POST", path: "/v1/organizations/organizationID-1/domains/domainID-1/verify", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Verify(ctx, "organizationID-1", "domainID-1")
		return err
	}},
	{name: "Domains.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/domains/domainID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.Update(ctx, "organizationID-1", "domainID-1", UpdateDomainParams{})
		return err
	}},
	{name: "Domains.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/domains/domainID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Domains.Delete(ctx, "organizationID-1", "domainID-1")
	}},
	{name: "Domains.LoginDomain", method: "GET", path: "/v1/organizations/organizationID-1/login-domain", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.LoginDomain(ctx, "organizationID-1")
		return err
	}},
	{name: "Domains.CheckLoginDomain", method: "POST", path: "/v1/organizations/organizationID-1/login-domain/check", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.CheckLoginDomain(ctx, "organizationID-1")
		return err
	}},
	{name: "Domains.ListBreakGlassMembers", method: "GET", path: "/v1/organizations/organizationID-1/break-glass-members", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Domains.ListBreakGlassMembers(ctx, "organizationID-1")
		return err
	}},
	{name: "Domains.AddBreakGlassMember", method: "PUT", path: "/v1/organizations/organizationID-1/break-glass-members/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Domains.AddBreakGlassMember(ctx, "organizationID-1", "userID-1")
	}},
	{name: "Domains.RemoveBreakGlassMember", method: "DELETE", path: "/v1/organizations/organizationID-1/break-glass-members/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Domains.RemoveBreakGlassMember(ctx, "organizationID-1", "userID-1")
	}},
	{name: "Branding.Get", method: "GET", path: "/v1/organizations/organizationID-1/branding", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Branding.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "Branding.Update", method: "PUT", path: "/v1/organizations/organizationID-1/branding", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Branding.Update(ctx, "organizationID-1", Branding{})
		return err
	}},
	{name: "Roles.List", method: "GET", path: "/v1/organizations/organizationID-1/roles", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Roles.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Roles.Create", method: "POST", path: "/v1/organizations/organizationID-1/roles", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Roles.Create(ctx, "organizationID-1", CreateRoleParams{})
		return err
	}},
	{name: "Roles.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/roles/roleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Roles.Update(ctx, "organizationID-1", "roleID-1", UpdateRoleParams{})
		return err
	}},
	{name: "Roles.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/roles/roleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Roles.Delete(ctx, "organizationID-1", "roleID-1") }},
	{name: "Roles.ListPermissions", method: "GET", path: "/v1/roles/roleID-1/permissions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Roles.ListPermissions(ctx, "roleID-1")
		return err
	}},
	{name: "Roles.SetPermissions", method: "PUT", path: "/v1/organizations/organizationID-1/roles/roleID-1/permissions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Roles.SetPermissions(ctx, "organizationID-1", "roleID-1", SetRolePermissionsParams{})
		return err
	}},
	{name: "Roles.GrantPermission", method: "POST", path: "/v1/roles/roleID-1/permissions/permissionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Roles.GrantPermission(ctx, "roleID-1", "permissionID-1")
	}},
	{name: "Roles.RevokePermission", method: "DELETE", path: "/v1/roles/roleID-1/permissions/permissionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Roles.RevokePermission(ctx, "roleID-1", "permissionID-1")
	}},
	{name: "Roles.Holders", method: "GET", path: "/v1/organizations/organizationID-1/roles/roleID-1/holders", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Roles.Holders(ctx, "organizationID-1", "roleID-1", nil)
		return err
	}},
	{name: "Roles.AssignToUser", method: "POST", path: "/v1/users/userID-1/organizations/organizationID-1/roles/roleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Roles.AssignToUser(ctx, "userID-1", "organizationID-1", "roleID-1")
	}},
	{name: "Roles.UnassignFromUser", method: "DELETE", path: "/v1/users/userID-1/organizations/organizationID-1/roles/roleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Roles.UnassignFromUser(ctx, "userID-1", "organizationID-1", "roleID-1")
	}},
	{name: "Permissions.List", method: "GET", path: "/v1/permissions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Permissions.List(ctx); return err }},
	{name: "Permissions.Check", method: "GET", path: "/v1/check", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Permissions.Check(ctx, nil); return err }},
	{name: "Permissions.Effective", method: "GET", path: "/v1/users/userID-1/organizations/organizationID-1/permissions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Permissions.Effective(ctx, "userID-1", "organizationID-1")
		return err
	}},
	{name: "ABACPolicies.List", method: "GET", path: "/v1/organizations/organizationID-1/abac/policies", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.List(ctx, "organizationID-1")
		return err
	}},
	{name: "ABACPolicies.Create", method: "POST", path: "/v1/organizations/organizationID-1/abac/policies", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Create(ctx, "organizationID-1", ABACPolicyParams{})
		return err
	}},
	{name: "ABACPolicies.Get", method: "GET", path: "/v1/organizations/organizationID-1/abac/policies/policyID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Get(ctx, "organizationID-1", "policyID-1")
		return err
	}},
	{name: "ABACPolicies.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/abac/policies/policyID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Update(ctx, "organizationID-1", "policyID-1", ABACPolicyParams{})
		return err
	}},
	{name: "ABACPolicies.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/abac/policies/policyID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.ABACPolicies.Delete(ctx, "organizationID-1", "policyID-1")
	}},
	{name: "ABACPolicies.ListVersions", method: "GET", path: "/v1/organizations/organizationID-1/abac/policies/policyID-1/versions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.ListVersions(ctx, "organizationID-1", "policyID-1")
		return err
	}},
	{name: "ABACPolicies.RestoreVersion", method: "POST", path: "/v1/organizations/organizationID-1/abac/policies/policyID-1/versions/version-1/restore", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.RestoreVersion(ctx, "organizationID-1", "policyID-1", "version-1")
		return err
	}},
	{name: "ABACPolicies.Evaluate", method: "POST", path: "/v1/organizations/organizationID-1/abac/evaluate", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ABACPolicies.Evaluate(ctx, "organizationID-1", ABACEvaluateParams{}, nil)
		return err
	}},
	{name: "Relationships.List", method: "GET", path: "/v1/organizations/organizationID-1/relation-tuples", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.List(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Relationships.Create", method: "POST", path: "/v1/organizations/organizationID-1/relation-tuples", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.Create(ctx, "organizationID-1", CreateRelationTupleParams{})
		return err
	}},
	{name: "Relationships.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/relation-tuples/tupleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Relationships.Delete(ctx, "organizationID-1", "tupleID-1")
	}},
	{name: "Relationships.Check", method: "POST", path: "/v1/organizations/organizationID-1/relation-tuples/check", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.Check(ctx, "organizationID-1", CheckRelationParams{}, nil)
		return err
	}},
	{name: "Relationships.ListTypes", method: "GET", path: "/v1/organizations/organizationID-1/relation-tuples/types", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.ListTypes(ctx, "organizationID-1")
		return err
	}},
	{name: "Relationships.ListObjects", method: "GET", path: "/v1/organizations/organizationID-1/relation-tuples/objects", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.ListObjects(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Relationships.AllObjects", method: "GET", path: "/v1/organizations/organizationID-1/relation-tuples/objects", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Relationships.AllObjects(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Relationships.Graph", method: "GET", path: "/v1/organizations/organizationID-1/relation-tuples/graph", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Relationships.Graph(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "AuthZEN.Evaluate", method: "POST", path: "/v1/organizations/organizationID-1/access/v1/evaluation", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthZEN.Evaluate(ctx, "organizationID-1", EvaluationParams{})
		return err
	}},
	{name: "AuthZEN.EvaluateBatch", method: "POST", path: "/v1/organizations/organizationID-1/access/v1/evaluations", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthZEN.EvaluateBatch(ctx, "organizationID-1", BatchEvaluationParams{})
		return err
	}},
	{name: "AuthZEN.ListDecisions", method: "GET", path: "/v1/organizations/organizationID-1/authorization/decisions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthZEN.ListDecisions(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "AuthZEN.AllDecisions", method: "GET", path: "/v1/organizations/organizationID-1/authorization/decisions", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.AuthZEN.AllDecisions(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "AuthPolicy.Get", method: "GET", path: "/v1/organizations/organizationID-1/auth-policy", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthPolicy.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "AuthPolicy.Update", method: "PUT", path: "/v1/organizations/organizationID-1/auth-policy", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthPolicy.Update(ctx, "organizationID-1", AuthPolicy{})
		return err
	}},
	{name: "AuthPolicy.Usage", method: "GET", path: "/v1/organizations/organizationID-1/auth-policy/usage", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthPolicy.Usage(ctx, "organizationID-1")
		return err
	}},
	{name: "SecurityPolicy.Get", method: "GET", path: "/v1/organizations/organizationID-1/policy", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecurityPolicy.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "SecurityPolicy.Replace", method: "PUT", path: "/v1/organizations/organizationID-1/policy", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecurityPolicy.Replace(ctx, "organizationID-1", SecurityPolicy{})
		return err
	}},
	{name: "SecurityPolicy.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/policy", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SecurityPolicy.Update(ctx, "organizationID-1", UpdateSecurityPolicyParams{})
		return err
	}},
	{name: "IPRules.List", method: "GET", path: "/v1/organizations/organizationID-1/ip-rules", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.List(ctx, "organizationID-1")
		return err
	}},
	{name: "IPRules.Create", method: "POST", path: "/v1/organizations/organizationID-1/ip-rules", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.Create(ctx, "organizationID-1", CreateIPRuleParams{})
		return err
	}},
	{name: "IPRules.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/ip-rules/ruleID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.IPRules.Delete(ctx, "organizationID-1", "ruleID-1")
	}},
	{name: "IPRules.Check", method: "POST", path: "/v1/organizations/organizationID-1/ip-rules/check", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.Check(ctx, "organizationID-1", CheckIPParams{})
		return err
	}},
	{name: "IPRules.SetEnforcement", method: "PUT", path: "/v1/organizations/organizationID-1/ip-rules/config", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.IPRules.SetEnforcement(ctx, "organizationID-1", SetIPEnforcementParams{})
		return err
	}},
	{name: "Threats.Bots", method: "GET", path: "/v1/organizations/organizationID-1/security/bots", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.Bots(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Threats.BotSettings", method: "GET", path: "/v1/organizations/organizationID-1/security/bots/settings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.BotSettings(ctx, "organizationID-1")
		return err
	}},
	{name: "Threats.UpdateBotSettings", method: "PUT", path: "/v1/organizations/organizationID-1/security/bots/settings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.UpdateBotSettings(ctx, "organizationID-1", BotSettings{})
		return err
	}},
	{name: "Threats.RiskSettings", method: "GET", path: "/v1/organizations/organizationID-1/security/risk-settings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.RiskSettings(ctx, "organizationID-1")
		return err
	}},
	{name: "Threats.UpdateRiskSettings", method: "PUT", path: "/v1/organizations/organizationID-1/security/risk-settings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.UpdateRiskSettings(ctx, "organizationID-1", RiskSettings{})
		return err
	}},
	{name: "Threats.ListAnomalies", method: "GET", path: "/v1/organizations/organizationID-1/security/anomalies", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.ListAnomalies(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Threats.AllAnomalies", method: "GET", path: "/v1/organizations/organizationID-1/security/anomalies", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Threats.AllAnomalies(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Threats.AnomalySummary", method: "GET", path: "/v1/organizations/organizationID-1/security/anomalies/summary", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Threats.AnomalySummary(ctx, "organizationID-1")
		return err
	}},
	{name: "Threats.ResolveAnomaly", method: "POST", path: "/v1/organizations/organizationID-1/security/anomalies/anomalyID-1/resolve", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Threats.ResolveAnomaly(ctx, "organizationID-1", "anomalyID-1")
	}},
	{name: "OIDCClients.List", method: "GET", path: "/v1/organizations/organizationID-1/oidc/clients", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.List(ctx, "organizationID-1")
		return err
	}},
	{name: "OIDCClients.Create", method: "POST", path: "/v1/organizations/organizationID-1/oidc/clients", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.Create(ctx, "organizationID-1", CreateOIDCClientParams{})
		return err
	}},
	{name: "OIDCClients.Get", method: "GET", path: "/v1/organizations/organizationID-1/oidc/clients/clientID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.Get(ctx, "organizationID-1", "clientID-1")
		return err
	}},
	{name: "OIDCClients.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/oidc/clients/clientID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.Update(ctx, "organizationID-1", "clientID-1", UpdateOIDCClientParams{})
		return err
	}},
	{name: "OIDCClients.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/oidc/clients/clientID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.OIDCClients.Delete(ctx, "organizationID-1", "clientID-1")
	}},
	{name: "OIDCClients.RotateSecret", method: "POST", path: "/v1/organizations/organizationID-1/oidc/clients/clientID-1/rotate-secret", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.RotateSecret(ctx, "organizationID-1", "clientID-1")
		return err
	}},
	{name: "OIDCClients.AuthFailures", method: "GET", path: "/v1/organizations/organizationID-1/oidc/clients/clientID-1/auth-failures", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.AuthFailures(ctx, "organizationID-1", "clientID-1")
		return err
	}},
	{name: "OIDCClients.ListShadowAI", method: "GET", path: "/v1/organizations/organizationID-1/oidc/clients/shadow-ai", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCClients.ListShadowAI(ctx, "organizationID-1")
		return err
	}},
	{name: "OIDCClients.AcknowledgeShadowAI", method: "POST", path: "/v1/organizations/organizationID-1/oidc/clients/clientID-1/review", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.OIDCClients.AcknowledgeShadowAI(ctx, "organizationID-1", "clientID-1")
	}},
	{name: "OAuthResources.List", method: "GET", path: "/v1/organizations/organizationID-1/oidc/resources", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuthResources.List(ctx, "organizationID-1")
		return err
	}},
	{name: "OAuthResources.Create", method: "POST", path: "/v1/organizations/organizationID-1/oidc/resources", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuthResources.Create(ctx, "organizationID-1", CreateOAuthResourceParams{})
		return err
	}},
	{name: "OAuthResources.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/oidc/resources/resourceID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.OAuthResources.Delete(ctx, "organizationID-1", "resourceID-1")
	}},
	{name: "OAuthGrants.List", method: "GET", path: "/v1/organizations/organizationID-1/oauth/grants", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuthGrants.List(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "OAuthGrants.All", method: "GET", path: "/v1/organizations/organizationID-1/oauth/grants", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.OAuthGrants.All(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "OAuthGrants.Revoke", method: "DELETE", path: "/v1/organizations/organizationID-1/oauth/grants/grantID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.OAuthGrants.Revoke(ctx, "organizationID-1", "grantID-1")
	}},
	{name: "OAuthGrants.RevokeMany", method: "POST", path: "/v1/organizations/organizationID-1/oauth/grants/revoke", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OAuthGrants.RevokeMany(ctx, "organizationID-1", RevokeOAuthGrantsParams{})
		return err
	}},
	{name: "DeviceAuthorizations.List", method: "GET", path: "/v1/organizations/organizationID-1/oauth/devices", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.DeviceAuthorizations.List(ctx, "organizationID-1")
		return err
	}},
	{name: "DeviceAuthorizations.Revoke", method: "DELETE", path: "/v1/organizations/organizationID-1/oauth/devices/authorizationID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.DeviceAuthorizations.Revoke(ctx, "organizationID-1", "authorizationID-1")
	}},
	{name: "ServicePrincipals.Create", method: "POST", path: "/v1/service-principals", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.ServicePrincipals.Create(ctx, CreateServicePrincipalParams{})
		return err
	}},
	{name: "ServicePrincipals.List", method: "GET", path: "/v1/organizations/organizationID-1/service-principals", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ServicePrincipals.List(ctx, "organizationID-1")
		return err
	}},
	{name: "ServicePrincipals.Disable", method: "DELETE", path: "/v1/service-principals/principalID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.ServicePrincipals.Disable(ctx, "principalID-1") }},
	{name: "ServicePrincipals.Enable", method: "POST", path: "/v1/organizations/organizationID-1/service-principals/principalID-1/enable", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.ServicePrincipals.Enable(ctx, "organizationID-1", "principalID-1")
	}},
	{name: "ServicePrincipals.RotateSecret", method: "POST", path: "/v1/organizations/organizationID-1/service-principals/principalID-1/rotate-secret", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.ServicePrincipals.RotateSecret(ctx, "organizationID-1", "principalID-1", RotateSecretParams{})
		return err
	}},
	{name: "SocialProviders.List", method: "GET", path: "/v1/organizations/organizationID-1/social/providers", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SocialProviders.List(ctx, "organizationID-1")
		return err
	}},
	{name: "SocialProviders.Upsert", method: "POST", path: "/v1/social/providers", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SocialProviders.Upsert(ctx, UpsertSocialProviderParams{})
		return err
	}},
	{name: "SocialProviders.Update", method: "PATCH", path: "/v1/social/providers/provider-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SocialProviders.Update(ctx, "provider-1", UpdateSocialProviderParams{})
		return err
	}},
	{name: "SocialProviders.Delete", method: "DELETE", path: "/v1/social/providers/provider-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.SocialProviders.Delete(ctx, "provider-1") }},
	{name: "SocialProviders.UnlinkIdentity", method: "DELETE", path: "/v1/social/identities/identityID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SocialProviders.UnlinkIdentity(ctx, "identityID-1")
	}},
	{name: "SAMLConnections.List", method: "GET", path: "/v1/organizations/organizationID-1/saml", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.List(ctx, "organizationID-1")
		return err
	}},
	{name: "SAMLConnections.Create", method: "POST", path: "/v1/organizations/organizationID-1/saml", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Create(ctx, "organizationID-1", SAMLConnectionParams{})
		return err
	}},
	{name: "SAMLConnections.Import", method: "POST", path: "/v1/organizations/organizationID-1/saml/import", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Import(ctx, "organizationID-1", ImportSAMLConnectionParams{})
		return err
	}},
	{name: "SAMLConnections.Get", method: "GET", path: "/v1/organizations/organizationID-1/saml/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Get(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "SAMLConnections.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/saml/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Update(ctx, "organizationID-1", "connectionID-1", UpdateSAMLConnectionParams{})
		return err
	}},
	{name: "SAMLConnections.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/saml/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLConnections.Delete(ctx, "organizationID-1", "connectionID-1")
	}},
	{name: "SAMLConnections.ApplyMetadata", method: "POST", path: "/v1/organizations/organizationID-1/saml/connectionID-1/metadata", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.ApplyMetadata(ctx, "organizationID-1", "connectionID-1", ApplySAMLMetadataParams{})
		return err
	}},
	{name: "SAMLConnections.AddCertificate", method: "POST", path: "/v1/organizations/organizationID-1/saml/connectionID-1/certificates", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.AddCertificate(ctx, "organizationID-1", "connectionID-1", AddSAMLCertificateParams{})
		return err
	}},
	{name: "SAMLConnections.RemoveCertificate", method: "DELETE", path: "/v1/organizations/organizationID-1/saml/connectionID-1/certificates/certificateID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLConnections.RemoveCertificate(ctx, "organizationID-1", "connectionID-1", "certificateID-1")
	}},
	{name: "SAMLConnections.Test", method: "POST", path: "/v1/organizations/organizationID-1/saml/connectionID-1/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.Test(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "SAMLConnections.ListGroupMappings", method: "GET", path: "/v1/organizations/organizationID-1/saml/connectionID-1/group-mappings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.ListGroupMappings(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "SAMLConnections.CreateGroupMapping", method: "POST", path: "/v1/organizations/organizationID-1/saml/connectionID-1/group-mappings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLConnections.CreateGroupMapping(ctx, "organizationID-1", "connectionID-1", CreateGroupMappingParams{})
		return err
	}},
	{name: "SAMLConnections.DeleteGroupMapping", method: "DELETE", path: "/v1/organizations/organizationID-1/saml/connectionID-1/group-mappings/mappingID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLConnections.DeleteGroupMapping(ctx, "organizationID-1", "connectionID-1", "mappingID-1")
	}},
	{name: "SAMLServiceProviders.List", method: "GET", path: "/v1/organizations/organizationID-1/saml-providers", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.List(ctx, "organizationID-1")
		return err
	}},
	{name: "SAMLServiceProviders.Create", method: "POST", path: "/v1/organizations/organizationID-1/saml-providers", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.Create(ctx, "organizationID-1", SAMLServiceProviderParams{})
		return err
	}},
	{name: "SAMLServiceProviders.Get", method: "GET", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.Get(ctx, "organizationID-1", "providerID-1")
		return err
	}},
	{name: "SAMLServiceProviders.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.Update(ctx, "organizationID-1", "providerID-1", SAMLServiceProviderParams{})
		return err
	}},
	{name: "SAMLServiceProviders.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.Delete(ctx, "organizationID-1", "providerID-1")
	}},
	{name: "SAMLServiceProviders.ListAssignments", method: "GET", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1/assignments", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SAMLServiceProviders.ListAssignments(ctx, "organizationID-1", "providerID-1")
		return err
	}},
	{name: "SAMLServiceProviders.AssignUser", method: "POST", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1/assignments/users/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.AssignUser(ctx, "organizationID-1", "providerID-1", "userID-1")
	}},
	{name: "SAMLServiceProviders.UnassignUser", method: "DELETE", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1/assignments/users/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.UnassignUser(ctx, "organizationID-1", "providerID-1", "userID-1")
	}},
	{name: "SAMLServiceProviders.AssignGroup", method: "POST", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1/assignments/groups/groupID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.AssignGroup(ctx, "organizationID-1", "providerID-1", "groupID-1")
	}},
	{name: "SAMLServiceProviders.UnassignGroup", method: "DELETE", path: "/v1/organizations/organizationID-1/saml-providers/providerID-1/assignments/groups/groupID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SAMLServiceProviders.UnassignGroup(ctx, "organizationID-1", "providerID-1", "groupID-1")
	}},
	{name: "OIDCConnections.List", method: "GET", path: "/v1/organizations/organizationID-1/oidc-connections", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.List(ctx, "organizationID-1")
		return err
	}},
	{name: "OIDCConnections.Create", method: "POST", path: "/v1/organizations/organizationID-1/oidc-connections", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.Create(ctx, "organizationID-1", CreateOIDCConnectionParams{})
		return err
	}},
	{name: "OIDCConnections.Get", method: "GET", path: "/v1/organizations/organizationID-1/oidc-connections/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.Get(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "OIDCConnections.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/oidc-connections/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.OIDCConnections.Update(ctx, "organizationID-1", "connectionID-1", UpdateOIDCConnectionParams{})
		return err
	}},
	{name: "OIDCConnections.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/oidc-connections/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.OIDCConnections.Delete(ctx, "organizationID-1", "connectionID-1")
	}},
	{name: "LDAP.List", method: "GET", path: "/v1/organizations/organizationID-1/ldap", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.LDAP.List(ctx, "organizationID-1"); return err }},
	{name: "LDAP.Create", method: "POST", path: "/v1/organizations/organizationID-1/ldap", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Create(ctx, "organizationID-1", CreateLDAPConnectionParams{})
		return err
	}},
	{name: "LDAP.Get", method: "GET", path: "/v1/organizations/organizationID-1/ldap/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Get(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "LDAP.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/ldap/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Update(ctx, "organizationID-1", "connectionID-1", UpdateLDAPConnectionParams{})
		return err
	}},
	{name: "LDAP.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/ldap/connectionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.LDAP.Delete(ctx, "organizationID-1", "connectionID-1")
	}},
	{name: "LDAP.Test", method: "POST", path: "/v1/organizations/organizationID-1/ldap/connectionID-1/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.Test(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "LDAP.ListGroupMappings", method: "GET", path: "/v1/organizations/organizationID-1/ldap/connectionID-1/group-mappings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.ListGroupMappings(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "LDAP.CreateGroupMapping", method: "POST", path: "/v1/organizations/organizationID-1/ldap/connectionID-1/group-mappings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LDAP.CreateGroupMapping(ctx, "organizationID-1", "connectionID-1", CreateGroupMappingParams{})
		return err
	}},
	{name: "LDAP.DeleteGroupMapping", method: "DELETE", path: "/v1/organizations/organizationID-1/ldap/connectionID-1/group-mappings/mappingID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.LDAP.DeleteGroupMapping(ctx, "organizationID-1", "connectionID-1", "mappingID-1")
	}},
	{name: "Directory.ListConnections", method: "GET", path: "/v1/organizations/organizationID-1/directory/connections", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListConnections(ctx, "organizationID-1")
		return err
	}},
	{name: "Directory.ListMappings", method: "GET", path: "/v1/organizations/organizationID-1/directory/connections/connectionID-1/mappings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListMappings(ctx, "organizationID-1", "connectionID-1")
		return err
	}},
	{name: "Directory.CreateMapping", method: "POST", path: "/v1/organizations/organizationID-1/directory/connections/connectionID-1/mappings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.CreateMapping(ctx, "organizationID-1", "connectionID-1", DirectoryMapping{})
		return err
	}},
	{name: "Directory.UpdateMapping", method: "PUT", path: "/v1/organizations/organizationID-1/directory/connections/connectionID-1/mappings/mappingID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.UpdateMapping(ctx, "organizationID-1", "connectionID-1", "mappingID-1", DirectoryMapping{})
		return err
	}},
	{name: "Directory.DeleteMapping", method: "DELETE", path: "/v1/organizations/organizationID-1/directory/connections/connectionID-1/mappings/mappingID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Directory.DeleteMapping(ctx, "organizationID-1", "connectionID-1", "mappingID-1")
	}},
	{name: "Directory.PreviewMappings", method: "POST", path: "/v1/organizations/organizationID-1/directory/connections/connectionID-1/mappings/preview", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.PreviewMappings(ctx, "organizationID-1", "connectionID-1", PreviewMappingsParams{})
		return err
	}},
	{name: "Directory.ListErrors", method: "GET", path: "/v1/organizations/organizationID-1/directory/errors", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListErrors(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Directory.ResolveError", method: "POST", path: "/v1/organizations/organizationID-1/directory/errors/eventID-1/resolve", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ResolveError(ctx, "organizationID-1", "eventID-1")
		return err
	}},
	{name: "Directory.RetryError", method: "POST", path: "/v1/organizations/organizationID-1/directory/errors/eventID-1/retry", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.RetryError(ctx, "organizationID-1", "eventID-1")
		return err
	}},
	{name: "Directory.ListActivity", method: "GET", path: "/v1/organizations/organizationID-1/directory/activity", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListActivity(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Directory.ListAlerts", method: "GET", path: "/v1/organizations/organizationID-1/directory/alerts", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.ListAlerts(ctx, "organizationID-1")
		return err
	}},
	{name: "Directory.CreateAlert", method: "POST", path: "/v1/organizations/organizationID-1/directory/alerts", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.CreateAlert(ctx, "organizationID-1", CreateDirectoryAlertParams{})
		return err
	}},
	{name: "Directory.UpdateAlert", method: "PATCH", path: "/v1/organizations/organizationID-1/directory/alerts/alertID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.UpdateAlert(ctx, "organizationID-1", "alertID-1", UpdateDirectoryAlertParams{})
		return err
	}},
	{name: "Directory.DeleteAlert", method: "DELETE", path: "/v1/organizations/organizationID-1/directory/alerts/alertID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Directory.DeleteAlert(ctx, "organizationID-1", "alertID-1")
	}},
	{name: "Directory.Profile", method: "GET", path: "/v1/organizations/organizationID-1/directory/profiles/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Directory.Profile(ctx, "organizationID-1", "userID-1")
		return err
	}},
	{name: "SCIM.Config", method: "GET", path: "/v1/organizations/organizationID-1/scim", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.Config(ctx, "organizationID-1")
		return err
	}},
	{name: "SCIM.ListTokens", method: "GET", path: "/v1/organizations/organizationID-1/scim/tokens", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.ListTokens(ctx, "organizationID-1")
		return err
	}},
	{name: "SCIM.CreateToken", method: "POST", path: "/v1/organizations/organizationID-1/scim/tokens", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.CreateToken(ctx, "organizationID-1", CreateSCIMTokenParams{})
		return err
	}},
	{name: "SCIM.RotateToken", method: "POST", path: "/v1/organizations/organizationID-1/scim/tokens/tokenID-1/rotate", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.RotateToken(ctx, "organizationID-1", "tokenID-1", RotateSCIMTokenParams{})
		return err
	}},
	{name: "SCIM.RevokeToken", method: "DELETE", path: "/v1/organizations/organizationID-1/scim/tokens/tokenID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.SCIM.RevokeToken(ctx, "organizationID-1", "tokenID-1")
	}},
	{name: "SCIM.ReplaceAllTokens", method: "POST", path: "/v1/organizations/organizationID-1/scim/token", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.ReplaceAllTokens(ctx, "organizationID-1")
		return err
	}},
	{name: "SCIM.RevokeAllTokens", method: "DELETE", path: "/v1/organizations/organizationID-1/scim/token", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.SCIM.RevokeAllTokens(ctx, "organizationID-1") }},
	{name: "SCIM.ListUsers", method: "GET", path: "/v1/organizations/organizationID-1/scim/users", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.ListUsers(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "SCIM.AllUsers", method: "GET", path: "/v1/organizations/organizationID-1/scim/users", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.SCIM.AllUsers(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "SCIM.ListGroups", method: "GET", path: "/v1/organizations/organizationID-1/scim/groups", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.SCIM.ListGroups(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "SCIM.AllGroups", method: "GET", path: "/v1/organizations/organizationID-1/scim/groups", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.SCIM.AllGroups(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "AdminPortalLinks.Create", method: "POST", path: "/v1/organizations/organizationID-1/admin-portal/links", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AdminPortalLinks.Create(ctx, "organizationID-1", CreateAdminPortalLinkParams{})
		return err
	}},
	{name: "AdminPortalLinks.List", method: "GET", path: "/v1/organizations/organizationID-1/admin-portal/links", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AdminPortalLinks.List(ctx, "organizationID-1")
		return err
	}},
	{name: "AdminPortalLinks.Revoke", method: "DELETE", path: "/v1/organizations/organizationID-1/admin-portal/links/linkID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.AdminPortalLinks.Revoke(ctx, "organizationID-1", "linkID-1")
	}},
	{name: "APIKeys.Create", method: "POST", path: "/v1/api-keys", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.APIKeys.Create(ctx, CreateAPIKeyParams{})
		return err
	}},
	{name: "APIKeys.List", method: "GET", path: "/v1/organizations/organizationID-1/api-keys", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.APIKeys.List(ctx, "organizationID-1")
		return err
	}},
	{name: "APIKeys.Rotate", method: "POST", path: "/v1/organizations/organizationID-1/api-keys/keyID-1/rotate", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.APIKeys.Rotate(ctx, "organizationID-1", "keyID-1", RotateSecretParams{})
		return err
	}},
	{name: "APIKeys.Revoke", method: "DELETE", path: "/v1/api-keys/keyID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.APIKeys.Revoke(ctx, "keyID-1") }},
	{name: "Agents.Create", method: "POST", path: "/v1/organizations/organizationID-1/agents", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Create(ctx, "organizationID-1", CreateAgentParams{})
		return err
	}},
	{name: "Agents.List", method: "GET", path: "/v1/organizations/organizationID-1/agents", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Agents.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/agents/agentID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Agents.Update(ctx, "organizationID-1", "agentID-1", UpdateAgentParams{})
	}},
	{name: "Agents.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/agents/agentID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Agents.Delete(ctx, "organizationID-1", "agentID-1")
	}},
	{name: "Agents.Suspend", method: "POST", path: "/v1/organizations/organizationID-1/agents/agentID-1/suspend", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Suspend(ctx, "organizationID-1", "agentID-1")
		return err
	}},
	{name: "Agents.Resume", method: "POST", path: "/v1/organizations/organizationID-1/agents/agentID-1/resume", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Resume(ctx, "organizationID-1", "agentID-1")
		return err
	}},
	{name: "Agents.Decommission", method: "POST", path: "/v1/organizations/organizationID-1/agents/agentID-1/decommission", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.Decommission(ctx, "organizationID-1", "agentID-1")
		return err
	}},
	{name: "Agents.SuspendAll", method: "POST", path: "/v1/organizations/organizationID-1/agents/kill-all", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.SuspendAll(ctx, "organizationID-1")
		return err
	}},
	{name: "Agents.ResumeAll", method: "POST", path: "/v1/organizations/organizationID-1/agents/resume-all", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.ResumeAll(ctx, "organizationID-1")
		return err
	}},
	{name: "Agents.RotateSecret", method: "POST", path: "/v1/organizations/organizationID-1/agents/agentID-1/rotate-secret", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.RotateSecret(ctx, "organizationID-1", "agentID-1", RotateSecretParams{})
		return err
	}},
	{name: "Agents.ListSponsoredBy", method: "GET", path: "/v1/organizations/organizationID-1/agents/sponsored-by/userID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.ListSponsoredBy(ctx, "organizationID-1", "userID-1")
		return err
	}},
	{name: "Agents.TransferSponsorship", method: "POST", path: "/v1/organizations/organizationID-1/agents/sponsored-by/userID-1/transfer", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Agents.TransferSponsorship(ctx, "organizationID-1", "userID-1", TransferSponsorshipParams{})
		return err
	}},
	{name: "Credentials.Issue", method: "POST", path: "/v1/organizations/organizationID-1/credentials", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Credentials.Issue(ctx, "organizationID-1", IssueCredentialParams{})
		return err
	}},
	{name: "Credentials.List", method: "GET", path: "/v1/organizations/organizationID-1/credentials", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Credentials.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Credentials.Revoke", method: "POST", path: "/v1/organizations/organizationID-1/credentials/credentialID-1/revoke", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Credentials.Revoke(ctx, "organizationID-1", "credentialID-1")
	}},
	{name: "Credentials.Verify", method: "POST", path: "/v1/credentials/verify", public: true, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Credentials.Verify(ctx, VerifyCredentialParams{})
		return err
	}},
	{name: "AuthHooks.List", method: "GET", path: "/v1/organizations/organizationID-1/auth-hooks", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.List(ctx, "organizationID-1")
		return err
	}},
	{name: "AuthHooks.Create", method: "POST", path: "/v1/organizations/organizationID-1/auth-hooks", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.Create(ctx, "organizationID-1", CreateAuthHookParams{})
		return err
	}},
	{name: "AuthHooks.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/auth-hooks/hookID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.Update(ctx, "organizationID-1", "hookID-1", UpdateAuthHookParams{})
		return err
	}},
	{name: "AuthHooks.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/auth-hooks/hookID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.AuthHooks.Delete(ctx, "organizationID-1", "hookID-1")
	}},
	{name: "AuthHooks.Test", method: "POST", path: "/v1/organizations/organizationID-1/auth-hooks/hookID-1/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.AuthHooks.Test(ctx, "organizationID-1", "hookID-1")
		return err
	}},
	{name: "Webhooks.Create", method: "POST", path: "/v1/webhooks", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.Create(ctx, CreateWebhookParams{})
		return err
	}},
	{name: "Webhooks.Get", method: "GET", path: "/v1/webhooks/webhookID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Webhooks.Get(ctx, "webhookID-1"); return err }},
	{name: "Webhooks.List", method: "GET", path: "/v1/organizations/organizationID-1/webhooks", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Webhooks.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/webhooks/webhookID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.Update(ctx, "organizationID-1", "webhookID-1", UpdateWebhookParams{})
		return err
	}},
	{name: "Webhooks.Disable", method: "DELETE", path: "/v1/webhooks/webhookID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.Webhooks.Disable(ctx, "webhookID-1") }},
	{name: "Webhooks.Enable", method: "POST", path: "/v1/organizations/organizationID-1/webhooks/webhookID-1/enable", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.Enable(ctx, "organizationID-1", "webhookID-1")
		return err
	}},
	{name: "Webhooks.RotateSecret", method: "POST", path: "/v1/organizations/organizationID-1/webhooks/webhookID-1/rotate-secret", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.RotateSecret(ctx, "organizationID-1", "webhookID-1")
		return err
	}},
	{name: "Webhooks.Test", method: "POST", path: "/v1/webhooks/webhookID-1/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Webhooks.Test(ctx, "webhookID-1"); return err }},
	{name: "Webhooks.ListEventTypes", method: "GET", path: "/v1/organizations/organizationID-1/webhooks/events", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.ListEventTypes(ctx, "organizationID-1")
		return err
	}},
	{name: "Webhooks.ListDeliveries", method: "GET", path: "/v1/webhooks/webhookID-1/deliveries", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.ListDeliveries(ctx, "webhookID-1", nil)
		return err
	}},
	{name: "Webhooks.AllDeliveries", method: "GET", path: "/v1/webhooks/webhookID-1/deliveries", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Webhooks.AllDeliveries(ctx, "webhookID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Webhooks.ListAttempts", method: "GET", path: "/v1/organizations/organizationID-1/webhooks/webhookID-1/deliveries/deliveryID-1/attempts", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.ListAttempts(ctx, "organizationID-1", "webhookID-1", "deliveryID-1")
		return err
	}},
	{name: "Webhooks.RetryDelivery", method: "POST", path: "/v1/webhooks/webhookID-1/deliveries/deliveryID-1/retry", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Webhooks.RetryDelivery(ctx, "webhookID-1", "deliveryID-1")
	}},
	{name: "Webhooks.RetryGivenUp", method: "POST", path: "/v1/organizations/organizationID-1/webhooks/webhookID-1/deliveries/retry", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Webhooks.RetryGivenUp(ctx, "organizationID-1", "webhookID-1")
		return err
	}},
	{name: "Secrets.List", method: "GET", path: "/v1/organizations/organizationID-1/secrets", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.List(ctx, "organizationID-1")
		return err
	}},
	{name: "Secrets.Create", method: "POST", path: "/v1/organizations/organizationID-1/secrets", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.Create(ctx, "organizationID-1", CreateSecretParams{})
		return err
	}},
	{name: "Secrets.Update", method: "PATCH", path: "/v1/organizations/organizationID-1/secrets/secretID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.Update(ctx, "organizationID-1", "secretID-1", UpdateSecretParams{})
		return err
	}},
	{name: "Secrets.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/secrets/secretID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Secrets.Delete(ctx, "organizationID-1", "secretID-1")
	}},
	{name: "Secrets.Reveal", method: "POST", path: "/v1/organizations/organizationID-1/secrets/secretID-1/reveal", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Secrets.Reveal(ctx, "organizationID-1", "secretID-1")
		return err
	}},
	{name: "Secrets.GetByName", method: "GET", path: "/v1/vault/name-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Secrets.GetByName(ctx, "name-1"); return err }},
	{name: "TokenVault.RegisterProvider", method: "POST", path: "/v1/organizations/organizationID-1/vault/tokens/providers", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.RegisterProvider(ctx, "organizationID-1", RegisterTokenVaultProviderParams{})
		return err
	}},
	{name: "TokenVault.ListProviders", method: "GET", path: "/v1/organizations/organizationID-1/vault/tokens/providers", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.ListProviders(ctx, "organizationID-1")
		return err
	}},
	{name: "TokenVault.DeleteProvider", method: "DELETE", path: "/v1/organizations/organizationID-1/vault/tokens/providers/provider-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.TokenVault.DeleteProvider(ctx, "organizationID-1", "provider-1")
	}},
	{name: "TokenVault.ListConnections", method: "GET", path: "/v1/vault/tokens", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.TokenVault.ListConnections(ctx); return err }},
	{name: "TokenVault.Connect", method: "GET", path: "/v1/vault/tokens/provider-1/connect", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.Connect(ctx, "provider-1")
		return err
	}},
	{name: "TokenVault.AccessToken", method: "GET", path: "/v1/vault/tokens/provider-1/access-token", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.TokenVault.AccessToken(ctx, "provider-1")
		return err
	}},
	{name: "TokenVault.Disconnect", method: "DELETE", path: "/v1/vault/tokens/provider-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.TokenVault.Disconnect(ctx, "provider-1") }},
	{name: "Audit.List", method: "GET", path: "/v1/organizations/organizationID-1/audit", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.List(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Audit.All", method: "GET", path: "/v1/organizations/organizationID-1/audit", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Audit.All(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Audit.Verify", method: "GET", path: "/v1/organizations/organizationID-1/audit/verify", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.Verify(ctx, "organizationID-1")
		return err
	}},
	{name: "Audit.LastVerification", method: "GET", path: "/v1/organizations/organizationID-1/audit/verification", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.LastVerification(ctx, "organizationID-1")
		return err
	}},
	{name: "Audit.ListAnomalies", method: "GET", path: "/v1/organizations/organizationID-1/audit/anomalies", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.ListAnomalies(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Audit.AllAnomalies", method: "GET", path: "/v1/organizations/organizationID-1/audit/anomalies", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.Audit.AllAnomalies(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "Audit.AnomalySummary", method: "GET", path: "/v1/organizations/organizationID-1/audit/anomalies/summary", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.AnomalySummary(ctx, "organizationID-1")
		return err
	}},
	{name: "Audit.ResolveAnomaly", method: "POST", path: "/v1/organizations/organizationID-1/audit/anomalies/anomalyID-1/resolve", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Audit.ResolveAnomaly(ctx, "organizationID-1", "anomalyID-1")
	}},
	{name: "Audit.AnomalySettings", method: "GET", path: "/v1/organizations/organizationID-1/audit/anomaly-settings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.AnomalySettings(ctx, "organizationID-1")
		return err
	}},
	{name: "Audit.UpdateAnomalySettings", method: "PUT", path: "/v1/organizations/organizationID-1/audit/anomaly-settings", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Audit.UpdateAnomalySettings(ctx, "organizationID-1", UpdateAuditAnomalySettingsParams{})
		return err
	}},
	{name: "Analytics.Overview", method: "GET", path: "/v1/organizations/organizationID-1/analytics/overview", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Analytics.Overview(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Compliance.GenerateEvidence", method: "POST", path: "/v1/organizations/organizationID-1/compliance/framework-1/evidence", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.GenerateEvidence(ctx, "organizationID-1", "framework-1")
		return err
	}},
	{name: "Compliance.ListEvidence", method: "GET", path: "/v1/organizations/organizationID-1/compliance/framework-1/evidence", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.ListEvidence(ctx, "organizationID-1", "framework-1")
		return err
	}},
	{name: "Compliance.GetEvidence", method: "GET", path: "/v1/organizations/organizationID-1/compliance/evidence/runID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.GetEvidence(ctx, "organizationID-1", "runID-1")
		return err
	}},
	{name: "Compliance.ListControlOwners", method: "GET", path: "/v1/organizations/organizationID-1/compliance/controls/owners", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.ListControlOwners(ctx, "organizationID-1")
		return err
	}},
	{name: "Compliance.SetControlOwner", method: "PUT", path: "/v1/organizations/organizationID-1/compliance/controls/controlID-1/owner", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Compliance.SetControlOwner(ctx, "organizationID-1", "controlID-1", SetControlOwnerParams{})
	}},
	{name: "Compliance.RemoveControlOwner", method: "DELETE", path: "/v1/organizations/organizationID-1/compliance/controls/controlID-1/owner", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Compliance.RemoveControlOwner(ctx, "organizationID-1", "controlID-1")
	}},
	{name: "Compliance.ListAudits", method: "GET", path: "/v1/organizations/organizationID-1/compliance/audits", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.ListAudits(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "Compliance.CreateAudit", method: "POST", path: "/v1/organizations/organizationID-1/compliance/audits", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Compliance.CreateAudit(ctx, "organizationID-1", CreateExternalAuditParams{})
		return err
	}},
	{name: "Compliance.DeleteAudit", method: "DELETE", path: "/v1/organizations/organizationID-1/compliance/audits/auditID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.Compliance.DeleteAudit(ctx, "organizationID-1", "auditID-1")
	}},
	{name: "Retention.Get", method: "GET", path: "/v1/organizations/organizationID-1/retention", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Retention.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "Retention.Update", method: "PUT", path: "/v1/organizations/organizationID-1/retention", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Retention.Update(ctx, "organizationID-1", RetentionPolicy{})
		return err
	}},
	{name: "Retention.Preview", method: "POST", path: "/v1/organizations/organizationID-1/retention/preview", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Retention.Preview(ctx, "organizationID-1")
		return err
	}},
	{name: "Retention.Run", method: "POST", path: "/v1/organizations/organizationID-1/retention/run", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Retention.Run(ctx, "organizationID-1")
		return err
	}},
	{name: "GDPR.CreatePurge", method: "POST", path: "/v1/gdpr/purge", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.CreatePurge(ctx, CreatePurgeRequestParams{})
		return err
	}},
	{name: "GDPR.CancelPurge", method: "DELETE", path: "/v1/gdpr/purge/requestID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.GDPR.CancelPurge(ctx, "requestID-1") }},
	{name: "GDPR.ListPurges", method: "GET", path: "/v1/organizations/organizationID-1/gdpr/purge", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.ListPurges(ctx, "organizationID-1")
		return err
	}},
	{name: "GDPR.CreateExport", method: "POST", path: "/v1/gdpr/export", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.CreateExport(ctx, CreateExportRequestParams{})
		return err
	}},
	{name: "GDPR.ListExports", method: "GET", path: "/v1/organizations/organizationID-1/gdpr/export", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.ListExports(ctx, "organizationID-1")
		return err
	}},
	{name: "GDPR.GetExport", method: "GET", path: "/v1/organizations/organizationID-1/gdpr/export/requestID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.GDPR.GetExport(ctx, "organizationID-1", "requestID-1")
		return err
	}},
	{name: "LogSinks.List", method: "GET", path: "/v1/organizations/organizationID-1/log-sinks", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.List(ctx, "organizationID-1")
		return err
	}},
	{name: "LogSinks.Create", method: "POST", path: "/v1/organizations/organizationID-1/log-sinks", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.Create(ctx, "organizationID-1", LogSinkParams{})
		return err
	}},
	{name: "LogSinks.SetEnabled", method: "PATCH", path: "/v1/organizations/organizationID-1/log-sinks/sinkID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.LogSinks.SetEnabled(ctx, "organizationID-1", "sinkID-1", SetLogSinkEnabledParams{})
	}},
	{name: "LogSinks.Delete", method: "DELETE", path: "/v1/organizations/organizationID-1/log-sinks/sinkID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.LogSinks.Delete(ctx, "organizationID-1", "sinkID-1")
	}},
	{name: "LogSinks.Test", method: "POST", path: "/v1/organizations/organizationID-1/log-sinks/sinkID-1/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.Test(ctx, "organizationID-1", "sinkID-1")
		return err
	}},
	{name: "LogSinks.TestSettings", method: "POST", path: "/v1/organizations/organizationID-1/log-sinks/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.LogSinks.TestSettings(ctx, "organizationID-1", LogSinkParams{})
		return err
	}},
	{name: "RateLimits.Get", method: "GET", path: "/v1/organizations/organizationID-1/rate-limits", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.RateLimits.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "RateLimits.Update", method: "PUT", path: "/v1/organizations/organizationID-1/rate-limits", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.RateLimits.Update(ctx, "organizationID-1", UpdateRateLimitsParams{})
		return err
	}},
	{name: "RateLimits.Reset", method: "DELETE", path: "/v1/organizations/organizationID-1/rate-limits", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { return c.RateLimits.Reset(ctx, "organizationID-1") }},
	{name: "EmailTemplates.List", method: "GET", path: "/v1/organizations/organizationID-1/email-templates", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.List(ctx, "organizationID-1")
		return err
	}},
	{name: "EmailTemplates.Get", method: "GET", path: "/v1/organizations/organizationID-1/email-templates/key-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Get(ctx, "organizationID-1", "key-1")
		return err
	}},
	{name: "EmailTemplates.Update", method: "PUT", path: "/v1/organizations/organizationID-1/email-templates/key-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Update(ctx, "organizationID-1", "key-1", UpdateEmailTemplateParams{})
		return err
	}},
	{name: "EmailTemplates.Reset", method: "DELETE", path: "/v1/organizations/organizationID-1/email-templates/key-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Reset(ctx, "organizationID-1", "key-1")
		return err
	}},
	{name: "EmailTemplates.Preview", method: "POST", path: "/v1/organizations/organizationID-1/email-templates/key-1/preview", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.Preview(ctx, "organizationID-1", "key-1", PreviewEmailTemplateParams{})
		return err
	}},
	{name: "EmailTemplates.SendTest", method: "POST", path: "/v1/organizations/organizationID-1/email-templates/key-1/test", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailTemplates.SendTest(ctx, "organizationID-1", "key-1", SendTestEmailParams{})
		return err
	}},
	{name: "EmailDelivery.ListMessages", method: "GET", path: "/v1/organizations/organizationID-1/email/messages", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailDelivery.ListMessages(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "EmailDelivery.ListSuppressions", method: "GET", path: "/v1/organizations/organizationID-1/email/suppressions", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.EmailDelivery.ListSuppressions(ctx, "organizationID-1", nil)
		return err
	}},
	{name: "EmailDelivery.AllSuppressions", method: "GET", path: "/v1/organizations/organizationID-1/email/suppressions", public: false, call: func(ctx context.Context, c *Client) error {
		for _, err := range c.EmailDelivery.AllSuppressions(ctx, "organizationID-1", nil) {
			if err != nil {
				return err
			}
		}
		return nil
	}},
	{name: "EmailDelivery.DeleteSuppression", method: "DELETE", path: "/v1/organizations/organizationID-1/email/suppressions/suppressionID-1", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		return c.EmailDelivery.DeleteSuppression(ctx, "organizationID-1", "suppressionID-1")
	}},
	{name: "Billing.ListPlans", method: "GET", path: "/v1/billing/plans", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.Billing.ListPlans(ctx); return err }},
	{name: "Billing.Subscription", method: "GET", path: "/v1/organizations/organizationID-1/billing/subscription", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Subscription(ctx, "organizationID-1")
		return err
	}},
	{name: "Billing.Subscribe", method: "PUT", path: "/v1/organizations/organizationID-1/billing/subscription", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Subscribe(ctx, "organizationID-1", SubscribeParams{})
		return err
	}},
	{name: "Billing.CancelSubscription", method: "POST", path: "/v1/organizations/organizationID-1/billing/subscription/cancel", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.CancelSubscription(ctx, "organizationID-1")
		return err
	}},
	{name: "Billing.StartTrial", method: "POST", path: "/v1/organizations/organizationID-1/billing/trial", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.StartTrial(ctx, "organizationID-1", StartTrialParams{})
		return err
	}},
	{name: "Billing.Checkout", method: "POST", path: "/v1/organizations/organizationID-1/billing/checkout", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Checkout(ctx, "organizationID-1", CheckoutParams{})
		return err
	}},
	{name: "Billing.ListInvoices", method: "GET", path: "/v1/organizations/organizationID-1/billing/invoices", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.ListInvoices(ctx, "organizationID-1")
		return err
	}},
	{name: "Billing.PayInvoice", method: "POST", path: "/v1/organizations/organizationID-1/billing/invoices/invoiceID-1/pay", public: false, idempotent: true, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.PayInvoice(ctx, "organizationID-1", "invoiceID-1", PayInvoiceParams{})
		return err
	}},
	{name: "Billing.PaymentMethod", method: "GET", path: "/v1/organizations/organizationID-1/billing/payment-method", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.PaymentMethod(ctx, "organizationID-1")
		return err
	}},
	{name: "Billing.Profile", method: "GET", path: "/v1/organizations/organizationID-1/billing/profile", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.Profile(ctx, "organizationID-1")
		return err
	}},
	{name: "Billing.UpdateProfile", method: "PUT", path: "/v1/organizations/organizationID-1/billing/profile", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Billing.UpdateProfile(ctx, "organizationID-1", BillingProfile{})
		return err
	}},
	{name: "Entitlements.Get", method: "GET", path: "/v1/organizations/organizationID-1/entitlements", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Entitlements.Get(ctx, "organizationID-1")
		return err
	}},
	{name: "Entitlements.Usage", method: "GET", path: "/v1/organizations/organizationID-1/entitlements/usage", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error {
		_, err := c.Entitlements.Usage(ctx, "organizationID-1")
		return err
	}},
	{name: "SigningKeys.List", method: "GET", path: "/v1/oidc/signing-keys", public: false, idempotent: false, call: func(ctx context.Context, c *Client) error { _, err := c.SigningKeys.List(ctx); return err }},
}
