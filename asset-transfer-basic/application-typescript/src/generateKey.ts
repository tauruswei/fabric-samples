import * as elliptic from 'elliptic';
// const HDKey = require('@ont-dev/hdkey-secp256r1');
import {HDKey} from '@ont-dev/hdkey-secp256r1';
// const jsrsa = require('jsrsasign');
// const {KEYUTIL} = jsrsa;
// var ec = new elliptic.ec('p256');
import * as buffer from 'buffer';
import { Gateway, GatewayOptions } from 'fabric-network';
// let seed = 'a0c42a9c3ac6abf2ba6a9946ae83af18f51bf1c9fa7dacc4c92513cc4dd015834341c775dcd4c0fac73547c5662d81a9e9361a0aac604a73a321bd9103bce8af';
const bip39 = require('bip39');
const mnemonic = bip39.generateMnemonic();
console.log(mnemonic);
let seed = mnemonic;
let hdkey = HDKey.fromMasterSeed(new Buffer(seed, 'hex'));
// console.log(hdkey.toJSON());
// console.log(Buffer.from(ec.keyFromPrivate(hdkey.privateKey).getPublic(false, 'hex')).toString("hex"));
// let publicKdy = Buffer.from(ec.keyFromPrivate(hdkey.privateKey).getPublic(false, 'hex'));

// console.log(Buffer.from(ec.keyFromPrivate(hdkey.privateKey).getPublic(true, 'hex')).toString("hex"));
// let publicKdy1 = Buffer.from(ec.keyFromPrivate(hdkey.privateKey).getPublic(true, 'hex'));
// console.log(hdkey.privateKey.toString("hex"));
// // => 'xprv9s21ZrQH143K2SKJK9EYRW3Vsg8tWVHRS54hAJasj1eGsQXeWDHLeuu5hpLHRbeKedDJM4Wj9wHHMmuhPF8dQ3bzyup6R7qmMQ1i1FtzNEW'
// console.log(hdkey.publicKey.toString("hex"));

// => 'xpub661MyMwAqRbcEvPmRAmYndzERhyNux1GoHzHxgzVHMBFkCro3kbbCiDZZ5XabZDyXPj5mH3hktvkjhhUdCQxie5e1g4t2GuAWNbPmsSfDp2'
// console.log(Buffer.from(hdkey.privateKey).toString("hex"));
// console.log(Buffer.from(`30770201010420${Buffer.from(hdkey.privateKey).toString("hex")}A00A06082A8648CE3D030107A144034200${Buffer.from(hdkey.publicKey).toString("hex")}`, 'hex').toString('base64'));
console.log(`-----BEGIN PRIVATE KEY-----
${Buffer.from(`30770201010420${Buffer.from(hdkey.privateKey).toString("hex")}A00A06082A8648CE3D030107A144034200${Buffer.from(hdkey.publicKey).toString("hex")}`, 'hex').toString('base64')}
-----END PRIVATE KEY-----`);
// console.log(Buffer.from(hdkey.publicKey).toString("hex"));
// console.log(`-----BEGIN PRIVATE KEY-----
// ${Buffer.from(`308184020100301006072a8648ce3d020106052b8104000a046d306b0201010420${Buffer.from(hdkey.privateKey).toString("hex")}a144034200${Buffer.from(hdkey.publicKey).toString("hex")}`, 'hex').toString('base64')}
// -----END PRIVATE KEY-----`);

// const a = hdkey.sign(Buffer.alloc(32, 0));
// console.log(a.toString("hex"));

// const pem = KEYUTIL.getPEM(hdkey.privateKey, 'PKCS8PRV');
// console.log(pem);

// // const { Buffer } = require('buffer');
// const crypto = require('crypto');
// const keyPair = crypto.createECDH('prime256v1');
//
// // keyPair.computeSecret("123");
//
// let buffer1 = keyPair.generateKeys();
// console.log(buffer1.toString("hex"));
//
// console.log(keyPair.getPrivateKey("hex"));
// console.log(keyPair.getPublicKey("hex"));
// // Print the PEM-encoded private key
// console.log(`-----BEGIN PRIVATE KEY-----
// ${Buffer.from(`308184020100301006072a8648ce3d020106052b8104000a046d306b0201010420${keyPair.getPrivateKey('hex')}a144034200${keyPair.getPublicKey('hex')}`, 'hex').toString('base64')}
// -----END PRIVATE KEY-----`);
//
// // Print the PEM-encoded public key
// console.log(`-----BEGIN PUBLIC KEY-----
// ${Buffer.from(`3056301006072a8648ce3d020106052b8104000a034200${keyPair.getPublicKey('hex')}`, 'hex').toString('base64')}
// -----END PUBLIC KEY-----`);
//
// // Print the PEM-encoded public key
// console.log(`-----BEGIN PUBLIC KEY-----
// ${Buffer.from(`3056301006072a8648ce3d020106052b8104000a034200${hdkey.publicKey.toString('hex')}`, 'hex').toString('base64')}
// -----END PUBLIC KEY-----`);
