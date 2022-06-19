
let url = window.location.href
let urlBase = window.location.origin
let searchPath = "/search"
let detailPath = "/details"
let searchField = document.getElementById('search-field')
let searchButton = document.getElementById('search-btn');
let searchDropdown = document.getElementById('search-dropdown')
let searchDiv = document.getElementById('search-div')

function removeAllChildNodes(parent) {
  while (parent.firstChild) {
      parent.removeChild(parent.firstChild);
  }
}


function makeSearchResult(symbol, name) {
  let li = document.createElement('li');
  let row = document.createElement('div')
  let col3 = document.createElement('div')
  let col9 = document.createElement('div')
  li.classList.add("list-group-item", "select")
  row.classList.add("row");
  col3.classList.add("col-4")
  col3.textContent = symbol
  col9.classList.add("col-8", "text-end")
  col9.textContent = name
  row.appendChild(col3)
  row.appendChild(col9)
  li.appendChild(row)
  li.onclick = function () {
   window.location.replace(urlBase + detailPath + "/" + symbol) 
  }
  return li
}


function search() {
  fetch(urlBase + searchPath, {
    method: "post",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },

    body: JSON.stringify({term: searchField.value})
  })
  .then(response => response.json())
  .then(data => {
    for(let i = 0; i < data.length; i++) {
      searchDropdown.appendChild(makeSearchResult(data[i]["symbol"], data[i]["name"]))
    }
  });

}

searchButton.onclick = search
searchField.addEventListener('input', function() {
  removeAllChildNodes(searchDropdown)
});