// commented variables are already initiated in scripts imported before
url = window.location.href;
urlBase = window.location.origin
let assetCode = url.substring(url.lastIndexOf('/') + 1)
let assetPath = '/asset'
let buyPath = '/buy'
let sellPath = 'sell'

let buyButton = document.getElementById('buy-button')
let sellButton = document.getElementById('sell-button')
let messagesDiv = document.getElementById('messages-div')
// let buySlider = document.getElementById('buy-slider')
// let sellSlider = document.getElementById('sell-slider')

function buyAsset() {
  let assetAmount = Number(buySlider.value)
  let data = {
    code: assetCode,
    amount: assetAmount
  }
  fetch(urlBase + assetPath + buyPath, {
    method: "post",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data),
  })
  .then(response => response.json())
  .then(data => {
    messagesDiv.innerText = data.messages[0]
    buySlider.max = data.max_buy
    sellSlider.max = data.max_sell
    buyLabel.innerText = buySlider.value
    sellLabel.innerText = sellSlider.value
  })
  
}


function sellAsset() {
  let assetAmount = Number(sellSlider.value)
  let data = {
    code: assetCode,
    amount: assetAmount
  }
  fetch(urlBase + assetPath + '/' + sellPath, {
    method: "post",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    credentials: 'include',
    body: JSON.stringify(data),
  })
  .then(response => response.json())
  .then(data => {
    
    messagesDiv.innerText = data.messages[0]
    buySlider.max = data.max_buy
    sellSlider.max = data.max_sell
    buyLabel.innerText = buySlider.value
    sellLabel.innerText = sellSlider.value
  })

}


buyButton.addEventListener('click', buyAsset)

sellButton.addEventListener('click', sellAsset)