const elem = document.getElementById("json")
const obj = JSON.parse(elem.textContent)
elem.textContent = JSON.stringify(obj, undefined, 4);