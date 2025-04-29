package main

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io"
    "math"
    "net/http"
    "net/url"
    "os"
    "strings"
    "time"
)

var enableDebug = false

type NodeData struct {
    Lat           float64
    Lon           float64
    DisplayName   string
    BuildingHt    float64
    GroundElev    float64
    AntennaHeight float64
}

type ProfilePoint struct {
    Distance  float64
    Elevation float64
    Lat       float64
    Lon       float64
}

func main() {
    printBanner()

    addr1 := flag.String("addr1", "", "Address of first point")
    addr2 := flag.String("addr2", "", "Address of second point")
    lat1 := flag.Float64("lat1", 0, "Latitude of first point")
    lon1 := flag.Float64("lon1", 0, "Longitude of first point")
    lat2 := flag.Float64("lat2", 0, "Latitude of second point")
    lon2 := flag.Float64("lon2", 0, "Longitude of second point")
    freq := flag.Float64("freq", 2400, "Frequency in MHz (e.g., 2400 for 2.4GHz)")
    debug := flag.Bool("debug", false, "Enable debug output")
    flag.Parse()

    enableDebug = *debug

    var n1Lat, n1Lon, n2Lat, n2Lon float64
    var name1, name2 string
    var err error

    if *addr1 != "" {
        n1Lat, n1Lon, err = resolveAddress(*addr1)
        checkErr(err)
        name1 = *addr1
    } else if *lat1 != 0 && *lon1 != 0 {
        n1Lat, n1Lon = *lat1, *lon1
        name1 = fmt.Sprintf("(%.6f, %.6f)", n1Lat, n1Lon)
    } else {
        fmt.Println("Error: Must specify either addr1 or lat1/lon1")
        os.Exit(1)
    }

    if *addr2 != "" {
        n2Lat, n2Lon, err = resolveAddress(*addr2)
        checkErr(err)
        name2 = *addr2
    } else if *lat2 != 0 && *lon2 != 0 {
        n2Lat, n2Lon = *lat2, *lon2
        name2 = fmt.Sprintf("(%.6f, %.6f)", n2Lat, n2Lon)
    } else {
        fmt.Println("Error: Must specify either addr2 or lat2/lon2")
        os.Exit(1)
    }

    node1, err := enrichNodeData(n1Lat, n1Lon, name1)
    checkErr(err)
    node2, err := enrichNodeData(n2Lat, n2Lon, name2)
    checkErr(err)

    distance := calculateSurfaceDistance(node1, node2)

    fmt.Println("Fetching Elevation Profile...")
    profile, err := fetchElevationProfile(node1, node2, 50)
    checkErr(err)

    fmt.Println("Checking Points:")
    losOK, clearancePercent, blockReason := checkLOS(node1, node2, profile, *freq)

    fmt.Println("\n# 📡 Julian's Point Mapper (jpmapper) Report\n")
    fmt.Printf("**Site 1**: %s\n", node1.DisplayName)
    fmt.Printf("- Building Height: **%.2f m**\n", node1.BuildingHt)
    fmt.Printf("- Ground Elevation: **%.2f m**\n", node1.GroundElev)
    fmt.Printf("- Total Antenna Height: **%.2f m**\n\n", node1.AntennaHeight)

    fmt.Printf("**Site 2**: %s\n", node2.DisplayName)
    fmt.Printf("- Building Height: **%.2f m**\n", node2.BuildingHt)
    fmt.Printf("- Ground Elevation: **%.2f m**\n", node2.GroundElev)
    fmt.Printf("- Total Antenna Height: **%.2f m**\n\n", node2.AntennaHeight)

    fmt.Printf("📏 **Surface Distance**: **%.2f meters**\n", distance)

    if losOK {
        fmt.Printf("\n🔭 **Line of Sight (LOS)**: ✅ Clear\n")
    } else {
        fmt.Printf("\n🔭 **Line of Sight (LOS)**: ❌ Blocked (%s)\n", blockReason)
    }
    fmt.Printf("🌐 **Minimum Clearance**: **%.1f%% of Fresnel zone**\n", clearancePercent)

    fmt.Println("\n---\n✅ Done.\n")
}

