# loFiLibrarian

So I'm there standing at Charlie Byrnes browsing for some new books and when I pick one up I (quickly and easily as possible) want to know its genres, if its in a series then if its the first and its general ratings on GoodReads. Of course it just so happens that my data is horrendous when I'm there so looking up books on GoodReads directly takes an eternity. 

To fix this I've made this little application. I scan the barcode (which is the ISBN13) and this backend queries GoodReads and sends back a tiny JSON object with the key points that I want to know. No data is wasted on the client side loading UI, js frameworks or images that I don't care about from GoodReads and I can therefore query more books faster. I've also added in some functionality to visually highlight the genres that I love and the ones that I hate

| Idle   | Detect an ISBN | Retrieve Information |  
| ----------- | ----------- | ----------- |
| ![](https://i.imgur.com/7FSbi0t.png)     | ![](https://i.imgur.com/IrPybVU.png)  |  ![](https://i.imgur.com/kPByKMB.png)  

## Installation

`docker compose up -d`

If you're looking to go about running this yourself there are some quirks since its purely written to serve myself. First you'll want to edit the "good" and "bad" genres [here at the bottom of main.js](https://github.com/IamCathal/loFiLibrarian/blob/master/static/javascript/main.js#L187-L198). When running locally it's also limited in scope in terms of what works. Even if you have a webcam it doesn't guarentee that the barcode scanning will work as I've linked below theres only a small subset of browsers that currently support it. Still from a locally run instance you can pop in ISBNs manually and do the same lookups (after also manually changing the references from `wss://` to `ws://` in [`main.js`](https://github.com/IamCathal/loFiLibrarian/blob/master/static/javascript/main.js)) since locally you'll not be using a HTTPS cert.

[Heres a table](https://developer.mozilla.org/en-US/docs/Web/API/Barcode_Detection_API#browser_compatibility) showing what device types and browsers support the barcode detection feature (basically its just chrome on mobile). On top of that this feature might only work when the client is communicating via HTTPS with the backend. But you can still get some use out of the project without the barcode detection, theres a manual barcode input option that still works but its a bit more tedious