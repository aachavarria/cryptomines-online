-- ============================================================================
-- Phase 1 MVP Migration - Cryptomines Online
-- Tables: players, planets, building_types, buildings, resources,
--         tech_types, technologies, + all building level reference tables
-- Seed data: building_types with exact GO2 wiki values, tech_types
-- ============================================================================

-- =========================
-- 1. PLAYERS
-- =========================
CREATE TABLE players (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    anonymous_id TEXT UNIQUE NOT NULL,
    username TEXT UNIQUE,
    level INTEGER NOT NULL DEFAULT 1,
    experience BIGINT NOT NULL DEFAULT 0,
    mall_points INTEGER NOT NULL DEFAULT 0,
    vouchers INTEGER NOT NULL DEFAULT 0,
    honor_points INTEGER NOT NULL DEFAULT 0,
    champion_points INTEGER NOT NULL DEFAULT 0,
    badges INTEGER NOT NULL DEFAULT 0,
    corsairs_gold INTEGER NOT NULL DEFAULT 0,
    tutorial_step INTEGER NOT NULL DEFAULT 0,
    is_online BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_login TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_players_anonymous_id ON players (anonymous_id);
CREATE INDEX idx_players_username ON players (username) WHERE username IS NOT NULL;
CREATE INDEX idx_players_last_login ON players (last_login);

-- =========================
-- 2. PLANETS
-- =========================
CREATE TABLE planets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    name TEXT NOT NULL DEFAULT 'New Colony',
    position_x INTEGER NOT NULL,
    position_y INTEGER NOT NULL,
    is_homeworld BOOLEAN NOT NULL DEFAULT false,
    is_rbp BOOLEAN NOT NULL DEFAULT false,
    rbp_level INTEGER NOT NULL DEFAULT 0,
    controlling_corp_id UUID,
    protection_until TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_planet_position UNIQUE (position_x, position_y)
);

CREATE INDEX idx_planets_player_id ON planets (player_id);
CREATE INDEX idx_planets_position ON planets (position_x, position_y);
CREATE INDEX idx_planets_rbp ON planets (is_rbp) WHERE is_rbp = true;

-- =========================
-- 3. BUILDING TYPES (Reference Data)
-- =========================
CREATE TABLE building_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('resource', 'core', 'military', 'defense', 'space', 'decorative')),
    base TEXT NOT NULL CHECK (base IN ('ground', 'space')) DEFAULT 'ground',
    base_cost_metal BIGINT NOT NULL DEFAULT 0,
    base_cost_he3 BIGINT NOT NULL DEFAULT 0,
    base_cost_gold BIGINT NOT NULL DEFAULT 0,
    base_time_seconds INTEGER NOT NULL DEFAULT 60,
    cost_multiplier NUMERIC(6,4) NOT NULL DEFAULT 3.0300,
    time_multiplier NUMERIC(6,4) NOT NULL DEFAULT 2.8700,
    base_production_per_hour INTEGER NOT NULL DEFAULT 0,
    production_multiplier NUMERIC(6,4) NOT NULL DEFAULT 1.1340,
    max_level INTEGER NOT NULL DEFAULT 24,
    max_count_per_planet INTEGER NOT NULL DEFAULT 1,
    prerequisite_building TEXT,
    prerequisite_level INTEGER NOT NULL DEFAULT 0,
    civic_center_req_per_level BOOLEAN NOT NULL DEFAULT true,
    description TEXT NOT NULL DEFAULT '',

    CONSTRAINT chk_positive_costs CHECK (
        base_cost_metal >= 0 AND base_cost_he3 >= 0 AND base_cost_gold >= 0
    ),
    CONSTRAINT chk_positive_multipliers CHECK (
        cost_multiplier > 0 AND time_multiplier > 0
    )
);

-- =========================
-- 4. BUILDINGS
-- =========================
CREATE TABLE buildings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planet_id UUID NOT NULL REFERENCES planets(id) ON DELETE CASCADE,
    building_type INTEGER NOT NULL REFERENCES building_types(id),
    grid_col INTEGER NOT NULL,
    grid_row INTEGER NOT NULL,
    level INTEGER NOT NULL DEFAULT 1,
    is_upgrading BOOLEAN NOT NULL DEFAULT false,
    upgrade_finish_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_level_positive CHECK (level >= 0),
    CONSTRAINT chk_grid_bounds CHECK (grid_col >= 0 AND grid_row >= 0),
    CONSTRAINT chk_upgrade_consistency CHECK (
        (is_upgrading = true AND upgrade_finish_at IS NOT NULL) OR
        (is_upgrading = false AND upgrade_finish_at IS NULL)
    )
);

CREATE INDEX idx_buildings_planet_id ON buildings (planet_id);
CREATE INDEX idx_buildings_type ON buildings (building_type);
CREATE INDEX idx_buildings_upgrading ON buildings (is_upgrading) WHERE is_upgrading = true;

-- =========================
-- 5. RESOURCES
-- =========================
CREATE TABLE resources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    planet_id UUID UNIQUE NOT NULL REFERENCES planets(id) ON DELETE CASCADE,
    metal BIGINT NOT NULL DEFAULT 0,
    he3 BIGINT NOT NULL DEFAULT 0,
    gold BIGINT NOT NULL DEFAULT 0,
    metal_per_hour BIGINT NOT NULL DEFAULT 0,
    he3_per_hour BIGINT NOT NULL DEFAULT 0,
    gold_per_hour BIGINT NOT NULL DEFAULT 0,
    storage_capacity BIGINT NOT NULL DEFAULT 100000,
    last_collected_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT chk_resources_non_negative CHECK (
        metal >= 0 AND he3 >= 0 AND gold >= 0
    ),
    CONSTRAINT chk_rates_non_negative CHECK (
        metal_per_hour >= 0 AND he3_per_hour >= 0 AND gold_per_hour >= 0
    )
);

-- =========================
-- 6. TECH TYPES (Reference Data)
-- =========================
CREATE TABLE tech_types (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL,
    tree TEXT NOT NULL CHECK (tree IN (
        'logistics_construction', 'planetary_defense', 'ballistics_science',
        'directional_science', 'missile_science', 'ship_based_science', 'ship_defense_science'
    )),
    prerequisites_json JSONB NOT NULL DEFAULT '[]',
    base_cost_metal BIGINT NOT NULL DEFAULT 0,
    base_cost_he3 BIGINT NOT NULL DEFAULT 0,
    base_cost_gold BIGINT NOT NULL DEFAULT 0,
    cost_multiplier NUMERIC(6,4) NOT NULL DEFAULT 1.5000,
    base_time_seconds INTEGER NOT NULL DEFAULT 120,
    time_multiplier NUMERIC(6,4) NOT NULL DEFAULT 1.5000,
    max_level INTEGER NOT NULL DEFAULT 10,
    effects_json JSONB NOT NULL DEFAULT '{}',
    description TEXT NOT NULL DEFAULT ''
);

