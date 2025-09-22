import csv
from collections import defaultdict

INPUT_FILE = "data/stops.txt"
OUTPUT_FILE = "../config/stations.go"

stations = defaultdict(list)
with open(INPUT_FILE, newline='', encoding='utf-8') as csvfile:
    reader = csv.DictReader(csvfile)
    for row in reader:
        stop_id = row["stop_id"]
        stop_name = row["stop_name"]
        stations[stop_name].append(stop_id)

with open(OUTPUT_FILE, "w", encoding="utf-8") as out:
    out.write("package config\n\n")
    out.write("// StationStops maps station name -> stop IDs (N/S directions etc)\n")
    out.write("var StationStops = map[string][]string{\n")

    for name, ids in sorted(stations.items()):  
        sorted_ids = sorted(ids)  
        ids_list = ", ".join([f"\"{i}\"" for i in sorted_ids])
        out.write(f'\t"{name}": {{{ids_list}}},\n')

    out.write("}\n")

print(f"Generated {OUTPUT_FILE} with {len(stations)} stations.")
