function delete_cookie(name) {
  document.cookie = name +'=; Path=/; Expires=Thu, 01 Jan 1970 00:00:01 GMT;';
}


let logoutBtn = document.getElementById('logout-btn')
logoutBtn.addEventListener('click', function() {
  delete_cookie("identity")
  window.location.reload()
})