func printBanner() {
    banner := `
  _   _   _   _   _   _   _   _  
 / \ / \ / \ / \ / \ / \ / \ / \ 
( j | p | m | a | p | p | e | r )
 \_/ \_/ \_/ \_/ \_/ \_/ \_/ \_/ 

Julian's Point Mapper (jpmapper)

Usage Examples:
  ./jpmapper -addr1 "Empire State Building, NYC" -addr2 "One World Trade Center, NYC" -freq 5800
  ./jpmapper -lat1 40.7484 -lon1 -73.9857 -lat2 40.7127 -lon2 -74.0134 -freq 2400 -debug
  ./jpmapper -addr1 "144 Spencer St, Brooklyn, NY" -addr2 "303 Vernon Ave, Brooklyn, NY" -freq 5800

Notes:
        * addresses and place names will be resolved to lat/lon
        * freq is in MHz
        * debug shows obstruction work so you can troubleshoot buildings
---
`
    fmt.Println(banner)
}

func resolveAddress(address string) (float64, float64, error) {
    query := url.QueryEscape(address)
    apiURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", query)

    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(apiURL)
    if err != nil {
        return 0, 0, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return 0, 0, err
    }

    var results []struct {
        Lat string `json:"lat"`
        Lon string `json:"lon"`
    }
    if err := json.Unmarshal(body, &results); err != nil {
        return 0, 0, err
    }

    if len(results) == 0 {
        return 0, 0, fmt.Errorf("address not found: %s", address)
    }

    var lat, lon float64
    fmt.Sscanf(results[0].Lat, "%f", &lat)
    fmt.Sscanf(results[0].Lon, "%f", &lon)
    return lat, lon, nil
}

func enrichNodeData(lat, lon float64, name string) (NodeData, error) {
    elevation, err := fetchElevation(lat, lon)
    if err != nil {
        return NodeData{}, fmt.Errorf("elevation fetch error: %w", err)
    }
    bHeight, err := fetchBuildingHeight(lat, lon)
    if err != nil {
        return NodeData{}, fmt.Errorf("building height fetch error: %w", err)
    }
    return NodeData{
        Lat: lat,
        Lon: lon,
        DisplayName: name,
        GroundElev: elevation,
        BuildingHt: bHeight,
        AntennaHeight: elevation + bHeight,
    }, nil
}

func fetchBuildingHeight(lat, lon float64) (float64, error) {
    overpassQuery := fmt.Sprintf(`[out:json];way(around:10,%.6f,%.6f)["building"]["height"];out center;`, lat, lon)
    body := bytes.NewBufferString("data=" + overpassQuery)
    resp, err := doPost("https://overpass-api.de/api/interpreter", "application/x-www-form-urlencoded", body)
    if err != nil {
        return 0, err
    }

    var result struct {
        Elements []struct {
            Tags map[string]string `json:"tags"`
        } `json:"elements"`
    }
    if err := json.Unmarshal(resp, &result); err != nil {
        return 0, err
    }

    if len(result.Elements) == 0 {
        return 30.0, nil
    }

    hStr := result.Elements[0].Tags["height"]
    var height float64
    fmt.Sscanf(hStr, "%f", &height)
    if height == 0 {
        return 30.0, nil
    }
    return height, nil
}

