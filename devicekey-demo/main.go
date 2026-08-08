package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

type stateFile struct {
	DeviceID string `json:"device_id"`
	Name     string `json:"name"`
}

type device struct {
	stateFile
	Public  ed25519.PublicKey
	Private ed25519.PrivateKey
}

func newDevice(name string) (*device, error) {
	id, err := newDeviceID()
	if err != nil {
		return nil, err
	}
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generate device key: %w", err)
	}
	return &device{
		stateFile: stateFile{DeviceID: id, Name: name},
		Public:    public,
		Private:   private,
	}, nil
}

func newDeviceID() (string, error) {
	// A real ULID also encodes a timestamp. For this lesson, a random,
	// 26-character Crockford-base32 string has the same copyable-ID property.
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("generate device ID: %w", err)
	}
	encoding := base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)
	return encoding.EncodeToString(raw[:]), nil
}

func assertionMessage(deviceID, action string) []byte {
	return []byte("device-key-demo/v1\n" + deviceID + "\n" + action)
}

type legacyServer struct {
	tokenHash [sha256.Size]byte
}

func newLegacyServer(token string) *legacyServer {
	return &legacyServer{tokenHash: sha256.Sum256([]byte(token))}
}

func (s *legacyServer) authenticate(token string) bool {
	presented := sha256.Sum256([]byte(token))
	return subtle.ConstantTimeCompare(presented[:], s.tokenHash[:]) == 1
}

func step1() error {
	alice, err := newDevice("alice-laptop")
	if err != nil {
		return err
	}
	attacker, err := newDevice("attacker-machine")
	if err != nil {
		return err
	}

	stateJSON, err := json.Marshal(alice.stateFile)
	if err != nil {
		return fmt.Errorf("marshal state: %w", err)
	}
	var copied stateFile
	if err := json.Unmarshal(stateJSON, &copied); err != nil {
		return fmt.Errorf("copy state: %w", err)
	}

	message := assertionMessage(alice.DeviceID, "sync")
	aliceSignature := ed25519.Sign(alice.Private, message)
	attackerSignature := ed25519.Sign(attacker.Private, message)

	fmt.Println("STEP 1 — Assertion vs proof")
	fmt.Printf("alice state.json: %s\n", stateJSON)
	fmt.Printf("attacker copies state.json; claimed device ID matches: %t\n", copied.DeviceID == alice.DeviceID)
	fmt.Println("old server checks only the claimed ID: accepted=true")
	fmt.Printf("alice signs the claim; enrolled public key verifies it: %t\n", ed25519.Verify(alice.Public, message, aliceSignature))
	fmt.Printf("attacker signs the same claim with a different private key: verifies=%t\n", ed25519.Verify(alice.Public, message, attackerSignature))
	fmt.Println("The ID names the device; possession of its private key proves the identity.")
	return nil
}

func step2() error {
	const token = "long-lived-human-copied-secret"
	server := newLegacyServer(token)

	fmt.Println("STEP 2 — The bearer status quo")
	fmt.Printf("alice-laptop presents the long-lived token: accepted=%t\n", server.authenticate(token))
	fmt.Printf("different-machine replays the exact token: accepted=%t\n", server.authenticate(token))
	fmt.Println("Both requests are the same identity to the server: the token holder.")
	fmt.Println("It cannot revoke one machine while letting the other keep this token.")
	return nil
}

func run(step int) error {
	switch step {
	case 1:
		return step1()
	case 2:
		return step2()
	case 3, 4, 5:
		return fmt.Errorf("step %d is reserved for the next slice and is not implemented yet", step)
	default:
		return errors.New("choose a step from 1 through 5 with -step N (steps 1 and 2 are implemented)")
	}
}

func main() {
	step := flag.Int("step", 0, "lesson to run (1-5; currently 1-2)")
	flag.Parse()
	if err := run(*step); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
}
