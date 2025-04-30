# jpmapper-osm

  _   _   _   _   _   _   _   _  
 / \ / \ / \ / \ / \ / \ / \ / \ 
( j | p | m | a | p | p | e | r )
 \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ 

📡  🏛🌲🏗🌲 <<>> 🏢🌲🏚  📡

Julian's Point Mapper

(The OpenStreetMap version)

If you don't want to use LIDAR data (more accurate, more obstructions considered), this application is faster and lighter by a lot. If accuracy of results is paramount, there's a python application which takes geotiff or las files as input and performs similar analyses using first bounce data, https://github.com/davidemerson/jpmapper-lidar

A golang cli which
 * takes address or lat/lon
 * uses the openstreetmap to calculate:
    * distance,
    * sample elevation
    * check terrain
    * check buildings and obstructions mid-path
    * report on fresnel zone obstruction.

## build
* grab jpmapper.go
* install golang
* `go build jpmapper.go`
* run the resultant binary

## freshness
The OpenStreetMap Overpass API is used to get building heights, which are as fresh as a few hours old in some cities, but not more than a year old generally.

The Open Elevation API is more stale, since that's topo maps. The ground doesn't change a lot, which is great. Won't be older than about Y2K.

The OpenStreetMap Geocoding (address lookup from lat/lon) is updated daily.

## example
```
% ./jpmapper -addr1 "144 Spencer St, Brooklyn, NY" -addr2 "303 Vernon Ave, Brooklyn, NY" -freq 5800 -debug

Fetching Elevation Profile...
Checking Points:
Checking Points: [--------------------] 0%
[DEBUG] Checking point (40.694180, -73.955217) elevation: 15.0 m
[DEBUG] Checking point (40.694215, -73.954907) elevation: 15.0 m
[DEBUG] Checking point (40.694250, -73.954597) elevation: 15.0 m
[DEBUG] Checking point (40.694285, -73.954287) elevation: 15.0 m
[DEBUG] Checking point (40.694320, -73.953977) elevation: 5.0 m
[DEBUG] Checking point (40.694355, -73.953667) elevation: 5.0 m
[DEBUG] Checking point (40.694391, -73.953357) elevation: 5.0 m
[DEBUG] Checking point (40.694426, -73.953047) elevation: 5.0 m
[DEBUG] Checking point (40.694461, -73.952737) elevation: 5.0 m
[DEBUG] Checking point (40.694496, -73.952428) elevation: 5.0 m
[DEBUG] Checking point (40.694531, -73.952118) elevation: 5.0 m
[DEBUG] Checking point (40.694566, -73.951808) elevation: 8.0 m
[DEBUG] Checking point (40.694601, -73.951498) elevation: 8.0 m
[DEBUG] Checking point (40.694636, -73.951188) elevation: 8.0 m
[DEBUG] Checking point (40.694671, -73.950878) elevation: 8.0 m
[DEBUG] Checking point (40.694707, -73.950568) elevation: 8.0 m
[DEBUG] Checking point (40.694742, -73.950258) elevation: 8.0 m
[DEBUG] Checking point (40.694777, -73.949948) elevation: 12.0 m
[DEBUG] Checking point (40.694812, -73.949638) elevation: 12.0 m
[DEBUG] Checking point (40.694847, -73.949328) elevation: 12.0 m
[DEBUG] Checking point (40.694882, -73.949018) elevation: 12.0 m
[DEBUG] Checking point (40.694917, -73.948708) elevation: 12.0 m
[DEBUG] Checking point (40.694952, -73.948398) elevation: 12.0 m
[DEBUG] Checking point (40.694987, -73.948088) elevation: 12.0 m
[DEBUG] Checking point (40.695022, -73.947778) elevation: 13.0 m
Checking Points: [##########----------] 50%68) elevation: 13.0 m
[DEBUG] Checking point (40.695093, -73.947158) elevation: 13.0 m
[DEBUG] Checking point (40.695128, -73.946848) elevation: 13.0 m
[DEBUG] Checking point (40.695163, -73.946539) elevation: 13.0 m
[DEBUG] Checking point (40.695198, -73.946229) elevation: 13.0 m
Checking Points: [############--------] 60%19) elevation: 13.0 m
[DEBUG] Checking point (40.695268, -73.945609) elevation: 15.0 m
[DEBUG] Checking point (40.695303, -73.945299) elevation: 15.0 m
[DEBUG] Checking point (40.695338, -73.944989) elevation: 15.0 m
[DEBUG] Checking point (40.695374, -73.944679) elevation: 15.0 m
Checking Points: [##############------] 70%69) elevation: 15.0 m
[DEBUG] Checking point (40.695444, -73.944059) elevation: 15.0 m
[DEBUG] Checking point (40.695479, -73.943749) elevation: 17.0 m
[DEBUG] Checking point (40.695514, -73.943439) elevation: 17.0 m
[DEBUG] Checking point (40.695549, -73.943129) elevation: 17.0 m
Checking Points: [################----] 80%19) elevation: 17.0 m
[DEBUG] Checking point (40.695619, -73.942509) elevation: 17.0 m
[DEBUG] Checking point (40.695654, -73.942199) elevation: 17.0 m
[DEBUG] Checking point (40.695690, -73.941889) elevation: 17.0 m
[DEBUG] Checking point (40.695725, -73.941579) elevation: 21.0 m
Checking Points: [##################--] 90%69) elevation: 21.0 m
[DEBUG] Checking point (40.695795, -73.940959) elevation: 21.0 m
[DEBUG] Checking point (40.695830, -73.940649) elevation: 21.0 m
[DEBUG] Checking point (40.695865, -73.940340) elevation: 23.0 m
[DEBUG] Checking point (40.695900, -73.940030) elevation: 23.0 m
[DEBUG] Checking point (40.695935, -73.939720) elevation: 23.0 m
# 📡 Julian's Point Mapper (jpmapper) Report

**Site 1**: 144 Spencer St, Brooklyn, NY
- Building Height: **23.10 m**
- Ground Elevation: **15.00 m**
- Total Antenna Height: **38.10 m**

**Site 2**: 303 Vernon Ave, Brooklyn, NY
- Building Height: **65.00 m**
- Ground Elevation: **23.00 m**
- Total Antenna Height: **88.00 m**

📏 **Surface Distance**: **1321.05 meters**

🔭 **Line of Sight (LOS)**: ✅ Clear
🌐 **Minimum Clearance**: **100.0% of Fresnel zone**

---
✅ Done.
```
