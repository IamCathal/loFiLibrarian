let videoStream = null
let bookWasFoundDontScanAgainInInterval = false
let lastFoundBookID;
var barcodeDetector;

try {
  barcodeDetector = new BarcodeDetector();
} catch (error) {
  console.log("Barcode detection is not support by your browser. See https://developer.mozilla.org/en-US/docs/Web/API/Barcode_Detection_API#browser_compatibility for support details")
  hideWebcamElements()
}

webcam()
giveSwayaaangBordersToItems()
setInterval(tryToDetectISBN, 150)

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
            if (barcodes[0].rawValue == lastFoundBookID) {
              document.getElementById("detectionInfo").textContent = `This book was just looked up`
            } else {
              bookWasFoundDontScanAgainInInterval = true
              document.getElementById("detectionInfo").textContent = `Detected: ${barcodes[0].rawValue}`
              lookUpIdWs(barcodes[0].rawValue)
              lastFoundBookID = barcodes[0].rawValue

              setTimeout(() => {
                document.getElementById("detectionInfo").textContent = `Looking for ISBN...`
                bookWasFoundDontScanAgainInInterval = false
              }, 500)
            }
          
          }
        })
        .catch((err) => {
          console.error(err);
          document.getElementById("detectionInfo").textContent += err;
        });
      })
  }
}

document.getElementById("searchButton").addEventListener("click", () => {
    lookUpIdWs(document.getElementById("bookIDInput").value)
})

function lookUpIdWs(id){
    return new Promise((resolve, reject) => {
        addSearchButtonSkeleton()
        lookUpWs(id)
    })
}

document.getElementById("clearButton").addEventListener("click", (ev) => {
  clearBooks()
})

function renderPartialBookBreadcrumb(bookInfo, timeTaken, timeTakenForInitialRequest) {
  const timeTakenFormatted = timeTakenString(timeTaken)
  document.getElementById("bookInfoDiv").innerHTML = 
  `
              <div class="row pt-1 pb-1 pl-2 pr-2" id="${bookInfo.isbn}">
                  <div class="col">
                    <div class="row">
                      <div class="col-3 pl-1 pr-1">
                      <a href="${bookInfo.link}">
                        <img 
                          src="${bookInfo.mainCover}"
                          style="width: 90%"
                        >
                      </a>
                    </div>
                    <div class="col">
                      <div class="row bookTitle">
                          ${bookInfo.title}
                      </div>
                      <div class="row bookSubInfo bookSeries skeleton" id="${bookInfo.isbn}-series" style="width: 9.2rem; height: 1.2rem">
                          ${bookInfo.series}
                      </div>
                      <div class="row bookSubInfo">
                          ${bookInfo.author}
                      </div>
                      <div class="row bookPagesAndReview">
                        <div class="col-4 pl-0">
                          ${bookInfo.pages.toLocaleString()}
                        </div>
                        <div class="col pl-0">
                          ${bookInfo.rating}
                        </div>
                        <div class="col pl-0">
                          ${bookInfo.ratingsCount.toLocaleString()}
                        </div>
                      </div>
                      <div class="row bookSubInfo" id="${bookInfo.isbn}-genres">
                          <div class="m-1 pl-1 pr-1 genreBox skeleton" style="width: 9.2rem; height: 1.2rem"> </div>
                          <div class="m-1 pl-1 pr-1 genreBox skeleton" style="width: 5.1rem; height: 1.2rem"> </div>
                          <div class="m-1 pl-1 pr-1 genreBox skeleton" style="width: 2.8rem; height: 1.2rem"> </div>
                          <div class="m-1 pl-1 pr-1 genreBox skeleton" style="width: 6.5rem; height: 1.2rem"> </div>
                          <div class="m-1 pl-1 pr-1 genreBox skeleton" style="width: 3.5rem; height: 1.2rem"> </div>
                      </div>
                    </div>
                    </div>

                    <div class="row pt-1">
                        <div class="col text-center timeTakenText" id="${bookInfo.isbn}-firstRequestTimeTaken">
                          ${timeTakenString(timeTakenForInitialRequest)}
                        </div>
                        <div class="col text-center timeTakenText" id="${bookInfo.isbn}-apiLookUpTimeTaken">
                          ${timeTakenString(timeTaken)}
                        </div>
                        <div class="col text-center timeTakenText" id="${bookInfo.isbn}-pageLookupTimeTaken">

                        </div>
                    </div>
                  </div>
              </div>

              <hr class="mt-0 mb-4" style="background-color: #c0c0c0"/>
  
  ` + document.getElementById("bookInfoDiv").innerHTML 
}

