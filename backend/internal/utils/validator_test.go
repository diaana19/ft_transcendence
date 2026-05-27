package utils

import (
	"testing"
	"time"
)

func TestCheckUserAge_Over13(t *testing.T) {
	birthDate := time.Now().AddDate(-20, 0, 0)
	if !CheckUserAge(birthDate) {
		t.Error("20 year old should pass age check")
	}
}

func TestCheckUserAge_Exactly13(t *testing.T) {
	birthDate := time.Now().AddDate(-13, 0, 0)
	if !CheckUserAge(birthDate) {
		t.Error("exactly 13 year old should pass age check")
	}
}

func TestCheckUserAge_Under13(t *testing.T) {
	birthDate := time.Now().AddDate(-12, 0, 0)
	if CheckUserAge(birthDate) {
		t.Error("12 year old should not pass age check")
	}
}

func TestCheckUserAge_BirthdayTomorrow(t *testing.T) {
	birthDate := time.Now().AddDate(-13, 0, 1)
	if CheckUserAge(birthDate) {
		t.Error("user turning 13 tomorrow should not pass")
	}
}

func TestCheckUserAge_Newborn(t *testing.T) {
	birthDate := time.Now()
	if CheckUserAge(birthDate) {
		t.Error("newborn should not pass age check")
	}
}

func TestCheckEmailFormat_Valid(t *testing.T) {
	validEmails := []string{
		"user@example.com",
		"user.name@example.com",
		"user+tag@example.co.uk",
		"user123@test.org",
	}

	for _, email := range validEmails {
		if !CheckEmailFormat(email) {
			t.Errorf("expected valid email: %s", email)
		}
	}
}

func TestCheckEmailFormat_Invalid(t *testing.T) {
	invalidEmails := []string{
		"invalid-email",
		"test@",
		"@example.com",
		"test@.com",
		"",
		"user@example",
	}

	for _, email := range invalidEmails {
		if CheckEmailFormat(email) {
			t.Errorf("expected invalid email: %s", email)
		}
	}
}

func TestCheckPasswordFormat_Valid(t *testing.T) {
	ok, code := CheckPasswordFormat("StrongPass123!", "testuser")
	if !ok {
		t.Errorf("expected valid password, got error code %d", code)
	}
}

func TestCheckPasswordFormat_TooShort(t *testing.T) {
	ok, code := CheckPasswordFormat("Ab1!", "user")
	if ok || code != 1 {
		t.Errorf("expected error code 1 (too short), got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_NoLowercase(t *testing.T) {
	ok, code := CheckPasswordFormat("ALLUPPERS123!", "user")
	if ok || code != 2 {
		t.Errorf("expected error code 2 (no lowercase), got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_NoUppercase(t *testing.T) {
	ok, code := CheckPasswordFormat("alllowers123!", "user")
	if ok || code != 3 {
		t.Errorf("expected error code 3 (no uppercase), got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_NoDigit(t *testing.T) {
	ok, code := CheckPasswordFormat("NoDigitsHere!", "user")
	if ok || code != 4 {
		t.Errorf("expected error code 4 (no digit), got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_NoSpecial(t *testing.T) {
	ok, code := CheckPasswordFormat("NoSpecial123", "user")
	if ok || code != 5 {
		t.Errorf("expected error code 5 (no special), got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_ContainsUsername(t *testing.T) {
	ok, code := CheckPasswordFormat("myOrionPass123!", "orion123")
	if ok || code != 0 {
		t.Errorf("expected error code 0 (contains username), got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_ShortUsername_NoSubstringCheck(t *testing.T) {
	ok, _ := CheckPasswordFormat("AbcPass123!", "abc")
	if !ok {
		t.Error("username shorter than 4 chars should skip substring check")
	}
}

func TestCheckPasswordFormat_CaseInsensitiveUsername(t *testing.T) {
	ok, code := CheckPasswordFormat("myORIOPass123!", "orion")
	if ok || code != 0 {
		t.Errorf("username substring check should be case insensitive, got ok=%v code=%d", ok, code)
	}
}

func TestCheckPasswordFormat_ExactlyMinLength(t *testing.T) {
	ok, _ := CheckPasswordFormat("Abcde1!", "x")
	if ok {
		t.Error("7 char password should fail (min is 8)")
	}

	ok, _ = CheckPasswordFormat("Abcdef1!", "x")
	if !ok {
		t.Error("8 char password meeting all criteria should pass")
	}
}
