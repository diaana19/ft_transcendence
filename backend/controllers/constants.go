/*
** File: constants.go
** Description: Shared response-message constants for the controllers package.
** Responsibilities:
** - Centralise repeated user-facing error/status strings so they stay
**   consistent across handlers.
 */

package controllers

const (
	msgPostNotFound   = "Post not found"
	msgUserNotFound   = "User not found"
	msgInvalid2FACode = "invalid 2FA code"
)
