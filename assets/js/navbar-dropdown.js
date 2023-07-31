let hamburgerIsOpen = false;

function openDropdown() {
  let hamburgerNavContainer = document.getElementById(
    "dropdown-nav-container"
  );

  if (!hamburgerIsOpen) {
    hamburgerNavContainer.style.display = "block";
    hamburgerIsOpen = true;
  } else {
    hamburgerNavContainer.style.display = "none";
    hamburgerIsOpen = false;
  }
}