# API Design Recommendations for BattleScribe Data

## Overview

This document outlines recommendations for building an API to serve the Warhammer 40,000 10th Edition BattleScribe XML data. The API should provide easy access to game data including units, weapons, abilities, and army building information.

---

## Architecture Options

### Option 1: Direct XML Parsing (Recommended for MVP)

**Pros:**
- Fast to implement
- No database setup required
- Always serves latest data from XML files
- Lower infrastructure complexity

**Cons:**
- Slower response times (XML parsing overhead)
- No advanced querying capabilities
- Memory intensive for large files

**Best For:** Prototype, small-scale deployment, when data changes frequently

### Option 2: Database-Backed API (Recommended for Production)

**Pros:**
- Fast query performance
- Advanced search and filtering
- Caching and indexing
- Better scalability

**Cons:**
- Requires database setup and maintenance
- Need data synchronization when XML files update
- More complex architecture

**Best For:** Production deployment, high traffic, complex queries

### Option 3: Hybrid Approach (Recommended for Scale)

**Pros:**
- Best of both worlds
- Fast queries via database
- Easy updates from XML files
- Can serve raw XML when needed

**Cons:**
- Most complex to implement
- Requires both parsing and database layers

**Best For:** Production with frequent updates, large user base

---

## Recommended Technology Stack

### Option A: Python (FastAPI) - **Recommended**

**Why:**
- Excellent XML parsing libraries (`lxml`, `xml.etree`)
- FastAPI provides automatic OpenAPI docs
- Easy to work with complex nested data
- Great async support for concurrent requests

**Stack:**
- **Framework**: FastAPI
- **XML Parsing**: `lxml` or `xml.etree.ElementTree`
- **Database** (if needed): PostgreSQL with `sqlalchemy` or MongoDB
- **Caching**: Redis (optional)
- **Deployment**: Docker + uvicorn/gunicorn

### Option B: Node.js (Express/Fastify)

**Why:**
- Good XML parsing (`xml2js`, `fast-xml-parser`)
- Fast runtime
- Easy JSON transformation
- Great for real-time features

**Stack:**
- **Framework**: Express.js or Fastify
- **XML Parsing**: `fast-xml-parser` or `xml2js`
- **Database**: PostgreSQL with Prisma or MongoDB with Mongoose
- **Caching**: Redis
- **Deployment**: Docker + PM2

### Option C: Go

**Why:**
- Excellent performance
- Built-in XML parsing
- Great concurrency
- Small memory footprint

**Stack:**
- **Framework**: Gin or Echo
- **XML Parsing**: `encoding/xml` (standard library)
- **Database**: PostgreSQL with GORM
- **Caching**: Redis
- **Deployment**: Docker

---

## API Design

### REST API Structure

#### Base URL
```
https://api.example.com/v1
```

#### Endpoints

##### 1. Game System
```
GET /game-system
GET /game-system/profile-types
GET /game-system/categories
GET /game-system/cost-types
```

##### 2. Catalogues
```
GET /catalogues                    # List all catalogues
GET /catalogues/{id}               # Get catalogue details
GET /catalogues/{id}/units         # Get all units in catalogue
GET /catalogues/{id}/detachments   # Get detachments
```

##### 3. Units
```
GET /units                         # List all units (with filters)
GET /units/{id}                    # Get unit details
GET /units/{id}/weapons            # Get unit's weapons
GET /units/{id}/abilities          # Get unit's abilities
GET /units/search?q={query}        # Search units
```

##### 4. Factions
```
GET /factions                      # List all factions
GET /factions/{faction}/units      # Get units by faction
GET /factions/{faction}/catalogues # Get catalogues for faction
```

##### 5. Army Building
```
POST /army-lists                   # Create army list
GET /army-lists/{id}               # Get army list
PUT /army-lists/{id}               # Update army list
POST /army-lists/{id}/validate     # Validate army list
GET /army-lists/{id}/points        # Calculate points
```

##### 6. Search
```
GET /search?q={query}&type={unit|weapon|ability}
GET /search/units?name={name}&faction={faction}
GET /search/weapons?name={name}&type={ranged|melee}
```

