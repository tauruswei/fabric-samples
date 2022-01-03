var Eckles = require('eckles');
// var pem = require('fs')
//     .readFileSync('/Users/fengxiaoxiao/work/go/src/github.com/hyperledger/fabric-samples/asset-transfer-basic/application-typescript/src/mywallet/testKey_sk', 'ascii');
const pem = "-----BEGIN PRIVATE KEY-----\n" +
    "MHcCAQEEIGqxQ3Bkd7kge8FJV02vZ/NN5+99JatjIG13cVJ8bfV3oAoGCCqGSM49AwEHoUQDQgAEVlhWOXdUZkrh4ns49SNV1OjlYFomf4jNYyUvR4XFifyNGHJlfCzHqKnbhQMl7GsYyZTnAEr+QFOha1+7dBb2mg==\n" +
    "-----END PRIVATE KEY-----";

Eckles.import({ pem: pem }).then(function (jwk:string) {
    // console.log(jwk);
    Eckles.export({ jwk: jwk, format: 'pkcs8' }).then(function (pem:string) {
        // PEM in PKCS#8 format
        console.log(pem);
    });
});

Eckles.import({ pem: pem }).then(function (jwk:string) {
    // console.log(jwk);
    Eckles.export({ jwk: jwk, public: true }).then(function (pem:string) {
        // PEM in PKCS#8 format
        console.log(pem);
        console.log(Buffer.from(pem).toString('base64'));
        var str=Buffer.from('LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNURENDQWZPZ0F3SUJBZ0lSQU1xcXl1VkdJZHJkaGs3MjY3dFRYSFF3Q2dZSUtvWkl6ajBFQXdJd1ZqRUwKTUFrR0ExVUVCaE1DUTA0eEVEQU9CZ05WQkFnVEIwSmxhVXBwYm1jeEVEQU9CZ05WQkFjVEIwSmxhVXBwYm1jeApEekFOQmdOVkJBb1RCbUZoTG01d1l6RVNNQkFHQTFVRUF4TUpZMkV1WVdFdWJuQmpNQjRYRFRJeE1USXlPVEV3Ck1qUXdNRm9YRFRNeE1USXlOekV3TWpRd01Gb3dYREVMTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnVENrTmgKYkdsbWIzSnVhV0V4RmpBVUJnTlZCQWNURFZOaGJpQkdjbUZ1WTJselkyOHhEekFOQmdOVkJBc1RCbU5zYVdWdQpkREVQTUEwR0ExVUVBeE1HYzNSeWFXNW5NRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVWbGhXCk9YZFVaa3JoNG5zNDlTTlYxT2psWUZvbWY0ak5ZeVV2UjRYRmlmeU5HSEpsZkN6SHFLbmJoUU1sN0dzWXlaVG4KQUVyK1FGT2hhMSs3ZEJiMm1xT0JtekNCbURBT0JnTlZIUThCQWY4RUJBTUNCNEF3SFFZRFZSMGxCQll3RkFZSQpLd1lCQlFVSEF3SUdDQ3NHQVFVRkJ3TUJNQThHQTFVZEV3RUIvd1FGTUFNQkFmOHdLUVlEVlIwT0JDSUVJRUIwCjErZ2hwWGY5Q1lmb20vV1dtRS9FR1R6cEc3SDB2aGlXOUREOTN1M3pNQ3NHQTFVZEl3UWtNQ0tBSUNNUGxBQTgKZm5DUWc3SHZadllTNDhLT2ExSnVkaEtGb2J0UkVMZVdPZzNCTUFvR0NDcUdTTTQ5QkFNQ0EwY0FNRVFDSUJXegpRcUJlT3ZlbTQ5QUw2aEYyR0NSRWhRbzlsdTh6VHBHY0VkOVc1WE5qQWlBN0hWWnEzOERzenhCMW8zakorK3ZpCkZOeHRUK1Q3eTZ4TnUwQXNsRXIvd1E9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg==',"base64").toString();
        console.log(str);
    });
});
