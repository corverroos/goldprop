package main

// generate this with help of `goldprop/cmd/areagen`
var areas = map[string]int{
	"Adriaanse":                  8128,
	"Airlie":                     10184,
	"Airport Industria":          11637,
	"Alphen":                     10056,
	"Athlone":                    8734,
	"Athlone Industrial":         8735,
	"Bakoven":                    11012,
	"Bantry Bay":                 11013,
	"Barbarosa Estate":           10058,
	"Bel Ombre":                  10059,
	"Belgravia":                  8736,
	"Belhar":                     8132,
	"Belle Constantia":           16246,
	"Belthorne Estate":           8783,
	"Bergvliet":                  10189,
	"Bishop Lavis":               8121,
	"Bishopscourt":               8661,
	"Bishopscourt Village":       14241,
	"Black River":                8662,
	"Bo Kaap":                    9136,
	"Bonteheuwel":                8739,
	"Bridgetown":                 8723,
	"Browns Farm":                16526,
	"Buckingham":                 16443,
	"Cafda Village":              10191,
	"Camps Bay":                  11014,
	"Cape Farms":                 12693,
	"Cape Town":                  432,
	"Cape Town City Centre":      9138,
	"Capricorn":                  12788,
	"Charlesville":               8155,
	"Claremont":                  11741,
	"Claremont Upper":            14225,
	"Clarkes Estate":             8133,
	"Clifton":                    11015,
	"Coniston Park":              9066,
	"Constantia":                 11742,
	"Constantia Heights":         10112,
	"Constantia Hill Estate":     10065,
	"Constantia Vale":            16257,
	"Constantia Village":         16413,
	"Costa Da Gama":              13856,
	"Crawford":                   8787,
	"Crossroads":                 8741,
	"De Waterkant":               9141,
	"Dennedal":                   11743,
	"Devil's Peak Estate":        16541,
	"Diep River":                 10195,
	"Dreyersdal Estate":          9049,
	"Eagles Nest":                10111,
	"Eden Klein":                 16444,
	"Edward":                     16452,
	"Elfindale":                  10170,
	"Epping Industrial":          8098,
	"Erica Township":             16397,
	"Factreton":                  8006,
	"Fairways":                   10068,
	"Ferness Estate":             10069,
	"Fisantekraal":               9529,
	"Fisantekraal Industrial":    16421,
	"Foreshore":                  9143,
	"Frere Estate":               8788,
	"Fresnaye":                   11016,
	"Frogmore Estate":            10203,
	"Gardens":                    9145,
	"Gatesville":                 8746,
	"Gaylands":                   10071,
	"Gleemoor":                   8747,
	"Grassy Park":                10204,
	"Green Point":                11017,
	"Greenville Garden City":     16424,
	"Guguletu":                   8696,
	"Hanover Park":               8790,
	"Harfield Village":           8666,
	"Hatton":                     16261,
	"Hazendal":                   8717,
	"Heathfield":                 10207,
	"Heideveld":                  8726,
	"Heinz Park":                 16284,
	"Higgovale":                  9149,
	"High Constantia":            16248,
	"Kalk Bay":                   9067,
	"Kalksteenfontein":           8753,
	"Kenilworth":                 8669,
	"Kenilworth Upper":           14224,
	"Kensington":                 8010,
	"Kenwyn":                     8792,
	"Kewtown":                    8754,
	"Kirstenhof":                 10178,
	"Knole Park":                 16285,
	"Kreupelbosch":               11757,
	"Lakeside":                   10174,
	"Langa":                      9117,
	"Lansdowne":                  8778,
	"Lavender Hill":              10211,
	"Llandudno":                  9118,
	"Lotus River":                10212,
	"Maitland":                   9153,
	"Maitland Garden Village":    9119,
	"Manenberg":                  8718,
	"Marina Da Gama":             10213,
	"Matroosfontein":             8154,
	"Meadowridge":                10052,
	"Montana":                    8165,
	"Montevideo":                 8157,
	"Morning Star":               7906,
	"Mouille Point":              11018,
	"Mountview":                  8716,
	"Mowbray":                    8677,
	"Muizenberg":                 9025,
	"Ndabeni":                    8014,
	"Nerissa Estate":             8798,
	"Newfields":                  8727,
	"Newlands":                   8679,
	"Newlands Upper":             14240,
	"Nieuwe Steenberg":           16008,
	"Nooitgedacht":               8123,
	"Nova Constantia":            16249,
	"Nyanga":                     9132,
	"Observatory":                10157,
	"Ocean View":                 9069,
	"Oranjezicht":                9155,
	"Ottery":                     10090,
	"Ottery East":                16453,
	"Otyhouse Estate":            16454,
	"Parktown":                   8756,
	"Parkwood":                   10179,
	"Pelican Heights":            12195,
	"Pelican Park":               9029,
	"Penlyn Estate":              8800,
	"Philadelphia":               9168,
	"Philippi":                   8818,
	"Philippi East":              16287,
	"Pinati":                     8779,
	"Pinelands":                  8017,
	"Plumstead":                  10094,
	"Primrose Park":              8758,
	"Punts Estate":               9032,
	"Retreat":                    9034,
	"Retreat Industrial":         25704,
	"Romp Vallei":                16445,
	"Rondebosch":                 8682,
	"Rondebosch East":            8806,
	"Rondebosch Park Estate":     32960,
	"Rondebosch Village":         8660,
	"Rosebank":                   8683,
	"Rylands":                    8730,
	"Salt River":                 10158,
	"Sand Industria":             16183,
	"Schaap Kraal":               12527,
	"Sea Point":                  11021,
	"Seawinds":                   10180,
	"Sheraton Park":              15769,
	"Silvertown":                 8731,
	"Silvertree Estate":          15132,
	"Silwersteen Estate":         16720,
	"Southfield":                 10096,
	"Springfield":                16288,
	"St James":                   9039,
	"Steenberg":                  9040,
	"Steenberg Golf Estate":      15183,
	"Stonehurst Mountain Estate": 12840,
	"Sunlands":                   16499,
	"Sunnyside":                  8761,
	"Surrey Estate":              8762,
	"Sweet Home":                 16289,
	"Sybrand Park":               8686,
	"Tamboerskloof":              9163,
	"The Vines Estate":           10098,
	"Thornton":                   8112,
	"Three Anchor Bay":           11022,
	"Tokai":                      9044,
	"Trovato":                    8687,
	"Turf Hall":                  16446,
	"Turf Hall Estate":           16455,
	"University Estate":          10161,
	"Valhalla Park":              8127,
	"Vanguard":                   8766,
	"Vredehoek":                  9166,
	"Vrygrond":                   16635,
	"Walmer Estate":              10163,
	"Waterfront":                 9169,
	"Welcome Estate":             8771,
	"Western Cape":               9,
	"Westlake":                   11740,
	"Wetton":                     8812,
	"Windermere":                 8019,
	"Woodstock":                  10164,
	"World View":                 16875,
	"Wynberg":                    10109,
	"Wynberg Upper":              10102,
	"Yorkshire Estate":           8814,
	"Zandvlei":                   9036,
	"Zeekoevlei":                 9047,
	"Zonnebloem":                 10166,
	"Zwaanswyk":                  9057,
	"Zwartdam":                   8733,
}
