# Grimoire Web

React frontend for the Grimoire Warhammer 40K API.

## Development

```bash
# Install dependencies
npm install

# Start development server
npm run dev

# Build for production
npm run build

# Preview production build
npm run preview
```

## Environment Variables

- `VITE_API_URL` - API base URL (defaults to `http://localhost:8080/api/v1`)

## Project Structure

```
web/
├── src/
│   ├── components/     # Reusable React components
│   ├── pages/         # Page components
│   ├── hooks/         # Custom React hooks
│   ├── utils/         # Utility functions (API client, etc.)
│   ├── types/         # TypeScript type definitions
│   ├── App.tsx        # Main app component
│   └── main.tsx       # Entry point
├── public/            # Static assets
└── dist/              # Production build output
```

## Features

- Browse units, catalogues, and factions
- Search functionality
- Unit detail pages with profiles, abilities, and weapons
- Responsive design

