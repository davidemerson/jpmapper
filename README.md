# jpmapper
A golang cli which
 * takes address or lat/lon
 * uses the openstreetmap to calculate:
    * distance,
    * sample elevation
    * check terrain
    * check buildings and obstructions mid-path
    * report on fresnel zone obstruction.


## freshness
The OpenStreetMap Overpass API is used to get building heights, which are as fresh as a few hours old in some cities, but not more than a year old generally.

The Open Elevation API is more stale, since that's topo maps. The ground doesn't change a lot, which is great. Won't be older than about Y2K.

The OpenStreetMap Geocoding (address lookup from lat/lon) is updated daily.
