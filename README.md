# jpmapper
A golang cli which
 * takes address or lat/lon
 * uses the openstreetmap to calculate:
    * distance,
    * sample elevation
    * check terrain
    * check buildings and obstructions mid-path
    * report on fresnel zone obstruction.

Source | Used for | How Recent? | Notes
OpenStreetMap (Overpass API) | Building heights | Depends on user edits; typically updated every few days to few weeks | Crowdsourced; very fresh in cities, slower in rural areas.
Open Elevation API (based on SRTM, ASTER) | Ground elevation | ~2000-2013 (depending on region) | Derived from NASA Shuttle Radar Topography Mission (SRTM) and ASTER data.
Nominatim (OpenStreetMap Geocoding) | Address → lat/lon lookup | Real-time (from latest OSM edits) | Updated daily as new OSM edits flow in.
