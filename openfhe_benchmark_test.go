//go:build openfhe && cgo

package openfhe

import "testing"

const (
	benchDepth = 2
	benchT     = 65537
	benchMsg   = uint64(42)
)

func BenchmarkContextCreation(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ctx := NewBGVRNS(benchDepth, benchT)
		ctx.Free()
	}
}

func BenchmarkKeyGeneration(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pk, sk, err := ctx.KeyGenPtr()
		if err != nil {
			b.Fatalf("KeyGenPtr failed: %v", err)
		}
		pk.Free()
		sk.Free()
	}
}

func BenchmarkEncryption(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct, err := ctx.EncryptU64ToPtr(pk, benchMsg)
		if err != nil {
			b.Fatalf("EncryptU64ToPtr failed: %v", err)
		}
		ct.Free()
	}
}

func BenchmarkDecryption(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	ct, err := ctx.EncryptU64ToPtr(pk, benchMsg)
	if err != nil {
		b.Fatalf("EncryptU64ToPtr failed: %v", err)
	}
	defer ct.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ctx.DecryptU64FromPtr(sk, ct)
		if err != nil {
			b.Fatalf("DecryptU64FromPtr failed: %v", err)
		}
	}
}

func BenchmarkEncryptDecrypt(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct, err := ctx.EncryptU64ToPtr(pk, benchMsg)
		if err != nil {
			b.Fatalf("EncryptU64ToPtr failed: %v", err)
		}
		_, err = ctx.DecryptU64FromPtr(sk, ct)
		if err != nil {
			b.Fatalf("DecryptU64FromPtr failed: %v", err)
		}
		ct.Free()
	}
}

func BenchmarkPublicKeySerialization(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := pk.Serialize()
		if err != nil {
			b.Fatalf("PublicKey.Serialize failed: %v", err)
		}
	}
}

func BenchmarkSecretKeySerialization(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := sk.Serialize()
		if err != nil {
			b.Fatalf("SecretKey.Serialize failed: %v", err)
		}
	}
}

func BenchmarkPublicKeyDeserialization(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	pkBytes, err := pk.Serialize()
	if err != nil {
		b.Fatalf("PublicKey.Serialize failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		pk2 := DeserializePublicKey(ctx, pkBytes)
		if pk2 == nil {
			b.Fatal("DeserializePublicKey failed")
		}
		pk2.Free()
	}
}

func BenchmarkSecretKeyDeserialization(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	skBytes, err := sk.Serialize()
	if err != nil {
		b.Fatalf("SecretKey.Serialize failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sk2 := DeserializeSecretKey(ctx, skBytes)
		if sk2 == nil {
			b.Fatal("DeserializeSecretKey failed")
		}
		sk2.Free()
	}
}

func BenchmarkCiphertextSerialization(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	ct, err := ctx.EncryptU64ToPtr(pk, benchMsg)
	if err != nil {
		b.Fatalf("EncryptU64ToPtr failed: %v", err)
	}
	defer ct.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ct.Serialize()
		if err != nil {
			b.Fatalf("Ciphertext.Serialize failed: %v", err)
		}
	}
}

func BenchmarkCiphertextDeserialization(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	ct, err := ctx.EncryptU64ToPtr(pk, benchMsg)
	if err != nil {
		b.Fatalf("EncryptU64ToPtr failed: %v", err)
	}
	defer ct.Free()

	ctBytes, err := ct.Serialize()
	if err != nil {
		b.Fatalf("Ciphertext.Serialize failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ct2 := DeserializeCiphertext(ctx, ctBytes)
		if ct2 == nil {
			b.Fatal("DeserializeCiphertext failed")
		}
		ct2.Free()
	}
}

func BenchmarkHomomorphicAddition(b *testing.B) {
	ctx := NewBGVRNS(benchDepth, benchT)
	defer ctx.Free()

	pk, sk, err := ctx.KeyGenPtr()
	if err != nil {
		b.Fatalf("KeyGenPtr failed: %v", err)
	}
	defer pk.Free()
	defer sk.Free()

	ct1, err := ctx.EncryptU64ToPtr(pk, 10)
	if err != nil {
		b.Fatalf("EncryptU64ToPtr failed: %v", err)
	}
	defer ct1.Free()

	ct2, err := ctx.EncryptU64ToPtr(pk, 32)
	if err != nil {
		b.Fatalf("EncryptU64ToPtr failed: %v", err)
	}
	defer ct2.Free()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create copies for the benchmark to avoid modifying the original
		ct1Copy := DeserializeCiphertext(ctx, func() []byte {
			bytes, _ := ct1.Serialize()
			return bytes
		}())
		if ct1Copy == nil {
			b.Fatal("Failed to create ciphertext copy")
		}

		if err := ctx.EvalAdd(ct1Copy, ct2); err != nil {
			b.Fatalf("EvalAdd failed: %v", err)
		}
		ct1Copy.Free()
	}
}

func BenchmarkFullWorkflow(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create context
		ctx := NewBGVRNS(benchDepth, benchT)

		// Generate keys
		pk, sk, err := ctx.KeyGenPtr()
		if err != nil {
			b.Fatalf("KeyGenPtr failed: %v", err)
		}

		// Encrypt
		ct, err := ctx.EncryptU64ToPtr(pk, benchMsg)
		if err != nil {
			b.Fatalf("EncryptU64ToPtr failed: %v", err)
		}

		// Serialize ciphertext
		ctBytes, err := ct.Serialize()
		if err != nil {
			b.Fatalf("Ciphertext.Serialize failed: %v", err)
		}

		// Deserialize ciphertext
		ct2 := DeserializeCiphertext(ctx, ctBytes)
		if ct2 == nil {
			b.Fatal("DeserializeCiphertext failed")
		}

		// Decrypt
		_, err = ctx.DecryptU64FromPtr(sk, ct2)
		if err != nil {
			b.Fatalf("DecryptU64FromPtr failed: %v", err)
		}

		// Cleanup
		ct.Free()
		ct2.Free()
		pk.Free()
		sk.Free()
		ctx.Free()
	}
}
