package ggps

import "math"

const earthRadius = 6378100 // in meters

// MidpointCoordinates determines coordinates of the midpoint
// points should have latitude and longitude
func MidpointCoordinates(p1, p2 []float64) []float64 {
	if len(p1) < 2 || len(p2) < 2 {
		return nil
	}
	lat1 := p1[0] * math.Pi / 180.0
	lat2 := p2[0] * math.Pi / 180.0
	lon1 := p1[1] * math.Pi / 180.0
	dLon := (p2[1] - p1[1]) * math.Pi / 180.0
	bx := math.Cos(lat2) * math.Cos(dLon)
	by := math.Cos(lat2) * math.Sin(dLon)
	lat3Rad := math.Atan2(
		math.Sin(lat1)+math.Sin(lat2),
		math.Sqrt(math.Pow(math.Cos(lat1)+bx, 2)+math.Pow(by, 2)),
	)
	lon3Rad := lon1 + math.Atan2(by, math.Cos(lat1)+bx)
	return []float64{lat3Rad * 180.0 / math.Pi, lon3Rad * 180.0 / math.Pi}
}

// DistanceBetweenPoints calculates distance between two points
// points should have latitude and longitude
func DistanceBetweenPoints(p1, p2 []float64) float64 {
	if len(p1) < 2 || len(p2) < 2 {
		return 0
	}
	p11 := p1[0]
	p21 := p2[0]
	dLat := (p21 - p11) * (math.Pi / 180.0)
	dLon := (p2[1] - p1[1]) * (math.Pi / 180.0)
	p11 = p11 * (math.Pi / 180.0)
	p21 = p21 * (math.Pi / 180.0)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(p11)*math.Cos(p21)
	return earthRadius * 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

// ShortestDistanceFromPointToLine calculates shortest distance from point to line
// and coordinates of the point this distance is measured to
// each point should have latitude and longitude
func ShortestDistanceFromPointToLine(o []float64, ab [][]float64) (float64, []float64) {
	if len(o) < 2 || len(ab) < 2 {
		return 0, nil
	}
	for _, v := range ab {
		if len(v) < 2 {
			return 0, nil
		}
	}
	abc := [][]float64{ab[0], MidpointCoordinates(ab[0], ab[1]), ab[1]}
	var (
		a, b, c      = 0, 1, 2
		epsilon      = 0.01
		distance     = 0.0
		closestPoint = abc[a]
		oa, ob, oc   float64
	)
	for i := 0; i < 30; i++ {
		oa = DistanceBetweenPoints(abc[a], o)
		ob = DistanceBetweenPoints(abc[b], o)
		oc = DistanceBetweenPoints(abc[c], o)
		distance = math.Max(oa, math.Max(ob, oc))
		switch {
		case distance == oa:
			abc[a] = MidpointCoordinates(abc[b], abc[c])
			closestPoint = abc[a]
		case distance == ob:
			abc[b] = MidpointCoordinates(abc[a], abc[c])
			closestPoint = abc[b]
		case distance == oc:
			abc[c] = MidpointCoordinates(abc[a], abc[b])
			closestPoint = abc[c]
		}
		if math.Abs(distance-oa) <= epsilon && math.Abs(distance-ob) <= epsilon && math.Abs(distance-oc) <= epsilon {
			break
		}
	}
	return distance, closestPoint
}

// ContainsLocation determines whether the point is inside the n-sized polygon.
// each point of the polygon should have latitude and longitude
func ContainsLocation(point []float64, polygon [][]float64) bool {
	size := len(polygon)
	if size < 3 || len(point) < 2 {
		return false
	}
	for _, v := range polygon {
		if len(v) != 2 {
			return false
		}
	}
	var (
		lat2, lng2, dLng3 float64
	)
	lat3 := toRadians(point[0])
	lng3 := toRadians(point[1])
	prev := polygon[size-1]
	lat1 := toRadians(prev[0])
	lng1 := toRadians(prev[1])
	nIntersect := 0
	for _, v := range polygon {
		dLng3 = wrap(lng3-lng1, -math.Pi, math.Pi)
		// special case: point equal to vertex is inside.
		if lat3 == lat1 && dLng3 == 0 {
			return true
		}
		lat2 = toRadians(v[0])
		lng2 = toRadians(v[1])
		// offset longitudes by -lng1.
		if intersects(lat1, lat2, wrap(lng2-lng1, -math.Pi, math.Pi), lat3, dLng3, true) {
			nIntersect++
		}
		lat1 = lat2
		lng1 = lng2
	}
	return (nIntersect & 1) != 0
}

func toRadians(p float64) float64 {
	return p * (math.Pi / 180.0)
}

func wrap(n, min, max float64) float64 {
	if n >= min && n < max {
		return n
	}
	return mod(n-min, max-min) + min
}

func mod(x, m float64) float64 {
	return math.Remainder(math.Remainder(x, m)+m, m)
}

func intersects(lat1, lat2, lng2, lat3, lng3 float64, geodesic bool) bool {
	// Both ends on the same side of lng3.
	if (lng3 >= 0 && lng3 >= lng2) || (lng3 < 0 && lng3 < lng2) {
		return false
	}
	// Point is South Pole.
	if lat3 <= -math.Pi/2 {
		return false
	}
	// Any segment end is a pole.
	if lat1 <= -math.Pi/2 || lat2 <= -math.Pi/2 || lat1 >= math.Pi/2 || lat2 >= math.Pi/2 {
		return false
	}
	if lng2 <= -math.Pi {
		return false
	}
	linearLat := (lat1*(lng2-lng3) + lat2*lng3) / lng2
	// northern hemisphere and point under lat-lng line.
	if lat1 >= 0 && lat2 >= 0 && lat3 < linearLat {
		return false
	}
	// southern hemisphere and point above lat-lng line.
	if lat1 <= 0 && lat2 <= 0 && lat3 >= linearLat {
		return true
	}
	// north pole.
	if lat3 >= math.Pi/2 {
		return true
	}

	// Compare lat3 with latitude on the GC/Rhumb segment corresponding to lng3.
	// Compare through a strictly-increasing function (tan() or mercator()) as convenient.
	if geodesic {
		return math.Tan(lat3) >= tanLatGC(lat1, lat2, lng2, lng3)
	}
	return mercator(lat3) >= mercatorLatRhumb(lat1, lat2, lng2, lng3)
}

func tanLatGC(lat1, lat2, lng2, lng3 float64) float64 {
	return (math.Tan(lat1)*math.Sin(lng2-lng3) + math.Tan(lat2)*math.Sin(lng3)) / math.Sin(lng2)
}

func mercator(lat float64) float64 {
	return math.Log(math.Tan(lat*0.5 + math.Pi/4))
}

func mercatorLatRhumb(lat1, lat2, lng2, lng3 float64) float64 {
	return (mercator(lat1)*(lng2-lng3) + mercator(lat2)*lng3) / lng2
}
