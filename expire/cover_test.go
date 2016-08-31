package expire

import (
	"reflect"
	"testing"

	"github.com/omniscale/imposm3/geom/geojson"
)

// https://github.com/mapbox/tile-cover/blob/master/test/fixtures/building.geojson
var Building = geojson.Polygon{geojson.LineString{
	{-77.15269088745116, 38.87153962460514},
	{-77.1521383523941, 38.871322446566325},
	{-77.15196132659912, 38.87159391901113},
	{-77.15202569961546, 38.87162315444336},
	{-77.1519023180008, 38.87179021382536},
	{-77.15266406536102, 38.8727758561868},
	{-77.1527713537216, 38.87274662122871},
	{-77.15282499790192, 38.87282179681094},
	{-77.15323269367218, 38.87267562199469},
	{-77.15313613414764, 38.87254197618533},
	{-77.15270698070526, 38.87236656567917},
	{-77.1523904800415, 38.87198233162923},
	{-77.15269088745116, 38.87153962460514},
}}

var Donut = geojson.Polygon{geojson.LineString{
	{-76.165286, 45.479514},
	{-76.140095, 45.457437},
	{-76.162348, 45.444872},
	{-76.168656, 45.441087},
	{-76.201963, 45.420225},
	{-76.213668, 45.429276},
	{-76.214261, 45.429917},
	{-76.227477, 45.440383},
	{-76.263056, 45.467983},
	{-76.245084, 45.468609},
	{-76.240206, 45.471202},
	{-76.238518, 45.475254},
	{-76.233483, 45.507829},
	{-76.227816, 45.511836},
	{-76.212117, 45.51623},
	{-76.191776, 45.50154},
	{-76.174016, 45.486911},
	{-76.165286, 45.479514},
}, geojson.LineString{
	{-76.227618, 45.489247},
	{-76.232113, 45.486983},
	{-76.232151, 45.486379},
	{-76.231812, 45.485106},
	{-76.230698, 45.483236},
	{-76.225664, 45.477365},
	{-76.223568, 45.475174},
	{-76.202829, 45.458815},
	{-76.200229, 45.458822},
	{-76.199069, 45.459164},
	{-76.188361, 45.465784},
	{-76.204505, 45.479018},
	{-76.215555, 45.488534},
	{-76.220249, 45.492175},
	{-76.221154, 45.493315},
	{-76.22631, 45.490189},
	{-76.226543, 45.489754},
	{-76.227618, 45.489247},
}}

func comparePointTile(t *testing.T, p geojson.Point, expected Tile) {
	actual := CoverPoint(p, expected.Z)
	if actual != expected {
		t.Error("Expected ", expected, ", got ", actual)
	}
}

func TestPointToTile(t *testing.T) {
	// Knie's Kinderzoo in Rapperswil, Switzerland
	comparePointTile(t, geojson.Point{8.8223, 47.2233}, Tile{8593, 5747, 14})

}

func TestLongCoverLinestring(t *testing.T) {
	points := []geojson.Point{
		{-106.21719360351562, 28.592359801121567},
		{-106.1004638671875, 28.791130513231813},
		{-105.87661743164062, 28.864519767126602},
		{-105.82374572753905, 28.60743139267596},
	}
	tiles, _ := CoverLinestring(points, 10)

	expectedTiles := FromTiles([]Tile{
		Tile{209, 427, 10},
		Tile{209, 426, 10},
		Tile{210, 426, 10},
		Tile{210, 427, 10},
	})

	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles", tiles.ToTiles())
	}
}

// Test a spiked geojson.Polygon with many intersections
// https://github.com/mapbox/tile-cover/blob/master/test/fixtures/spiked.geojson
func TestCoverSpikePolygon(t *testing.T) {
	points := geojson.Polygon{geojson.LineString{
		{16.611328125, 8.667918002363134},
		{13.447265624999998, 3.381823735328289},
		{15.3369140625, -6.0968598188879355},
		{16.7431640625, 1.0546279422758869},
		{18.193359375, -10.314919285813147},
		{19.248046875, -1.4061088354351468},
		{20.698242187499996, -4.565473550710278},
		{22.587890625, 0.3515602939922709},
		{24.2138671875, -11.73830237143684},
		{29.091796875, 5.003394345022162},
		{26.4990234375, 9.752370139173285},
		{26.0595703125, 7.623886853120036},
		{24.9169921875, 9.44906182688142},
		{22.587890625, 6.751896464843375},
		{21.665039062499996, 12.597454504832017},
		{20.9619140625, 8.189742344383703},
		{18.193359375, 14.3069694978258},
		{16.611328125, 8.667918002363134},
	}}
	tiles := CoverPolygon(points, 6)
	expectedTiles := FromTiles([]Tile{
		Tile{35, 29, 6},
		Tile{34, 30, 6},
		Tile{35, 30, 6},
		Tile{36, 30, 6},
		Tile{37, 30, 6},
		Tile{34, 31, 6},
		Tile{35, 31, 6},
		Tile{36, 31, 6},
		Tile{37, 31, 6},
		Tile{34, 32, 6},
		Tile{35, 32, 6},
		Tile{36, 32, 6},
		Tile{34, 33, 6},
		Tile{35, 33, 6},
		Tile{36, 33, 6},
		Tile{36, 34, 6},
	})
	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles", tiles.ToTiles())
	}
}

