package main

import (
	"crypto/ed25519"
	"encoding/json"
	"testing"
)

func TestSignatureVerifiesAndForgeryFails(t *testing.T) {
	alice, err := newDevice("alice-laptop")
	if err != nil {
		t.Fatal(err)
	}
	attacker, err := newDevice("attacker-machine")
	if err != nil {
		t.Fatal(err)
	}

	message := assertionMessage(alice.DeviceID, "sync")
	valid := ed25519.Sign(alice.Private, message)
	forged := ed25519.Sign(attacker.Private, message)

	if !ed25519.Verify(alice.Public, message, valid) {
		t.Fatal("signature from the enrolled device did not verify")
	}
	if ed25519.Verify(alice.Public, message, forged) {
		t.Fatal("signature from the attacker's key verified as Alice")
	}
}

func TestCopiedStatePerfectlyCopiesClaimedIdentity(t *testing.T) {
	alice, err := newDevice("alice-laptop")
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(alice.stateFile)
	if err != nil {
		t.Fatal(err)
	}
	var copied stateFile
	if err := json.Unmarshal(encoded, &copied); err != nil {
		t.Fatal(err)
	}
	if copied != alice.stateFile {
		t.Fatalf("copied state = %#v, want %#v", copied, alice.stateFile)
	}
}

func TestLegacyBearerTokenCanBeReplayed(t *testing.T) {
	const token = "long-lived-human-copied-secret"
	server := newLegacyServer(token)

	if !server.authenticate(token) {
		t.Fatal("original token was rejected")
	}
	// There is deliberately no machine parameter: the replay is
	// indistinguishable from the original request.
	if !server.authenticate(token) {
		t.Fatal("verbatim replay was rejected")
	}
	if server.authenticate("wrong-token") {
		t.Fatal("wrong token was accepted")
	}
}