---

## Data Transformation: XML to JSON

### Unit JSON Structure

```json
{
  "id": "828d-840a-9a67-9074",
  "name": "Asurmen",
  "type": "model",
  "publication": {
    "id": "442c-a37c-3908-d701",
    "name": "Codex: Aeldari",
    "page": 7
  },
  "profiles": {
    "unit": {
      "movement": "7\"",
      "toughness": 3,
      "save": "2+",
      "wounds": 5,
      "leadership": "6+",
      "objectiveControl": 1
    },
    "abilities": [
      {
        "name": "Hand of Asuryan",
        "description": "Once per battle, when this model is selected to shoot..."
      },
      {
        "name": "Invulnerable Save (Asurmen)",
        "description": "Asurmen has a 4+ Invulnerable save"
      }
    ]
  },
  "weapons": {
    "ranged": [
      {
        "name": "The Bloody Twins",
        "range": "24\"",
        "attacks": 6,
        "ballisticSkill": "2+",
        "strength": 5,
        "armorPenetration": -1,
        "damage": 2,
        "keywords": ["Assault", "Pistol"]
      }
    ],
    "melee": [
      {
        "name": "The Sword of Asur",
        "range": "Melee",
        "attacks": 6,
        "weaponSkill": "2+",
        "strength": 6,
        "armorPenetration": -3,
        "damage": 3,
        "keywords": ["Devastating Wounds"]
      }
    ]
  },
  "categories": [
    {
      "id": "cf47-a0d7-7207-29dc",
      "name": "Infantry",
      "primary": false
    },
    {
      "id": "4f3a-f0f7-6647-348d",
      "name": "Epic Hero",
      "primary": true
    },
    {
      "id": "9cfd-1c32-585f-7d5c",
      "name": "Character",
      "primary": false
    }
  ],
  "rules": [
    {
      "id": "b4dd-3e1f-41cb-218f",
      "name": "Leader"
    },
    {
      "id": "c324-e193-e23c-7d2e",
      "name": "Battle Focus"
    }
  ],
  "costs": {
    "points": 125,
    "crusadePoints": 0
  },
  "constraints": {
    "maxPerRoster": 1
  },
  "faction": {
    "id": "4378-1827-4988-be4e",
    "name": "Faction: Asuryani"
  },
  "catalogue": {
    "id": "dfcf-1214-b57-2205",
    "name": "Aeldari - Aeldari Library"
  }
}
```

### Catalogue JSON Structure

```json
{
  "id": "34a5-8c7e-f468-82d1",
  "name": "Xenos - Aeldari",
  "revision": 8,
  "library": false,
  "gameSystemId": "sys-352e-adc2-7639-d6a9",
  "linkedCatalogues": [
    {
      "id": "dfcf-1214-b57-2205",
      "name": "Aeldari - Aeldari Library",
      "importRootEntries": true
    }
  ],
  "units": [
    {
      "id": "b425-6079-7cf8-c3fd",
      "name": "Asurmen",
      "targetId": "828d-840a-9a67-9074",
      "costs": {
        "points": 125
      }
    }
  ],
  "publications": [
    {
      "id": "442c-a37c-3908-d701",
      "name": "Codex: Aeldari",
      "shortName": "C:Aeldari",
      "publicationDate": "2025-02-08"
    }
  ]
}
```

---

## Implementation Approach

### Phase 1: XML Parser (MVP)

Create a parser that:
1. Loads XML files
2. Resolves entry links to library entries
3. Transforms XML to JSON
4. Serves via REST API

**Key Components:**
- XML Parser module
- Link Resolver (handles entryLinks → selectionEntries)
- Data Transformer (XML → JSON)
- API Routes

### Phase 2: Database Integration

1. Parse XML files on startup or via admin endpoint
2. Store transformed data in database
3. Serve from database with caching
4. Provide sync endpoint to update from XML

**Database Schema:**
- `catalogues` table
- `units` table
- `weapons` table
- `abilities` table
- `categories` table
- `rules` table
- `army_lists` table (for user-created lists)

### Phase 3: Advanced Features

