package header

import (
	"context"
	"net/textproto"
	"strings"
)

var (
	HeaderAsertoTenantID  = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Tenant-Id"))
	HeaderAsertoAccountID = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Account-Id"))
	HeaderAsertoTenantKey = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Api-Key"))
	HeaderAsertoRequestID = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Request-Id"))
	HeaderAsertoSessionID = CtxKey(textproto.CanonicalMIMEHeaderKey("Aserto-Session-Id"))

	HeaderAsertoTenantIDLowercase  = CtxKey(strings.ToLower(string(HeaderAsertoTenantID)))
	HeaderAsertoAccountIDLowercase = CtxKey(strings.ToLower(string(HeaderAsertoAccountID)))
	HeaderAsertoTenantKeyLowercase = CtxKey(strings.ToLower(string(HeaderAsertoTenantKey)))
	HeaderAsertoRequestIDLowercase = CtxKey(strings.ToLower(string(HeaderAsertoRequestID)))
	HeaderAsertoSessionIDLowercase = CtxKey(strings.ToLower(string(HeaderAsertoSessionID)))
)

func ContextWithAccountID(ctx context.Context, accountID string) context.Context {
	return context.WithValue(ctx, HeaderAsertoAccountID, accountID)
}

func ContextWithTenantID(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, HeaderAsertoTenantID, tenantID)
}

func ContextWithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, HeaderAsertoRequestID, requestID)
}

func ContextWithSessionID(ctx context.Context, sessionID string) context.Context {
	return context.WithValue(ctx, HeaderAsertoSessionID, sessionID)
}

func ExtractAccountID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoAccountID).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractTenantID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoTenantID).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractTenantKey(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoTenantKey).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractRequestID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoRequestID).(string)
	if !ok {
		return ""
	}

	return id
}

func ExtractSessionID(ctx context.Context) string {
	id, ok := ctx.Value(HeaderAsertoSessionID).(string)
	if !ok {
		return ""
	}

	return id
}