CREATE INDEX idx_tech_types_tree ON tech_types (tree);

-- =========================
-- 7. TECHNOLOGIES
-- =========================
CREATE TABLE technologies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    player_id UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    tech_type INTEGER NOT NULL REFERENCES tech_types(id),
    level INTEGER NOT NULL DEFAULT 0,
    is_researching BOOLEAN NOT NULL DEFAULT false,
    research_finish_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT uq_player_tech UNIQUE (player_id, tech_type),
    CONSTRAINT chk_level_non_negative CHECK (level >= 0),
    CONSTRAINT chk_research_consistency CHECK (
        (is_researching = true AND research_finish_at IS NOT NULL) OR
        (is_researching = false AND research_finish_at IS NULL)
    )
);

CREATE INDEX idx_technologies_player_id ON technologies (player_id);
CREATE INDEX idx_technologies_researching ON technologies (is_researching) WHERE is_researching = true;

-- ============================================================================
-- BUILDING LEVEL REFERENCE TABLES
-- All data is exact from GO2 wiki
-- ============================================================================

-- =========================
-- 8. HE3 EXTRACTOR LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/He3_Extractor
-- =========================
CREATE TABLE he3_extractor_levels (
    level INTEGER PRIMARY KEY,
    he3_output_per_hour INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    civic_center_req INTEGER NOT NULL
);

INSERT INTO he3_extractor_levels VALUES
(1,  1180,   95,       80,       95,       40,      1),
(2,  1215,   163,      138,      163,      69,      1),
(3,  1264,   283,      238,      283,      119,     2),
(4,  1327,   492,      414,      492,      207,     2),
(5,  1407,   861,      725,      861,      362,     3),
(6,  1505,   1515,     1276,     1515,     638,     3),
(7,  1626,   2681,     2258,     2681,     1129,    4),
(8,  1772,   4773,     4019,     4773,     2010,    4),
(9,  1949,   8544,     7195,     8544,     3597,    5),
(10, 2164,   15379,    12950,    15379,    6475,    5),
(11, 2423,   27835,    23440,    27835,    11720,   6),
(12, 2738,   50660,    42661,    50660,    21331,   6),
(13, 3122,   92708,    78070,    92708,    39035,   7),
(14, 3590,   170583,   143649,   170583,   71824,   7),
(15, 4164,   315579,   265750,   315579,   132875,  8),
(16, 4872,   586976,   494296,   586976,   247148,  8),
(17, 5749,   1097645,  924333,   1097645,  462167,  9),
(18, 6842,   2063573,  1737746,  2063573,  868873,  9),
(19, 8210,   3900154,  3284340,  3900154,  1642170, 10),
(20, 9934,   7410292,  6240246,  7410292,  3120123, 10),
(21, 12120,  14153658, 11918870, 14153658, 5959435, 11),
(22, 14907,  27175024, 22884230, 27175024, 11442115, 11),
(23, 18485,  52447796, 44166565, 52447796, 22083283, 12),
(24, 23106,  101748723,85683135, 101748723,42841568, 12);

-- =========================
-- 9. METAL COLLECTOR LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Metal_Collector
-- =========================
CREATE TABLE metal_collector_levels (
    level INTEGER PRIMARY KEY,
    metal_output_per_hour INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    civic_center_req INTEGER NOT NULL
);

INSERT INTO metal_collector_levels VALUES
(1,  1080,   85,       106,      85,       40,      1),
(2,  1112,   146,      182,      146,      69,      1),
(3,  1157,   253,      315,      253,      119,     2),
(4,  1215,   440,      549,      440,      207,     2),
(5,  1288,   770,      960,      770,      362,     3),
(6,  1378,   1355,     1690,     1355,     638,     3),
(7,  1488,   2399,     2992,     2399,     1129,    4),
(8,  1622,   4271,     5326,     4271,     2010,    4),
(9,  1784,   7644,     9533,     7644,     3597,    5),
(10, 1980,   13760,    17159,    13760,    6475,    5),
(11, 2218,   24905,    31058,    24905,    11720,   6),
(12, 2506,   45328,    56526,    45328,    21331,   6),
(13, 2857,   82949,    103443,   82949,    39035,   7),
(14, 3286,   152627,   190335,   152627,   71824,   7),
(15, 3812,   282360,   352119,   282360,   132875,  8),
(16, 4459,   525189,   654942,   525189,   247148,  8),
(17, 5262,   982104,   1224741,  982104,   462167,  9),
(18, 6262,   1846355,  2302514,  1846355,  868873,  9),
(19, 7514,   3489611,  4351751,  3489611,  1642170, 10),
(20, 9092,   6630261,  8268326,  6630261,  3120123, 10),
(21, 11093,  12663799, 15792503, 12663799, 5959435, 11),
(22, 13644,  24314495, 30321605, 24314495, 11442115, 11),
(23, 16919,  46926975, 58520698, 46926975, 22083283, 12),
(24, 21148,  91038331, 113530154,91038331, 42841568, 12);

-- =========================
-- 10. RESIDENTIAL AREA LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Residential_Area
-- =========================
CREATE TABLE residential_area_levels (
    level INTEGER PRIMARY KEY,
    gold_output_per_hour INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    civic_center_req INTEGER NOT NULL
);

INSERT INTO residential_area_levels VALUES
(1,  1300,   78,       72,       65,       40,      1),
(2,  1339,   134,      124,      112,      69,      1),
(3,  1393,   232,      214,      193,      119,     2),
(4,  1462,   404,      373,      337,      207,     2),
(5,  1550,   707,      652,      589,      362,     3),
(6,  1658,   1244,     1148,     1037,     638,     3),
(7,  1791,   2202,     2032,     1835,     1129,    4),
(8,  1952,   3919,     3617,     3266,     2010,    4),
(9,  2148,   7015,     6475,     5846,     3597,    5),
(10, 2384,   12627,    11655,    10522,    6475,    5),
(11, 2670,   22854,    21096,    19045,    11720,   6),
(12, 3017,   41595,    38395,    34662,    21331,   6),
(13, 3439,   76118,    70263,    63432,    39035,   7),
(14, 3955,   140058,   129284,   116715,   71824,   7),
(15, 4588,   259107,   239175,   215922,   132875,  8),
(16, 5368,   481938,   444866,   401615,   247148,  8),
(17, 6334,   901225,   831900,   751021,   462167,  9),
(18, 7538,   1694302,  1563971,  1411919,  868873,  9),
(19, 9045,   3202232,  2955906,  2668526,  1642170, 10),
(20, 10945,  6084240,  5616221,  5070200,  3120123, 10),
(21, 13352,  11620898, 10726983, 9684082,  5959435, 11),
(22, 16423,  22312125, 20595807, 18593437, 11442115, 11),
(23, 20365,  43062401, 39749908, 35885334, 22083283, 12),
(24, 25456,  83541057, 77114822, 69617548, 42841568, 12);

