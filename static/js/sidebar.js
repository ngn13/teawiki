// adds some responsivity to the sidebar

const sidebar = document.getElementById("sidebar")
const showbar = document.getElementById("showbar")

showbar.addEventListener("click", (event) => {
  if (sidebar.checkVisibility()) sidebar.style.display = "none"
  else sidebar.style.display = "flex"

  event.preventDefault()
})

window.addEventListener("resize", () => {
  if (screen.width >= 700) sidebar.style.display = "flex"
  else sidebar.style.display = "none"
})
