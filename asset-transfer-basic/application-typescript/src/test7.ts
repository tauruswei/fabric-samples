//Obj of data to send in future like a dummyDb

import fetch from "node-fetch";

const data = { username: 'example' };

//POST request with body equal on data in JSON format
fetch('https://example.com/profile', {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
    },
    body: JSON.stringify(data),
})
    .then((response) => response.json())
    //Then with the data from the response in JSON...
    .then((data) => {
        console.log('Success:', data);
    })
    //Then with the error genereted...
    .catch((error) => {
        console.error('Error:', error);
    });
