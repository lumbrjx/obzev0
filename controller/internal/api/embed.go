package api

import "embed"

// uiFS holds the compiled React application served at /.
// Build with: cd ui && npm run build
// Vite outputs to controller/internal/api/ui/dist (configured in vite.config.ts).
//
//go:embed ui/dist
var uiFS embed.FS
