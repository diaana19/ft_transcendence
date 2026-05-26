# Browser Compatibility

## Supported browsers

This project is built for modern browsers and has been tested for consistent behavior across:

- Google Chrome
- Mozilla Firefox
- Apple Safari
- Microsoft Edge

## Compatibility notes

- The frontend uses React, Vite, Tailwind CSS, and modern JavaScript syntax (`async/await`, `fetch`, `URL.createObjectURL`). These are supported by current versions of Chrome, Firefox, Safari, and Edge.
- Image upload and preview use `URL.createObjectURL(file)`, which is supported in current modern browsers.
- The application uses `BrowserRouter` from `react-router-dom`, which relies on the HTML5 History API and is supported in modern browsers.

## Known limitations

- Older browsers such as Internet Explorer are not supported.
- Very old versions of Safari or Edge Legacy may not fully support the modern JavaScript features used in this app.
- Browser compatibility is intended for current stable releases of Chrome, Firefox, Safari, and Edge.

## Testing guidance

To verify browser support, open the app in each supported browser and confirm these behaviors:

- user login and registration
- post creation, including image upload
- feed rendering and media display
- profile navigation
- responsive layout and UI consistency

If any browser-specific issues are found, document them here and update the code or styles accordingly.
