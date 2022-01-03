// // const ethers = require('ethers');
// // let mnemonic = "derive razor possible melt gas approve coyote choose high side final choice";
// // let mnemonicWallet = ethers.Wallet.fromMnemonic(mnemonic);
// // console.log(mnemonicWallet.privateKey);
//
// // const elliptic = require('elliptic');
// // const { KEYUTIL } = require('jsrsasign');
// //
// // const privateKeyPEM = '<The PEM encoded private key>';
// // const { prvKeyHex } = KEYUTIL.getKey(privateKeyPEM); // convert the pem encoded key to hex encoded private key
// // const EdDSA = elliptic.eddsa;
// // const ec = new EdDSA('ed25519');
//
// // const EC = elliptic.eddsa;
// // // const ecdsaCurve = elliptic.curves['ed25519'];
// //
// // const ecdsa = new EC("ed25519");
// // // const genKeyPair = ecdsa.genKeyPair();
// // // console.log(genKeyPair.toString());
// // let keyPair = ecdsa.keyFromSecret("secret");
// // console.log(keyPair.toString());
// //
// // const ECDSA = require('ecdsa-secp256r1');
// //
// // const privateKey = ECDSA.generateKey();
//
// window.crypto.subtle.generateKey(
//     {
//         name: "ECDH",
//         namedCurve: "P-256", //can be "P-256", "P-384", or "P-521"
//     },
//     false, //whether the key is extractable (i.e. can be used in exportKey)
//     ["deriveKey", "deriveBits"] //can be any combination of "deriveKey" and "deriveBits"
// )
//     .then(function(key){
//         //returns a keypair object
//         console.log(key);
//         console.log(key.publicKey);
//         console.log(key.privateKey);
//     })
//     .catch(function(err){
//         console.error(err);
//     });
//