function renderRemainingBookBreadcrumbDetails(bookInfo, timeTaken) {
    // fill in genres, seriesText and maybe the new profiler
    document.querySelectorAll(".genreBox").forEach(ev => {
      ev.classList.remove("skeleton")
    })
    document.querySelectorAll(".bookSeries").forEach(ev => {
      ev.classList.remove("skeleton")
      ev.style.width = "100%"
    })


    document.getElementById(`${bookInfo.isbn}-genres`).innerHTML = fillInGenreBlocks(bookInfo.genres)
    document.getElementById(`${bookInfo.isbn}-series`).innerHTML = bookInfo.series
    document.getElementById(`${bookInfo.isbn}-pageLookupTimeTaken`).innerHTML = timeTakenString(timeTaken)
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

function lookUpWs(bookId) {
    const startTime = new Date()
    const socket = new WebSocket(`wss://${getCurrentHostname()}/lookupws`);
    socket.onopen = function(ev) {
      const lookUpRequest = {
        "id": crypto.randomUUID(),
        "bookId": bookId
      }
      socket.send(JSON.stringify(lookUpRequest))
    }

    let partialBookBreadcrumbReceived = false
    let timeTakenForInitialRequest = 0


    socket.onmessage = function(ev) {
      const response = JSON.parse(ev.data)
      console.log(`Latency is ${timeSince(new Date(response.time))}`)
      console.log(response)

      switch (getMessageType(response)) {
        case "message":
            console.log("Message type")
            console.log(response.msg)
            timeTakenForInitialRequest = new Date().getTime() - startTime.getTime()
            break

        case "bookInfo":
          console.log(response.bookInfo)
            if (!partialBookBreadcrumbReceived) {
              partialBookBreadcrumbReceived = true
              renderPartialBookBreadcrumb(response.bookInfo, response.timeTaken, timeTakenForInitialRequest)
            } else {
              renderRemainingBookBreadcrumbDetails(response.bookInfo, response.timeTaken)
            }
            break
        case "error":
            document.getElementById("bookInfoDiv").innerHTML = response.errorMessage
            break

        default:
          console.error("no type")
      }
    }

    socket.onclose = function(ev) {
      console.log("socket closed!")
      removeSearchButtonSkeleton()
    }
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

function getMessageType(messageObj) {
    if (messageObj.msg != undefined) {
      return "message"
    } else if (messageObj.bookInfo != undefined) {
      return "bookInfo"
    } else if (messageObj.errormessage != undefined) {
      return "error"
    }

    console.error("Cant determine message type")
    return ""
}

function addSearchButtonSkeleton() {
  document.getElementById("searchButton").classList.add("skeleton")
  document.getElementById("searchButton").style.color = "#22242f"
}


function removeSearchButtonSkeleton() {
  document.getElementById("searchButton").classList.remove("skeleton")
  document.getElementById("searchButton").style.color = "#c0c0c0"
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

function timeSince(targetDate) {
    let diffInMs = Math.abs(targetDate.getTime() - new Date())

    if (diffInMs > 1000) {
        return Math.floor(diffInMs / 1000)+"s"
    }
    return diffInMs+"ms"
}

function timeTakenString(timeTakenMs) {
  return timeTakenMs >= 1000 ? `${timeTakenMs/1000}s` : `${timeTakenMs}ms`
}

function getCurrentHostname() {
  return new URL(window.location.href).host
}