import * as buffer from 'buffer';
import fetch from "cross-fetch";
async function test() {
    try {
        const response = await fetch("http://127.0.0.1:9445/platform/signCertificate", {
            method: 'POST',
            body: JSON.stringify({nickName: "string",publicKey:"LS0tLS1CRUdJTiBQVUJMSUMgS0VZLS0tLS0KTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFVmxoV09YZFVaa3JoNG5zNDlTTlYxT2psWUZvbQpmNGpOWXlVdlI0WEZpZnlOR0hKbGZDekhxS25iaFFNbDdHc1l5WlRuQUVyK1FGT2hhMSs3ZEJiMm1nPT0KLS0tLS1FTkQgUFVCTElDIEtFWS0tLS0t"}),
            // headers: {'Content-Type': 'application/json', 'Authorization': 'key='+API_KEY}
            headers: {'Content-Type': 'application/json'}
        });

        if (!response.ok) {
            console.error("Error");
        } else if (response.status >= 400) {
            console.error('HTTP Error: '+response.status+' - '+response.json());
        } else {
            const result = await response.json();
            console.log(Buffer.from(result.data,"base64").toString());
        }
    }catch(err){
        console.error(err);
    }
}
test();