func fetchElevation(lat, lon float64) (float64, error) {
    url := fmt.Sprintf("https://api.open-elevation.com/api/v1/lookup?locations=%.6f,%.6f", lat, lon)
    body, err := doGet(url)
    if err != nil {
        return 0, err
    }

    var result struct {
        Results []struct {
            Elevation float64 `json:"elevation"`
        } `json:"results"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return 0, err
    }

    if len(result.Results) == 0 {
        return 0, fmt.Errorf("no elevation data found")
    }

    return result.Results[0].Elevation, nil
}

func fetchElevationProfile(n1, n2 NodeData, samples int) ([]ProfilePoint, error) {
    var points []string
    for i := 0; i <= samples; i++ {
        frac := float64(i) / float64(samples)
        lat := n1.Lat + frac*(n2.Lat-n1.Lat)
        lon := n1.Lon + frac*(n2.Lon-n1.Lon)
        points = append(points, fmt.Sprintf("%.6f,%.6f", lat, lon))
    }
    url := fmt.Sprintf("https://api.open-elevation.com/api/v1/lookup?locations=%s", strings.Join(points, "|"))

    body, err := doGet(url)
    if err != nil {
        return nil, err
    }

    var result struct {
        Results []struct {
            Elevation float64 `json:"elevation"`
        } `json:"results"`
    }
    if err := json.Unmarshal(body, &result); err != nil {
        return nil, err
    }

    distTotal := calculateSurfaceDistance(n1, n2)
    var profile []ProfilePoint
    for i, r := range result.Results {
        frac := float64(i) / float64(samples)
        lat := n1.Lat + frac*(n2.Lat-n1.Lat)
        lon := n1.Lon + frac*(n2.Lon-n1.Lon)
        profile = append(profile, ProfilePoint{
            Distance:  distTotal * frac,
            Elevation: r.Elevation,
            Lat:       lat,
            Lon:       lon,
        })
    }
    return profile, nil
}

func calculateSurfaceDistance(n1, n2 NodeData) float64 {
    const R = 6371000
    phi1 := n1.Lat * math.Pi / 180
    phi2 := n2.Lat * math.Pi / 180
    dPhi := (n2.Lat - n1.Lat) * math.Pi / 180
    dLambda := (n2.Lon - n1.Lon) * math.Pi / 180

    a := math.Sin(dPhi/2)*math.Sin(dPhi/2) +
        math.Cos(phi1)*math.Cos(phi2)*math.Sin(dLambda/2)*math.Sin(dLambda/2)
    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}

func checkLOS(n1, n2 NodeData, profile []ProfilePoint, freqMHz float64) (bool, float64, string) {
    surfaceDist := calculateSurfaceDistance(n1, n2)
    minClearance := 100.0
    blockReason := ""
    var lastPercent int = -1

    for idx, p := range profile {
        frac := p.Distance / surfaceDist
        antennaLine := n1.AntennaHeight + frac*(n2.AntennaHeight-n1.AntennaHeight)
        fresnel := calculateFresnelRadius(surfaceDist, freqMHz)
        clearance := (antennaLine - p.Elevation) / fresnel * 100

        if clearance < minClearance {
            minClearance = clearance
            if clearance < 60 {
                blockReason = "terrain obstruction"
            }
        }

        // Progress bar updating
        percent := (idx * 100) / len(profile)
        if percent%10 == 0 && percent != lastPercent {
            fmt.Printf("\rChecking Points: [%s%s] %d%%", strings.Repeat("#", percent/5), strings.Repeat("-", 20-percent/5), percent)
            lastPercent = percent
        }

        if enableDebug {
            fmt.Printf("\n[DEBUG] Checking point (%.6f, %.6f) elevation: %.1f m", p.Lat, p.Lon, p.Elevation)
        }

        time.Sleep(150 * time.Millisecond)
    }

    return minClearance >= 60, minClearance, blockReason
}

func calculateFresnelRadius(d, freqMHz float64) float64 {
    return 17.32 * math.Sqrt((d/1000)/(4*freqMHz))
}

func checkErr(err error) {
    if err != nil {
        fmt.Println("Fatal error:", err)
        os.Exit(1)
    }
}

func doGet(url string) ([]byte, error) {
    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}

func doPost(url, contentType string, body *bytes.Buffer) ([]byte, error) {
    client := http.Client{Timeout: 10 * time.Second}
    resp, err := client.Post(url, contentType, body)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    return io.ReadAll(resp.Body)
}


