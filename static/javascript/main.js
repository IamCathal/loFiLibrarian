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

document.getElementById("screenshotButton").addEventListener("click", (ev) => {
  const context = document.getElementById("canvas").getContext("2d");
  context.drawImage(document.getElementById("vidOutput"), 0, 0, 200, 140);
  const data = document.getElementById("canvas").toDataURL("image/png");
})