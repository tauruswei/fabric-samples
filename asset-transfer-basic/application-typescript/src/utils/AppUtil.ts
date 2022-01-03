/*
 * Copyright IBM Corp. All Rights Reserved.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

import { Wallet, Wallets } from 'fabric-network';
import * as fs from 'fs';
import * as path from 'path';

// const adminUserId = 'cert0';
// const adminUserPasswd = 'adminpw';

// 读取配置文件：connection-org1.json
const buildCCPOrg1 = (ccpPath: string): Record<string, any> => {
    // load the common connection configuration file
    // const ccpPath = path.resolve(__dirname, '..', 'connection-npc.json');
    // const ccpPath = path.resolve(__dirname, '..', '..', '..', '..', 'test-network',
    //     'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
    const fileExists = fs.existsSync(ccpPath);
    if (!fileExists) {
        throw new Error(`no such file or directory: ${ccpPath}`);
    }
    const contents = fs.readFileSync(ccpPath, 'utf8');

    // build a JSON object from the file contents
    const ccp = JSON.parse(contents);

    console.log(`Loaded the network configuration located at ${ccpPath}`);
    return ccp;
};

const buildCCPOrg2 = (): Record<string, any> => {
    // load the common connection configuration file
    const ccpPath = path.resolve(__dirname, '..', '..', '..', '..', 'test-network',
        'organizations', 'peerOrganizations', 'org2.example.com', 'connection-org2.json');
    const fileExists = fs.existsSync(ccpPath);
    if (!fileExists) {
        throw new Error(`no such file or directory: ${ccpPath}`);
    }
    const contents = fs.readFileSync(ccpPath, 'utf8');

    // build a JSON object from the file contents
    const ccp = JSON.parse(contents);

    console.log(`Loaded the network configuration located at ${ccpPath}`);
    return ccp;
};

// const buildWallet = async (walletPath: string): Promise<Wallet> => {
//     // Create a new  wallet : Note that wallet is for managing identities.
//     let wallet: Wallet;
//     if (walletPath) {
//         wallet = await Wallets.newFileSystemWallet(walletPath);
//         console.log(`Built a file system wallet at ${walletPath}`);
//     } else {
//         wallet = await Wallets.newInMemoryWallet();
//         console.log('Built an in memory wallet');
//     }
//
//     return wallet;
// };

const buildWallet = async (walletPath: string, mspOrgId: string, userName: string, certPemPath: string, keyPemPath: string): Promise<Wallet> => {
    // Create a new  wallet : Note that wallet is for managing identities.
    let wallet: Wallet;
    if (walletPath) {
        wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Built a file system wallet at ${walletPath}`);
    } else {
        wallet = await Wallets.newInMemoryWallet();
        console.log('Built an in memory wallet');
    }
    try {
        const certPath = path.resolve(__dirname, certPemPath);
        const cert = fs.readFileSync(certPath, 'utf8');
        console.log(cert);
        // const parse = JSON.parse(s);
        // console.log(parse);
        const keyPath = path.resolve(__dirname, keyPemPath);
        const key = fs.readFileSync(keyPath, 'utf8');
        console.log(key);
        const x509Identity = {
            credentials: {
                certificate: cert,
                privateKey: key,
            },
            mspId: mspOrgId,
            type: 'X.509',
        };
        await wallet.put(userName, x509Identity);
    } catch (error) {
        console.error(`Failed to enroll admin user : ${error}`);
    }
    return wallet;
};

const prettyJSONString = (inputString: string): string => {
    if (inputString) {
         return JSON.stringify(JSON.parse(inputString), null, 2);
    } else {
         return inputString;
    }
};

export {
    buildCCPOrg1,
    buildCCPOrg2,
    buildWallet,
    prettyJSONString,
};