func TestCoverLinestringLinearRing(t *testing.T) {
	outerRing := Building[0]
	tiles, ring := CoverLinestring(outerRing, 20)
	expectedTiles := FromTiles([]Tile{
		Tile{299564, 401224, 20},
		Tile{299564, 401225, 20},
		Tile{299565, 401225, 20},
		Tile{299566, 401225, 20},
		Tile{299566, 401224, 20},
		Tile{299566, 401223, 20},
		Tile{299566, 401222, 20},
		Tile{299565, 401222, 20},
		Tile{299565, 401221, 20},
		Tile{299564, 401221, 20},
		Tile{299564, 401220, 20},
		Tile{299563, 401220, 20},
		Tile{299562, 401220, 20},
		Tile{299563, 401221, 20},
		Tile{299564, 401222, 20},
		Tile{299565, 401223, 20},
		Tile{299564, 401223, 20},
	})

	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles", tiles.ToTiles())
	}
	expectedRing := []TileFraction{
		TileFraction{299564, 401224},
		TileFraction{299564, 401225},
		TileFraction{299566, 401224},
		TileFraction{299566, 401223},
		TileFraction{299566, 401222},
		TileFraction{299565, 401221},
		TileFraction{299564, 401220},
		TileFraction{299563, 401221},
		TileFraction{299564, 401222},
		TileFraction{299565, 401223},
	}
	if !reflect.DeepEqual(ring, expectedRing) {
		t.Error("Unexpected ring", tiles.ToTiles())
	}
}

// Test a building polygon
func TestCoverPolygonBuilding(t *testing.T) {
	tiles := CoverPolygon(Building, 20)
	expectedTiles := FromTiles([]Tile{
		Tile{299565, 401224, 20},
		Tile{299564, 401224, 20},
		Tile{299564, 401225, 20},
		Tile{299565, 401225, 20},
		Tile{299566, 401225, 20},
		Tile{299566, 401224, 20},
		Tile{299566, 401223, 20},
		Tile{299566, 401222, 20},
		Tile{299565, 401222, 20},
		Tile{299565, 401221, 20},
		Tile{299564, 401221, 20},
		Tile{299564, 401220, 20},
		Tile{299563, 401220, 20},
		Tile{299562, 401220, 20},
		Tile{299563, 401221, 20},
		Tile{299564, 401222, 20},
		Tile{299565, 401223, 20},
		Tile{299564, 401223, 20},
	})
	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles for building polygon", tiles.ToTiles())
	}
}

func TestCoverPolygon(t *testing.T) {
	poly := geojson.Polygon{geojson.LineString{
		{-79.9365234375, 32.77212032198862},
		{-79.9306869506836, 32.77212032198862},
		{-79.9306869506836, 32.776811185047144},
		{-79.9365234375, 32.776811185047144},
		{-79.9365234375, 32.77212032198862},
	}}
	tiles := CoverPolygon(poly, 16)

	expectedTiles := FromTiles([]Tile{
		Tile{18216, 26447, 16},
		Tile{18217, 26447, 16},
		Tile{18217, 26446, 16},
		Tile{18216, 26446, 16},
	})
	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles", tiles.ToTiles())
	}
}

func TestCoverLinestring(t *testing.T) {
	points := []geojson.Point{
		{8.826313018798828, 47.22796198584928},
		{8.82596969604492, 47.20755789924751},
		{8.826141357421875, 47.194845099780174},
	}

	tiles, _ := CoverLinestring(points, 14)
	expectedTiles := FromTiles([]Tile{
		Tile{8593, 5747, 14},
		Tile{8593, 5748, 14},
		Tile{8593, 5749, 14},
	})

	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles", tiles.ToTiles())
	}
}

func TestCoverPolygonDonut(t *testing.T) {
	tiles := CoverPolygon(Donut, 14)
	expectedTiles := FromTiles([]Tile{
		Tile{4725, 5862, 14},
		Tile{4725, 5863, 14},
		Tile{4726, 5863, 14},
		Tile{4726, 5864, 14},
		Tile{4725, 5864, 14},
		Tile{4725, 5865, 14},
		Tile{4724, 5865, 14},
		Tile{4724, 5866, 14},
		Tile{4723, 5866, 14},
		Tile{4723, 5865, 14},
		Tile{4722, 5865, 14},
		Tile{4722, 5864, 14},
		Tile{4721, 5864, 14},
		Tile{4721, 5863, 14},
		Tile{4722, 5863, 14},
		Tile{4722, 5862, 14},
		Tile{4722, 5861, 14},
		Tile{4722, 5860, 14},
		Tile{4723, 5860, 14},
		Tile{4724, 5860, 14},
		Tile{4724, 5861, 14},
		Tile{4725, 5861, 14},
		Tile{4723, 5863, 14},
		Tile{4723, 5864, 14},
		Tile{4724, 5864, 14},
		Tile{4724, 5863, 14},
		Tile{4724, 5862, 14},
		Tile{4723, 5862, 14},
		Tile{4723, 5861, 14},
	})

	if !reflect.DeepEqual(tiles, expectedTiles) {
		t.Error("Unexpected tiles", tiles.ToTiles())
	}
}