- Full-text search
- Army list validation
- Point calculation engine
- Constraint checking
- Modifier application

---

## Example Implementation: Python/FastAPI

### Project Structure

```
api/
├── app/
│   ├── __init__.py
│   ├── main.py                 # FastAPI app
│   ├── config.py               # Configuration
│   ├── models/
│   │   ├── __init__.py
│   │   ├── unit.py             # Unit data models
│   │   ├── catalogue.py         # Catalogue models
│   │   └── weapon.py           # Weapon models
│   ├── parsers/
│   │   ├── __init__.py
│   │   ├── xml_parser.py        # XML parsing logic
│   │   ├── link_resolver.py     # Resolve entryLinks
│   │   └── transformer.py      # XML to JSON
│   ├── routers/
│   │   ├── __init__.py
│   │   ├── units.py             # Unit endpoints
│   │   ├── catalogues.py        # Catalogue endpoints
│   │   ├── factions.py          # Faction endpoints
│   │   └── search.py            # Search endpoints
│   └── services/
│       ├── __init__.py
│       ├── catalogue_service.py
│       └── unit_service.py
├── data/
│   └── wh40k-10e/              # XML files directory
├── requirements.txt
├── Dockerfile
└── README.md
```

### Core Parser Example

```python
# app/parsers/xml_parser.py
from lxml import etree
from typing import Dict, List, Optional
from pathlib import Path

class BattleScribeParser:
    def __init__(self, data_dir: Path):
        self.data_dir = data_dir
        self.catalogues: Dict[str, etree.Element] = {}
        self.libraries: Dict[str, etree.Element] = {}
        self.game_system: Optional[etree.Element] = None
        
    def load_game_system(self):
        """Load the game system file"""
        gst_file = self.data_dir / "Warhammer 40,000.gst"
        self.game_system = etree.parse(str(gst_file)).getroot()
        
    def load_catalogue(self, filename: str):
        """Load a catalogue file"""
        cat_file = self.data_dir / filename
        root = etree.parse(str(cat_file)).getroot()
        cat_id = root.get('id')
        
        if root.get('library') == 'true':
            self.libraries[cat_id] = root
        else:
            self.catalogues[cat_id] = root
            
    def resolve_entry_link(self, entry_link: etree.Element, catalogue_id: str) -> Optional[etree.Element]:
        """Resolve an entryLink to its actual selectionEntry"""
        target_id = entry_link.get('targetId')
        
        # Check linked libraries first
        catalogue = self.catalogues.get(catalogue_id)
        if catalogue is not None:
            # Find catalogueLinks
            for cat_link in catalogue.findall('.//{http://www.battlescribe.net/schema/catalogueSchema}catalogueLink'):
                linked_id = cat_link.get('targetId')
                library = self.libraries.get(linked_id)
                if library is not None:
                    entry = self._find_entry_by_id(library, target_id)
                    if entry is not None:
                        return entry
                        
        # Check libraries directly
        for library in self.libraries.values():
            entry = self._find_entry_by_id(library, target_id)
            if entry is not None:
                return entry
                
        return None
        
    def _find_entry_by_id(self, root: etree.Element, entry_id: str) -> Optional[etree.Element]:
        """Find a selectionEntry by ID"""
        ns = {'bs': 'http://www.battlescribe.net/schema/catalogueSchema'}
        return root.find(f'.//{{http://www.battlescribe.net/schema/catalogueSchema}}selectionEntry[@id="{entry_id}"]')
```

### API Endpoint Example

