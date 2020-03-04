package main

import (
	"testing"
)

//TODO: enter proper credentials here
const validHost = "http://localhost"
const validToken = "4a6244683f471d73abf96870aa3b41d025c90aba"

const invalidHost = "http://example.com"
const invalidToken = "invalidtoken"

func TestValidateHostValid(t *testing.T) {
	err := ValidateHost(validHost)
	if err != nil {
		t.Error("Given valid host, ValidateHost returned ", err)
	}
}

func TestValidateHostInvalid(t *testing.T) {
	err := ValidateHost(invalidHost)
	if err == nil {
		t.Error("Given invalid host, ValidateHost returned ", err)
	}
}

func TestValidateTokenBothValid(t *testing.T) {
	err := ValidateToken(validHost, validToken)
	if err != nil {
		t.Error("Given valid credentials, ValidateToken returned ", err)
	}
}

func TestValidateTokenInvalidHost(t *testing.T) {
	err := ValidateToken(invalidHost, validToken)
	if err == nil {
		t.Error("Given invalid host, ValidateToken returned ", err)
	}
}

func TestValidateTokenInvalidToken(t *testing.T) {
	err := ValidateToken(validHost, invalidToken)
	if err == nil {
		t.Error("Given invalid token, ValidateToken returned ", err)
	}
}

func TestValidateTokenBothInvalid(t *testing.T) {
	err := ValidateToken(invalidHost, invalidToken)
	if err == nil {
		t.Error("Given invalid host and token, ValidateToken returned ", err)
	}
}
