let videoStream = null
let bookWasFoundDontScanAgainInInterval = false
var barcodeDetector;
try {
  barcodeDetector = new BarcodeDetector();
} catch (error) {
  console.log("Barcode detection is not support by your browser. See https://developer.mozilla.org/en-US/docs/Web/API/Barcode_Detection_API#browser_compatibility for support details")
  hideWebcamElements()
}


webcam()
giveSwayaaangBordersToItems()
setInterval(tryToDetectISBN, 250)


function hideWebcamElements() {
  document.getElementById("webcamElements").style.display = "none";
  document.getElementById("manualEntryDetail").setAttribute("open", true)
}

function webcam() {
  navigator.mediaDevices.getUserMedia({
    video: {
      height: {
        min: 200,
        max: 320,
        ideal: 300
      },
      width: {
        min: 80,
        max: 160,
        ideal: 140
      },
      facingMode: "environment"
    }
  })
  .then(stream => {
    document.getElementById("vidOutput").srcObject = stream;
    videoStream = stream
  });
}

function tryToDetectISBN() {
  if (videoStream != null && bookWasFoundDontScanAgainInInterval == false) {
    let capturer = new ImageCapture(videoStream.getVideoTracks()[0])
    capturer.grabFrame().then(bitMap => {
      barcodeDetector
        .detect(bitMap)
        .then((barcodes) => {
          if (barcodes.length >= 1) {
            bookWasFoundDontScanAgainInInterval = true
            document.getElementById("detectionInfo").textContent = `Detected: ${barcodes[0].rawValue}`
            lookUpId(barcodes[0].rawValue)

            setTimeout(() => {
              document.getElementById("detectionInfo").textContent = `Looking for ISBN...`
              bookWasFoundDontScanAgainInInterval = false
            }, 2600)
          }
        })
        .catch((err) => {
          console.error(err);
          document.getElementById("detectionInfo").textContent += err;
        });
      })
  }
}

document.getElementById("searchButton").addEventListener("click", (ev) => {
  lookUpId(document.getElementById("bookIDInput").value)
})

function lookUpId(id) {
  return new Promise((resolve, reject) => {
    document.getElementById("searchButton").classList.add("skeleton")
    document.getElementById("searchButton").style.color = "#22242f"
    lookUp(id).then(res => {
      fillInBookInfo(res)
      document.getElementById("searchButton").classList.remove("skeleton")
      document.getElementById("searchButton").style.color = "#c0c0c0"
    }, (err) => {
      console.error(err)
      document.getElementById("searchButton").classList.remove("skeleton")
      document.getElementById("searchButton").style.color = "#c0c0c0"
    })
  })
}

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