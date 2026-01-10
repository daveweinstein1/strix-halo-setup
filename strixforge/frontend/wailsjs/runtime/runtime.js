// Wails runtime stub
// This file provides a fallback when running outside Wails context

(function () {
    // Check if already running in Wails
    if (window.go) return;

    // Create stub for non-Wails environments
    console.log('Wails runtime not available - running in standalone mode');
})();