-- =========================
-- 11. RESOURCE WAREHOUSE LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Resource_Warehouse
-- =========================
CREATE TABLE resource_warehouse_levels (
    level INTEGER PRIMARY KEY,
    storage_capacity BIGINT NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL,
    civic_center_req INTEGER NOT NULL
);

INSERT INTO resource_warehouse_levels VALUES
(1,  10000,     380,       370,       480,       35,      1),
(2,  30000,     654,       636,       826,       60,      1),
(3,  60000,     1131,      1101,      1428,      104,     2),
(4,  90000,     1967,      1916,      2485,      181,     2),
(5,  120000,    3443,      3352,      4349,      317,     3),
(6,  150000,    6060,      5900,      7654,      558,     3),
(7,  200000,    10726,     10444,     13548,     988,     4),
(8,  250000,    19092,     18590,     24116,     1758,    4),
(9,  300000,    34175,     33275,     43168,     3148,    5),
(10, 350000,    61514,     59896,     77702,     5666,    5),
(11, 400000,    111341,    108411,    140641,    10255,   6),
(12, 500000,    202641,    197308,    255967,    18664,   6),
(13, 600000,    370833,    361074,    468420,    34156,   7),
(14, 700000,    682332,    664376,    861893,    62846,   7),
(15, 800000,    1262314,   1229096,   1594502,   116266,  8),
(16, 1000000,   2347905,   2286118,   2965774,   216254,  8),
(17, 2000000,   4390582,   4275040,   5545998,   404396,  9),
(18, 3000000,   8254294,   8037075,   10426476,  760264,  9),
(19, 4000000,   15600615,  15190073,  19706040,  1436899, 10),
(20, 6000000,   29641169,  28861138,  37441476,  2730108, 10),
(21, 8000000,   56614633,  55124774,  71513220,  5214506, 11),
(22, 10000000,  108700094, 105839566, 137305382, 10011851, 11),
(23, 15000000,  209791182, 204270362, 264999388, 19322852, 12),
(24, 20000000,  406994893, 396284502, 514098813, 37486372, 12);

