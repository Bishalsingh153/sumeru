package orm_test

import (
	"context"
	"strings"
	"testing"

	"sumeru/core/orm"
)

func TestBeginOwnTOTPEnrollment_requiresAuth(t *testing.T) {
	_, err := orm.BeginOwnTOTPEnrollment(context.Background(), 0, 0)
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBeginOwnTOTPEnrollment_requiresOwnAccount(t *testing.T) {
	ctx := orm.ContextWithUID(context.Background(), 2)
	_, err := orm.BeginOwnTOTPEnrollment(ctx, 2, 3)
	if err == nil || !strings.Contains(err.Error(), "own account") {
		t.Fatalf("err=%v", err)
	}
}

func TestDisableUserTOTP_requiresAdmin(t *testing.T) {
	ctx := orm.ContextWithUID(context.Background(), 2)
	err := orm.DisableUserTOTP(ctx, 2, 3)
	if err == nil || !strings.Contains(err.Error(), "system administrator") {
		t.Fatalf("err=%v", err)
	}
}

func TestConfirmOwnTOTPEnrollment_requiresAuth(t *testing.T) {
	err := orm.ConfirmOwnTOTPEnrollment(context.Background(), 0, 0, "123456")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestDisableOwnTOTP_requiresAuth(t *testing.T) {
	err := orm.DisableOwnTOTP(context.Background(), 0, 0, "123456")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestUserTOTPCredentialsForLogin_withoutDB(t *testing.T) {
	secret, enabled := orm.UserTOTPCredentialsForLogin(context.Background(), 1)
	if secret != "" || enabled {
		t.Fatalf("secret=%q enabled=%v", secret, enabled)
	}
	_, ok, err := orm.OwnTOTPEnrollmentPending(context.Background(), 1, 1)
	if err == nil || ok {
		t.Fatalf("ok=%v err=%v", ok, err)
	}
	enabled, err = orm.UserTOTPEnabled(context.Background(), 1)
	if err == nil || enabled {
		t.Fatalf("enabled=%v err=%v", enabled, err)
	}
}
