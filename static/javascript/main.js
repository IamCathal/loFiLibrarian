function webcam() {
  navigator.mediaDevices.getUserMedia({
    video: {
      width: 200,
      height: 140,
      facingMode: "environment"
    }
  })
  .then(stream => {
    document.getElementById("vidOutput").srcObject = stream;
  });
}

document.getElementById("screenshotButton").addEventListener("click", (ev) => {
  const context = document.getElementById("canvas").getContext("2d");
  context.drawImage(document.getElementById("vidOutput"), 0, 0, 200, 140);
  const data = document.getElementById("canvas").toDataURL("image/png");

  const barcodeDetector = new BarcodeDetector();
  barcodeDetector.detect(data).then(data => {
    document.getElementById("detectionInfo").textContent = data;
  });
}, (err) => {
    document.getElementById("detectionInfo").textContent = err;
})

webcam()
giveSwayaaangBordersToItems()

document.getElementById("searchButton").addEventListener("click", (ev) => {
  const bookID = document.getElementById("bookIDInput").value
  document.getElementById("searchButton").classList.add("skeleton")
  document.getElementById("searchButton").style.color = "#22242f"
  lookUp(bookID).then(res => {
    fillInBookInfo(res)
    document.getElementById("searchButton").classList.remove("skeleton")
    document.getElementById("searchButton").style.color = "#c0c0c0"
  }, (err) => {
    console.error(err)
    document.getElementById("searchButton").classList.remove("skeleton")
    document.getElementById("searchButton").style.color = "#c0c0c0"
  })
})

document.getElementById("clearButton").addEventListener("click", (ev) => {
  clearBooks()
})

function fillInBookInfo(bookInfo) {
  document.getElementById("bookInfoDiv").innerHTML = 
  `
              <div class="row p-3">
                <div class="col-3 pl-1 pr-1">
                  <img 
                      src="${bookInfo.mainCover}"
                      style="width: 90%"
                  >
                </div>
                <div class="col">
                  <div class="row bookTitle">
                      ${bookInfo.title}
                  </div>
                  <div class="row bookSubInfo">
                      ${bookInfo.series}
                  </div>
                  <div class="row bookSubInfo">
                      ${bookInfo.author}
                  </div>
                  <div class="row bookPagesAndReview">
                    <div class="col-4 pl-0">
                      ${bookInfo.pages}
                    </div>
                    <div class="col pl-0">
                      ${bookInfo.rating}
                    </div>
                    <div class="col pl-0">
                      ${bookInfo.ratingsCount.toLocaleString()}
                    </div>
                  </div>
                  <div class="row bookSubInfo">
                      ${fillInGenreBlocks(bookInfo.genres)}
                  </div>
                </div>
              </div>
  
  ` + document.getElementById("bookInfoDiv").innerHTML 
}

function fillInGenreBlocks(genres) {
  let output = ""
  genres.forEach(genre => {
    output += `<div class="m-1 pl-1 pr-1 genreBox ${getGenreClass(genre)}"> ${genre} </div>`
  })
  return output
}

function clearBooks() {
  document.getElementById("bookInfoDiv").innerHTML = ""
}

function lookUp(bookID) {
  return new Promise((resolve, reject) => {
    fetch(`/lookup?id=${bookID}`, {
        method: "GET",
        headers: {
            "Content-Type": "application/json",
            "Accept": "application/json"
        },
    }).then((res) => res.json())
    .then((res) => {
        resolve(res)
    }, (err) => {
        reject(err)
    });
  })
}

function giveSwayaaangBordersToItems() {
  document.getElementById("bookInfoDiv").style += swayaaangBorders(0.8)
  document.getElementById("clearButton").style += swayaaangBorders(0.8)
  document.getElementById("searchButton").style += swayaaangBorders(0.8)
}

function swayaaangBorders(borderRadius) {
  const borderArr = [
      `border-top-right-radius: ${borderRadius}rem;`, 
      `border-bottom-right-radius: ${borderRadius}rem;`,
      `border-top-left-radius: ${borderRadius}rem;`,
      `border-bottom-left-radius: ${borderRadius}rem;`,
  ]

  let borderRadiuses = "";
  for (let k = 0; k < 4; k++) {
      const randNum = Math.floor(Math.random() * 2)
      if (randNum % 2 == 0) {
          borderRadiuses += borderArr[k]
      }
  } 
  return borderRadiuses
}

function getGenreClass(genre) {
  const goodGenres = ["Fantasy", "Epic Fantasy", "High Fantasy", "Science Fiction", "Magic",
      "Adult", "Adventure", "Fiction", "Fantasy"]
  const badGenres = ["Young Adult", "Romance", "Teen", "Family", "Sociology"]
  if (goodGenres.includes(genre)) {
    return "goodGenre"
  }
  if (badGenres.includes(genre)) {
    return "badGenre"
  }
  return "normalGenre"
}