-- =========================
-- 12. CIVIC CENTER LEVELS
-- Source: GO2 wiki (exact data from GDD)
-- =========================
CREATE TABLE civic_center_levels (
    level INTEGER PRIMARY KEY,
    space_station_req INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO civic_center_levels VALUES
(1,  1,  550,       480,       600,       300),
(2,  1,  1661,      1450,      1812,      846),
(3,  2,  5033,      4392,      5490,      2394),
(4,  3,  15300,     13353,     16691,     6799),
(5,  4,  46664,     40725,     50907,     19378),
(6,  5,  142793,    124619,    155774,    55422),
(7,  6,  438375,    382582,    478227,    159063),
(8,  7,  1350194,   1178351,   1472939,   457500),
(9,  8,  4172100,   3641105,   4551382,   1323909),
(10, 9,  12933509,  11287426,  14109283,  3838937),
(11, 10, 40223214,  35103895,  43879869,  11172471),
(12, 11, 125496426, 109524154, 136905192, 32623616);

-- =========================
-- 13. TECHNOLOGY CENTER LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Technology_Center
-- =========================
CREATE TABLE technology_center_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    research_time_decrease_pct INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO technology_center_levels VALUES
(1,  1,  3,  450,       420,       650,       100),
(2,  2,  6,  1269,      1184,      1833,      302),
(3,  3,  9,  3591,      3352,      5187,      915),
(4,  4,  12, 10199,     9519,      14732,     2782),
(5,  5,  15, 29068,     27130,     41987,     8484),
(6,  6,  18, 83134,     77591,     120082,    25962),
(7,  7,  21, 238594,    222688,    344636,    79704),
(8,  8,  24, 687150,    641340,    992550,    245490),
(9,  9,  27, 1985864,   1853473,   2868470,   758564),
(10, 10, 30, 5759006,   5375072,   8318564,   2351527),
(11, 11, 33, 16758707,  15641460,  24207021,  7314512),
(12, 12, 36, 48935424,  45673062,  70684501,  22825932);

-- =========================
-- 14. COMMAND CENTER LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Command_Center
-- =========================
CREATE TABLE command_center_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    recruitment_cooldown_seconds INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO command_center_levels VALUES
(1,  1,  10800, 600,       450,       500,       40),
(2,  2,  10200, 1692,      1269,      1410,      121),
(3,  3,  9600,  4788,      3591,      3990,      366),
(4,  4,  9000,  13599,     10199,     11332,     1113),
(5,  5,  8400,  38757,     29068,     32297,     3394),
(6,  6,  7800,  110845,    83134,     92371,     10385),
(7,  7,  7200,  318125,    238594,    265104,    31882),
(8,  8,  6600,  916200,    687150,    763500,    98196),
(9,  9,  6000,  2647819,   1985864,   2206516,   303425),
(10, 10, 5400,  7678674,   5759006,   6398895,   940619),
(11, 11, 4800,  22344943,  16758707,  18620785,  2929525),
(12, 12, 4200,  65247232,  48935424,  54372693,  9126953);

-- =========================
-- 15. SPACE STATION LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Space_Station
-- =========================
CREATE TABLE space_station_levels (
    level INTEGER PRIMARY KEY,
    hp INTEGER NOT NULL,
    attack INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO space_station_levels VALUES
(1,  100,  0,      650,       600,       850,       200),
(2,  200,  0,      2093,      1932,      2737,      364),
(3,  300,  500,    6760,      6240,      8841,      1830),
(4,  400,  1500,   21904,     20219,     28643,     5564),
(5,  500,  4000,   71187,     65711,     93091,     16969),
(6,  600,  8000,   232069,    214218,    303475,    51925),
(7,  700,  16000,  758867,    700492,    992364,    159409),
(8,  800,  20000,  2489083,   2297615,   3254954,   491100),
(9,  900,  24000,  8189082,   7559153,   10708799,  1517127),
(10, 1000, 28000,  27023970,  24945203,  35339038,  4703094),
(11, 1100, 32000,  89449342,  82568623,  116972216, 14626623),
(12, 1200, 36000,  296971815, 274127829, 388347758, 45635064);

-- =========================
-- 16. WEAPON RESEARCH CENTER LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Weapon_Research_Center
-- =========================
CREATE TABLE weapon_research_center_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    weapon_research_time_decrease_pct INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO weapon_research_center_levels VALUES
(1,  1,  3,  500,       300,       450,       40),
(2,  2,  6,  1410,      846,       1269,      121),
(3,  3,  9,  3990,      2394,      3591,      366),
(4,  4,  12, 11332,     6799,      10199,     1113),
(5,  5,  15, 32297,     19378,     29068,     3394),
(6,  6,  18, 92371,     55422,     83134,     10385),
(7,  7,  21, 265104,    159063,    238594,    31882),
(8,  8,  24, 763500,    458100,    687150,    98196),
(9,  9,  27, 2206516,   1323909,   1985864,   303425),
(10, 10, 30, 6398895,   3839337,   5759006,   940619),
(11, 11, 33, 18620785,  11172471,  16758707,  2929525),
(12, 12, 36, 54372693,  32623616,  48935424,  9126953);

-- =========================
-- 17. ALLIANCE CENTER LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Alliance_Center
-- =========================
CREATE TABLE alliance_center_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    tech_center_req INTEGER NOT NULL,
    corps_support_speed_pct INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO alliance_center_levels VALUES
(1,  2,  1,  0,  500,       400,       550,       300),
(2,  3,  2,  10, 1410,      1128,      1551,      906),
(3,  4,  3,  20, 3990,      3192,      4389,      2745),
(4,  5,  4,  30, 11332,     9066,      12466,     8345),
(5,  6,  5,  40, 32297,     25838,     35527,     25453),
(6,  7,  6,  50, 92371,     73897,     101608,    77887),
(7,  8,  7,  60, 265104,    212083,    291615,    239113),
(8,  9,  8,  70, 763500,    610800,    839850,    736470),
(9,  10, 9,  80, 2206516,   1765213,   2427167,   2275691),
(10, 11, 10, 85, 6398895,   5119116,   7038785,   7055841),
(11, 12, 11, 90, 18620785,  14896628,  20482864,  21939935);

-- =========================
-- 18. TRADING CENTER LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Trading_Center
-- =========================
CREATE TABLE trading_center_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    weapon_research_req INTEGER NOT NULL,
    tech_center_req INTEGER NOT NULL,
    tax_rate_pct INTEGER NOT NULL,
    max_transactions INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO trading_center_levels VALUES
(1,  2,  1,  1,  10, 1, 1200,      1100,      1500,      400),
(2,  3,  2,  2,  10, 2, 4224,      3872,      5280,      1208),
(3,  4,  3,  3,  9,  3, 14911,     13668,     18638,     3660),
(4,  5,  4,  4,  8,  4, 52784,     48385,     65980,     11127),
(5,  6,  5,  5,  7,  5, 187383,    171768,    234229,    33938),
(6,  7,  6,  6,  6,  5, 667084,    611493,    833854,    103850),
(7,  8,  7,  7,  5,  5, 2381488,   2183031,   2976860,   318818),
(8,  9,  8,  8,  4,  5, 8525728,   7815251,   10657160,  982559),
(9,  10, 9,  9,  3,  5, 30607363,  28056750,  38259204,  3034254);

-- =========================
-- 19. RADAR LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Radar
-- =========================
CREATE TABLE radar_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    tech_center_req INTEGER NOT NULL,
    detect_advance_seconds INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO radar_levels VALUES
(1,  2,  1,  1800,  450,       400,       550,       60),
(2,  3,  2,  3600,  1269,      1128,      1551,      181),
(3,  4,  3,  5400,  3591,      3192,      4389,      549),
(4,  5,  4,  7200,  10199,     9066,      12466,     1699),
(5,  6,  5,  9000,  29068,     25838,     35527,     5091),
(6,  7,  6,  10800, 83134,     73897,     101608,    15577),
(7,  8,  7,  12600, 238594,    212083,    291615,    47823),
(8,  9,  8,  14400, 687150,    610800,    839850,    147294),
(9,  10, 9,  16200, 1985864,   1765213,   2427167,   455138);

-- =========================
-- 20. SPACEDOCK LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Spacedock
-- =========================
CREATE TABLE spacedock_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    ship_factory_req INTEGER NOT NULL,
    repair_pct INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO spacedock_levels VALUES
(1,  1,  1,  1,  900,       675,       750,       44),
(2,  2,  3,  2,  2538,      1903,      2115,      133),
(3,  3,  5,  3,  7182,      5386,      5985,      402),
(4,  4,  7,  4,  20398,     15298,     16998,     1224),
(5,  5,  9,  6,  58135,     43602,     48445,     3733),
(6,  6,  11, 8,  166267,    124701,    138556,    11423),
(7,  7,  13, 10, 477187,    357891,    397656,    35070),
(8,  8,  15, 12, 1374300,   1030725,   1145250,   108015),
(9,  9,  17, 14, 3971728,   2978796,   3309774,   333767),
(10, 10, 19, 16, 11518011,  8638509,   9598342,   1034680),
(11, 11, 21, 18, 33517414,  25138060,  27931177,  3218457),
(12, 12, 23, 20, 97870848,  73403136,  81559039,  10039714);

-- =========================
-- 21. RECYCLING PLANT LEVELS (partial data - only Lv1 and Lv11 have full costs on wiki)
-- Source: https://galaxyonlineii.fandom.com/wiki/Recycling_Plant
-- =========================
CREATE TABLE recycling_plant_levels (
    level INTEGER PRIMARY KEY,
    civic_center_req INTEGER NOT NULL,
    ship_factory_req INTEGER NOT NULL DEFAULT 0,
    recovery_pct INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL DEFAULT 0,
    he3_cost BIGINT NOT NULL DEFAULT 0,
    gold_cost BIGINT NOT NULL DEFAULT 0,
    build_time_seconds INTEGER NOT NULL DEFAULT 0
);

INSERT INTO recycling_plant_levels VALUES
(1,  2,  1,  5,  500,       400,       550,       200),
(2,  3,  3,  10, 0,         0,         0,         0),
(3,  4,  5,  15, 0,         0,         0,         0),
(4,  5,  7,  20, 0,         0,         0,         0),
(5,  6,  9,  25, 0,         0,         0,         0),
(6,  7,  11, 30, 0,         0,         0,         0),
(7,  8,  13, 35, 0,         0,         0,         0),
(8,  9,  15, 40, 0,         0,         0,         0),
(9,  10, 0,  45, 0,         0,         0,         0),
(10, 11, 0,  50, 0,         0,         0,         0),
(11, 12, 15, 55, 19941409,  15953127,  21935549,  14626623);

-- =========================
-- 22. METEOR STAR LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Meteor_Star
-- =========================
CREATE TABLE meteor_star_levels (
    level INTEGER PRIMARY KEY,
    space_station_req INTEGER NOT NULL,
    hp BIGINT NOT NULL,
    defense INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO meteor_star_levels VALUES
(1,  2,  200,       50,  50,        45,        48,        20),
(2,  2,  1000,      50,  141,       127,       135,       56),
(3,  3,  2000,      50,  399,       359,       383,       160),
(4,  3,  10000,     50,  1133,      1020,      1088,      453),
(5,  4,  20000,     50,  3230,      2907,      3101,      1292),
(6,  5,  40000,     50,  9237,      8313,      8868,      3695),
(7,  6,  100000,    50,  26510,     23859,     25450,     10604),
(8,  7,  200000,    50,  76350,     68715,     73296,     30540),
(9,  8,  400000,    50,  220652,    198586,    211826,    88261),
(10, 9,  30000000,  120, 639890,    575901,    614294,    255956),
(11, 10, 60000000,  120, 1862079,   1675871,   1787595,   744831),
(12, 11, 90000000,  120, 5437269,   4893542,   5219779,   2174908);

-- =========================
-- 23. PARTICLE CANNON LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Particle_Cannon
-- =========================
CREATE TABLE particle_cannon_levels (
    level INTEGER PRIMARY KEY,
    space_station_req INTEGER NOT NULL,
    attack INTEGER NOT NULL,
    hp BIGINT NOT NULL,
    range INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO particle_cannon_levels VALUES
(1,  2,  5000,    10000,     5,  450,       360,       520,       45),
(2,  2,  10000,   30000,     6,  1359,      1087,      1570,      136),
(3,  3,  15000,   100000,    7,  4118,      3294,      4758,      412),
(4,  3,  20000,   200000,    8,  12518,     10014,     14465,     1252),
(5,  4,  25000,   400000,    9,  38180,     30544,     44119,     3818),
(6,  5,  30000,   800000,    10, 116831,    93465,     135004,    11683),
(7,  6,  40000,   1000000,   11, 358670,    286936,    414463,    35867),
(8,  7,  50000,   1200000,   12, 1104704,   883763,    1276547,   110470),
(9,  8,  70000,   1400000,   13, 3413536,   2730829,   3944531,   341354),
(10, 9,  1000000, 5000000,   14, 10581962,  8465570,   12228045,  1058196),
(11, 10, 1500000, 10000000,  15, 32909902,  26327922,  38029220,  3290990),
(12, 11, 2000000, 20000000,  16, 102678894, 82143115,  118651167, 10267889);

-- =========================
-- 24. ANTI-AIRCRAFT GUN LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Anti-aircraft_Gun
-- =========================
CREATE TABLE anti_aircraft_gun_levels (
    level INTEGER PRIMARY KEY,
    space_station_req INTEGER NOT NULL,
    attack INTEGER NOT NULL,
    hp BIGINT NOT NULL,
    range INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO anti_aircraft_gun_levels VALUES
(1,  4,  25000,   20000,     3,  650,       600,       850,       45),
(2,  4,  50000,   40000,     4,  1963,      1812,      2567,      136),
(3,  5,  75000,   200000,    5,  5948,      5490,      7778,      412),
(4,  5,  100000,  400000,    6,  18082,     16691,     23645,     1252),
(5,  6,  125000,  1000000,   7,  55149,     50907,     72118,     3818),
(6,  6,  150000,  2000000,   8,  168755,    155774,    220680,    11683),
(7,  7,  175000,  2400000,   9,  518079,    478227,    677488,    35867),
(8,  8,  200000,  2800000,   10, 1595684,   1472939,   2086664,   110470),
(9,  9,  225000,  3200000,   11, 4930663,   4551382,   6447790,   341354),
(10, 10, 2000000, 10000000,  12, 15285056,  14109283,  19988151,  1058196),
(11, 11, 2500000, 30000000,  13, 47536525,  43879869,  58506492,  3290990),
(12, 12, 3000000, 50000000,  14, 148313958, 136905192, 182540256, 10267889);

-- =========================
-- 25. THOR'S CANNON LEVELS
-- Source: https://galaxyonlineii.fandom.com/wiki/Thor%27s_Cannon
-- =========================
CREATE TABLE thors_cannon_levels (
    level INTEGER PRIMARY KEY,
    space_station_req INTEGER NOT NULL,
    attack INTEGER NOT NULL,
    cooldown INTEGER NOT NULL,
    hp BIGINT NOT NULL,
    range INTEGER NOT NULL,
    metal_cost BIGINT NOT NULL,
    he3_cost BIGINT NOT NULL,
    gold_cost BIGINT NOT NULL,
    build_time_seconds INTEGER NOT NULL
);

INSERT INTO thors_cannon_levels VALUES
(1,  6,  10000,    2, 200000,    10, 80000,     90000,     125000,    40000),
(2,  6,  50000,    2, 400000,    11, 161600,    181800,    252500,    60800),
(3,  7,  100000,   2, 1000000,   12, 328048,    369054,    512575,    93024),
(4,  7,  150000,   2, 2000000,   13, 669218,    752870,    1045653,   143257),
(5,  8,  200000,   2, 3000000,   14, 1371897,   1543384,   2143589,   222048),
(6,  8,  300000,   2, 4000000,   15, 2826107,   3179371,   4415793,   346395),
(7,  9,  400000,   2, 5000000,   16, 5850042,   6581297,   9140691,   543841),
(8,  9,  500000,   2, 6000000,   17, 12168087,  13689098,  19012637,  589268),
(9,  10, 800000,   2, 7000000,   18, 25431303,  28610216,  39736411,  1366237),
(10, 10, 3000000,  1, 30000000,  21, 53405736,  60081453,  83446462,  2185978),
(11, 11, 4500000,  1, 60000000,  22, 112686103, 126771866, 176072036, 3519425),
(12, 12, 6000000,  1, 90000000,  23, 238894538, 268756355, 373272716, 5701489);

-- ============================================================================
-- SEED DATA: BUILDING TYPES
-- Uses exact Lv1 base costs from GO2 wiki lookup tables
-- ============================================================================
INSERT INTO building_types (name, display_name, category, base, base_cost_metal, base_cost_he3, base_cost_gold, base_time_seconds, cost_multiplier, time_multiplier, base_production_per_hour, production_multiplier, max_level, max_count_per_planet, prerequisite_building, prerequisite_level, description) VALUES
-- Ground Base: Resource Buildings (exact Lv1 costs from wiki)
('metal_collector',       'Metal Collector',       'resource',  'ground', 85,   106,  85,   40,  1.7500, 1.7500, 1080, 1.1340, 24, 8,  'civic_center', 1,  'Produces Metal'),
('he3_extractor',         'He3 Extractor',         'resource',  'ground', 95,   80,   95,   40,  1.7500, 1.7500, 1180, 1.1340, 24, 8,  'civic_center', 1,  'Produces He3 (Helium-3)'),
('residential_area',      'Residential Area',      'resource',  'ground', 78,   72,   65,   40,  1.7500, 1.7500, 1300, 1.1340, 24, 8,  'civic_center', 1,  'Produces Gold (highest output)'),
('resource_warehouse',    'Resource Warehouse',    'resource',  'ground', 380,  370,  480,  35,  1.7500, 1.7500, 0,    1.0000, 24, 1,  'civic_center', 1,  'Stores all resources; increases capacity'),
-- Ground Base: Core / Administrative Buildings
('civic_center',          'Civic Center',          'core',      'ground', 550,  480,  600,  300, 3.0300, 2.8700, 0,    1.0000, 12, 1,  NULL,           0,  'Main hub; determines max level of all other buildings'),
('technology_center',     'Technology Center',     'core',      'ground', 450,  420,  650,  100, 3.0300, 2.8700, 0,    1.0000, 12, 1,  'civic_center', 1,  'Research facility (7 science trees); 3% research time reduction per level'),
('alliance_center',       'Alliance Center',       'core',      'ground', 500,  400,  550,  300, 3.0300, 2.8700, 0,    1.0000, 11, 1,  'civic_center', 2,  'Enables Corp membership and features'),
('trading_center',        'Trading Center',        'core',      'ground', 1200, 1100, 1500, 400, 3.0300, 2.8700, 0,    1.0000, 9,  1,  'civic_center', 2,  'Player-to-player trading and auctions'),
('galaxy_transporter',    'Galaxy Transporter',    'core',      'ground', 350,  300,  450,  100, 3.0300, 2.8700, 0,    1.0000, 12, 1,  'civic_center', 1,  'Inter-system resource transport'),
('compound_center',       'Compound Center',       'core',      'ground', 400,  350,  500,  100, 3.0300, 2.8700, 0,    1.0000, 12, 1,  'civic_center', 3,  'Commander card merging and enhancement'),
('radar',                 'Radar',                 'core',      'ground', 450,  400,  550,  60,  3.0300, 2.8700, 0,    1.0000, 9,  1,  'civic_center', 2,  'Detection of incoming attacks'),
-- Ground Base: Military Buildings
('ship_factory',          'Ship Factory',          'military',  'ground', 600,  450,  500,  200, 3.0300, 2.8700, 0,    1.0000, 24, 1,  'civic_center', 1,  'Constructs ships; holds 20 designs, 5 production slots'),
('spacedock',             'Spacedock',             'military',  'ground', 900,  675,  750,  44,  3.0300, 2.8700, 0,    1.0000, 12, 1,  'civic_center', 1,  'Ship berthing and fleet management'),
('command_center',        'Command Center',        'military',  'ground', 600,  450,  500,  40,  3.0300, 2.8700, 0,    1.0000, 12, 1,  'civic_center', 1,  'Commander recruitment (60 max at Lv71+)'),
('weapon_research_center','Weapon Research Center','military',  'ground', 500,  300,  450,  40,  3.0300, 2.8700, 0,    1.0000, 12, 1,  'civic_center', 1,  'Develops weapons and modules; 3% weapon research time reduction per level'),
('recycling_plant',       'Recycling Plant',       'military',  'ground', 500,  400,  550,  200, 3.0300, 2.8700, 0,    1.0000, 11, 1,  'civic_center', 2,  'Recovers resources from scrapped ships'),
-- Space Base Buildings
('space_station',         'Space Station',         'space',     'space',  650,  600,  850,  200, 3.0300, 2.8700, 0,    1.0000, 12, 1,  NULL,           0,  'Main orbital structure; must stay within 1 level of Civic Center'),
('meteor_star',           'Meteor Star',           'defense',   'space',  50,   45,   48,   20,  3.0300, 2.8700, 0,    1.0000, 12, 40, 'space_station', 2,  'Orbital defense structure'),
('particle_cannon',       'Particle Cannon',       'defense',   'space',  450,  360,  520,  45,  3.0300, 2.8700, 0,    1.0000, 12, 20, 'space_station', 2,  'Energy weapon defense; no cooldown, fires every round'),
('anti_aircraft_gun',     'Anti-Aircraft Gun',     'defense',   'space',  650,  600,  850,  45,  3.0300, 2.8700, 0,    1.0000, 12, 15, 'space_station', 4,  'Anti-air defense; AoE damage'),
('thors_cannon',          'Thor''s Cannon',        'defense',   'space',  80000,90000,125000,40000,3.0300,2.8700, 0,    1.0000, 12, 3,  'space_station', 6,  'Advanced heavy defense cannon; cooldown drops from 2 to 1 at Lv10'),
('celestial_base',        'Celestial Base',        'space',     'space',  600,  500,  600,  200, 3.0300, 2.8700, 0,    1.0000, 12, 1,  'space_station', 6,  'Advanced orbital facility');

-- ============================================================================
-- SEED DATA: TECH TYPES (Logistics Construction + Ship Defense Science for MVP)
-- Exact costs and times from GO2 wiki research
-- ============================================================================

-- Logistics Construction Science Tree
INSERT INTO tech_types (name, display_name, tree, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, max_level, effects_json, description) VALUES
('concurrent_construction', 'Concurrent Construction', 'logistics_construction',
    '[]',
    0, 0, 1000, 1.0000, 20, 1.0000, 1,
    '{"construction_slots": 1}',
    'Adds 1 additional construction slot'),
('construction_boost', 'Construction Boost', 'logistics_construction',
    '[{"tech": "concurrent_construction", "level": 1}]',
    0, 0, 2400, 1.9600, 480, 1.9600, 10,
    '{"construction_speed_pct": [1, 2, 3, 5, 6, 8, 9, 11, 13, 15]}',
    'Increases building construction speed'),
('quality_materials', 'Quality Materials', 'logistics_construction',
    '[{"tech": "construction_boost", "level": 3}]',
    0, 0, 1200, 1.9600, 240, 1.9600, 10,
    '{"building_cost_reduction_pct": [1, 2, 3, 5, 6, 8, 9, 11, 13, 15]}',
    'Reduces building resource costs'),
('ship_building_boost', 'Ship Building Boost', 'logistics_construction',
    '[]',
    0, 0, 990, 2.2500, 198, 2.2500, 10,
    '{"ship_building_speed_pct": [1, 2, 3, 5, 6, 8, 9, 11, 13, 15]}',
    'Increases shipbuilding speed'),
('ship_building_logistics', 'Ship Building Logistics', 'logistics_construction',
    '[{"tech": "ship_building_boost", "level": 2}]',
    0, 0, 1386, 2.2500, 277, 2.2500, 10,
    '{"ship_cost_reduction_pct": [1, 2, 3, 5, 6, 8, 9, 11, 13, 15]}',
    'Reduces ship resource costs'),
('sync_shipbuilding', 'Sync Shipbuilding', 'logistics_construction',
    '[{"tech": "ship_building_logistics", "level": 4}]',
    0, 0, 174000, 1.0000, 34800, 1.0000, 1,
    '{"ship_building_slots": 1}',
    'Adds 1 additional shipbuilding slot (5th slot)'),
('repair_technology', 'Repair Technology', 'logistics_construction',
    '[{"tech": "sync_shipbuilding", "level": 1}]',
    0, 0, 5310, 1.9600, 1062, 1.9600, 10,
    '{"repair_pct": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}',
    'Increases percentage of ships that can be repaired after battle'),
('high_yield_mining', 'High Yield Mining', 'logistics_construction',
    '[]',
    0, 0, 1740, 1.9600, 348, 1.9600, 10,
    '{"metal_production_pct": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}',
    'Increases Metal output'),
('high_yield_chemistry', 'High Yield Chemistry', 'logistics_construction',
    '[{"tech": "high_yield_mining", "level": 2}]',
    0, 0, 2400, 1.9600, 480, 1.9600, 10,
    '{"he3_production_pct": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}',
    'Increases He3 output'),
('high_yield_investing', 'High Yield Investing', 'logistics_construction',
    '[{"tech": "high_yield_chemistry", "level": 2}]',
    0, 0, 3570, 1.9600, 714, 1.9600, 10,
    '{"gold_production_pct": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}',
    'Increases Gold output'),
('expand_capacity', 'Expand Capacity', 'logistics_construction',
    '[{"tech": "high_yield_investing", "level": 4}]',
    0, 0, 3540, 1.9600, 708, 1.9600, 10,
    '{"warehouse_capacity_pct": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]}',
    'Increases warehouse storage capacity');

-- Ship Defense Science Tree
INSERT INTO tech_types (name, display_name, tree, prerequisites_json, base_cost_metal, base_cost_he3, base_cost_gold, cost_multiplier, base_time_seconds, time_multiplier, max_level, effects_json, description) VALUES
('ship_defense_tech', 'Ship Defense Tech', 'ship_defense_science',
    '[]',
    0, 0, 2846, 1.1580, 3240, 1.3400, 2,
    '{"base_shield_pct": [2, 5], "base_structure_pct": [2, 5], "base_agility_pct": [2, 5], "base_defense_pct": [2, 5], "stability_pct": [5, 10]}',
    'Base ship defense technology'),
('shield_research', 'Shield Research', 'ship_defense_science',
    '[{"tech": "ship_defense_tech", "level": 1}]',
    0, 0, 10392, 1.4200, 5400, 1.4200, 5,
    '{"base_shield_pct": [1, 2, 3, 4, 5]}',
    'Increases base shield'),
('energy_diffusion', 'Energy Diffusion', 'ship_defense_science',
    '[{"tech": "ship_defense_tech", "level": 2}, {"tech": "shield_research", "level": 3}]',
    0, 0, 46765, 1.4200, 10800, 1.4200, 3,
    '{"shield_damage_reduction": [1, 2, 3]}',
    'Each shield module reduces damage'),
('penetration_resistance', 'Penetration Resistance', 'ship_defense_science',
    '[{"tech": "shield_research", "level": 5}, {"tech": "energy_diffusion", "level": 1}]',
    0, 0, 80000, 1.4200, 14400, 1.4200, 2,
    '{"enemy_shield_pen_reduction_pct": [3, 7]}',
    'Reduces enemy shield penetration chance'),
('augment_shield', 'Augment Shield', 'ship_defense_science',
    '[{"tech": "energy_diffusion", "level": 2}, {"tech": "penetration_resistance", "level": 1}]',
    0, 0, 120000, 1.4200, 21600, 1.4200, 3,
    '{"base_shield_pct": [6, 12, 20]}',
    'Increases base shield significantly'),
('restoration', 'Restoration', 'ship_defense_science',
    '[{"tech": "energy_diffusion", "level": 3}, {"tech": "augment_shield", "level": 2}]',
    0, 0, 200000, 1.4200, 36000, 1.4200, 2,
    '{"shield_restore_pct": [30, 60], "interception_pct": [1, 2]}',
    'Restores shields per round and increases interception'),
('augment_absorption', 'Augment Absorption', 'ship_defense_science',
    '[{"tech": "shield_research", "level": 5}, {"tech": "energy_diffusion", "level": 1}]',
    0, 0, 80000, 1.4200, 14400, 1.4200, 2,
    '{"shield_damage_reduction": [2, 5]}',
    'Shield modules reduce more damage'),
('energy_conservation', 'Energy Conservation', 'ship_defense_science',
    '[{"tech": "energy_diffusion", "level": 2}, {"tech": "augment_absorption", "level": 1}]',
    0, 0, 120000, 1.4200, 21600, 1.4200, 3,
    '{"absorb_no_he3_pct": [3, 6, 10]}',
    'Chance to absorb damage without He3 cost'),
('electronic_barrier', 'Electronic Barrier', 'ship_defense_science',
    '[{"tech": "energy_diffusion", "level": 3}, {"tech": "energy_conservation", "level": 2}]',
    0, 0, 200000, 1.4200, 36000, 1.4200, 2,
    '{"reflect_damage_pct": [5, 10]}',
    'Reflects damage before shields drop to 0'),
('damage_mitigation', 'Damage Mitigation', 'ship_defense_science',
    '[{"tech": "augment_shield", "level": 3}, {"tech": "restoration", "level": 2}, {"tech": "energy_conservation", "level": 3}, {"tech": "electronic_barrier", "level": 2}]',
    0, 0, 500000, 1.4200, 72000, 1.4200, 3,
    '{"absorb_double_pct": [10, 20, 30], "lower_collateral_pct": [15, 30, 45]}',
    'Absorbs double damage and lowers collateral'),
('ship_structural_analysis', 'Ship Structural Analysis', 'ship_defense_science',
    '[{"tech": "ship_defense_tech", "level": 1}]',
    0, 0, 10392, 1.4200, 5400, 1.4200, 2,
    '{"base_structure_pct": [1, 2]}',
    'Increases base structure'),
('ship_reinforcement', 'Ship Reinforcement', 'ship_defense_science',
    '[{"tech": "ship_defense_tech", "level": 2}, {"tech": "ship_structural_analysis", "level": 2}]',
    0, 0, 46765, 1.4200, 10800, 1.4200, 3,
    '{"structure_damage_reduction": [1, 2, 3]}',
    'Each structure module reduces damage'),
('resilience', 'Resilience', 'ship_defense_science',
    '[{"tech": "ship_structural_analysis", "level": 2}, {"tech": "ship_reinforcement", "level": 1}]',
    0, 0, 80000, 1.4200, 14400, 1.4200, 2,
    '{"enemy_structure_pen_reduction_pct": [3, 7]}',
    'Reduces enemy structure penetration chance'),
('structure_improvement', 'Structure Improvement', 'ship_defense_science',
    '[{"tech": "ship_reinforcement", "level": 2}, {"tech": "resilience", "level": 1}]',
    0, 0, 120000, 1.4200, 21600, 1.4200, 3,
    '{"base_structure_pct": [6, 12, 20]}',
    'Increases base structure significantly'),
('fast_repair', 'Fast Repair', 'ship_defense_science',
    '[{"tech": "ship_reinforcement", "level": 3}, {"tech": "structure_improvement", "level": 2}]',
    0, 0, 200000, 1.4200, 36000, 1.4200, 2,
    '{"structure_restore_pct": [30, 60]}',
    'Restores structure per round'),
('reaction_armor_improvement', 'Reaction Armor Improvement', 'ship_defense_science',
    '[{"tech": "ship_structural_analysis", "level": 2}, {"tech": "ship_reinforcement", "level": 1}]',
    0, 0, 80000, 1.4200, 14400, 1.4200, 2,
    '{"structure_damage_reduction": [2, 5]}',
    'Structure modules reduce more damage'),
('defense_improvement', 'Defense Improvement', 'ship_defense_science',
    '[{"tech": "ship_reinforcement", "level": 2}, {"tech": "reaction_armor_improvement", "level": 1}]',
    0, 0, 120000, 1.4200, 21600, 1.4200, 3,
    '{"absorb_no_he3_pct": [3, 6, 10]}',
    'Chance to absorb damage without He3 cost (structure)'),
('reflection_mastery', 'Reflection Mastery', 'ship_defense_science',
    '[{"tech": "ship_reinforcement", "level": 3}, {"tech": "defense_improvement", "level": 2}]',
    0, 0, 200000, 1.4200, 36000, 1.4200, 2,
    '{"reflect_damage_pct": [5, 10]}',
    'Reflects damage before structure drops to 0'),
('stability_mastery', 'Stability Mastery', 'ship_defense_science',
    '[{"tech": "structure_improvement", "level": 3}, {"tech": "fast_repair", "level": 2}, {"tech": "defense_improvement", "level": 3}, {"tech": "reflection_mastery", "level": 2}]',
    0, 0, 500000, 1.4200, 72000, 1.4200, 3,
    '{"absorb_double_pct": [10, 20, 30], "lower_collateral_pct": [15, 30, 45]}',
    'Absorbs double damage and lowers collateral (structure)');

-- ============================================================================
-- TRIGGERS: Auto-update updated_at
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER tr_players_updated_at
    BEFORE UPDATE ON players
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_planets_updated_at
    BEFORE UPDATE ON planets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_buildings_updated_at
    BEFORE UPDATE ON buildings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_resources_updated_at
    BEFORE UPDATE ON resources
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER tr_technologies_updated_at
    BEFORE UPDATE ON technologies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- ROW LEVEL SECURITY (RLS) POLICIES
-- ============================================================================
ALTER TABLE players ENABLE ROW LEVEL SECURITY;
ALTER TABLE planets ENABLE ROW LEVEL SECURITY;
ALTER TABLE buildings ENABLE ROW LEVEL SECURITY;
ALTER TABLE resources ENABLE ROW LEVEL SECURITY;
ALTER TABLE technologies ENABLE ROW LEVEL SECURITY;

-- Players can read their own data
CREATE POLICY players_select_own ON players
    FOR SELECT USING (auth.uid() = id);

-- Players can read their own planets
CREATE POLICY planets_select_own ON planets
    FOR SELECT USING (player_id = auth.uid());

-- Players can read buildings on their own planets
CREATE POLICY buildings_select_own ON buildings
    FOR SELECT USING (
        planet_id IN (SELECT id FROM planets WHERE player_id = auth.uid())
    );

-- Players can read their own resources
CREATE POLICY resources_select_own ON resources
    FOR SELECT USING (
        planet_id IN (SELECT id FROM planets WHERE player_id = auth.uid())
    );

-- Players can read their own technologies
CREATE POLICY technologies_select_own ON technologies
    FOR SELECT USING (player_id = auth.uid());

-- Reference tables are readable by all authenticated users
ALTER TABLE building_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE tech_types ENABLE ROW LEVEL SECURITY;
ALTER TABLE he3_extractor_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE metal_collector_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE residential_area_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE resource_warehouse_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE civic_center_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE technology_center_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE command_center_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE space_station_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE weapon_research_center_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE alliance_center_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE trading_center_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE radar_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE spacedock_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE recycling_plant_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE meteor_star_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE particle_cannon_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE anti_aircraft_gun_levels ENABLE ROW LEVEL SECURITY;
ALTER TABLE thors_cannon_levels ENABLE ROW LEVEL SECURITY;

CREATE POLICY building_types_select_all ON building_types FOR SELECT USING (true);
CREATE POLICY tech_types_select_all ON tech_types FOR SELECT USING (true);
CREATE POLICY he3_extractor_levels_select_all ON he3_extractor_levels FOR SELECT USING (true);
CREATE POLICY metal_collector_levels_select_all ON metal_collector_levels FOR SELECT USING (true);
CREATE POLICY residential_area_levels_select_all ON residential_area_levels FOR SELECT USING (true);
CREATE POLICY resource_warehouse_levels_select_all ON resource_warehouse_levels FOR SELECT USING (true);
CREATE POLICY civic_center_levels_select_all ON civic_center_levels FOR SELECT USING (true);
CREATE POLICY technology_center_levels_select_all ON technology_center_levels FOR SELECT USING (true);
CREATE POLICY command_center_levels_select_all ON command_center_levels FOR SELECT USING (true);
CREATE POLICY space_station_levels_select_all ON space_station_levels FOR SELECT USING (true);
CREATE POLICY weapon_research_center_levels_select_all ON weapon_research_center_levels FOR SELECT USING (true);
CREATE POLICY alliance_center_levels_select_all ON alliance_center_levels FOR SELECT USING (true);
CREATE POLICY trading_center_levels_select_all ON trading_center_levels FOR SELECT USING (true);
CREATE POLICY radar_levels_select_all ON radar_levels FOR SELECT USING (true);
CREATE POLICY spacedock_levels_select_all ON spacedock_levels FOR SELECT USING (true);
CREATE POLICY recycling_plant_levels_select_all ON recycling_plant_levels FOR SELECT USING (true);
CREATE POLICY meteor_star_levels_select_all ON meteor_star_levels FOR SELECT USING (true);
CREATE POLICY particle_cannon_levels_select_all ON particle_cannon_levels FOR SELECT USING (true);
CREATE POLICY anti_aircraft_gun_levels_select_all ON anti_aircraft_gun_levels FOR SELECT USING (true);
CREATE POLICY thors_cannon_levels_select_all ON thors_cannon_levels FOR SELECT USING (true);
