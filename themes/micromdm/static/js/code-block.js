(function () {
  var liveRegion;

  function announce(message) {
    if (!liveRegion) {
      liveRegion = document.createElement('div');
      liveRegion.className = 'visually-hidden';
      liveRegion.setAttribute('role', 'status');
      liveRegion.setAttribute('aria-live', 'polite');
      document.body.appendChild(liveRegion);
    }
    liveRegion.textContent = message;
  }

  function showCopied(button) {
    button.classList.add('code-block__copy--copied');
    announce('Code copied to clipboard');
    setTimeout(function () {
      button.classList.remove('code-block__copy--copied');
    }, 2000);
  }

  function fallbackCopy(text, button) {
    var textarea = document.createElement('textarea');
    textarea.value = text;
    textarea.setAttribute('readonly', '');
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand('copy');
      showCopied(button);
    } catch (e) {
      // Ignore copy failures.
    }
    document.body.removeChild(textarea);
  }

  function copyCode(button) {
    var block = button.closest('.code-block');
    var code = block.querySelector('code');
    var text = code ? code.textContent : '';

    if (navigator.clipboard && navigator.clipboard.writeText) {
      navigator.clipboard.writeText(text).then(function () {
        showCopied(button);
      }, function () {
        fallbackCopy(text, button);
      });
    } else {
      fallbackCopy(text, button);
    }
  }

  document.addEventListener('click', function (event) {
    var button = event.target.closest('.code-block__copy');
    if (button) {
      copyCode(button);
    }
  });
})();
