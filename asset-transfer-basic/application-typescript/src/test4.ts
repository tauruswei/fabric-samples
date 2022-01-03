var Eckles = require('eckles');
// var pem = require('fs')
//     .readFileSync('/Users/fengxiaoxiao/work/go/src/github.com/hyperledger/fabric-samples/asset-transfer-basic/application-typescript/src/mywallet/testKey_sk', 'ascii');

//
// Eckles.import({ pem: pem }).then(function (jwk:string) {
//     // console.log(jwk);
//     Eckles.export({ jwk: jwk, format: 'pkcs8' }).then(function (pem:string) {
//         // PEM in PKCS#8 format
//         console.log(pem);
//     });
// });
// tslint:disable-next-line:no-shadowed-variable no-empty
async function ImportPrivateKey(pem: string)(Promise<string>){
    return Eckles.import({ pem });
}
async function main() {
    const pem = "-----BEGIN PRIVATE KEY-----\n" +
        "MHcCAQEEIGqxQ3Bkd7kge8FJV02vZ/NN5+99JatjIG13cVJ8bfV3oAoGCCqGSM49AwEHoUQDQgAEVlhWOXdUZkrh4ns49SNV1OjlYFomf4jNYyUvR4XFifyNGHJlfCzHqKnbhQMl7GsYyZTnAEr+QFOha1+7dBb2mg==\n" +
        "-----END PRIVATE KEY-----";
    let newVar = await ImportPrivateKey(pem);
    console.log(newVar);
}

// Eckles.import({ pem: pem }).then(function (jwk:string) {
//     // console.log(jwk);
//     Eckles.export({ jwk: jwk, public: true }).then(function (pem:string) {
//         // PEM in PKCS#8 format
//         console.log(pem);
//     });
// });

