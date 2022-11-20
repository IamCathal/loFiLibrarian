# loFiLibrarian

So I'm standing there at Charlie Byrnes browsing for some new books and when I pick one up most of the time I'm completely lost when it comes to important details like the genres a book consists of and if its the first in a possible series. Of course it just so happens that my data is horrendous when I'm there so looking up books on GoodReads directly takes an eternity. 

To fix this I've made this little application. I pop in the ISBN (found on the barcode of all books) and the backend then queries GoodReads and sends back a tiny JSON object with the key points that I want to know. No data is wasted on loading UI, js frameworks or images that I don't care about and I can therefore query more books faster. I've also added in some functionality to visually highlight the genres that I love and the ones that I hate

| Idle   | Detect an ISBN | Retrieve Information |  
| ----------- | ----------- | ----------- |
| ![](https://i.imgur.com/7FSbi0t.png)     | ![](https://i.imgur.com/IrPybVU.png)  |  ![](https://i.imgur.com/kPByKMB.png)  

## Installation

`docker compose up -d`

If you're looking to go about running this yourself theres some quirks. This is a weekend project so I won't be fleshing this out fully. First you'll want to edit the "good" and "bad" genres [here at the bottom of main.js](https://github.com/IamCathal/loFiLibrarian/blob/master/static/javascript/main.js#L187-L198). The video feature is only available on devices with one of course but the barcode detection itself is even more limited.

[Heres a table](https://developer.mozilla.org/en-US/docs/Web/API/Barcode_Detection_API#browser_compatibility) showing what device types and browsers support the barcode detection feature (basically its just chrome on mobile). On top of that this feature might only work when the client is communicating via HTTPS with the backend. But you can still get some use out of the project without the barcode detection, theres a manual barcode input option that still works but its a bit more tedious