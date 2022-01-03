import { cryptico, RSAKey } from '@daotl/cryptico';
const bip39 = require('bip39');
const mnemonic = bip39.generateMnemonic();
console.log(mnemonic);

// var EC = require('elliptic').ec;
// var ec = new EC('secp256k1');
// // const SHA256 = require("crypto-js/sha256");
//
// let generateKey = ec.generateKeyPair(mnemonic);
// console.log(generateKey.toString())
// ec.getCurves();

const elliptic = require('elliptic');

const EC = elliptic.ec;
const ecdsaCurve = elliptic.curves['p384'];

const ecdsa = new EC(ecdsaCurve);
const genKeyPair = ecdsa.genKeyPair();
console.log(genKeyPair.toString());

// const cryptico = require("cryptico");

// function cryptoObj(passPhrase:string) {
//     this.bits = 1024; //2048;
//     this.passPhrase = passPhrase;
//     this.rsaKey = cryptico.generateRSAKey(this.passPhrase,this.bits);
//     this.rsaPublicKey = cryptico.publicKeyString(this.rsaKey);
//
//     this.encrypt = function(message:string){
//         var result = cryptico.encrypt(message,this.rsaPublicKey);
//         return result.cipher;
//     };
//
//     this.decrypt = function(message:string){
//         var result = cryptico.decrypt(message, this.rsaKey);
//         return result.plaintext;
//     };
// }
// class CryptoObj {
//     bits：number; //2048;
//     passPhrase: string;
//     rsaKey = cryptico.generateRSAKey(passPhrase,bits);
//     rsaPublicKey = cryptico.publicKeyString(rsaKey);
//
//     encrypt(message:string){
//         let result = cryptico.encrypt(message, rsaPublicKey);
//     };
//
//     decrypt function(message:string){
//         let result = cryptico.decrypt(message, rsaKey);
//         // return result.plaintext;/
//     };
// }
//
// console.log('---------------------------------------------------------');
// const localEncryptor = cryptoObj("XXyour secret txt or number hereXX");
//
// var encryptedMessage = localEncryptor.encrypt('new message or json code here');
// var decryptedMessage = localEncryptor.decrypt(encryptedMessage);
//
// console.log('');
// console.log('>>> Encrypted Message: '+encryptedMessage);
// console.log('');
// console.log('>>> Decrypted Message: '+decryptedMessage);

const key: RSAKey = cryptico.generateRSAKey('Made with love by DAOT Labs', 512);
// const key: RSAKey = cryptico.generate('Made with love by DAOT Labs', 512);
console.log(key.toJSON());

const crypto = require('crypto');
const RSA = 'ec';

const passphrase = 'I had learned that some things are best kept secret.';

let options = {

    // modulusLength: 256,
    // publicKeyEncoding: {
    //     type: 'spki',
    //     format: 'pem',
    // },
    // privateKeyEncoding: {
    //     type: 'pkcs8',
    //     format: 'pem',
    //     // cipher: 'aes-256-cbc',
    //     // passphrase: passphrase,
    // }
    namedCurve: 'sec256k1',

};

let start = Date.now();

let myCallback = (err: any, publicKey: any, privateKey: any) => {

    if (!err) {

        console.log('\n');
        console.log(publicKey);
        console.log(privateKey);

        let end = Date.now();
        console.log("\n> Process completed successfully in " + (end - start) + " milliseconds.");

    } else {
        throw err;

    }

};

crypto.generateKeyPair(RSA, options, myCallback);

