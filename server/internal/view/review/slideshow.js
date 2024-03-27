;(function() {
  "use strict";

  const keyCodeToProp = {
    "ArrowLeft": "previousElementSibling",
    "ArrowRight": "nextElementSibling",
  };

  function s(selector) {
    return document.querySelector(selector);
  }

  function cl(el, o) {
    const classes = el.classList;
    Object.entries(o).forEach(function(entry) {
      classes[entry[1] ? "add" : "remove"](entry[0]);
    });
  }

  s("#play-full-screen").addEventListener("click", function() {
    s("#slideshow").requestFullscreen();
  });

  window.addEventListener("keyup", function(e) {
    if (e.defaultPrevented) {
      return;
    }

    const action = keyCodeToProp[e.key];
    if (!action) {
      return;
    }

    e.preventDefault();

    const current = s(".slide.visible"),
          next = current[action];

    if (next != null) {
      cl(current, { visible: false, hidden: true });
      cl(next, { visible: true, hidden: false });
    }
  });
})();