```python
# app/routers/units.py
from fastapi import APIRouter, HTTPException, Query
from typing import List, Optional
from app.models.unit import Unit, UnitSummary
from app.services.unit_service import UnitService

router = APIRouter(prefix="/units", tags=["units"])

@router.get("/", response_model=List[UnitSummary])
async def list_units(
    faction: Optional[str] = Query(None, description="Filter by faction"),
    category: Optional[str] = Query(None, description="Filter by category"),
    search: Optional[str] = Query(None, description="Search by name"),
    limit: int = Query(100, le=1000),
    offset: int = Query(0, ge=0)
):
    """List all units with optional filters"""
    service = UnitService()
    units = await service.list_units(
        faction=faction,
        category=category,
        search=search,
        limit=limit,
        offset=offset
    )
    return units

@router.get("/{unit_id}", response_model=Unit)
async def get_unit(unit_id: str):
    """Get detailed unit information"""
    service = UnitService()
    unit = await service.get_unit(unit_id)
    if unit is None:
        raise HTTPException(status_code=404, detail="Unit not found")
    return unit

@router.get("/{unit_id}/weapons")
async def get_unit_weapons(unit_id: str):
    """Get all weapons for a unit"""
    service = UnitService()
    unit = await service.get_unit(unit_id)
    if unit is None:
        raise HTTPException(status_code=404, detail="Unit not found")
    return {
        "ranged": unit.weapons.get("ranged", []),
        "melee": unit.weapons.get("melee", [])
    }
```

### Main Application

```python
# app/main.py
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.routers import units, catalogues, factions, search

app = FastAPI(
    title="Warhammer 40K 10th Edition API",
    description="API for accessing BattleScribe data",
    version="1.0.0"
)

# CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # Configure appropriately
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(units.router)
app.include_router(catalogues.router)
app.include_router(factions.router)
app.include_router(search.router)

@app.get("/")
async def root():
    return {
        "name": "Warhammer 40K 10th Edition API",
        "version": "1.0.0",
        "endpoints": {
            "units": "/units",
            "catalogues": "/catalogues",
            "factions": "/factions",
            "search": "/search"
        }
    }

@app.get("/health")
async def health():
    return {"status": "healthy"}
```

---

## Performance Considerations

### Caching Strategy

1. **In-Memory Cache**: Cache parsed XML trees and resolved entries
2. **Response Cache**: Cache JSON responses (Redis or in-memory)
3. **Database Cache**: If using database, cache frequently accessed data

### Optimization Tips

1. **Lazy Loading**: Only parse XML files when needed
2. **Indexing**: Create indexes on frequently queried fields
3. **Pagination**: Always paginate list endpoints
4. **Compression**: Enable gzip compression for responses
5. **CDN**: Use CDN for static catalogue metadata

### Expected Performance

- **Direct XML Parsing**: 50-200ms per request
- **Database Query**: 5-20ms per request
- **Cached Response**: <5ms per request

---

## Deployment Recommendations

### Development
- Run locally with `uvicorn app.main:app --reload`
- Use SQLite for database (if using database)

### Production
- **Container**: Docker with multi-stage builds
- **WSGI Server**: Gunicorn with uvicorn workers
- **Reverse Proxy**: Nginx
- **Database**: PostgreSQL or MongoDB
- **Caching**: Redis
- **Monitoring**: Prometheus + Grafana
- **Logging**: Structured logging (JSON)

### Docker Example

```dockerfile
FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application
COPY app/ ./app/
COPY data/ ./data/

# Expose port
EXPOSE 8000

# Run application
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]
```

---

## Security Considerations

1. **Rate Limiting**: Implement rate limiting on all endpoints
2. **Authentication**: Add API keys for production use
3. **Input Validation**: Validate all query parameters
4. **XML Security**: Use secure XML parsing (prevent XXE attacks)
5. **CORS**: Configure CORS appropriately
6. **HTTPS**: Always use HTTPS in production

---

## Next Steps

1. **Choose Architecture**: Direct XML parsing vs Database
2. **Choose Stack**: Python/FastAPI recommended
3. **Implement Parser**: Start with XML parser and link resolver
4. **Create API Endpoints**: Begin with units and catalogues
5. **Add Search**: Implement search functionality
6. **Add Validation**: Implement army list validation
7. **Deploy**: Set up deployment pipeline

---

## Additional Resources

- [FastAPI Documentation](https://fastapi.tiangolo.com/)
- [lxml Documentation](https://lxml.de/)
- [BattleScribe Schema](http://www.battlescribe.net/schema/)
- [REST API Best Practices](https://restfulapi.net/)

---

*This document provides a comprehensive guide for building an API to serve BattleScribe XML data. Choose the approach that best fits your needs and scale requirements.*

