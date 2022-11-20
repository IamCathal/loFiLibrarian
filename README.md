# loFiLibrarian

So I'm standing there at Charlie Byrnes browsing for some new books and when I pick one up most of the time I'm completely lost when it comes to important details like the genres a book consists of and if its the first in a possible series. Of course it just so happens that my data is horrendous when I'm there so looking up books on GoodReads directly takes an eternity. 

To fix this I've made this little application. I pop in the ISBN (found on the barcode of all books) and the backend then queries GoodReads and sends back a tiny JSON object with the key points that I want to know. No data is wasted on loading UI, js frameworks or images that I don't care about and I can therefore query more books faster. I've also added in some functionality to visually highlight the genres that I love and the ones that I hate

| Idle   | Detect an ISBN | Retrieve Information |  
| ----------- | ----------- | ----------- |
| ![](https://i.imgur.com/7FSbi0t.png)     | ![](https://i.imgur.com/IrPybVU.png)  |  ![](https://i.imgur.com/kPByKMB.png)  