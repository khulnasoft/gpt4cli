---
sidebar_position: 12
sidebar_label: Collaboration / Orgs
---

# Collaboration and Orgs

While so far Gpt4cli is mainly focused on a single-user experience, we plan to add features for sharing, collaboration, and team management in the future, and some groundwork has already been done. **Orgs** are the basis for collaboration in Gpt4cli.

## Multiple Users

Orgs are helpful already if you have multiple users using Gpt4cli in the same project. Because Gpt4cli outputs a `.gpt4cli` file containing a bit of non-sensitive config data in each directory a plan is created in, you'll have problems with multiple users unless you either get each user into the same org or put `.gpt4cli` in your `.gitignore` file. Otherwise, each user will overwrite other users' `.gpt4cli` files on every push, and no one will be happy.

## Domain Access

When starting out with Gpt4cli and creating a new org, you have the option of automatically granting access to anyone with an email address on your domain.

## Invitations

If you choose not to grant access to your whole domain, or you want to invite someone from outside your email domain, you can use `gpt4cli invite`:

```bash
gpt4cli invite
```

## Joining an Org

To join an org you've been invited to, use `gpt4cli sign-in`:

```bash
gpt4cli sign-in
```

## Listing Users and Invites

To list users and pending invites, use `gpt4cli users`:

```bash
gpt4cli users
```

## Revoking Users and Invites

To revoke an invite or remove a user, use `gpt4cli revoke`:

```bash
gpt4cli revoke
```
