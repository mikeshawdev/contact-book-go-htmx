document.addEventListener("DOMContentLoaded", (event) => {
  document.body.addEventListener('htmx:beforeSwap', function (event) {
    if (event.detail.xhr.status === 422) {
      // allow 422 responses to swap as we are using this as a signal that
      // a form was submitted with bad data and want to rerender with the
      // errors
      //
      // set isError to false to avoid error logging in console
      event.detail.shouldSwap = true;
      event.detail.isError = false;
    }

    if (event.detail.xhr.status === 204) {
      // allow 204 responses to swap as we are using this as a signal that
      // a delete was successful and want to remove the element from the
      // DOM
      event.detail.shouldSwap = true;
    }
  });
})
