
let buyLabel = document.getElementById('buy-label')
let buySlider = document.getElementById('buy-slider')

let sellLabel = document.getElementById('sell-label')
let sellSlider = document.getElementById('sell-slider')

buyLabel.innerHTML = buySlider.value

buySlider.oninput = function() {
  buyLabel.innerHTML = this.value;
}

sellLabel.innerHTML = sellSlider.value

sellSlider.oninput = function() {
  sellLabel.innerHTML = this.value;
}