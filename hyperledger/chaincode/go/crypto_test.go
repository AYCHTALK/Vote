package main

import (
	"crypto/ecdsa"
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
)

func createZKP(userId string, x *big.Int, v *big.Int, publicKey *ecdsa.PublicKey) []*big.Int {

	bitCurve := crypto.S256()

	Gx := bitCurve.Params().Gx
	Gy := bitCurve.Params().Gy

	if !bitCurve.IsOnCurve(publicKey.X, publicKey.Y) {
		// raise exception
		logger.Error("xG is not on curve")
		return nil
	}

	vGx, vGy := bitCurve.ScalarMult(Gx, Gy, v.Bytes())

	// Get c = H(g, g^{x}, g^{v});
	data := Append([]byte(userId), Gx, Gy, publicKey.X, publicKey.Y, vGx, vGy)
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	n := bitCurve.Params().N
	xc := mulMod(x, c, n)
	r := subMod(v, xc, n)

	return []*big.Int{r, vGx, vGy, v}
}

// from secp256k1.
func generateKeyPair() (pubkey *ecdsa.PublicKey, privkey *ecdsa.PrivateKey) {

	key, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	return &key.PublicKey, key
}

func create1outof2ZKPYesVote(
	userId string,
	xG *ecdsa.PublicKey,
	yG *ecdsa.PublicKey,
	w *big.Int,
	r1 *big.Int,
	d1 *big.Int,
	x *big.Int) ([]*big.Int, []*big.Int) {
	// Return values
	var res []*big.Int
	var res2 []*big.Int

	// Curve
	curve := crypto.S256()

	// Curve base point
	Gx := curve.Params().Gx
	Gy := curve.Params().Gy

	var temp []*big.Int
	var temp1 []*big.Int
	var temp2 []*big.Int

	// y = h^{x} * g
	temp1[0], temp1[1] = curve.ScalarMult(yG.X, yG.Y, x.Bytes())
	temp1[0], temp1[1] = curve.Add(temp1[0], temp1[1], Gx, Gy)
	res[0] = temp1[0]
	res[1] = temp1[1]

	// a1 = g^{r1} * x^{d1}
	temp1[0], temp1[1] = curve.ScalarMult(Gx, Gy, r1.Bytes())
	temp2[0], temp2[1] = curve.ScalarMult(xG.X, xG.Y, d1.Bytes())
	temp1[0], temp1[1] = curve.Add(temp1[0], temp1[1], temp2[0], temp2[1])
	res[2] = temp1[0]
	res[3] = temp1[1]

	// b1 = h^{r1} * y^{d1} (temp = affine 'y')
	temp1[0], temp1[1] = curve.ScalarMult(yG.X, yG.Y, r1.Bytes())

	// Setting temp to 'y'
	temp[0] = res[0]
	temp[1] = res[1]

	temp2[0], temp2[1] = curve.ScalarMult(temp[0], temp[1], d1.Bytes())
	temp1[0], temp1[1] = curve.Add(temp1[0], temp1[1], temp2[0], temp2[1])
	res[4] = temp1[0]
	res[5] = temp1[1]

	// a2 = g^{w}
	temp1[0], temp1[1] = curve.ScalarMult(Gx, Gy, w.Bytes())
	res[6] = temp1[0]
	res[7] = temp1[1]

	// b2 = h^{w} (where h = g^{y})
	temp1[0], temp1[1] = curve.ScalarMult(yG.X, yG.Y, w.Bytes())
	res[8] = temp1[0]
	res[9] = temp1[1]

	// Get c = H(id, xG, Y, a1, b1, a2, b2);
	// id is H(round, voter_index, voter_address, contract_address)...
	data := Append([]byte(userId)[:], xG.X, xG.Y, res[0], res[1], res[2], res[3], res[4], res[5], res[6], res[7], res[8], res[8])
	hashBytes := sha256.Sum256(data)
	c := new(big.Int)
	c.SetBytes(hashBytes[:])

	// d2 = c - d1 mod q
	n := curve.Params().N
	temp[0] = subMod(c, d1, n)
	// r2 = w - (x * d2)
	temp[1] = subMod(w, mulMod(x, temp[0], n), n)

	res2[0] = d1
	res2[1] = temp[0]
	res2[2] = r1
	res2[3] = temp[1]

	return res, res2
}

func Test_VerifyYesVote(t *testing.T) {
	// Registered Key (public keys)
	// var xG = [voter[0][0], voter[0][1]];
	// Reconstructed Key (private key?)
	// var yG = [voter[1][0], voter[1][1]];
	//       if (choice == 1) {
	//           choice_text = "YES";
	//           result = cryptoAddr.create1outof2ZKPYesVote.call(xG, yG, w, r, d, x, {
	//               from: web3.eth.accounts[accounts_index]
	//           });
	//       var y = [result[0][0], result[0][1]];
	//       var a1 = [result[0][2], result[0][3]];
	//       var b1 = [result[0][4], result[0][5]];
	//       var a2 = [result[0][6], result[0][7]];
	//       var b2 = [result[0][8], result[0][9]];
	//       var params = [result[1][0], result[1][1], result[1][2], result[1][3]];
	//       result = anonymousvotingAddr.verify1outof2ZKP.call(params, y, a1, b1, a2, b2, {
	//           from: web3.eth.accounts[accounts_index]
	//       });

	publicKeyECDSA, privateKeyECDSA := generateKeyPair()
	result1, result2 := create1outof2ZKPYesVote(userId, )

	
	y := result1[]

}

func Test_verifyZKP(t *testing.T) {
	publicKeyECDSA, privateKeyECDSA := generateKeyPair()
	userID := "someUserId"

	// We abuse this as a PRNG (pseudo-random number generator)
	_, vECDSA := generateKeyPair()
	v := vECDSA.D

	zkp := createZKP(userID, privateKeyECDSA.D, v, publicKeyECDSA)

	r := zkp[0]
	vG := []*big.Int{zkp[1], zkp[2], zkp[3]}

	scc := new(SmartContract)
	assert.True(t, scc.verifyZKP(userID, publicKeyECDSA, r, vG))
}
