
url = window.location.href;
let detailsPath = "/details"
let dataPath = "/data"
let code = '/' + url.substring(url.lastIndexOf('/') + 1)

let datapoints

function parseData(entry) {
  return {
    Date: d3.timeParse("%Y-%m-%dT%H:%M:%SZ")(entry.date),
    Sales: +entry.value
  }
  }

function get() {
  fetch(urlBase + detailsPath + code + dataPath, {
    method: "get",
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
  })
  .then(response => response.json())
  .then(data => {
    data = data.map(parseData);
      
  var svg = d3.select("svg"),
    margin = {top: 20, right: 20, bottom: 30, left: 50},
    width = +svg.attr("width") - margin.left - margin.right,
    height = +svg.attr("height") - margin.top - margin.bottom,
    g = svg.append("g").attr("transform", "translate(" + margin.left + "," + margin.top + ")");
      
  var x = d3.scaleTime()
    .rangeRound([0, width]);
      
  var y = d3.scaleLinear()
    .rangeRound([height, 0]);
      
  var line = d3.line()
    .x(function(d) { return x(d.Date); })
    .y(function(d) { return y(d.Sales); });
      
  x.domain(d3.extent(data, function(d) { return d.Date; }));
  y.domain(d3.extent(data, function(d) { return d.Sales; }));
      
  g.append("g")
    .attr("transform", "translate(0," + height + ")")
    .call(d3.axisBottom(x))
      
  g.append("g")
    .call(d3.axisLeft(y))
      
  g.append("path")
    .datum(data)
    .attr("fill", "none")
    .attr("stroke", "steelblue")
    .attr("stroke-linejoin", "round")
    .attr("stroke-linecap", "round")
    .attr("stroke-width", 1.5)
    .attr("d", line);
  })}

get()
