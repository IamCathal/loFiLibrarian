# loFiLibrarian

So I'm there standing at Charlie Byrnes browsing for some new books and when I pick one up I (quickly and easily as possible) want to know its genres, if its in a series then if its the first in that series and its general ratings on GoodReads. Of course it just so happens that my data is horrendous when I'm there so looking up books on GoodReads directly takes an eternity. 

To fix that problem I've made this. I scan the barcode (which thankfully is the ISBN13) and this backend queries GoodReads and parses the response to then send back a tiny JSON object with the key points that I want to know. From the perspective of my own phone I might have to wait a second or two but I'm only downloading less than a kilobyte per book lookup. No data is wasted on the client side loading UI, js frameworks or images that I don't care about from GoodReads and I can therefore query more books faster since the backend is the one doing the call to GoodReads. I've also added in some functionality to visually highlight the genres that I love and the ones that I hate

| Idle   | Detect an ISBN | Retrieve Information |  
| ----------- | ----------- | ----------- |
| ![](https://i.imgur.com/7FSbi0t.png)     | ![](https://i.imgur.com/IrPybVU.png)  |  ![](https://i.imgur.com/kPByKMB.png)  

## Installation

`docker compose up -d`

If you're looking to go about running this yourself there are some quirks since its purely written to serve myself and I do not intend on making this a general use application. First you'll want to edit the "good" and "bad" genres [here at the bottom of main.js](https://github.com/IamCathal/loFiLibrarian/blob/master/static/javascript/main.js#L187-L198). When running locally it's unfortunately also limited in scope in terms of what works. Even if you have a webcam it doesn't guarentee that the barcode scanning will work as I've linked below theres only a small subset of browsers that currently support barcode detection. Still from a locally run instance you can pop in ISBNs manually and do the same lookups (after also manually changing the references from `wss://` to `ws://` in [`main.js`](https://github.com/IamCathal/loFiLibrarian/blob/master/static/javascript/main.js)) since locally you'll not be using a HTTPS cert. Unfortunately trying to connect to `wss://` first and then falling back to `ws://` just isn't possible as you might expect.

[Heres a table](https://developer.mozilla.org/en-US/docs/Web/API/Barcode_Detection_API#browser_compatibility) showing what device types and browsers support the barcode detection feature (basically its just chrome on mobile at the time of writing). On top of that this feature might only work when the client is communicating via HTTPS with the backend. But you can still get some use out of the project without the barcode detection, theres a manual barcode input option that still works but its a bit more tedious


## Configuration

`OPT_RABBITMQ_ENABLE` (defaults to false) controls whether or not the successful book lookups are sent to a RabbitMQ instance. I like to push all of my book lookups through to a RabbitMQ instance which then has other applications that read from this queue and persist them for later viewing. If you enable this you'll also then need to define valid values for 

| Env var      | Description | Example value   |
| ----------- | ----------- |----------- |
| OPT_RABBITMQ_USER      | Username of rabbitMQ account with write permissions to whatever queue used       | `my_producer`  |
| OPT_RABBITMQ_PASSWORD   | Password of rabbitMQ account        | `greenDieselHai`  |
| OPT_RABBITMQ_URL   | Full URL of the rabbitMQ instance with port        | `192.168.0.190:5672/`  |

The format of the messages it sends on this topic will change whenever I need to updat the schema but currently its in the following format:
```json
{
    "id": "dfbba440-d9cc-4d36-b685-63b904e6d9c4",
    "timestamp": 1704462430,
    "type": "lofilibrarian",
    "level": "INFO",
    "msg": "{\"title\":\"Guinness: Celebrating 250 Remarkable Years\",\"author\":\"Hartley-paul\",\"series\":\"\",\"mainCover\":\"https://images-na.ssl-images-amazon.com/images/S/compressed.photo.goodreads.com/books/1349085250i/8383928.jpg\",\"otherCovers\":[],\"pages\":128,\"link\":\"\",\"rating\":4.05,\"ratingsCount\":40,\"genres\":[\"Nonfiction\",\"Beer\",\"Ireland\"],\"isbn\":\"9780600619888\"}",
}
```

The book in its native JSON format would look like this:
```json
{
    "title": "Guinness: Celebrating 250 Remarkable Years",
    "author": "Hartley-paul",
    "series": "",
    "mainCover": "https://images-na.ssl-images-amazon.com/images/S/compressed.photo.goodreads.com/books/1349085250i/8383928.jpg",
    "otherCovers": [],
    "pages": 128,
    "link": "",
    "rating": 4.05,
    "ratingsCount": 40,
    "genres": [
        "Nonfiction",
        "Beer",
        "Ireland"
    ],
    "isbn": "9780600619888"
}
```