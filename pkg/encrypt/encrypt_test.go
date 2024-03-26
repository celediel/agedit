package encrypt

import (
	"io/fs"
	"os"
	"testing"

	"filippo.io/age"
	"git.burning.moe/celediel/agedit/pkg/tmpfile"
)

var (
	generator        = tmpfile.NewGenerator("test_", ".txt", 18)
	strings_to_write = []string{
		"hello world",
		"hola mundo",
		"مرحبا بالعالم",
		"こんにちは世界",
		"你好世界",
		"Γειά σου Κόσμε",
		"Привіт Світ",
		"Բարեւ աշխարհ",
		"გამარჯობა მსოფლიო",
		"अभिवादन पृथ्वी",
	}
)

// TestEncryptionDecryption writes a string to a file, encrypts it, then decrypts it, and reads the string.
func TestEncryptionDecryption(t *testing.T) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatal(err)
	}

	for _, str := range strings_to_write {
		var (
			outname           string = generator.GenerateFullPath()
			encrypted_outname string = outname + ".age"
			b                 []byte
			err               error
		)

		t.Run("testing writing "+str, func(t *testing.T) {
			if err = os.WriteFile(outname, []byte(str), fs.FileMode(0600)); err != nil {
				t.Fatal(err)
			}

			if b, err = os.ReadFile(outname); err != nil {
				t.Fatal(err)
			}

			if err = Encrypt(b, encrypted_outname, id.Recipient()); err != nil {
				t.Fatal(err)
			}

			if b, err = Decrypt(encrypted_outname, id); err != nil {
				t.Fatal(err)
			}

			if string(b) != str {
				t.Fatal(string(b) + " isn't the same as " + str)
			}

			if err = os.Remove(outname); err != nil {
				t.Fatal(err)
			}

			if err = os.Remove(encrypted_outname); err != nil {
				t.Fatal(err)
			}

		})
	}
}

func TestMultipleIdentities(t *testing.T) {
	var (
		identities []age.Identity
		recipients []age.Recipient
	)

	for range 10 {
		id, err := age.GenerateX25519Identity()
		if err != nil {
			t.Fatalf("age broke: %v", err)
		}
		identities = append(identities, id)
		recipients = append(recipients, id.Recipient())
	}

	for _, str := range strings_to_write {
		var (
			outname           string = generator.GenerateFullPath()
			encrypted_outname string = outname + ".age"
			b                 []byte
			err               error
		)

		t.Run("testing writing "+str, func(t *testing.T) {
			if err = os.WriteFile(outname, []byte(str), fs.FileMode(0600)); err != nil {
				t.Fatal(err)
			}

			if b, err = os.ReadFile(outname); err != nil {
				t.Fatal(err)
			}

			if err = Encrypt(b, encrypted_outname, recipients...); err != nil {
				t.Fatal(err)
			}

			// try decrypting with each identity
			for _, id := range identities {
				if b, err = Decrypt(encrypted_outname, id); err != nil {
					t.Fatal(err)
				}
				if string(b) != str {
					t.Fatal(string(b) + " isn't the same as " + str)
				}
			}

			// then all of them because why not
			if b, err = Decrypt(encrypted_outname, identities...); err != nil {
				t.Fatal(err)
			}

			if string(b) != str {
				t.Fatal(string(b) + " isn't the same as " + str)
			}

			if err = os.Remove(outname); err != nil {
				t.Fatal(err)
			}

			if err = os.Remove(encrypted_outname); err != nil {
				t.Fatal(err)
			}

		})
	}
}

// TestNewIdentity creats a new identity, writes it to file, then re-reads it back from the file.
func TestNewIdentity(t *testing.T) {
	for range 1000 {
		outfile := generator.GenerateFullPath()

		identity, err := NewIdentity()
		if err != nil {
			t.Fatal(err)
		}

		err = WriteIdentityToFile(identity, outfile)
		if err != nil {
			t.Fatal(err)
		}

		other_identity, err := ReadIdentityFromFile(outfile)
		if err != nil {
			t.Fatal(err)
		}

		if identity.Recipient().String() != other_identity.Recipient().String() && identity.String() != other_identity.String() {
			t.Fatal("Identities don't match!", identity.Recipient(), "!=", identity.Recipient())
		}
		os.Remove(outfile)
	}
}
